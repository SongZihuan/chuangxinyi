package accessrecord

import (
	"bufio"
	"bytes"
	"context"
	"database/sql"
	"fmt"
	"gitee.com/wuntsong-auth/backend/src/config"
	"gitee.com/wuntsong-auth/backend/src/global/peername"
	"gitee.com/wuntsong-auth/backend/src/ip"
	"gitee.com/wuntsong-auth/backend/src/logger"
	"gitee.com/wuntsong-auth/backend/src/model/db"
	"gitee.com/wuntsong-auth/backend/src/mysql"
	"gitee.com/wuntsong-auth/backend/src/service/v1/gen/types"
	"gitee.com/wuntsong-auth/backend/src/service/v1/src/record"
	"gitee.com/wuntsong-auth/backend/src/service/v1/src/respmsg"
	"gitee.com/wuntsong-auth/backend/src/utils"
	"github.com/gorilla/websocket"
	"github.com/wuntsong-org/go-zero-plus/rest/httpx"
	errors "github.com/wuntsong-org/wterrors"
	"io"
	"net"
	"net/http"
	"strings"
	"time"
	"unicode/utf8"
)

func AccessRecordHandle(w http.ResponseWriter, r *http.Request, next http.HandlerFunc, notRecordPath []string) {
	ctx := r.Context()
	notRecord := utils.StringIn(r.URL.Path, notRecordPath)
	contentType := r.Header.Get("Content-Type")
	method := r.Method
	path := r.URL.Path
	queryValues := r.URL.Query()
	host := r.Host
	scheme := r.Header.Get("X-Forwarded-Proto")
	if len(scheme) == 0 {
		scheme = "http"
	}

	if websocket.IsWebSocketUpgrade(r) {
		scheme = fmt.Sprintf("%s/ws", scheme)
	}

	if len(contentType) > 80 {
		contentType = contentType[:80]
	}

	if len(method) > 15 {
		method = method[:15]
	}

	if len(scheme) > 15 {
		scheme = scheme[:15]
	}

	if len(path) > 100 {
		path = path[:100]
	}

	queryMap := make(map[string][]string, len(queryValues))
QUERY:
	for k, v := range queryValues {
		if len(v) == 0 {
			continue
		}

		if len(v) > 5 {
			continue
		}

		for _, vv := range v {
			if len(vv) > 100 {
				continue QUERY
			}
		}

		queryMap[k] = v
	}
	queryByte, jsonErr := utils.JsonMarshal(queryMap)
	if jsonErr != nil {
		queryByte = []byte("{}")
	}
	query := string(queryByte)
	if len(query) > 6000 {
		query = "{}"
	}

	if len(host) > 500 {
		host = host[:500]
	}

	headerMap := make(map[string]string, len(config.BackendConfig.User.AllowHeader))
	for _, h := range config.BackendConfig.User.AllowHeader {
		v := r.Header.Get(h)
		if len(v) == 0 {
			continue
		}

		if len(v) > 100 {
			continue
		}

		headerMap[h] = v
	}
	headerByte, jsonErr := utils.JsonMarshal(headerMap)
	if jsonErr != nil {
		headerByte = []byte("{}")
	}
	header := string(headerByte)
	if len(header) > 6000 {
		header = "{}"
	}

	bodyByte, err := io.ReadAll(r.Body)
	if err != nil {
		accessRecordError(w, r, errors.WarpQuick(err), "获取请求体错误")
		return
	}
	r.Body = io.NopCloser(bytes.NewBuffer(bodyByte)) // 塞回去

	var body string
	err = r.ParseMultipartForm(10 << 20) // 限制最大10MB大小的表单数据
	if err == nil {
		type File struct {
			FileName string `json:"fileName"`
			FileSize int64  `json:"fileSize"`
		}

		data := struct {
			Value map[string][]string `json:"value"`
			File  map[string][]File
		}{Value: r.MultipartForm.Value, File: make(map[string][]File, len(r.MultipartForm.File))}

		for n, k := range r.MultipartForm.File {
			FileList := make([]File, 0, len(k))
			for _, h := range k {
				FileList = append(FileList, File{
					FileName: h.Filename,
					FileSize: h.Size,
				})
			}
			data.File[n] = FileList
		}

		bodyByte, err := utils.JsonMarshal(data)
		if err != nil {
			accessRecordError(w, r, err, "打包form请求体错误")
			return
		}

		body = string(bodyByte)
	} else {
		var data map[string]interface{}
		err := utils.JsonUnmarshal(bodyByte, &data)
		if err == nil {
			for k, v := range data {
				vString, ok := v.(string)
				if !ok {
					continue
				}

				if strings.HasPrefix(vString, "base64:") {
					data[k] = ""
				}
			}

			newBody, err := utils.JsonMarshal(data)
			if err != nil {
				newBody = bodyByte
			}

			body = string(newBody)
			if len(body) > 65000 {
				body = body[:65000]
			}
		} else {
			if utf8.Valid(bodyByte) {
				body = string(bodyByte)
				if len(body) > 65000 {
					body = body[:65000]
				}
			} else {
				body = fmt.Sprintf("<media type: %s bytes: %d>", utils.GetMediaType(bodyByte), len(bodyByte))
			}
		}
	}

	realIp := utils.GetTargetIP(r)
	var code, geo string
	code, geo, err = ip.GetGeo(ctx, realIp)
	if err != nil {
		code = ip.UnknownGeoCode
		geo = "未知"
	}

	var requestIDPrefix string
	err = func() (resErr errors.WTError) {
		// 请求前缀可以重复，因为请求ID还有唯一的后缀数字
		defer utils.Recover(logger.Logger, &resErr, "create requests id prefix error")
		unixNano := time.Now().UnixNano()
		text := fmt.Sprintf("%s\n%d\n%s\n%s\n", realIp, unixNano, r.URL.String(), peername.PeerName)
		requestIDPrefix = fmt.Sprintf("%s-%d", utils.HashSHA256WithBase62(text), unixNano)
		return nil
	}()
	if err != nil {
		accessRecordError(w, r, errors.WarpQuick(err), "创建请求前缀错误")
		return
	}

	access := &db.AccessRecord{
		RequestIdPrefix: requestIDPrefix,
		ServerName:      peername.PeerName,
		Ip:              realIp,
		GeoCode:         code,
		Geo:             geo,
		Method:          method,
		Path:            path,
		Scheme:          scheme,
		Host:            host,
		Query:           query,
		ContentType:     contentType,
		RequestsBody:    body,
		RequestsHeader:  header,
	}

	accessModel := db.NewAccessRecordModel(mysql.MySQLConn)

	if !notRecord {
		res, err := accessModel.Insert(context.Background(), access)
		if err != nil {
			logger.Logger.Error("mysql error: %s", err.Error())
			return
		}

		access.Id, err = res.LastInsertId()
		if err != nil {
			logger.Logger.Error("mysql error: %s", err.Error())
			return
		}
	}

	w.Header().Set("X-Server", peername.PeerName)

	writer := MakeNewWriter(w)

	recordData := &record.Record{RequestsID: fmt.Sprintf("%s-%d", access.RequestIdPrefix, access.Id)}

	ctx = context.WithValue(ctx, "X-Record", recordData)
	ctx = context.WithValue(ctx, "X-Real-IP-Geo-Code", code)
	ctx = context.WithValue(ctx, "X-Real-IP-Geo", geo)
	ctx = context.WithValue(ctx, "X-Real-IP", realIp)

	var startTime time.Time
	var endTime time.Time

	nextFuncErr := func() (resErr errors.WTError) {
		defer utils.Recover(logger.Logger, &resErr, "next func error")
		defer func() {
			endTime = time.Now()
		}()
		startTime = time.Now()
		next(writer, r.WithContext(ctx))
		return nil
	}()

	go func() {
		if recordData.User != nil {
			access.UserId = sql.NullInt64{
				Valid: true,
				Int64: recordData.User.Id,
			}
			access.UserUid = sql.NullString{
				Valid:  true,
				String: recordData.User.Uid,
			}
		}

		if len(recordData.UserToken) != 0 {
			access.UserToken = sql.NullString{
				Valid:  true,
				String: recordData.UserToken,
			}
		}

		if recordData.Role != nil {
			access.RoleId = sql.NullInt64{
				Valid: true,
				Int64: recordData.Role.ID,
			}
			access.RoleName = sql.NullString{
				Valid:  true,
				String: recordData.Role.Name,
			}
			access.RoleSign = sql.NullString{
				Valid:  true,
				String: recordData.Role.Sign,
			}
		}

		if recordData.Website != nil {
			access.WebId = sql.NullInt64{
				Valid: true,
				Int64: recordData.Website.ID,
			}
			access.WebName = sql.NullString{
				Valid:  true,
				String: recordData.Website.Name,
			}
		}

		if recordData.RequestWebsite != nil {
			access.RequestsWebId = sql.NullInt64{
				Valid: true,
				Int64: recordData.RequestWebsite.ID,
			}
			access.RequestsWebName = sql.NullString{
				Valid:  true,
				String: recordData.RequestWebsite.Name,
			}
		}

		if !strings.HasPrefix(path, "/api/v1/admin/accessrecord") {
			var body string
			if utf8.Valid(writer.Body) {
				body = string(writer.Body)
				if len(body) > 65000 {
					body = body[:65000]
				}
			} else {
				body = fmt.Sprintf("<media type: %s bytes: %d>", utils.GetMediaType(writer.Body), len(writer.Body))
			}

			access.ResponseBody = sql.NullString{
				Valid:  true,
				String: body,
			}
		}

		headerByte, err := utils.JsonMarshal(writer.Header())
		if err != nil {
			headerByte = []byte("{}")
		}
		header := string(headerByte)
		if len(header) > 6000 {
			header = "{}"
		}
		access.ResponseHeader = sql.NullString{
			Valid:  true,
			String: header,
		}

		writeError := writer.WriteError
		if len(writeError) > 2000 {
			writeError = writeError[:2000]
		}
		access.ResponseBodyError = sql.NullString{
			Valid:  len(writer.WriteError) != 0,
			String: writeError,
		}

		access.StatusCode = sql.NullInt64{
			Valid: writer.Status != 0,
			Int64: writer.Status,
		}

		if nextFuncErr != nil {
			panicError := nextFuncErr.Error()
			if len(panicError) > 2000 {
				panicError = panicError[:2000]
			}
			access.PanicError = sql.NullString{
				Valid:  true,
				String: panicError,
			}
		}

		var msgByte []byte
		if recordData.Err == nil {
			msgByte, err = utils.JsonMarshal(struct {
				Msg   string `json:"msg"`
				Stack string `json:"stack"`
			}{
				Msg:   recordData.Msg,
				Stack: recordData.Stack,
			})
		} else {
			msgByte, err = utils.JsonMarshal(struct {
				Msg        string `json:"msg"`
				ErrorMsg   string `json:"errorMsg"`
				ErrorCode  string `json:"errorCode"`
				ErrorStack string `json:"errorStack"`
				Stack      string `json:"stack"`
			}{
				Msg:        recordData.Msg,
				ErrorMsg:   recordData.Err.Message(),
				ErrorCode:  recordData.Err.Code(),
				ErrorStack: recordData.Err.Stack(),
				Stack:      recordData.Stack,
			})
		}
		if err != nil {
			msgByte, err = utils.JsonMarshal(struct {
				Msg   string `json:"msg"`
				Stack string `json:"stack"`
			}{
				Msg:   recordData.Msg,
				Stack: recordData.Stack,
			})
			if err != nil {
				msgByte = []byte("{}")
			}
		}
		msg := string(msgByte)
		if len(msg) > 6000 {
			msg = msg[:6000]
		}
		access.Message = sql.NullString{
			Valid:  true,
			String: msg,
		}

		use := endTime.Sub(startTime).Milliseconds() // 毫秒
		access.UseTime = sql.NullInt64{
			Valid: true,
			Int64: use,
		}

		access.StartAt = sql.NullTime{
			Valid: true,
			Time:  startTime,
		}

		access.EndAt = sql.NullTime{
			Valid: true,
			Time:  endTime,
		}

		if !notRecord {
			accessModel := db.NewAccessRecordModel(mysql.MySQLConn)
			err := accessModel.Update(context.Background(), access)
			if err != nil {
				logger.Logger.Error("mysql error: %s", err)
			}
		}
	}()
}

func accessRecordError(w http.ResponseWriter, r *http.Request, err errors.WTError, m string) {
	httpx.OkJsonCtx(r.Context(), w, &types.RespEmpty{
		Resp: respmsg.GetRespByError(r.Context(), respmsg.UnknownError, errors.WarpQuick(err), m),
	})
}

type NewWriter struct {
	http.ResponseWriter
	Status     int64
	Body       []byte
	WriteError string
}

func MakeNewWriter(w http.ResponseWriter) *NewWriter {
	return &NewWriter{
		ResponseWriter: w,
		Status:         0,
		Body:           []byte{},
		WriteError:     "",
	}
}

func (w *NewWriter) Header() http.Header {
	return w.ResponseWriter.Header()
}

func (w *NewWriter) WriteHeader(statusCode int) {
	w.Status = int64(statusCode)
	w.ResponseWriter.WriteHeader(statusCode)
}

func (w *NewWriter) Write(data []byte) (int, error) {
	res, err := w.ResponseWriter.Write(data)
	w.Body = data
	if err != nil {
		w.WriteError = err.Error()
	}
	return res, err
}

func (w *NewWriter) Hijack() (net.Conn, *bufio.ReadWriter, error) {
	h, ok := w.ResponseWriter.(http.Hijacker)
	if !ok {
		panic("response does not implement http.Hijacker")
	}

	return h.Hijack()
}

func (w *NewWriter) CloseNotify() <-chan bool {
	h, ok := w.ResponseWriter.(http.CloseNotifier)
	if !ok {
		panic("response does not implement http.CloseNotifier")
	}

	return h.CloseNotify()
}

func (w *NewWriter) Flush() {
	h, ok := w.ResponseWriter.(http.Flusher)
	if !ok {
		panic("response does not implement http.Flusher")
	}

	h.Flush()
}

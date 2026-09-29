package ocr

import (
	"bytes"
	"gitee.com/wuntsong-auth/backend/src/logger"
	"gitee.com/wuntsong-auth/backend/src/utils"
	ocr_api "github.com/alibabacloud-go/ocr-api-20210707/client"
	errors "github.com/wuntsong-org/wterrors"
)

type IDCard struct {
	Name string
	ID   string
}

func CheckIDCard(pic []byte) (idcard IDCard, resErr errors.WTError) {
	defer utils.Recover(logger.Logger, &resErr, "bad idcard")

	request := &ocr_api.RecognizeIdcardRequest{}
	request.SetOutputFigure(true)
	request.SetOutputQualityInfo(true)
	picReader := bytes.NewReader(pic)

	request.SetBody(picReader)

	response, err := OcrClient.RecognizeIdcard(request)
	if err != nil {
		return IDCard{}, errors.WarpQuick(err)
	}

	if response.Body.Code != nil || response.Body.Data == nil {
		return IDCard{}, errors.Errorf("bad idcard")
	}

	type FaceData struct {
		Name     string `json:"name"`
		IDNumber string `json:"idNumber"`
	}

	type Warning struct {
		CompletenessScore float64 `json:"completenessScore"` // 完整度
		IsCopy            int64   `json:"isCopy"`            // 是否复印
		IsReshoot         int64   `json:"isReshoot"`         // 是否翻拍
		QualityScore      float64 `json:"qualityScore"`      // 质量
		TamperScore       float64 `json:"tamperScore"`       // 篡改指数（越大越可能篡改）
	}

	type Face struct {
		Data    FaceData `json:"data"`
		Warning Warning  `json:"warning"`
	}

	type Data struct {
		Face Face `json:"face"`
	}

	type Body struct {
		Data Data `json:"data"`
	}

	data := Body{}
	err = utils.JsonUnmarshal([]byte(*response.Body.Data), &data)
	if err != nil {
		return IDCard{}, errors.WarpQuick(err)
	}

	if data.Data.Face.Warning.IsReshoot == 1 || data.Data.Face.Warning.IsCopy == 1 {
		return IDCard{}, errors.Errorf("bad idcard")
	}

	if data.Data.Face.Warning.TamperScore > 60 {
		return IDCard{}, errors.Errorf("bad idcard")
	}

	if data.Data.Face.Warning.QualityScore < 55 || data.Data.Face.Warning.CompletenessScore < 90 {
		return IDCard{}, errors.Errorf("bad idcard")
	}

	name := data.Data.Face.Data.Name
	idNumber := data.Data.Face.Data.IDNumber

	if !utils.IsValidChineseName(name) || !utils.IsValidIDCard(idNumber) {
		return IDCard{}, errors.Errorf("bad idcard")
	}

	return IDCard{
		Name: name,
		ID:   idNumber,
	}, nil
}

func CheckIDCardBack(pic []byte) (resErr errors.WTError) {
	defer utils.Recover(logger.Logger, &resErr, "bad idcard")

	request := &ocr_api.RecognizeIdcardRequest{}
	request.SetOutputFigure(true)
	request.SetOutputQualityInfo(true)
	picReader := bytes.NewReader(pic)

	request.SetBody(picReader)

	response, err := OcrClient.RecognizeIdcard(request)
	if err != nil {
		return errors.WarpQuick(err)
	}

	if response.Body.Code != nil || response.Body.Data == nil {
		return errors.Errorf("bad idcard")
	}

	type BackData struct {
		IssueAuthority string `json:"issueAuthority"`
		ValidPeriod    string `json:"validPeriod"`
	}

	type Warning struct {
		CompletenessScore float64 `json:"completenessScore"` // 完整度
		IsCopy            int64   `json:"isCopy"`            // 是否复印
		IsReshoot         int64   `json:"isReshoot"`         // 是否翻拍
		QualityScore      float64 `json:"qualityScore"`      // 质量
		TamperScore       float64 `json:"tamperScore"`       // 篡改指数（越大越可能篡改）
	}

	type Back struct {
		Data    BackData `json:"data"`
		Warning Warning  `json:"warning"`
	}

	type Data struct {
		Back Back `json:"back"`
	}

	type Body struct {
		Data Data `json:"data"`
	}

	data := Body{}
	err = utils.JsonUnmarshal([]byte(*response.Body.Data), &data)
	if err != nil {
		return errors.WarpQuick(err)
	}

	if data.Data.Back.Warning.IsReshoot == 1 || data.Data.Back.Warning.IsCopy == 1 {
		return errors.Errorf("bad idcard")
	}

	if data.Data.Back.Warning.TamperScore > 60 {
		return errors.Errorf("bad idcard")
	}

	if data.Data.Back.Warning.QualityScore < 55 || data.Data.Back.Warning.CompletenessScore < 90 {
		return errors.Errorf("bad idcard")
	}

	if len(data.Data.Back.Data.IssueAuthority) == 0 {
		return errors.Errorf("bad idcard")
	}

	if !utils.IsValidPeriod(data.Data.Back.Data.ValidPeriod) {
		return errors.Errorf("bad idcard")
	}

	return nil
}

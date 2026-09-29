package email

import (
	"bytes"
	"context"
	"crypto/tls"
	"database/sql"
	"fmt"
	"gitee.com/wuntsong-auth/backend/src/config"
	"gitee.com/wuntsong-auth/backend/src/model/db"
	"gitee.com/wuntsong-auth/backend/src/mysql"
	errors "github.com/wuntsong-org/wterrors"
	"html/template"
	"io"
	"net"
	"net/mail"
	"net/smtp"
	"path"
	"time"
)

var CodeTemplate *template.Template
var ImportCodeTemplate *template.Template
var ChangeTemplate *template.Template
var ChangeEmailTemplate *template.Template
var DeleteTemplate *template.Template
var BindTemplate *template.Template

func InitSmtp() errors.WTError {
	if len(config.BackendConfig.Smtp.Addr) == 0 {
		return errors.Errorf("smtp addr must be given")
	}

	if len(config.BackendConfig.Smtp.FromEmail) == 0 {
		return errors.Errorf("smtp from email must be given")
	}

	if len(config.BackendConfig.Smtp.UserName) == 0 {
		return errors.Errorf("smtp username must be given")
	}

	if len(config.BackendConfig.Smtp.Password) == 0 {
		return errors.Errorf("smtp pssword must be given")
	}

	if len(config.BackendConfig.Smtp.Sender) == 0 {
		return errors.Errorf("smtp sender must be given")
	}

	var err error
	CodeTemplate, err = template.ParseFiles(path.Join(config.BackendConfig.Smtp.TemplateFilePath, "code.txt"))
	if err != nil {
		return errors.WarpQuick(err)
	}

	ImportCodeTemplate, err = template.ParseFiles(path.Join(config.BackendConfig.Smtp.TemplateFilePath, "importcode.txt"))
	if err != nil {
		return errors.WarpQuick(err)
	}

	ChangeTemplate, err = template.ParseFiles(path.Join(config.BackendConfig.Smtp.TemplateFilePath, "change.txt"))
	if err != nil {
		return errors.WarpQuick(err)
	}

	ChangeEmailTemplate, err = template.ParseFiles(path.Join(config.BackendConfig.Smtp.TemplateFilePath, "changeEmail.txt"))
	if err != nil {
		return errors.WarpQuick(err)
	}

	DeleteTemplate, err = template.ParseFiles(path.Join(config.BackendConfig.Smtp.TemplateFilePath, "delete.txt"))
	if err != nil {
		return errors.WarpQuick(err)
	}

	BindTemplate, err = template.ParseFiles(path.Join(config.BackendConfig.Smtp.TemplateFilePath, "bind.txt"))
	if err != nil {
		return errors.WarpQuick(err)
	}

	return nil
}

func SendCode(code int64, email string) errors.WTError {
	if code > 999999 || code < 0 {
		return errors.Errorf("bad code")
	}

	data := struct {
		Sig  string
		Code string
	}{Sig: config.BackendConfig.Smtp.Sig, Code: fmt.Sprintf("%06d", code)}

	var buf bytes.Buffer
	err := CodeTemplate.Execute(&buf, data)
	if err != nil {
		return errors.WarpQuick(err)
	}

	content, err := io.ReadAll(&buf)
	if err != nil {
		return errors.WarpQuick(err)
	}

	err = SendEmail(context.Background(), fmt.Sprintf("%s验证码", config.BackendConfig.User.ReadableName), string(content), email, config.BackendConfig.Smtp.Sender, 0)
	if err != nil {
		return errors.WarpQuick(err)
	}

	return nil
}

func SendImportCode(code int64, email string) errors.WTError {
	if code > 999999 || code < 0 {
		return errors.Errorf("bad code")
	}

	data := struct {
		Sig  string
		Code string
	}{Sig: config.BackendConfig.Smtp.Sig, Code: fmt.Sprintf("%06d", code)}

	var buf bytes.Buffer
	err := ImportCodeTemplate.Execute(&buf, data)
	if err != nil {
		return errors.WarpQuick(err)
	}

	content, err := io.ReadAll(&buf)
	if err != nil {
		return errors.WarpQuick(err)
	}

	err = SendEmail(context.Background(), fmt.Sprintf("%s验证码", config.BackendConfig.User.ReadableName), string(content), email, config.BackendConfig.Smtp.Sender, 0)
	if err != nil {
		return errors.WarpQuick(err)
	}

	return nil
}

func SendChange(project string, email string) errors.WTError {
	data := struct {
		Sig     string
		Project string
	}{Sig: config.BackendConfig.Smtp.Sig, Project: project}

	var buf bytes.Buffer
	err := ChangeTemplate.Execute(&buf, data)
	if err != nil {
		return errors.WarpQuick(err)
	}

	content, err := io.ReadAll(&buf)
	if err != nil {
		return errors.WarpQuick(err)
	}

	err = SendEmail(context.Background(), "用户信息变更提醒", string(content), email, config.BackendConfig.Smtp.Sender, 0)
	if err != nil {
		return errors.WarpQuick(err)
	}

	return nil
}

func SendEmailChange(newEmail string, email string) errors.WTError {
	data := struct {
		Sig   string
		Email string
	}{Sig: config.BackendConfig.Smtp.Sig, Email: newEmail}

	var buf bytes.Buffer
	err := ChangeEmailTemplate.Execute(&buf, data)
	if err != nil {
		return errors.WarpQuick(err)
	}

	content, err := io.ReadAll(&buf)
	if err != nil {
		return errors.WarpQuick(err)
	}

	err = SendEmail(context.Background(), "用户信息变更提醒", string(content), email, config.BackendConfig.Smtp.Sender, 0)
	if err != nil {
		return errors.WarpQuick(err)
	}

	return nil
}

func SendDelete(email string) errors.WTError {
	data := struct {
		Sig string
	}{Sig: config.BackendConfig.Smtp.Sig}

	var buf bytes.Buffer
	err := DeleteTemplate.Execute(&buf, data)
	if err != nil {
		return errors.WarpQuick(err)
	}

	content, err := io.ReadAll(&buf)
	if err != nil {
		return errors.WarpQuick(err)
	}

	err = SendEmail(context.Background(), "用户注销提醒", string(content), email, config.BackendConfig.Smtp.Sender, 0)
	if err != nil {
		return errors.WarpQuick(err)
	}

	return nil
}

func SendBind(email string) errors.WTError {
	data := struct {
		Sig string
	}{Sig: config.BackendConfig.Smtp.Sig}

	var buf bytes.Buffer
	err := BindTemplate.Execute(&buf, data)
	if err != nil {
		return errors.WarpQuick(err)
	}

	content, err := io.ReadAll(&buf)
	if err != nil {
		return errors.WarpQuick(err)
	}

	err = SendEmail(context.Background(), "用户绑定提醒", string(content), email, config.BackendConfig.Smtp.Sender, 0)
	if err != nil {
		return errors.WarpQuick(err)
	}

	return nil
}

func SendEmail(ctx context.Context, subject, content string, toAddress string, sender string, senderID int64) errors.WTError {
	emailMessageModel := db.NewEmailMessageModel(mysql.MySQLConn)
	emailMessage := &db.EmailMessage{
		Email:    toAddress,
		Subject:  subject,
		Content:  content,
		Sender:   sender,
		SenderId: senderID,
		Success:  true,
	}

	host, _, err := net.SplitHostPort(config.BackendConfig.Smtp.Addr)
	if err != nil {
		return errors.WarpQuick(err)
	}

	auth := smtp.PlainAuth("",
		config.BackendConfig.Smtp.UserName,
		config.BackendConfig.Smtp.Password,
		host)

	tlsConfig := &tls.Config{
		InsecureSkipVerify: true,
		ServerName:         host,
	}

	from := mail.Address{Name: sender, Address: config.BackendConfig.Smtp.FromEmail}
	to := mail.Address{Name: "", Address: toAddress}

	headers := make(map[string]string)
	headers["From"] = from.String()
	headers["To"] = to.String()
	headers["Subject"] = subject
	headers["Content-Type"] = "text/plain; charset=utf-8"
	headers["Date"] = time.Now().Format("Mon, 02 Jan 2006 15:04:05 -0700")

	message := ""
	for k, v := range headers {
		message += fmt.Sprintf("%s: %s\r\n", k, v)
	}
	message += "\r\n" + content

	conn, err := tls.Dial("tcp", config.BackendConfig.Smtp.Addr, tlsConfig)
	if err != nil {
		return errors.WarpQuick(err)
	}

	c, err := smtp.NewClient(conn, host)
	if err != nil {
		return errors.WarpQuick(err)
	}

	// Auth
	err = c.Auth(auth)
	if err != nil {
		emailMessage.Success = false
		emailMessage.ErrorMsg = sql.NullString{
			Valid:  true,
			String: err.Error(),
		}
		_, _ = emailMessageModel.Insert(ctx, emailMessage)
		return errors.WarpQuick(err)
	}

	// To && From
	err = c.Mail(from.Address)
	if err != nil {
		emailMessage.Success = false
		emailMessage.ErrorMsg = sql.NullString{
			Valid:  true,
			String: err.Error(),
		}
		_, _ = emailMessageModel.Insert(ctx, emailMessage)
		return errors.WarpQuick(err)
	}

	err = c.Rcpt(to.Address)
	if err != nil {
		emailMessage.Success = false
		emailMessage.ErrorMsg = sql.NullString{
			Valid:  true,
			String: err.Error(),
		}
		_, _ = emailMessageModel.Insert(ctx, emailMessage)
		return errors.WarpQuick(err)
	}

	// Data
	w, err := c.Data()
	if err != nil {
		emailMessage.Success = false
		emailMessage.ErrorMsg = sql.NullString{
			Valid:  true,
			String: err.Error(),
		}
		_, _ = emailMessageModel.Insert(ctx, emailMessage)
		return errors.WarpQuick(err)
	}

	_, err = w.Write([]byte(message))
	if err != nil {
		emailMessage.Success = false
		emailMessage.ErrorMsg = sql.NullString{
			Valid:  true,
			String: err.Error(),
		}
		_, _ = emailMessageModel.Insert(ctx, emailMessage)
		return errors.WarpQuick(err)
	}

	err = w.Close()
	if err != nil {
		emailMessage.Success = false
		emailMessage.ErrorMsg = sql.NullString{
			Valid:  true,
			String: err.Error(),
		}
		_, _ = emailMessageModel.Insert(ctx, emailMessage)
		return errors.WarpQuick(err)
	}

	err = c.Quit()
	if err != nil {
		emailMessage.Success = false
		emailMessage.ErrorMsg = sql.NullString{
			Valid:  true,
			String: err.Error(),
		}
		_, _ = emailMessageModel.Insert(ctx, emailMessage)
		return errors.WarpQuick(err)
	}

	_, _ = emailMessageModel.Insert(ctx, emailMessage)
	return nil
}

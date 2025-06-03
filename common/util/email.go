package util

import (
	"bytes"
	"context"
	"email/common/errcode"
	"email/common/log"
	"encoding/base64"
	"encoding/json"
	"gopkg.in/gomail.v2"
	"html/template"
)

func SendEmail(smtpHost string, smtpPort int, tagName, smtpUser, smtpPass, from, to, subject, contentType, content string) error {
	// 创建邮件对象
	m := gomail.NewMessage()
	m.SetHeader("From", from)       // 发件人
	m.SetHeader("To", to)           // 收件人
	m.SetHeader("Subject", subject) // 邮件主题
	m.SetBody(contentType, content) // 邮件内容

	if tagName != "" {
		trace := map[string]string{
			"OpenTrace": "1",     // 打开邮件跟踪
			"LinkTrace": "1",     // 点击邮件里的URL跟踪
			"TagName":   tagName, // 控制台创建的标签tagname
		}
		jsonTrace, err := json.Marshal(trace)
		if err != nil {
			return err
		}
		base64Trace := base64.StdEncoding.EncodeToString(jsonTrace)
		m.SetHeader("X-AliDM-Trace", base64Trace)
	}

	d := gomail.NewDialer(smtpHost, smtpPort, smtpUser, smtpPass)
	d.SSL = true

	//发送邮件
	if err := d.DialAndSend(m); err != nil {
		return err // 返回错误信息
	}

	return nil // 成功时返回nil
}

// RenderEmailContent 根据模板和数据生成邮件内容
func RenderEmailContent(ctx context.Context, templateStr string, data map[string]interface{}) (string, error) {
	// 解析模板
	tmpl, err := template.New("emailTemplate").Parse(templateStr)
	if err != nil {
		log.New(ctx).Error("解析模板失败", "template", templateStr, "error", err)
		return "", errcode.ErrServer.WithCause(err).AppendMsg("解析模板失败")
	}

	// 渲染模板
	var buf bytes.Buffer
	err = tmpl.Execute(&buf, data)
	if err != nil {
		log.New(ctx).Error("渲染模板失败", "template", templateStr, "data", data, "error", err)
		return "", err
	}

	return buf.String(), nil
}

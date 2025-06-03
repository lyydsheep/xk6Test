package enum

import (
	"strings"
)

// SmtpCodeCategory 表示邮件发送结果的分类
type SmtpCodeCategory int

const (
	Success SmtpCodeCategory = iota
	// TemporaryFailure 可以重试
	TemporaryFailure
	// PermanentFailure 不可以重试
	PermanentFailure
	Unknown
)

// SmtpCode 统一表示一个 SMTP 返回码的结构
type SmtpCode struct {
	Name     string
	Code     string           // 原始返回码，如 "450", "550"
	Category SmtpCodeCategory // 分类
	Desc     string           // 描述信息
}

// AliyunSMTPCodes 阿里云 DM 的所有 SMTP 返回码集合
var AliyunSMTPCodes = []SmtpCode{
	{
		Name:     "Success",
		Code:     "250",
		Category: Success,
		Desc:     "邮件发送成功",
	},
	{
		Name:     "IP Exceed",
		Code:     "421 4.7.0 [TSS04] Messages from",
		Category: TemporaryFailure,
		Desc:     "该邮件内容涉嫌大量群发，被判为垃圾邮件或被多数收件人投诉为垃圾邮件。",
	},
	{
		Name:     "IP Exceed",
		Code:     "421  4.7.28 review our Bulk Email Senders Guidelines",
		Category: TemporaryFailure,
		Desc:     "该邮件内容涉嫌大量群发，被判为垃圾邮件或被多数收件人投诉为垃圾邮件。",
	},
	{
		Name:     "Access Denied",
		Code:     "423 Dns resolve error",
		Category: PermanentFailure,
		Desc:     "收信域名 MX 解析查询失败。",
	},
	{
		Name:     "Domain Frequency Limited",
		Code:     "427 Socks Connect to UNREACHABLE host",
		Category: PermanentFailure,
		Desc:     "目标主机不可到达，可能是对端拒绝连接，或者接收域名的邮件解析（MX）记录不存在",
	},
	{
		Name:     "Temporary Problem",
		Code:     "451 Temporary local problem - please try later",
		Category: TemporaryFailure,
		Desc:     "接收方系统临时故障。请稍后重试投递。如果重试投递失败，请直接反馈接收方邮件服务商检查、处理。",
	},
	{
		Name:     "RateLimited Due To Reputation",
		Code:     "451 4.7.650 The mail server",
		Category: TemporaryFailure,
		Desc:     "由于 IP 信誉被限流",
	},
	{
		Name:     "Mailbox Full",
		Code:     "452",
		Category: PermanentFailure,
		Desc:     "收信人邮箱已满。",
	},
	{
		Name:     "Dns Resolve Failed",
		Code:     "524 Host not found by dns resolve",
		Category: PermanentFailure,
		Desc:     "收信域名 MX 解析查询失败。",
	},

	{
		Name:     "Dns Resolve Failed",
		Code:     "526 No data by dns resolve",
		Category: PermanentFailure,
		Desc:     "收信域名 MX 解析查询失败。",
	},
	{
		Name:     "IP blacklisting",
		Code:     "554 5.7.1",
		Category: TemporaryFailure,
		Desc:     "发信 IP 被列入黑名单。",
	},
	{
		Name:     "Invalid address",
		Code:     "550  5.1.1",
		Category: PermanentFailure,
		Desc:     "请检查输入的收件人地址是否有误，如检查是否有多余的空格或特殊字符。",
	},
	{
		Name:     "Invalid address",
		Code:     "550  5.2.1  https://support.google.com/mail",
		Category: PermanentFailure,
		Desc:     "请检查输入的收件人地址是否有误，如检查是否有多余的空格或特殊字符。",
	},
	{
		Name:     "Invalid address",
		Code:     "550  5.4.1 Recipient address rejected: Access denied. For more information see https://aka.ms/EXOSmtpErrors",
		Category: PermanentFailure,
		Desc:     "请检查输入的收件人地址是否有误，如检查是否有多余的空格或特殊字符。",
	},
	{
		Name:     "IP blacklisting",
		Code:     "550  5.7.1",
		Category: TemporaryFailure,
		Desc:     "发信 IP 被列入黑名单。",
	},
	{
		Name:     "Frequency Limited",
		Code:     "550  Domain frequency limited",
		Category: PermanentFailure,
		Desc:     "发信域名频率超限。请暂停发信，稍后降低频率重新尝试发信。",
	},
	{
		Name:     "Mailbox Full",
		Code:     "552  Mailbox limit exeeded for this email address",
		Category: PermanentFailure,
		Desc:     "收信人邮箱已满。",
	},
	{
		Name:     "Invalid Address",
		Code:     "553",
		Category: PermanentFailure,
		Desc:     "请检查输入的收件人地址是否有误，如检查是否有多余的空格或特殊字符。",
	},
	{
		Name:     "IP blacklisting",
		Code:     "554  5.7.1",
		Category: PermanentFailure,
		Desc:     "发信 IP 被列入黑名单。",
	},
	{
		Name:     "Mailbox is disabled",
		Code:     "554  30 Sorry, your message to",
		Category: PermanentFailure,
		Desc:     "废弃邮箱",
	},
	{
		Name:     "Low bounce",
		Code:     "563 Rcptto is on the account-level",
		Category: TemporaryFailure,
		Desc:     "Rcptto is on the account-level bounce suppression list",
	},
}

// UnifiedCodeMap 是所有 SMTP 返回码的统一映射表
var UnifiedCodeMap = map[string]SmtpCode{}

func init() {
	register := func(code SmtpCode) {
		UnifiedCodeMap[code.Code] = code
	}
	for _, code := range AliyunSMTPCodes {
		register(code)
	}
}

// GetSmtpCodeCategory 获取给定 SMTP 返回码的分类
func GetSmtpCodeCategory(code string) SmtpCodeCategory {
	if smtpCode, ok := UnifiedCodeMap[code]; ok {
		return smtpCode.Category
	}
	return Unknown
}

// GetSmtpCode 获取原始 SMTP 返回码对应的结构体
func GetSmtpCode(message string) SmtpCode {
	for key, v := range UnifiedCodeMap {
		if strings.HasPrefix(message, key) {
			return v
		}
	}
	return SmtpCode{
		Name:     "Unknown",
		Code:     message,
		Desc:     "未知的 SMTP 返回码",
		Category: Unknown,
	}
}

package utils

import (
	"net"
	"net/url"
	"regexp"
	"unicode"
)

var PhoneRegex = regexp.MustCompile(`^1[3456789]\d{9}$`)
var PhoneCallRegex = regexp.MustCompile(`[0-9\-\s]+`)
var EmailRegex = regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
var Sha256Regex = regexp.MustCompile("^[a-fA-F0-9]{64}$")
var UIDRegex = regexp.MustCompile(`^[a-zA-Z0-9\-]{6,50}$`)
var IDCardRegex = regexp.MustCompile(`^[0-9Xx]+$`)
var CreditCodeRegex = regexp.MustCompile(`^[0-9A-Z]+$`)
var TotpSecretRegex = regexp.MustCompile(`^[A-Za-z0-9]{64,128}$`)
var WeChatRegex = regexp.MustCompile(`^[0-9A-Za-z_-]+$`)
var QQRegex = regexp.MustCompile(`^[0-9]+$`)
var ValidPeriodRegex = regexp.MustCompile(`^\d{4}\.\d{2}\.\d{2}-\d{4}\.\d{2}\.\d{2}$`)
var UUIDRegex = regexp.MustCompile(`^[0-9a-f]{8}(-[0-9a-f]{4}){3}-[0-9a-f]{12}$`)

func IsUUID(input string) bool {
	return UUIDRegex.MatchString(input)
}

func IsPhoneNumber(input string) bool {
	return PhoneRegex.MatchString(input)
}

func IsPhoneCall(input string) bool {
	return PhoneCallRegex.MatchString(input)
}

func IsWeChat(input string) bool {
	return WeChatRegex.MatchString(input)
}

func IsQQ(input string) bool {
	return QQRegex.MatchString(input)
}

func IsValidPeriod(input string) bool {
	return ValidPeriodRegex.MatchString(input)
}

func IsEmailAddress(input string) bool {
	return EmailRegex.MatchString(input)
}

func IsSha256(input string) bool {
	return Sha256Regex.MatchString(input)
}

func IsUID(uid string) bool {
	return UIDRegex.MatchString(uid)
}

func IsHttpOrHttps(input string) bool {
	u, err := url.Parse(input)
	if err != nil {
		return false
	}

	if u.Scheme != "http" && u.Scheme != "https" {
		return false
	}

	return true
}

func IsIP(input string) bool {
	ip := net.ParseIP(input)
	if ip == nil {
		return false
	}

	return true
}

func IsCIDR(input string) bool {
	_, _, err := net.ParseCIDR(input)
	if err != nil {
		return false
	}

	return true
}

func IsValidIDCard(idCard string) bool {
	// 判断身份证号码长度是否正确
	if len(idCard) != 15 && len(idCard) != 18 {
		return false
	}

	// 判断身份证号码是否由合法字符组成
	if !IDCardRegex.MatchString(idCard) {
		return false
	}
	return true
}

func IsValidCreditCode(creditCode string) bool {
	// 判断社会统一信用代码长度是否正确
	if len(creditCode) != 18 {
		return false
	}

	// 判断社会统一信用代码是否由合法字符组成
	if !CreditCodeRegex.MatchString(creditCode) {
		return false
	}
	return true
}

func IsValidChineseName(name string) bool {
	if len(name) < 2 {
		return false
	}

	for _, ch := range name {
		if !unicode.Is(unicode.Scripts["Han"], ch) && ch != '·' {
			return false
		}
	}
	return true
}

func IsValidChineseCompanyName(name string) bool {
	if len(name) < 2 {
		return false
	}

	for _, ch := range name {
		if !unicode.Is(unicode.Scripts["Han"], ch) && !unicode.IsLetter(ch) && !unicode.IsDigit(ch) && ch != ' ' && ch != '-' && ch != '(' && ch != ')' && ch != '（' && ch != '）' {
			return false
		}
	}
	return true
}

func IsTotpSecret(secret string) bool {
	return TotpSecretRegex.MatchString(secret)
}

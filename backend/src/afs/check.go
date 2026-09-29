package afs

import (
	afs "github.com/alibabacloud-go/afs-20180112/client"
)

func CheckCAPTCHA(SessionID, Token, Sig, Scene string, AppKey string, RemoteIP string) bool {
	request := new(afs.AuthenticateSigRequest)
	request.SetSig(Sig)
	request.SetSessionId(SessionID)
	request.SetToken(Token)
	request.SetRemoteIp(RemoteIP)
	request.SetScene(Scene)
	request.SetAppKey(AppKey)
	response, err := AFSClient.AuthenticateSig(request)
	if err != nil {
		return false
	}

	if response.Code == nil {
		return false
	} else if *response.Code == 100 {
		return true
	} else if *response.Code == 900 {
		return false
	} else {
		return false
	}
}

const (
	Pass       = 1
	CheckAgain = 2
	Banned     = 3
)

func CheckSilenceCAPTCHA(nvc string) int {
	request := new(afs.AnalyzeNvcRequest)
	request.SetData(nvc)
	request.SetScoreJsonStr("{\"200\":\"PASS\",\"400\":\"NC\",\"800\":\"BLOCK\",\"900\":\"BLOCK\"}")
	response, err := AFSClient.AnalyzeNvc(request)
	if err != nil {
		return Banned
	}

	if response.BizCode == nil {
		return Banned
	} else if *response.BizCode == "200" || *response.BizCode == "100" {
		return Pass
	} else if *response.BizCode == "900" || *response.BizCode == "800" {
		return Banned
	} else if *response.BizCode == "400" {
		return CheckAgain
	} else {
		return Banned
	}
}

package utils

import (
	"net/http"
	"strings"
)

var MediaTypeSuffixMap = map[string]string{
	"image/png":       "png",
	"image/jpeg":      "jpeg",
	"image/bmp":       "bmp",
	"video/avi":       "avi",
	"video/mp4":       "mp4",
	"audio/mpeg":      "mp3",
	"audio/wave":      "wav",
	"application/pdf": "pdf",
}

func GetMediaType(data []byte) string {
	contentTypeLst := strings.Split(http.DetectContentType(data), ";")
	if len(contentTypeLst) == 0 {
		return "application/octet-stream"
	}

	contentType := contentTypeLst[0]
	if len(contentType) == 0 {
		return "application/octet-stream"
	}

	return contentType
}

func IsAcceptType(fileType string) bool {
	_, ok := MediaTypeSuffixMap[fileType]
	return ok
}

func IsPdf(fileType string) bool {
	return fileType == "application/pdf"
}

func IsImage(fileType string) bool {
	return IsAcceptType(fileType) && JustImage(fileType)
}

func IsVideo(fileType string) bool {
	return IsAcceptType(fileType) && JustVideo(fileType)
}

func IsAudio(fileType string) bool {
	return IsAcceptType(fileType) && JustAudio(fileType)
}

func JustImage(fileType string) bool {
	return strings.HasPrefix(fileType, "image/")
}

func JustVideo(fileType string) bool {
	return strings.HasPrefix(fileType, "video/")
}

func JustAudio(fileType string) bool {
	return strings.HasPrefix(fileType, "audio/")
}

package yundun

import (
	"gitee.com/wuntsong-auth/backend/src/utils"
	green20220302 "github.com/alibabacloud-go/green-20220302/client"
	util "github.com/alibabacloud-go/tea-utils/v2/service"
	"github.com/alibabacloud-go/tea/tea"
	errors "github.com/wuntsong-org/wterrors"
	"net/http"
	"strings"
)

var TextLabel = []string{
	"political_content",
	"profanity",
	"contraband",
	"sexual_content",
	"violence",
	"nonsense",
	"negative_content",
	"religion",
	"cyberbullying",
	"C_customized",
}

func invokeText(text string, service string) (*green20220302.TextModerationResponse, errors.WTError) {
	var err error

	contentMap := map[string]interface{}{
		"content": text,
	}

	serviceParameters, err := utils.JsonMarshal(contentMap)
	if err != nil {
		return nil, errors.WarpQuick(err)
	}

	textModerationRequest := &green20220302.TextModerationRequest{
		Service:           tea.String(service),
		ServiceParameters: tea.String(string(serviceParameters)),
	}

	runtime := &util.RuntimeOptions{}

	response, err := YunDunClient.TextModerationWithOptions(textModerationRequest, runtime)
	if err != nil {
		return nil, errors.WarpQuick(err)
	}

	return response, nil
}

func CheckName(name string) (bool, errors.WTError) {
	var err error
	name = strings.TrimSpace(name)

	if len(name) > 600 {
		return false, nil
	} else if len(name) == 0 {
		return true, nil
	}

	response, err := invokeText(name, "nickname_detection")
	if err != nil {
		return false, errors.WarpQuick(err)
	} else if response == nil {
		return false, errors.Errorf("empty response")
	} else if *response.StatusCode != http.StatusOK || *response.Body.Code != http.StatusOK {
		if response.Body.Message != nil {
			return false, errors.Errorf("response code %d. response:%s", *response.Body.Code, *response.Body.Message)
		} else {
			return false, errors.Errorf("response code %d.", *response.Body.Code)
		}
	}

	body := response.Body
	imageModerationResponseData := body.Data
	labelString := imageModerationResponseData.Labels
	if len(*labelString) == 0 {
		return true, nil
	}

	labelLst := strings.Split(*labelString, ",")
	for _, l := range labelLst {
		if utils.InList(TextLabel, l) {
			return false, nil
		}
	}

	return true, nil
}

func CheckText(data string) (bool, errors.WTError) {
	var err error
	data = strings.TrimSpace(data)

	if len(data) > 600 {
		return false, nil
	} else if len(data) == 0 {
		return true, nil
	}

	response, err := invokeText(data, "comment_detection")
	if err != nil {
		return false, errors.WarpQuick(err)
	} else if response == nil {
		return false, errors.Errorf("empty response")
	} else if *response.StatusCode != http.StatusOK || *response.Body.Code != http.StatusOK {
		if response.Body.Message != nil {
			return false, errors.Errorf("response code %d. response:%s", *response.Body.Code, *response.Body.Message)
		} else {
			return false, errors.Errorf("response code %d.", *response.Body.Code)
		}
	}

	body := response.Body
	imageModerationResponseData := body.Data
	labelString := imageModerationResponseData.Labels
	if len(*labelString) == 0 {
		return true, nil
	}

	labelLst := strings.Split(*labelString, ",")
	for _, l := range labelLst {
		if utils.InList(TextLabel, l) {
			return false, nil
		}
	}

	return true, nil
}

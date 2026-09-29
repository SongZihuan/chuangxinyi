package yundun

import (
	"bytes"
	"fmt"
	"gitee.com/wuntsong-auth/backend/src/utils"
	green20220302 "github.com/alibabacloud-go/green-20220302/client"
	util "github.com/alibabacloud-go/tea-utils/v2/service"
	"github.com/alibabacloud-go/tea/tea"
	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	"github.com/google/uuid"
	errors "github.com/wuntsong-org/wterrors"
	"net/http"
	"time"
)

var PicBaseLineLabel = []string{
	"pornographic_adultContent",
	"pornographic_adultContent_tii",
	"sexual_suggestiveContent",
	"sexual_partialNudity",
	"sexual_affectionDisplay",
	"political_historicalNihility",
	"political_historicalNihility_tii",
	"political_politicalFigure",
	"political_politicalFigure_name_tii",
	"political_politicalFigure_metaphor_tii",
	"political_prohibitedPerson",
	"political_prohibitedPerson_tii",
	"political_taintedCelebrity",
	"political_taintedCelebrity_tii",
	"political_flag",
	"political_map",
	"political_TVLogo",
	"political_outfit",
	"political_religion_tii",
	"political_racism_tii",
	"violent_explosion",
	"violent_armedForces",
	"violent_crowding",
	"violent_gunKnives",
	"violent_gunKnives_tii",
	"violent_horrificContent",
	"violent_horrific_tii",
	"contraband_drug",
	"contraband_drug_tii",
	"contraband_gamble",
	"contraband_gamble_tii",
	"fraud_videoAbuse",
	"fraud_playerAbuse",
}

var PicProfilePhotoLabel = []string{
	"pornographic_adultContent",
	"pornographic_adultToys",
	"pornographic_artwork",
	"pornographic_adultContent_tii",
	"sexual_suggestiveContent",
	"sexual_breastBump",
	"sexual_cleavage",
	"sexual_femaleUnderwear",
	"sexual_femaleShoulder",
	"sexual_femaleLeg",
	"sexual_maleTopless",
	"sexual_cartoon",
	"sexual_pregnancy",
	"sexual_underage",
	"political_historicalNihility",
	"political_historicalNihility_tii",
	"political_politicalFigure_1",
	"political_politicalFigure_2",
	"political_politicalFigure_3",
	"political_politicalFigure_4",
	"political_politicalFigure_name_tii",
	"political_politicalFigure_metaphor_tii",
	"political_prohibitedPerson_1",
	"political_prohibitedPerson_2",
	"political_taintedCelebrity",
	"political_taintedCelebrity_tii",
	"political_Chinaflag",
	"political_otherflag",
	"political_Chinamap",
	"political_logo",
	"political_outfit",
	"political_medicalOutfit",
	"political_badge",
	"political_racism_tii",
	"violent_explosion",
	"violent_burning",
	"violent_armedForces",
	"violent_crowding",
	"violent_gun",
	"violent_Knives",
	"violent_gunKnives_tii",
	"violent_blood",
	"violent_horrific",
	"violent_horrific_tii",
	"contraband_drug",
	"contraband_drug_tii",
	"contraband_gamble",
	"contraband_gamble_tii",
	"contraband_certificate_tii",
	"religion_funeral",
	"religion_buddhism",
	"religion_christianity",
	"religion_muslim",
	"religion_tii",
	"racism_tii",
	"PDA_kiss",
	"PDA_physicalContact",
	"object_landmark",
	"object_rmb",
	"object_foreignCurrency",
	"object_wn",
	"object_carcrash",
	"object_candle",
	"object_flood",
	"pt_logotoSocialNetwork",
	"pt_qrCode",
	"pt_programCode",
	"pt_toDirectContact_tii",
	"pt_toSocialNetwork_tii",
	"pt_toShortVideos_tii",
	"pt_investment_tii",
	"pt_recruitment_tii",
	"pt_certificate_tii",
	"inappropriate_smoking",
	"inappropriate_drinking",
	"inappropriate_tattoo",
	"inappropriate_middleFinger",
	"inappropriate_foodWasting",
	"quality_meaningless",
	"logo_brand",
	"profanity_oral_tii",
	"profanity_offensive_tii",
	"meme_vulgar",
	"meme_metaphor",
}

// 文件上传token
var tokenData *green20220302.DescribeUploadTokenResponseBodyData = nil

// 创建上传文件客户端
func createOssClient(tokenData *green20220302.DescribeUploadTokenResponseBodyData) (*oss.Bucket, errors.WTError) {
	ossClient, err := oss.New(tea.StringValue(tokenData.OssInternetEndPoint), tea.StringValue(tokenData.AccessKeyId), tea.StringValue(tokenData.AccessKeySecret), oss.SecurityToken(tea.StringValue(tokenData.SecurityToken)))
	if err != nil {
		return nil, errors.WarpQuick(err)
	}
	bucket, err := ossClient.Bucket(tea.StringValue(tokenData.BucketName))
	if err != nil {
		return nil, errors.WarpQuick(err)
	}
	return bucket, nil
}

// 上传文件
func uploadFile(file []byte, fileType string) (string, errors.WTError) {
	var err error

	if tokenData == nil || tea.Int32Value(tokenData.Expiration) <= int32(time.Now().Unix()) {
		//获取文件上传临时token
		uploadTokenResponse, err := YunDunClient.DescribeUploadToken()
		if err != nil {
			return "", errors.WarpQuick(err)
		}
		tokenData = uploadTokenResponse.Body.Data
	}

	bucket, err := createOssClient(tokenData)
	if err != nil {
		return "", errors.WarpQuick(err)
	}

	suffix, ok := utils.MediaTypeSuffixMap[fileType]
	if !ok {
		return "", errors.Errorf("bad file type")
	}

	key, err := uuid.NewRandom()
	if err != nil {
		return "", errors.WarpQuick(err)
	}

	objectName := fmt.Sprintf("%s%s.%s", tea.StringValue(tokenData.FileNamePrefix), key, suffix)

	err = bucket.PutObject(objectName, bytes.NewReader(file))
	if err != nil {
		return "", errors.WarpQuick(err)
	}
	return objectName, nil
}

func invokePic(file []byte, fileType string, service string) (*green20220302.ImageModerationResponse, errors.WTError) {
	runtime := &util.RuntimeOptions{}
	var objectName, _ = uploadFile(file, fileType)

	//构建图片检测请求。
	serviceParameters, _ := utils.JsonMarshal(map[string]interface{}{
		"ossBucketName": tea.StringValue(tokenData.BucketName),
		"ossObjectName": objectName,
	},
	)
	imageModerationRequest := &green20220302.ImageModerationRequest{
		//图片检测service：内容安全控制台图片增强版规则配置的serviceCode，示例：baselineCheck
		Service:           tea.String(service),
		ServiceParameters: tea.String(string(serviceParameters)),
	}

	res, err := YunDunClient.ImageModerationWithOptions(imageModerationRequest, runtime)
	if err != nil {
		return nil, errors.WarpQuick(err)
	}

	return res, nil
}

func CheckBaseLinePic(file []byte, fileType string) (bool, errors.WTError) {
	var err error
	response, err := invokePic(file, fileType, "baselineCheck")
	if err != nil {
		return false, errors.WarpQuick(err)
	} else if response == nil {
		return false, errors.Errorf("empty response")
	} else if *response.StatusCode != http.StatusOK || *response.Body.Code != http.StatusOK {
		if response.Body.Msg != nil {
			return false, errors.Errorf("response code %d. response:%s", *response.Body.Code, *response.Body.Msg)
		} else {
			return false, errors.Errorf("response code %d.", *response.Body.Code)
		}
	}

	body := response.Body
	imageModerationResponseData := body.Data
	result := imageModerationResponseData.Result
	for _, r := range result {
		if r == nil || r.Label == nil || r.Confidence == nil {
			continue
		}
		if !utils.InList(PicBaseLineLabel, *r.Label) {
			continue
		}

		if *r.Confidence > 65 {
			return false, nil
		}
	}

	return true, nil
}

func CheckHeaderPic(file []byte, fileType string) (bool, errors.WTError) {
	var err error
	response, err := invokePic(file, fileType, "profilePhotoCheck")
	if err != nil {
		return false, errors.WarpQuick(err)
	} else if response == nil {
		return false, errors.Errorf("empty response")
	} else if *response.StatusCode != http.StatusOK || *response.Body.Code != http.StatusOK {
		if response.Body.Msg != nil {
			return false, errors.Errorf("response code %d. response:%s", *response.Body.Code, *response.Body.Msg)
		} else {
			return false, errors.Errorf("response code %d.", *response.Body.Code)
		}
	}

	body := response.Body
	imageModerationResponseData := body.Data
	result := imageModerationResponseData.Result
	for _, r := range result {
		if r == nil || r.Label == nil || r.Confidence == nil {
			continue
		}
		if !utils.InList(PicProfilePhotoLabel, *r.Label) {
			continue
		}

		if *r.Confidence > 65 {
			return false, nil
		}
	}

	return true, nil
}

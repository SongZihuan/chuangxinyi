package ocr

import (
	"bytes"
	"gitee.com/wuntsong-auth/backend/src/logger"
	"gitee.com/wuntsong-auth/backend/src/utils"
	ocr_api "github.com/alibabacloud-go/ocr-api-20210707/client"
	errors "github.com/wuntsong-org/wterrors"
)

type Company struct {
	LegalPerson string
	Name        string
	ID          string
}

func CheckCompany(pic []byte) (company Company, resErr errors.WTError) {
	defer utils.Recover(logger.Logger, &resErr, "bad license")

	request := &ocr_api.RecognizeBusinessLicenseRequest{}
	picReader := bytes.NewReader(pic)

	request.SetBody(picReader)

	response, err := OcrClient.RecognizeBusinessLicense(request)
	if err != nil {
		return Company{}, errors.WarpQuick(err)
	}

	if response.Body.Code != nil || response.Body.Data == nil {
		return Company{}, errors.Errorf("bad license 1")
	}

	type Data struct {
		CompanyName string `json:"companyName"`
		CreditCode  string `json:"creditCode"`
		LegalPerson string `json:"legalPerson"`
	}

	type Body struct {
		Data  Data  `json:"data"`
		FType int64 `json:"ftype"` // 是否复印件
	}

	data := Body{}
	err = utils.JsonUnmarshal([]byte(*response.Body.Data), &data)
	if err != nil {
		return Company{}, errors.WarpQuick(err)
	}

	if data.FType == 1 {
		return Company{}, errors.Errorf("bad license 2")
	}

	name := data.Data.CompanyName
	id := data.Data.CreditCode
	legalPerson := data.Data.LegalPerson

	if !utils.IsValidChineseCompanyName(name) || !utils.IsValidCreditCode(id) || !utils.IsValidChineseName(legalPerson) {
		return Company{}, errors.Errorf("bad license 3")
	}

	return Company{
		Name:        name,
		ID:          id,
		LegalPerson: legalPerson,
	}, nil
}

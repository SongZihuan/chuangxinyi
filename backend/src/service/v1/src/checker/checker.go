package checker

import (
	"gitee.com/wuntsong-auth/backend/src/coupons"
	"gitee.com/wuntsong-auth/backend/src/discount"
	"gitee.com/wuntsong-auth/backend/src/global/jwt"
	"gitee.com/wuntsong-auth/backend/src/global/model"
	"gitee.com/wuntsong-auth/backend/src/model/db"
	"gitee.com/wuntsong-auth/backend/src/service/v1/gen/types"
	"gitee.com/wuntsong-auth/backend/src/service/v1/src/action"
	"gitee.com/wuntsong-auth/backend/src/service/v1/src/warp"
	"gitee.com/wuntsong-auth/backend/src/utils"
	"github.com/wuntsong-org/checker"
	"github.com/wuntsong-org/go-zero-plus/rest/httpx"
	errors "github.com/wuntsong-org/wterrors"
	"net/http"
	"reflect"
)

type Validator struct {
	ch *checker.Checker
}

func (v *Validator) Validate(r *http.Request, data any) error {
	err := v.ch.Check(data)
	if err != nil {
		return err
	}
	return nil
}

func GetValidator() (httpx.Validator, errors.WTError) {
	ch := checker.NewChecker()

	ch.AddStringChecker("PhoneString", func(field *reflect.StructField, s string) errors.WTError {
		if !utils.IsPhoneNumber(s) {
			return errors.Errorf("is not a phone")
		}
		return nil
	})

	ch.AddStringChecker("EmailString", func(field *reflect.StructField, s string) errors.WTError {
		if !utils.IsEmailAddress(s) {
			return errors.Errorf("is not a email")
		}
		return nil
	})

	ch.AddStringChecker("ChineseNameString", func(field *reflect.StructField, s string) errors.WTError {
		if !utils.IsValidChineseName(s) {
			return errors.Errorf("is not a name")
		}
		return nil
	})

	ch.AddStringChecker("ChineseIDCardString", func(field *reflect.StructField, s string) errors.WTError {
		if !utils.IsValidIDCard(s) {
			return errors.Errorf("is not a name")
		}
		return nil
	})

	ch.AddStringChecker("UUIDString", func(field *reflect.StructField, s string) errors.WTError {
		if !utils.IsUUID(s) {
			return errors.Errorf("is not a uuid")
		}
		return nil
	})

	ch.AddStringChecker("ChineseCompanyNameString", func(field *reflect.StructField, s string) errors.WTError {
		if !utils.IsValidChineseCompanyName(s) {
			return errors.Errorf("is not a company name")
		}
		return nil
	})

	ch.AddStringChecker("ChineseCompanyIDString", func(field *reflect.StructField, s string) errors.WTError {
		if !utils.IsValidCreditCode(s) {
			return errors.Errorf("is not a company id")
		}
		return nil
	})

	ch.AddIntChecker("WebsiteFundingType", func(field *reflect.StructField, s int64) errors.WTError {
		if !db.IsWebsiteFundingType(s) {
			return errors.Errorf("is not a website funding type")
		}
		return nil
	})

	ch.AddStringChecker("WeChatString", func(field *reflect.StructField, s string) errors.WTError {
		if !utils.IsWeChat(s) {
			return errors.Errorf("is not a wechat")
		}
		return nil
	})

	ch.AddStringChecker("QQString", func(field *reflect.StructField, s string) errors.WTError {
		if !utils.IsQQ(s) {
			return errors.Errorf("is not a qq")
		}
		return nil
	})

	ch.AddIntChecker("PayStatus", func(field *reflect.StructField, i int64) errors.WTError {
		if !db.IsPayStatus(i) {
			return errors.Errorf("is not a pay status")
		}
		return nil
	})

	ch.AddIntChecker("DefrayStatus", func(field *reflect.StructField, i int64) errors.WTError {
		if !db.IsDefrayStatus(i) {
			return errors.Errorf("is not a defray status")
		}
		return nil
	})

	ch.AddIntChecker("WithdrawStatus", func(field *reflect.StructField, i int64) errors.WTError {
		if !db.IsWithdrawStatus(i) {
			return errors.Errorf("is not a withdraw status")
		}
		return nil
	})

	ch.AddIntChecker("InvoiceType", func(field *reflect.StructField, i int64) errors.WTError {
		if !db.IsInvoiceType(i) {
			return errors.Errorf("is not a invoice type")
		}
		return nil
	})

	ch.AddStringChecker("UserStatus", func(field *reflect.StructField, i string) errors.WTError {
		if !db.IsUserStatus(i) {
			return errors.Errorf("is not a user status")
		}
		return nil
	})

	ch.AddIntChecker("CouponsType", func(field *reflect.StructField, i int64) errors.WTError {
		if !coupons.IsCouponsType(i) {
			return errors.Errorf("is not a coupons type")
		}
		return nil
	})

	ch.AddIntChecker("DiscountType", func(field *reflect.StructField, i int64) errors.WTError {
		if !discount.IsDiscountType(i) {
			return errors.Errorf("is not a discount type")
		}
		return nil
	})

	ch.AddIntChecker("InvoiceStatus", func(field *reflect.StructField, i int64) errors.WTError {
		if !db.IsInvoiceStatus(i) {
			return errors.Errorf("is not a invoice status")
		}
		return nil
	})

	ch.AddIntChecker("WalletRecordType", func(field *reflect.StructField, i int64) errors.WTError {
		if !db.IsWalletRecordType(i) {
			return errors.Errorf("is not a wallet record type")
		}
		return nil
	})

	ch.AddStringChecker("HttpString", func(field *reflect.StructField, s string) errors.WTError {
		if !utils.IsHttpOrHttps(s) {
			return errors.Errorf("is not a http or https url")
		}
		return nil
	})

	ch.AddStringChecker("IPString", func(field *reflect.StructField, s string) errors.WTError {
		if !utils.IsIP(s) {
			return errors.Errorf("is not a ip")
		}
		return nil
	})

	ch.AddStringChecker("IPOrCIDRString", func(field *reflect.StructField, s string) errors.WTError {
		if !utils.IsIP(s) && !utils.IsCIDR(s) {
			return errors.Errorf("is not a ip or cidr")
		}
		return nil
	})

	ch.AddIntChecker("RoleStatus", func(field *reflect.StructField, i int64) errors.WTError {
		if !db.IsRoleStatus(i) {
			return errors.Errorf("not a role status")
		}
		return nil
	})

	ch.AddIntChecker("WebsiteStatus", func(field *reflect.StructField, i int64) errors.WTError {
		if !db.IsWebsiteStatus(i) {
			return errors.Errorf("not a website status")
		}
		return nil
	})

	ch.AddIntChecker("WebID", func(field *reflect.StructField, i int64) errors.WTError {
		if i == 0 {
			return nil
		}

		web := action.GetWebsite(i)
		if web.ID == warp.UnknownWebsite {
			return errors.Errorf("not a website")
		}
		return nil
	})

	ch.AddIntChecker("WebIDWithAll", func(field *reflect.StructField, i int64) errors.WTError {
		if i == 0 || i == -1 {
			return nil
		}

		web := action.GetWebsite(i)
		if web.ID == warp.UnknownWebsite {
			return errors.Errorf("not a website")
		}
		return nil
	})

	ch.AddIntChecker("WebIDNotCenter", func(field *reflect.StructField, i int64) errors.WTError {
		web := action.GetWebsite(i)
		if web.ID == warp.UnknownWebsite || web.ID == warp.UserCenterWebsite {
			return errors.Errorf("not a website")
		}
		return nil
	})

	ch.AddIntChecker("DBTimeType", func(field *reflect.StructField, i int64) errors.WTError {
		if !db.IsTimeType(i) {
			return errors.Errorf("not a time type")
		}
		return nil
	})

	ch.AddIntChecker("WorkOrderStatus", func(field *reflect.StructField, i int64) errors.WTError {
		if !db.IsWorkOrderStatus(i) {
			return errors.Errorf("not a work order status")
		}
		return nil
	})

	ch.AddIntChecker("PermissionStatus", func(field *reflect.StructField, i int64) errors.WTError {
		if !db.IsPolicyStatus(i) {
			return errors.Errorf("not a work permission status")
		}
		return nil
	})

	ch.AddIntChecker("WebsitePathMode", func(field *reflect.StructField, i int64) errors.WTError {
		if !db.IsWebsitePathMode(i) {
			return errors.Errorf("not a path mode")
		}
		return nil
	})

	ch.AddIntChecker("PathMode", func(field *reflect.StructField, i int64) errors.WTError {
		if !db.IsPathMode(i) {
			return errors.Errorf("not a path mode")
		}
		return nil
	})

	ch.AddIntChecker("CorsMode", func(field *reflect.StructField, i int64) errors.WTError {
		if !db.IsPathCorsModel(i) {
			return errors.Errorf("not a cors mode")
		}
		return nil
	})

	ch.AddIntChecker("AdminMode", func(field *reflect.StructField, i int64) errors.WTError {
		if !db.IsPathAdminMode(i) {
			return errors.Errorf("not a admin mode")
		}
		return nil
	})

	ch.AddIntChecker("BusyMode", func(field *reflect.StructField, i int64) errors.WTError {
		if !db.IsPathBusyMode(i) {
			return errors.Errorf("not a busy mode")
		}
		return nil
	})

	ch.AddIntChecker("CaptchaModel", func(field *reflect.StructField, i int64) errors.WTError {
		if !db.IsCaptchaMode(i) {
			return errors.Errorf("not a captcha mode")
		}
		return nil
	})

	ch.AddSliceChecker("HttpMethod", func(field *reflect.StructField, i any) errors.WTError {
		lst := i.([]string)

		for _, m := range lst {
			if m != http.MethodGet && m != http.MethodPost {
				return errors.Errorf("bad method")
			}
		}
		return nil
	})

	ch.AddIntChecker("MenuStatus", func(field *reflect.StructField, i int64) errors.WTError {
		if !db.IsMenuStatus(i) {
			return errors.Errorf("not a work menu status")
		}
		return nil
	})

	ch.AddIntChecker("WebsitePermissionStatus", func(field *reflect.StructField, i int64) errors.WTError {
		if !db.IsWebsitePolicyStatus(i) {
			return errors.Errorf("not a menu status")
		}
		return nil
	})

	ch.AddIntChecker("WebsitePathStatus", func(field *reflect.StructField, i int64) errors.WTError {
		if !db.IsWebsitePathStatus(i) {
			return errors.Errorf("not a menu status")
		}
		return nil
	})

	ch.AddIntChecker("PathStatus", func(field *reflect.StructField, i int64) errors.WTError {
		if !db.IsPathStatus(i) {
			return errors.Errorf("not a path status")
		}
		return nil
	})

	ch.AddIntChecker("ApplicationStatus", func(field *reflect.StructField, i int64) errors.WTError {
		if !db.IsApplicationStatus(i) {
			return errors.Errorf("not a application status")
		}
		return nil
	})

	ch.AddSliceChecker("PermissionSlice", func(field *reflect.StructField, i any) errors.WTError {
		lst, ok := i.([]string)
		if !ok {
			return errors.Errorf("not a slice")
		}

		for _, p := range lst {
			_, ok := model.PermissionsSign()[p]
			if !ok {
				return errors.Errorf("not a policy")
			}
		}
		return nil
	})

	ch.AddSliceChecker("WebsitePermissionSlice", func(field *reflect.StructField, i any) errors.WTError {
		lst, ok := i.([]string)
		if !ok {
			return errors.Errorf("not a slice")
		}

		for _, p := range lst {
			_, ok := model.WebsitePermissionsSign()[p]
			if !ok {
				return errors.Errorf("not a policy")
			}
		}
		return nil
	})

	ch.AddSliceChecker("SubPermissionSlice", func(field *reflect.StructField, i any) errors.WTError {
		lst, ok := i.([]string)
		if !ok {
			return errors.Errorf("not a slice")
		}

		for _, p := range lst {
			_, ok := jwt.UserSubTokenStringMap[p]
			if !ok {
				return errors.Errorf("not a sub policy")
			}
		}
		return nil
	})

	ch.AddMapChecker("WebsiteKeyMap", func(field *reflect.StructField, i any) errors.WTError {
		lst, ok := i.([]types.LabelValueRecord)
		if !ok {
			return errors.Errorf("not a slice")
		}

		for _, p := range lst {
			if len(p.Label) > 10 {
				return errors.Errorf("label too long")
			}
		}
		return nil
	})

	ch.AddIntChecker("OrderStatus", func(field *reflect.StructField, i int64) errors.WTError {
		if db.IsWorkOrderStatus(i) {
			return nil
		}
		return errors.Errorf("not a order status")
	})

	return &Validator{
		ch: ch,
	}, nil
}

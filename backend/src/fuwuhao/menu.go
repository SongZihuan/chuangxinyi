package fuwuhao

import (
	"gitee.com/wuntsong-auth/backend/src/config"
	"gitee.com/wuntsong-auth/backend/src/utils"
	"github.com/fastwego/offiaccount/apis/menu"
	errors "github.com/wuntsong-org/wterrors"
)

type SubButton struct {
	Name string `json:"name"`
	Type string `json:"type"`
	Url  string `json:"url,omitempty"`
	Key  string `json:"key,omitempty"`
}

type Button struct {
	Name      string      `json:"name"`
	SubButton []SubButton `json:"sub_button"`
}

type CreateReq struct {
	Button []Button `json:"button"`
}

const AboutUsContactKey = "aboutUsContact"
const AboutUsKefuKey = "aboutUsKefu"
const BindAuthKey = "bind"

func CreateMenu() errors.WTError {
	var err error
	// 修改自定义菜单不需要delete
	//resp1, err := menu.Delete(OffiAccount)
	//if err != nil {
	//	return errors.WarpQuick(err)
	//}
	//
	//var resp1Data Resp
	//err = utils.JsonUnmarshal(resp1, &resp1Data)
	//if err != nil {
	//	return errors.WarpQuick(err)
	//}
	//
	//if resp1Data.Errcode != 0 {
	//	return errors.Errorf("%d: %s", resp1Data.Errcode, resp1Data.Errmsg)
	//}

	req2, err := utils.JsonMarshal(CreateReq{
		Button: []Button{
			{
				Name: "关于我们",
				SubButton: []SubButton{
					{
						Name: "官方网站",
						Type: "view",
						Url:  config.BackendConfig.FuWuHao.Menu.AboutUsWebsite,
					},
					{
						Name: "联系我们",
						Type: "click",
						Key:  AboutUsContactKey,
					},
					{
						Name: "客服热线",
						Type: "click",
						Key:  AboutUsKefuKey,
					},
				},
			},
			{
				Name: "主要产品",
				SubButton: []SubButton{
					{
						Name: "创思域变",
						Type: "view",
						Url:  config.BackendConfig.FuWuHao.Menu.ProductVxwk,
					},
				},
			},
			{
				Name: "账号",
				SubButton: []SubButton{
					{
						Name: "绑定账号",
						Type: "click",
						Key:  BindAuthKey,
					},
					{
						Name: "我的账号",
						Type: "view",
						Url:  config.BackendConfig.FuWuHao.Menu.Auth,
					},
				},
			},
		},
	})

	resp2, err := menu.Create(OffiAccount, req2)
	if err != nil {
		return errors.WarpQuick(err)
	}

	var resp2Data Resp
	err = utils.JsonUnmarshal(resp2, &resp2Data)
	if err != nil {
		return errors.WarpQuick(err)
	}

	if resp2Data.Errcode != 0 {
		return errors.Errorf("%d: %s", resp2Data.Errcode, resp2Data.Errmsg)
	}

	return nil
}

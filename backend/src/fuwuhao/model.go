package fuwuhao

type TemplateValue struct {
	Value string `json:"value"`
}

type TemplateMsgReq struct {
	ToUser     string                   `json:"touser"`
	TemplateId string                   `json:"template_id"`
	Url        string                   `json:"url"`
	Data       map[string]TemplateValue `json:"data"`
}

type Resp struct {
	Errcode int    `json:"errcode"`
	Errmsg  string `json:"errmsg"`
}

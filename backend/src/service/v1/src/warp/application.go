package warp

import "gitee.com/wuntsong-auth/backend/src/service/v1/gen/types"

type Application struct {
	types.AdminApplication
}

func (a *Application) GetAdminApplicationType() types.AdminApplication {
	return a.AdminApplication
}

func (a *Application) GetApplicationType() types.Application {
	return types.Application{
		Name:    a.Name,
		WebName: a.WebName,
		WebUID:  a.WebUID,
		Url:     a.Url,
		Icon:    a.Icon,
		Sort:    a.Sort,
	}
}

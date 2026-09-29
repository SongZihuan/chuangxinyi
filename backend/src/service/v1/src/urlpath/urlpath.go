package urlpath

import (
	"gitee.com/wuntsong-auth/backend/src/config"
	"gitee.com/wuntsong-auth/backend/src/global/jwt"
	"gitee.com/wuntsong-auth/backend/src/global/model"
	"gitee.com/wuntsong-auth/backend/src/global/websocket"
	"gitee.com/wuntsong-auth/backend/src/model/db"
	"gitee.com/wuntsong-auth/backend/src/peers/wstype"
	"gitee.com/wuntsong-auth/backend/src/permission"
	"gitee.com/wuntsong-auth/backend/src/service/v1/gen/logic/center/userwstype"
	"gitee.com/wuntsong-auth/backend/src/service/v1/src/action"
	"gitee.com/wuntsong-auth/backend/src/service/v1/src/warp"
	"gitee.com/wuntsong-auth/backend/src/utils"
	"strings"
)

func UpdateRoleByMsg(roleID int64, msg websocket.WSMessage) {
	websocket.RoleConnMapMutex.Lock()
	defer websocket.RoleConnMapMutex.Unlock()

	lst, ok := websocket.RoleConnMap[roleID]
	if !ok {
		return
	}

	for _, i := range lst {
		websocket.WriteMessage(i, msg)
	}
}

func UpdateRoleByDB(role *db.Role, ch chan websocket.WSMessage) {
	var lst []chan websocket.WSMessage

	websocket.RoleConnMapMutex.Lock()
	defer websocket.RoleConnMapMutex.Unlock()

	if ch == nil {
		var ok bool
		lst, ok = websocket.RoleConnMap[role.Id]
		if !ok {
			lst = []chan websocket.WSMessage{}
		}
	} else {
		lst = []chan websocket.WSMessage{ch}
	}

	r := action.GetRole(role.Id, false)

	msg := websocket.WSMessage{
		Code: userwstype.UpdateRoleInfo,
		Data: r.GetRole(),
	}

	websocket.WritePeersMessage(wstype.PeersUpdateRoleInfo, struct {
		RoleID int64 `json:"roleID"`
	}{RoleID: r.ID}, msg)

	for _, i := range lst {
		websocket.WriteMessage(i, msg)
	}
}

func UpdateRole(role warp.Role, ch chan websocket.WSMessage) {
	var lst []chan websocket.WSMessage

	websocket.RoleConnMapMutex.Lock()
	defer websocket.RoleConnMapMutex.Unlock()

	if ch == nil {
		var ok bool
		lst, ok = websocket.RoleConnMap[role.ID]
		if !ok {
			lst = []chan websocket.WSMessage{}
		}
	} else {
		lst = []chan websocket.WSMessage{ch}
	}

	msg := websocket.WSMessage{
		Code: userwstype.UpdateRoleInfo,
		Data: role.GetRole(),
	}

	if ch == nil {
		websocket.WritePeersMessage(wstype.PeersUpdateRoleInfo, struct {
			RoleID int64 `json:"roleID"`
		}{RoleID: role.ID}, msg)
	}

	for _, i := range lst {
		websocket.WriteMessage(i, msg)
	}
}

func ChangeRole(userID int64, roleID int64) {
	websocket.UserConnMapMutex.Lock()
	defer websocket.UserConnMapMutex.Unlock()

	lst, ok := websocket.UserConnMap[userID]
	if !ok {
		lst = []chan websocket.WSMessage{}
	}

	role := action.GetRole(roleID, false)

	msg := websocket.WSMessage{
		Code: userwstype.RoleChange,
		Data: role.GetRole(),
	}

	websocket.WritePeersMessage(wstype.PeersRoleChange, struct {
		UserID int64 `json:"userID"`
	}{UserID: userID}, msg)

	for _, i := range lst {
		websocket.WriteMessage(i, msg)
	}
}

func ChangeRoleByMsg(userID int64, msg websocket.WSMessage) {
	websocket.UserConnMapMutex.Lock()
	defer websocket.UserConnMapMutex.Unlock()

	lst, ok := websocket.UserConnMap[userID]
	if !ok {
		return
	}

	for _, i := range lst {
		websocket.WriteMessage(i, msg)
	}
}

func ChangeRoleByDeleteByMsg(roleID int64, msg websocket.WSMessage) {
	websocket.RoleConnMapMutex.Lock()
	defer websocket.RoleConnMapMutex.Unlock()

	lst, ok := websocket.RoleConnMap[roleID]
	if !ok {
		return
	}

	for _, i := range lst {
		websocket.WriteMessage(i, msg)
	}
}

func ChangeRoleByDelete(roleID int64) {
	websocket.RoleConnMapMutex.Lock()
	defer websocket.RoleConnMapMutex.Unlock()

	lst, ok := websocket.RoleConnMap[roleID]
	if !ok {
		lst = []chan websocket.WSMessage{}
	}

	role := action.GetRoleBySign(config.BackendConfig.Admin.AnonymousRole.RoleSign, false)

	msg := websocket.WSMessage{
		Code: userwstype.RoleChange,
		Data: role.GetRole(),
	}

	websocket.WritePeersMessage(wstype.PeersRoleChangeByDelete, struct {
		RoleID int64 `json:"roleID"`
	}{RoleID: roleID}, msg)

	for _, i := range lst {
		websocket.WriteMessage(i, msg)
	}
}

func CheckUrlPath(path string, method string) (*warp.UrlPath, bool) {
	var matchLen int
	var matcher *warp.UrlPath

	methodPermission, ok := db.PathMethodStringMap[method]
	if !ok {
		methodPermission = db.PathPost
	}

	for _, u := range model.UrlPathMap() {
		if u.Status == db.PathStatusDelete {
			continue
		}

		waitPath := u.Path

		if len(waitPath) == 0 {
			continue
		}

		if !strings.HasPrefix(waitPath, "/") {
			waitPath = "/" + waitPath
		}

		if !strings.HasPrefix(waitPath, "/api/v1") {
			waitPath = "/api/v1" + waitPath
		}

		if strings.HasSuffix(waitPath, "/") {
			waitPath = waitPath[0 : len(waitPath)-1]
		}

		switch u.Mode {
		case db.PathModeComplete:
			if path != waitPath {
				continue
			}

			if len(waitPath) == 0 {
				continue
			}

			if len(waitPath) < matchLen {
				continue
			}

			if !permission.CheckPermissionInt64(u.MethodPermission, methodPermission) {
				continue
			}

			matchLen = len(waitPath)
			matcher = utils.GetPointer(u)
		case db.PathModePrefix:
			if path != waitPath && !strings.HasPrefix(path, waitPath+"/") {
				continue
			}

			if len(waitPath) < matchLen {
				continue
			}

			if len(waitPath) == matchLen && matcher != nil && matcher.Mode == db.PathModeComplete {
				continue
			}

			if !permission.CheckPermissionInt64(u.MethodPermission, methodPermission) {
				continue
			}

			matchLen = len(waitPath)
			matcher = utils.GetPointer(u)
		case db.PathModeRegex:
			l := u.Regex.FindString(path)
			if len(l) == 0 {
				continue
			}

			if len(l) < matchLen {
				continue
			}

			if len(l) == matchLen && matcher != nil && (matcher.Mode == db.PathModeComplete || matcher.Mode == db.PathModePrefix) {
				continue
			}

			if !permission.CheckPermissionInt64(u.MethodPermission, methodPermission) {
				continue
			}

			matchLen = len(l)
			matcher = utils.GetPointer(u)
		}
	}

	if matcher == nil {
		return nil, false
	}

	return matcher, true
}

func CheckMatherPermission(matcher *warp.UrlPath, role warp.Role, subType int) bool {
	if role.Sign == config.BackendConfig.Admin.RootRole.RoleSign {
		return true
	}

	if !CheckMatherPermissionOptions(matcher) {
		return false
	}

	subTypePermission, ok := jwt.UserSubTokenPermissionMap[subType]
	if !ok {
		return false
	}

	if matcher.IsOr {
		if !permission.HasOnePermission(role.PolicyPermissions, matcher.PolicyPermission) {
			return false
		}
	} else {
		if !permission.HasAllPermission(role.PolicyPermissions, matcher.PolicyPermission) {
			return false
		}
	}

	if !permission.CheckPermissionInt64(matcher.SubPolicyPermission, subTypePermission) {
		return false
	}

	return true
}

func CheckMatherPermissionOptions(matcher *warp.UrlPath) bool {
	if matcher.Status == db.PathStatusBanned {
		return false
	}

	return true
}

func CheckWebsiteUrlPath(path string, method string) (*warp.WebsiteUrlPath, bool) {
	var matchLen int
	var matcher *warp.WebsiteUrlPath

	methodPermission, ok := db.WebsitePathMethodStringMap[method]
	if !ok {
		methodPermission = db.PathPost
	}

	for _, u := range model.WebsiteUrlPathMap() {
		if u.Status == db.WebsitePathStatusDelete {
			continue
		}

		waitPath := u.Path

		if len(waitPath) == 0 {
			continue
		}

		if !strings.HasPrefix(waitPath, "/") {
			waitPath = "/" + waitPath
		}

		if !strings.HasPrefix(waitPath, "/api/v1") {
			waitPath = "/api/v1" + waitPath
		}

		if strings.HasSuffix(waitPath, "/") {
			waitPath = waitPath[0 : len(waitPath)-1]
		}

		switch u.Mode {
		case db.WebsitePathModeComplete:
			if path != waitPath {
				continue
			}

			if len(path) == 0 {
				continue
			}

			if len(path) < matchLen {
				continue
			}

			if !permission.CheckPermissionInt64(u.MethodPermission, methodPermission) {
				continue
			}

			matchLen = len(path)
			matcher = utils.GetPointer(u)
		case db.WebsitePathModePrefix:
			if path != waitPath && !strings.HasPrefix(path, waitPath+"/") {
				continue
			}

			if len(waitPath) == 0 {
				continue
			}

			if len(waitPath) < matchLen {
				continue
			}

			if len(waitPath) == matchLen && matcher != nil && matcher.Mode == db.WebsitePathModeComplete {
				continue
			}

			if !permission.CheckPermissionInt64(u.MethodPermission, methodPermission) {
				continue
			}

			matchLen = len(waitPath)
			matcher = utils.GetPointer(u)
		case db.WebsitePathModeRegex:
			l := u.Regex.FindString(path)
			if len(l) == 0 {
				continue
			}

			if len(l) < matchLen {
				continue
			}

			if len(l) == matchLen && matcher != nil && (matcher.Mode == db.WebsitePathModeComplete || matcher.Mode == db.WebsitePathModePrefix) {
				continue
			}

			if !permission.CheckPermissionInt64(u.MethodPermission, methodPermission) {
				continue
			}

			matchLen = len(l)
			matcher = utils.GetPointer(u)
		}
	}

	if matcher == nil {
		return nil, false
	}

	return matcher, true
}

func CheckWebsiteMatherPermission(matcher *warp.WebsiteUrlPath, web warp.Website) bool {
	if matcher.Status == db.WebsitePathStatusDelete {
		return false
	}

	if matcher.IsOr {
		if !permission.HasOnePermission(web.PolicyPermissions, matcher.PolicyPermission) {
			return false
		}
	} else {
		if !permission.HasAllPermission(web.PolicyPermissions, matcher.PolicyPermission) {
			return false
		}
	}

	return true
}

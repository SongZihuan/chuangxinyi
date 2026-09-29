package permission

import (
	"math/big"
)

func GetPermissions(permissions ...big.Int) big.Int {
	var p big.Int
	for _, i := range permissions {
		p = AddPermission(p, i)
	}
	return p
}

func AddPermissionInt64(permission, p int64) int64 {
	return permission | p
}

func AddPermission(permission, p big.Int) big.Int {
	var result big.Int
	ok := result.Or(&permission, &p)
	if ok == nil {
		return *big.NewInt(0)
	}
	return result
}

func CheckPermissionInt64(permission, p int64) bool {
	return permission&p != 0
}

func CheckPermission(permission, p big.Int) bool {
	var result big.Int
	ok := result.And(&permission, &p)
	if ok == nil {
		return false
	}
	return result.Cmp(big.NewInt(0)) != 0
}

func HasOnePermission(permission, pSet big.Int) bool {
	var result big.Int
	ok := result.And(&permission, &pSet)
	if ok == nil {
		return false
	}
	return result.Cmp(big.NewInt(0)) != 0
}

func HasAllPermission(permission, pSet big.Int) bool {
	var result big.Int
	ok := result.And(&permission, &pSet)
	if ok == nil {
		return false
	}
	return result.Cmp(&pSet) == 0
}

func ClearPermission(allPermission, pSet big.Int) big.Int {
	var result big.Int
	ok := result.And(&allPermission, &pSet)
	if ok == nil {
		return *big.NewInt(0)
	}
	return result
}

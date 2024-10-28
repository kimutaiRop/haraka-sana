package config

import (
	permissionModels "haraka-sana/permissions/models"
	"reflect"
)

var Permissions = newPermissionsRegistry()

func newPermissionsRegistry() *PermissionsRegistry {
	return &PermissionsRegistry{
		VIEW_ORGANIZATION: "VIEW ORGANIZATION",

		CREATE_ORDERS: "CREATE ORDERS",
		VIEW_ORDERS:   "VIEW ORDERS",
		EDIT_ORDERS:   "EDIT ORDERS",

		VIEW_STAFF:   "VIEW STAFF",
		CREATE_STAFF: "CREATE STAFF",
		EDIT_STAFF:   "EDIT STAFF",
	}
}

type PermissionsRegistry struct {
	VIEW_ORGANIZATION string

	CREATE_ORDERS string
	VIEW_ORDERS   string
	EDIT_ORDERS   string

	VIEW_PERMISSION string
	EDIT_PERMISSION string
	CREATE_POSITION string

	CREATE_STAFF string
	VIEW_STAFF   string
	EDIT_STAFF   string
}

func SeedPermissions() {
	val := reflect.ValueOf(Permissions).Elem()

	for i := 0; i < val.NumField(); i++ {
		permissionName := val.Field(i).Interface().(string)
		permission := permissionModels.Permission{Name: permissionName}
		DB.Where(permissionModels.Permission{Name: permission.Name}).FirstOrCreate(&permission)
	}

	adminPosition := permissionModels.Position{Name: "admin"}
	DB.Where(permissionModels.Position{Name: adminPosition.Name}).FirstOrCreate(&adminPosition)

	if adminPosition.ID != 0 {
		for i := 0; i < val.NumField(); i++ {
			permissionName := val.Field(i).Interface().(string)
			var p permissionModels.Permission
			DB.Where(permissionModels.Permission{Name: permissionName}).First(&p)
			if p.ID == 0 {
				continue
			}
			positionPermission := permissionModels.PositionPermission{
				PositionID:   adminPosition.ID,
				PermissionID: p.ID,
				Active:       true,
			}
			DB.Where(permissionModels.PositionPermission{PositionID: positionPermission.PositionID, PermissionID: positionPermission.PermissionID}).FirstOrCreate(&positionPermission)
		}
	}
}

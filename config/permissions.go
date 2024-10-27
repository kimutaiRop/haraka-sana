package config

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

	CREATE_STAFF string
	VIEW_STAFF   string
	EDIT_STAFF   string
}

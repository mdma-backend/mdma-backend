package permission

type Permission string

const (
	MeshNodeCreate Permission = "mesh_node_create"
	MeshNodeRead   Permission = "mesh_node_read"
	MeshNodeUpdate Permission = "mesh_node_update"
	MeshNodeDelete Permission = "mesh_node_delete"

	MeshNodeUpdateCreate Permission = "mesh_node_update_create"
	MeshNodeUpdateRead   Permission = "mesh_node_update_read"
	MeshNodeUpdateDelete Permission = "mesh_node_update_delete"

	DataCreate Permission = "data_create"
	DataRead   Permission = "data_read"
	DataDelete Permission = "data_delete"

	UserAccountCreate Permission = "user_account_create"
	UserAccountRead   Permission = "user_account_read"
	UserAccountUpdate Permission = "user_account_update"
	UserAccountDelete Permission = "user_account_delete"

	ServiceAccountCreate Permission = "service_account_create"
	ServiceAccountRead   Permission = "service_account_read"
	ServiceAccountUpdate Permission = "service_account_update"
	ServiceAccountDelete Permission = "service_account_delete"
)

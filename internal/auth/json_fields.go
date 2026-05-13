package auth

// JSON:API field names and resource type strings used when building
// user-resource responses. Extracted to satisfy goconst and to keep the
// API contract in one place — JSON keys never change without coordination.
const (
	resourceTypeUser = "users"

	attrDisplayName = "display_name"
	attrEmail       = "email"
	attrIsAdmin     = "is_admin"
	attrAuthSource  = "auth_source"
	attrRole        = "role"
)

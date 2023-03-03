package pivot

import "github.com/galexrt/arpanet/pkg/permify/models"

// UserRoles represents the database model of user roles relationships
type UserRoles struct {
	UserID uint `gorm:"primary_key" json:"user_id"`
	RoleID uint `gorm:"primary_key" json:"role_id"`
}

// TableName sets the table name
func (UserRoles) TableName() string {
	return models.TablePrefix + "user_roles"
}

package permissions

import (
	"database/sql"
	"testing"

	"github.com/galexrt/arpanet/query"
)

func TestGetAllPermissionsOfUser(t *testing.T) {
	query.DB, _ = sql.Open("mysql", "arpanet:changeme@tcp(127.0.0.1:3306)/arpanet?collation=utf8mb4_unicode_ci&parseTime=True&loc=Local")
	GetAllPermissionsOfUser(26061)
}

// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.

package query

import (
	"context"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"gorm.io/gorm/schema"

	"gorm.io/gen"
	"gorm.io/gen/field"

	"gorm.io/plugin/dbresolver"

	"github.com/galexrt/arpanet/model"
	"github.com/galexrt/arpanet/pkg/permify/models"
)

func newUser(db *gorm.DB, opts ...gen.DOOption) user {
	_user := user{}

	_user.userDo.UseDB(db, opts...)
	_user.userDo.UseModel(&model.User{})

	tableName := _user.userDo.TableName()
	_user.ALL = field.NewAsterisk(tableName)
	_user.ID = field.NewInt32(tableName, "id")
	_user.Identifier = field.NewString(tableName, "identifier")
	_user.Job = field.NewString(tableName, "job")
	_user.JobGrade = field.NewInt(tableName, "job_grade")
	_user.Firstname = field.NewString(tableName, "firstname")
	_user.Lastname = field.NewString(tableName, "lastname")
	_user.Dateofbirth = field.NewString(tableName, "dateofbirth")
	_user.Sex = field.NewField(tableName, "sex")
	_user.Height = field.NewString(tableName, "height")
	_user.Jail = field.NewInt32(tableName, "jail")
	_user.PhoneNumber = field.NewString(tableName, "phone_number")
	_user.Accounts = field.NewField(tableName, "accounts")
	_user.Disabled = field.NewBool(tableName, "disabled")
	_user.Visum = field.NewInt32(tableName, "visum")
	_user.Playtime = field.NewInt32(tableName, "playtime")
	_user.CreatedAt = field.NewTime(tableName, "created_at")
	_user.UpdatedAt = field.NewTime(tableName, "last_seen")
	_user.UserProps = userHasOneUserProps{
		db: db.Session(&gorm.Session{}),

		RelationField: field.NewRelation("UserProps", "model.UserProps"),
	}

	_user.UserLicenses = userHasManyUserLicenses{
		db: db.Session(&gorm.Session{}),

		RelationField: field.NewRelation("UserLicenses", "model.UserLicense"),
	}

	_user.Documents = userHasManyDocuments{
		db: db.Session(&gorm.Session{}),

		RelationField: field.NewRelation("Documents", "model.Document"),
		Responses: struct {
			field.RelationField
			Responses struct {
				field.RelationField
			}
			Mentions struct {
				field.RelationField
			}
			JobAccess struct {
				field.RelationField
			}
			UserAccess struct {
				field.RelationField
			}
		}{
			RelationField: field.NewRelation("Documents.Responses", "model.Document"),
			Responses: struct {
				field.RelationField
			}{
				RelationField: field.NewRelation("Documents.Responses.Responses", "model.Document"),
			},
			Mentions: struct {
				field.RelationField
			}{
				RelationField: field.NewRelation("Documents.Responses.Mentions", "model.DocumentMentions"),
			},
			JobAccess: struct {
				field.RelationField
			}{
				RelationField: field.NewRelation("Documents.Responses.JobAccess", "model.DocumentJobAccess"),
			},
			UserAccess: struct {
				field.RelationField
			}{
				RelationField: field.NewRelation("Documents.Responses.UserAccess", "model.DocumentUserAccess"),
			},
		},
		Mentions: struct {
			field.RelationField
		}{
			RelationField: field.NewRelation("Documents.Mentions", "model.DocumentMentions"),
		},
		JobAccess: struct {
			field.RelationField
		}{
			RelationField: field.NewRelation("Documents.JobAccess", "model.DocumentJobAccess"),
		},
		UserAccess: struct {
			field.RelationField
		}{
			RelationField: field.NewRelation("Documents.UserAccess", "model.DocumentUserAccess"),
		},
	}

	_user.Roles = userManyToManyRoles{
		db: db.Session(&gorm.Session{}),

		RelationField: field.NewRelation("Roles", "models.Role"),
		Permissions: struct {
			field.RelationField
		}{
			RelationField: field.NewRelation("Roles.Permissions", "models.Permission"),
		},
	}

	_user.Permissions = userManyToManyPermissions{
		db: db.Session(&gorm.Session{}),

		RelationField: field.NewRelation("Permissions", "models.Permission"),
	}

	_user.fillFieldMap()

	return _user
}

type user struct {
	userDo

	ALL         field.Asterisk
	ID          field.Int32
	Identifier  field.String
	Job         field.String
	JobGrade    field.Int
	Firstname   field.String
	Lastname    field.String
	Dateofbirth field.String
	Sex         field.Field
	Height      field.String
	Jail        field.Int32
	PhoneNumber field.String
	Accounts    field.Field
	Disabled    field.Bool
	Visum       field.Int32
	Playtime    field.Int32
	CreatedAt   field.Time
	UpdatedAt   field.Time
	UserProps   userHasOneUserProps

	UserLicenses userHasManyUserLicenses

	Documents userHasManyDocuments

	Roles userManyToManyRoles

	Permissions userManyToManyPermissions

	fieldMap map[string]field.Expr
}

func (u user) Table(newTableName string) *user {
	u.userDo.UseTable(newTableName)
	return u.updateTableName(newTableName)
}

func (u user) As(alias string) *user {
	u.userDo.DO = *(u.userDo.As(alias).(*gen.DO))
	return u.updateTableName(alias)
}

func (u *user) updateTableName(table string) *user {
	u.ALL = field.NewAsterisk(table)
	u.ID = field.NewInt32(table, "id")
	u.Identifier = field.NewString(table, "identifier")
	u.Job = field.NewString(table, "job")
	u.JobGrade = field.NewInt(table, "job_grade")
	u.Firstname = field.NewString(table, "firstname")
	u.Lastname = field.NewString(table, "lastname")
	u.Dateofbirth = field.NewString(table, "dateofbirth")
	u.Sex = field.NewField(table, "sex")
	u.Height = field.NewString(table, "height")
	u.Jail = field.NewInt32(table, "jail")
	u.PhoneNumber = field.NewString(table, "phone_number")
	u.Accounts = field.NewField(table, "accounts")
	u.Disabled = field.NewBool(table, "disabled")
	u.Visum = field.NewInt32(table, "visum")
	u.Playtime = field.NewInt32(table, "playtime")
	u.CreatedAt = field.NewTime(table, "created_at")
	u.UpdatedAt = field.NewTime(table, "last_seen")

	u.fillFieldMap()

	return u
}

func (u *user) GetFieldByName(fieldName string) (field.OrderExpr, bool) {
	_f, ok := u.fieldMap[fieldName]
	if !ok || _f == nil {
		return nil, false
	}
	_oe, ok := _f.(field.OrderExpr)
	return _oe, ok
}

func (u *user) fillFieldMap() {
	u.fieldMap = make(map[string]field.Expr, 22)
	u.fieldMap["id"] = u.ID
	u.fieldMap["identifier"] = u.Identifier
	u.fieldMap["job"] = u.Job
	u.fieldMap["job_grade"] = u.JobGrade
	u.fieldMap["firstname"] = u.Firstname
	u.fieldMap["lastname"] = u.Lastname
	u.fieldMap["dateofbirth"] = u.Dateofbirth
	u.fieldMap["sex"] = u.Sex
	u.fieldMap["height"] = u.Height
	u.fieldMap["jail"] = u.Jail
	u.fieldMap["phone_number"] = u.PhoneNumber
	u.fieldMap["accounts"] = u.Accounts
	u.fieldMap["disabled"] = u.Disabled
	u.fieldMap["visum"] = u.Visum
	u.fieldMap["playtime"] = u.Playtime
	u.fieldMap["created_at"] = u.CreatedAt
	u.fieldMap["last_seen"] = u.UpdatedAt

}

func (u user) clone(db *gorm.DB) user {
	u.userDo.ReplaceConnPool(db.Statement.ConnPool)
	return u
}

func (u user) replaceDB(db *gorm.DB) user {
	u.userDo.ReplaceDB(db)
	return u
}

type userHasOneUserProps struct {
	db *gorm.DB

	field.RelationField
}

func (a userHasOneUserProps) Where(conds ...field.Expr) *userHasOneUserProps {
	if len(conds) == 0 {
		return &a
	}

	exprs := make([]clause.Expression, 0, len(conds))
	for _, cond := range conds {
		exprs = append(exprs, cond.BeCond().(clause.Expression))
	}
	a.db = a.db.Clauses(clause.Where{Exprs: exprs})
	return &a
}

func (a userHasOneUserProps) WithContext(ctx context.Context) *userHasOneUserProps {
	a.db = a.db.WithContext(ctx)
	return &a
}

func (a userHasOneUserProps) Model(m *model.User) *userHasOneUserPropsTx {
	return &userHasOneUserPropsTx{a.db.Model(m).Association(a.Name())}
}

type userHasOneUserPropsTx struct{ tx *gorm.Association }

func (a userHasOneUserPropsTx) Find() (result *model.UserProps, err error) {
	return result, a.tx.Find(&result)
}

func (a userHasOneUserPropsTx) Append(values ...*model.UserProps) (err error) {
	targetValues := make([]interface{}, len(values))
	for i, v := range values {
		targetValues[i] = v
	}
	return a.tx.Append(targetValues...)
}

func (a userHasOneUserPropsTx) Replace(values ...*model.UserProps) (err error) {
	targetValues := make([]interface{}, len(values))
	for i, v := range values {
		targetValues[i] = v
	}
	return a.tx.Replace(targetValues...)
}

func (a userHasOneUserPropsTx) Delete(values ...*model.UserProps) (err error) {
	targetValues := make([]interface{}, len(values))
	for i, v := range values {
		targetValues[i] = v
	}
	return a.tx.Delete(targetValues...)
}

func (a userHasOneUserPropsTx) Clear() error {
	return a.tx.Clear()
}

func (a userHasOneUserPropsTx) Count() int64 {
	return a.tx.Count()
}

type userHasManyUserLicenses struct {
	db *gorm.DB

	field.RelationField
}

func (a userHasManyUserLicenses) Where(conds ...field.Expr) *userHasManyUserLicenses {
	if len(conds) == 0 {
		return &a
	}

	exprs := make([]clause.Expression, 0, len(conds))
	for _, cond := range conds {
		exprs = append(exprs, cond.BeCond().(clause.Expression))
	}
	a.db = a.db.Clauses(clause.Where{Exprs: exprs})
	return &a
}

func (a userHasManyUserLicenses) WithContext(ctx context.Context) *userHasManyUserLicenses {
	a.db = a.db.WithContext(ctx)
	return &a
}

func (a userHasManyUserLicenses) Model(m *model.User) *userHasManyUserLicensesTx {
	return &userHasManyUserLicensesTx{a.db.Model(m).Association(a.Name())}
}

type userHasManyUserLicensesTx struct{ tx *gorm.Association }

func (a userHasManyUserLicensesTx) Find() (result []*model.UserLicense, err error) {
	return result, a.tx.Find(&result)
}

func (a userHasManyUserLicensesTx) Append(values ...*model.UserLicense) (err error) {
	targetValues := make([]interface{}, len(values))
	for i, v := range values {
		targetValues[i] = v
	}
	return a.tx.Append(targetValues...)
}

func (a userHasManyUserLicensesTx) Replace(values ...*model.UserLicense) (err error) {
	targetValues := make([]interface{}, len(values))
	for i, v := range values {
		targetValues[i] = v
	}
	return a.tx.Replace(targetValues...)
}

func (a userHasManyUserLicensesTx) Delete(values ...*model.UserLicense) (err error) {
	targetValues := make([]interface{}, len(values))
	for i, v := range values {
		targetValues[i] = v
	}
	return a.tx.Delete(targetValues...)
}

func (a userHasManyUserLicensesTx) Clear() error {
	return a.tx.Clear()
}

func (a userHasManyUserLicensesTx) Count() int64 {
	return a.tx.Count()
}

type userHasManyDocuments struct {
	db *gorm.DB

	field.RelationField

	Responses struct {
		field.RelationField
		Responses struct {
			field.RelationField
		}
		Mentions struct {
			field.RelationField
		}
		JobAccess struct {
			field.RelationField
		}
		UserAccess struct {
			field.RelationField
		}
	}
	Mentions struct {
		field.RelationField
	}
	JobAccess struct {
		field.RelationField
	}
	UserAccess struct {
		field.RelationField
	}
}

func (a userHasManyDocuments) Where(conds ...field.Expr) *userHasManyDocuments {
	if len(conds) == 0 {
		return &a
	}

	exprs := make([]clause.Expression, 0, len(conds))
	for _, cond := range conds {
		exprs = append(exprs, cond.BeCond().(clause.Expression))
	}
	a.db = a.db.Clauses(clause.Where{Exprs: exprs})
	return &a
}

func (a userHasManyDocuments) WithContext(ctx context.Context) *userHasManyDocuments {
	a.db = a.db.WithContext(ctx)
	return &a
}

func (a userHasManyDocuments) Model(m *model.User) *userHasManyDocumentsTx {
	return &userHasManyDocumentsTx{a.db.Model(m).Association(a.Name())}
}

type userHasManyDocumentsTx struct{ tx *gorm.Association }

func (a userHasManyDocumentsTx) Find() (result []*model.Document, err error) {
	return result, a.tx.Find(&result)
}

func (a userHasManyDocumentsTx) Append(values ...*model.Document) (err error) {
	targetValues := make([]interface{}, len(values))
	for i, v := range values {
		targetValues[i] = v
	}
	return a.tx.Append(targetValues...)
}

func (a userHasManyDocumentsTx) Replace(values ...*model.Document) (err error) {
	targetValues := make([]interface{}, len(values))
	for i, v := range values {
		targetValues[i] = v
	}
	return a.tx.Replace(targetValues...)
}

func (a userHasManyDocumentsTx) Delete(values ...*model.Document) (err error) {
	targetValues := make([]interface{}, len(values))
	for i, v := range values {
		targetValues[i] = v
	}
	return a.tx.Delete(targetValues...)
}

func (a userHasManyDocumentsTx) Clear() error {
	return a.tx.Clear()
}

func (a userHasManyDocumentsTx) Count() int64 {
	return a.tx.Count()
}

type userManyToManyRoles struct {
	db *gorm.DB

	field.RelationField

	Permissions struct {
		field.RelationField
	}
}

func (a userManyToManyRoles) Where(conds ...field.Expr) *userManyToManyRoles {
	if len(conds) == 0 {
		return &a
	}

	exprs := make([]clause.Expression, 0, len(conds))
	for _, cond := range conds {
		exprs = append(exprs, cond.BeCond().(clause.Expression))
	}
	a.db = a.db.Clauses(clause.Where{Exprs: exprs})
	return &a
}

func (a userManyToManyRoles) WithContext(ctx context.Context) *userManyToManyRoles {
	a.db = a.db.WithContext(ctx)
	return &a
}

func (a userManyToManyRoles) Model(m *model.User) *userManyToManyRolesTx {
	return &userManyToManyRolesTx{a.db.Model(m).Association(a.Name())}
}

type userManyToManyRolesTx struct{ tx *gorm.Association }

func (a userManyToManyRolesTx) Find() (result []*models.Role, err error) {
	return result, a.tx.Find(&result)
}

func (a userManyToManyRolesTx) Append(values ...*models.Role) (err error) {
	targetValues := make([]interface{}, len(values))
	for i, v := range values {
		targetValues[i] = v
	}
	return a.tx.Append(targetValues...)
}

func (a userManyToManyRolesTx) Replace(values ...*models.Role) (err error) {
	targetValues := make([]interface{}, len(values))
	for i, v := range values {
		targetValues[i] = v
	}
	return a.tx.Replace(targetValues...)
}

func (a userManyToManyRolesTx) Delete(values ...*models.Role) (err error) {
	targetValues := make([]interface{}, len(values))
	for i, v := range values {
		targetValues[i] = v
	}
	return a.tx.Delete(targetValues...)
}

func (a userManyToManyRolesTx) Clear() error {
	return a.tx.Clear()
}

func (a userManyToManyRolesTx) Count() int64 {
	return a.tx.Count()
}

type userManyToManyPermissions struct {
	db *gorm.DB

	field.RelationField
}

func (a userManyToManyPermissions) Where(conds ...field.Expr) *userManyToManyPermissions {
	if len(conds) == 0 {
		return &a
	}

	exprs := make([]clause.Expression, 0, len(conds))
	for _, cond := range conds {
		exprs = append(exprs, cond.BeCond().(clause.Expression))
	}
	a.db = a.db.Clauses(clause.Where{Exprs: exprs})
	return &a
}

func (a userManyToManyPermissions) WithContext(ctx context.Context) *userManyToManyPermissions {
	a.db = a.db.WithContext(ctx)
	return &a
}

func (a userManyToManyPermissions) Model(m *model.User) *userManyToManyPermissionsTx {
	return &userManyToManyPermissionsTx{a.db.Model(m).Association(a.Name())}
}

type userManyToManyPermissionsTx struct{ tx *gorm.Association }

func (a userManyToManyPermissionsTx) Find() (result []*models.Permission, err error) {
	return result, a.tx.Find(&result)
}

func (a userManyToManyPermissionsTx) Append(values ...*models.Permission) (err error) {
	targetValues := make([]interface{}, len(values))
	for i, v := range values {
		targetValues[i] = v
	}
	return a.tx.Append(targetValues...)
}

func (a userManyToManyPermissionsTx) Replace(values ...*models.Permission) (err error) {
	targetValues := make([]interface{}, len(values))
	for i, v := range values {
		targetValues[i] = v
	}
	return a.tx.Replace(targetValues...)
}

func (a userManyToManyPermissionsTx) Delete(values ...*models.Permission) (err error) {
	targetValues := make([]interface{}, len(values))
	for i, v := range values {
		targetValues[i] = v
	}
	return a.tx.Delete(targetValues...)
}

func (a userManyToManyPermissionsTx) Clear() error {
	return a.tx.Clear()
}

func (a userManyToManyPermissionsTx) Count() int64 {
	return a.tx.Count()
}

type userDo struct{ gen.DO }

type IUserDo interface {
	gen.SubQuery
	Debug() IUserDo
	WithContext(ctx context.Context) IUserDo
	WithResult(fc func(tx gen.Dao)) gen.ResultInfo
	ReplaceDB(db *gorm.DB)
	ReadDB() IUserDo
	WriteDB() IUserDo
	As(alias string) gen.Dao
	Session(config *gorm.Session) IUserDo
	Columns(cols ...field.Expr) gen.Columns
	Clauses(conds ...clause.Expression) IUserDo
	Not(conds ...gen.Condition) IUserDo
	Or(conds ...gen.Condition) IUserDo
	Select(conds ...field.Expr) IUserDo
	Where(conds ...gen.Condition) IUserDo
	Order(conds ...field.Expr) IUserDo
	Distinct(cols ...field.Expr) IUserDo
	Omit(cols ...field.Expr) IUserDo
	Join(table schema.Tabler, on ...field.Expr) IUserDo
	LeftJoin(table schema.Tabler, on ...field.Expr) IUserDo
	RightJoin(table schema.Tabler, on ...field.Expr) IUserDo
	Group(cols ...field.Expr) IUserDo
	Having(conds ...gen.Condition) IUserDo
	Limit(limit int) IUserDo
	Offset(offset int) IUserDo
	Count() (count int64, err error)
	Scopes(funcs ...func(gen.Dao) gen.Dao) IUserDo
	Unscoped() IUserDo
	Create(values ...*model.User) error
	CreateInBatches(values []*model.User, batchSize int) error
	Save(values ...*model.User) error
	First() (*model.User, error)
	Take() (*model.User, error)
	Last() (*model.User, error)
	Find() ([]*model.User, error)
	FindInBatch(batchSize int, fc func(tx gen.Dao, batch int) error) (results []*model.User, err error)
	FindInBatches(result *[]*model.User, batchSize int, fc func(tx gen.Dao, batch int) error) error
	Pluck(column field.Expr, dest interface{}) error
	Delete(...*model.User) (info gen.ResultInfo, err error)
	Update(column field.Expr, value interface{}) (info gen.ResultInfo, err error)
	UpdateSimple(columns ...field.AssignExpr) (info gen.ResultInfo, err error)
	Updates(value interface{}) (info gen.ResultInfo, err error)
	UpdateColumn(column field.Expr, value interface{}) (info gen.ResultInfo, err error)
	UpdateColumnSimple(columns ...field.AssignExpr) (info gen.ResultInfo, err error)
	UpdateColumns(value interface{}) (info gen.ResultInfo, err error)
	UpdateFrom(q gen.SubQuery) gen.Dao
	Attrs(attrs ...field.AssignExpr) IUserDo
	Assign(attrs ...field.AssignExpr) IUserDo
	Joins(fields ...field.RelationField) IUserDo
	Preload(fields ...field.RelationField) IUserDo
	FirstOrInit() (*model.User, error)
	FirstOrCreate() (*model.User, error)
	FindByPage(offset int, limit int) (result []*model.User, count int64, err error)
	ScanByPage(result interface{}, offset int, limit int) (count int64, err error)
	Scan(result interface{}) (err error)
	Returning(value interface{}, columns ...string) IUserDo
	UnderlyingDB() *gorm.DB
	schema.Tabler
}

func (u userDo) Debug() IUserDo {
	return u.withDO(u.DO.Debug())
}

func (u userDo) WithContext(ctx context.Context) IUserDo {
	return u.withDO(u.DO.WithContext(ctx))
}

func (u userDo) ReadDB() IUserDo {
	return u.Clauses(dbresolver.Read)
}

func (u userDo) WriteDB() IUserDo {
	return u.Clauses(dbresolver.Write)
}

func (u userDo) Session(config *gorm.Session) IUserDo {
	return u.withDO(u.DO.Session(config))
}

func (u userDo) Clauses(conds ...clause.Expression) IUserDo {
	return u.withDO(u.DO.Clauses(conds...))
}

func (u userDo) Returning(value interface{}, columns ...string) IUserDo {
	return u.withDO(u.DO.Returning(value, columns...))
}

func (u userDo) Not(conds ...gen.Condition) IUserDo {
	return u.withDO(u.DO.Not(conds...))
}

func (u userDo) Or(conds ...gen.Condition) IUserDo {
	return u.withDO(u.DO.Or(conds...))
}

func (u userDo) Select(conds ...field.Expr) IUserDo {
	return u.withDO(u.DO.Select(conds...))
}

func (u userDo) Where(conds ...gen.Condition) IUserDo {
	return u.withDO(u.DO.Where(conds...))
}

func (u userDo) Exists(subquery interface{ UnderlyingDB() *gorm.DB }) IUserDo {
	return u.Where(field.CompareSubQuery(field.ExistsOp, nil, subquery.UnderlyingDB()))
}

func (u userDo) Order(conds ...field.Expr) IUserDo {
	return u.withDO(u.DO.Order(conds...))
}

func (u userDo) Distinct(cols ...field.Expr) IUserDo {
	return u.withDO(u.DO.Distinct(cols...))
}

func (u userDo) Omit(cols ...field.Expr) IUserDo {
	return u.withDO(u.DO.Omit(cols...))
}

func (u userDo) Join(table schema.Tabler, on ...field.Expr) IUserDo {
	return u.withDO(u.DO.Join(table, on...))
}

func (u userDo) LeftJoin(table schema.Tabler, on ...field.Expr) IUserDo {
	return u.withDO(u.DO.LeftJoin(table, on...))
}

func (u userDo) RightJoin(table schema.Tabler, on ...field.Expr) IUserDo {
	return u.withDO(u.DO.RightJoin(table, on...))
}

func (u userDo) Group(cols ...field.Expr) IUserDo {
	return u.withDO(u.DO.Group(cols...))
}

func (u userDo) Having(conds ...gen.Condition) IUserDo {
	return u.withDO(u.DO.Having(conds...))
}

func (u userDo) Limit(limit int) IUserDo {
	return u.withDO(u.DO.Limit(limit))
}

func (u userDo) Offset(offset int) IUserDo {
	return u.withDO(u.DO.Offset(offset))
}

func (u userDo) Scopes(funcs ...func(gen.Dao) gen.Dao) IUserDo {
	return u.withDO(u.DO.Scopes(funcs...))
}

func (u userDo) Unscoped() IUserDo {
	return u.withDO(u.DO.Unscoped())
}

func (u userDo) Create(values ...*model.User) error {
	if len(values) == 0 {
		return nil
	}
	return u.DO.Create(values)
}

func (u userDo) CreateInBatches(values []*model.User, batchSize int) error {
	return u.DO.CreateInBatches(values, batchSize)
}

// Save : !!! underlying implementation is different with GORM
// The method is equivalent to executing the statement: db.Clauses(clause.OnConflict{UpdateAll: true}).Create(values)
func (u userDo) Save(values ...*model.User) error {
	if len(values) == 0 {
		return nil
	}
	return u.DO.Save(values)
}

func (u userDo) First() (*model.User, error) {
	if result, err := u.DO.First(); err != nil {
		return nil, err
	} else {
		return result.(*model.User), nil
	}
}

func (u userDo) Take() (*model.User, error) {
	if result, err := u.DO.Take(); err != nil {
		return nil, err
	} else {
		return result.(*model.User), nil
	}
}

func (u userDo) Last() (*model.User, error) {
	if result, err := u.DO.Last(); err != nil {
		return nil, err
	} else {
		return result.(*model.User), nil
	}
}

func (u userDo) Find() ([]*model.User, error) {
	result, err := u.DO.Find()
	return result.([]*model.User), err
}

func (u userDo) FindInBatch(batchSize int, fc func(tx gen.Dao, batch int) error) (results []*model.User, err error) {
	buf := make([]*model.User, 0, batchSize)
	err = u.DO.FindInBatches(&buf, batchSize, func(tx gen.Dao, batch int) error {
		defer func() { results = append(results, buf...) }()
		return fc(tx, batch)
	})
	return results, err
}

func (u userDo) FindInBatches(result *[]*model.User, batchSize int, fc func(tx gen.Dao, batch int) error) error {
	return u.DO.FindInBatches(result, batchSize, fc)
}

func (u userDo) Attrs(attrs ...field.AssignExpr) IUserDo {
	return u.withDO(u.DO.Attrs(attrs...))
}

func (u userDo) Assign(attrs ...field.AssignExpr) IUserDo {
	return u.withDO(u.DO.Assign(attrs...))
}

func (u userDo) Joins(fields ...field.RelationField) IUserDo {
	for _, _f := range fields {
		u = *u.withDO(u.DO.Joins(_f))
	}
	return &u
}

func (u userDo) Preload(fields ...field.RelationField) IUserDo {
	for _, _f := range fields {
		u = *u.withDO(u.DO.Preload(_f))
	}
	return &u
}

func (u userDo) FirstOrInit() (*model.User, error) {
	if result, err := u.DO.FirstOrInit(); err != nil {
		return nil, err
	} else {
		return result.(*model.User), nil
	}
}

func (u userDo) FirstOrCreate() (*model.User, error) {
	if result, err := u.DO.FirstOrCreate(); err != nil {
		return nil, err
	} else {
		return result.(*model.User), nil
	}
}

func (u userDo) FindByPage(offset int, limit int) (result []*model.User, count int64, err error) {
	result, err = u.Offset(offset).Limit(limit).Find()
	if err != nil {
		return
	}

	if size := len(result); 0 < limit && 0 < size && size < limit {
		count = int64(size + offset)
		return
	}

	count, err = u.Offset(-1).Limit(-1).Count()
	return
}

func (u userDo) ScanByPage(result interface{}, offset int, limit int) (count int64, err error) {
	count, err = u.Count()
	if err != nil {
		return
	}

	err = u.Offset(offset).Limit(limit).Scan(result)
	return
}

func (u userDo) Scan(result interface{}) (err error) {
	return u.DO.Scan(result)
}

func (u userDo) Delete(models ...*model.User) (result gen.ResultInfo, err error) {
	return u.DO.Delete(models)
}

func (u *userDo) withDO(do gen.Dao) *userDo {
	u.DO = *do.(*gen.DO)
	return u
}

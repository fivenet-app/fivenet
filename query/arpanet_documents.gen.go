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
)

func newDocument(db *gorm.DB, opts ...gen.DOOption) document {
	_document := document{}

	_document.documentDo.UseDB(db, opts...)
	_document.documentDo.UseModel(&model.Document{})

	tableName := _document.documentDo.TableName()
	_document.ALL = field.NewAsterisk(tableName)
	_document.ID = field.NewUint(tableName, "id")
	_document.CreatedAt = field.NewTime(tableName, "created_at")
	_document.UpdatedAt = field.NewTime(tableName, "updated_at")
	_document.DeletedAt = field.NewField(tableName, "deleted_at")
	_document.Type = field.NewString(tableName, "content_type")
	_document.Title = field.NewString(tableName, "title")
	_document.Content = field.NewString(tableName, "content")
	_document.CreatorID = field.NewInt32(tableName, "creator")
	_document.CreatorJob = field.NewString(tableName, "creator_job")
	_document.Public = field.NewBool(tableName, "public")
	_document.ResponseID = field.NewUint(tableName, "response_id")
	_document.Responses = documentHasManyResponses{
		db: db.Session(&gorm.Session{}),

		RelationField: field.NewRelation("Responses", "model.Document"),
		Responses: struct {
			field.RelationField
		}{
			RelationField: field.NewRelation("Responses.Responses", "model.Document"),
		},
		Mentions: struct {
			field.RelationField
			Users struct {
				field.RelationField
				UserLicenses struct {
					field.RelationField
				}
				Documents struct {
					field.RelationField
				}
				Roles struct {
					field.RelationField
					Permissions struct {
						field.RelationField
					}
				}
				Permissions struct {
					field.RelationField
				}
			}
		}{
			RelationField: field.NewRelation("Responses.Mentions", "model.DocumentMentions"),
			Users: struct {
				field.RelationField
				UserLicenses struct {
					field.RelationField
				}
				Documents struct {
					field.RelationField
				}
				Roles struct {
					field.RelationField
					Permissions struct {
						field.RelationField
					}
				}
				Permissions struct {
					field.RelationField
				}
			}{
				RelationField: field.NewRelation("Responses.Mentions.Users", "model.User"),
				UserLicenses: struct {
					field.RelationField
				}{
					RelationField: field.NewRelation("Responses.Mentions.Users.UserLicenses", "model.UserLicense"),
				},
				Documents: struct {
					field.RelationField
				}{
					RelationField: field.NewRelation("Responses.Mentions.Users.Documents", "model.Document"),
				},
				Roles: struct {
					field.RelationField
					Permissions struct {
						field.RelationField
					}
				}{
					RelationField: field.NewRelation("Responses.Mentions.Users.Roles", "models.Role"),
					Permissions: struct {
						field.RelationField
					}{
						RelationField: field.NewRelation("Responses.Mentions.Users.Roles.Permissions", "models.Permission"),
					},
				},
				Permissions: struct {
					field.RelationField
				}{
					RelationField: field.NewRelation("Responses.Mentions.Users.Permissions", "models.Permission"),
				},
			},
		},
		JobAccess: struct {
			field.RelationField
		}{
			RelationField: field.NewRelation("Responses.JobAccess", "model.DocumentJobAccess"),
		},
		UserAccess: struct {
			field.RelationField
		}{
			RelationField: field.NewRelation("Responses.UserAccess", "model.DocumentUserAccess"),
		},
	}

	_document.Mentions = documentHasManyMentions{
		db: db.Session(&gorm.Session{}),

		RelationField: field.NewRelation("Mentions", "model.DocumentMentions"),
	}

	_document.JobAccess = documentHasManyJobAccess{
		db: db.Session(&gorm.Session{}),

		RelationField: field.NewRelation("JobAccess", "model.DocumentJobAccess"),
	}

	_document.UserAccess = documentHasManyUserAccess{
		db: db.Session(&gorm.Session{}),

		RelationField: field.NewRelation("UserAccess", "model.DocumentUserAccess"),
	}

	_document.fillFieldMap()

	return _document
}

type document struct {
	documentDo

	ALL        field.Asterisk
	ID         field.Uint
	CreatedAt  field.Time
	UpdatedAt  field.Time
	DeletedAt  field.Field
	Type       field.String
	Title      field.String
	Content    field.String
	CreatorID  field.Int32
	CreatorJob field.String
	Public     field.Bool
	ResponseID field.Uint
	Responses  documentHasManyResponses

	Mentions documentHasManyMentions

	JobAccess documentHasManyJobAccess

	UserAccess documentHasManyUserAccess

	fieldMap map[string]field.Expr
}

func (d document) Table(newTableName string) *document {
	d.documentDo.UseTable(newTableName)
	return d.updateTableName(newTableName)
}

func (d document) As(alias string) *document {
	d.documentDo.DO = *(d.documentDo.As(alias).(*gen.DO))
	return d.updateTableName(alias)
}

func (d *document) updateTableName(table string) *document {
	d.ALL = field.NewAsterisk(table)
	d.ID = field.NewUint(table, "id")
	d.CreatedAt = field.NewTime(table, "created_at")
	d.UpdatedAt = field.NewTime(table, "updated_at")
	d.DeletedAt = field.NewField(table, "deleted_at")
	d.Type = field.NewString(table, "content_type")
	d.Title = field.NewString(table, "title")
	d.Content = field.NewString(table, "content")
	d.CreatorID = field.NewInt32(table, "creator")
	d.CreatorJob = field.NewString(table, "creator_job")
	d.Public = field.NewBool(table, "public")
	d.ResponseID = field.NewUint(table, "response_id")

	d.fillFieldMap()

	return d
}

func (d *document) GetFieldByName(fieldName string) (field.OrderExpr, bool) {
	_f, ok := d.fieldMap[fieldName]
	if !ok || _f == nil {
		return nil, false
	}
	_oe, ok := _f.(field.OrderExpr)
	return _oe, ok
}

func (d *document) fillFieldMap() {
	d.fieldMap = make(map[string]field.Expr, 15)
	d.fieldMap["id"] = d.ID
	d.fieldMap["created_at"] = d.CreatedAt
	d.fieldMap["updated_at"] = d.UpdatedAt
	d.fieldMap["deleted_at"] = d.DeletedAt
	d.fieldMap["content_type"] = d.Type
	d.fieldMap["title"] = d.Title
	d.fieldMap["content"] = d.Content
	d.fieldMap["creator"] = d.CreatorID
	d.fieldMap["creator_job"] = d.CreatorJob
	d.fieldMap["public"] = d.Public
	d.fieldMap["response_id"] = d.ResponseID

}

func (d document) clone(db *gorm.DB) document {
	d.documentDo.ReplaceConnPool(db.Statement.ConnPool)
	return d
}

func (d document) replaceDB(db *gorm.DB) document {
	d.documentDo.ReplaceDB(db)
	return d
}

type documentHasManyResponses struct {
	db *gorm.DB

	field.RelationField

	Responses struct {
		field.RelationField
	}
	Mentions struct {
		field.RelationField
		Users struct {
			field.RelationField
			UserLicenses struct {
				field.RelationField
			}
			Documents struct {
				field.RelationField
			}
			Roles struct {
				field.RelationField
				Permissions struct {
					field.RelationField
				}
			}
			Permissions struct {
				field.RelationField
			}
		}
	}
	JobAccess struct {
		field.RelationField
	}
	UserAccess struct {
		field.RelationField
	}
}

func (a documentHasManyResponses) Where(conds ...field.Expr) *documentHasManyResponses {
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

func (a documentHasManyResponses) WithContext(ctx context.Context) *documentHasManyResponses {
	a.db = a.db.WithContext(ctx)
	return &a
}

func (a documentHasManyResponses) Model(m *model.Document) *documentHasManyResponsesTx {
	return &documentHasManyResponsesTx{a.db.Model(m).Association(a.Name())}
}

type documentHasManyResponsesTx struct{ tx *gorm.Association }

func (a documentHasManyResponsesTx) Find() (result []*model.Document, err error) {
	return result, a.tx.Find(&result)
}

func (a documentHasManyResponsesTx) Append(values ...*model.Document) (err error) {
	targetValues := make([]interface{}, len(values))
	for i, v := range values {
		targetValues[i] = v
	}
	return a.tx.Append(targetValues...)
}

func (a documentHasManyResponsesTx) Replace(values ...*model.Document) (err error) {
	targetValues := make([]interface{}, len(values))
	for i, v := range values {
		targetValues[i] = v
	}
	return a.tx.Replace(targetValues...)
}

func (a documentHasManyResponsesTx) Delete(values ...*model.Document) (err error) {
	targetValues := make([]interface{}, len(values))
	for i, v := range values {
		targetValues[i] = v
	}
	return a.tx.Delete(targetValues...)
}

func (a documentHasManyResponsesTx) Clear() error {
	return a.tx.Clear()
}

func (a documentHasManyResponsesTx) Count() int64 {
	return a.tx.Count()
}

type documentHasManyMentions struct {
	db *gorm.DB

	field.RelationField
}

func (a documentHasManyMentions) Where(conds ...field.Expr) *documentHasManyMentions {
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

func (a documentHasManyMentions) WithContext(ctx context.Context) *documentHasManyMentions {
	a.db = a.db.WithContext(ctx)
	return &a
}

func (a documentHasManyMentions) Model(m *model.Document) *documentHasManyMentionsTx {
	return &documentHasManyMentionsTx{a.db.Model(m).Association(a.Name())}
}

type documentHasManyMentionsTx struct{ tx *gorm.Association }

func (a documentHasManyMentionsTx) Find() (result []*model.DocumentMentions, err error) {
	return result, a.tx.Find(&result)
}

func (a documentHasManyMentionsTx) Append(values ...*model.DocumentMentions) (err error) {
	targetValues := make([]interface{}, len(values))
	for i, v := range values {
		targetValues[i] = v
	}
	return a.tx.Append(targetValues...)
}

func (a documentHasManyMentionsTx) Replace(values ...*model.DocumentMentions) (err error) {
	targetValues := make([]interface{}, len(values))
	for i, v := range values {
		targetValues[i] = v
	}
	return a.tx.Replace(targetValues...)
}

func (a documentHasManyMentionsTx) Delete(values ...*model.DocumentMentions) (err error) {
	targetValues := make([]interface{}, len(values))
	for i, v := range values {
		targetValues[i] = v
	}
	return a.tx.Delete(targetValues...)
}

func (a documentHasManyMentionsTx) Clear() error {
	return a.tx.Clear()
}

func (a documentHasManyMentionsTx) Count() int64 {
	return a.tx.Count()
}

type documentHasManyJobAccess struct {
	db *gorm.DB

	field.RelationField
}

func (a documentHasManyJobAccess) Where(conds ...field.Expr) *documentHasManyJobAccess {
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

func (a documentHasManyJobAccess) WithContext(ctx context.Context) *documentHasManyJobAccess {
	a.db = a.db.WithContext(ctx)
	return &a
}

func (a documentHasManyJobAccess) Model(m *model.Document) *documentHasManyJobAccessTx {
	return &documentHasManyJobAccessTx{a.db.Model(m).Association(a.Name())}
}

type documentHasManyJobAccessTx struct{ tx *gorm.Association }

func (a documentHasManyJobAccessTx) Find() (result []*model.DocumentJobAccess, err error) {
	return result, a.tx.Find(&result)
}

func (a documentHasManyJobAccessTx) Append(values ...*model.DocumentJobAccess) (err error) {
	targetValues := make([]interface{}, len(values))
	for i, v := range values {
		targetValues[i] = v
	}
	return a.tx.Append(targetValues...)
}

func (a documentHasManyJobAccessTx) Replace(values ...*model.DocumentJobAccess) (err error) {
	targetValues := make([]interface{}, len(values))
	for i, v := range values {
		targetValues[i] = v
	}
	return a.tx.Replace(targetValues...)
}

func (a documentHasManyJobAccessTx) Delete(values ...*model.DocumentJobAccess) (err error) {
	targetValues := make([]interface{}, len(values))
	for i, v := range values {
		targetValues[i] = v
	}
	return a.tx.Delete(targetValues...)
}

func (a documentHasManyJobAccessTx) Clear() error {
	return a.tx.Clear()
}

func (a documentHasManyJobAccessTx) Count() int64 {
	return a.tx.Count()
}

type documentHasManyUserAccess struct {
	db *gorm.DB

	field.RelationField
}

func (a documentHasManyUserAccess) Where(conds ...field.Expr) *documentHasManyUserAccess {
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

func (a documentHasManyUserAccess) WithContext(ctx context.Context) *documentHasManyUserAccess {
	a.db = a.db.WithContext(ctx)
	return &a
}

func (a documentHasManyUserAccess) Model(m *model.Document) *documentHasManyUserAccessTx {
	return &documentHasManyUserAccessTx{a.db.Model(m).Association(a.Name())}
}

type documentHasManyUserAccessTx struct{ tx *gorm.Association }

func (a documentHasManyUserAccessTx) Find() (result []*model.DocumentUserAccess, err error) {
	return result, a.tx.Find(&result)
}

func (a documentHasManyUserAccessTx) Append(values ...*model.DocumentUserAccess) (err error) {
	targetValues := make([]interface{}, len(values))
	for i, v := range values {
		targetValues[i] = v
	}
	return a.tx.Append(targetValues...)
}

func (a documentHasManyUserAccessTx) Replace(values ...*model.DocumentUserAccess) (err error) {
	targetValues := make([]interface{}, len(values))
	for i, v := range values {
		targetValues[i] = v
	}
	return a.tx.Replace(targetValues...)
}

func (a documentHasManyUserAccessTx) Delete(values ...*model.DocumentUserAccess) (err error) {
	targetValues := make([]interface{}, len(values))
	for i, v := range values {
		targetValues[i] = v
	}
	return a.tx.Delete(targetValues...)
}

func (a documentHasManyUserAccessTx) Clear() error {
	return a.tx.Clear()
}

func (a documentHasManyUserAccessTx) Count() int64 {
	return a.tx.Count()
}

type documentDo struct{ gen.DO }

type IDocumentDo interface {
	gen.SubQuery
	Debug() IDocumentDo
	WithContext(ctx context.Context) IDocumentDo
	WithResult(fc func(tx gen.Dao)) gen.ResultInfo
	ReplaceDB(db *gorm.DB)
	ReadDB() IDocumentDo
	WriteDB() IDocumentDo
	As(alias string) gen.Dao
	Session(config *gorm.Session) IDocumentDo
	Columns(cols ...field.Expr) gen.Columns
	Clauses(conds ...clause.Expression) IDocumentDo
	Not(conds ...gen.Condition) IDocumentDo
	Or(conds ...gen.Condition) IDocumentDo
	Select(conds ...field.Expr) IDocumentDo
	Where(conds ...gen.Condition) IDocumentDo
	Order(conds ...field.Expr) IDocumentDo
	Distinct(cols ...field.Expr) IDocumentDo
	Omit(cols ...field.Expr) IDocumentDo
	Join(table schema.Tabler, on ...field.Expr) IDocumentDo
	LeftJoin(table schema.Tabler, on ...field.Expr) IDocumentDo
	RightJoin(table schema.Tabler, on ...field.Expr) IDocumentDo
	Group(cols ...field.Expr) IDocumentDo
	Having(conds ...gen.Condition) IDocumentDo
	Limit(limit int) IDocumentDo
	Offset(offset int) IDocumentDo
	Count() (count int64, err error)
	Scopes(funcs ...func(gen.Dao) gen.Dao) IDocumentDo
	Unscoped() IDocumentDo
	Create(values ...*model.Document) error
	CreateInBatches(values []*model.Document, batchSize int) error
	Save(values ...*model.Document) error
	First() (*model.Document, error)
	Take() (*model.Document, error)
	Last() (*model.Document, error)
	Find() ([]*model.Document, error)
	FindInBatch(batchSize int, fc func(tx gen.Dao, batch int) error) (results []*model.Document, err error)
	FindInBatches(result *[]*model.Document, batchSize int, fc func(tx gen.Dao, batch int) error) error
	Pluck(column field.Expr, dest interface{}) error
	Delete(...*model.Document) (info gen.ResultInfo, err error)
	Update(column field.Expr, value interface{}) (info gen.ResultInfo, err error)
	UpdateSimple(columns ...field.AssignExpr) (info gen.ResultInfo, err error)
	Updates(value interface{}) (info gen.ResultInfo, err error)
	UpdateColumn(column field.Expr, value interface{}) (info gen.ResultInfo, err error)
	UpdateColumnSimple(columns ...field.AssignExpr) (info gen.ResultInfo, err error)
	UpdateColumns(value interface{}) (info gen.ResultInfo, err error)
	UpdateFrom(q gen.SubQuery) gen.Dao
	Attrs(attrs ...field.AssignExpr) IDocumentDo
	Assign(attrs ...field.AssignExpr) IDocumentDo
	Joins(fields ...field.RelationField) IDocumentDo
	Preload(fields ...field.RelationField) IDocumentDo
	FirstOrInit() (*model.Document, error)
	FirstOrCreate() (*model.Document, error)
	FindByPage(offset int, limit int) (result []*model.Document, count int64, err error)
	ScanByPage(result interface{}, offset int, limit int) (count int64, err error)
	Scan(result interface{}) (err error)
	Returning(value interface{}, columns ...string) IDocumentDo
	UnderlyingDB() *gorm.DB
	schema.Tabler
}

func (d documentDo) Debug() IDocumentDo {
	return d.withDO(d.DO.Debug())
}

func (d documentDo) WithContext(ctx context.Context) IDocumentDo {
	return d.withDO(d.DO.WithContext(ctx))
}

func (d documentDo) ReadDB() IDocumentDo {
	return d.Clauses(dbresolver.Read)
}

func (d documentDo) WriteDB() IDocumentDo {
	return d.Clauses(dbresolver.Write)
}

func (d documentDo) Session(config *gorm.Session) IDocumentDo {
	return d.withDO(d.DO.Session(config))
}

func (d documentDo) Clauses(conds ...clause.Expression) IDocumentDo {
	return d.withDO(d.DO.Clauses(conds...))
}

func (d documentDo) Returning(value interface{}, columns ...string) IDocumentDo {
	return d.withDO(d.DO.Returning(value, columns...))
}

func (d documentDo) Not(conds ...gen.Condition) IDocumentDo {
	return d.withDO(d.DO.Not(conds...))
}

func (d documentDo) Or(conds ...gen.Condition) IDocumentDo {
	return d.withDO(d.DO.Or(conds...))
}

func (d documentDo) Select(conds ...field.Expr) IDocumentDo {
	return d.withDO(d.DO.Select(conds...))
}

func (d documentDo) Where(conds ...gen.Condition) IDocumentDo {
	return d.withDO(d.DO.Where(conds...))
}

func (d documentDo) Exists(subquery interface{ UnderlyingDB() *gorm.DB }) IDocumentDo {
	return d.Where(field.CompareSubQuery(field.ExistsOp, nil, subquery.UnderlyingDB()))
}

func (d documentDo) Order(conds ...field.Expr) IDocumentDo {
	return d.withDO(d.DO.Order(conds...))
}

func (d documentDo) Distinct(cols ...field.Expr) IDocumentDo {
	return d.withDO(d.DO.Distinct(cols...))
}

func (d documentDo) Omit(cols ...field.Expr) IDocumentDo {
	return d.withDO(d.DO.Omit(cols...))
}

func (d documentDo) Join(table schema.Tabler, on ...field.Expr) IDocumentDo {
	return d.withDO(d.DO.Join(table, on...))
}

func (d documentDo) LeftJoin(table schema.Tabler, on ...field.Expr) IDocumentDo {
	return d.withDO(d.DO.LeftJoin(table, on...))
}

func (d documentDo) RightJoin(table schema.Tabler, on ...field.Expr) IDocumentDo {
	return d.withDO(d.DO.RightJoin(table, on...))
}

func (d documentDo) Group(cols ...field.Expr) IDocumentDo {
	return d.withDO(d.DO.Group(cols...))
}

func (d documentDo) Having(conds ...gen.Condition) IDocumentDo {
	return d.withDO(d.DO.Having(conds...))
}

func (d documentDo) Limit(limit int) IDocumentDo {
	return d.withDO(d.DO.Limit(limit))
}

func (d documentDo) Offset(offset int) IDocumentDo {
	return d.withDO(d.DO.Offset(offset))
}

func (d documentDo) Scopes(funcs ...func(gen.Dao) gen.Dao) IDocumentDo {
	return d.withDO(d.DO.Scopes(funcs...))
}

func (d documentDo) Unscoped() IDocumentDo {
	return d.withDO(d.DO.Unscoped())
}

func (d documentDo) Create(values ...*model.Document) error {
	if len(values) == 0 {
		return nil
	}
	return d.DO.Create(values)
}

func (d documentDo) CreateInBatches(values []*model.Document, batchSize int) error {
	return d.DO.CreateInBatches(values, batchSize)
}

// Save : !!! underlying implementation is different with GORM
// The method is equivalent to executing the statement: db.Clauses(clause.OnConflict{UpdateAll: true}).Create(values)
func (d documentDo) Save(values ...*model.Document) error {
	if len(values) == 0 {
		return nil
	}
	return d.DO.Save(values)
}

func (d documentDo) First() (*model.Document, error) {
	if result, err := d.DO.First(); err != nil {
		return nil, err
	} else {
		return result.(*model.Document), nil
	}
}

func (d documentDo) Take() (*model.Document, error) {
	if result, err := d.DO.Take(); err != nil {
		return nil, err
	} else {
		return result.(*model.Document), nil
	}
}

func (d documentDo) Last() (*model.Document, error) {
	if result, err := d.DO.Last(); err != nil {
		return nil, err
	} else {
		return result.(*model.Document), nil
	}
}

func (d documentDo) Find() ([]*model.Document, error) {
	result, err := d.DO.Find()
	return result.([]*model.Document), err
}

func (d documentDo) FindInBatch(batchSize int, fc func(tx gen.Dao, batch int) error) (results []*model.Document, err error) {
	buf := make([]*model.Document, 0, batchSize)
	err = d.DO.FindInBatches(&buf, batchSize, func(tx gen.Dao, batch int) error {
		defer func() { results = append(results, buf...) }()
		return fc(tx, batch)
	})
	return results, err
}

func (d documentDo) FindInBatches(result *[]*model.Document, batchSize int, fc func(tx gen.Dao, batch int) error) error {
	return d.DO.FindInBatches(result, batchSize, fc)
}

func (d documentDo) Attrs(attrs ...field.AssignExpr) IDocumentDo {
	return d.withDO(d.DO.Attrs(attrs...))
}

func (d documentDo) Assign(attrs ...field.AssignExpr) IDocumentDo {
	return d.withDO(d.DO.Assign(attrs...))
}

func (d documentDo) Joins(fields ...field.RelationField) IDocumentDo {
	for _, _f := range fields {
		d = *d.withDO(d.DO.Joins(_f))
	}
	return &d
}

func (d documentDo) Preload(fields ...field.RelationField) IDocumentDo {
	for _, _f := range fields {
		d = *d.withDO(d.DO.Preload(_f))
	}
	return &d
}

func (d documentDo) FirstOrInit() (*model.Document, error) {
	if result, err := d.DO.FirstOrInit(); err != nil {
		return nil, err
	} else {
		return result.(*model.Document), nil
	}
}

func (d documentDo) FirstOrCreate() (*model.Document, error) {
	if result, err := d.DO.FirstOrCreate(); err != nil {
		return nil, err
	} else {
		return result.(*model.Document), nil
	}
}

func (d documentDo) FindByPage(offset int, limit int) (result []*model.Document, count int64, err error) {
	result, err = d.Offset(offset).Limit(limit).Find()
	if err != nil {
		return
	}

	if size := len(result); 0 < limit && 0 < size && size < limit {
		count = int64(size + offset)
		return
	}

	count, err = d.Offset(-1).Limit(-1).Count()
	return
}

func (d documentDo) ScanByPage(result interface{}, offset int, limit int) (count int64, err error) {
	count, err = d.Count()
	if err != nil {
		return
	}

	err = d.Offset(offset).Limit(limit).Scan(result)
	return
}

func (d documentDo) Scan(result interface{}) (err error) {
	return d.DO.Scan(result)
}

func (d documentDo) Delete(models ...*model.Document) (result gen.ResultInfo, err error) {
	return d.DO.Delete(models)
}

func (d *documentDo) withDO(do gen.Dao) *documentDo {
	d.DO = *do.(*gen.DO)
	return d
}

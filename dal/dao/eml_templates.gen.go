// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.

package dao

import (
	"context"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"gorm.io/gorm/schema"

	"gorm.io/gen"
	"gorm.io/gen/field"

	"gorm.io/plugin/dbresolver"

	"email/dal/model"
)

func newEmlTemplate(db *gorm.DB, opts ...gen.DOOption) emlTemplate {
	_emlTemplate := emlTemplate{}

	_emlTemplate.emlTemplateDo.UseDB(db, opts...)
	_emlTemplate.emlTemplateDo.UseModel(&model.EmlTemplate{})

	tableName := _emlTemplate.emlTemplateDo.TableName()
	_emlTemplate.ALL = field.NewAsterisk(tableName)
	_emlTemplate.ID = field.NewInt64(tableName, "id")
	_emlTemplate.Cid = field.NewInt64(tableName, "cid")
	_emlTemplate.Category = field.NewString(tableName, "category")
	_emlTemplate.ContentType = field.NewString(tableName, "content_type")
	_emlTemplate.TemplateName = field.NewString(tableName, "template_name")
	_emlTemplate.Priority = field.NewInt32(tableName, "priority")
	_emlTemplate.MaxRetry = field.NewInt32(tableName, "max_retry")
	_emlTemplate.SubjectZh = field.NewString(tableName, "subject_zh")
	_emlTemplate.SubjectEn = field.NewString(tableName, "subject_en")
	_emlTemplate.FromName = field.NewString(tableName, "from_name")
	_emlTemplate.ContentZh = field.NewString(tableName, "content_zh")
	_emlTemplate.ContentEn = field.NewString(tableName, "content_en")
	_emlTemplate.Description = field.NewString(tableName, "description")
	_emlTemplate.UpdateTime = field.NewTime(tableName, "update_time")
	_emlTemplate.CreateTime = field.NewTime(tableName, "create_time")

	_emlTemplate.fillFieldMap()

	return _emlTemplate
}

type emlTemplate struct {
	emlTemplateDo emlTemplateDo

	ALL          field.Asterisk
	ID           field.Int64
	Cid          field.Int64  // 关联商户 ID
	Category     field.String // 模板类型，目前有：普通（common）、推荐奖励（redemption）
	ContentType  field.String // 模板类型(text/HTML 或 text/plain)
	TemplateName field.String // 模板名称
	Priority     field.Int32  // 优先级（越小越优先）
	MaxRetry     field.Int32  // 最大重试次数
	SubjectZh    field.String // 邮件中文标题
	SubjectEn    field.String // 邮件英文标题
	FromName     field.String // 发送人名称
	ContentZh    field.String // 中文模板
	ContentEn    field.String // 英文模板
	Description  field.String // 模板描述（非必需）
	UpdateTime   field.Time   // 更新时间
	CreateTime   field.Time   // 创建时间

	fieldMap map[string]field.Expr
}

func (e emlTemplate) Table(newTableName string) *emlTemplate {
	e.emlTemplateDo.UseTable(newTableName)
	return e.updateTableName(newTableName)
}

func (e emlTemplate) As(alias string) *emlTemplate {
	e.emlTemplateDo.DO = *(e.emlTemplateDo.As(alias).(*gen.DO))
	return e.updateTableName(alias)
}

func (e *emlTemplate) updateTableName(table string) *emlTemplate {
	e.ALL = field.NewAsterisk(table)
	e.ID = field.NewInt64(table, "id")
	e.Cid = field.NewInt64(table, "cid")
	e.Category = field.NewString(table, "category")
	e.ContentType = field.NewString(table, "content_type")
	e.TemplateName = field.NewString(table, "template_name")
	e.Priority = field.NewInt32(table, "priority")
	e.MaxRetry = field.NewInt32(table, "max_retry")
	e.SubjectZh = field.NewString(table, "subject_zh")
	e.SubjectEn = field.NewString(table, "subject_en")
	e.FromName = field.NewString(table, "from_name")
	e.ContentZh = field.NewString(table, "content_zh")
	e.ContentEn = field.NewString(table, "content_en")
	e.Description = field.NewString(table, "description")
	e.UpdateTime = field.NewTime(table, "update_time")
	e.CreateTime = field.NewTime(table, "create_time")

	e.fillFieldMap()

	return e
}

func (e *emlTemplate) WithContext(ctx context.Context) *emlTemplateDo {
	return e.emlTemplateDo.WithContext(ctx)
}

func (e emlTemplate) TableName() string { return e.emlTemplateDo.TableName() }

func (e emlTemplate) Alias() string { return e.emlTemplateDo.Alias() }

func (e emlTemplate) Columns(cols ...field.Expr) gen.Columns { return e.emlTemplateDo.Columns(cols...) }

func (e *emlTemplate) GetFieldByName(fieldName string) (field.OrderExpr, bool) {
	_f, ok := e.fieldMap[fieldName]
	if !ok || _f == nil {
		return nil, false
	}
	_oe, ok := _f.(field.OrderExpr)
	return _oe, ok
}

func (e *emlTemplate) fillFieldMap() {
	e.fieldMap = make(map[string]field.Expr, 15)
	e.fieldMap["id"] = e.ID
	e.fieldMap["cid"] = e.Cid
	e.fieldMap["category"] = e.Category
	e.fieldMap["content_type"] = e.ContentType
	e.fieldMap["template_name"] = e.TemplateName
	e.fieldMap["priority"] = e.Priority
	e.fieldMap["max_retry"] = e.MaxRetry
	e.fieldMap["subject_zh"] = e.SubjectZh
	e.fieldMap["subject_en"] = e.SubjectEn
	e.fieldMap["from_name"] = e.FromName
	e.fieldMap["content_zh"] = e.ContentZh
	e.fieldMap["content_en"] = e.ContentEn
	e.fieldMap["description"] = e.Description
	e.fieldMap["update_time"] = e.UpdateTime
	e.fieldMap["create_time"] = e.CreateTime
}

func (e emlTemplate) clone(db *gorm.DB) emlTemplate {
	e.emlTemplateDo.ReplaceConnPool(db.Statement.ConnPool)
	return e
}

func (e emlTemplate) replaceDB(db *gorm.DB) emlTemplate {
	e.emlTemplateDo.ReplaceDB(db)
	return e
}

type emlTemplateDo struct{ gen.DO }

func (e emlTemplateDo) Debug() *emlTemplateDo {
	return e.withDO(e.DO.Debug())
}

func (e emlTemplateDo) WithContext(ctx context.Context) *emlTemplateDo {
	return e.withDO(e.DO.WithContext(ctx))
}

func (e emlTemplateDo) ReadDB() *emlTemplateDo {
	return e.Clauses(dbresolver.Read)
}

func (e emlTemplateDo) WriteDB() *emlTemplateDo {
	return e.Clauses(dbresolver.Write)
}

func (e emlTemplateDo) Session(config *gorm.Session) *emlTemplateDo {
	return e.withDO(e.DO.Session(config))
}

func (e emlTemplateDo) Clauses(conds ...clause.Expression) *emlTemplateDo {
	return e.withDO(e.DO.Clauses(conds...))
}

func (e emlTemplateDo) Returning(value interface{}, columns ...string) *emlTemplateDo {
	return e.withDO(e.DO.Returning(value, columns...))
}

func (e emlTemplateDo) Not(conds ...gen.Condition) *emlTemplateDo {
	return e.withDO(e.DO.Not(conds...))
}

func (e emlTemplateDo) Or(conds ...gen.Condition) *emlTemplateDo {
	return e.withDO(e.DO.Or(conds...))
}

func (e emlTemplateDo) Select(conds ...field.Expr) *emlTemplateDo {
	return e.withDO(e.DO.Select(conds...))
}

func (e emlTemplateDo) Where(conds ...gen.Condition) *emlTemplateDo {
	return e.withDO(e.DO.Where(conds...))
}

func (e emlTemplateDo) Order(conds ...field.Expr) *emlTemplateDo {
	return e.withDO(e.DO.Order(conds...))
}

func (e emlTemplateDo) Distinct(cols ...field.Expr) *emlTemplateDo {
	return e.withDO(e.DO.Distinct(cols...))
}

func (e emlTemplateDo) Omit(cols ...field.Expr) *emlTemplateDo {
	return e.withDO(e.DO.Omit(cols...))
}

func (e emlTemplateDo) Join(table schema.Tabler, on ...field.Expr) *emlTemplateDo {
	return e.withDO(e.DO.Join(table, on...))
}

func (e emlTemplateDo) LeftJoin(table schema.Tabler, on ...field.Expr) *emlTemplateDo {
	return e.withDO(e.DO.LeftJoin(table, on...))
}

func (e emlTemplateDo) RightJoin(table schema.Tabler, on ...field.Expr) *emlTemplateDo {
	return e.withDO(e.DO.RightJoin(table, on...))
}

func (e emlTemplateDo) Group(cols ...field.Expr) *emlTemplateDo {
	return e.withDO(e.DO.Group(cols...))
}

func (e emlTemplateDo) Having(conds ...gen.Condition) *emlTemplateDo {
	return e.withDO(e.DO.Having(conds...))
}

func (e emlTemplateDo) Limit(limit int) *emlTemplateDo {
	return e.withDO(e.DO.Limit(limit))
}

func (e emlTemplateDo) Offset(offset int) *emlTemplateDo {
	return e.withDO(e.DO.Offset(offset))
}

func (e emlTemplateDo) Scopes(funcs ...func(gen.Dao) gen.Dao) *emlTemplateDo {
	return e.withDO(e.DO.Scopes(funcs...))
}

func (e emlTemplateDo) Unscoped() *emlTemplateDo {
	return e.withDO(e.DO.Unscoped())
}

func (e emlTemplateDo) Create(values ...*model.EmlTemplate) error {
	if len(values) == 0 {
		return nil
	}
	return e.DO.Create(values)
}

func (e emlTemplateDo) CreateInBatches(values []*model.EmlTemplate, batchSize int) error {
	return e.DO.CreateInBatches(values, batchSize)
}

// Save : !!! underlying implementation is different with GORM
// The method is equivalent to executing the statement: db.Clauses(clause.OnConflict{UpdateAll: true}).Create(values)
func (e emlTemplateDo) Save(values ...*model.EmlTemplate) error {
	if len(values) == 0 {
		return nil
	}
	return e.DO.Save(values)
}

func (e emlTemplateDo) First() (*model.EmlTemplate, error) {
	if result, err := e.DO.First(); err != nil {
		return nil, err
	} else {
		return result.(*model.EmlTemplate), nil
	}
}

func (e emlTemplateDo) Take() (*model.EmlTemplate, error) {
	if result, err := e.DO.Take(); err != nil {
		return nil, err
	} else {
		return result.(*model.EmlTemplate), nil
	}
}

func (e emlTemplateDo) Last() (*model.EmlTemplate, error) {
	if result, err := e.DO.Last(); err != nil {
		return nil, err
	} else {
		return result.(*model.EmlTemplate), nil
	}
}

func (e emlTemplateDo) Find() ([]*model.EmlTemplate, error) {
	result, err := e.DO.Find()
	return result.([]*model.EmlTemplate), err
}

func (e emlTemplateDo) FindInBatch(batchSize int, fc func(tx gen.Dao, batch int) error) (results []*model.EmlTemplate, err error) {
	buf := make([]*model.EmlTemplate, 0, batchSize)
	err = e.DO.FindInBatches(&buf, batchSize, func(tx gen.Dao, batch int) error {
		defer func() { results = append(results, buf...) }()
		return fc(tx, batch)
	})
	return results, err
}

func (e emlTemplateDo) FindInBatches(result *[]*model.EmlTemplate, batchSize int, fc func(tx gen.Dao, batch int) error) error {
	return e.DO.FindInBatches(result, batchSize, fc)
}

func (e emlTemplateDo) Attrs(attrs ...field.AssignExpr) *emlTemplateDo {
	return e.withDO(e.DO.Attrs(attrs...))
}

func (e emlTemplateDo) Assign(attrs ...field.AssignExpr) *emlTemplateDo {
	return e.withDO(e.DO.Assign(attrs...))
}

func (e emlTemplateDo) Joins(fields ...field.RelationField) *emlTemplateDo {
	for _, _f := range fields {
		e = *e.withDO(e.DO.Joins(_f))
	}
	return &e
}

func (e emlTemplateDo) Preload(fields ...field.RelationField) *emlTemplateDo {
	for _, _f := range fields {
		e = *e.withDO(e.DO.Preload(_f))
	}
	return &e
}

func (e emlTemplateDo) FirstOrInit() (*model.EmlTemplate, error) {
	if result, err := e.DO.FirstOrInit(); err != nil {
		return nil, err
	} else {
		return result.(*model.EmlTemplate), nil
	}
}

func (e emlTemplateDo) FirstOrCreate() (*model.EmlTemplate, error) {
	if result, err := e.DO.FirstOrCreate(); err != nil {
		return nil, err
	} else {
		return result.(*model.EmlTemplate), nil
	}
}

func (e emlTemplateDo) FindByPage(offset int, limit int) (result []*model.EmlTemplate, count int64, err error) {
	result, err = e.Offset(offset).Limit(limit).Find()
	if err != nil {
		return
	}

	if size := len(result); 0 < limit && 0 < size && size < limit {
		count = int64(size + offset)
		return
	}

	count, err = e.Offset(-1).Limit(-1).Count()
	return
}

func (e emlTemplateDo) ScanByPage(result interface{}, offset int, limit int) (count int64, err error) {
	count, err = e.Count()
	if err != nil {
		return
	}

	err = e.Offset(offset).Limit(limit).Scan(result)
	return
}

func (e emlTemplateDo) Scan(result interface{}) (err error) {
	return e.DO.Scan(result)
}

func (e emlTemplateDo) Delete(models ...*model.EmlTemplate) (result gen.ResultInfo, err error) {
	return e.DO.Delete(models)
}

func (e *emlTemplateDo) withDO(do gen.Dao) *emlTemplateDo {
	e.DO = *do.(*gen.DO)
	return e
}

package tpls

import (
	"bytes"
	"text/template"
)

const repoTpl = `
package repo

import (
	"context"

	"{{.AppPkgPath}}/entity"
	"{{.AppPkgPath}}/repo/dbal"

	"{{.ProjectName}}/common/tools/filterx"
)

// {{.EntityName}}Repo . @GI
type {{.EntityName}}Repo struct {
	DBAL *dbal.{{.EntityName}}RepoDBAL
}

func New{{.EntityName}}Repo() *{{.EntityName}}Repo {
	return &{{.EntityName}}Repo{
		DBAL: dbal.New{{.EntityName}}RepoDBAL(),
	}
}

func (r *{{.EntityName}}Repo) Query(ctx context.Context, query filterx.FilteringList, pg *filterx.Page) (entity.{{.EntityName}}List, int, error) {
	return r.DBAL.Query(ctx,query,pg)
}

func (r *{{.EntityName}}Repo) Count(ctx context.Context, query filterx.FilteringList) (int64, error) {
	return r.DBAL.Count(ctx, query)
}

func (r *{{.EntityName}}Repo) QueryOne(ctx context.Context, query filterx.FilteringList) (*entity.{{.EntityName}}, error) {
	return r.DBAL.QueryOne(ctx, query)
}

func (r *{{.EntityName}}Repo) Create(ctx context.Context, input *entity.{{.EntityName}}) (*entity.{{.EntityName}}, error) {
	return r.DBAL.Create(ctx,input)
}

func (r *{{.EntityName}}Repo) Save(ctx context.Context, input *entity.{{.EntityName}}) (*entity.{{.EntityName}}, error) {
	return r.DBAL.Save(ctx,input)
}


func (r *{{.EntityName}}Repo) Transaction(ctx context.Context, executeFunc func(tx *gorm.DB) error) error {
	return r.DBAL.Transaction(ctx, executeFunc)
}

func (r *{{.EntityName}}Repo) UpdateByQuery(ctx context.Context, query filterx.FilteringList, updates map[string]any) error {
	return r.DBAL.UpdateByQuery(ctx, query, updates)
}


func (r *{{.EntityName}}Repo) DeleteByQuery(ctx context.Context, query filterx.FilteringList) error {
	return r.DBAL.DeleteByQuery(ctx, query)
}


{{- if .HasID}}
func (r *{{.EntityName}}Repo) GetByID(ctx context.Context, id int64) (*entity.{{.EntityName}}, error) {
	return r.DBAL.GetByID(ctx,id)
}

func (r *{{.EntityName}}Repo) GetListByIDs(ctx context.Context, ids []int64) (entity.{{.EntityName}}List, error) {
	return r.DBAL.GetListByIDs(ctx,ids)
}

func (r *{{.EntityName}}Repo) UpdateByID(ctx context.Context, id int64, updates map[string]any) error {
	return r.DBAL.UpdateByID(ctx,id,updates)
}

func (r *{{.EntityName}}Repo) UpdateByIDs(ctx context.Context, ids []int64, updates map[string]any) error {
	return r.DBAL.UpdateByIDs(ctx,ids,updates)
}

func (r *{{.EntityName}}Repo) DeleteByID(ctx context.Context, id int64) error {
	return r.DBAL.DeleteByID(ctx,id)
}

func (r *{{.EntityName}}Repo) DeleteByIDs(ctx context.Context, ids []int64) error {
	return r.DBAL.DeleteByIDs(ctx,ids)
}
{{- end}}

`

const RepoDBALTpl = `package dbal

import (
	"context"

	"gorm.io/gorm"

	"{{.AppPkgPath}}/config"
	"{{.AppPkgPath}}/entity"
	"{{.AppPkgPath}}/repo/converter"
	"{{.AppPkgPath}}/repo/do"

	"{{.ProjectName}}/common/lib/db"
	"{{.ProjectName}}/common/tools/filterx"
)

// {{.EntityName}}RepoDBAL .
type {{.EntityName}}RepoDBAL struct {
}

func New{{.EntityName}}RepoDBAL() *{{.EntityName}}RepoDBAL {
	return &{{.EntityName}}RepoDBAL{}
}

func (impl *{{.EntityName}}RepoDBAL) NewReadSession(ctx context.Context) *gorm.DB {
	return impl.NewCreateSession(ctx)
}

func (impl *{{.EntityName}}RepoDBAL) NewUpdateSession(ctx context.Context) *gorm.DB {
	return impl.NewCreateSession(ctx)
}

func (impl *{{.EntityName}}RepoDBAL) NewCreateSession(ctx context.Context) *gorm.DB {
	session := config.GetDB().NewSession(ctx)
	session = session.Table(do.TableName{{.EntityName}}Do)
	return session
}

func (impl *{{.EntityName}}RepoDBAL) NewTransactionSession(ctx context.Context) *gorm.DB {
	session := config.GetDB().NewSession(ctx)
	return session
}

func (impl *{{.EntityName}}RepoDBAL) findPage(session *gorm.DB) (do.{{.EntityName}}DoList, int, error) {
	result := make(do.{{.EntityName}}DoList, 0)
	err := session.Find(&result).Error
	if err != nil {
		return nil, 0, errors.Wrapf(err, "{{.EntityName}} FindPage failed 数据库错误")
	}
	delete(session.Statement.Clauses, "LIMIT")
	var count int64
	err = session.Count(&count).Error
	if err != nil {
		return nil, 0, errors.Wrapf(err, "{{.EntityName}} FindPage failed 数据库错误")
	}
	return result, int(count), nil
}

func (impl *{{.EntityName}}RepoDBAL) findAll(session *gorm.DB) (do.{{.EntityName}}DoList, error) {
	result := make(do.{{.EntityName}}DoList, 0)
	err := session.Find(&result).Error
	if err != nil {
		return nil, errors.Wrapf(err, "{{.EntityName}} FindAll failed 数据库错误")
	}
	return result, nil
}

func (impl *{{.EntityName}}RepoDBAL) Query(ctx context.Context, query filterx.FilteringList, pg *filterx.Page) (entity.{{.EntityName}}List, int, error) {
	session := impl.NewReadSession(ctx)
	session, err := query.GormOption(session)
	if err != nil {
		return nil, 0, err
	}
	session, noCount := filterx.PageGormOption(session, pg)
	var (
		doList do.{{.EntityName}}DoList
		count  int
	)
	if noCount {
		doList, err = impl.findAll(session)
	} else {
		doList, count, err = impl.findPage(session)
	}
	if err != nil {
		return nil, 0, err
	}
	return converter.To{{.EntityName}}List(doList), count, nil
}

func (impl *{{.EntityName}}RepoDBAL) Count(ctx context.Context, query filterx.FilteringList) (int64, error) {
	session := impl.NewReadSession(ctx)
	session, err := query.GormOption(session)
	if err != nil {
		return 0, err
	}
	var count int64
	err = session.Count(&count).Error
	if err != nil {
		return 0, errors.Wrapf(err, "{{.EntityName}} Count failed 数据库错误")
	}
	return count, nil
}

func (impl *{{.EntityName}}RepoDBAL) QueryOne(ctx context.Context, query filterx.FilteringList) (*entity.{{.EntityName}}, error) {
	session := impl.NewReadSession(ctx)
	session, err := query.GormOption(session)
	if err != nil {
		return nil, err
	}
	_do := &do.{{.EntityName}}Do{}
	err = session.First(_do).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, errors.Wrapf(err, "{{.EntityName}} QueryOne failed")
	}
	return converter.To{{.EntityName}}Entity(_do), nil
}

func (impl *{{.EntityName}}RepoDBAL) Create(ctx context.Context, input *entity.{{.EntityName}}) (*entity.{{.EntityName}}, error) {
	session := impl.NewCreateSession(ctx)
	_do := converter.From{{.EntityName}}Entity(input)
	err := session.Create(_do).Error
	if err != nil {
		return nil, errors.Wrapf(err, "{{.EntityName}} Create failed")
	}
	output := converter.To{{.EntityName}}Entity(_do)
	return output, err
}

func (impl *{{.EntityName}}RepoDBAL) Save(ctx context.Context, input *entity.{{.EntityName}}) (*entity.{{.EntityName}}, error) {
	session := impl.NewCreateSession(ctx)
	_do := converter.From{{.EntityName}}Entity(input)
	err := session.Save(_do).Error
	if err != nil {
		return nil, errors.Wrapf(err, "{{.EntityName}} Save failed")
	}
	output := converter.To{{.EntityName}}Entity(_do)
	return output, err
}

func (impl *{{.EntityName}}RepoDBAL) Transaction(ctx context.Context, executeFunc func(tx *gorm.DB) error) (err error) {
	// 注意使用的场景（分库分表情况慎用）
	session := impl.NewTransactionSession(ctx)
	err = session.Transaction(executeFunc)
	return err
}

func (impl *{{.EntityName}}RepoDBAL) DemoTransactionWithFunc(ctx context.Context, withFunList []func() error) (err error) {
	// 这是例子 请针对性业务定制 注意使用的场景
	session := impl.NewTransactionSession(ctx)
	err = session.Transaction(func(tx *gorm.DB) error {
		//do something
		for _, fun := range withFunList {
			err = fun()
			if err != nil {
				return err
			}
		}
		//do something
		return nil
	})
	return err
}

func (impl *{{.EntityName}}RepoDBAL) UpdateByQuery(ctx context.Context, query filterx.FilteringList, updates map[string]any) error {
	session := impl.NewUpdateSession(ctx)
	session, err := query.GormOption(session)
	if err != nil {
		return err
	}
	err = session.Updates(updates).Error
	if err != nil {
		return errors.Wrapf(err, "{{.EntityName}} Update failed")
	}
	return err
}

func (impl *{{.EntityName}}RepoDBAL) DeleteByQuery(ctx context.Context, query filterx.FilteringList) error {
	session := impl.NewUpdateSession(ctx)
	session, err := query.GormOption(session)
	if err != nil {
		return err
	}
	err = session.Delete(&do.{{.EntityName}}Do{}).Error
	if  err != nil {
		return errors.Wrapf(err, "{{.EntityName}} Delete failed")
	}
	return err
}

{{- if .HasID}}

func (impl *{{.EntityName}}RepoDBAL) GetByID(ctx context.Context, id int64) (*entity.{{.EntityName}}, error) {
	session := impl.NewReadSession(ctx)
	session = session.Where("id = ?",id)
	_do := &do.{{.EntityName}}Do{}
	err := session.First(_do).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, errors.Wrapf(err, "{{.EntityName}} GetByID failed")
	}
	return converter.To{{.EntityName}}Entity(_do), nil
}

func (impl *{{.EntityName}}RepoDBAL) GetListByIDs(ctx context.Context, ids []int64) (entity.{{.EntityName}}List, error) {
	session := impl.NewReadSession(ctx)
	session = session.Where("id in ?", ids)
	_doList, err := impl.findAll(session)
	if err != nil {
		return nil, err
	}
	return converter.To{{.EntityName}}List(_doList), nil
}

func (impl *{{.EntityName}}RepoDBAL) UpdateByID(ctx context.Context, id int64, updates map[string]any) error {
	session := impl.NewUpdateSession(ctx)
	session = session.Where("id = ?",id)
	err := session.Updates(updates).Error
	if err != nil {
		return errors.Wrapf(err, "{{.EntityName}} UpdateByID failed")
	}
	return err
}

func (impl *{{.EntityName}}RepoDBAL) UpdateByIDs(ctx context.Context, ids []int64, updates map[string]any) error {
	session := impl.NewUpdateSession(ctx)
	session = session.Where("id in ?",ids)
	err := session.Updates(updates).Error
	if err != nil {
		return errors.Wrapf(err, "{{.EntityName}} UpdateByIDs failed")
	}
	return err
}

func (impl *{{.EntityName}}RepoDBAL) DeleteByID(ctx context.Context, id int64) error {
	session := impl.NewUpdateSession(ctx)
	session = session.Where("id = ?",id)
	err := session.Delete(&do.{{.EntityName}}Do{}).Error
	if  err != nil {
		return errors.Wrapf(err, "{{.EntityName}} DeleteByID failed")
	}
	return err
}

func (impl *{{.EntityName}}RepoDBAL) DeleteByIDs(ctx context.Context, ids []int64) error {
	session := impl.NewUpdateSession(ctx)
	session = session.Where("id in ?", ids)
	err := session.Delete(&do.{{.EntityName}}Do{}).Error
	if  err != nil {
		return errors.Wrapf(err, "{{.EntityName}} DeleteByID failed")
	}
	return err
}

{{- end}}
`

type Repo struct {
	ProjectName string
	AppPkgPath  string
	EntityName  string
	TableName   string
	HasID       bool
}

func (s *Repo) Execute() ([]byte, error) {
	buf := new(bytes.Buffer)
	tmpl, err := template.New("repl").Parse(repoTpl)
	if err != nil {
		return nil, err
	}
	if err := tmpl.Execute(buf, s); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func (s *Repo) ExecuteImpl() ([]byte, error) {
	buf := new(bytes.Buffer)
	tmpl, err := template.New("repl.impl").Parse(RepoDBALTpl)
	if err != nil {
		return nil, err
	}
	if err := tmpl.Execute(buf, s); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

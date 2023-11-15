package tpls

import (
	"bytes"
	"text/template"
)

const DaoHeaderCodes = `// Code generated by go3gen. DO NOT EDIT.
package %s

import (
	"github.com/pkg/errors"
	"gorm.io/gorm"
    "%s/repo/dbal/do"
)
`

const daoTpl = `
type {{.DaoName}} struct {
}

func New{{.DaoName}}() *{{.DaoName}} {
  return &{{.DaoName}}{}
}

{{if .PkName}}func (dao *{{.DaoName}}) GetById(session *gorm.DB, id {{.PkType}}) (*do.{{.EntityName}}, error) {
	result := &do.{{.EntityName}}{}
	err := session.Where("{{.PkCol}} = ?", id).First(result).Error
	if err != nil {
		return nil, errors.Wrapf(err, "{{.DaoName}} GetById failed")
	}
	return result, nil
} 
{{end}}

{{if .PkName}}func (dao *{{.DaoName}}) GetByIdList(session *gorm.DB, idList []{{.PkType}}) (do.{{.EntityListName}}, error) {
	result := make([]*do.{{.EntityName}}, 0)
	if err := session.Where("{{.PkCol}} in (?)", idList).Find(&result).Error; err != nil {
		return nil, errors.Wrapf(err, "{{.DaoName}} GetByIdList failed")
	}
	return result, nil
} 
{{end}}

func (dao *{{.DaoName}}) Create(session *gorm.DB, data *do.{{.EntityName}}) error {
	err := session.Create(data).Error
	if err != nil {
		return errors.Wrapf(err, "{{.DaoName}} Create failed")
	}
	return nil
}


func (dao *{{.DaoName}}) Save(session *gorm.DB, data *do.{{.EntityName}}) error {
	err := session.Save(data).Error
	if err != nil {
		return errors.Wrapf(err, "{{.DaoName}} Save failed")
	}
	return nil
}


func (dao *{{.DaoName}}) CreateBatch(session *gorm.DB, data do.{{.EntityListName}}) error {
	err := session.CreateInBatches(data, len(data)).Error
	if err != nil {
		return errors.Wrapf(err, "{{.DaoName}} CreateBatch failed")
	}
	return nil
}

func (dao *{{.DaoName}}) Update(session *gorm.DB,updates map[string]any) error {
	err := session.Updates(updates).Error
	if err != nil {
		return errors.Wrapf(err, "{{.DaoName}} Update failed")
	}
	return nil
}

func (dao *{{.DaoName}}) Delete(session *gorm.DB) error {
	err := session.Delete(&do.{{.EntityName}}{}).Error
	if  err != nil {
		return errors.Wrapf(err, "{{.DaoName}} Delete failed")
	}
	return nil
}

func (dao *{{.DaoName}}) FindPage(session *gorm.DB) (do.{{.EntityListName}}, int, error) {

	result := make([]*do.{{.EntityName}}, 0)
	err := session.Find(&result).Error
	if err != nil {
		return nil, 0, errors.Wrapf(err, "{{.DaoName}} FindPage failed 数据库错误")
	}
	delete(session.Statement.Clauses, "LIMIT")
	var count int64
	err = session.Count(&count).Error
	if err != nil {
		return nil, 0, errors.Wrapf(err, "{{.DaoName}} FindPage failed 数据库错误")
	}
	return result, int(count), nil
}

func (dao *{{.DaoName}}) Count(session *gorm.DB) (int64, error) {
	var count int64
	err := session.Count(&count).Error
	if err != nil {
		return 0, errors.Wrapf(err, "{{.DaoName}} Count failed 数据库错误")
	}
	return count, nil
}


func (dao *{{.DaoName}}) FindAll(session *gorm.DB) (do.{{.EntityListName}}, error) {
	result := make([]*do.{{.EntityName}}, 0)
	err := session.Find(&result).Error
	if err != nil {
		return nil, errors.Wrapf(err, "{{.DaoName}} FindAll failed 数据库错误")
	}
	return result, nil
}

func (dao *{{.DaoName}}) Get(session *gorm.DB) (*do.{{.EntityName}}, error) {
	result := &do.{{.EntityName}}{}
	err := session.First(result).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, errors.New("记录获取失败")
	}
	return result, nil
}`

type Dao struct {
	EntityName     string
	DaoName        string
	EntityListName string
	TableName      string
	PkName         string
	PkType         string
	PkCol          string
}

func (s *Dao) Execute() ([]byte, error) {
	buf := new(bytes.Buffer)
	tmpl, err := template.New("Dao").Parse(daoTpl)
	if err != nil {
		return nil, err
	}
	if err := tmpl.Execute(buf, s); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

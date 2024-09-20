package tpls

import (
	"bytes"
	"text/template"
)

const repoExcelTpl = `
package repo

import (
	"context"

	"github.com/pkg/errors"

	"{{.AppPkgPath}}/entity"
	"{{.ProjectName}}/common/tools"
)

// {{.EntityName}}Repo . @GI
type {{.EntityName}}Repo struct {
	Table string
	mux sync.Mutex
	memCache entity.{{.EntityName}}List
}

func New{{.EntityName}}Repo() *{{.EntityName}}Repo {
	r := &{{.EntityName}}Repo{
		memCache: entity.{{.EntityName}}List{},
		mux: sync.Mutex{},
		Table: "{{.TableName}}",
	}
	r.loadCache()
	go func() {
		ticker := time.NewTicker(time.Second * 5)
		for range ticker.C {
			r.loadCache()
		}
	}()
	return r
}

func (r *{{.EntityName}}Repo) loadCache() {
	r.mux.Lock()
	defer r.mux.Unlock()
	data, err := tool_excel.ExcelSheetData(fmt.Sprintf("%s/%s", config.CCfg.ExcelPath, r.Table), "")
	if err != nil {
		panic(err)
	}
	buf, err := tools.JSONFuzzy.Marshal(data)
	if err != nil {
		panic(err)
	}
	list := entity.{{.EntityName}}List{}
	err = tools.JSONFuzzy.Unmarshal(buf, &list)
	if err != nil {
		panic(err)
	}
	r.memCache = list
}

func (r *{{.EntityName}}Repo) GetCaches() entity.{{.EntityName}}List {
	if r.memCache == nil || len(r.memCache) == 0 {
		r.loadCache()
    }
	return r.memCache
}
`

type RepoEXCEL struct {
	ProjectName string
	AppPkgPath  string
	EntityName  string
	TableName   string
}

func (s *RepoEXCEL) Execute() ([]byte, error) {
	buf := new(bytes.Buffer)
	tmpl, err := template.New("eepoEXCEL").Parse(repoExcelTpl)
	if err != nil {
		return nil, err
	}
	if err := tmpl.Execute(buf, s); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

package tpls

import (
	"bytes"
	"text/template"
)

const repoYAMLTpl = `
package repo

import (
	"fmt"
	"os"
	"sync"

	"gopkg.in/yaml.v3"
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
	r.LoadCache()
	return r
}

func (r *{{.EntityName}}Repo) Filename() string {
	return r.Table
}

func (r *{{.EntityName}}Repo) LoadCache() {
	r.mux.Lock()
	defer r.mux.Unlock()
	buf, err := os.ReadFile(path.Join(config.CCfg.ExcelPath, r.Table))
	list := entity.{{.EntityName}}List{}
	err = yaml.Unmarshal(buf, &list)
	if err != nil {
		panic(err)
	}
	r.memCache = list
}

func (r *{{.EntityName}}Repo) GetCaches() entity.{{.EntityName}}List {
	if r.memCache == nil || len(r.memCache) == 0 {
		r.LoadCache()
    }
	return r.memCache
}
`

type RepoYAML struct {
	ProjectName string
	AppPkgPath  string
	EntityName  string
	TableName   string
}

func (s *RepoYAML) Execute() ([]byte, error) {
	buf := new(bytes.Buffer)
	tmpl, err := template.New("RepoYAML").Parse(repoYAMLTpl)
	if err != nil {
		return nil, err
	}
	if err := tmpl.Execute(buf, s); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

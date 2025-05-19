package tpls

import (
	"bytes"
	"text/template"
)

const repoJSONTpl = `
package repo

import (
	"fmt"
	"os"
	"sync"

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
	buf, err := os.ReadFile(path.Join(config.CCfg.BCfgPath, r.Table))
	list := entity.{{.EntityName}}List{}
	err = tools.JSONFuzzy.Unmarshal(buf, &list)
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

type RepoJSON struct {
	ProjectName string
	AppPkgPath  string
	EntityName  string
	TableName   string
}

func (s *RepoJSON) Execute() ([]byte, error) {
	buf := new(bytes.Buffer)
	tmpl, err := template.New("repoJSON").Parse(repoJSONTpl)
	if err != nil {
		return nil, err
	}
	if err := tmpl.Execute(buf, s); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

package tpls

import (
	"bytes"
	"text/template"
)

const repoBCfgTpl = `
package repo

import (
	"context"

	"github.com/pkg/errors"

	

	"{{.AppPkgPath}}/entity"

	"{{.ProjectName}}/common/tools/filterx"
	"{{.ProjectName}}/common/core/log"
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
	go func() {
		tick := time.NewTicker(5 * time.Second)
		for range tick.C {
			r.loadCache()
		}
	}()
	return r
}

func (r *{{.EntityName}}Repo) loadCache() {
	r.mux.Lock()
	defer r.mux.Unlock()
	list, _, err := r.Query(context.Background(), nil, nil)
	if err != nil {
		return
	}
	r.memCache = list
}

func (r *{{.EntityName}}Repo) GetCaches() entity.{{.EntityName}}List {
	if r.memCache == nil || len(r.memCache) == 0 {
		r.loadCache()
    }
	return r.memCache
}


func (r *{{.EntityName}}Repo) Query(ctx context.Context, query filterx.FilteringList, pg *filterx.Page) (entity.{{.EntityName}}List, int, error) {
	var rs = make(entity.{{.EntityName}}List, 0)
	_tab, err := TableRepoInstance().GetByTable(ctx, r.Table)
	if err != nil {
		return rs, 0, err
	}
	if _tab == nil {
		return rs, 0, errors.New("not found table "+r.Table)
	}
	if pg == nil {
		pg = &filterx.Page{
			Page:     1,
			PageSize: 1000,
			OrderBy:  "sort asc,id asc",
		}
	}
	list, total, err := TableRowsRepoInstance().QueryByTable(ctx, _tab.ID, pg)
	if err != nil {
		return nil, 0, err
	}
	for _, it := range list {
		_row := &entity.{{.EntityName}}{}
		err := tools.JSON.UnmarshalFromString(it.Row, &_row)
		if err != nil {
			log.Errorf("%s Unmarshal err: %v, id: %d",r.Table, err, it.ID)
			continue
		}
		_row.ID = it.ID
		_row.Sort = it.Sort
		rs = append(rs, _row)
	}
	return rs, total, nil
}


func (r *{{.EntityName}}Repo) Create(ctx context.Context, input *entity.{{.EntityName}}) (*entity.{{.EntityName}}, error) {
	_tab, err := TableRepoInstance().GetByTable(ctx, r.Table)
	if err != nil {
		return nil, err
	}
	if _tab == nil {
		return nil, errors.New("not found table "+r.Table)
	}
	rs, err := TableRowsRepoInstance().Create(ctx, &entity.TableRows{
		TID:    _tab.ID,
		Row:    input.String(),
		Sort:   input.Sort,
		Status: 1,
	})
	if err != nil {
		return nil, err
	}
	input.ID = rs.ID
	return input, nil
}

func (r *{{.EntityName}}Repo) GetByID(ctx context.Context, id int64) (*entity.{{.EntityName}}, error) {
	rs, err := TableRowsRepoInstance().GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if rs == nil {
		return nil, errors.Errorf("not found table %s, id: %d", r.Table, id)
	}
	_row := &entity.{{.EntityName}}{}
	err = tools.JSON.UnmarshalFromString(rs.Row, &_row)
	if err != nil {
		return nil, errors.Wrapf(err, "%s Unmarshal id: %d", r.Table, rs.ID)
	}
	_row.ID = rs.ID
	_row.Sort = rs.Sort
	return _row, nil
}

func (r *{{.EntityName}}Repo) Save(ctx context.Context, input *entity.{{.EntityName}}) (*entity.{{.EntityName}}, error) {
	rs, err := TableRowsRepoInstance().GetByID(ctx, input.ID)
	if err != nil {
		return nil, err
	}
	if rs == nil {
		return nil, errors.Errorf("not found table %s id: %d", r.Table, input.ID)
	}
	rs.Sort = input.Sort
	rs.Row, _ = tools.JSON.MarshalToString(input)
	_, err = TableRowsRepoInstance().Save(ctx, rs)
	if err != nil {
		return nil, err
	}
	return input, nil
}


func (r *{{.EntityName}}Repo) DeleteByID(ctx context.Context, id int64) error {
	return TableRowsRepoInstance().DeleteByID(ctx, id)
}
`

type RepoBCFG struct {
	ProjectName string
	AppPkgPath  string
	EntityName  string
	TableName   string
}

func (s *RepoBCFG) Execute() ([]byte, error) {
	buf := new(bytes.Buffer)
	tmpl, err := template.New("repoBCfg").Parse(repoBCfgTpl)
	if err != nil {
		return nil, err
	}
	if err := tmpl.Execute(buf, s); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

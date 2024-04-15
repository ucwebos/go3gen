package tpls

import (
	"bytes"
	"text/template"
)

const AdminAPIItemTpl = `package {{.PkgName}}

import (
	"github.com/gin-gonic/gin"

	"{{.Project}}/panel/types"
	
	"{{.Project}}/micro/{{.AppName}}/entity"
	"{{.Project}}/micro/{{.AppName}}/repo"
)


type {{.Name}}ListReq struct {
	Page *types.Page ` + "`json:\"page\"`" + `
	Form types.Form  ` + "`json:\"form\"`" + `
	Sort *types.Sort ` + "`json:\"sort\"`" + `
}

type {{.Name}}ListResp struct {
	CurrentPage int             ` + "`json:\"currentPage\"`" + `
	PageSize    int             ` + "`json:\"pageSize\"`" + `
	Total       int             ` + "`json:\"total\"`" + `
	Records     []*entity.{{.Name}} ` + "`json:\"records\"`" + `
	Form        types.Form      ` + "`json:\"form\"`" + `
	Sort        *types.Sort     ` + "`json:\"sort\"`" + `
}

func {{.Name}}List(ctx *gin.Context) {
	var (
		req  = &{{.Name}}ListReq{}
		resp = &{{.Name}}ListResp{}
	)
	if err := types.BindBody(ctx, &req); err != nil {
		types.JSONError(ctx, types.ErrorParams, err.Error())
		return
	}

	list, total, err := repo.{{.Name}}RepoInstance().Query(ctx, req.Form.ToFilteringList(), nil)
	if err != nil {
		types.JSONError(ctx, types.ErrorSys, err.Error())
		return
	}
	resp.Records = list
	resp.Total = total
	resp.CurrentPage = req.Page.CurrentPage
	resp.PageSize = req.Page.PageSize
	types.JSONSuccess(ctx, resp)
}

type {{.Name}}AddReq map[string]any

func (req {{.Name}}AddReq) ToEntity() *entity.{{.Name}} {
	out := &entity.{{.Name}}{}
	out.FromMap(req)
	return out
}

type {{.Name}}AddResp struct {
	Target *entity.{{.Name}} ` + "`json:\"target\"`" + `
}

func {{.Name}}Add(ctx *gin.Context) {
	var (
		req  = &{{.Name}}AddReq{}
		resp = &{{.Name}}AddResp{}
	)
	if err := types.BindBody(ctx, &req); err != nil {
		types.JSONError(ctx, types.ErrorParams, err.Error())
		return
	}

	rs, err := repo.{{.Name}}RepoInstance().Create(ctx, req.ToEntity())
	if err != nil {
		types.JSONError(ctx, types.ErrorSys, err.Error())
		return
	}
	resp.Target = rs
	types.JSONSuccess(ctx, resp)
}

type {{.Name}}EditReq map[string]any

func (req {{.Name}}EditReq) ToEntity() *entity.{{.Name}} {
	out := &entity.{{.Name}}{}
	out.FromMap(req)
	return out
}

type {{.Name}}EditResp struct {
	Target *entity.{{.Name}} ` + "`json:\"target\"`" + `
}

func {{.Name}}Edit(ctx *gin.Context) {
	var (
		req  = &{{.Name}}EditReq{}
		resp = &{{.Name}}EditResp{}
	)
	if err := types.BindBody(ctx, &req); err != nil {
		types.JSONError(ctx, types.ErrorParams, err.Error())
		return
	}
	rs := req.ToEntity()
	old, err := repo.{{.Name}}RepoInstance().GetByID(ctx, rs.ID)
	if err != nil {
		types.JSONError(ctx, types.ErrorSys, err.Error())
		return
	}
	if old == nil {
		types.JSONError(ctx, types.ErrorParams, "not found")
		return
	}
	rs, err = repo.{{.Name}}RepoInstance().Save(ctx, rs)
	if err != nil {
		types.JSONError(ctx, types.ErrorSys, err.Error())
		return
	}
	resp.Target = rs
	types.JSONSuccess(ctx, resp)
}

type {{.Name}}DeleteReq struct {
	ID int64 ` + "`json:\"id\"`" + `
}

type {{.Name}}DeleteResp struct {
}

func {{.Name}}Delete(ctx *gin.Context) {
	var (
		req  = &{{.Name}}DeleteReq{}
		resp = &{{.Name}}DeleteResp{}
	)
	if err := types.BindBody(ctx, &req); err != nil {
		types.JSONError(ctx, types.ErrorParams, err.Error())
		return
	}
	err := repo.{{.Name}}RepoInstance().DeleteByID(ctx, req.ID)
	if err != nil {
		types.JSONError(ctx, types.ErrorSys, err.Error())
		return
	}
	types.JSONSuccess(ctx, resp)
}`

const adminConvIoTpl = `
func From{{.NameMark}}{{.Name}}Entity(input *entity.{{.Name}}) *types.{{.Name}}{
	if input == nil {
		return nil
	}
	output := &types.{{.Name}}{}
{{- range .Fields }}
	{{- if eq .SType 1}}
	output.{{.Name}} = From{{$.NameMark}}{{.Type2}}Entity(input.{{.Name}})
	{{- else if eq .SType 2}}
	if input.{{.Name}} != nil {
		{{- if .Type2Entity}}
		output.{{.Name}} = From{{$.NameMark}}{{.Type2}}List(input.{{.Name}})
		{{- else}}
		output.{{.Name}} = input.{{.Name}}
		{{- end}}
	}
	{{- else if eq .SType 4}}
		if !input.{{.Name}}.IsZero() {
			output.{{.Name}} = tool_time.TimeToDateTimeString(input.{{.Name}})
		}
	{{- else}}
	output.{{.Name}} = input.{{.Name}}
	{{- end}}
{{- end}}
	return output
}

func To{{.NameMark}}{{.Name}}Entity(input *types.{{.Name}}) *entity.{{.Name}}{
	if input == nil {
		return nil
	}
	output := &entity.{{.Name}}{}
{{- range .Fields }}
	{{- if eq .SType 1}} 
	output.{{.Name}} = To{{$.NameMark}}{{.Type2}}Entity(input.{{.Name}})
	{{- else if eq .SType 2}}
		{{- if .Type2Entity}}
		output.{{.Name}} = To{{$.NameMark}}{{.Type2}}List(input.{{.Name}})
		{{- else}}
		output.{{.Name}} = input.{{.Name}}
		{{- end}}
	{{- else if eq .SType 4}}
		if ts := tool_time.ParseDateTime(input.{{.Name}}); !ts.IsZero() {
			output.{{.Name}} = ts
		}
	{{- else}}
	output.{{.Name}} = input.{{.Name}}
	{{- end}}
{{- end}}
	return output
}

func From{{.NameMark}}{{.Name}}List(input entity.{{.Name}}List) []*types.{{.Name}} {
	if input == nil {
		return nil
	}
	output := make([]*types.{{.Name}}, 0, len(input))
	for _, item := range input {
		resultItem := From{{.NameMark}}{{.Name}}Entity(item)
		output = append(output, resultItem)
	}
	return output
}

func To{{.NameMark}}{{.Name}}List(input []*types.{{.Name}}) entity.{{.Name}}List {
	if input == nil || len(input) == 0 {
		return nil
	}
	output := make(entity.{{.Name}}List, 0, len(input))
	for _, item := range input {
		resultItem := To{{.NameMark}}{{.Name}}Entity(item)
		output = append(output, resultItem)
	}
	return output
}

`

type AdminAPIItem struct {
	Project string
	AppName string
	PkgName string
	Name    string
	NameVal string
}

func (s *AdminAPIItem) Execute() ([]byte, error) {
	buf := new(bytes.Buffer)
	tmpl, err := template.New("AdminAPIItem").Parse(AdminAPIItemTpl)
	if err != nil {
		return nil, err
	}
	if err := tmpl.Execute(buf, s); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

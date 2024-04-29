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

	list, total, err := repo.{{.Name}}RepoInstance().Query(ctx, req.Form.ToFilteringList(), &filterx.Page{
		Page:     int32(req.Page.CurrentPage),
		PageSize: int32(req.Page.PageSize),
	})
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

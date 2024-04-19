package tpls

import (
	"bytes"
	"text/template"
)

const AdminAPIRouteTpl = `package micro

import (
	"github.com/gin-gonic/gin"

{{- range .Groups}}
	ms{{.Name}} "{{$.Project}}/micro/{{.Name}}"
{{- end}}

{{- range .Groups}}
	"{{$.Project}}/panel/micro/{{.Name}}"
{{- end}}
)

func Init() {
{{- range .Groups}}
	ms{{.Name}}.Init()
{{- end}}
}

func Route(r gin.IRoutes) {
{{- range $g := .Groups}}
	// ----------------- {{.Name}} ---------------
	{{- range $g.CrudList}}
	// {{.NameVal}}
	r.POST("/micro/{{$g.Name}}/{{.NameVal}}-list", {{$g.Name}}.{{.Name}}List)
	r.POST("/micro/{{$g.Name}}/{{.NameVal}}-add", {{$g.Name}}.{{.Name}}Add)
	r.POST("/micro/{{$g.Name}}/{{.NameVal}}-edit", {{$g.Name}}.{{.Name}}Edit)
	r.POST("/micro/{{$g.Name}}/{{.NameVal}}-delete", {{$g.Name}}.{{.Name}}Delete)
	{{- end}}
{{- end}}
}`

type AdminAPIRoute struct {
	Project string
	Groups  []*AdminGroup
}

type AdminGroup struct {
	Name     string
	NameVal  string
	CrudList []CrudItem
}

type CrudItem struct {
	Group   string
	Name    string
	NameVal string
}

func (s *AdminAPIRoute) Execute() ([]byte, error) {
	buf := new(bytes.Buffer)
	tmpl, err := template.New("AdminAPIRoute").Parse(AdminAPIRouteTpl)
	if err != nil {
		return nil, err
	}
	if err := tmpl.Execute(buf, s); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

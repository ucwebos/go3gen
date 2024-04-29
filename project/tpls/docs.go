package tpls

import (
	"bytes"
	"text/template"
)

const docsItemTpl = `# {{.Name}}

> {{.RoutePath}}

## 参数

> [c#] {{.RequestClass}}.cs 

| 字段     | 类型     | 是否必填 | 含义  |
|--------|--------|------|-----|
{{- range .Request}}
| {{.Name}} | {{.Type}} | {{.Must}} | {{.Comment}} |
{{- end}}

## 响应

> [c#] {{.ResponseClass}}.cs 

| 字段          | 类型     | 含义     |
|-------------|--------|--------|
{{- range .Response}}
| {{.Name}} | {{.Type}} | {{.Comment}} |
{{- end}}

## 响应例子
`

type DocsItem struct {
	Name          string
	RoutePath     string
	RequestClass  string
	ResponseClass string
	Request       []DocsItemField
	Response      []DocsItemField
	ExpJSON       []byte
}

type DocsItemField struct {
	Name    string
	Type    string
	Must    string
	Comment string
}

func (s *DocsItem) Execute() ([]byte, error) {
	buf := new(bytes.Buffer)
	tmpl, err := template.New("DocsItem").Parse(docsItemTpl)
	if err != nil {
		return nil, err
	}
	if err := tmpl.Execute(buf, s); err != nil {
		return nil, err
	}
	return append(buf.Bytes(), s.ExpJSON...), nil

}

const docsSidebarTpl = `
{{range $it := .Groups}}
* {{$it.GroupName}}
{{- range $v:=$it.FunList}}
* * [{{$v.FunMark}}]({{$v.URI2}})
{{- end}}
{{end}}
`

type DocsSidebar struct {
	Entry  string
	Groups []*EntryGroup
}

func (s *DocsSidebar) Execute() ([]byte, error) {
	buf := new(bytes.Buffer)
	tmpl, err := template.New("DocsSidebar").Parse(docsSidebarTpl)
	if err != nil {
		return nil, err
	}
	if err := tmpl.Execute(buf, s); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil

}

package tpls

import (
	"bytes"
	"text/template"
)

const HttpTypesTpl = `
package types
{{range .FunList}}
// {{.ReqName}} {{.FunMark}}参数 
type {{.ReqName}} struct {}
// {{.RespName}} {{.FunMark}}响应
type {{.RespName}} struct {}
{{end}}
`

const HttpTypesAppendTpl = `
{{range .FunList}}
// {{.ReqName}} {{.FunMark}}参数 
type {{.ReqName}} struct {}
// {{.RespName}} {{.FunMark}}响应
type {{.RespName}} struct {}
{{end}}
`

type HandlerTypes struct {
	EntryPath string
	Entry     string
	Group     string
	FunList   []EntryFunItem
}

type HandlerTypesAppend struct {
	Body    []byte
	FunList []EntryFunItem
}

func (s *HandlerTypes) Execute() ([]byte, error) {
	buf := new(bytes.Buffer)
	tmpl, err := template.New("HandlerTypes").Parse(HttpTypesTpl)
	if err != nil {
		return nil, err
	}
	if err := tmpl.Execute(buf, s); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func (s *HandlerTypesAppend) Execute() ([]byte, error) {
	buf := new(bytes.Buffer)
	tmpl, err := template.New("HandlerTypesAppend").Parse(HttpTypesAppendTpl)
	if err != nil {
		return nil, err
	}
	if err := tmpl.Execute(buf, s); err != nil {
		return nil, err
	}
	return append(s.Body, buf.Bytes()...), nil
}

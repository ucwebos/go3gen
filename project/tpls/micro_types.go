package tpls

import (
	"bytes"
	"text/template"
)

const microTypesTpl = `
package types_{{.AppName}}

import (
	"time"
)

{{- range .FunList}}

// {{.ReqName}} .
type {{.ReqName}} struct {}

// {{.RespName}} .
type {{.RespName}} struct {}

{{- end}}
`

const microTypesAppendTpl = `
{{- range .FunList}}

// {{.ReqName}} .
type {{.ReqName}} struct {}

// {{.RespName}} .
type {{.RespName}} struct {}

{{- end}}
`

type MicroTypesAppend struct {
	Body    []byte
	FunList []MicroFunItem
}

func (s *MicroTypesAppend) Execute() ([]byte, error) {
	buf := new(bytes.Buffer)
	tmpl, err := template.New("microTypesAppend").Parse(microTypesAppendTpl)
	if err != nil {
		return nil, err
	}
	if err := tmpl.Execute(buf, s); err != nil {
		return nil, err
	}
	return append(s.Body, buf.Bytes()...), nil
}

type MicroTypes struct {
	AppName    string
	AppPkgPath string
	FunList    []MicroFunItem
}

func (s *MicroTypes) Execute() ([]byte, error) {
	buf := new(bytes.Buffer)
	tmpl, err := template.New("microTypes").Parse(microTypesTpl)
	if err != nil {
		return nil, err
	}
	if err := tmpl.Execute(buf, s); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

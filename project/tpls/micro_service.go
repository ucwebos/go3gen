package tpls

import (
	"bytes"
	"text/template"
)

const MServiceFuncInitTpl = `package service

import (
	"context"

	"{{.AppPkgPath}}/types_{{.AppName}}"
)

// {{.Service}} @GI
type {{.Service}} struct {
}

func New{{.Service}}() *{{.Service}} {
	return &{{.Service}}{}
}

{{range .FunList}}
// {{.Method}} {{.FunMark}} 
func (s *{{.Service}}){{.Method}}(ctx context.Context, input *types_{{$.AppName}}.{{.ReqName}}) (*types_{{$.AppName}}.{{.RespName}}, error) {
	var (
		output = &types_{{$.AppName}}.{{.RespName}}{}
	)

	// todo ...
	return output, nil
}
{{end}}
`

type MServiceFunc struct {
	AppName    string
	AppPkgPath string
	Service    string
	FunList    []MicroFunItem
}

const MServiceFuncAppendTpl = `
{{range .FunList}}
func (s *{{.Service}}){{.Method}}(ctx context.Context, input *types_{{$.AppName}}.{{.ReqName}}) (*types_{{$.AppName}}.{{.RespName}}, error) {
	var (
		output = &types_{{$.AppName}}.{{.RespName}}{}
	)

	// todo ...
	return output, nil
}
{{end}}
`

type MServiceFuncAppend struct {
	Body    []byte
	AppName string
	FunList []MicroFunItem
}

func (s *MServiceFuncAppend) Execute() ([]byte, error) {
	buf := new(bytes.Buffer)
	tmpl, err := template.New("MServiceFuncAppend").Parse(MServiceFuncAppendTpl)
	if err != nil {
		return nil, err
	}
	if err := tmpl.Execute(buf, s); err != nil {
		return nil, err
	}
	return append(s.Body, buf.Bytes()...), nil
}

func (s *MServiceFunc) Execute() ([]byte, error) {
	buf := new(bytes.Buffer)
	tmpl, err := template.New("MServiceFunc").Parse(MServiceFuncInitTpl)
	if err != nil {
		return nil, err
	}
	if err := tmpl.Execute(buf, s); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil

}

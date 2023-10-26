package tpls

import (
	"bytes"
	"text/template"
)

const HandlerFuncInitTpl = `package handler

import (
	"context"

	"{{.EntryPath}}/types"
)
{{range .FunList}}
func {{.FunName}}(ctx context.Context, req *types.{{.ReqName}}) (*types.{{.RespName}}, error) {
	var (
		resp = &types.{{.RespName}}{}
	)
	// todo ...
	return resp, nil
}
{{end}}
`

type HandlerFunc struct {
	EntryPath string
	Entry     string
	Group     string
	FunList   []EntryFunItem
}

const HandlerFuncAppendTpl = `
{{range .FunList}}
func {{.FunName}}(ctx context.Context, req *types.{{.ReqName}}) (*types.{{.RespName}}, error) {
	var (
		resp = &types.{{.RespName}}{}
	)
	// todo ...
	return resp, nil
}
{{end}}
`

type HandlerFuncAppend struct {
	Body    []byte
	FunList []EntryFunItem
}

func (s *HandlerFuncAppend) Execute() ([]byte, error) {
	buf := new(bytes.Buffer)
	tmpl, err := template.New("HandlerFuncAppend").Parse(HandlerFuncAppendTpl)
	if err != nil {
		return nil, err
	}
	if err := tmpl.Execute(buf, s); err != nil {
		return nil, err
	}
	return append(s.Body, buf.Bytes()...), nil
}

func (s *HandlerFunc) Execute() ([]byte, error) {
	buf := new(bytes.Buffer)
	tmpl, err := template.New("HandlerFunc").Parse(HandlerFuncInitTpl)
	if err != nil {
		return nil, err
	}
	if err := tmpl.Execute(buf, s); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil

}

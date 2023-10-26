package tpls

import (
	"bytes"
	"text/template"
)

const MicroTestingTpl = `package service

import (
	"context"
	"testing"

	"{{.AppPkgPath}}/types_{{.AppName}}"
)

{{range .FunList}}

func Test{{.Service}}_{{.Method}}(t *testing.T) {
	// 一些依赖
	// ...

	ctx := context.Background()
	//ctx = context.WithValue(ctx, common.UIDKey, int64(1))
	resp, err := {{.Service}}Instance.{{.Method}}(ctx, &types_{{$.AppName}}.{{.ReqName}}{
		// ...
	})
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("resp: %+v \n", *resp)
}

func Benchmark{{.Service}}_{{.Method}}(b *testing.B)  {
	// todo...
}

{{end}}`

type MicroTesting struct {
	AppName    string
	AppNameUF  string
	AppPkgPath string
	Service    string
	FunList    []MicroFunItem
}

func (s *MicroTesting) Execute() ([]byte, error) {
	buf := new(bytes.Buffer)
	tmpl, err := template.New("microEntry").Parse(MicroTestingTpl)
	if err != nil {
		return nil, err
	}
	if err := tmpl.Execute(buf, s); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

package tpls

import (
	"bytes"
	"text/template"
)

const MicroEntryTpl = `// Code generated by go3gen. DO NOT EDIT.
package provider

import (
	"context"

	"{{.AppPkgPath}}/types_{{.AppName}}"
	
)

type {{.AppNameUF}}Interface interface {
	{{- range .FunList}}
		// {{.FunName}} {{.FunMark}} 
		{{.FunName}}(ctx context.Context, input *types_{{$.AppName}}.{{.ReqName}}) (*types_{{$.AppName}}.{{.RespName}},error)
	{{- end}}
}

// Register{{.AppNameUF}}
func Register{{.AppNameUF}}({{.AppName}} {{.AppNameUF}}Interface) {
	{{.AppNameUF}} = {{.AppName}}
}

`

const MicroServiceTpl = `// Code generated by go3gen. DO NOT EDIT.
package {{.AppName}}

import (
	"context"
	"time"

	"{{.Project}}/common"
	"{{.Project}}/common/core/log"
	"{{.Project}}/common/core/otel/prometheus"
	"{{.Project}}/common/core/otel/tracing"

	"{{.AppPkgPath}}/types_{{.AppName}}"
	"{{.AppPkgPath}}/service"
	
)

type {{.AppNameUF}} struct {
}

func New{{.AppNameUF}}() *{{.AppNameUF}} {
	return &{{.AppNameUF}}{}
}

{{- range .FunList}}

// {{.FunName}} {{.FunMark}} 
func (s *{{$.AppNameUF}}){{.FunName}}(ctx context.Context, input *types_{{$.AppName}}.{{.ReqName}}) (*types_{{$.AppName}}.{{.RespName}},error) {
	var (
		st = time.Now()
		resp = &types_{{$.AppName}}.{{.RespName}}{}
		err error
	)
	ctx, span := tracing.StartSpan(ctx, "micro:{{$.AppName}}_{{.FunName}}")
	defer func() {
		span.End()
		prometheus.HistogramVec.Timing("micro_seconds", map[string]string{
			"micro":  "{{$.AppName}}",
			"func":   "{{.FunName}}",
			"ret":    prometheus.RetLabel(err),
		}, st)
		//log.With().TraceID(ctx).Field("req", input).Field("resp", resp).Field("err", err).Debug("on-micro")
	}()
	
	resp, err = service.{{.Service}}Instance().{{.Method}}(ctx,input)
	return resp, err
}	

{{- end}}
`

type MicroEntry struct {
	Project    string
	AppName    string
	AppNameUF  string
	AppPkgPath string
	FunList    []MicroFunItem
}

type MicroFunItem struct {
	Service  string
	Method   string
	FunName  string
	FunMark  string
	ReqName  string
	RespName string
}

func (s *MicroEntry) Execute(t string) ([]byte, error) {
	buf := new(bytes.Buffer)
	tmpl, err := template.New("microEntry").Parse(t)
	if err != nil {
		return nil, err
	}
	if err := tmpl.Execute(buf, s); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

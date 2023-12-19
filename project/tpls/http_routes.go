package tpls

import (
	"bytes"
	"text/template"
)

const HttpRouteTpl = `// Code generated by go3gen. DO NOT EDIT.
package route

import (
	"time"

	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/otel/attribute"

	"{{.Project}}/common"
	"{{.Project}}/common/core/log"
	"{{.Project}}/common/core/otel/prometheus"
	"{{.Project}}/common/core/otel/tracing"


	"{{.AppPkgPath}}/cmd/{{.EntryName}}/handler"
	"{{.AppPkgPath}}/cmd/{{.EntryName}}/middleware"
	"{{.AppPkgPath}}/cmd/{{.EntryName}}/types"
)

func generated(r gin.IRoutes) {

	{{- range .Groups}}
	// ----------------------------------- {{.GroupName}} -----------------------------------
	{{- range $it := .FunList}}
	// {{.FunMark}}
	r.POST("{{.URI}}"{{with $it.Middlewares}},{{range $it.Middlewares}}middleware.{{.}}{{end}}{{end}}, func(ctx *gin.Context) {
		var (
			st = time.Now()
			_ctx = common.HTTPMetadata(ctx)
			req = &types.{{$it.ReqName}}{}
			resp = &types.{{$it.RespName}}{}
			err error
		)
		reqRaw,raw,_common, err := reqJSON(ctx)
		if err != nil {
			JSONError(ctx, common.ErrParams)
			return
		}
		if _common != nil {
			_ctx = context.WithValue(_ctx,common.IOCommonParamsKey,_common)
		}
		_ctx, span := tracing.StartSpan(_ctx, "http:{{$it.URI}}")
		defer func() {
			span.End()
			prometheus.HistogramVec.Timing("http_seconds", map[string]string{
				"entry":  "{{$.EntryName}}",
				"api":    "{{$it.URI}}",
				"ret":    prometheus.RetLabel(err),
			}, st)
			_resp,_ := tools.JSON.Marshal(resp)
			log.With().TraceID(_ctx).Field("common",_common).Field("uri", "{{$it.URI}}").Field("req", raw).Field("resp", _resp).Field("err", err).Info("on-http")
		}()
		if err = common.BindBody(reqRaw, &req); err != nil {
			JSONError(ctx, common.ErrParams)
			return
		}
		resp, err = handler.{{$it.FunName}}(_ctx, req)
		JSON(ctx, resp, err)
	})
	{{- end}}
	{{- end}}
}
`

type HttpEntry struct {
	Project      string
	AppName      string
	AppNameUF    string
	AppPkgPath   string
	EntryName    string
	EntryPath    string
	EntryPkgPath string
	Groups       []*EntryGroup
}

type EntryGroup struct {
	Group        string
	GroupUFirst  string
	GroupName    string
	GMiddlewares []string
	FunList      []EntryFunItem
}

type EntryFunItem struct {
	FunName     string
	Method      string
	FunMark     string
	ReqName     string
	RespName    string
	Middlewares []string
	URI         string
	URI2        string
}

func (s *HttpEntry) Execute(t string) ([]byte, error) {
	buf := new(bytes.Buffer)
	tmpl, err := template.New("httpEntry").Parse(t)
	if err != nil {
		return nil, err
	}
	if err := tmpl.Execute(buf, s); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

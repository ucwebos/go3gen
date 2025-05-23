package tpls

const SocketRouteTpl = `// Code generated by go3gen. DO NOT EDIT.
package route

import (
	"context"
	"time"
	"strings"

	"github.com/topfreegames/pitaya/v2"
	"github.com/topfreegames/pitaya/v2/component"
	"go.opentelemetry.io/otel/attribute"

	"{{.Project}}/common"
	"{{.Project}}/common/core/log"
	"{{.Project}}/common/core/otel/prometheus"
	"{{.Project}}/common/core/otel/tracing"


	"{{.EntryPkgPath}}/handler"
	"{{.EntryPkgPath}}/middleware"
	"{{.EntryPkgPath}}/types"
)

{{- range $x := .Groups}}
// ----------------------------------- {{.GroupName}} -----------------------------------
type {{$.SocketTypeUF}}{{$x.GroupUFirst}} struct {
	component.Base
	app pitaya.Pitaya
}

{{- range $it := .FunList}}
	// {{.FunMark}}
	func (m *{{$.SocketTypeUF}}{{$x.GroupUFirst}}) {{$it.Method}}(ctx context.Context, req *types.{{$it.ReqName}}) (*types.{{$it.RespName}}, error) {
		var (
			st = time.Now()
			sess = m.app.GetSessionFromCtx(ctx)
			resp = &types.{{$it.RespName}}{}
			err error
		)
		// 前置处理
		release, _common, err := middleware.SocBefore(ctx,"{{$.SocketType}}.{{$it.URI}}",sess);
		if err != nil {
			if e, ok := err.(common.ErrCode); ok {
				return nil, pitaya.Error(e,fmt.Sprintf("PI-%d",e.Code),map[string]string{})
				switch e.Code {
				case common.ErrClosureIPCode, common.ErrClosureUserCode, common.ErrMaintainingCode:
					sess.Kick(context.Background())
				}
			}
			return nil, err
		}
		defer release()
		ctx = context.WithValue(ctx, common.IOCommonParamsKey, _common)
		ctx = context.WithValue(ctx, common.IOSessionKey, sess)
		ctx, span := tracing.StartSpan(ctx, "socket:{{$.SocketType}}.{{$it.URI}}")
		defer func() {
			span.End()
			prometheus.HistogramVec.Timing("socket_seconds", map[string]string{
				"entry": "{{$.EntryName}}",
				"route": "{{$.SocketType}}.{{$it.URI}}",
				"ret":   prometheus.RetLabel(err),
			}, st)
			_resp, _ := tools.JSON.MarshalToString(resp)
			_req,_ := tools.JSON.MarshalToString(req)
			log.With().TraceID(ctx).Int("uid",int(_common.UID)).String("uri", "{{$it.URI}}").String("req", _req).String("resp", _resp).Info("ioReply")
		}()
		resp, err = handler.{{$it.FunName}}(ctx, req)
		if err != nil {
			if e, ok := err.(common.ErrCode); ok {
				return nil, pitaya.Error(e, fmt.Sprintf("PI-%d", e.Code), map[string]string{})
			}
		}
		return resp, err
	}
{{- end}}
{{- end}}

func {{.SocketType}}Generated(app pitaya.Pitaya) {
{{- range .Groups}}
	app.Register(
		&{{$.SocketTypeUF}}{{.GroupUFirst}}{app: app}, component.WithName("{{.Group}}"), component.WithNameFunc(strings.ToLower),
	)
{{- end}}
}
`

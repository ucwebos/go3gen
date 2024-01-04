package tpls

import (
	"bytes"
	"text/template"
)

const AdminRouteTpl = `export default [
  {
    title: "{{.Title}}",
    name: "{{.Name}}",
    meta: {
      fixedAside: true,
      showOnHeader: false,
      icon: "{{.Icon}}"
    },
    path: "/micro/{{.NameVal}}",
    redirect: "/micro/{{.NameVal}}",
    children: [
  {{- range .CrudList}}
      {
        title: "{{.Title}}",
        name: "{{$.Name}}{{.Name}}",
        path: "/micro/{{$.NameVal}}/{{.NameVal}}",
        meta: {
          icon: "{{.Icon}}"
        },
        component: "/micro/{{$.NameVal}}/{{.NameVal}}/index.vue"
      },
  {{- end}}
    ]
  }
];
`

type AdminRoutes struct {
	Groups []AdminGroup
}

type AdminGroup struct {
	Name     string
	NameVal  string
	Title    string
	Icon     string
	CrudList []CrudItem
}

type CrudItem struct {
	ApiBaseURL string
	Group      string
	Name       string
	NameVal    string
	Title      string
	Icon       string
	Fields     []CrudItemField
}

type CrudItemField struct {
	Name  string `json:"name"`
	Title string `json:"title"`
}

func (s *AdminGroup) Execute(t string) ([]byte, error) {
	buf := new(bytes.Buffer)
	tmpl, err := template.New("AdminRoutes").Parse(t)
	if err != nil {
		return nil, err
	}
	if err := tmpl.Execute(buf, s); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

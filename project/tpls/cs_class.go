package tpls

import (
	"bytes"
	"text/template"
)

const csClassTpl = `namespace DataModel
{
	[System.Serializable]
	public class {{.Name}}
	{
	{{- range .Fields}}
		/// <summary>
		/// {{.Comment}}
		/// </summary>
		public {{.CsType}} {{.JSONTag}};
	{{- end}}
	}
}`

type IOLang struct {
	Name   string
	Fields []IOLangFields
}

type IOLangFields struct {
	CsType  string
	JSONTag string
	Comment string
}

func (s *IOLang) GenCs() ([]byte, error) {
	buf := new(bytes.Buffer)
	tmpl, err := template.New("IOLang").Parse(csClassTpl)
	if err != nil {
		return nil, err
	}
	if err := tmpl.Execute(buf, s); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

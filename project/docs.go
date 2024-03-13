package project

import (
	"fmt"
	jsoniter "github.com/json-iterator/go"
	"github.com/ucwebos/go3gen/project/parser"
	"github.com/ucwebos/go3gen/project/tpls"
	"github.com/ucwebos/go3gen/utils"
	"github.com/xbitgo/core/tools/tool_file"
	"log"
	"strings"
)

func (a *App) docsItem(dir string, f tpls.EntryFunItem, structList map[string]parser.XST) {
	var (
		uris     = strings.Split(f.URI2, "/")
		filename = fmt.Sprintf("%s/%s.md", dir, uris[len(uris)-1])
		reqXST   = structList[f.ReqName]
		respXST  = structList[f.RespName]
	)
	request := a.toDocsItemFields(reqXST.FieldList, structList, "")
	response := a.toDocsItemFields(respXST.FieldList, structList, "")

	t := &tpls.DocsItem{
		Name:      f.FunMark,
		RoutePath: f.URI,
		Request:   request,
		Response:  response,
		ExpJSON:   []byte{},
	}
	body := a.getJSON(respXST.FieldList, structList)
	sb, _ := jsoniter.MarshalIndent(body, "", "  ")

	t.ExpJSON = append(t.ExpJSON, []byte("```\n")...)
	t.ExpJSON = append(t.ExpJSON, sb...)
	t.ExpJSON = append(t.ExpJSON, []byte("\n```")...)

	buf, err := t.Execute()
	err = tool_file.WriteFile(filename, buf)
	if err != nil {
		log.Printf("doc item gen [%s] write file err: %v \n", filename, err)
	}
	return
}

func (a *App) getJSON(fields map[string]parser.XField, structList map[string]parser.XST) *utils.OrderMap {
	body := make(map[string]interface{})
	for _, it := range fields {
		j := it.GetTag("json")
		if j != nil {
			body[j.Name] = a.getJSONVal(it, structList)
		}
	}
	om := utils.NewOrderMap(utils.DefaultOrderMapKeySort)
	_ = om.LoadStringMap(body)
	return om
}

func (a *App) getJSONVal(field parser.XField, structList map[string]parser.XST) interface{} {
	switch field.SType {
	case parser.STypeStruct:
		sk := strings.TrimPrefix(field.Type, "*")
		if v, ok := structList[sk]; ok {
			return a.getJSON(v.FieldList, structList)
		}
	case parser.STypeSlice:
		sk := strings.ReplaceAll(strings.ReplaceAll(field.Type, "*", ""), "[]", "")
		if v, ok := structList[sk]; ok {
			return []interface{}{
				a.getJSON(v.FieldList, structList),
			}
		} else {
			return []interface{}{a.getZeroVal(field.Type)}
		}
	default:
		return a.getZeroVal(field.Type)
	}
	return ""
}
func (a *App) getZeroVal(xType string) interface{} {
	if xType == "bool" {
		return true
	}
	if strings.Contains(xType, "int") {
		return 0
	}
	if strings.Contains(xType, "float") {
		return 0.1
	}
	return ""
}

func (a *App) toDocsItemFields(fields map[string]parser.XField, structList map[string]parser.XST, prefix string) []tpls.DocsItemField {
	_fields := a.sortFields(fields)
	request := make([]tpls.DocsItemField, 0)
	for _, field := range _fields {
		j := field.GetTag("json")
		name := ""
		if j != nil {
			name = prefix + strings.ReplaceAll(j.Name, ",omitempty", "")
		}
		if j.Txt == "-" {
			continue
		}
		_type := field.Type
		switch field.SType {
		case parser.STypeStruct:
			prefix = strings.ReplaceAll(prefix, "[i].", "")
			_type = "object"
			request = append(request, tpls.DocsItemField{
				Name:    name,
				Type:    _type,
				Must:    "Y",
				Comment: field.Comment,
			})
			sk := strings.TrimPrefix(field.Type, "*")
			if v, ok := structList[sk]; ok {
				r := a.toDocsItemFields(v.FieldList, structList, prefix+"&emsp;&emsp;")
				request = append(request, r...)
			}
		case parser.STypeSlice:
			prefix = strings.ReplaceAll(prefix, "[i].", "")
			_type = "array"
			request = append(request, tpls.DocsItemField{
				Name:    name,
				Type:    _type,
				Must:    "Y",
				Comment: field.Comment,
			})
			sk := strings.ReplaceAll(strings.ReplaceAll(field.Type, "*", ""), "[]", "")
			if v, ok := structList[sk]; ok {
				r := a.toDocsItemFields(v.FieldList, structList, prefix+"&emsp;&emsp;[i].")
				request = append(request, r...)
			}
		default:
			request = append(request, tpls.DocsItemField{
				Name:    name,
				Type:    _type,
				Must:    "Y",
				Comment: field.Comment,
			})
		}

	}
	return request
}

func (a *App) sortFields(fields map[string]parser.XField) []parser.XField {
	r := make([]parser.XField, len(fields))
	for _, field := range fields {
		r[field.Idx] = field
	}
	return r

}

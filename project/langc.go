package project

import (
	"fmt"
	"github.com/ucwebos/go3gen/cfg"
	"github.com/ucwebos/go3gen/project/parser"
	"github.com/ucwebos/go3gen/project/tpls"
	"github.com/xbitgo/core/tools/tool_file"
	"github.com/xbitgo/core/tools/tool_str"
	"log"
	"os"
	"path"
	"strings"
)

func (a *App) LangC() {
	for _, tf := range a.typesGenFiles() {
		typesDir := path.Join(tf.Entry, "types")
		if !tool_file.Exists(typesDir) {
			continue
		}
		a.LangCDir(path.Join(cfg.C.RootPath, "panel", "docs", a.Name, "cs-models"), typesDir, "")
	}
}

func (a *App) LangCDir(outDir string, typesDir string, pkg string) {
	os.MkdirAll(outDir, 0777)
	fileInfos, err := os.ReadDir(typesDir)
	if err != nil {
		log.Fatal(err)
	}
	for _, fi := range fileInfos {
		if fi.IsDir() {
			pwd := path.Join(typesDir, fi.Name())
			a.LangCDir(outDir, pwd, fi.Name())
		}
	}
	ips, err := parser.Scan(typesDir, parser.ParseTypeWatch)
	if err != nil {
		log.Fatal(err)
	}
	for _, xst := range ips.StructList {
		cName := fmt.Sprintf("%s%s", tool_str.ToUFirst(tool_str.ToCamelCase(pkg)), xst.Name)
		fields := make([]tpls.IOLangFields, 0)
		for _, field := range xst.FieldList {
			if v, ok := a.toCsField(field, tool_str.ToUFirst(tool_str.ToCamelCase(pkg))); ok {
				fields = append(fields, v)
			}

		}
		t := &tpls.IOLang{
			Name:   cName,
			Fields: fields,
		}
		buf, err := t.GenCs()
		if err != nil {
			log.Fatal(err)
		}
		tool_file.WriteFile(path.Join(outDir, cName+".cs"), buf)
	}
}

func (a *App) toCsField(input parser.XField, pkg string) (tpls.IOLangFields, bool) {
	j := input.GetTag("json")
	if j.Txt == "-" {
		return tpls.IOLangFields{}, false
	}
	csType := ""
	switch input.SType {
	case parser.STypeBasic:
		if strings.Contains(input.Type, "int") {
			csType = "long"
		} else if strings.Contains(input.Type, "float") {
			csType = "double"
		} else {
			csType = input.Type
		}
	case parser.STypeStruct:
		_type := strings.TrimPrefix(input.Type, "*")
		arr := strings.Split(_type, ".")
		if len(arr) == 2 {
			csType = fmt.Sprintf("%s%s", pkg+tool_str.ToUFirst(tool_str.ToCamelCase(arr[0])), arr[1])
		} else {
			csType = _type
		}
	case parser.STypeSlice:
		_type := strings.ReplaceAll(input.Type, "*", "")
		_type = strings.ReplaceAll(_type, "[]", "")
		arr := strings.Split(_type, ".")
		if len(arr) == 2 {
			csType = fmt.Sprintf("%s%s", tool_str.ToUFirst(tool_str.ToCamelCase(arr[0])), arr[1]) + "[]"
		} else {
			csType = pkg + _type + "[]"
		}
	case parser.STypeMap:
		//csType = strings.Replace(input.Type,"map[","Dictionary<string, int>")
	case parser.STypeTime:
		csType = "string"
	}
	return tpls.IOLangFields{
		CsType:  csType,
		JSONTag: j.Name,
		Comment: input.Comment,
	}, true
}

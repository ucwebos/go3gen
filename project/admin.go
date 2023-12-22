package project

import (
	"github.com/ucwebos/go3gen/cfg"
	"github.com/ucwebos/go3gen/project/parser"
	"github.com/ucwebos/go3gen/project/tpls"
	"github.com/xbitgo/core/tools/tool_file"
	"github.com/xbitgo/core/tools/tool_str"
	"log"
	"os"
	"path"
	"sort"
)

func (a *AdminGroup) GenUI() {
	tg := a.ToTpl()
	// API
	a.adminAPI(tg)

	// UI
	a.adminUI(tg)

	// route
	a.adminRoute(tg)
}

type AdminGroup struct {
	Type      int
	Name      string
	Path      string
	AdminRoot string
	Title     string      `json:"title"`
	Icon      string      `json:"icon"`
	Items     []AdminItem `json:"items"`
}

type AdminItem struct {
	Entity string `json:"entity"`
	Title  string `json:"title"`
	Rule   string `json:"rule"`
	Icon   string `json:"icon"`
}

func (a *AdminGroup) ToTpl() *tpls.AdminGroup {
	entityDir := path.Join(a.Path, "entity")
	ipr, err := parser.Scan(entityDir, parser.ParseTypeWatch)
	if err != nil {
		log.Fatalf("GenUI: parse dir[%s], err: %v", entityDir, err)
	}
	tg := &tpls.AdminGroup{
		Name:     a.Name,
		NameVal:  tool_str.ToSnakeCase(a.Name),
		Title:    a.Title,
		Icon:     a.Icon,
		CrudList: make([]tpls.CrudItem, 0),
	}
	for _, it := range a.Items {
		if v, ok := ipr.StructList[it.Entity]; ok {
			var (
				fields    = make([]tpls.CrudItemField, 0)
				fieldList = make([]parser.XField, 0)
			)
			for _, field := range v.FieldList {
				fieldList = append(fieldList, field)
			}
			sort.SliceStable(fieldList, func(i, j int) bool {
				if fieldList[i].Idx < fieldList[j].Idx {
					return true
				}
				return false
			})

			for _, field := range fieldList {
				jv := field.GetTag("json")
				if jv != nil && jv.Name != "-" {
					fields = append(fields, tpls.CrudItemField{
						Name:  jv.Name,
						Title: field.Comment,
					})
				}
			}
			tg.CrudList = append(tg.CrudList, tpls.CrudItem{
				ApiBaseURL: "/micro/" + tg.Name,
				Name:       it.Entity,
				NameVal:    tool_str.ToSnakeCase(it.Entity),
				Title:      it.Title,
				Icon:       it.Icon,
				Fields:     fields,
			})
		}
	}
	return tg
}

func (a *AdminGroup) adminAPI(tg *tpls.AdminGroup) {
	dir := path.Join(cfg.C.RootPath, "panel", "admin", "micro", tg.Name)
	for _, item := range tg.CrudList {
		_t := &tpls.AdminAPIItem{
			Project: cfg.C.Project,
			AppName: tg.Name,
			PkgName: tg.Name,
			Name:    item.Name,
			NameVal: item.NameVal,
		}
		buf, err := _t.Execute()
		if err != nil {
			panic(err)
		}
		tool_file.WriteFile(path.Join(dir, item.NameVal+".go"), buf)
	}

}

func (a *AdminGroup) adminRoute(tg *tpls.AdminGroup) {
	routeFile := path.Join(a.AdminRoot, a.Name, "routes.ts")
	buf, err := tg.Execute(tpls.AdminRouteTpl)
	if err != nil {
		panic(err)
	}
	tool_file.WriteFile(routeFile, buf)
}

func (a *AdminGroup) adminUI(tg *tpls.AdminGroup) {
	for _, item := range tg.CrudList {
		dir := path.Join(a.AdminRoot, a.Name, item.NameVal)
		os.MkdirAll(dir, 0777)
		// api
		filename := path.Join(dir, "api.ts")
		buf, err := item.Execute(tpls.AdminUIApiTpl)
		if err != nil {
			panic(err)
		}
		tool_file.WriteFile(filename, buf)

		// crud
		filename = path.Join(dir, "crud.ts")
		buf, err = item.Execute(tpls.AdminUICrud)
		if err != nil {
			panic(err)
		}
		tool_file.WriteFile(filename, buf)

		//
		filename = path.Join(dir, "index.vue")
		buf, err = item.Execute(tpls.AdminUIIndex)
		if err != nil {
			panic(err)
		}
		tool_file.WriteFile(filename, buf)
	}

}

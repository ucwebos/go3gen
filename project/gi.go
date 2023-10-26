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
	"regexp"
	"strings"
)

func (a *App) GI() {
	for _, ipr := range a.scanGIList() {
		a._gi(ipr)
	}
}

func (a *App) _gi(iParser *parser.IParser) {
	gdi := tpls.GI{
		Pkg:  iParser.Package,
		List: make([]tpls.GItem, 0),
	}
	for _, xst := range iParser.StructList {
		if xst.GI {
			it := tpls.GItem{
				Name:    xst.Name,
				NameVal: tool_str.ToLFirst(xst.Name),
			}
			if nMth, ok := iParser.NewFuncList[xst.Name]; ok {
				it.NewReturnsLen = len(nMth.Results)
			}
			gdi.List = append(gdi.List, it)
		}
	}
	if len(gdi.List) == 0 {
		return
	}
	filename := path.Join(iParser.Pwd, "gi_gen.go")
	buf, err := gdi.Execute()
	if err != nil {
		panic(err)
	}
	buf = a.format(buf, filename)
	err = tool_file.WriteFile(filename, buf)
	if err != nil {
		log.Printf("gen gi [%s] write file err: %v \n", filename, err)
	}
	log.Printf("gen gi file %s \n", filename)
}

func (a *App) MicroFun() bool {
	// 读取micro 入口配置
	doc, err := parser.MicroEntryDoc(path.Join(a.Path, "entry.go"))
	if err == nil && doc != "" {
		items := a.MicroEntryDocParse(doc)
		tpl := tpls.MicroEntry{
			Project:    cfg.C.Project,
			AppName:    a.Name,
			AppNameUF:  tool_str.ToUFirst(a.Name),
			AppPkgPath: a.appPkgPath(),
			FunList:    items,
		}
		buf, err := tpl.Execute(tpls.MicroEntryTpl)
		if err != nil {
			log.Printf("gen MicroEntry %s err: %v \n", a.Name, err)
			return false
		}
		filename := path.Join(cfg.C.RootPath, "provider", a.Name+"_gen.go")
		buf = a.format(buf, filename)
		_ = tool_file.WriteFile(filename, buf)

		buf, err = tpl.Execute(tpls.MicroServiceTpl)
		if err != nil {
			log.Printf("gen MicroEntry %s err: %v \n", a.Name, err)
			return false
		}
		filename = path.Join(a.Path, "service_gen.go")
		buf = a.format(buf, filename)
		_ = tool_file.WriteFile(filename, buf)
		log.Printf("gen micro-fun file %s \n", filename)
		// 生成IO types
		a._microFunIO(items)
		// 生成service 方法
		a._microService(items)
		// 生成service 单元测试
		a._microServiceUniTest(items)
		return true
	}
	return false
}

func (a *App) _microFunIO(items []tpls.MicroFunItem) {
	tpl := tpls.MicroTypes{
		AppName:    a.Name,
		AppPkgPath: a.appPkgPath(),
		FunList:    items,
	}
	buf, err := tpl.Execute()
	if err != nil {
		log.Printf("gen MicroTypes %s err: %v \n", a.Name, err)
		return
	}
	filename := path.Join(a.Path, "types_"+a.Name, "types.go")
	buf = a.format(buf, filename)
	_ = tool_file.WriteFile(filename, buf)
	log.Printf("gen micro-types file %s \n", filename)
}

func (a *App) _microService(items []tpls.MicroFunItem) {
	itemsMap := map[string][]tpls.MicroFunItem{}
	for _, item := range items {
		if _, ok := itemsMap[item.Service]; !ok {
			itemsMap[item.Service] = []tpls.MicroFunItem{}
		}
		itemsMap[item.Service] = append(itemsMap[item.Service], item)
	}

	ipr, err := parser.Scan(path.Join(a.Path, "service"), parser.ParseTypeImpl)
	if err != nil {
		log.Fatal(err)
	}
	for service, funItems := range itemsMap {
		var (
			filename = path.Join(a.Path, "service", strings.ToLower(service)+".go")
			buf      []byte
			err      error
		)
		//fmt.Println(filename)
		if xst, ok := ipr.StructList[service]; ok {
			_buf, err := os.ReadFile(filename)
			if err != nil {
				log.Fatal(err)
			}
			tpl := tpls.MServiceFuncAppend{
				Body:    _buf,
				AppName: a.Name,
				FunList: make([]tpls.MicroFunItem, 0),
			}
			for _, it := range funItems {
				if _, ok := xst.Methods[it.Method]; !ok {
					tpl.FunList = append(tpl.FunList, it)
				}
			}
			if len(tpl.FunList) > 0 {
				buf, err = tpl.Execute()
				if err != nil {
					log.Fatal(err)
				}
				log.Printf("append micro-service file %s \n", filename)
			}

		} else {
			funList := make([]tpls.MicroFunItem, 0)
			for _, item := range items {
				if item.Service == service {
					funList = append(funList, item)
				}
			}
			tpl := tpls.MServiceFunc{
				AppName:    a.Name,
				AppPkgPath: a.appPkgPath(),
				Service:    service,
				FunList:    funList,
			}
			buf, err = tpl.Execute()
			if err != nil {
				log.Fatal(err)
			}
		}
		buf = a.format(buf, filename)
		err = tool_file.WriteFile(filename, buf)
		log.Printf("gen micro-service file %s \n", filename)
		if err != nil {
			log.Fatal(err)
		}
	}
}

func (a *App) _microServiceUniTest(items []tpls.MicroFunItem) {
	itemsMap := map[string][]tpls.MicroFunItem{}
	for _, item := range items {
		if _, ok := itemsMap[item.Service]; !ok {
			itemsMap[item.Service] = []tpls.MicroFunItem{}
		}
		itemsMap[item.Service] = append(itemsMap[item.Service], item)
	}

	for service, funItems := range itemsMap {
		var (
			filename = path.Join(a.Path, "service", strings.ToLower(service)+"_test.go")
			buf      []byte
			err      error
		)
		tpl := tpls.MicroTesting{
			AppName:    a.Name,
			AppNameUF:  tool_str.ToUFirst(a.Name),
			AppPkgPath: a.appPkgPath(),
			Service:    service,
			FunList:    funItems,
		}
		buf, err = tpl.Execute()
		if err != nil {
			log.Fatal(err)
		}
		buf = a.format(buf, filename)
		err = tool_file.WriteFile(filename, buf)
		if err != nil {
			log.Fatal(err)
		}
		log.Printf("gen unit-test file %s \n", filename)

	}
}

var microFunExp = regexp.MustCompile(`(\w+)\s+(.+)\s+\[([\w|.]+)]`)

func (a *App) MicroEntryDocParse(doc string) []tpls.MicroFunItem {
	items := make([]tpls.MicroFunItem, 0)
	lines := strings.Split(doc, "\n")
	for _, line := range lines {
		r := microFunExp.FindStringSubmatch(line)
		if len(r) == 4 {
			items = append(items, tpls.MicroFunItem{
				Service:  strings.Split(r[3], ".")[0],
				Method:   strings.Split(r[3], ".")[1],
				FunName:  r[1],
				FunMark:  r[2],
				ReqName:  r[1] + "Req",
				RespName: r[1] + "Resp",
			})
		}
	}

	return items
}

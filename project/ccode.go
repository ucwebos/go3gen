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
	"regexp"
	"strings"
)

func (a *App) _cRepo(xstList []parser.XST) {
	entityList := make([]string, 0)
	bcfgList := make([]string, 0)
	hasIDMap := map[string]bool{}
	for _, it := range xstList {
		if it.BCFG {
			bcfgList = append(bcfgList, it.Name)
			continue
		}
		for _, field := range it.FieldList {
			tag := field.GetTag("db")
			if tag != nil && tag.Txt != "-" {
				entityList = append(entityList, it.Name)
				break
			}
		}
		for _, field := range it.FieldList {
			if field.Name == "ID" {
				hasIDMap[it.Name] = true
				break
			}
		}
	}
	for _, s := range entityList {
		tpl := tpls.Repo{
			ProjectName: cfg.C.Project,
			AppPkgPath:  a.appPkgPath(),
			EntityName:  s,
			TableName:   tool_str.ToSnakeCase(s),
			HasID:       hasIDMap[s],
		}
		buf, err := tpl.Execute()
		if err != nil {
			log.Printf("gen Repo %s err: %v \n", s, err)
			return
		}
		filename := path.Join(a.Path, "repo", fmt.Sprintf("%s_repo.go", tool_str.ToSnakeCase(s)))
		if !tool_file.Exists(filename) {
			buf = a.format(buf, filename)
			log.Printf("gen repo file %s \n", filename)
			err = tool_file.WriteFile(filename, buf)
			if err != nil {
				return
			}
		}
		buf, err = tpl.ExecuteImpl()
		if err != nil {
			log.Printf("gen Repo.dbal %s err: %v \n", s, err)
			return
		}
		filename = path.Join(a.Path, "repo", "dbal", fmt.Sprintf("%s_dbal.go", tool_str.ToSnakeCase(s)))
		if !tool_file.Exists(filename) {
			buf = a.format(buf, filename)
			log.Printf("gen dbal file %s \n", filename)
			err = tool_file.WriteFile(filename, buf)
			if err != nil {
				return
			}
		}
	}

	for _, s := range bcfgList {
		tpl := tpls.RepoBCFG{
			ProjectName: cfg.C.Project,
			AppPkgPath:  a.appPkgPath(),
			EntityName:  s,
			TableName:   tool_str.ToSnakeCase(s),
		}
		buf, err := tpl.Execute()
		if err != nil {
			log.Printf("gen RepoBCFG %s err: %v \n", s, err)
			return
		}
		filename := path.Join(a.Path, "repo", fmt.Sprintf("%s_repo.go", tool_str.ToSnakeCase(s)))
		if !tool_file.Exists(filename) {
			buf = a.format(buf, filename)
			log.Printf("gen repo file %s \n", filename)
			err = tool_file.WriteFile(filename, buf)
			if err != nil {
				return
			}
		}
	}
}

func (a *App) _cHandler(ef tpls.HttpEntry) {
	dir := path.Join(ef.EntryPath, "handler")
	ips, err := parser.Scan(dir, parser.ParseTypeImpl)
	if err != nil {
		log.Panic(err)
	}
	hasFuncMap := ips.FuncList
	for _, group := range ef.Groups {
		filename := path.Join(ef.EntryPath, "handler", group.Group+".go")
		if tool_file.Exists(filename) {
			buf, err := os.ReadFile(filename)
			if err != nil {
				log.Panic(err)
			}
			_t := tpls.HandlerFuncAppend{
				Body:    buf,
				FunList: make([]tpls.EntryFunItem, 0),
			}
			for _, it := range group.FunList {
				if _, ok := hasFuncMap[it.FunName]; !ok {
					_t.FunList = append(_t.FunList, it)
				}
			}
			buf, err = _t.Execute()
			if err != nil {
				panic(err)
			}
			buf = a.format(buf, filename)
			log.Printf("append handler file %s \n", filename)
			err = tool_file.WriteFile(filename, buf)
			if err != nil {
				return
			}
		} else {
			_t := tpls.HandlerFunc{
				EntryPath: ef.EntryPkgPath,
				Entry:     ef.EntryName,
				Group:     group.Group,
				FunList:   group.FunList,
			}
			buf, err := _t.Execute()
			if err != nil {
				panic(err)
			}
			buf = a.format(buf, filename)
			log.Printf("gen handler file %s \n", filename)
			err = tool_file.WriteFile(filename, buf)
			if err != nil {
				return
			}
		}
	}
}

func (a *App) _cHttpTypes(ef tpls.HttpEntry) {
	dir := path.Join(ef.EntryPath, "types")
	ips, err := parser.Scan(dir, parser.ParseTypeWatch)
	if err != nil {
		log.Panic(err)
	}
	hasStructMap := ips.StructList
	for _, group := range ef.Groups {
		filename := path.Join(ef.EntryPath, "types", "io_"+group.Group+".go")
		if tool_file.Exists(filename) {
			buf, err := os.ReadFile(filename)
			if err != nil {
				log.Panic(err)
			}
			_t := tpls.HandlerTypesAppend{
				Body:    buf,
				FunList: make([]tpls.EntryFunItem, 0),
			}
			for _, it := range group.FunList {
				if _, ok := hasStructMap[it.ReqName]; !ok {
					_t.FunList = append(_t.FunList, it)
				}
			}
			buf, err = _t.Execute()
			if err != nil {
				panic(err)
			}
			buf = a.format(buf, filename)
			log.Printf("append io-types file %s \n", filename)
			err = tool_file.WriteFile(filename, buf)
			if err != nil {
				return
			}
		} else {
			_t := tpls.HandlerTypes{
				EntryPath: ef.EntryPkgPath,
				Entry:     ef.EntryName,
				Group:     group.Group,
				FunList:   group.FunList,
			}
			buf, err := _t.Execute()
			if err != nil {
				panic(err)
			}
			buf = a.format(buf, filename)
			log.Printf("gen io-types file %s \n", filename)
			err = tool_file.WriteFile(filename, buf)
			if err != nil {
				return
			}
		}

	}
}

func (a *App) _cHttpDocs(ef tpls.HttpEntry) {
	ips, err := parser.Scan(path.Join(ef.EntryPath, "types"), parser.ParseTypeWatch)
	if err != nil {
		log.Panic(err)
	}

	// modules
	typesPath := path.Join(ef.EntryPath, "types")
	fileInfos, err := os.ReadDir(typesPath)
	if err != nil {
		log.Fatal(err)
	}
	for _, fi := range fileInfos {
		if fi.IsDir() {
			pwd := path.Join(typesPath, fi.Name())
			_ips, err := parser.Scan(pwd, parser.ParseTypeWatch)
			if err != nil {
				log.Panic(err)
			}
			for s, xst := range _ips.StructList {
				ips.StructList[fi.Name()+"."+s] = xst
			}
		}
	}
	for _, group := range ef.Groups {
		groupURIName := strings.ReplaceAll(tool_str.ToSnakeCase(group.Group), "_", "-")
		dir := path.Join(cfg.C.RootPath, "panel", "docs", a.Name, ef.EntryName, groupURIName)
		if a.Type == TypeBFF {
			dir = path.Join(cfg.C.RootPath, "panel", "docs", a.Name, groupURIName)
		}
		os.MkdirAll(dir, 0777)
		for _, fun := range group.FunList {
			a.docsItem(dir, fun, ips.StructList)
		}
	}
	filename := path.Join(cfg.C.RootPath, "panel", "docs", a.Name, "_sidebar.md")
	sider, err := os.ReadFile(filename)
	siderStr := string(sider)
	idx1 := strings.Index(siderStr, "---")

	t := &tpls.DocsSidebar{
		Entry:  ef.EntryName,
		Groups: ef.Groups,
	}
	buf, err := t.Execute()
	if siderStr == "" {
		return
	}
	_str := siderStr[:idx1+3] + string(buf)
	tool_file.WriteFile(filename, []byte(_str))
}

func (a *App) CHandlerAndDoc() {
	for _, tf := range a.typesGenFiles() {
		remark, err := parser.MicroEntryDoc(path.Join(tf.Entry, "route", "routes.go"), "_gen")
		if err == nil && remark != "" {
			groups := a.BizEntryDocParse(tf.EntryName, remark, false)
			_tp := tpls.HttpEntry{
				Project:      cfg.C.Project,
				AppName:      a.Name,
				AppNameUF:    tool_str.ToUFirst(a.Name),
				AppPkgPath:   a.appPkgPath(),
				EntryPath:    tf.Entry,
				EntryPkgPath: tf.EntryPkgPath,
				EntryName:    tf.EntryName,
				Groups:       groups,
			}
			buf, err := _tp.Execute(tpls.HttpRouteTpl)
			if err != nil {
				panic(err)
			}
			filename := path.Join(tf.Entry, "route", "routes_gen.go")
			buf = a.format(buf, filename)
			log.Printf("gen routes file %s \n", filename)
			err = tool_file.WriteFile(filename, buf)
			if err != nil {
				return
			}
			a._cHandler(_tp)
			a._cHttpTypes(_tp)
			a._cHttpDocs(_tp)
		}
	}
}

func (a *App) CSocketHandlerAndDoc() {
	for _, tf := range a.typesGenFiles() {
		for _, socketType := range []string{
			"default", "temporary",
		} {
			remark, err := parser.MicroEntryDoc(path.Join(tf.Entry, "route", fmt.Sprintf("soc_%s_routes.go", socketType)), fmt.Sprintf("_gen%s", tool_str.ToUFirst(socketType)))
			if err == nil && remark != "" {
				groups := a.BizEntryDocParse(tf.EntryName, remark, true)
				_tp := tpls.HttpEntry{
					Project:      cfg.C.Project,
					AppName:      a.Name,
					AppNameUF:    tool_str.ToUFirst(a.Name),
					AppPkgPath:   a.appPkgPath(),
					EntryPath:    tf.Entry,
					EntryPkgPath: tf.EntryPkgPath,
					EntryName:    tf.EntryName,
					SocketType:   socketType,
					SocketTypeUF: tool_str.ToUFirst(socketType),
					Groups:       groups,
				}
				buf, err := _tp.Execute(tpls.SocketRouteTpl)
				if err != nil {
					panic(err)
				}
				filename := path.Join(tf.Entry, "route", fmt.Sprintf("soc_%s_routes_gen.go", socketType))
				buf = a.format(buf, filename)
				log.Printf("gen socket routes file %s \n", filename)
				err = tool_file.WriteFile(filename, buf)
				if err != nil {
					return
				}
				a._cHandler(_tp)
				a._cHttpTypes(_tp)
				a._cHttpDocs(_tp)
			}
		}

	}
}

func (a *App) CRepos() {
	// 解析 entity
	ipr := a.scanEntity()
	xstList := ipr.GetStructList()
	a._cRepo(xstList)
}

const (
	NLOG = "@NLOG"
	NTLP = "@NTLP"
)

var (
	handlerFunExp    = regexp.MustCompile(`(.+)\s+\[(\w+)]`)
	handlerGroupExp  = regexp.MustCompile(`#(\S+)\s+(\w+)`)
	handlerMiddleExp = regexp.MustCompile(`@M\(([\w|,|\(|\)]+)\)`)
)

func (a *App) BizEntryDocParse(entryName string, doc string, socket bool) []*tpls.EntryGroup {
	groups := make([]*tpls.EntryGroup, 0)
	lines := strings.Split(doc, "\n")
	for _, line := range lines {
		rg := handlerGroupExp.FindStringSubmatch(line)
		if len(rg) == 3 {
			group := &tpls.EntryGroup{
				Group:        rg[2],
				GroupUFirst:  tool_str.ToUFirst(rg[2]),
				GroupName:    rg[1],
				FunList:      make([]tpls.EntryFunItem, 0),
				GMiddlewares: make([]string, 0),
			}
			rm := handlerMiddleExp.FindStringSubmatch(line)
			if len(rm) == 2 {
				group.GMiddlewares = strings.Split(rm[1], ",")
			}
			groups = append(groups, group)
			continue
		}
		if len(groups) == 0 {
			continue
		}
		group := groups[len(groups)-1]
		r := handlerFunExp.FindStringSubmatch(line)
		withLog := true
		if strings.Contains(line, NLOG) {
			withLog = false
		}
		if len(r) == 3 {
			method := tool_str.ToUFirst(r[2])
			fun := tool_str.ToUFirst(group.Group) + method
			m := strings.ReplaceAll(tool_str.ToSnakeCase(method), "_", "-")
			groupURIName := strings.ReplaceAll(tool_str.ToSnakeCase(group.Group), "_", "-")
			item := tpls.EntryFunItem{
				FunName:     fun,
				Method:      method,
				FunMark:     r[1],
				ReqName:     fun + "Req",
				RespName:    fun + "Resp",
				WithLog:     withLog,
				Middlewares: make([]string, 0),
				URI:         fmt.Sprintf("/%s/%s/%s", entryName, groupURIName, m),
				URI2:        fmt.Sprintf("/%s/%s", groupURIName, m),
			}
			if socket {
				item.URI = fmt.Sprintf("%s.%s", group.Group, strings.ToLower(method))
			}
			rm := handlerMiddleExp.FindStringSubmatch(line)
			if len(rm) == 2 {
				item.Middlewares = strings.Split(rm[1], ",")
				if len(group.GMiddlewares) > 0 {
					middlewareMap := map[string]struct{}{}
					middlewares := make([]string, 0)
					for _, middleware := range group.GMiddlewares {
						middlewareMap[middleware] = struct{}{}
						middlewares = append(middlewares, middleware)
					}
					for _, middleware := range item.Middlewares {
						if _, ok := middlewareMap[middleware]; !ok {
							middlewares = append(middlewares, middleware)
						}
					}
					item.Middlewares = middlewares
				}
			} else {
				item.Middlewares = group.GMiddlewares
			}
			group.FunList = append(group.FunList, item)
		}
	}
	return groups
}

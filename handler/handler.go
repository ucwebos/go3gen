package handler

import (
	"bufio"
	"fmt"
	jsoniter "github.com/json-iterator/go"
	"github.com/spf13/cobra"
	"github.com/ucwebos/go3gen/cfg"
	"github.com/ucwebos/go3gen/project"
	"github.com/ucwebos/go3gen/project/tpls"
	"github.com/xbitgo/core/tools/tool_file"
	"github.com/xbitgo/core/tools/tool_str"
	"go/format"
	"gopkg.in/yaml.v3"
	"log"
	"os"
	"path"
	"strings"
)

var pwd, _ = os.Getwd()

func CmdList() []*cobra.Command {
	return []*cobra.Command{
		{
			Use:   "generate",
			Short: "生成go代码",
			Long:  "生成go代码",
			Run:   generate,
		},
		{
			Use:   "add-m",
			Short: "添加一个微模块",
			Long:  "添加一个微模块 go3gen add-m {name}",
			Run:   addMs,
		},
		{
			Use:   "add-bff",
			Short: "添加一个BFF",
			Long:  "添加一个BFF go3gen add-bff {name}",
			Run:   addBff,
		},
		{
			Use:   "admin",
			Short: "生成微模块CRUD后台",
			Long:  "生成微模块CRUD后台",
			Run:   admin,
		},
		{
			Use:   "sql",
			Short: "生成SQL文件",
			Long:  "生成SQL文件",
			Run:   sql,
		},
	}
}

func generate(cmd *cobra.Command, args []string) {
	_project()
	microList := make([]string, 0)
	for _, app := range _scanMicros() {
		// ETO 生成
		app.ETO()
		// CRepos 生成
		app.CRepos()
		// MicroFun 生成
		if app.MicroFun() {
			microList = append(microList, tool_str.ToUFirst(app.Name))
		}
		// GI
		app.GI()
	}
	tpl := tpls.MicroProvider{MicroList: microList}
	buf, err := tpl.Execute()
	if err != nil {
		log.Fatalf("gen MicroProvider err: %v \n", err)
		return
	}
	filename := path.Join(cfg.C.RootPath, "provider", "provider_gen.go")
	buf, err = format.Source(buf)
	if err != nil {
		log.Fatalf("gen x MicroProvider format err: %v \n", err)
	}
	_ = tool_file.WriteFile(filename, buf)

	for _, app := range _scanBusiness() {
		// ETO 生成
		app.ETO()
		// CRepos 生成
		app.CRepos()
		// GI
		app.GI()
		// CHandlerAndDoc
		app.CHandlerAndDoc()
		app.CSocketHandlerAndDoc()
		// Proto 生成
		app.Protoc()
	}

	for _, app := range _scanBFF() {
		app.BffModulesTypes()
		// GI
		app.GI()
		// CHandlerAndDoc
		app.CHandlerAndDoc()
		app.CSocketHandlerAndDoc()
		// Proto 生成
		app.Protoc()
	}
}

func addMs(cmd *cobra.Command, args []string) {
	_project()
	if len(args) < 1 {
		log.Fatalf("请输入模块名 ")
		return
	}
	name := args[0]
	project.AddMicro(name)
}

func addBff(cmd *cobra.Command, args []string) {
	_project()
	if len(args) < 1 {
		log.Fatalf("请输入BFF名 ")
		return
	}
	name := args[0]
	project.AddBFF(name)
}

func sql(cmd *cobra.Command, args []string) {
	_project()
	for _, app := range _scanBusiness() {
		dsn, ok := cfg.C.DBMaps[app.Name]
		if !ok {
			dsn = cfg.C.DB
		}
		if dsn == "" {
			log.Fatalf("db not set!")
			return
		}
		app.GenSql(dsn)
	}

	for _, app := range _scanMicros() {
		dsn, ok := cfg.C.DBMaps[app.Name]
		if !ok {
			dsn = cfg.C.DB
		}
		if dsn == "" {
			log.Fatalf("db not set!")
			return
		}
		app.GenSql(dsn)
	}
}

type ApiGroup struct {
	Micro  string   `json:"micro"`
	Entity []string `json:"entity"`
}

func admin(cmd *cobra.Command, args []string) {
	_project()
	var (
		genGroups = make([]*tpls.AdminGroup, 0)
	)
	genFile := path.Join(cfg.C.RootPath, "panel", "micro", "gen.json")
	if !tool_file.Exists(genFile) {
		log.Fatalf("genFile: %s not exists", genFile)
		return
	}
	JSONBuf, err := os.ReadFile(genFile)
	if err != nil {
		log.Fatalf("genFile %s read err: %v", genFile, err)
		return
	}
	groups := make([]ApiGroup, 0)
	err = jsoniter.Unmarshal(JSONBuf, &groups)
	if err != nil {
		log.Fatalf("genFile %s unmarshal err: %v", genFile, err)
		return
	}
	for _, group := range groups {
		items := make([]tpls.CrudItem, 0)
		for _, s := range group.Entity {
			items = append(items, tpls.CrudItem{
				Group:   tool_str.ToUFirst(group.Micro),
				Name:    s,
				NameVal: tool_str.ToSnakeCase(s),
			})
		}
		genGroups = append(genGroups, &tpls.AdminGroup{
			Name:     group.Micro,
			NameVal:  tool_str.ToSnakeCase(group.Micro),
			CrudList: items,
		})
	}
	for _, group := range genGroups {
		adminAPI(group)
	}
	aaRoute := &tpls.AdminAPIRoute{
		Project: cfg.C.Project,
		Groups:  genGroups,
	}
	buf, err := aaRoute.Execute()
	if err != nil {
		panic(err)
	}
	tool_file.WriteFile(path.Join(cfg.C.RootPath, "panel", "micro", "route.go"), buf)
}

func adminAPI(tg *tpls.AdminGroup) {
	dir := path.Join(cfg.C.RootPath, "panel", "micro", tg.Name)
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
		if !tool_file.Exists(path.Join(dir, item.NameVal+".go")) {
			tool_file.WriteFile(path.Join(dir, item.NameVal+".go"), buf)
		}
	}
}

func _project() {
	cfg.C.RootPath = pwd
	goMod := fmt.Sprintf("%s/go.mod", pwd)
	f, err := os.Open(goMod)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()
	br := bufio.NewReader(f)
	for {
		line, _, err := br.ReadLine()
		if err != nil {
			log.Fatal(err)
		}
		module := strings.Replace(string(line), "module", "", 1)
		cfg.C.Project = strings.TrimSpace(module)
		break
	}
	if tool_file.Exists(cfg.C.RootPath + "/.go3gen.yaml") {
		buf, _ := os.ReadFile(cfg.C.RootPath + "/.go3gen.yaml")
		err := yaml.Unmarshal(buf, cfg.C)
		if err != nil {
			panic(err)
		}
	}
}

func _scanMicros() []*project.App {
	appList := make([]*project.App, 0)
	micro := path.Join(pwd, "micro")
	if tool_file.Exists(micro) {
		fileInfos, err := os.ReadDir(micro)
		if err != nil {
			log.Fatal(err)
		}
		for _, fi := range fileInfos {
			if fi.IsDir() {
				iPwd := path.Join(micro, fi.Name())
				app := project.NewApp(project.TypeMicro, fi.Name(), iPwd)
				appList = append(appList, app)
			}
		}
	}

	return appList
}

func _scanBusiness() []*project.App {
	appList := make([]*project.App, 0)
	biz := path.Join(pwd, "business")
	if tool_file.Exists(biz) {
		fileInfos, err := os.ReadDir(biz)
		if err != nil {
			log.Fatal(err)
		}
		for _, fi := range fileInfos {
			if fi.IsDir() {
				iPwd := path.Join(biz, fi.Name())
				app := project.NewApp(project.TypeAPI, fi.Name(), iPwd)
				appList = append(appList, app)
			}
		}
	}
	return appList
}

func _scanBFF() []*project.App {
	appList := make([]*project.App, 0)
	bff := path.Join(pwd, "bff")
	if tool_file.Exists(bff) {
		fileInfos, err := os.ReadDir(bff)
		if err != nil {
			log.Fatal(err)
		}
		for _, fi := range fileInfos {
			if fi.IsDir() {
				iPwd := path.Join(bff, fi.Name())
				app := project.NewApp(project.TypeBFF, fi.Name(), iPwd)
				appList = append(appList, app)
			}
		}
	}
	return appList
}

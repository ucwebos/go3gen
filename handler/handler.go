package handler

import (
	"bufio"
	"fmt"
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
			Short: "生成所有go代码 依次 eto > c.code > conv > GI",
			Long:  "生成所有go代码 依次 eto > c.code > conv > GI",
			Run:   generate,
		},
		{
			Use:   "sql",
			Short: "生成接口单元测试用例",
			Long:  "生成接口单元测试用例; 参数 {app}; app为应用名称 必须",
			Run:   sql,
		},
	}
}

func generate(cmd *cobra.Command, args []string) {
	_project()
	microList := make([]string, 0)
	for _, app := range _scanMicros() {
		//if app.Name != "mail" {
		//	continue
		//}
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
		app.CWsHandlerAndDoc()
		// Proto 生成
		app.Protoc()
	}
}

func sql(cmd *cobra.Command, args []string) {
	_project()
	var (
		dbSet = cfg.DBSet{}
	)
	buf, _ := os.ReadFile(cfg.C.RootPath + "/.go3gen.yaml")
	err := yaml.Unmarshal(buf, &dbSet)
	if err != nil {
		panic(err)
	}

	for _, app := range _scanBusiness() {
		fmt.Println(app.Name)
		dsn, ok := dbSet.DBMaps[app.Name]
		if !ok {
			dsn = dbSet.DB
		}
		if dsn == "" {
			log.Fatalf("db not set!")
			return
		}
		app.GenSql(dsn)
	}

	for _, app := range _scanMicros() {
		fmt.Println("micro:", app.Name)
		dsn, ok := dbSet.DBMaps[app.Name]
		if !ok {
			dsn = dbSet.DB
		}
		if dsn == "" {
			log.Fatalf("db not set!")
			return
		}
		app.GenSql(dsn)
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
				app := project.NewApp(project.TypeModule, fi.Name(), iPwd)
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

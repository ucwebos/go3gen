package project

import (
	"fmt"
	"github.com/ucwebos/go3gen/cfg"
	"github.com/ucwebos/go3gen/project/tpls"
	"github.com/xbitgo/core/tools/tool_file"
	"github.com/xbitgo/core/tools/tool_str"
	"log"
	"os"
)

func AddMicro(name string) {
	msDir := fmt.Sprintf("%s/micro/%s", cfg.C.RootPath, name)
	if tool_file.Exists(msDir) {
		log.Fatal("已存在该Micro 不能重复创建")
		return
	}
	// dir
	for _, item := range []string{
		"config",
		"entity",
		"repo/converter",
		"repo/do",
		"repo/sql",
		"repo/dbal",
		"service",
		fmt.Sprintf("types_%s", name),
	} {
		dir := fmt.Sprintf("%s/%s", msDir, item)
		err := os.MkdirAll(dir, os.ModePerm)
		if err != nil {
			panic(err)
		}
	}
	// file
	add := tpls.AddMicro{
		Project: cfg.C.Project,
		Name:    name,
		NameUF:  tool_str.ToUFirst(name),
	}
	for _fName, _buf := range map[string][]byte{
		"config/config.go": add.Execute(tpls.AddMsCfgTpl),
		"entry.go":         add.Execute(tpls.AddMsEntryTpl),
	} {
		filename := fmt.Sprintf("%s/%s", msDir, _fName)
		err := tool_file.WriteFile(filename, _buf)
		if err != nil {
			panic(err)
		}
	}
}

func AddBFF(name string) {
	msDir := fmt.Sprintf("%s/bff/%s", cfg.C.RootPath, name)
	if tool_file.Exists(msDir) {
		log.Fatal("已存在该BFF 不能重复创建")
		return
	}
	// dir
	for _, item := range []string{
		"config",
		"entity",
		"converter",
		"handler",
		"middleware",
		"push",
		"route",
		"service",
		"types", name,
	} {
		dir := fmt.Sprintf("%s/%s", msDir, item)
		err := os.MkdirAll(dir, os.ModePerm)
		if err != nil {
			panic(err)
		}
	}
	// file
	//add := tpls.AddMicro{
	//	Project: cfg.C.Project,
	//	Name:    name,
	//	NameUF:  tool_str.ToUFirst(name),
	//}
	//for _fName, _buf := range map[string][]byte{
	//	"config/config.go": add.Execute(tpls.AddBffCfgTpl),
	//	"route/conf.go":    add.Execute(tpls.AddBffRouteCfgTpl),
	//  "route/route.go":   add.Execute(tpls.AddBffRouteTpl),
	//} {
	//	filename := fmt.Sprintf("%s/%s", msDir, _fName)
	//	err := tool_file.WriteFile(filename, _buf)
	//	if err != nil {
	//		panic(err)
	//	}
	//}
}

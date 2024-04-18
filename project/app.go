package project

import (
	"github.com/ucwebos/go3gen/cfg"
	"github.com/ucwebos/go3gen/project/parser"
	"log"
	"os"
	"path"
	"strings"
)

const (
	TypeAPI   = 1
	TypeMicro = 2
	TypeBFF   = 3
)

type App struct {
	Type int // 1 business 2 micro
	Name string
	Path string
}

type MicroXSTList struct {
	Name    string
	PkgPath string
	XSTList []parser.XST
}

func NewApp(aType int, name string, pwd string) *App {
	return &App{
		Type: aType,
		Name: name,
		Path: pwd,
	}
}

func (a *App) scanEntity() *parser.IParser {
	ipr, err := parser.Scan(path.Join(a.Path, "entity"), parser.ParseTypeWatch)
	if err != nil {
		log.Fatalf("EntityTypeDef: parse dir[%s], err: %v", path.Join(a.Path, "entity"), err)
	}
	return ipr
}
func (a *App) scanDo() *parser.IParser {
	ipr, err := parser.Scan(path.Join(a.Path, "repo", "do"), parser.ParseTypeDo)
	if err != nil {
		log.Fatalf("EntityTypeDef: parse dir[%s], err: %v", path.Join(a.Path, "entity"), err)
	}
	return ipr
}

func (a *App) scanGIList() []*parser.IParser {
	iprs := make([]*parser.IParser, 0)
	ignoreDir := map[string]struct{}{
		"cmd":        {},
		"config":     {},
		"entity":     {},
		"jobs":       {},
		"scripts":    {},
		"converter":  {},
		"handler":    {},
		"middleware": {},
		"route":      {},
		"types":      {},
	}
	fileInfos, err := os.ReadDir(a.Path)
	if err != nil {
		log.Fatal(err)
	}
	for _, fi := range fileInfos {
		if fi.IsDir() {
			if _, ok := ignoreDir[fi.Name()]; !ok {
				ipr, err := parser.Scan(path.Join(a.Path, fi.Name()), parser.ParseTypeWatch)
				if err != nil {
					log.Fatal(err)
				}
				iprs = append(iprs, ipr)
			}
		}
	}
	return iprs
}

func (a *App) appPkgPath() string {
	return cfg.C.Project + strings.Replace(a.Path, cfg.C.RootPath, "", 1)
}

func (a *App) entityTypeDefGenFile() string {
	return path.Join(a.Path, "entity", "type_def_code_gen.go")
}

func (a *App) doTypeDefGenFile() string {
	return path.Join(a.Path, "repo", "do", "type_def_code_gen.go")
}

func (a *App) doGenFile() string {
	return path.Join(a.Path, "repo", "do", "do_gen.go")
}
func (a *App) doTableNameFile() string {
	return path.Join(a.Path, "repo", "do", "tables.go")
}

func (a *App) doConverterGenFile() string {
	return path.Join(a.Path, "repo", "converter", "converter_gen.go")
}

func (a *App) typesGenFiles() []typesGenFile {
	typesFiles := make([]typesGenFile, 0)
	switch a.Type {
	case TypeAPI:
		cmdPath := path.Join(a.Path, "cmd")
		fileInfos, err := os.ReadDir(cmdPath)
		if err != nil {
			log.Fatal(err)
		}
		for _, fi := range fileInfos {
			if fi.IsDir() {
				pwd := path.Join(cmdPath, fi.Name())
				typesFiles = append(typesFiles, typesGenFile{
					EntryName:    fi.Name(),
					Entry:        pwd,
					EntryPkgPath: a.appPkgPath() + "/cmd/" + fi.Name(),
					Header: a.GenFileHeader("types", []string{
						"time",
					}),
				})
			}
		}
	case TypeBFF:
		typesFiles = append(typesFiles, typesGenFile{
			EntryName:    a.Name,
			Entry:        a.Path,
			EntryPkgPath: a.appPkgPath(),
			Header: a.GenFileHeader("types", []string{
				"time",
			}),
		})

	}

	return typesFiles
}

type typesGenFile struct {
	EntryName    string `json:"entryName"`
	Entry        string `json:"entry"`
	EntryPkgPath string `json:"entryPkgPath"`
	Header       []byte `json:"header"`
}

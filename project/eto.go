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
	"sort"
	"strings"
)

var modules = []MicroXSTList{}

func (a *App) ETO() {
	// 解析 entity
	ipr := a.scanEntity()
	xstList := ipr.GetStructList()
	// 模块
	if a.Type == TypeModule {
		modules = append(modules, MicroXSTList{
			Name:    a.Name,
			PkgPath: a.appPkgPath(),
			XSTList: xstList,
		})
	}
	// 生成type_def
	a.eTypeDef(xstList)
	// 生成do 及 converter
	a.edo(xstList)
	// 根据do生成相关
	a.doNext()
	// 生成初始types及converter
	a.eTypes(xstList)
}

func (a *App) eTypeDef(xstList []parser.XST) {
	buf := []byte(fmt.Sprintf(tpls.EntityTypeDefCodes, "entity", cfg.C.Project))
	for _, xst := range xstList {
		_b, err := a._typedef(xst)
		if err != nil {
			log.Fatal(err)
		}
		buf = append(buf, _b...)
	}
	filename := a.entityTypeDefGenFile()
	buf = a.format(buf, filename)
	err := tool_file.WriteFile(filename, buf)
	if err != nil {
		log.Printf("type_def gen [%s] write file err: %v \n", filename, err)
	}
	log.Printf("gen type_def file %s \n", filename)
}

func (a *App) edo(xstList []parser.XST) {
	var (
		bufd = a.GenFileHeader("do", []string{
			"time",
			"gorm.io/gorm",
		})
		bufc = a.GenFileHeader("converter", []string{
			fmt.Sprintf("%s/common/tools", cfg.C.Project),
			fmt.Sprintf("%s/common/core/log", cfg.C.Project),
			fmt.Sprintf("%s/common/tools/tool_time", cfg.C.Project),
			fmt.Sprintf("%s/entity", a.appPkgPath()),
			fmt.Sprintf("%s/repo/dbal/do", a.appPkgPath()),
		})
	)
	for _, xst := range xstList {
		_b, _bc, err := a._do(xst)
		if err != nil {
			log.Fatal(err)
		}
		bufd = append(bufd, _b...)
		bufc = append(bufc, _bc...)
	}
	filename := a.doGenFile()
	bufd = a.format(bufd, filename)
	err := tool_file.WriteFile(filename, bufd)
	if err != nil {
		log.Printf("do gen [%s] write file err: %v \n", filename, err)
	}
	filename = a.doConverterGenFile()
	bufc = a.format(bufc, filename)
	err = tool_file.WriteFile(filename, bufc)
	if err != nil {
		log.Printf("conv gen [%s] write file err: %v \n", filename, err)
	}
	log.Printf("gen do file %s \n", filename)
}

func (a *App) doNext() {
	// 解析 do
	ipr := a.scanDo()
	xstList := ipr.GetStructList()
	// 生成 type_def
	buf := []byte(fmt.Sprintf(tpls.EntityTypeDefCodes, "do", cfg.C.Project))
	for _, xst := range xstList {
		_b, err := a._typedef(xst)
		if err != nil {
			log.Fatal(err)
		}
		buf = append(buf, _b...)
	}
	filename := a.doTypeDefGenFile()
	buf = a.format(buf, filename)
	err := tool_file.WriteFile(filename, buf)
	if err != nil {
		log.Printf("type_def gen [%s] write file err: %v \n", filename, err)
	}
	log.Printf("gen type_def file %s \n", filename)

	// 生成dao
	buf = []byte(fmt.Sprintf(tpls.DaoHeaderCodes, "dao", a.appPkgPath()))
	for _, xst := range xstList {
		_b, err := a._dao(xst)
		if err != nil {
			log.Fatal(err)
		}
		buf = append(buf, _b...)
	}
	filename = a.doDaoGenFile()
	buf = a.format(buf, filename)
	err = tool_file.WriteFile(filename, buf)
	if err != nil {
		log.Printf("dao gen [%s] write file err: %v \n", filename, err)
	}
	log.Printf("gen dao file %s \n", filename)
}

func (a *App) modulesTypes(tf typesGenFile) {
	// 模块注册模式 显性控制需要引入的模块
	// 读取micro配置
	usedMods := make(map[string]struct{})
	mainDoc, err := parser.MainFunDoc(path.Join(tf.Entry, "main.go"))
	if err == nil && mainDoc != "" {
		list := a.MicroParse(mainDoc)
		for _, s := range list {
			usedMods[s] = struct{}{}
		}
	}
	// module下的实体
	for _, module := range modules {
		if _, ok := usedMods[module.Name]; !ok {
			continue
		}
		var (
			xstMaps = map[string][]parser.XST{}

			bufc = a.GenFileHeaderWithAsName("converter", []string{
				fmt.Sprintf("%s/common/tools", cfg.C.Project),
				fmt.Sprintf("%s/common/core/log", cfg.C.Project),
				fmt.Sprintf("%s/common/tools/tool_time", cfg.C.Project),
				fmt.Sprintf("%s/entity", module.PkgPath),
			}, map[string]string{
				"types": fmt.Sprintf("%s/types/%s", tf.EntryPkgPath, "micro_"+module.Name),
			})
		)
		for _, xst := range module.XSTList {
			if _, ok := xstMaps[xst.File]; !ok {
				xstMaps[xst.File] = make([]parser.XST, 0)
			}
			xstMaps[xst.File] = append(xstMaps[xst.File], xst)
		}

		for _file, xstList := range xstMaps {
			var (
				_, _fname = path.Split(_file)
				bufd      = a.GenFileHeaderAllowEdit("micro_"+module.Name, []string{
					"time",
				})
			)

			ipts, err := parser.Scan(path.Join(tf.Entry, "types", "micro_"+module.Name), parser.ParseTypeWatch)
			if err != nil {
				log.Fatal(err)
			}

			for _, xst := range xstList {
				oldXst := ipts.StructList[xst.Name]
				_b, _bc, err := a._types(xst, oldXst, "json", "Micro"+tool_str.ToUFirst(module.Name))
				if err != nil {
					log.Fatal(err)
				}
				bufd = append(bufd, _b...)
				bufc = append(bufc, _bc...)
			}
			os.MkdirAll(path.Join(tf.Entry, "types", "micro_"+module.Name), 0777)
			filename := path.Join(tf.Entry, "types", "micro_"+module.Name, fmt.Sprintf("entity_%s", _fname))
			bufd = a.format(bufd, filename)
			err = tool_file.WriteFile(filename, bufd)
			if err != nil {
				log.Printf("types[micro] gen [%s] write file err: %v \n", filename, err)
			}
			log.Printf("gen types[micro] file %s \n", filename)
		}

		filename := path.Join(tf.Entry, "converter", "micro_"+module.Name+"_converter_gen.go")
		bufc = a.format(bufc, filename)
		err = tool_file.WriteFile(filename, bufc)
		if err != nil {
			log.Printf("types_conv gen [%s] write file err: %v \n", filename, err)
		}
		log.Printf("gen types_conv file %s \n", filename)
	}
}

func (a *App) eTypes(xstList []parser.XST) {
	if a.Type != TypeAPI {
		return
	}

	for _, tf := range a.typesGenFiles() {
		if !tool_file.Exists(path.Join(tf.Entry, "types")) || !tool_file.Exists(path.Join(tf.Entry, "converter")) {
			continue
		}
		// modules
		a.modulesTypes(tf)
		// business
		var (
			xstMaps = map[string][]parser.XST{}
			bufc    = a.GenFileHeader("converter", []string{
				fmt.Sprintf("%s/common/tools", cfg.C.Project),
				fmt.Sprintf("%s/common/core/log", cfg.C.Project),
				fmt.Sprintf("%s/common/tools/tool_time", cfg.C.Project),
				fmt.Sprintf("%s/entity", a.appPkgPath()),
				fmt.Sprintf("%s/types", tf.EntryPkgPath),
			})
		)
		for _, xst := range xstList {
			if _, ok := xstMaps[xst.File]; !ok {
				xstMaps[xst.File] = make([]parser.XST, 0)
			}
			xstMaps[xst.File] = append(xstMaps[xst.File], xst)
		}

		ipts, err := parser.Scan(path.Join(tf.Entry, "types"), parser.ParseTypeWatch)
		if err != nil {
			log.Fatal(err)
		}

		for _file, xstList := range xstMaps {
			var (
				_, _fname = path.Split(_file)
				bufd      = a.GenFileHeaderAllowEdit("types", []string{
					"time",
				})
			)
			for _, xst := range xstList {

				oldXst := ipts.StructList[xst.Name]
				_b, _bc, err := a._types(xst, oldXst, "json", "")
				if err != nil {
					log.Fatal(err)
				}
				bufd = append(bufd, _b...)
				bufc = append(bufc, _bc...)
			}
			filename := path.Join(tf.Entry, "types", fmt.Sprintf("entity_%s", _fname))
			bufd = a.format(bufd, filename)
			err := tool_file.WriteFile(filename, bufd)
			if err != nil {
				log.Printf("types gen [%s] write file err: %v \n", filename, err)
			}
			log.Printf("gen types file %s \n", filename)
		}
		filename := path.Join(tf.Entry, "converter", "entity_converter_gen.go")
		bufc = a.format(bufc, filename)
		err = tool_file.WriteFile(filename, bufc)
		if err != nil {
			log.Printf("types_conv gen [%s] write file err: %v \n", filename, err)
		}
		log.Printf("gen types_conv file %s \n", filename)

	}

}

func (a *App) _types(xst parser.XST, oldXst parser.XST, tagName string, nameMark string) ([]byte, []byte, error) {

	gio := tpls.IO{
		Name:   xst.Name,
		Fields: make([]tpls.IoField, 0),
	}
	fieldList := make([]parser.XField, 0)
	for _, field := range xst.FieldList {
		fieldList = append(fieldList, field)
	}
	sort.SliceStable(fieldList, func(i, j int) bool {
		return fieldList[i].Idx < fieldList[j].Idx
	})

	for _, field := range fieldList {
		tagJSON := field.GetTag("json")
		tagIO := field.GetTag(tagName)
		tags := ""

		if tagJSON == nil {
			continue
		}
		if tagJSON != nil {
			tags = fmt.Sprintf("`json:\"%s\"`", tagJSON.Name)
			if tagJSON.Name == "-" {
				tags = ""
			}
		}
		if tagIO != nil {
			if tagIO.Txt == "-" {
				continue
			}
			if tagIO.Txt != "" {
				tagJSON.Name = tagIO.Name
			}
		}

		type2 := ""
		type2Entity := false

		fType := field.Type
		switch field.SType {
		case 1:
			type2 = strings.Replace(field.Type, "*", "", 1)
			if strings.Contains(field.Type, "time.Time") {
				field.SType = parser.STypeTime
				fType = "string"
			}
		case 2:
			type2 = strings.Replace(field.Type, "[]", "", 1)
			type2 = strings.Replace(type2, "*", "", 1)
			if tool_str.UFirst(type2) {
				type2Entity = true
			}
		}
		gio.Fields = append(gio.Fields, tpls.IoField{
			Name:        field.Name,
			Type:        fType,
			Type2:       type2,
			Type2Entity: type2Entity,
			SType:       field.SType,
			Tag:         tags,
			Comment:     field.Comment,
		})
	}

	if len(gio.Fields) == 0 {
		return nil, nil, nil
	}

	convBuf, err := a._ioConv(gio, nameMark)
	if err != nil {
		return nil, nil, err
	}

	// 自定义字段
	for _, field := range oldXst.FieldList {
		if _, ok := xst.FieldList[field.Name]; !ok {
			tagJSON := field.GetTag("json")
			tagIO := field.GetTag(tagName)
			tags := ""

			if tagJSON == nil {
				continue
			}
			if tagJSON != nil {
				tags = fmt.Sprintf("`json:\"%s\"`", tagJSON.Name)
				if tagJSON.Name == "-" {
					tags = ""
				}
			}
			if tagIO != nil {
				if tagIO.Txt == "-" {
					continue
				}
				if tagIO.Txt != "" {
					tagJSON.Name = tagIO.Name
				}
			}

			type2 := ""
			type2Entity := false

			fType := field.Type
			switch field.SType {
			case 1:
				type2 = strings.Replace(field.Type, "*", "", 1)
				if strings.Contains(field.Type, "time.Time") {
					field.SType = parser.STypeTime
					fType = "string"
				}
			case 2:
				type2 = strings.Replace(field.Type, "[]", "", 1)
				type2 = strings.Replace(type2, "*", "", 1)
				if tool_str.UFirst(type2) {
					type2Entity = true
				}
			}
			gio.Fields = append(gio.Fields, tpls.IoField{
				Name:        field.Name,
				Type:        fType,
				Type2:       type2,
				Type2Entity: type2Entity,
				SType:       field.SType,
				Tag:         tags,
				Comment:     field.Comment,
			})
		}
	}

	buf, err := gio.Execute()
	if err != nil {
		return nil, nil, err
	}

	return buf, convBuf, nil
}

func (a *App) _ioConv(gio tpls.IO, nameMark string) ([]byte, error) {
	for idx, item := range gio.Fields {
		if item.Name == "" {
			gio.Fields[idx].Name = item.Type2
		}
	}
	convGen := tpls.IoConv{
		Name:     gio.Name,
		NameMark: nameMark,
		Fields:   gio.Fields,
	}
	buf, err := convGen.Execute()
	if err != nil {
		return nil, err
	}
	return buf, nil
}

func (a *App) _do(xst parser.XST) ([]byte, []byte, error) {
	gdo := tpls.Do{
		Name:     xst.Name,
		Fields:   make([]tpls.DoField, 0),
		DeleteAT: false,
	}
	if xst.NoDeleteAT {
		gdo.DeleteAT = false
	}
	fieldList := make([]parser.XField, 0)
	for _, field := range xst.FieldList {
		fieldList = append(fieldList, field)
	}
	sort.SliceStable(fieldList, func(i, j int) bool {
		return fieldList[i].Idx < fieldList[j].Idx
	})
	for _, field := range fieldList {
		tagDesc := field.GetTag("db")
		if tagDesc != nil && tagDesc.Txt != "-" {
			tag := tagDesc.Txt
			convSlice := false
			isPoint := false
			type2 := ""
			if tagDesc.Opts != nil && len(tagDesc.Opts) > 0 {
				if v, ok := tagDesc.Opts["conv"]; ok {
					tagConv := fmt.Sprintf("conv:%s", v)
					convSlice = true
					tag = strings.Replace(tag, tagConv+";", "", 1)
					tag = strings.Replace(tag, tagConv, "", 1)
				}
			}
			tags := fmt.Sprintf("`db:\"%s\" gorm:\"column:%s\"`", tagDesc.Name, tag)
			fType := field.Type
			switch field.SType {
			case 1:
				type2 = strings.Replace(field.Type, "*", "", 1)
				if strings.Contains(field.Type, "time.Time") {
					field.SType = parser.STypeTime
				} else {
					if strings.Index(field.Type, "*") == 0 {
						isPoint = true
					}
				}

			case 2:
				type2 = strings.Replace(field.Type, "[]", "", 1)
				if strings.Contains(type2, "[]") || strings.Index(type2, ".") > 0 {
					convSlice = false
				}
				if tool_str.UFirst(type2) || strings.Contains(type2, "map") {
					convSlice = false
				}
				fType = AddEntityPkg(fType)
			case 3:
				fType = AddEntityPkg(fType)
			}

			gdo.Fields = append(gdo.Fields, tpls.DoField{
				Name:      field.Name,
				Type:      fType,
				Type2:     type2,
				SType:     field.SType,
				Tag:       tags,
				ConvSlice: convSlice,
				IsPoint:   isPoint,
				Comment:   field.Comment,
			})
		}
	}
	if len(gdo.Fields) == 0 {
		return nil, nil, nil
	}
	buf, err := gdo.Execute()
	if err != nil {
		return nil, nil, err
	}

	convGen := tpls.DoConv{
		Name:   gdo.Name,
		Fields: gdo.Fields,
	}
	buf2, err := convGen.Execute()
	if err != nil {
		return nil, nil, err
	}

	return buf, buf2, nil
}

func (a *App) _dao(xst parser.XST) ([]byte, error) {
	var (
		pkName = ""
		pkType = ""
		pkCol  = ""
	)
	for _, field := range xst.FieldList {
		tag := field.GetTag("gorm")
		if tag != nil && strings.Contains(tag.Txt, "primaryKey") {
			pkName = field.Name
			pkType = field.Type
			pkCol = tag.Name
		}
	}
	tGen := tpls.Dao{
		EntityName:     xst.Name,
		DaoName:        strings.TrimSuffix(xst.Name, "Do") + "Dao",
		EntityListName: fmt.Sprintf("%sList", xst.Name),
		TableName:      fmt.Sprintf("do.TableName%s", xst.Name),
		PkName:         pkName,
		PkType:         pkType,
		PkCol:          pkCol,
	}
	buf, err := tGen.Execute()
	if err != nil {
		return buf, err
	}
	return buf, nil
}

func (a *App) _typedef(xst parser.XST) ([]byte, error) {
	tGen := tpls.EntityTypeMap{
		ProjectName:    cfg.C.Project,
		EntityName:     xst.Name,
		EntityListName: fmt.Sprintf("%sList", xst.Name),
		Field:          make([]tpls.Field, 0),
		HasCreator:     false,
		CreatorName:    "",
	}
	fieldList := make([]parser.XField, 0)
	for _, field := range xst.FieldList {
		fieldList = append(fieldList, field)
	}
	sort.SliceStable(fieldList, func(i, j int) bool {
		return fieldList[i].Idx < fieldList[j].Idx
	})
	feList := make([]tpls.Field, 0)

	for _, field := range fieldList {
		_type := field.Type
		tags := strings.Trim(field.Tag, "`")
		tagsMap := parseFieldTagMap(tags)
		dbTag := tagsMap["db"]
		if dbTag != "" && strings.Contains(dbTag, ";") {
			dbTag = strings.Split(dbTag, ";")[0]
		}
		if dbTag == "create_time" || dbTag == "update_time" || dbTag == "id" || dbTag == "deleted_at" {
			dbTag = ""
		}
		fe := tpls.Field{
			Field:           field.Name,
			FieldTag:        tags,
			FieldEscapedTag: fmt.Sprintf("%q", tags),
			FieldTagMap:     tagsMap,
			DBTag:           dbTag,
			Type:            _type,
			UseJSON:         false,
			NamedType:       "",
			TypeInName:      "",
			GenSliceFunc:    true,
			Nullable:        false,
			Comparable:      false,
		}
		if field.SType != 0 && field.SType != 4 {
			fe.UseJSON = true
		}
		if strings.Index(_type, "*") == 0 || field.SType >= 2 || _type == "interface{}" {
			fe.Nullable = true
		} else {
			fe.Comparable = true
		}
		switch _type {
		case "int":
			fe.TypeInName = "Int"
		case "int32":
			fe.TypeInName = "Int32"
		case "int64":
			fe.TypeInName = "Int64"
		case "string":
			fe.TypeInName = "String"
		default:
			fe.GenSliceFunc = false
		}
		feList = append(feList, fe)
	}

	tGen.Field = feList
	return tGen.Execute()
}

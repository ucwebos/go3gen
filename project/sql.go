package project

import (
	"fmt"
	"github.com/ucwebos/go3gen/project/parser"
	"github.com/ucwebos/go3gen/project/tpls"
	"github.com/ucwebos/go3gen/utils"
	"github.com/xbitgo/core/tools/tool_file"
	"log"
	"path"
	"sort"
	"strings"
	"time"
)

func (a *App) GenSql(dsn string) {
	var (
		db    = utils.GetDB(dsn)
		doDir = path.Join(a.Path, "repo", "dbal", "do")
	)
	fmt.Println(dsn)
	ipr, err := parser.Scan(doDir, parser.ParseTypeDo)
	if err != nil {
		log.Fatalf("genSql: parse dir[%s], err: %v", doDir, err)
	}
	fmt.Println(doDir)
	for s, xst := range ipr.StructList {
		if v, ok := ipr.ConstStrList["TableName"+s]; ok {
			err = a.createTableSQL(db, v, xst)
			if err != nil {
				panic(err)
			}
		}
	}
}

func (a *App) createTableSQL(db *utils.DB, tableName string, xst parser.XST) error {
	filename := fmt.Sprintf("%s/%s_create.sql", path.Join(a.Path, "repo", "dbal", "sql"), tableName)
	fmt.Println(filename)
	createSQL := db.TableCreateSQL(tableName)
	if createSQL != "" {
		tool_file.WriteFile(filename, []byte(createSQL))
		a.modifySQL(db, tableName, xst)
	} else {
		genSql := a.toGenSQL(tableName, xst)
		createSQL, err := genSql.CreateTable()
		if err != nil {
			return err
		}
		tool_file.WriteFile(filename, createSQL)
	}
	return nil
}

func (a *App) toGenSQL(tableName string, xst parser.XST) tpls.GenSQL {
	fieldList := make([]parser.XField, 0)
	for _, field := range xst.FieldList {
		fieldList = append(fieldList, field)
	}
	sort.SliceStable(fieldList, func(i, j int) bool {
		return fieldList[i].Idx < fieldList[j].Idx
	})
	primaryKey := ""
	last := ""
	sFields := make([]tpls.SQLField, 0)
	for _, field := range fieldList {
		dbTag := field.GetTag("db")
		if dbTag != nil {
			name := utils.AddStrSqlC(dbTag.Name, "`")
			if sf, ok := tpls.SpecialField[dbTag.Name]; ok {
				sf.TableName = utils.AddStrSqlC(tableName, "`")
				sf.After = last
				sf.SrcName = dbTag.Name
				sFields = append(sFields, sf)
				last = name
				continue
			}
			if strings.Contains(field.Tag, "primaryKey") {
				primaryKey = name
			}
			tt := tpls.TypeMap[strings.TrimPrefix(field.Type, "*")]
			sf := tpls.SQLField{
				TableName: utils.AddStrSqlC(tableName, "`"),
				After:     last,
				Name:      name,
				SrcName:   dbTag.Name,
				Type:      tt.Type,
				DataType:  tt.DataType,
				Default:   tt.Default,
				Comment:   field.Comment,
				NotNull:   "NOT NULL",
			}
			sFields = append(sFields, sf)
			last = name
		}
	}
	genSql := tpls.GenSQL{
		TableName:  utils.AddStrSqlC(tableName, "`"),
		PrimaryKey: primaryKey,
		Fields:     sFields,
	}
	return genSql
}

func (a *App) modifySQL(db *utils.DB, tableName string, xst parser.XST) error {
	columns, _ := db.TableColumns(tableName)
	genSql := a.toGenSQL(tableName, xst)
	addColumns := make([]tpls.SQLField, 0)
	for _, field := range genSql.Fields {
		if _, ok := columns[field.SrcName]; !ok {
			addColumns = append(addColumns, field)
		}
	}
	if len(addColumns) > 0 {
		filename := fmt.Sprintf("%s/%s_column_add_%s.sql", path.Join(a.Path, "repo", "dbal", "sql"), tableName, time.Now().Format("200601021504"))
		genSql.Fields = addColumns
		createSQL, err := genSql.AddColumns()
		if err != nil {
			fmt.Println(err)
			return err
		}
		tool_file.WriteFile(filename, createSQL)
	}
	return nil
}

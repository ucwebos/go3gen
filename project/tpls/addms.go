package tpls

import (
	"bytes"
	"text/template"
)

const AddMsCfgTpl = `package config

import (
	"context"

	"github.com/redis/go-redis/v9"

	"{{.Project}}/common/core/cfg"
	"{{.Project}}/common/core/log"
	"{{.Project}}/common/lib/db"
)

type Config struct {
	Env        string        ` + "`json:\"env\" yaml:\"env\"`" + `
	DB         *db.Config        ` + "`json:\"db\" yaml:\"db\"`" + `
	Redis      *cfg.Redis        ` + "`json:\"redis\" yaml:\"db\"`" + `
}

var CCfg = &Config{}

var dbInstance *db.Wrapper
var redisInstance *redis.Client

func MustInit() {
	cfg.MustLoad("micro-{{.Name}}", CCfg)
	// DB
	_db, err := db.NewWrapper(CCfg.DB)
	if err != nil {
		log.Panicf("micro-{{.Name}} NewDBWrapper error: %v, config: %+v", err, CCfg.DB)
	}
	dbInstance = _db
	// redis
	redisInstance = redis.NewClient(CCfg.Redis.ToRedisOption())
	_rs := redisInstance.Ping(context.Background())
	if err = _rs.Err(); err != nil {
		log.Panicf("micro-{{.Name}} NewRedis error: %v, config: %+v", err, CCfg.Redis)
	}
}

func GetDB() *db.Wrapper {
	return dbInstance
}

func GetRedis() *redis.Client {
	return redisInstance
}`

const AddMsEntryTpl = `package {{.Name}}

import "{{.Project}}/micro/{{.Name}}/config"

func Init() {
	config.MustInit()
	// provider.Register{{.NameUF}}(New{{.NameUF}}())
}

// Test 测试方法 [{{.NameUF}}.Test]
func _gen() {

}
`

type AddMicro struct {
	Project string
	Name    string
	NameUF  string
}

func (s *AddMicro) Execute(t string) []byte {
	buf := new(bytes.Buffer)
	tmpl, err := template.New("AddMicro").Parse(t)
	if err != nil {
		return nil
	}
	if err := tmpl.Execute(buf, s); err != nil {
		return nil
	}
	return buf.Bytes()
}

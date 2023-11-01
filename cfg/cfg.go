package cfg

type Cfg struct {
	Project  string `json:"project"`
	RootPath string `json:"root_path"`
}

type DBSet struct {
	DB     string            `json:"db" yaml:"db"`
	DBMaps map[string]string `json:"db_maps" yaml:"db_maps"`
}

var C = &Cfg{}

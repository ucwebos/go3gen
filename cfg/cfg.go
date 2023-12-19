package cfg

type Cfg struct {
	Project   string            `json:"project"`
	RootPath  string            `json:"root_path"`
	AdminRoot string            `json:"admin_root" yaml:"admin_root"`
	DB        string            `json:"db" yaml:"db"`
	DBMaps    map[string]string `json:"db_maps" yaml:"db_maps"`
}

var C = &Cfg{}

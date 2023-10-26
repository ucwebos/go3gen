package cfg

type Cfg struct {
	Project  string `json:"project"`
	RootPath string `json:"root_path"`
}

var C = &Cfg{}

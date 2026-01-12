package easyproxy

type Config struct {
	Listen string  `yaml:"listen" validate:"required"`
	Routes []Route `yaml:"routes" validate:"required"`
}

type Route struct {
	Path   string `yaml:"path"   validate:"required"`
	Target string `yaml:"target" validate:"required"`
}

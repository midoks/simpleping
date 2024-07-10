package conf

import (
)

var (
	App struct {
		Version string `ini:"-"`

		Name      string
		BrandName string
		RunUser   string
		RunMode   string

		DataPath  string
	}



	Ping struct {
		Timeout  int64
		Ip       string
		Day      int64
	}
)

type DatabaseOpts struct {
	Type         string
	Host         string
	Name         string
	User         string
	Password     string
	SslMode      string `ini:"ssl_mode"`
	Path         string
	Prefix       string
	Charset      string
	Timezone     string
	MaxOpenConns int
	MaxIdleConns int
}

// Database settings
var Database DatabaseOpts
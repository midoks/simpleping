package conf

import (
	"os"
	"path/filepath"

	"github.com/pkg/errors"
	"gopkg.in/ini.v1"
)

var File *ini.File

// 判断所给路径文件/文件夹是否存在 
func PathExists(path string) (bool,error) {
	_ , err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err){
		return false, nil
	}
	return false, err
}

func makeConf(default_cfg string){
	cfg := ini.Empty()

	cfg.Section("").Key("app_name").SetValue("simpleping")
	cfg.Section("").Key("brand_name").SetValue("simpleping")
	cfg.Section("").Key("run_mode").SetValue("prod")

	cfg.Section("database").Key("type").SetValue("sqlite3")
	cfg.Section("database").Key("name").SetValue("simpleping")
	cfg.Section("database").Key("prefix").SetValue("ping_")
	cfg.Section("database").Key("charset").SetValue("utf8mb4")
	cfg.Section("database").Key("path").SetValue("data/simpleping.db3")
	cfg.Section("database").Key("max_open_conns").SetValue("30")
	cfg.Section("database").Key("max_idle_conns").SetValue("30")
	
	cfg.Section("ping").Key("timeout").SetValue("1")
	cfg.Section("ping").Key("ip").SetValue("127.0.0.1")
	cfg.Section("ping").Key("day").SetValue("7")

	cfg.Section("ping_mysql_slave").Key("off").SetValue("1")
	cfg.Section("ping_mysql_slave").Key("timeout").SetValue("1")
	cfg.Section("ping_mysql_slave").Key("ip").SetValue("127.0.0.1")
	cfg.Section("ping_mysql_slave").Key("port").SetValue("3306")
	cfg.Section("ping_mysql_slave").Key("user").SetValue("root")
	cfg.Section("ping_mysql_slave").Key("pass").SetValue("root")
	cfg.Section("ping_mysql_slave").Key("day").SetValue("7")

	os.MkdirAll(filepath.Dir(default_cfg), os.ModePerm)
	if err := cfg.SaveTo(default_cfg); err != nil {
		return
	}
}

func InitConf() error {
	default_cfg := "conf/app.conf"
	b , err := PathExists(default_cfg)
	if !b {
		makeConf(default_cfg)
	}

	File, err := ini.Load(default_cfg)
    if err != nil {
    	return errors.Wrapf(err, "parse '%s'", default_cfg)
    }

    File.NameMapper = ini.TitleUnderscore
	if err = File.Section(ini.DefaultSection).MapTo(&App); err != nil {
		return errors.Wrap(err, "mapping default section")
	}

	App.DataPath = ensureAbs(App.DataPath)

	// ***************************
	// ----- Database settings -----
	// ***************************
	if err = File.Section("database").MapTo(&Database); err != nil {
		return errors.Wrap(err, "mapping [database] section")
	}

	// ***************************
	// ----- Ping settings -----
	// ***************************
	if err = File.Section("ping").MapTo(&Ping); err != nil {
		return errors.Wrap(err, "mapping [ping] section")
	}

	// ***************************
	// ----- MySQLPing settings -----
	// ***************************
	if err = File.Section("ping_mysql_slave").MapTo(&MySQLPing); err != nil {
		return errors.Wrap(err, "mapping [ping_mysql_slave] section")
	}

	return nil
}
package main

import (
	"time"

	_ "github.com/mattn/go-sqlite3"
	"github.com/spf13/viper"
	"xorm.io/xorm"
)

// File is the stuff we need...
type File struct {
	Md5         string
	FileID      string
	FileSize    int
	OutputName  string `xorm:"varchar(200)"`
	DownloadURL string
	Created     time.Time `xorm:"created"`
	Updated     time.Time `xorm:"updated"`
}

// ConnectDB ...
func ConnectDB() (*xorm.Engine, error) {
	engine, err := xorm.NewEngine("sqlite3", viper.GetString("DbPath"))
	if err != nil {
		return nil, err
	}
	engine.ShowSQL(true)
	// engine.Logger().SetLevel(log.LOG_DEBUG)
	err = engine.Sync2(new(File))
	if err != nil {
		return nil, err
	}
	return engine, nil
}

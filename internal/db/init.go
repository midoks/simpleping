package db

import (
    // "database/sql"
    "fmt"
    "log"
    "os"
    "path/filepath"
    "strings"
    "time"

    "simpleping/internal/conf"

    "github.com/pkg/errors"
    "gorm.io/driver/sqlite"
    "gorm.io/gorm"
    "gorm.io/gorm/logger"

)

var (
    db  *gorm.DB
    err error
)

var Tables = []interface{}{ new(Ping) }

func InitDb() error {

    newLogger := logger.New(
        log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
        logger.Config{
            SlowThreshold: time.Second,   // 慢 SQL 阈值
            LogLevel:      logger.Silent, // Log level
            Colorful:      false,         // 禁用彩色打印
        },
    )
    dbPath := conf.Database.Path
    if strings.EqualFold(conf.Database.Path, "data/simpleping.db3") {
        dbPath = conf.App.DataPath + "/" + conf.Database.Path
    }
    os.MkdirAll(filepath.Dir(dbPath), os.ModePerm)

    db, err = gorm.Open(sqlite.Open(dbPath), &gorm.Config{
        Logger:                 newLogger,
        SkipDefaultTransaction: true,
        PrepareStmt:            true,
    })

    if err != nil {
        return errors.Wrap(err, "failed to connect database")
    }

    db.Exec("PRAGMA synchronous = OFF;")

    for _, table := range Tables {
        if db.Migrator().HasTable(table) {
            continue
        }

        name := strings.TrimPrefix(fmt.Sprintf("%T", table), "*db.")
        err = db.Migrator().AutoMigrate(table)
        if err != nil {
            return errors.Wrapf(err, "auto migrate %q", name)
        }
    }

    sqlDB, err := db.DB()
    if err != nil {
        return err
    }

    sqlDB.SetMaxIdleConns(conf.Database.MaxIdleConns)
    sqlDB.SetMaxOpenConns(conf.Database.MaxOpenConns)
    sqlDB.SetConnMaxLifetime(time.Hour)

    return nil
}

func TablePrefix(tn string) string {
    return fmt.Sprintf("%s%s", "sp_", tn)
}
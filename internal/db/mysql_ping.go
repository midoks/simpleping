package db

import (
	// "fmt"
	"time"
)

type MySQLPing struct {
	Id       int64  `gorm:"primaryKey"`
	Ip       string `gorm:"index:,size:50;comment:IP"`
	Value    int64 `gorm:"index:,comment:速度"`
	Created     time.Time `gorm:"autoCreateTime;index:,comment:创建时间"`
	CreatedUnix int64     `gorm:"autoCreateTime;index:,;comment:创建时间"`
}

func (MySQLPing) TableName() string {
	return TablePrefix("mysql_ping")
}

func AddMySQLData(ip string, val int64) (err error) {
	u := MySQLPing{
		Ip:ip,
		Value:val,
	}

	result := db.Create(&u)
	return result.Error
}

func DeleteMySQLData(day int64){
	ping := MySQLPing{}
	t := time.Now().Unix()  - day*86400
	// fmt.Println(t)
	db.Where("created_unix < ?", t).Delete(&ping)
}
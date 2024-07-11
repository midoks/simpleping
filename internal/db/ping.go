package db

import (
	// "fmt"
	"time"
)

type Ping struct {
	Id       int64  `gorm:"primaryKey"`
	Ip       string `gorm:"index:,size:50;comment:IP"`
	Speed    int64 `gorm:"index:,comment:速度"`
	Created     time.Time `gorm:"autoCreateTime;index:,comment:创建时间"`
	CreatedUnix int64     `gorm:"autoCreateTime;index:,;comment:创建时间"`
}

func (Ping) TableName() string {
	return TablePrefix("ping")
}

func AddPingData(ip string, speed int64) (err error) {
	u := Ping{
		Ip:ip,
		Speed:speed,
	}

	result := db.Create(&u)
	return result.Error
}

func DeletePingData(day int64){
	ping := Ping{}
	t := time.Now().Unix()  - day*86400
	// fmt.Println(t)
	db.Where("created_unix < ?", t).Delete(&ping)
}
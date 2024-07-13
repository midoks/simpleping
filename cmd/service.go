package cmd

import (
	"fmt"
	"time"
	"runtime"
	"sync"
	"strings"
	"strconv"


	"simpleping/internal/db"
	"simpleping/internal/conf"

	"gorm.io/driver/mysql"
  	"gorm.io/gorm"

	"github.com/go-ping/ping"
	"github.com/urfave/cli"
	"github.com/jakecoffman/cron"
)


var Service = cli.Command{
	Name:        "service",
	Usage:       "This command starts web service",
	Description: `Start Service`,
	Action:      ServiceRun,
	Flags: []cli.Flag{
	},
}

func pingIp(ip string, timeout int64, wg *sync.WaitGroup) error {
	pinger, err := ping.NewPinger(ip)

	if err != nil {
		wg.Done()
		return err
	}

	pinger.OnFinish = func(stats *ping.Statistics) {
		db.AddPingData(ip, int64(stats.AvgRtt))
	}

	pinger.Size = 56
	pinger.Timeout = time.Duration(timeout) * time.Second
	err = pinger.Run()

	if err != nil {
		wg.Done()
		return err
	}
	wg.Done()
	return err
}

func autoPingIpRun(){
	var wg sync.WaitGroup
	ips := strings.Split(conf.Ping.Ip, ",")
	// fmt.Println(ips)
	for{
		for _, ip := range ips {
			wg.Add(1)
			go pingIp(ip, conf.Ping.Timeout, &wg)
		}
		wg.Wait()
		// time.Sleep(1*time.Second)
	}
}

func pingMySQLIp(ip string, port int, user string, pass string, wg *sync.WaitGroup) error {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/?charset=utf8mb4&parseTime=True&loc=Local", user, pass, ip, port)
	// fmt.Println(dsn);
  	mydb, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
  	sqlDB, err := mydb.DB()
    if err != nil {
    	wg.Done()
        return err
    }

    sqlDB.SetMaxIdleConns(3)
    sqlDB.SetMaxOpenConns(3)

  	if err == nil {
  		result := map[string]interface{}{}
  		mydb.Raw("select version() as ver").Scan(&result)
  		ver := result["ver"].(string)
  		ver_m_list := strings.Split(ver, ".")
  		ver_big, _ := strconv.Atoi(ver_m_list[0])
  		if (ver_big > 5){
  			replica := map[string]interface{}{}
  			mydb.Raw("SHOW REPLICA STATUS").Scan(&replica)
  			sbm := replica["Seconds_Behind_Master"].(uint64)
  			// sbm_int, _ := strconv.ParseInt(sbm,10,64)
  			// fmt.Println("sbm_int",sbm,sbm_int)
  			db.AddMySQLData(ip, int64(sbm))
  		} else{
  			slave := map[string]interface{}{}
  			mydb.Raw("SHOW SLAVE STATUS").Scan(&slave)
  			sbm := slave["Seconds_Behind_Master"].(uint64)
  			// sbm_int, _ := strconv.ParseInt(sbm,10,64)
  			// fmt.Println("sbm_int",sbm)
  			db.AddMySQLData(ip, int64(sbm))
  		}
  	}

  	sqlDB.Close()
  	wg.Done()
	return err
}

func autoMySQLPingIpRun(){
	var wg sync.WaitGroup
	ip := conf.MySQLPing.Ip
	port := conf.MySQLPing.Port
	user := conf.MySQLPing.User
	pass := conf.MySQLPing.Pass
	off := conf.MySQLPing.Off

	if off == 0 {
		for{
			// fmt.Println("run:",ip, port, user, pass)
			wg.Add(1)
			go pingMySQLIp(ip, port, user, pass, &wg)
			wg.Wait()
			time.Sleep(1*time.Second)
		}
	}	
}

func deleteExpiredPing(){
	//删除过期数据
	cron_task := cron.New()
	cron_task.AddFunc("*/30 * * * * *", func() {
		db.DeletePingData(conf.Ping.Day)
	}, "delete_expired_ping")

	//删除过期mysql数据
	cron_task.AddFunc("*/30 * * * * *", func() {
		db.DeleteMySQLData(conf.MySQLPing.Day)
	}, "delete_expired_mysqlping")
	cron_task.Start()	
}


func ServiceRun(c *cli.Context) error {
	runtime.GOMAXPROCS(runtime.NumCPU())

	err := conf.InitConf()
	if err !=nil{
		fmt.Println(err)
		return err
	}

	err = db.InitDb()
	if err !=nil{
		fmt.Println(err)
		return err
	}

	go deleteExpiredPing()

	//mysql ping 
	go autoMySQLPingIpRun()
	//ip ping
	autoPingIpRun()
	return nil
}



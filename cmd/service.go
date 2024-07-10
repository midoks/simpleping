package cmd

import (
	"fmt"
	"time"
	"runtime"
	"sync"
	"strings"


	"simpleping/internal/db"
	"simpleping/internal/conf"

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

func deleteExpiredPing(){
	//删除过期数据
	cron_task := cron.New()
	cron_task.AddFunc("*/30 * * * * *", func() {
		db.DeletePingData(conf.Ping.Day)
	}, "delete_expired_ping")
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

	deleteExpiredPing()

	autoPingIpRun()
	return nil
}



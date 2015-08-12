package main

import (
	"flag"
	"fmt"
	"github.com/open-falcon/agent/cron"
	"github.com/open-falcon/agent/funcs"
	"github.com/open-falcon/agent/g"
	"github.com/open-falcon/agent/http"
	"os"
)

func main() {

	cfg := flag.String("c", "cfg.json", "configuration file")
	version := flag.Bool("v", false, "show version")
	check := flag.Bool("check", false, "check collector")

	flag.Parse()

	if *version {
		fmt.Println(g.VERSION)
		os.Exit(0)
	}

	// 辅助功能：人工View一下Collector的状态
	if *check {
		funcs.CheckCollector()
		os.Exit(0)
	}

	g.ParseConfig(*cfg)

	g.InitRootDir()

	// 获取本地的Ip
	g.InitLocalIps()

	// 创建Hb Client& Transfer Clients
	g.InitRpcClients()

	funcs.BuildMappers()

	// 启动CPU, Disk的数据统计
	go cron.InitDataHistory()

	// 告诉Hbs，Agent的状态
	cron.ReportAgentStatus()

	// 从Hbs获取plugins, 统计数据
	cron.SyncMinePlugins()
	cron.SyncBuiltinMetrics()
	cron.SyncTrustableIps()

	// 正式开始Collect
	cron.Collect()

	go http.Start()

	select {}

}

package cron

import (
	"github.com/open-falcon/agent/g"
	"github.com/open-falcon/common/model"
	"log"
	"strconv"
	"strings"
	"time"
)

func SyncBuiltinMetrics() {
	// 通过心跳来获取所需要监控的指标
	if g.Config().Heartbeat.Enabled && g.Config().Heartbeat.Addr != "" {
		go syncBuiltinMetrics()
	}
}

//
// 从hbs Server获取当前Agent需要监控的指标；那么本地的Agent又是如何利用去获取这些指标的呢?
//
func syncBuiltinMetrics() {

	var timestamp int64 = -1
	var checksum string = "nil"

	duration := time.Duration(g.Config().Heartbeat.Interval) * time.Second

	for {
	REST:
		time.Sleep(duration)

		var ports = []int64{}
		var paths = []string{}
		var procs = make(map[string]map[int]string)

		hostname, err := g.Hostname()
		if err != nil {
			goto REST
		}

		req := model.AgentHeartbeatRequest{
			Hostname: hostname,
			Checksum: checksum,
		}

		var resp model.BuiltinMetricResponse

		// 参考: hbs项目中的定义
		//    Agent#BuiltinMetrics(AgentHeartbeatRequest,BuiltinMetricResponse)
		// 两个项目通过: open-falcon/common 来定义公共的Model
		//
		err = g.HbsClient.Call("Agent.BuiltinMetrics", req, &resp)
		if err != nil {
			log.Println("ERROR:", err)
			goto REST
		}

		// 如果timestamp和checksum没有改变，则继续等待
		if resp.Timestamp <= timestamp {
			goto REST
		}

		if resp.Checksum == checksum {
			goto REST
		}

		timestamp = resp.Timestamp
		checksum = resp.Checksum

		for _, metric := range resp.Metrics {
			// 读取到 "net.port.listen"
			if metric.Metric == "net.port.listen" {
				arr := strings.Split(metric.Tags, "=")
				if len(arr) != 2 {
					continue
				}

				// Parse Port: 64-bit(8-byte长的int)
				if port, err := strconv.ParseInt(arr[1], 10, 64); err == nil {
					ports = append(ports, port)
				} else {
					log.Println("metrics ParseInt failed:", err)
				}

				continue
			}

			if metric.Metric == "du.bs" {
				arr := strings.Split(metric.Tags, "=")
				if len(arr) != 2 {
					continue
				}

				paths = append(paths, strings.TrimSpace(arr[1]))
				continue
			}

			if metric.Metric == "proc.num" {
				arr := strings.Split(metric.Tags, ",")

				tmpMap := make(map[int]string)

				for i := 0; i < len(arr); i++ {
					if strings.HasPrefix(arr[i], "name=") {
						tmpMap[1] = strings.TrimSpace(arr[i][5:])
					} else if strings.HasPrefix(arr[i], "cmdline=") {
						tmpMap[2] = strings.TrimSpace(arr[i][8:])
					}
				}

				procs[metric.Tags] = tmpMap
			}
		}

		// 更新package g中的变量
		g.SetReportPorts(ports)
		g.SetReportProcs(procs)
		g.SetDuPaths(paths)

	}
}

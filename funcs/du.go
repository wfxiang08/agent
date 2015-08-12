package funcs

import (
	"github.com/open-falcon/agent/g"
	"github.com/open-falcon/common/model"
	"github.com/toolkits/sys"
	"log"
	"strconv"
	"strings"
)

// 统计给定路径的Du,大小为byte
func DuMetrics() (L []*model.MetricValue) {
	paths := g.DuPaths()
	for _, path := range paths {
		out, err := sys.CmdOutNoLn("du", "-bs", path)
		if err != nil {
			log.Println("du -bs", path, "fail", err)
			continue
		}

		arr := strings.Fields(out)
		if len(arr) == 1 {
			continue
		}

		size, err := strconv.ParseUint(arr[0], 10, 64)
		if err != nil {
			log.Println("cannot parse du -bs", path, "output")
			continue
		}

		// 统计结果：
		// metric, tags, value
		// Hostname在外部统一添加
		L = append(L, GaugeValue("du.bs", size, "path="+path))
	}

	return
}

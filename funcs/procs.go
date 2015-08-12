package funcs

import (
	"github.com/open-falcon/agent/g"
	"github.com/open-falcon/common/model"
	"github.com/toolkits/nux"
	"log"
	"strings"
)

func ProcMetrics() (L []*model.MetricValue) {

	reportProcs := g.ReportProcs()
	sz := len(reportProcs)
	if sz == 0 {
		return
	}

	ps, err := nux.AllProcs()
	if err != nil {
		log.Println(err)
		return
	}

	pslen := len(ps)

	for tags, m := range reportProcs {
		cnt := 0
		for i := 0; i < pslen; i++ {
			if is_a(ps[i], m) {
				cnt++
			}
		}

		L = append(L, GaugeValue("proc.num", cnt, tags))
	}

	return
}

// 两级匹配:
// Name必须一致，并且Cmdline必须包含val
func is_a(p *nux.Proc, m map[int]string) bool {
	// p:
	//  Pid
	//  Name 例如: usgi
	//  CmdLine: usgi-xsdfdsfds
	//
	// only one kv pair
	for key, val := range m {
		if key == 1 {
			// name
			if val != p.Name {
				return false
			}
		} else if key == 2 {
			// cmdline
			if !strings.Contains(p.Cmdline, val) {
				return false
			}
		}
	}
	return true
}

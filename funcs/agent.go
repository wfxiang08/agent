package funcs

import (
	"github.com/open-falcon/common/model"
)

// 获取Agent的状态: 能访问，就只有一个状态"活着"
func AgentMetrics() []*model.MetricValue {
	return []*model.MetricValue{GaugeValue("agent.alive", 1)}
}

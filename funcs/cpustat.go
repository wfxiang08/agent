package funcs

import (
	"github.com/open-falcon/common/model"
	"github.com/toolkits/nux"
	"sync"
)

const (
	historyCount int = 2
)

var (
	procStatHistory [historyCount]*nux.ProcStat
	psLock          = new(sync.RWMutex)
)

func UpdateCpuStat() error {
	ps, err := nux.CurrentProcStat()
	if err != nil {
		return err
	}

	// 通过固定长度的Buffer来记录: ProcStat的history的状态
	// 最多记住两个状态: historyCount
	psLock.Lock()
	defer psLock.Unlock()
	for i := historyCount - 1; i > 0; i-- {
		procStatHistory[i] = procStatHistory[i-1]
	}

	procStatHistory[0] = ps
	return nil
}

func deltaTotal() uint64 {
	if procStatHistory[1] == nil {
		return 0
	}
	return procStatHistory[0].Cpu.Total - procStatHistory[1].Cpu.Total
}

func CpuIdle() float64 {
	psLock.RLock()
	defer psLock.RUnlock()
	dt := deltaTotal()
	if dt == 0 {
		return 0.0
	}
	invQuotient := 100.00 / float64(dt)
	return float64(procStatHistory[0].Cpu.Idle-procStatHistory[1].Cpu.Idle) * invQuotient
}

func CpuUser() float64 {
	psLock.RLock()
	defer psLock.RUnlock()
	dt := deltaTotal()
	if dt == 0 {
		return 0.0
	}
	invQuotient := 100.00 / float64(dt)
	return float64(procStatHistory[0].Cpu.User-procStatHistory[1].Cpu.User) * invQuotient
}

func CpuNice() float64 {
	psLock.RLock()
	defer psLock.RUnlock()
	dt := deltaTotal()
	if dt == 0 {
		return 0.0
	}
	invQuotient := 100.00 / float64(dt)
	return float64(procStatHistory[0].Cpu.Nice-procStatHistory[1].Cpu.Nice) * invQuotient
}

func CpuSystem() float64 {
	psLock.RLock()
	defer psLock.RUnlock()
	dt := deltaTotal()
	if dt == 0 {
		return 0.0
	}
	invQuotient := 100.00 / float64(dt)
	return float64(procStatHistory[0].Cpu.System-procStatHistory[1].Cpu.System) * invQuotient
}

func CpuIowait() float64 {
	psLock.RLock()
	defer psLock.RUnlock()
	dt := deltaTotal()
	if dt == 0 {
		return 0.0
	}
	invQuotient := 100.00 / float64(dt)
	return float64(procStatHistory[0].Cpu.Iowait-procStatHistory[1].Cpu.Iowait) * invQuotient
}

func CpuIrq() float64 {
	psLock.RLock()
	defer psLock.RUnlock()
	dt := deltaTotal()
	if dt == 0 {
		return 0.0
	}
	invQuotient := 100.00 / float64(dt)
	return float64(procStatHistory[0].Cpu.Irq-procStatHistory[1].Cpu.Irq) * invQuotient
}

func CpuSoftIrq() float64 {
	psLock.RLock()
	defer psLock.RUnlock()
	dt := deltaTotal()
	if dt == 0 {
		return 0.0
	}
	invQuotient := 100.00 / float64(dt)
	return float64(procStatHistory[0].Cpu.SoftIrq-procStatHistory[1].Cpu.SoftIrq) * invQuotient
}

func CpuSteal() float64 {
	psLock.RLock()
	defer psLock.RUnlock()
	dt := deltaTotal()
	if dt == 0 {
		return 0.0
	}
	invQuotient := 100.00 / float64(dt)
	return float64(procStatHistory[0].Cpu.Steal-procStatHistory[1].Cpu.Steal) * invQuotient
}

func CpuGuest() float64 {
	psLock.RLock()
	defer psLock.RUnlock()
	dt := deltaTotal()
	if dt == 0 {
		return 0.0
	}
	invQuotient := 100.00 / float64(dt)
	return float64(procStatHistory[0].Cpu.Guest-procStatHistory[1].Cpu.Guest) * invQuotient
}

func CurrentCpuSwitches() uint64 {
	psLock.RLock()
	defer psLock.RUnlock()
	return procStatHistory[0].Ctxt
}

func CpuPrepared() bool {
	psLock.RLock()
	defer psLock.RUnlock()
	return procStatHistory[1] != nil
}

// 定义各种CPU Metrics
func CpuMetrics() []*model.MetricValue {
	if !CpuPrepared() {
		return []*model.MetricValue{}
	}

	cpuIdleVal := CpuIdle()
	idle := GaugeValue("cpu.idle", cpuIdleVal)
	busy := GaugeValue("cpu.busy", 100.0-cpuIdleVal)
	user := GaugeValue("cpu.user", CpuUser())
	nice := GaugeValue("cpu.nice", CpuNice())
	system := GaugeValue("cpu.system", CpuSystem())
	iowait := GaugeValue("cpu.iowait", CpuIowait())
	irq := GaugeValue("cpu.irq", CpuIrq())
	softirq := GaugeValue("cpu.softirq", CpuSoftIrq())
	steal := GaugeValue("cpu.steal", CpuSteal())
	guest := GaugeValue("cpu.guest", CpuGuest())
	switches := CounterValue("cpu.switches", CurrentCpuSwitches())
	return []*model.MetricValue{idle, busy, user, nice, system, iowait, irq, softirq, steal, guest, switches}
}

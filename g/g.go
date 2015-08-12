package g

import (
	"log"
	"runtime"
)

func init() {
	// 设置CPU
	runtime.GOMAXPROCS(runtime.NumCPU())
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)
}

package helpers

import (
	"fmt"
	"goscheduler/common/logger"
	"os"
	"path/filepath"
	"runtime/pprof"
	"time"
)

//GetDirectories - Gets the root and programs catalog
func GetDirectories(base string) (string, string, error) {

	prog, err := filepath.Abs(filepath.Dir(base))
	if err != nil {
		return "", "", err
	}
	root := filepath.Dir(prog)

	return root, prog, nil
}

//StartProfiler - starts profiler
func StartProfiler(log logger.AppLogger, profile string) {
	log.Info("Start profiler")
	f, _ := os.Create(fmt.Sprintf("%sprofile.pprof", profile))
	pprof.StartCPUProfile(f)

	go func() {
		time.Sleep(60 * time.Second)
		pprof.StopCPUProfile()
		f.Close()
		log.Info("profiler stoped")

	}()

}

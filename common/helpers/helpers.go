/*
Package helpers - contains miscelanous helpers methods
*/
package helpers

import (
	"fmt"
	"os"
	"overseer/common/logger"
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
	f, _ := os.Create(fmt.Sprintf("%s", profile))
	pprof.StartCPUProfile(f)

	go func() {
		time.Sleep(60 * time.Second)
		pprof.StopCPUProfile()
		f.Close()
		log.Info("profiler stoped")

	}()

}

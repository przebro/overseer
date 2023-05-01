/*
Package helpers - contains miscelanous helpers methods
*/
package helpers

import (
	"os"
	"path/filepath"
	"runtime/pprof"
	"time"

	"github.com/rs/zerolog"
)

// GetDirectories - Gets the root and programs catalog
func GetDirectories(base string) (string, string, error) {

	prog, err := filepath.Abs(filepath.Dir(base))
	if err != nil {
		return "", "", err
	}
	root := filepath.Dir(prog)

	return root, prog, nil
}

// StartProfiler - starts profiler
func StartProfiler(log zerolog.Logger, profile string) {
	log.Info().Msg("Start profiler")
	f, _ := os.Create(profile)
	pprof.StartCPUProfile(f)

	go func() {
		time.Sleep(60 * time.Second)
		pprof.StopCPUProfile()
		f.Close()
		log.Info().Msg("profiler stoped")

	}()

}

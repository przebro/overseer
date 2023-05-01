package logger

import (
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"
)

// fileRotateSync - rotating log
type fileRotateSync struct {
	file    *os.File
	limit   int
	writed  int
	dir     string
	fprefix string
	num     int
	loc     sync.Mutex
}

func newFile(dir, prefix string, num int) (*os.File, error) {

	fname := fmt.Sprintf("%s_%s_%04d.log", prefix, time.Now().Format("20060102_150405"), num)
	path := filepath.Join(dir, fname)

	if !filepath.IsAbs(path) {
		return nil, fmt.Errorf("filepath is not absolute:%s", path)

	}
	return os.Create(path)
}

// newRotateSync - creates a new fileRotateSync, it satisfies zapcore.WriteSyncer interface
func newRotateSync(dir, prefix string, limit int) (*fileRotateSync, error) {

	var file *os.File
	var err error

	if file, err = newFile(dir, prefix, 1); err != nil {
		return nil, err
	}

	rfile := &fileRotateSync{
		dir:     dir,
		fprefix: prefix,
		limit:   limit,
		num:     1,
		file:    file,
	}

	return rfile, nil
}

func (r *fileRotateSync) Write(p []byte) (n int, err error) {
	defer r.loc.Unlock()
	r.loc.Lock()

	if r.writed+len(p) > r.limit {

		r.num++
		r.writed = 0
		if r.file.Close(); err != nil {
			return 0, err
		}
		if r.file, err = newFile(r.dir, r.fprefix, r.num); err != nil {
			return 0, err
		}
	}
	r.writed += len(p)

	return r.file.Write(p)
}
func (r *fileRotateSync) Sync() error {
	defer r.loc.Unlock()
	r.loc.Lock()
	return r.file.Sync()

}

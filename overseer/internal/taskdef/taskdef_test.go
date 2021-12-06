package taskdef

import (
	"os"
	"path/filepath"
	"strings"
)

var defDirectory = "../../../def_test"

var managerPath string
var groupsDircetories int = 0

func init() {
	managerPath, _ = filepath.Abs(defDirectory)

	dirs, _ := os.ReadDir(managerPath)
	for _, n := range dirs {
		if n.IsDir() && !strings.HasPrefix(n.Name(), ".") {
			groupsDircetories++
		}

	}
}

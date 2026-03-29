package emitter

import (
	"fmt"
	"path/filepath"
	"runtime"
)

func repositoryCallerLocation(skip int) string {
	_, file, line, _ := runtime.Caller(skip)
	return fmt.Sprintf("%s:%d", filepath.Base(file), line)
}

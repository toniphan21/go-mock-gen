package meta

import (
	"fmt"
	"runtime"
	"strings"
)

func mockgenCaptureCleanStack(skip int) string {
	pc := make([]uintptr, 20)
	n := runtime.Callers(skip, pc)
	if n == 0 {
		return ""
	}

	pc = pc[:n]
	frames := runtime.CallersFrames(pc)
	var sb strings.Builder
	for {
		frame, more := frames.Next()

		if frame.Function == "" ||
			strings.Contains(frame.Function, "runtime.") ||
			strings.Contains(frame.Function, "testing.") {
			if !more {
				break
			}
			continue
		}

		fnName := frame.Function
		if lastDot := strings.LastIndex(fnName, "."); lastDot != -1 {
			fnName = fnName[lastDot+1:]
		}

		fnName = strings.TrimSuffix(fnName, ")")
		if lastOpen := strings.LastIndex(fnName, "("); lastOpen != -1 {
			fnName = fnName[lastOpen+1:]
		}

		if frame.File != "" && frame.Line != 0 {
			if _, err := fmt.Fprintf(&sb, "%s:%d: %s\n", frame.File, frame.Line, fnName); err != nil {
				continue
			}
		}

		if !more {
			break
		}
	}
	return sb.String()
}

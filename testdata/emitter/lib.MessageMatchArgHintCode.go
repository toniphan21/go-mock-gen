package emitter

import "fmt"

func repositoryMessageMatchArgHint() string {
	return fmt.Sprintf("\thint: check argument matching at %s\n\t\tor use STUB for fine-grained control", repositoryCallerLocation(3))
}

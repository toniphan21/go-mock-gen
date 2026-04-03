package mockgen

import (
	"golang.org/x/tools/go/packages"
	genlib "nhatp.com/go/gen-lib"
)

type Generator interface {
	Generate(pkg *packages.Package, configs []Config) error
}

func New(fileManager genlib.FileManager) Generator {
	impl := &generatorImpl{
		fileManager:     fileManager,
		emitter:         &DefaultEmitter{},
		pkgNameManagers: make(map[string]NameManager),
	}

	return impl
}

package mockgen

// For advanced usage beyond these helpers, see nhatp.com/go/gen-lib directly.

import (
	"golang.org/x/tools/go/packages"
	genlib "nhatp.com/go/gen-lib"
)

type Output = genlib.Output

type FileManager = genlib.FileManager
type EmitterContext = genlib.EmitterContext
type NameManager = genlib.NameManager

func NewFileManager(rootDir string, opts ...genlib.FileManagerOption) genlib.FileManager {
	return genlib.NewFileManager(rootDir, opts...)
}

func WithBinaryName(v string) genlib.FileManagerOption {
	return genlib.WithBinaryName(v)
}

func WithVersion(v string) genlib.FileManagerOption {
	return genlib.WithVersion(v)
}

func LoadPackages(dir string) ([]*packages.Package, error) {
	return genlib.LoadPackages(dir)
}

func NewEmitterContext(pkg *packages.Package, fm FileManager, gf *genlib.GenFile, varPrefix string) EmitterContext {
	return genlib.NewEmitterContext(pkg, fm, gf, varPrefix)
}

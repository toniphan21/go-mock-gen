package mockgen

import (
	"log/slog"

	"golang.org/x/tools/go/packages"
	genlib "nhatp.com/go/gen-lib"
)

type Generator interface {
	Generate(pkg *packages.Package, configs []Config) error
}

type generatorOptionFunc func(impl *generatorImpl)

func (f generatorOptionFunc) apply(impl *generatorImpl) {
	f(impl)
}

type Option interface {
	apply(generator *generatorImpl)
}

func New(fileManager genlib.FileManager, options ...Option) Generator {
	impl := &generatorImpl{
		fileManager:     fileManager,
		emitter:         &DefaultEmitter{},
		logger:          &DefaultLogger{Logger: slog.New(slog.DiscardHandler)},
		pkgNameManagers: make(map[string]NameManager),
	}

	for _, fn := range options {
		fn.apply(impl)
	}

	return impl
}

func WithLogger(logger Logger) Option {
	return generatorOptionFunc(func(impl *generatorImpl) {
		impl.logger = logger
	})
}

func WithLogPoints(points *LogPoints) Option {
	return WithLogger(NewLogger(points))
}

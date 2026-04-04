package mockgen

import (
	"fmt"
	"log/slog"

	"golang.org/x/tools/go/packages"
)

type GeneratedInfo struct {
	Config      Config
	Struct      string
	Constructor string
}

type Logger interface {
	Error(action string, err error)
	Parsed(pkg *packages.Package, parsed map[Config][]MethodInfo)
	StartGenerating(pkg *packages.Package, config Config)
	Generated(pkg *packages.Package, info GeneratedInfo)
}

type DefaultLogger struct {
	Logger *slog.Logger
}

func (l DefaultLogger) Error(action string, err error) {
	l.Logger.Error(fmt.Sprintf("%s: %s", action, err))
}

func (l DefaultLogger) Parsed(pkg *packages.Package, parsed map[Config][]MethodInfo) {}

func (l DefaultLogger) StartGenerating(pkg *packages.Package, config Config) {
	l.Logger.Info("start generating", slog.String("interface", config.InterfaceName))
}

func (l DefaultLogger) Generated(pkg *packages.Package, info GeneratedInfo) {
	l.Logger.Info("generated", slog.String("interface", info.Config.InterfaceName))
}

var _ Logger = (*DefaultLogger)(nil)

type LogPoints struct {
	Error           func(action string, err error)
	FilteredConfigs func(pkg *packages.Package, parsed map[Config][]MethodInfo)
	StartGenerating func(pkg *packages.Package, config Config)
	Generated       func(pkg *packages.Package, info GeneratedInfo)
}

type logPointsInternal struct {
	p *LogPoints
}

func (l logPointsInternal) Error(action string, err error) {
	if l.p.Error != nil {
		l.p.Error(action, err)
	}
}

func (l logPointsInternal) Parsed(pkg *packages.Package, parsed map[Config][]MethodInfo) {
	if l.p.FilteredConfigs != nil {
		l.p.FilteredConfigs(pkg, parsed)
	}
}

func (l logPointsInternal) StartGenerating(pkg *packages.Package, config Config) {
	if l.p.StartGenerating != nil {
		l.p.StartGenerating(pkg, config)
	}
}

func (l logPointsInternal) Generated(pkg *packages.Package, info GeneratedInfo) {
	if l.p.Generated != nil {
		l.p.Generated(pkg, info)
	}
}

var _ Logger = (*logPointsInternal)(nil)

// ---

func NewLogger(points *LogPoints) Logger {
	return &logPointsInternal{p: points}
}

package yaml

import (
	"bytes"
	"errors"
	"slices"
	"strings"

	yml "gopkg.in/yaml.v3"
	mockgen "nhatp.com/go/mock-gen"
)

const DefaultOutputFilename = "mockgen_test.go"

func makeError(lines ...string) error {
	return errors.New(strings.Join(lines, "\n"))
}

type Global struct {
	Version  string               `yaml:"version"`
	All      All                  `yaml:"all"`
	Packages map[string][]Package `yaml:"packages"`
}

func (g *Global) validate() error {
	if g.Version != "1" {
		return makeError(`unknown version, supported: "1"`)
	}
	if len(g.Packages) == 0 {
		return makeError(
			"no packages",
			"",
			"  hint: add packages configuration",
			"    packages:",
			"      github.com/you/project/path:",
			"        - Repository # mocked interface",
		)
	}
	return nil
}

type All struct {
	Output      Output `yaml:"output"`
	EmitExample bool   `yaml:"emit_example"`
	OmitExpect  bool   `yaml:"omit_expect"`
}

type Output struct {
	PkgName string `yaml:"pkg_name"`
	File    string `yaml:"file"`
}

type Package struct {
	Interface   string `yaml:"interface"`
	Struct      string `yaml:"struct"`
	Constructor string `yaml:"constructor"`
	OmitExpect  *bool  `yaml:"omit_expect"`
	EmitExample *bool  `yaml:"emit_example"`
	Output      Output `yaml:"output"`
}

func (p *Package) UnmarshalYAML(value *yml.Node) error {
	var name string
	if err := value.Decode(&name); err == nil {
		p.Interface = name
		return nil
	}

	type alias Package
	var aux alias
	if err := value.Decode(&aux); err != nil {
		return err
	}

	*p = Package(aux)
	return nil
}

func Parse(input []byte) (*Global, error) {
	cfg := Global{
		Version: "1",
		All: All{
			Output: Output{
				File: DefaultOutputFilename,
			},
		},
	}

	decoder := yml.NewDecoder(bytes.NewReader(input))
	decoder.KnownFields(true)

	err := decoder.Decode(&cfg)
	if err != nil {
		return nil, err
	}

	if err = cfg.validate(); err != nil {
		return nil, err
	}

	return &cfg, nil
}

func (g *Global) ToConfig() []mockgen.Config {
	var result []mockgen.Config
	for pkgPath, entries := range g.Packages {
		for _, entry := range entries {
			output := mergeOutput(g.All.Output, entry.Output)
			ee := g.All.EmitExample
			if entry.EmitExample != nil {
				ee = *entry.EmitExample
			}
			oe := g.All.OmitExpect
			if entry.OmitExpect != nil {
				oe = *entry.OmitExpect
			}

			result = append(result, mockgen.Config{
				PackagePath: pkgPath,
				Output: mockgen.Output{
					PackageName:  output.PkgName,
					TestFileName: output.File,
				},
				InterfaceName:   entry.Interface,
				StructName:      entry.Struct,
				ConstructorName: entry.Constructor,
				EmitExample:     ee,
				OmitExpect:      oe,
			})
		}
	}
	slices.SortFunc(result, func(a mockgen.Config, b mockgen.Config) int {
		return strings.Compare(a.PackagePath, b.PackagePath)
	})
	return result
}

func mergeOutput(outputs ...Output) Output {
	result := &Output{}
	for _, v := range outputs {
		if v.File != "" {
			result.File = v.File
		}
		if v.PkgName != "" {
			result.PkgName = v.PkgName
		}
	}
	return *result
}

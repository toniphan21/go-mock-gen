package mockgen

type Config struct {
	PackagePath     string
	Output          Output
	InterfaceName   string
	StructName      string // optional
	ConstructorName string // optional
	Namer           Namer  // optional
	EmitExample     bool
	OmitExpect      bool
	OmitGoDoc       bool
}

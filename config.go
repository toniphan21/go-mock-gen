package mockgen

type Config struct {
	PackagePath   string
	Output        Output
	InterfaceName string
	StructName    string // optional
	Namer         Namer  // optional
	SkipExpect    bool   // optional
	EmitExamples  bool
	OmitGoDoc     bool
}

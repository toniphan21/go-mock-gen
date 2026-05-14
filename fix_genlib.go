package mockgen

import (
	"go/types"

	"github.com/dave/jennifer/jen"
	genlib "nhatp.com/go/gen-lib"
)

// This is a temporary fix for the genlib function to handle known issues
// until the problem is resolved in nhatp.com/go/gen-lib.

// genlibTypeToJenCode is a temporary fix for genlib.TypeToJenCode.
func genlibTypeToJenCode(t types.Type) jen.Code {
	// An integration test for "type alias" covers this case;
	// there is no need for a unit test here.
	// The unit test should be implemented in nhatp.com/go/gen-lib.

	// Explicitly handle the buggy *types.Alias case; this block
	// will be removed once the upstream issue is fixed.
	if alias, ok := t.(*types.Alias); ok {
		obj := alias.Obj()
		pkg := obj.Pkg()

		if pkg != nil {
			return jen.Qual(pkg.Path(), obj.Name())
		}
		return jen.Id(obj.Name())
	}

	// Fall back to the upstream function for all other types.
	return genlib.TypeToJenCode(t)
}

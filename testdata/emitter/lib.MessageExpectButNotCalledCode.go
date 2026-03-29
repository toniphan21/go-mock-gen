package emitter

import (
	"fmt"
	"strings"
)

func repositoryMessageExpectButNotCalled(m repositoryMockMethod, want int, got int, index int) string {
	sb := &strings.Builder{}
	sb.WriteString(fmt.Sprintf("%s.%s was not called as expected\n", m.interfaceName(), m.methodName()))
	sb.WriteString(fmt.Sprintf("\twant: %d, got: %d\n\n", want, got))
	m.buildCallHistory(sb, "")
	sb.WriteString(fmt.Sprintf("\t#%d never called\n\n", index+1))
	sb.WriteString("\thint: add the missing call or remove the EXPECT above")
	return sb.String()
}

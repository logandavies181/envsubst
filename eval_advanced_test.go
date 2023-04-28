package envsubst

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEvalAdvanced(t *testing.T) {
	m := func(in string, n NodeInfo) (string, bool) {
		return os.Getenv(in), true
	}

	out, err := EvalAdvanced(`"${var:-5011}"`, m)
	assert.Nil(t, err)

	assert.Equal(t, `"5011"`, out)
}

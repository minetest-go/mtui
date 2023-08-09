package depanalyzer_test

import (
	"mtui/modmanager/depanalyzer"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSet1(t *testing.T) {
	a, err := depanalyzer.AnalyzeDeps("testdata/set1")
	assert.NoError(t, err)
	assert.NotNil(t, a)
	assert.Contains(t, a.Installed, "mymod")
	assert.Contains(t, a.Installed, "myothermod")
	assert.Contains(t, a.Installed, "somemod")
	assert.Equal(t, 3, len(a.Installed))

	assert.Contains(t, a.Missing, "missingmod")
	assert.Contains(t, a.Missing, "anothermissingmod")
	assert.Equal(t, 2, len(a.Missing))
}

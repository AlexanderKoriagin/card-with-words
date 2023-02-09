package random

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestRandom(t *testing.T) {
	testData := []struct {
		min, max int
	}{
		{0, 10000},
	}

	for _, td := range testData {
		for i := td.min; i < td.max; i++ {
			res := Random(td.min, td.max)
			assert.GreaterOrEqual(t, res, td.min)
			assert.Less(t, res, td.max)
		}
	}
}

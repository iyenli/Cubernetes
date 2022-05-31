package testing

import (
	"Cubernetes/pkg/utils/dag"
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestCheckCircle(t *testing.T) {
	nodes := make(map[string][]string)
	nodes["a"] = []string{"b", "c"}
	nodes["b"] = []string{"c", "d"}

	b, circle := dag.CheckCircle(nodes)
	fmt.Println(b, circle)
	assert.Equal(t, false, b)

	nodes["d"] = []string{"c", "a"}
	nodes["e"] = []string{"a", "b"}
	b, circle = dag.CheckCircle(nodes)
	fmt.Println(b, circle)
	assert.Equal(t, true, b)
}

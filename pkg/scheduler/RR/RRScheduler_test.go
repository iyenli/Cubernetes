package RR

import (
	"Cubernetes/pkg/scheduler/types"
	"github.com/stretchr/testify/assert"
	"strings"
	"testing"
)

func TestRRScheduler(t *testing.T) {
	rr := SchedulerRR{
		NumOfNodes:  0,
		NameOfNodes: nil,
		Next:        0,
	}

	_, err := rr.Schedule()
	assert.Error(t, err, ErrNoNodesToSchedule)

	err = rr.Init()
	assert.NoError(t, err)

	_, err = rr.Schedule()
	assert.Error(t, err, ErrNoNodesToSchedule)

	for i := 0; i < 10; i++ {
		info := types.NodeInfo{NodeUUID: strings.Repeat("s", i+1)}
		err = rr.AddNode(&info)
		assert.NoError(t, err)
	}

	for i := 0; i < 300; i++ {
		info, err := rr.Schedule()
		assert.NoError(t, err)
		assert.Equal(t, i%10+1, len(info.NodeUUID))
	}

	err = rr.RemoveNode(&types.NodeInfo{NodeUUID: "ss"})
	assert.NoError(t, err)
	err = rr.RemoveNode(&types.NodeInfo{NodeUUID: "ssss"})
	assert.NoError(t, err)

	_, err = rr.Schedule()
	assert.NoError(t, err)
}

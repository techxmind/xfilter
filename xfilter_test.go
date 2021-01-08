package xfilter

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/techxmind/filter/core"
)

func TestRun(t *testing.T) {
	filterItems := []interface{}{
		[]interface{}{"ctx.a", "=", "1"},
		[]interface{}{"ctx.b", "=", "1"},
		[]interface{}{"result", "=", map[string]interface{}{
			"foo": 1,
			"bar": 1,
		}},
	}

	ctx := core.WithContext(context.Background())
	ctx.Set("a", 1)
	ctx.Set("b", 1)

	f, _, err := Get(filterItems)
	assert.NoError(t, err)

	pluginExecuted := false
	RegisterResultObservers(NewResultObserver([]string{"result.foo"}, func(key string, value interface{}) bool {
		pluginExecuted = true
		return true
	}))

	data := make(map[string]interface{})
	err = Run(ctx, f, data)
	assert.NoError(t, err)
	assert.True(t, pluginExecuted)
	assert.EqualValues(t, map[string]interface{}{
		"result": map[string]interface{}{
			"bar": 1,
		},
	}, data)
}

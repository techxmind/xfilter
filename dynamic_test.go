package xfilter

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGet(t *testing.T) {
	ast := assert.New(t)

	filterItems := []interface{}{
		[]interface{}{"get.a", "=", "c"},
		[]interface{}{"result", "=", true},
	}

	filterJson, _ := json.Marshal(filterItems)

	filterGroupJson := []byte(`
[
	[
		["get.a", "=", "c"],
		["result", "=", true]
	],
	[
		["get.c", "=", "a"],
		["result", "=", true]
	]
]
	`)

	f1, old, err := Get(filterJson)
	ast.Nil(err)
	ast.NotNil(f1)
	ast.False(old)
	f2, old, err := Get(filterJson)
	ast.Nil(err)
	ast.True(old)
	ast.Equal(f1, f2)
	f3, old, err := Get(filterItems)
	ast.Nil(err)
	ast.True(old)
	ast.Equal(f1, f3)

	f4, old, err := Get(filterGroupJson)
	ast.Nil(err)
	ast.NotNil(f4)
	ast.False(old)
	f5, old, err := Get(filterGroupJson)
	ast.Nil(err)
	ast.True(old)
	ast.Equal(f4, f5)

	ast.NotEqual(f1, f4)
}

func BenchmarkGet(b *testing.B) {
	filterData := []byte(`
[
	[
		["get.a", "=", "c"],
		["get.a", "=", "c"],
		["get.a", "=", "c"],
		["get.a", "=", "c"],
		["get.a", "=", "c"],
		["get.a", "=", "c"],
		["get.a", "=", "c"],
		["get.a", "=", "c"],
		["get.a", "=", "c"],
		["get.a", "=", "c"],
		["result", "=", true]
	],
	[
		["get.a", "=", "c"],
		["get.a", "=", "c"],
		["get.a", "=", "c"],
		["get.a", "=", "c"],
		["get.a", "=", "c"],
		["get.a", "=", "c"],
		["get.a", "=", "c"],
		["get.a", "=", "c"],
		["get.a", "=", "c"],
		["get.a", "=", "c"],
		["get.a", "=", "c"],
		["get.a", "=", "c"],
		["get.a", "=", "c"],
		["get.a", "=", "c"],
		["get.a", "=", "c"],
		["get.a", "=", "c"],
		["get.a", "=", "c"],
		["get.a", "=", "c"],
		["get.a", "=", "c"],
		["get.a", "=", "c"],
		["get.a", "=", "c"],
		["get.a", "=", "c"],
		["get.c", "=", "a"],
		["result", "=", true]
	]
]
	`)

	for i := 0; i < b.N; i++ {
		Get(filterData)
	}
}

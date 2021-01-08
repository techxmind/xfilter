package xfilter

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/techxmind/config"
	"github.com/techxmind/filter/core"
)

func getTestMapConfig() *config.MapConfig {
	m, _ := config.JSONToMap([]byte(`
	{
		"l1" : {
			"l11" : "l11_value",
			"l12" : "l12_value"
		},
		"l2" : "l2_value",

		"filter" : [
			["ctx.foo", "=", "bar"],
			[
				["l1.l11", "=", "l11_value-mod"],
				["l2", "=", "l2_value-mod"]
			]
		]
	}
		`))

	return config.NewMapConfig(m)
}

func TestConfigWrapper(t *testing.T) {
	ast := assert.New(t)

	ctx := core.WithContext(context.Background())
	cfg := NewConfigWrapper(ctx, getTestMapConfig())

	ast.Equal("l11_value", cfg.String("l1.l11"))
	ast.Equal("l2_value", cfg.String("l2"))

	ctx.Set("foo", "bar")
	ast.Equal("l11_value-mod", cfg.String("l1.l11"))
	ast.Equal("l2_value-mod", cfg.String("l2"))
	ast.Equal(map[string]interface{}{
		"l11": "l11_value-mod",
		"l12": "l12_value",
	}, cfg.Get("l1"))
}

func BenchmarkConfigWrapper(b *testing.B) {
	ast := assert.New(b)

	for i := 0; i < b.N; i++ {
		ctx := core.WithContext(context.Background())
		cfg := NewConfigWrapper(ctx, getTestMapConfig())
		ctx.Set("foo", "bar")
		ast.Equal("l11_value-mod", cfg.String("l1.l11"))
	}
}

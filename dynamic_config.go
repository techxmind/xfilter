package xfilter

import (
	"context"

	"github.com/pkg/errors"
	"github.com/techxmind/config"
	"github.com/techxmind/filter"
	"github.com/techxmind/go-utils/object"
)

const (
	DynamicFilterKey = "filter"
)

// Execute filter when fetch config value, filter can overwrite original value.
//
//  config : {
//    "key1", "value1",
//    "key2" : {
//       "key21" : "value21",
//       "key22" : "value22"
//     },
//     "filter" : [
//       ["ctx.foo", "=", "bar"],
//       [
//         ["key1", "=", "value1-1"],
//         ["key2.key22", "=", "value22-1"]
//       ]
//     ]
//  }
//
//  w = NewConfigWrapper(ctx, cfg)
//  w.String("key1") // "value1"
//  w.JSON("key2") // {"key21":"value21", "key22":"value22"}
//  ifilter.ContextData(ctx).Set("foo", "bar")
//  w.String("key1") // "value1-1"
//  w.JSON("key2") // {"key21":"value21", "key22":"value22-1"}
//
type ConfigWrapper struct {
	config.ConfigHelper
}

type configWrapper struct {
	f   filter.Filter
	cfg config.Configer
	ctx context.Context

	config.ConfigHelper
}

func NewConfigWrapper(ctx context.Context, cfg config.Configer) *ConfigWrapper {
	w := &configWrapper{
		ctx: ctx,
		cfg: cfg,
	}

	filterData := cfg.Get(DynamicFilterKey)
	if filterData != nil {
		if f, _, err := Get(filterData); err != nil {
			logger.Errorf("config contains invalid dynamic filter:%v", err)
		} else {
			w.f = f
		}
	}

	return &ConfigWrapper{
		ConfigHelper: config.ConfigHelper{
			Configer: w,
		},
	}
}

func (w *configWrapper) Get(keyPath string) interface{} {

	if w.f == nil {
		return w.cfg.Get(keyPath)
	}

	result := make(map[string]interface{})
	err := Run(w.ctx, w.f, result)
	if err != nil {
		logger.Errorf("config filter wrapper execute filter err:%v", err)
		return w.cfg.Get(keyPath)
	}

	v, ok := object.GetValue(result, keyPath)
	if !ok {
		return w.cfg.Get(keyPath)
	}
	mv, ok := v.(map[string]interface{})
	if !ok {
		return v
	}

	ov := w.cfg.Get(keyPath)
	if ov == nil {
		return v
	}
	mov, ok := ov.(map[string]interface{})
	if !ok {
		return v
	}

	extendMap(mv, mov)

	return mv
}

func (w *configWrapper) Set(keyPath string, value interface{}) error {
	return errors.New("set method is unsupported by config filter wrapper")
}

func extendMap(originMap map[string]interface{}, extraMap map[string]interface{}) {
	for k, v := range extraMap {
		originV, originExist := originMap[k]
		if !originExist {
			originMap[k] = v
			continue
		}
		originSubMap, originSubMapOk := originV.(map[string]interface{})
		if !originSubMapOk {
			continue
		}
		if subMap, ok := v.(map[string]interface{}); ok {
			extendMap(originSubMap, subMap)
		}
	}
}

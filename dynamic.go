package xfilter

import (
	"encoding/json"

	"github.com/hashicorp/golang-lru"
	"github.com/spaolacci/murmur3"
	"github.com/techxmind/filter"
)

const (
	DefaultCacheSize = 1000
)

var (
	_cache Cache = NewCache(DefaultCacheSize)
)

type Cache interface {
	Get(key interface{}) (value interface{}, ok bool)
	Add(key, value interface{})
}

func SetCache(c Cache) {
	_cache = c
}

func NewCache(size int) Cache {
	c, err := lru.NewARC(size)

	if err != nil {
		logger.Errorf("Init cache err:%v", err)
		return nil
	}

	return c
}

// FromJSON get Filter by json data.
//
// It returns cached Filter after first call if the json data not changed.
// "filterID" is used for clearing old version data.
//
func Get(filterData interface{}, options ...filter.Option) (f filter.Filter, cache bool, err error) {
	var (
		hashID      uint64
		filterItems []interface{}
	)

	if v, ok := filterData.([]byte); ok {
		hashID = getHashID(v)
		if err = json.Unmarshal(v, &filterItems); err != nil {
			return
		}
	} else if v, ok := filterData.(string); ok {
		hashID = getHashID([]byte(v))
		if err = json.Unmarshal([]byte(v), &filterItems); err != nil {
			return
		}
	} else if v, ok := filterData.([]interface{}); ok {
		js, e := json.Marshal(v)
		if e != nil {
			err = e
			return
		}
		hashID = getHashID(js)
		filterItems = v
	} else {
		err = ErrInvalidFilterData
		return
	}

	if _cache != nil {
		if v, ok := _cache.Get(hashID); ok {
			return v.(filter.Filter), true, nil
		}
	}

	f, err = filter.New(filterItems, options...)
	if err != nil {
		return
	}

	logger.Debugf("new filter[%s]", f.Name())

	if _cache != nil {
		_cache.Add(hashID, f)
	}

	return
}

func getHashID(s []byte) uint64 {
	return murmur3.Sum64(s)
}

package xfilter

import (
	"strings"

	"github.com/techxmind/go-utils/object"
)

var (
	_resultObserverManager = &resultObserverManager{
		observers: make(map[string][]ResultObserver, 0),
		keys:      make([]string, 0),
	}
)

func RegisterResultObservers(observers ...ResultObserver) {
	_resultObserverManager.register(observers...)
}

// Observer if specified keys appears in filter result, then trigger the Callback
type ResultObserver interface {
	ObserveKeys() []string
	Callback(key string, value interface{}) (removeIt bool)
}

type resultObserverCreator struct {
	keys []string
	cb   func(key string, value interface{}) bool
}

func NewResultObserver(observeKeys []string, cb func(key string, value interface{}) bool) ResultObserver {
	return &resultObserverCreator{
		keys: observeKeys,
		cb:   cb,
	}
}

func (o *resultObserverCreator) ObserveKeys() []string {
	return o.keys
}

func (o *resultObserverCreator) Callback(key string, value interface{}) bool {
	return o.cb(key, value)
}

type resultObserverManager struct {
	observers map[string][]ResultObserver
	keys      []string
}

func (m *resultObserverManager) register(observers ...ResultObserver) {
	for _, observer := range observers {
		keys := observer.ObserveKeys()
		for _, key := range keys {
			if _, ok := m.observers[key]; !ok {
				m.observers[key] = make([]ResultObserver, 0, 1)
				m.keys = append(m.keys, key)
			}
			m.observers[key] = append(m.observers[key], observer)
		}
	}
}

func (m *resultObserverManager) observe(result interface{}) {
	for _, key := range m.keys {
		if v, ok := object.GetValue(result, key); ok && v != nil {
			removeIt := false
			for _, observer := range m.observers[key] {
				if observer.Callback(key, v) {
					removeIt = true
				}
			}
			if removeIt {
				segments := strings.Split(key, ".")
				if len(segments) > 1 {
					last := segments[len(segments)-1]
					parentKey := strings.Join(segments[:len(segments)-1], ".")
					if parent, ok := object.GetObject(result, parentKey, false); ok {
						if m, ok := parent.(map[string]interface{}); ok {
							delete(m, last)
						}
					}
				}
			}
		}
	}
}

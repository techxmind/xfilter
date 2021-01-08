package xfilter

func init() {
	RegisterResultObservers(NewResultObserver([]string{"metric"}, func(key string, value interface{}) bool {
		if m, ok := value.(map[string]interface{}); ok {
			for action, _ := range m {
				logger.Debugf("metric:%s", action)
				//imetrics.SimpleCounter.Add(1, action)
			}
		}
		return true
	}))
}

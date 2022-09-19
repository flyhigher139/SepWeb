package route

type matchInfo struct {
	N          *node
	PathParams map[string]string
}

func (m *matchInfo) addValue(key string, value string) {
	if m.PathParams == nil {
		m.PathParams = make(map[string]string, 1)
	}
	m.PathParams[key] = value
}

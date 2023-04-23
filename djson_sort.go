package djson

func (m *JSON) SortElement(isAsc bool, k ...interface{}) bool {
	var tArray *DA

	if len(k) == 0 {
		if m._Type == ARRAY {
			tArray = m._Array
		}
	}

	if len(k) > 0 {

		if m._Type == OBJECT {
			if key, ok := k[0].(string); ok {
				if da, ok := m._Object.Array(key); ok {
					tArray = da
				}
			}
		} else if m._Type == ARRAY {
			if idx, ok := k[0].(int); ok {
				if da, ok := m._Array.Array(idx); ok {
					tArray = da
				}
			}
		}
	}

	if tArray != nil {
		return tArray.Sort(isAsc)
	} else {
		return false
	}
}

func (m *JSON) SortAsc(k ...interface{}) bool {
	return m.SortElement(true, k...)
}

func (m *JSON) SortDesc(k ...interface{}) bool {
	return m.SortElement(false, k...)
}

func (m *JSON) SortArray(isAsc bool, k ...string) bool {
	if m._Type != ARRAY {
		return false
	}

	return m._Array.Sort(isAsc, k...)
}

func (m *JSON) SortArrayAsc(k ...string) bool {
	return m.SortArray(true, k...)
}

func (m *JSON) SortArrayDesc(k ...string) bool {
	return m.SortArray(false, k...)
}

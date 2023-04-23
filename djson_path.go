package djson

func (m *JSON) ObjectPath(path string) (*JSON, bool) {
	retJson := New()

	pok := m.DoPathFunc(path, nil,
		func(da *DA, idx int, v interface{}) {
			if obj, ok := da.Object(idx); ok {
				retJson._Object = obj
				retJson._Type = OBJECT
			}
		},
		func(do *DO, key string, v interface{}) {
			if obj, ok := do.Object(key); ok {
				retJson._Object = obj
				retJson._Type = OBJECT
			}
		},
	)

	if !pok || retJson._Type != OBJECT {
		return nil, false
	}

	return retJson, true
}

func (m *JSON) ArrayPath(path string) (*JSON, bool) {
	retJson := New()

	pok := m.DoPathFunc(path, nil,
		func(da *DA, idx int, v interface{}) {
			if arr, ok := da.Array(idx); ok {
				retJson._Array = arr
				retJson._Type = ARRAY
			}
		},
		func(do *DO, key string, v interface{}) {
			if arr, ok := do.Array(key); ok {
				retJson._Array = arr
				retJson._Type = ARRAY
			}
		},
	)

	if !pok || retJson._Type != ARRAY {
		return nil, false
	}

	return retJson, true
}

func (m *JSON) FloatPath(path string, dv ...float64) float64 {
	var ret float64
	var kok bool

	pok := m.DoPathFunc(path, nil,
		func(da *DA, idx int, v interface{}) {
			ret, kok = da.Float(idx)
		},
		func(do *DO, key string, v interface{}) {
			ret, kok = do.Float(key)
		},
	)

	if pok && kok {
		return ret
	}

	if len(dv) > 0 {
		return dv[0]
	}

	return 0
}

func (m *JSON) IntPath(path string, dv ...int64) int64 {
	var ret int64
	var kok bool

	pok := m.DoPathFunc(path, nil,
		func(da *DA, idx int, v interface{}) {
			ret, kok = da.Int(idx)
		},
		func(do *DO, key string, v interface{}) {
			ret, kok = do.Int(key)
		},
	)

	if pok && kok {
		return ret
	}

	if len(dv) > 0 {
		return dv[0]
	}

	return 0
}

func (m *JSON) BoolPath(path string, dv ...bool) bool {
	var ret bool
	var kok bool

	pok := m.DoPathFunc(path, nil,
		func(da *DA, idx int, v interface{}) {
			ret, kok = da.Bool(idx)
		},
		func(do *DO, key string, v interface{}) {
			ret, kok = do.Bool(key)
		},
	)

	if pok && kok {
		return ret
	}

	if len(dv) > 0 {
		return dv[0]
	}

	return false
}

func (m *JSON) StringPath(path string) string {
	var ret string

	_ = m.DoPathFunc(path, nil,
		func(da *DA, idx int, v interface{}) {
			ret = da.String(idx)
		},
		func(do *DO, key string, v interface{}) {
			ret = do.String(key)
		},
	)

	return ret
}

func (m *JSON) TypePath(path string) string {
	var pathType string

	_ = m.DoPathFunc(path, nil,
		func(da *DA, idx int, v interface{}) {
			pathType, _ = da.Type(idx)
		},
		func(do *DO, key string, v interface{}) {
			pathType, _ = do.Type(key)
		},
	)

	return pathType
}

func (m *JSON) SortAscPath(path string, k ...string) bool {
	return m.SortPath(path, true, k...)
}

func (m *JSON) SortDescPath(path string, k ...string) bool {
	return m.SortPath(path, false, k...)
}

func (m *JSON) SortPath(path string, isAsc bool, k ...string) bool {
	var isSorted bool

	pok := m.DoPathFunc(path, nil,
		func(da *DA, idx int, v interface{}) {
			if tda, ok := da.Array(idx); ok {
				isSorted = tda.Sort(isAsc, k...)
			} else {
				isSorted = false
			}
		},
		func(do *DO, key string, v interface{}) {
			if tda, ok := do.Array(key); ok {
				isSorted = tda.Sort(isAsc, k...)
			} else {
				isSorted = false
			}
		},
	)

	if !pok || !isSorted {
		return false
	} else {
		return true
	}
}

func (m *JSON) RemovePath(path string) bool {
	return m.DoPathFunc(path, nil,
		func(da *DA, idx int, v interface{}) {
			da.Remove(idx)
		},
		func(do *DO, key string, v interface{}) {
			do.Remove(key)
		},
	)
}

func (m *JSON) PutObjectToPath(path string, okey string, oval interface{}) bool {
	return m.DoPathFunc(path, oval,
		func(da *DA, idx int, v interface{}) {
			da.Insert(idx, Object{okey: v})
		},
		func(do *DO, key string, v interface{}) {
			do.Put(key, Object{okey: v})
		},
	)
}

// Replace or insert values as array

func (m *JSON) PutArrayToPath(path string, val ...interface{}) bool {
	return m.DoPathFunc(path, val,
		func(da *DA, idx int, v interface{}) {
			da.Insert(idx, v)
		},
		func(do *DO, key string, v interface{}) {
			do.Put(key, v)
		},
	)
}

// Pushback a value to array if possible.
// The path must indicate array.

func (m *JSON) PushBackToPath(path string, val interface{}) bool {
	return m.DoPathFunc(path, val,
		func(da *DA, idx int, v interface{}) {
			if dda, ok := da.Array(idx); ok {
				dda.PushBack(v)
			}
		},
		func(do *DO, key string, v interface{}) {
			if dda, ok := do.Array(key); ok {
				dda.PushBack(v)
			}
		},
	)
}

// Replace or insert a value

func (m *JSON) UpdatePath(path string, val interface{}) bool {
	return m.DoPathFunc(path, val,
		func(da *DA, idx int, v interface{}) {
			da.ReplaceAt(idx, v)
		},
		func(do *DO, key string, v interface{}) {
			do.Put(key, v)
		},
	)
}

func (m *JSON) doPathFuncCore(
	arrayTaskFunc func(da *DA, idx int, v interface{}),
	objectTaskFunc func(do *DO, key string, v interface{}),
	val interface{}, token ...interface{}) bool {

	jsonMode := m._Type
	dObject := m._Object
	dArray := m._Array

	tokenLen := len(token)

	for idx := range token {
		switch tkey := token[idx].(type) {
		case string:
			if jsonMode != OBJECT || dObject == nil {
				return false
			}

			if idx == tokenLen-1 {
				objectTaskFunc(dObject, tkey, val)
				return true
			} else {
				if _, ok := dObject.Map[tkey]; !ok {
					return false
				}

				switch t := dObject.Map[tkey].(type) {
				case *DO:
					dObject = t
					dArray = nil
					jsonMode = OBJECT
				case *DA:
					dObject = nil
					dArray = t
					jsonMode = ARRAY
				default:
					return false
				}
			}
		case int:
			if jsonMode != ARRAY || dArray == nil {
				return false
			}

			for dArray.Size() < tkey {
				dArray.PushBack(0)
			}

			if idx == tokenLen-1 {
				arrayTaskFunc(dArray, tkey, val)
				return true
			} else {
				switch t := dArray.Element[tkey].(type) {
				case *DO:
					dObject = t
					dArray = nil
					jsonMode = OBJECT
				case *DA:
					dObject = nil
					dArray = t
					jsonMode = ARRAY
				default:
					return false
				}
			}
		default:
			return false
		}
	}

	return false

}

func (m *JSON) DoPathFunc(path string, val interface{},
	arrayTaskFunc func(da *DA, idx int, v interface{}),
	objectTaskFunc func(do *DO, key string, v interface{})) bool {
	return m.doPathFuncCore(arrayTaskFunc, objectTaskFunc, val, PathTokenizer(path)...)
}

func (m *JSON) KeysPath(path string) ([]string, bool) {
	rk := make([]string, 0)

	pok := m.DoPathFunc(path, nil,
		func(da *DA, idx int, v interface{}) {
			if ddo, ok := da.Object(idx); ok {
				for k := range ddo.Map {
					rk = append(rk, k)
				}
			}
		},
		func(do *DO, key string, v interface{}) {
			if ddo, ok := do.Object(key); ok {
				for k := range ddo.Map {
					rk = append(rk, k)
				}
			}
		},
	)

	if !pok {
		return []string{}, false
	}

	return rk, true
}

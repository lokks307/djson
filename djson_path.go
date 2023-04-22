package djson

func (m *JSON) ObjectPath(path string) (*JSON, bool) {

	retJson := New()

	err := m.DoPathFunc(path, nil,
		func(da *DA, idx int, v interface{}) {
			if obj, ok := da.GetAsObject(idx); ok {
				retJson._Object = obj
				retJson._Type = OBJECT
			}
		},
		func(do *DO, key string, v interface{}) {
			if obj, ok := do.GetAsObject(key); ok {
				retJson._Object = obj
				retJson._Type = OBJECT
			}
		},
	)

	if err != nil || retJson._Type != OBJECT {
		return nil, false
	}

	return retJson, true
}

func (m *JSON) ArrayPath(path string) (*JSON, bool) {

	retJson := New()

	err := m.DoPathFunc(path, nil,
		func(da *DA, idx int, v interface{}) {
			if arr, ok := da.GetAsArray(idx); ok {
				retJson._Array = arr
				retJson._Type = ARRAY
			}
		},
		func(do *DO, key string, v interface{}) {
			if arr, ok := do.GetAsArray(key); ok {
				retJson._Array = arr
				retJson._Type = ARRAY
			}
		},
	)

	if err != nil || retJson._Type != ARRAY {
		return nil, false
	}

	return retJson, true
}

func (m *JSON) FloatPath(path string, defFloat ...float64) float64 {
	var retFloat float64
	var ok bool

	err := m.DoPathFunc(path, nil,
		func(da *DA, idx int, v interface{}) {
			retFloat, ok = da.GetAsFloat(idx)
		},
		func(do *DO, key string, v interface{}) {
			retFloat, ok = do.GetAsFloat(key)
		},
	)

	if err == nil && ok {
		return retFloat
	}

	if len(defFloat) > 0 {
		return defFloat[0]
	}

	return 0
}

func (m *JSON) IntPath(path string, defInt ...int64) int64 {
	var retInt int64
	var ok bool

	err := m.DoPathFunc(path, nil,
		func(da *DA, idx int, v interface{}) {
			retInt, ok = da.GetAsInt(idx)
		},
		func(do *DO, key string, v interface{}) {
			retInt, ok = do.GetAsInt(key)
		},
	)

	if err == nil && ok {
		return retInt
	}

	if len(defInt) > 0 {
		return defInt[0]
	}

	return 0
}

func (m *JSON) BoolPath(path string, defBool ...bool) bool {

	var retBool bool
	var ok bool

	err := m.DoPathFunc(path, nil,
		func(da *DA, idx int, v interface{}) {
			retBool, ok = da.GetAsBool(idx)
		},
		func(do *DO, key string, v interface{}) {
			retBool, ok = do.GetAsBool(key)
		},
	)

	if err == nil && ok {
		return retBool
	}

	if len(defBool) > 0 {
		return defBool[0]
	}

	return false
}

func (m *JSON) StringPath(path string) string {
	var retStr string

	_ = m.DoPathFunc(path, nil,
		func(da *DA, idx int, v interface{}) {
			retStr = da.GetAsString(idx)
		},
		func(do *DO, key string, v interface{}) {
			retStr = do.String(key)
		},
	)

	return retStr
}

func (m *JSON) TypePath(path string) string {
	var pathType string

	_ = m.DoPathFunc(path, nil,
		func(da *DA, idx int, v interface{}) {
			pathType, _ = da.GetType(idx)
		},
		func(do *DO, key string, v interface{}) {
			pathType, _ = do.GetType(key)
		},
	)

	return pathType
}

func (m *JSON) SortObjectArrayPath(path string, isAsc bool, okey string) error {
	var isSorted bool

	err := m.DoPathFunc(path, nil,
		func(da *DA, idx int, v interface{}) {
			if tda, ok := da.GetAsArray(idx); ok {
				isSorted = tda.SortObject(isAsc, okey)
			} else {
				isSorted = false
			}
		},
		func(do *DO, key string, v interface{}) {
			if tda, ok := do.GetAsArray(key); ok {
				isSorted = tda.SortObject(isAsc, okey)
			} else {
				isSorted = false
			}
		},
	)

	if err != nil || !isSorted {
		return failedToSortError
	} else {
		return nil
	}
}

func (m *JSON) SortObjectArrayAscPath(path string, key string) error {
	return m.SortObjectArrayPath(path, true, key)
}

func (m *JSON) SortObjectArrayDescPath(path string, key string) error {
	return m.SortObjectArrayPath(path, false, key)
}

func (m *JSON) SortPath(path string, isAsc bool) error {
	var isSorted bool

	err := m.DoPathFunc(path, nil,
		func(da *DA, idx int, v interface{}) {
			if tda, ok := da.GetAsArray(idx); ok {
				isSorted = tda.Sort(isAsc)
			} else {
				isSorted = false
			}
		},
		func(do *DO, key string, v interface{}) {
			if tda, ok := do.GetAsArray(key); ok {
				isSorted = tda.Sort(isAsc)
			} else {
				isSorted = false
			}
		},
	)

	if err != nil || !isSorted {
		return failedToSortError
	} else {
		return nil
	}
}

func (m *JSON) SortDescPath(path string) error {
	return m.SortPath(path, false)
}

func (m *JSON) SortAscPath(path string) error {
	return m.SortPath(path, true)
}

func (m *JSON) RemovePath(path string) error {
	return m.DoPathFunc(path, nil,
		func(da *DA, idx int, v interface{}) {
			da.Remove(idx)
		},
		func(do *DO, key string, v interface{}) {
			do.Remove(key)
		},
	)
}

func (m *JSON) PutPathNewObject(path string, okey string, oval interface{}) error {
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

func (m *JSON) PutPathNewArray(path string, val ...interface{}) error {
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

func (m *JSON) PushBackPath(path string, val interface{}) error {
	return m.DoPathFunc(path, val,
		func(da *DA, idx int, v interface{}) {
			if dda, ok := da.GetAsArray(idx); ok {
				dda.PushBack(v)
			}
		},
		func(do *DO, key string, v interface{}) {
			if dda, ok := do.GetAsArray(key); ok {
				dda.PushBack(v)
			}
		},
	)
}

// Replace or insert a value

func (m *JSON) UpdatePath(path string, val interface{}) error {
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
	val interface{}, token ...interface{}) error {

	jsonMode := m._Type
	dObject := m._Object
	dArray := m._Array

	tokenLen := len(token)

	for idx := range token {
		switch tkey := token[idx].(type) {
		case string:
			if jsonMode != OBJECT || dObject == nil {
				return invalidPathError
			}

			if idx == tokenLen-1 {
				objectTaskFunc(dObject, tkey, val)
				return nil
			} else {
				if _, ok := dObject.Map[tkey]; !ok {
					return invalidPathError
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
					return invalidPathError
				}
			}
		case int:
			if jsonMode != ARRAY || dArray == nil {
				return invalidPathError
			}

			for dArray.Size() < tkey {
				dArray.PushBack(0)
			}

			if idx == tokenLen-1 {
				arrayTaskFunc(dArray, tkey, val)
				return nil
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
					return invalidPathError
				}
			}
		default:
			return invalidPathError
		}
	}

	return invalidPathError

}

func (m *JSON) DoPathFunc(path string, val interface{},
	arrayTaskFunc func(da *DA, idx int, v interface{}),
	objectTaskFunc func(do *DO, key string, v interface{})) error {
	return m.doPathFuncCore(arrayTaskFunc, objectTaskFunc, val, PathTokenizer(path)...)
}

func (m *JSON) KeysPath(path string) ([]string, error) {
	rk := make([]string, 0)

	err := m.DoPathFunc(path, nil,
		func(da *DA, idx int, v interface{}) {
			if ddo, ok := da.GetAsObject(idx); ok {
				for k := range ddo.Map {
					rk = append(rk, k)
				}
			}
		},
		func(do *DO, key string, v interface{}) {
			if ddo, ok := do.GetAsObject(key); ok {
				for k := range ddo.Map {
					rk = append(rk, k)
				}
			}
		},
	)

	if err != nil {
		return []string{}, err
	}

	return rk, nil
}

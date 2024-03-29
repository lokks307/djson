package djson

import (
	"math"
	"reflect"

	"github.com/goccy/go-json"
	"github.com/volatiletech/null/v8"
)

type DO struct {
	Map map[string]interface{}
}

func NewDO() *DO {
	return &DO{
		Map: make(map[string]interface{}),
	}
}

func (m *DO) Put(key string, value interface{}) *DO {

	if IsFloatType(value) {
		switch t := value.(type) {
		case float32:
			if !math.IsNaN(float64(t)) && !math.IsInf(float64(t), 0) {
				m.Map[key] = t
			}
		case float64:
			if !math.IsNaN(t) && !math.IsInf(float64(t), 0) {
				m.Map[key] = t
			}
		}

		return m
	}

	if IsBaseType(value) {
		m.Map[key] = value
		return m
	}

	if n, ok := value.(json.Number); ok {
		if i, err := n.Int64(); err == nil {
			m.Map[key] = i
			return m
		}
		if f, err := n.Float64(); err == nil {
			// log.Println(math.IsNaN(f))
			m.Map[key] = f
			return m
		}
	}

	switch t := value.(type) {
	case null.String:
		if t.Valid {
			m.Map[key] = t.String
		} else {
			m.Map[key] = ""
		}
	case null.Bool:
		if t.Valid {
			m.Map[key] = t.Bool
		} else {
			m.Map[key] = false
		}
	case null.Int:
		if t.Valid {
			m.Map[key] = t.Int
		} else {
			m.Map[key] = 0
		}
	case null.Int8:
		if t.Valid {
			m.Map[key] = t.Int8
		} else {
			m.Map[key] = 0
		}
	case null.Int16:
		if t.Valid {
			m.Map[key] = t.Int16
		} else {
			m.Map[key] = 0
		}
	case null.Int32:
		if t.Valid {
			m.Map[key] = t.Int32
		} else {
			m.Map[key] = 0
		}
	case null.Int64:
		if t.Valid {
			m.Map[key] = t.Int64
		} else {
			m.Map[key] = 0
		}
	case null.Uint:
		if t.Valid {
			m.Map[key] = t.Uint
		} else {
			m.Map[key] = 0
		}
	case null.Uint8:
		if t.Valid {
			m.Map[key] = t.Uint8
		} else {
			m.Map[key] = 0
		}
	case null.Uint16:
		if t.Valid {
			m.Map[key] = t.Uint16
		} else {
			m.Map[key] = 0
		}
	case null.Uint32:
		if t.Valid {
			m.Map[key] = t.Uint32
		} else {
			m.Map[key] = 0
		}
	case null.Uint64:
		if t.Valid {
			m.Map[key] = t.Uint64
		} else {
			m.Map[key] = 0
		}
	case null.Float32:
		if t.Valid {
			m.Map[key] = t.Float32
		} else {
			m.Map[key] = float32(0.0)
		}
	case null.Float64:
		if t.Valid {
			m.Map[key] = t.Float64
		} else {
			m.Map[key] = float64(0.0)
		}
	case DO:
		m.Map[key] = &t
	case DA:
		m.Map[key] = &t
	case *DO:
		m.Map[key] = t
	case *DA:
		m.Map[key] = t
	case map[string]interface{}:
		m.Map[key] = MapToObject(t)
	case []interface{}:
		m.Map[key] = SliceToArray(t)
	case Object:
		m.Map[key] = MapToObject(t)
	case Array:
		m.Map[key] = SliceToArray(t)
	case JSON:
		m.Map[key] = t.Interface()
	case *JSON:
		m.Map[key] = t.Interface()
	case []string:
		m.Map[key] = PremitiveSliceToArray(t)
	case []bool:
		m.Map[key] = PremitiveSliceToArray(t)
	case []float32:
		m.Map[key] = PremitiveSliceToArray(t)
	case []float64:
		m.Map[key] = PremitiveSliceToArray(t)
	case []int:
		m.Map[key] = PremitiveSliceToArray(t)
	case []int8:
		m.Map[key] = PremitiveSliceToArray(t)
	case []int16:
		m.Map[key] = PremitiveSliceToArray(t)
	case []int32:
		m.Map[key] = PremitiveSliceToArray(t)
	case []int64:
		m.Map[key] = PremitiveSliceToArray(t)
	case []uint:
		m.Map[key] = PremitiveSliceToArray(t)
	case []uint8:
		m.Map[key] = PremitiveSliceToArray(t)
	case []uint16:
		m.Map[key] = PremitiveSliceToArray(t)
	case []uint32:
		m.Map[key] = PremitiveSliceToArray(t)
	case []uint64:
		m.Map[key] = PremitiveSliceToArray(t)
	case []null.String:
		m.Map[key] = PremitiveSliceToArray(t)
	case []null.Bool:
		m.Map[key] = PremitiveSliceToArray(t)
	case []null.Float32:
		m.Map[key] = PremitiveSliceToArray(t)
	case []null.Float64:
		m.Map[key] = PremitiveSliceToArray(t)
	case []null.Int:
		m.Map[key] = PremitiveSliceToArray(t)
	case []null.Int8:
		m.Map[key] = PremitiveSliceToArray(t)
	case []null.Int16:
		m.Map[key] = PremitiveSliceToArray(t)
	case []null.Int32:
		m.Map[key] = PremitiveSliceToArray(t)
	case []null.Int64:
		m.Map[key] = PremitiveSliceToArray(t)
	case []null.Uint:
		m.Map[key] = PremitiveSliceToArray(t)
	case []null.Uint8:
		m.Map[key] = PremitiveSliceToArray(t)
	case []null.Uint16:
		m.Map[key] = PremitiveSliceToArray(t)
	case []null.Uint32:
		m.Map[key] = PremitiveSliceToArray(t)
	case []null.Uint64:
		m.Map[key] = PremitiveSliceToArray(t)
	case nil:
		m.Map[key] = nil
	}

	return m
}

func (m *DO) PutArray(key string, array ...interface{}) *DO {
	nArray := NewDA()
	nArray.Put(array)
	m.Put(key, nArray)
	return m
}

func (m *DO) Append(obj map[string]interface{}) *DO {
	for k, v := range obj {
		m.Put(k, v)
	}

	return m
}

func (m *DO) HasKey(key string) bool {
	_, ok := m.Map[key]
	return ok
}

func (m *DO) String(key string) string {
	if key == "" {
		return ""
	}

	value, ok := m.Map[key]
	if !ok {
		return ""
	}

	switch t := value.(type) {
	case DO:
		return t.ToString()
	case DA:
		return t.ToString()
	case *DO:
		return t.ToString()
	case *DA:
		return t.ToString()
	case nil:
		return "null"
	}

	if str, ok := getStringBase(m.Map[key]); ok {
		return str
	}

	return ""
}

func (m *DO) String2(key string) (string, bool) {
	value, ok := m.Map[key]
	if !ok {
		return "", false
	}

	switch t := value.(type) {
	case DO:
		return t.ToString(), true
	case DA:
		return t.ToString(), true
	case *DO:
		return t.ToString(), true
	case *DA:
		return t.ToString(), true
	case nil:
		return "null", true
	}

	return getStringBase(m.Map[key])
}

func (m *DO) Get(key string) (interface{}, bool) {
	value, ok := m.Map[key]
	if !ok {
		return nil, false
	}

	return value, true
}

func (m *DO) Type(key string) (string, bool) {
	value, ok := m.Map[key]
	if !ok {
		return "", false
	}

	switch value.(type) {
	case DA, *DA:
		return "array", true
	case DO, *DO:
		return "object", true
	case int, uint, int8, uint8, int16, uint16, int32, uint32, int64, uint64:
		return "int", true
	case float32, float64:
		return "float", true
	case string:
		return "string", true
	case bool:
		return "bool", true
	case nil:
		return "null", true
	}

	return "", false
}

func (m *DO) Bool(key string) (bool, bool) {
	value, ok := m.Map[key]
	if !ok {
		return false, false
	}

	if boolVal, ok := getBoolBase(value); ok {
		return boolVal, true
	}

	return false, false
}

func (m *DO) Float(key string) (float64, bool) {
	value, ok := m.Map[key]
	if !ok {
		return 0, false
	}

	if floatVal, ok := getFloatBase(value); ok {
		return floatVal, true
	}

	return 0, false
}

func (m *DO) Int(key string) (int64, bool) {
	value, ok := m.Map[key]
	if !ok {
		return 0, false
	}

	if intVal, ok := getIntBase(value); ok {
		return intVal, true
	}

	return 0, false
}

func (m *DO) Object(key string) (*DO, bool) {
	value, ok := m.Map[key]
	if !ok {
		return nil, false
	}

	switch t := value.(type) {
	case DO:
		return &t, true
	case *DO:
		return t, true
	case **DO:
		return *t, true
	}

	return nil, false
}

func (m *DO) Array(key string) (*DA, bool) {

	value, ok := m.Map[key]
	if !ok {
		return nil, false
	}

	switch t := value.(type) {
	case DA:
		return &t, true
	case *DA:
		return t, true
	case **DA:
		return *t, true
	}

	return nil, false

}

func (m *DO) Remove(keys ...string) *DO {
	for idx := range keys {
		delete(m.Map, keys[idx])
	}
	return m
}

func (m *DO) ToStringPretty() string {
	jsonByte, _ := json.MarshalIndent(ObjectToMap(m), "", "   ")
	return string(jsonByte)
}

func (m *DO) ToString() string {
	jsonByte, err := json.Marshal(ObjectToMap(m))
	if err != nil {
		// log.Println(err)
		return ""
	}
	return string(jsonByte)
}

func (m *DO) Len() int {
	return len(m.Map)
}

func (m *DO) Size() int {
	return len(m.Map)
}

func (m *DO) Equal(t *DO) bool {
	if m.Size() != t.Size() {
		return false
	}

	for i := range m.Map {

		if m.Map[i] == nil || t.Map[i] == nil {
			if m.Map[i] == nil && t.Map[i] == nil {
				continue
			}
			return false
		}

		mtype := reflect.TypeOf(m.Map[i]).String()
		ttype := reflect.TypeOf(t.Map[i]).String()

		if mtype != ttype {
			return false
		}

		switch m.Map[i].(type) {
		case string:
			if m.Map[i].(string) != t.Map[i].(string) {
				return false
			}
		case bool:
			if m.Map[i].(bool) != t.Map[i].(bool) {
				return false
			}
		case int, uint, int8, uint8, int16, uint16, int32, uint32, int64, uint64:
			mInt, _ := m.Int(i)
			tInt, _ := t.Int(i)
			if mInt != tInt {
				return false
			}
		case float32, float64:
			mFloat, _ := m.Float(i)
			tFloat, _ := t.Float(i)
			if mFloat != tFloat {
				return false
			}
		case *DO:
			mdo := m.Map[i].(*DO)
			tdo := t.Map[i].(*DO)

			if !mdo.Equal(tdo) {
				return false
			}
		case *DA:
			mda := m.Map[i].(*DA)
			tda := t.Map[i].(*DA)

			if !mda.Equal(tda) {
				return false
			}
		case *JSON:
			mjson := m.Map[i].(*JSON)
			tjson := t.Map[i].(*JSON)

			if !mjson.Equal(tjson) {
				return false
			}
		}
	}

	return true
}

func (m *DO) Clone() *DO {

	t := NewDO()

	t.Map = make(map[string]interface{})

	for k := range m.Map {

		if m.Map[k] == nil {
			t.Map[k] = nil
			continue
		}

		switch m.Map[k].(type) {
		case string:
			t.Map[k] = m.Map[k].(string)
		case bool:
			t.Map[k] = m.Map[k].(bool)
		case int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64:
			t.Map[k], _ = m.Int(k)
		case float64:
			t.Map[k], _ = m.Float(k)
		case *DO:
			mdo := m.Map[k].(*DO)
			t.Map[k] = mdo.Clone()
		case *DA:
			mda := m.Map[k].(*DA)
			t.Map[k] = mda.Clone()
		case *JSON:
			mdjson := m.Map[k].(*JSON)
			t.Map[k] = mdjson.Clone()
		}
	}

	return t
}

func (m *DO) Rename(from, to string) bool {
	if !m.HasKey(from) || from == to {
		return false
	}

	m.Map[to] = m.Map[from]
	delete(m.Map, from)

	return true
}

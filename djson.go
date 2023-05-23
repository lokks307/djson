package djson

import (
	"reflect"
	"strconv"
	"strings"

	gov "github.com/asaskevich/govalidator"
)

const (
	NULL   = 0
	OBJECT = 1
	ARRAY  = 2
	STRING = 3
	INT    = 4
	FLOAT  = 5
	BOOL   = 6
)

type JSON struct {
	_Object *DO
	_Array  *DA
	_String string
	_Int    int64
	_Float  float64
	_Bool   bool
	_Type   int
}

func New(v ...int) *JSON {
	dj := JSON{
		_Type: NULL,
	}

	if len(v) == 1 {
		switch v[0] {
		case OBJECT:
			dj._Object = NewDO()
			dj._Type = OBJECT
		case ARRAY:
			dj._Array = NewDA()
			dj._Type = ARRAY
		case STRING:
			dj._Type = STRING
		case INT:
			dj._Type = INT
		case FLOAT:
			dj._Type = FLOAT
		case BOOL:
			dj._Type = BOOL
		}
	}

	return &dj
}

func NewString(v ...interface{}) *JSON {
	dj := New(STRING)
	if len(v) > 0 {
		dj.Put(v[0])
	}

	return dj
}

func NewInt(v ...interface{}) *JSON {
	dj := New(INT)
	if len(v) > 0 {
		dj.Put(v[0])
	}

	return dj
}

func NewBool(v ...interface{}) *JSON {
	dj := New(BOOL)
	if len(v) > 0 {
		dj.Put(v[0])
	}

	return dj
}

func NewFloat(v ...interface{}) *JSON {
	dj := New(FLOAT)
	if len(v) > 0 {
		dj.Put(v[0])
	}

	return dj
}

func NewObject(v ...interface{}) *JSON {
	dj := New(OBJECT)

	var key string
	var ok bool
	for idx := range v {
		if idx%2 == 0 {
			if key, ok = v[idx].(string); !ok {
				return dj
			}
		} else {
			dj.Put(key, v[idx])
		}
	}

	return dj
}

func NewArray(v ...interface{}) *JSON {
	dj := New(ARRAY)

	for idx := range v {
		dj.Put(v[idx])
	}

	return dj
}

func (m *JSON) SetToObject() *JSON {
	m._Object = NewDO()
	m._Array = nil
	m._Type = OBJECT

	return m
}

func (m *JSON) SetToArray() *JSON {
	m._Array = NewDA()
	m._Object = nil
	m._Type = ARRAY

	return m
}

func (m *JSON) Parse(doc string) *JSON {

	if m._Type != NULL {
		return m
	}

	tdoc := strings.TrimSpace(doc)
	if tdoc == "" {
		return m
	}

	var err error

	if tdoc[0] == '{' {
		m._Object, err = ParseToObject(tdoc)
		if err == nil {
			m._Type = OBJECT
		}
	} else if tdoc[0] == '[' {
		m._Array, err = ParseToArray(tdoc)
		if err == nil {
			m._Type = ARRAY
		}
	} else {
		if strings.EqualFold(tdoc, "null") {
			m._Type = NULL
		} else if strings.EqualFold(tdoc, "true") || strings.EqualFold(tdoc, "false") {
			m._Type = BOOL
			m._Bool, _ = gov.ToBoolean(tdoc)
		} else {
			if gov.IsNumeric(tdoc) {
				if gov.IsInt(tdoc) {
					m._Int, _ = strconv.ParseInt(tdoc, 10, 64)
					m._Type = INT
				} else {
					m._Float, _ = strconv.ParseFloat(tdoc, 64)
					m._Type = FLOAT
				}
			} else {
				m._String = tdoc
				m._Type = STRING
			}
		}
	}

	return m
}

func (m *JSON) Put(v ...interface{}) *JSON {

	if IsEmptyArg(v) {
		return m
	}

	if len(v) == 2 {

		if key, ok := v[0].(string); ok {
			m.PutObject(key, v[1])
		} else {
			for idx := range v {
				m.PutArray(v[idx])
			}
		}

		return m
	}

	if len(v) >= 3 { // must be array
		for idx := range v {
			m.PutArray(v[idx])
		}

		return m
	}

	// length of v must be 1

	if v[0] == nil {
		m._Array = nil
		m._Object = nil
		m._Type = NULL
		return m
	}

	if IsIntType(v[0]) {
		if m._Type == NULL || m._Type == INT {
			m._Int, _ = getIntBase(v[0])
			m._Array = nil
			m._Object = nil
			m._Type = INT
		} else {
			m.PutArray(v[0]) // best effort
		}
		return m
	}

	if IsFloatType(v[0]) {
		if m._Type == NULL || m._Type == FLOAT {
			m._Float, _ = getFloatBase(v[0])
			m._Array = nil
			m._Object = nil
			m._Type = FLOAT
		} else {
			m.PutArray(v[0]) // best effort
		}
		return m
	}

	if IsBoolType(v[0]) {
		if m._Type == NULL || m._Type == BOOL {
			m._Bool, _ = getBoolBase(v[0])
			m._Array = nil
			m._Object = nil
			m._Type = BOOL
		} else {
			m.PutArray(v[0]) // best effort
		}
		return m
	}

	if IsStringType(v[0]) {
		if m._Type == NULL || m._Type == STRING {
			m._String, _ = getStringBase(v[0])
			m._Array = nil
			m._Object = nil
			m._Type = STRING
		} else {
			m.PutArray(v[0]) // best effort
		}
		return m
	}

	switch t := v[0].(type) {
	case map[string]interface{}:
		if m._Type == OBJECT {
			for key := range t {
				m._Object.Put(key, t[key])
			}
		} else {
			m._Object = MapToObject(t)
			m._Array = nil
			m._Type = OBJECT
		}
	case Object:
		if m._Type == OBJECT {
			for key := range map[string]interface{}(t) {
				m._Object.Put(key, t[key])
			}
		} else {
			m._Object = MapToObject(t)
			m._Array = nil
			m._Type = OBJECT
		}
	case *DO:
		if m._Type == OBJECT {
			for key := range t.Map {
				m._Object.Put(key, t.Map[key])
			}
		} else {
			m._Object = t
			m._Array = nil
			m._Type = OBJECT
		}
	case DO:
		if m._Type == OBJECT {
			for key := range t.Map {
				m._Object.Put(key, t.Map[key])
			}
		} else {
			m._Object = &t
			m._Array = nil
			m._Type = OBJECT
		}
	case []interface{}:
		if m._Type == ARRAY {
			m._Array.Put(t)
		} else {
			m._Array = SliceToArray(t)
			m._Object = nil
			m._Type = ARRAY
		}

	case Array:
		if m._Type == ARRAY {
			m._Array.Put([]interface{}(t))
		} else {
			m._Array = SliceToArray(t)
			m._Object = nil
			m._Type = ARRAY
		}
	case *DA:
		if m._Type == ARRAY {
			m._Array.Put(t.Element)
		} else {
			m._Array = t
			m._Object = nil
			m._Type = ARRAY
		}
	case DA:
		if m._Type == ARRAY {
			m._Array.Put(t.Element)
		} else {
			m._Array = &t
			m._Object = nil
			m._Type = ARRAY
		}
	case JSON:
		m = &t
	default:
		if m._Type == ARRAY {
			m._Array.Put(t)
		}
	}

	return m
}

func (m *JSON) PutArray(value ...interface{}) *JSON {
	if m._Type == NULL {
		m._Array = NewDA()
		m._Type = ARRAY
	}

	if m._Type == ARRAY {
		m._Array.Put(value)
	}

	return m
}

func (m *JSON) PutObject(key string, value interface{}) *JSON {
	if m._Type == NULL {
		m._Object = NewDO()
		m._Type = OBJECT
	}

	if m._Type == OBJECT {
		m._Object.Put(key, value)
	}

	return m
}

func (m *JSON) Remove(key interface{}) *JSON {
	switch tkey := key.(type) {
	case string:
		if m._Type == OBJECT {
			m._Object.Remove(tkey)
		}
	case int:
		if m._Type == ARRAY {
			m._Array.Remove(tkey)
		}
	}

	return m
}

func (m *JSON) Interface(key ...interface{}) interface{} {
	if IsEmptyArg(key) {
		switch m._Type {
		case NULL:
			return nil
		case STRING:
			return m._String
		case BOOL:
			return m._Bool
		case INT:
			return m._Int
		case FLOAT:
			return m._Float
		case OBJECT:
			return m._Object
		case ARRAY:
			return m._Array
		}

		return nil
	} else {

		switch tkey := key[0].(type) {
		case string:
			if m._Type == OBJECT {
				if obj, ok := m._Object.Get(tkey); ok {
					return obj
				}
			}
		case int:
			if m._Type == ARRAY {
				if arr, ok := m._Array.Get(tkey); ok {
					return arr
				}
			}
		}
	}

	return nil
}

func (m *JSON) Get(key ...interface{}) (*JSON, bool) {
	if IsEmptyArg(key) {
		return m, true
	} else {

		r := New()
		var element interface{}
		var retOk bool

		switch tkey := key[0].(type) {
		case string:
			if m._Type == OBJECT {
				element, retOk = m._Object.Get(tkey)
			}
		case int:
			if m._Type == ARRAY {
				element, retOk = m._Array.Get(tkey)
			}
		}

		if !retOk {
			return nil, false
		}

		eVal := reflect.ValueOf(element)

		switch t := element.(type) {
		case nil:
			r._Type = NULL
		case string:
			r._String = t
			r._Type = STRING
		case bool:
			r._Bool = t
			r._Type = BOOL
		case uint8, uint16, uint32, uint64, uint:
			intVal := int64(eVal.Uint())
			r._Int = intVal
			r._Type = INT
		case int8, int16, int32, int64, int:
			intVal := eVal.Int()
			r._Int = intVal
			r._Type = INT
		case float32, float64:
			floatVal := eVal.Float()
			r._Float = floatVal
			r._Type = FLOAT
		case DA:
			r._Array = &t
			r._Type = ARRAY
		case DO:
			r._Object = &t
			r._Type = OBJECT
		case *DA:
			r._Array = t
			r._Type = ARRAY
		case *DO:
			r._Object = t
			r._Type = OBJECT
		default:
			return nil, false
		}

		return r, true
	}
}

// The DJSON as return shared Object.

func (m *JSON) Object(key ...interface{}) (*JSON, bool) {

	if m._Type != OBJECT && m._Type != ARRAY {
		return nil, false
	}

	if IsEmptyArg(key) {
		if m._Type == OBJECT {
			return m, true
		}
	} else {

		var ok bool
		var newObject *DO

		switch tkey := key[0].(type) {
		case string:
			if m._Type == OBJECT {
				newObject, ok = m._Object.Object(tkey)
			}
		default:
			kint, oki := getIntBase(key[0])
			if oki && m._Type == ARRAY {
				newObject, ok = m._Array.Object(int(kint))
			}
		}

		if !ok {
			return nil, false
		}

		if newObject != nil {
			return &JSON{
				_Object: newObject,
				_Array:  nil,
				_Type:   OBJECT,
			}, true
		}
	}

	return nil, false
}

// The DJSON as return shared Array.

func (m *JSON) Array(key ...interface{}) (*JSON, bool) {

	if m._Type != OBJECT && m._Type != ARRAY {
		return nil, false
	}

	if IsEmptyArg(key) {
		if m._Type == ARRAY {
			return m, true
		}
	} else {

		var ok bool
		var newArray *DA

		switch tkey := key[0].(type) {
		case string:
			if m._Type == OBJECT {
				newArray, ok = m._Object.Array(tkey)
			}
		default:
			kint, oki := getIntBase(key[0])
			if oki && m._Type == ARRAY {
				newArray, ok = m._Array.Array(int(kint))
			}
		}

		if !ok {
			return nil, false
		}

		if newArray != nil {
			return &JSON{
				_Object: nil,
				_Array:  newArray,
				_Type:   ARRAY,
			}, true
		}

	}

	return nil, false
}

func (m *JSON) Int(key ...interface{}) int64 {

	if IsEmptyArg(key) {

		switch m._Type {
		case ARRAY, OBJECT, NULL:
			return 0
		case BOOL:
			if m._Bool {
				return 1
			}
			return 0
		case STRING:
			if iVal, err := strconv.ParseInt(m._String, 10, 64); err == nil {
				return iVal
			}
			return 0
		case INT:
			return m._Int
		case FLOAT:
			return int64(m._Float)
		}

	} else {

		var dv int64

		if len(key) >= 2 {
			v, ok := getIntBase(key[1])
			if ok {
				dv = v
			}
		}

		switch tkey := key[0].(type) {
		case string:
			if m._Type == OBJECT {
				if iVal, ok := m._Object.Int(tkey); ok {
					return iVal
				}
			}
		default:
			kint, ok := getIntBase(key[0])
			if ok && m._Type == ARRAY {
				if iVal, ok := m._Array.Int(int(kint)); ok {
					return iVal
				}
			}
		}

		return dv
	}

	return 0 // zero value
}

func (m *JSON) Bool(key ...interface{}) bool {
	if IsEmptyArg(key) {

		switch m._Type {
		case NULL, FLOAT, ARRAY, OBJECT:
			return false
		case STRING:
			if strings.EqualFold(m._String, "true") {
				return true
			}
			return false
		case INT:
			return m._Int == 1
		case BOOL:
			return m._Bool
		}

	} else {

		var dv bool

		if len(key) >= 2 {
			dv = key[1].(bool)
		}

		switch tkey := key[0].(type) {
		case string:
			if m._Type == OBJECT {
				if bVal, ok := m._Object.Bool(tkey); ok {
					return bVal
				}
			}
		default:
			kint, ok := getIntBase(key[0])
			if ok && m._Type == ARRAY {
				if bVal, ok := m._Array.Bool(int(kint)); ok {
					return bVal
				}
			}
		}

		return dv

	}

	return false // zero value
}

func (m *JSON) Float(key ...interface{}) float64 {
	if IsEmptyArg(key) {

		switch m._Type {
		case NULL, ARRAY, OBJECT:
			return 0
		case STRING:
			if fVal, err := strconv.ParseFloat(m._String, 64); err == nil {
				return fVal
			}
			return 0
		case BOOL:
			if m._Bool {
				return 1
			}
			return 0
		case INT:
			return float64(m._Int)
		case FLOAT:
			return m._Float
		}

	} else {

		var dv float64

		if len(key) >= 2 {
			v, ok := getFloatBase(key[1])
			if ok {
				dv = v
			}
		}

		switch tkey := key[0].(type) {
		case string:
			if m._Type == OBJECT {
				if fVal, ok := m._Object.Float(tkey); ok {
					return fVal
				}
			}
		default:
			kint, ok := getIntBase(key[0])
			if ok && m._Type == ARRAY {
				if fVal, ok := m._Array.Float(int(kint)); ok {
					return fVal
				}
			}
		}

		return dv
	}

	return 0 // zero value
}

func (m *JSON) String(key ...interface{}) string {

	if IsEmptyArg(key) {
		return m.ToString()
	} else {
		var dv string

		if len(key) >= 2 {
			dv = key[1].(string)
		}

		switch tkey := key[0].(type) {
		case string:
			if m._Type == OBJECT {
				if m._Object.HasKey(tkey) {
					return m._Object.String(tkey)
				} else {
					return dv
				}
			} else {
				return tkey // maybe default
			}
		default:
			kint, ok := getIntBase(key[0])
			if ok {
				if m._Type == ARRAY {
					if iVal, ok := m._Array.String2(int(kint)); ok {
						return iVal
					} else {
						return dv
					}
				} else {
					return strconv.Itoa(int(kint)) // maybe default
				}
			}
		}

		return dv

	}
}

func (m *JSON) ToString() string {

	switch m._Type {
	case NULL:
		return "null"
	case STRING:
		return m._String
	case INT:
		intStr, ok := getStringBase(m._Int)
		if !ok {
			return ""
		}
		return intStr
	case FLOAT:
		floatStr, ok := getStringBase(m._Float)
		if !ok {
			return ""
		}
		return floatStr
	case BOOL:
		return gov.ToString(m._Bool)
	case OBJECT:
		return m._Object.ToString()
	case ARRAY:
		return m._Array.ToString()
	}

	return "" // zero value
}

func (m *JSON) Rename(from, to string) bool {
	if m._Type != OBJECT {
		return false
	}

	return m._Object.Rename(from, to)
}

func (m *JSON) ReplaceAt(k interface{}, v interface{}) *JSON {
	switch tkey := k.(type) {
	case string:
		if m._Type == OBJECT {
			if m._Object.HasKey(tkey) {
				m._Object.Put(tkey, v)
			}
		}
	default:
		kint, ok := getIntBase(k)
		if ok && m._Type == ARRAY {
			if m._Array.Size() > int(kint) {
				m._Array.ReplaceAt(int(kint), v)
			}
		}
	}

	return m
}

func (m *JSON) Seek(seekp ...int) bool {
	if m._Type == ARRAY {
		m._Array.Seek(seekp...)
		return true
	}

	return false
}

func (m *JSON) Next() bool {
	if m._Type == ARRAY {
		return m._Array.Next()
	}
	return false
}

func (m *JSON) Scan() *JSON {
	if m._Type == ARRAY {
		v, ok := m._Array.Scan()
		if !ok {
			return nil
		}

		ret := New()
		switch t := v.(type) {
		case string:
			ret._Type = STRING
			ret._String = t
		case bool:
			ret._Type = BOOL
			ret._Bool = t
		case int, int8, int16, int32, int64:
			ret._Type = INT
			ret._Int = reflect.ValueOf(t).Int()
		case uint, uint8, uint16, uint32, uint64:
			ret._Type = INT
			ret._Int = int64(reflect.ValueOf(t).Uint())
		case float32, float64:
			ret._Type = FLOAT
			ret._Float = reflect.ValueOf(t).Float()
		case *DA:
			ret._Type = ARRAY
			ret._Array = t
		case *DO:
			ret._Type = OBJECT
			ret._Object = t
		case DA:
			ret._Type = ARRAY
			ret._Array = &t
		case DO:
			ret._Type = OBJECT
			ret._Object = &t
		case *JSON:
			ret = t
		case JSON:
			ret = &t
		case nil:
			ret._Type = NULL
		}

		return ret
	}

	return nil
}

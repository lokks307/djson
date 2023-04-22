package djson

import (
	"reflect"
	"strings"

	"github.com/volatiletech/null/v8"
)

func (m *JSON) Size() int {
	return m.Len()
}

func (m *JSON) Len() int {
	if m._Type == NULL {
		return 0
	}

	if m._Type == ARRAY {
		return m._Array.Length()
	}

	if m._Type == OBJECT {
		return m._Object.Length()
	}

	return 1
}

func (m *JSON) HasKey(key interface{}) bool {
	switch tkey := key.(type) {
	case string:
		if m._Type == OBJECT {
			return m._Object.HasKey(tkey)
		}
	case int:
		if m._Type == ARRAY {
			return tkey >= 0 && m._Array.Size() > tkey
		}
	}

	return false
}

func (m *JSON) toFieldsValue(val reflect.Value, tags ...string) {

	for i := 0; i < val.NumField(); i++ {
		eval := val.Field(i)
		eachType := val.Type().Field(i)
		eachTag := eachType.Tag.Get("json")

		if !eval.CanSet() || !m.HasKey(eachTag) {
			continue
		}

		if len(tags) > 0 && !inTags(eachTag, tags...) {
			continue
		}

		eachKind := eachType.Type.Kind()

		if eachKind == reflect.Struct {

			switch eachType.Type.String() {
			case "null.String":
				eval.FieldByName("String").SetString(m.String(eachTag))
				eval.FieldByName("Valid").SetBool(true)
			case "null.Bool":
				eval.FieldByName("Bool").SetBool(m.Bool(eachTag))
				eval.FieldByName("Valid").SetBool(true)
			case "null.Float32":
				eval.FieldByName("Float32").SetFloat(m.Float(eachTag))
				eval.FieldByName("Valid").SetBool(true)
			case "null.Float64":
				eval.FieldByName("Float64").SetFloat(m.Float(eachTag))
				eval.FieldByName("Valid").SetBool(true)
			case "null.Int":
				eval.FieldByName("Int").SetInt(m.Int(eachTag))
				eval.FieldByName("Valid").SetBool(true)
			case "null.Int8":
				eval.FieldByName("Int8").SetInt(m.Int(eachTag))
				eval.FieldByName("Valid").SetBool(true)
			case "null.Int16":
				eval.FieldByName("Int16").SetInt(m.Int(eachTag))
				eval.FieldByName("Valid").SetBool(true)
			case "null.Int32":
				eval.FieldByName("Int32").SetInt(m.Int(eachTag))
				eval.FieldByName("Valid").SetBool(true)
			case "null.Int64":
				eval.FieldByName("Int64").SetInt(m.Int(eachTag))
				eval.FieldByName("Valid").SetBool(true)
			case "null.Uint":
				eval.FieldByName("Uint").SetUint(uint64(m.Int(eachTag)))
				eval.FieldByName("Valid").SetBool(true)
			case "null.Uint8":
				eval.FieldByName("Uint8").SetUint(uint64(m.Int(eachTag)))
				eval.FieldByName("Valid").SetBool(true)
			case "null.Uint16":
				eval.FieldByName("Uint16").SetUint(uint64(m.Int(eachTag)))
				eval.FieldByName("Valid").SetBool(true)
			case "null.Uint32":
				eval.FieldByName("Uint32").SetUint(uint64(m.Int(eachTag)))
				eval.FieldByName("Valid").SetBool(true)
			case "null.Uint64":
				eval.FieldByName("Uint64").SetUint(uint64(m.Int(eachTag)))
				eval.FieldByName("Valid").SetBool(true)
			default:

				if oJson, ok := m.Object(eachTag); ok {
					oJson.toFieldsValue(eval, downDepthWW(tags)...)
				}

			}

		} else {

			switch eachType.Type.String() {
			case "int", "int8", "int16", "int32", "int64":
				eval.SetInt(m.Int(eachTag))
			case "uint", "uint8", "uint16", "uint32", "uint64":
				eval.SetUint(uint64(m.Int(eachTag)))
			case "float32", "float64":
				eval.SetFloat(m.Float(eachTag))
			case "string":
				eval.SetString(m.String(eachTag))
			case "bool":
				eval.SetBool(m.Bool(eachTag))
			}
		}
	}
}

func (m *JSON) ToFields(st interface{}, tags ...string) {
	target := reflect.ValueOf(st)
	elements := target.Elem()
	m.toFieldsValue(elements, tags...)
}

func (m *JSON) fromFieldsValue(val reflect.Value, tags ...string) {

	kind := val.Type().Kind()

	if kind == reflect.Array || kind == reflect.Slice {

		for i := 0; i < val.Len(); i++ {
			eachVal := val.Index(i)
			eachType := eachVal.Type()

			switch eachVal.Kind() {
			case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
				m.PutArray(eachVal.Int())
			case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
				m.PutArray(eachVal.Uint())
			case reflect.Bool:
				m.PutArray(eachVal.Bool())
			case reflect.String:
				m.PutArray(eachVal.String())
			case reflect.Float32, reflect.Float64:
				m.PutArray(eachVal.Float())
			case reflect.Array, reflect.Slice:
				sJson := New()
				sJson.SetToArray()
				sJson.fromFieldsValue(eachVal, downDepthWW(tags)...)
				m.PutArray(sJson)
			case reflect.Struct, reflect.Map:
				switch eachType.String() {
				case "null.String":
					m.PutArray(eachVal.FieldByName("String").String())
				case "null.Bool":
					m.PutArray(eachVal.FieldByName("Bool").Bool())
				case "null.Float32":
					m.PutArray(eachVal.FieldByName("Float32").Float())
				case "null.Float64":
					m.PutArray(eachVal.FieldByName("Float64").Float())
				case "null.Int":
					m.PutArray(eachVal.FieldByName("Int").Int())
				case "null.Int8":
					m.PutArray(eachVal.FieldByName("Int8").Int())
				case "null.Int16":
					m.PutArray(eachVal.FieldByName("Int16").Int())
				case "null.Int32":
					m.PutArray(eachVal.FieldByName("Int32").Int())
				case "null.Int64":
					m.PutArray(eachVal.FieldByName("Int64").Int())
				case "null.Uint":
					m.PutArray(eachVal.FieldByName("Uint").Uint())
				case "null.Uint8":
					m.PutArray(eachVal.FieldByName("Uint8").Uint())
				case "null.Uint16":
					m.PutArray(eachVal.FieldByName("Uint16").Uint())
				case "null.Uint32":
					m.PutArray(eachVal.FieldByName("Uint32").Uint())
				case "null.Uint64":
					m.PutArray(eachVal.FieldByName("Uint64").Uint())
				default:
					sJson := New()
					sJson.SetToObject()
					sJson.fromFieldsValue(eachVal, downDepthWW(tags)...)
					m.PutArray(sJson)
				}
			default:
				m.PutArray(nil)
			}

		}

	} else if kind == reflect.Struct {

		for i := 0; i < val.NumField(); i++ {
			eachVal := val.Field(i)
			eachType := val.Type().Field(i)
			eachTag := eachType.Tag.Get("json")

			if len(tags) > 0 && !inTags(eachTag, tags...) {
				continue
			}

			eachKind := eachType.Type.Kind()

			if eachKind == reflect.Struct || eachKind == reflect.Map {

				switch eachType.Type.String() {
				case "null.String":
					m.Put(eachTag, eachVal.FieldByName("String").String())
				case "null.Bool":
					m.Put(eachTag, eachVal.FieldByName("Bool").Bool())
				case "null.Float32":
					m.Put(eachTag, eachVal.FieldByName("Float32").Float())
				case "null.Float64":
					m.Put(eachTag, eachVal.FieldByName("Float64").Float())
				case "null.Int":
					m.Put(eachTag, eachVal.FieldByName("Int").Int())
				case "null.Int8":
					m.Put(eachTag, eachVal.FieldByName("Int8").Int())
				case "null.Int16":
					m.Put(eachTag, eachVal.FieldByName("Int16").Int())
				case "null.Int32":
					m.Put(eachTag, eachVal.FieldByName("Int32").Int())
				case "null.Int64":
					m.Put(eachTag, eachVal.FieldByName("Int64").Int())
				case "null.Uint":
					m.Put(eachTag, eachVal.FieldByName("Uint").Uint())
				case "null.Uint8":
					m.Put(eachTag, eachVal.FieldByName("Uint8").Uint())
				case "null.Uint16":
					m.Put(eachTag, eachVal.FieldByName("Uint16").Uint())
				case "null.Uint32":
					m.Put(eachTag, eachVal.FieldByName("Uint32").Uint())
				case "null.Uint64":
					m.Put(eachTag, eachVal.FieldByName("Uint64").Uint())
				default:
					sJson := New()
					sJson.SetToObject()
					sJson.fromFieldsValue(eachVal, downDepthWW(tags)...)
					m.Put(eachTag, sJson)
				}
			} else if eachKind == reflect.Array || eachKind == reflect.Slice {

				sJson := New()
				sJson.SetToArray()
				sJson.fromFieldsValue(eachVal, downDepthWW(tags)...)
				m.Put(eachTag, sJson)

			} else {

				switch eachType.Type.String() {
				case "int", "int8", "int16", "int32", "int64":
					m.Put(eachTag, eachVal.Int())
				case "uint", "uint8", "uint16", "uint32", "uint64":
					m.Put(eachTag, eachVal.Uint())
				case "float32", "float64":
					m.Put(eachTag, eachVal.Float())
				case "string":
					m.Put(eachTag, eachVal.String())
				case "bool":
					m.Put(eachTag, eachVal.Bool())
				}
			}
		}
	} else if kind == reflect.Map {

		for _, e := range val.MapKeys() {
			eachKey, ok := e.Interface().(string)
			if !ok {
				continue
			}

			if len(tags) > 0 && !inTags(eachKey, tags...) {
				continue
			}

			eachVal := val.MapIndex(e)

			switch t := eachVal.Interface().(type) {
			case int:
				m.Put(eachKey, t)
			case int8:
				m.Put(eachKey, t)
			case int16:
				m.Put(eachKey, t)
			case int32:
				m.Put(eachKey, t)
			case int64:
				m.Put(eachKey, t)
			case uint:
				m.Put(eachKey, t)
			case uint8:
				m.Put(eachKey, t)
			case uint16:
				m.Put(eachKey, t)
			case uint32:
				m.Put(eachKey, t)
			case uint64:
				m.Put(eachKey, t)
			case float32:
				m.Put(eachKey, t)
			case float64:
				m.Put(eachKey, t)
			case string:
				m.Put(eachKey, t)
			case bool:
				m.Put(eachKey, t)
			case nil:
				m.Put(eachKey, t)
			case null.String:
				m.Put(eachKey, t.String)
			case null.Bool:
				m.Put(eachKey, t.Bool)
			case null.Int:
				m.Put(eachKey, t.Int)
			case null.Int8:
				m.Put(eachKey, t.Int8)
			case null.Int16:
				m.Put(eachKey, t.Int16)
			case null.Int32:
				m.Put(eachKey, t.Int32)
			case null.Int64:
				m.Put(eachKey, t.Int64)
			case null.Uint:
				m.Put(eachKey, t.Uint)
			case null.Uint8:
				m.Put(eachKey, t.Uint8)
			case null.Uint16:
				m.Put(eachKey, t.Uint16)
			case null.Uint32:
				m.Put(eachKey, t.Uint32)
			case null.Uint64:
				m.Put(eachKey, t.Uint64)
			case null.Float32:
				m.Put(eachKey, t.Float32)
			case null.Float64:
				m.Put(eachKey, t.Float64)
			default:

				skind := reflect.ValueOf(t).Type().Kind()

				if skind == reflect.Struct || skind == reflect.Map {
					sJson := New()
					sJson.SetToObject()
					sJson.FromFields(t, downDepthWW(tags)...)
					m.Put(eachKey, sJson)
				}

			}

		}
	}
}

func (m *JSON) FromFields(st interface{}, tags ...string) *JSON {
	baseValue := reflect.ValueOf(st)

	kind := baseValue.Type().Kind()

	if kind == reflect.Array || kind == reflect.Slice {

		m.SetToArray()
		m.fromFieldsValue(baseValue, tags...)

	} else if kind == reflect.Struct || kind == reflect.Map {

		m.SetToObject()
		m.fromFieldsValue(baseValue, tags...)

	}

	return m
}

func downDepthWW(tags []string) []string {
	tags2 := make([]string, 0)
	for idx := range tags {
		tmp := strings.Split(tags[idx], ".")
		tmp2 := strings.Join(tmp[1:], ".")
		if tmp2 != "" {
			tags2 = append(tags2, tmp2)
		} else {
			tags2 = append(tags2, "")
		}
	}

	return tags2
}

func inTags(idv string, tags ...string) bool {
	for idx := range tags {
		tmp := strings.Split(tags[idx], ".")
		if tmp[0] == idv {
			return true
		}
	}

	return false
}

func (m *JSON) doSort(isAsc bool, k ...interface{}) bool {
	var tArray *DA

	if len(k) == 0 {
		if m._Type == ARRAY {
			tArray = m._Array
		}
	}

	if len(k) > 0 {

		if m._Type == OBJECT {
			if key, ok := k[0].(string); ok {
				if da, ok := m._Object.GetAsArray(key); ok {
					tArray = da
				}
			}
		} else if m._Type == ARRAY {
			if idx, ok := k[0].(int); ok {
				if da, ok := m._Array.GetAsArray(idx); ok {
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
	return m.doSort(true, k...)
}

func (m *JSON) SortDesc(k ...interface{}) bool {
	return m.doSort(false, k...)
}

func (m *JSON) SortObjectArray(isAsc bool, key string) bool {
	if m._Type != ARRAY {
		return false
	}

	return m._Array.SortObject(isAsc, key)
}

func (m *JSON) SortObjectArrayAsc(key string) bool {
	return m.SortObjectArray(true, key)
}

func (m *JSON) SortObjectArrayDesc(key string) bool {
	return m.SortObjectArray(false, key)
}

func (m *JSON) Equal(t *JSON) bool {
	if m._Type != t._Type {
		return false
	}

	switch m._Type {
	case NULL:
		return true
	case BOOL:
		return m._Bool == t._Bool
	case INT:
		return m._Int == t._Int
	case FLOAT:
		return m._Float == t._Float
	case STRING:
		return m._String == t._String
	case OBJECT:
		return m._Object.Equal(t._Object)
	case ARRAY:
		return m._Array.Equal(t._Array)
	}

	return false
}

func (m *JSON) Clone() *JSON {
	t := New(m._Type)

	switch m._Type {
	case NULL:
	case BOOL:
		t._Bool = m._Bool
	case INT:
		t._Int = m._Int
	case FLOAT:
		t._Float = m._Float
	case STRING:
		t._String = m._String
	case OBJECT:
		t._Object = m._Object.Clone()
	case ARRAY:
		t._Array = m._Array.Clone()
	}

	return t
}

func (m *JSON) HasKeys(k ...interface{}) bool {
	for i := range k {
		if !m.HasKey(k[i]) {
			return false
		}
	}

	return true
}

func (m *JSON) GetKeys(k ...interface{}) []string {
	rk := make([]string, 0)

	if IsEmptyArg(k) {
		if m._Type != OBJECT {
			return rk
		}

		for k := range m._Object.Map {
			rk = append(rk, k)
		}

		return rk
	}

	if t, ok := m.Object(k[0]); ok {
		return t.GetKeys()
	}

	return rk
}

func (m *JSON) Find(key string, val string) *JSON {
	if key == "" || m._Type != ARRAY {
		return nil
	}

	for i := 0; i < m.Len(); i++ {
		each, ok := m.Object(i)
		if !ok {
			continue
		}

		if each.String(key) == val {
			return each
		}
	}

	return nil
}

func (m *JSON) Append(arrJson *JSON) *JSON {
	if arrJson == nil || m._Type != ARRAY || !arrJson.IsArray() {
		return m
	}

	for i := 0; i < arrJson.Len(); i++ {
		m.PutArray(arrJson._Array.Element[i])
	}

	return m
}

func IsEmptyArg(key []interface{}) bool {
	return len(key) == 0 || (len(key) == 1 && key[0] == "")
}

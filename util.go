package djson

import (
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
	"strconv"
	"strings"

	gov "github.com/asaskevich/govalidator"
)

func MapToObject(dmap map[string]interface{}) *DO {
	nObj := NewDO()
	for k, v := range dmap {
		nObj.Put(k, v)
	}
	return nObj
}

func SliceToArray(dslice []interface{}) *DA {
	nArr := NewDA()
	nArr.Put(dslice)
	return nArr
}

func ObjectToMap(obj *DO) map[string]interface{} {
	wMap := make(map[string]interface{})

	if obj == nil {
		return wMap
	}

	for k, v := range obj.Map {
		switch t := v.(type) {
		case DA:
			wMap[k] = ArrayToSlice(&t)
		case DO:
			wMap[k] = ObjectToMap(&t)
		case *DA:
			wMap[k] = ArrayToSlice(t)
		case *DO:
			wMap[k] = ObjectToMap(t)
		default:
			wMap[k] = v
		}
	}

	return wMap
}

func ArrayToSlice(arr *DA) []interface{} {

	wArray := make([]interface{}, 0)

	if arr == nil {
		return wArray
	}

	for idx := range arr.Element {
		switch t := arr.Element[idx].(type) {
		case DA:
			wArray = append(wArray, ArrayToSlice(&t))
		case DO:
			wArray = append(wArray, ObjectToMap(&t))
		case *DA:
			wArray = append(wArray, ArrayToSlice(t))
		case *DO:
			wArray = append(wArray, ObjectToMap(t))
		default:
			wArray = append(wArray, t)
		}
	}

	return wArray
}

func getStringBase(v interface{}) (string, bool) {
	if v == nil {
		return "nil", true
	}

	if IsInTypes(v, "string", "bool", "float32", "float64") {
		return fmt.Sprintf("%v", v), true
	}

	if IsInTypes(v, "int", "uint", "int8", "uint8", "int16", "uint16", "int32", "uint32", "int64", "uint64") {
		return fmt.Sprintf("%d", v), true
	}

	return "", false
}

func getBoolBase(v interface{}) (bool, bool) {
	if IsInTypes(v, "int", "uint", "int8", "uint8", "int16", "uint16", "int32", "uint32", "int64", "uint64") {
		intVal, _ := gov.ToInt(v)
		if intVal == 0 {
			return false, true
		}
	}

	if IsInTypes(v, "string") {
		if strVal, ok := v.(string); ok {
			if strings.EqualFold(strVal, "true") {
				return true, true
			} else if strings.EqualFold(strVal, "false") {
				return false, true
			}
		}
	}

	if IsInTypes(v, "bool") {
		if boolVal, ok := v.(bool); ok {
			return boolVal, true
		}
	}

	return false, false
}

func getFloatBase(v interface{}) (float64, bool) {
	if floatVal, err := gov.ToFloat(v); err != nil {
		return 0, false
	} else {
		return floatVal, true
	}
}

func getIntBase(v interface{}) (int64, bool) {
	if intVal, err := gov.ToInt(v); err != nil {
		return 0, false
	} else {
		return intVal, true
	}
}

func IsBaseType(v interface{}) bool {
	return IsInTypes(v, "string", "bool", "int", "uint", "int8", "uint8", "int16", "uint16", "int32", "uint32", "int64", "uint64", "float32", "float64")
}

func IsIntType(v interface{}) bool {
	return IsInTypes(v, "int", "uint", "int8", "uint8", "int16", "uint16", "int32", "uint32", "int64", "uint64")
}

func IsFloatType(v interface{}) bool {
	return IsInTypes(v, "float32", "float64")
}

func IsBoolType(v interface{}) bool {
	return IsInTypes(v, "bool")
}

func IsStringType(v interface{}) bool {
	return IsInTypes(v, "string")
}

func IsInTypes(v interface{}, types ...string) bool {
	var vTypeStr string
	if v == nil {
		vTypeStr = "nil"
	} else {
		vTypeStr = reflect.TypeOf(v).String()
	}

	for idx := range types {
		if vTypeStr == types[idx] {
			return true
		}
	}

	return false
}

func ParseToObject(doc string) (*DO, error) {
	var data map[string]interface{}

	d := json.NewDecoder(strings.NewReader(doc))
	d.UseNumber()

	if err := d.Decode(&data); err != nil {
		return nil, errors.New("not Object")
	}

	return ParseObject(data), nil

}

func ParseToArray(doc string) (*DA, error) {
	var data []interface{}

	d := json.NewDecoder(strings.NewReader(doc))
	d.UseNumber()

	if err := d.Decode(&data); err != nil {
		return nil, errors.New("not Array")
	}

	return ParseArray(data), nil
}

func ParseObject(data map[string]interface{}) *DO {
	obj := NewDO()
	for k, v := range data {
		if IsBaseType(v) {
			obj.Put(k, v)
			continue
		}

		if n, ok := v.(json.Number); ok {
			if i, err := n.Int64(); err == nil {
				obj.Put(k, i)
				continue
			}
			if f, err := n.Float64(); err == nil {
				obj.Put(k, f)
				continue
			}
		}

		switch tValue := v.(type) {
		case []interface{}: // Array
			nArr := ParseArray(tValue)
			obj.Put(k, nArr)
		case map[string]interface{}: // Object
			nObj := ParseObject(tValue)
			obj.Put(k, nObj)
		case nil: // null
			obj.Put(k, nil)
		}
	}

	return obj
}

func ParseArray(data []interface{}) *DA {
	arr := NewDA()

	for idx := range data {
		if IsBaseType(data[idx]) {
			arr.Put(data[idx])
			continue
		}

		if n, ok := data[idx].(json.Number); ok {
			if i, err := n.Int64(); err == nil {
				arr.Put(i)
				continue
			}
			if f, err := n.Float64(); err == nil {
				arr.Put(f)
				continue
			}
		}

		switch tValue := data[idx].(type) {
		case []interface{}: // Array
			nArr := ParseArray(tValue)
			arr.PutArray(nArr)
		case map[string]interface{}: // Object
			nObj := ParseObject(tValue)
			arr.Put(nObj)
		case nil: // null
			arr.Put(nil)
		}
	}

	return arr
}

func PathTokenizer(path string) []interface{} {
	rstack := NewRuneStack()
	token := make([]rune, 0)
	inTokens := make([]string, 0)

	prev := rune(0)
	var depthL int

	for _, each := range path {

		peek := rstack.Peek()

		if depthL == 0 {
			if each == '[' && prev != '\\' {
				rstack.Push(each)
				token = make([]rune, 0)
				depthL = 1
			} else {
				token = append(token, each)
			}
		} else if depthL == 1 {
			if peek == '[' && each == ']' && prev != '\\' {
				if len(token) > 0 {
					inTokens = append(inTokens, string(token))
					token = make([]rune, 0)
				}
				rstack.Pop()
				depthL = 0
			} else if (each == '"' || each == '\'') && prev != '\\' {
				rstack.Push(each)
				depthL = 2
			} else {
				token = append(token, each)
			}
		} else if depthL == 2 {

			if (peek == '"' && each == '"' && prev != '\\') || (peek == '\'' && each == '\'' && prev != '\\') {
				if len(token) > 0 {
					inTokens = append(inTokens, string(token))
					token = make([]rune, 0)
				}
				rstack.Pop()
				depthL = 1
			} else {
				token = append(token, each)
			}
		}

		prev = each
	}

	outTokens := make([]interface{}, 0)
	for idx := range inTokens {
		if intVal, err := strconv.Atoi(inTokens[idx]); err == nil {
			outTokens = append(outTokens, intVal)
		} else {
			outTokens = append(outTokens, inTokens[idx])
		}
	}

	return outTokens
}

func MustSome(opt *JSON, keys ...interface{}) bool {
	if opt == nil {
		return false
	}

	for _, key := range keys {
		if opt.String(key) == "" {
			return false
		}
	}

	return true
}

func MustGetObject(o *JSON, key interface{}) *JSON {
	emptyObjectJson := New(OBJECT)
	if o == nil || key == "" {
		return emptyObjectJson
	}

	r, ok := o.Object(key)
	if !ok || r == nil {
		return emptyObjectJson
	}
	return r
}

func MustGetArray(o *JSON, key interface{}) *JSON {
	emptyArrayJson := New(ARRAY)
	if o == nil || key == "" {
		return emptyArrayJson
	}

	r, ok := o.Array(key)
	if !ok || r == nil {
		return emptyArrayJson
	}
	return r
}

func MustGetString(o *JSON, key interface{}) string {
	if o == nil || key == "" {
		return ""
	}
	return o.String(key)
}

func MustGetBool(o *JSON, key interface{}) bool {
	if o == nil || key == "" {
		return false
	}
	return o.Bool(key)
}

func MustGetInt(o *JSON, key interface{}) int64 {
	if o == nil || key == "" {
		return 0
	}
	return o.Int(key)
}

// o = object djson
func MustGetStringSlice(o *JSON, key interface{}) []string {
	emptyArraySlice := make([]string, 0)
	if o == nil || key == "" {
		return emptyArraySlice
	}

	r, ok := o.Array(key)
	if !ok || r == nil {
		return emptyArraySlice
	}
	return ArrayJsonToStringSlice(r)
}

func MustGetInt64Slice(o *JSON, key interface{}) []int64 {
	emptyArraySlice := make([]int64, 0)
	if o == nil || key == "" {
		return emptyArraySlice
	}

	r, ok := o.Array(key)
	if !ok || r == nil {
		return emptyArraySlice
	}
	return JsonToInt64Slice(r)
}

func MustGetIntSlice(o *JSON, key interface{}) []int {
	emptyArraySlice := make([]int, 0)
	if o == nil || key == "" {
		return emptyArraySlice
	}

	r, ok := o.Array(key)
	if !ok || r == nil {
		return emptyArraySlice
	}
	return JsonToIntSlice(r)
}

func StringSliceToJson(ss []string) *JSON {
	retJson := New(ARRAY)
	for _, s := range ss {
		if s != "" {
			retJson.PutArray(s)
		}
	}

	return retJson
}

func StringMapToArray(ss map[string]string) Array {
	var retArray Array
	for _, s := range ss {
		if s != "" {
			retArray = append(retArray, s)
		}
	}

	return retArray
}

func Int64SliceToArray(is []int64) Array {
	var retArray Array
	for _, s := range is {
		retArray = append(retArray, s)
	}

	return retArray
}

func StringSliceToArray(ss []string) Array {
	var retArray Array
	for _, s := range ss {
		if s != "" {
			retArray = append(retArray, s)
		}
	}

	return retArray
}

func IntSliceToJson(ss []int) *JSON {
	retJson := New(ARRAY)
	for _, s := range ss {
		retJson.PutArray(s)
	}

	return retJson
}

// js must be array json
func ArrayJsonToStringSlice(js *JSON, key ...string) []string {
	if js == nil || !js.IsArray() || js.Size() == 0 {
		return []string{}
	}

	ss := make([]string, 0)
	js.Seek()
	if len(key) > 0 {
		for js.Next() {
			ec := js.Scan()
			ss = append(ss, ec.String(key[0]))
		}
	} else {
		for js.Next() {
			ec := js.Scan()
			ss = append(ss, ec.ToString())
		}
	}
	return ss
}

func JsonToIntSlice(js *JSON, key ...string) []int {
	if js == nil || !js.IsArray() || js.Size() == 0 {
		return []int{}
	}

	ss := make([]int, 0)
	js.Seek()
	if len(key) > 0 {
		for js.Next() {
			ec := js.Scan()
			ss = append(ss, int(ec.Int(key[0])))
		}
	} else {
		for js.Next() {
			ec := js.Scan()
			ss = append(ss, int(ec.Int()))
		}
	}
	return ss
}

func JsonToInt64Slice(js *JSON, key ...string) []int64 {
	if js == nil || !js.IsArray() || js.Size() == 0 {
		return []int64{}
	}

	ss := make([]int64, 0)
	js.Seek()
	if len(key) > 0 {
		for js.Next() {
			ec := js.Scan()
			ss = append(ss, ec.Int(key[0]))
		}
	} else {
		for js.Next() {
			ec := js.Scan()
			ss = append(ss, ec.Int())
		}
	}
	return ss
}

func JsonFilter(js *JSON, keys ...string) *JSON {
	ret := New(ARRAY)
	if js == nil || !js.IsArray() {
		return ret
	}

	js.Seek()
	for js.Next() {
		ec := js.Scan()
		nec := New()
		for _, k := range keys {
			if k == "" {
				continue
			}

			v, ok := ec.Get(k)
			if ok {
				nec.Put(k, v)
			} else {
				nec.Put(k, nil)
			}
		}

		ret.PutArray(nec)
	}

	return ret
}

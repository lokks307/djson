package djson

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
	return JsonToStringSlice(r)
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
func JsonToStringSlice(js *JSON, key ...string) []string {
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

func ArrayJsonToStringSlice(js *JSON, key ...string) []string {
	return JsonToStringSlice(js, key...)
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

func SliceToJsonString[T string | numbers](ss []T) string {
	return NewArray().Put(ss).ToString()
}

package djson

func (m *JSON) getTypeSimple(key interface{}) string {

	switch m._Type {
	case OBJECT:
		if key, tok := key.(string); tok {
			if typeStr, ok := m._Object.Type(key); ok {
				return typeStr
			}
		}
	case ARRAY:
		if idx, tok := key.(int); tok {
			if typeStr, ok := m._Array.Type(idx); ok {
				return typeStr
			}
		}
	}

	return ""
}

func (m *JSON) isSameType(key interface{}, inTypeStr string) bool {

	return m.getTypeSimple(key) == inTypeStr
}

func (m *JSON) IsBool(key ...interface{}) bool {
	if IsEmptyArg(key) {
		return m._Type == BOOL
	}

	return m.isSameType(key[0], "bool")
}

func (m *JSON) IsInt(key ...interface{}) bool {
	if IsEmptyArg(key) {
		return m._Type == INT
	}

	return m.isSameType(key[0], "int")
}

func (m *JSON) IsNumeric(key ...interface{}) bool {
	if IsEmptyArg(key) {
		return m._Type == FLOAT || m._Type == INT
	}

	return m.isSameType(key[0], "int") || m.isSameType(key[0], "float")
}

func (m *JSON) IsFloat(key ...interface{}) bool {
	if IsEmptyArg(key) {
		return m._Type == FLOAT
	}

	return m.isSameType(key[0], "float")
}

func (m *JSON) IsString(key ...interface{}) bool {
	if IsEmptyArg(key) {
		return m._Type == STRING
	}

	return m.isSameType(key[0], "string")
}

func (m *JSON) IsNull(key ...interface{}) bool {
	if IsEmptyArg(key) {
		return m._Type == NULL
	}

	return m.isSameType(key[0], "null")
}

func (m *JSON) IsObject(key ...interface{}) bool {
	if IsEmptyArg(key) {
		return m._Type == OBJECT
	}

	return m.isSameType(key[0], "object")
}

func (m *JSON) IsArray(key ...interface{}) bool {
	if IsEmptyArg(key) {
		return m._Type == ARRAY
	}

	return m.isSameType(key[0], "array")
}

func (m *JSON) Type(key ...interface{}) string {
	if IsEmptyArg(key) {
		switch m._Type {
		case NULL:
			return "null"
		case OBJECT:
			return "object"
		case ARRAY:
			return "array"
		case STRING:
			return "string"
		case INT:
			return "int"
		case FLOAT:
			return "float"
		case BOOL:
			return "bool"
		}

		return ""
	}

	return m.getTypeSimple(key[0])
}

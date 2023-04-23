package djson

import (
	"log"
	"testing"
)

func TestParseJson(t *testing.T) {
	jsonDoc := `{
		"name": null,
		"age": 9223372036854775807,
		"address" : [
			"Seoul", "Korea"
		],
		"family" : {
			"father": "John",
			"mother": "Jane"
		}
	}`

	obj, err := ParseToObject(jsonDoc)
	if err == nil {
		log.Println(obj.String("name"))
		log.Println(obj.String("age"))
		arr, ok := obj.Array("address")
		if !ok {
			log.Fatal("no such key")
		}
		log.Println(arr.String(1))
	}
}

func TestParseArray(t *testing.T) {

	jsonDoc := `[
		{
			"name":"Ricardo Longa",
			"idade":28,
			"skills":[
				"Golang","Android"
			]
		},
		{
			"name":"Hery Victor",
			"idade":32,
			"skills":[
				"Golang",
				"Java"
			]
		}
	]`

	arr, err := ParseToArray(jsonDoc)
	if err == nil {
		obj, ok := arr.Object(0)
		if !ok {
			log.Fatal("no such index")
		}
		log.Println(obj.String("name"))
	} else {
		log.Fatal("not array")
	}

}

func TestParseDJSON(t *testing.T) {
	jsonDoc := `[
		{
			"name":"Ricardo Longa",
			"idade":28,
			"skills":[
				"Golang","Android"
			]
		},
		{
			"name":"Hery Victor",
			"idade":32,
			"skills":[
				"Golang",
				"Java"
			]
		}
	]`

	aJson := New().Parse(jsonDoc)

	bJson, ok := aJson.Object(1)
	if !ok {
		log.Fatal("not object")
	}

	log.Println(bJson.Int("skills"))
	log.Println(bJson.String())
}

func TestPutDJSON(t *testing.T) {
	aJson := New().Put(
		Array{
			Object{
				"name":  "Ricardo Longa",
				"idade": 28,
				"skills": Array{
					"Golang",
					"Android",
				},
			},
			Object{
				"name":  "Hery Victor",
				"idade": 32,
				"skills": Array{
					"Golang",
					"Java",
				},
			},
		},
	)

	bJson, ok := aJson.Object(1)
	if !ok {
		log.Fatal("not object")
	}

	log.Println(bJson.Int("skills"))
	log.Println(bJson.String())

	log.Println(bJson.HasKey("name"))
}

func TestPutDArrayDJSON(t *testing.T) {
	aJson := New()
	aJson.Put(
		Array{
			Array{
				1, 2, 3, 4,
			},
			Array{
				5, 6, 7, 8,
			},
		},
	)

	bJson, ok := aJson.Get(1)
	if !ok {
		log.Fatal("not array")
	}

	log.Println(bJson.Int(1))
	log.Println(bJson.String())

	log.Println(bJson.HasKey(1))
}

func TestArrayAppendDJSON(t *testing.T) {
	aJson := New()
	aJson.Put(
		Array{
			1, 2, 3, 4,
		},
	)
	aJson.Put( // append array
		Array{
			5, 6, 7, 8,
		},
	)

	log.Println(aJson.String())
}

func TestObjectAppendDJSON(t *testing.T) {
	aJson := New()
	aJson.Put(
		Object{
			"name": "Hery Victor",
		},
	)
	aJson.Put( // append
		Object{
			"idade": 32,
		},
	)

	aJson.Put("not appended")

	log.Println(aJson.String())
}

func TestUpdateDJSON(t *testing.T) {
	jsonDoc := `[
		{
			"name":"Ricardo Longa",
			"idade":28,
			"skills":[
				"Golang","Android"
			]
		},
		{
			"name":"Hery Victor",
			"idade":32,
			"skills":[
				"Golang",
				"Java"
			]
		}
	]`

	aJson := New().Parse(jsonDoc)

	bJson, ok := aJson.Object(1)
	if !ok {
		log.Fatal("not object")
	}
	bJson.Put(Object{"hobbies": Array{"game"}})
	_ = bJson.UpdatePath(`["hobbies"][1]`, "running")
	_ = bJson.UpdatePath(`["hobbies"][0]`, "art")

	log.Println(aJson.String())
}

func TestHandleDJSON(t *testing.T) {
	jsonDoc := `{
		"name":"Ricardo Longa",
		"idade":28,
		"skills":[
			"Golang","Android"
		]
	}`

	mJson := New().Parse(jsonDoc)

	aJson := New()
	// aJson.Put("name", mJson.GetAsInterface("name"))
	// aJson.Put("idade", mJson.GetAsInterface("idade"))

	aJson.PutObject("name", mJson.Interface("name"))
	aJson.PutObject("idade", mJson.Interface("idade"))

	log.Println(aJson.ToString())

}

func TestFastDeclare(t *testing.T) {
	dJson := NewObject(
		"key", "value", "key2", "value2",
	)

	log.Println(dJson.ToString())

	aJson := NewArray(
		1, 2, 3, 4, 5, 6, 7,
	)

	log.Println(aJson.ToString())

	tJson := New(ARRAY)
	tJson.Put(1)

	log.Println(tJson.ToString())

	pJson := New()
	pJson.Put(1)

	log.Println(pJson.ToString())
}

func TestWrapArray(t *testing.T) {
	mJson := New(ARRAY)

	bb := []int{1, 2, 3, 4, 5, 6, 7, 8}

	mJson.Put(bb)

	log.Println(mJson.ToString())
}

func TestSeekNext(t *testing.T) {

	jsonDoc := `[
		{"name": "Ricardo Longa"},
		{"name": "Hery Victor"},
		3,4
	]`

	aJson := New().Parse(jsonDoc)

	bJson := aJson.Scan()
	if bJson == nil {
		log.Fatal("not object")
	}

	log.Println(bJson.ToString())

	bJson = aJson.Scan()
	if bJson == nil {
		log.Fatal("not object")
	}

	log.Println(bJson.ToString())

	bJson = aJson.Scan()
	if bJson == nil {
		log.Fatal("not object")
	}

	log.Println(bJson.ToString())
}

func TestFloat(t *testing.T) {

	jsonDoc := `{"number": "NaN"}`

	aJson := New().Parse(jsonDoc)

	log.Println(aJson.ToString())

	bJson := New(OBJECT)
	bJson.Put("number", aJson.Float("number"))
	bJson.Put("number2", 1)

	log.Println(bJson.ToString())
	log.Println("2")

}

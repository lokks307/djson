package djson

import (
	"log"
	"testing"
)

func TestGetType(t *testing.T) {

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

	mJson := New().Parse(jsonDoc)

	log.Println(mJson.Type())
	log.Println(mJson.Type(1))
	log.Println(mJson.IsObject(1))
	log.Println(mJson.TypePath(`[1]["name"]`))
	log.Println(mJson.TypePath(`[1]["idade"]`))
	log.Println(mJson.TypePath(`[1]["skills"]`))

	log.Println(mJson.Float(1, 0.7))
}

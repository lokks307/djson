package djson

import (
	"log"
	"testing"
)

func TestTokenizer(t *testing.T) {

	log.Println(PathTokenizer(`["aa"][1][b_b]`))  // [aa 1 b_b]
	log.Println(PathTokenizer(`["a'a"][1][b]b]`)) // [a'a 1 b]
}

func TestParse(t *testing.T) {
	doc := `[[1,2,3]]`
	tdjson := New().Parse(doc)
	log.Println(tdjson.IntPath(`[0][0]`))
}

func TestPremitiveArray(t *testing.T) {
	tdjson := New().Put(Object{"array": []string{"1", "2", "3"}})
	log.Println(tdjson.ToString())

	xdjson := NewArray().PutArray(Object{"array": []string{"1", "2", "3"}})
	log.Println(xdjson.ToString())

	adjson := New().Put([]string{"1", "2", "3"})
	log.Println(adjson.ToString())

	bdjson := NewArray().PutArray([]string{"1", "2", "3"}, []string{"4", "5", "6"}, "7", "8", "9")
	log.Println(bdjson.ToString())

	cdjson := NewObject().Put([]string{"1", "2", "3"})
	log.Println(cdjson.ToString())
}

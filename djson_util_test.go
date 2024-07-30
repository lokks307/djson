package djson

import (
	"fmt"
	"log"
	"testing"

	"github.com/volatiletech/null/v8"
)

func TestToFieldTag(t *testing.T) {
	type User struct {
		Id    string      `json:"id"`
		Name  string      `json:"name"`
		Email null.String `json:"email"`
	}

	var user User

	mJson := New().Put(
		Object{
			"id":    "id-1234",
			"name":  "Ricardo Longa",
			"email": "longa@test.com",
		},
	)

	mJson.ToFields(&user, "id", "email")

	log.Println(user)
}

func TestFromFieldTag(t *testing.T) {
	type User struct {
		Id    string      `json:"id"`
		Name  string      `json:"name"`
		Email null.String `json:"email"`
	}

	var user = User{
		Id:   "id-1234",
		Name: "Ricardo Longa",
		Email: null.String{
			String: "longa@test.com",
			Valid:  true,
		},
	}

	mJson := New()
	mJson.FromFields(user)

	log.Println(mJson.ToString())
}

func TestFromFieldMapTest(t *testing.T) {

	type Name struct {
		First  string `json:"first"`
		Family string `json:"family"`
	}

	user := make(map[string]interface{})

	user["id"] = "id-1234"
	user["name"] = Name{
		First:  "Ricardo",
		Family: "Longa",
	}

	user["email"] = null.String{
		String: "longa@test.com",
		Valid:  true,
	}

	mJson := New()
	mJson.FromFields(user, "name.first", "email")

	log.Println(mJson.ToString())
}

func TestSortingArray(t *testing.T) {
	mJson := NewArray(5, 6, 7, 8, 1, 2, 3, 4)
	if ok := mJson.SortAsc(); !ok {
		log.Fatal("sorting test failed")
	}

	log.Println(mJson.ToString())

	tJson := New().Put(Object{
		"d": "aaa",
		"a": Array{
			5, 6, 7, 8, 1, 2, 3, 4,
		},
	})

	if ok := tJson.SortAsc("a"); !ok {
		log.Fatal("sorting test failed")
	}

	log.Println(tJson.ToString())

	ok := tJson.SortDescPath(`["a"]`)

	if !ok {
		log.Fatal("sorting path failed")
	}

	log.Println(tJson.ToString())

	pJson := New().Put(
		Array{
			Object{
				"k": "1",
				"v": "1",
			},
			Object{
				"k": "22",
				"v": "2",
			},
			Object{
				"k": "4444",
				"v": "4",
			},
			Object{
				"k": "333",
				"v": "3",
			},
		},
	)

	pJson.SortArrayAsc("k")

	p2Json := New().Put(
		Object{
			"kk": Array{
				Object{
					"k": null.String{
						String: "1",
					},
					"v": "1",
				},
				Object{
					"k": null.String{
						String: "22",
					},
					"v": "2",
				},
				Object{
					"k": "9",
					"v": "4",
				},
				Object{
					"k": "1",
					"v": "3",
				},
			},
		},
	)

	p2Json.SortDescPath(`[kk]`, "k")

	log.Println(pJson.ToString())
	log.Println(p2Json.ToString())

}

func TestReflectType(t *testing.T) {

	aJson := New().Put(Object{
		"name": "yu",
		"skill": Array{
			"running", "playing",
		},
	})

	bJson := New().Put(Object{
		"skill": Array{
			"running", "playing",
		},
		"name": "yu",
	})

	if aJson.Equal(bJson) {
		log.Println("jsons are the same")
	} else {
		log.Println("jsons are not the same")
	}
}

func TestClone(t *testing.T) {
	aJson := New().Put(Object{
		"name": "yu",
		"skill": Array{
			"running", "playing",
		},
	})

	bJson := aJson.Clone()

	aJson.UpdatePath(`[skill][2]`, "swimming")

	bJson.PutObject("name", "not you")
	bJson.UpdatePath(`[skill][2]`, "studying")

	log.Println(aJson.ToString())
	log.Println(bJson.ToString())
}

func TestFind(t *testing.T) {
	aJson := New().Put(Array{
		Object{
			"name":  "1",
			"skill": "apple",
		},
		Object{
			"name":  "2",
			"skill": "banana",
		},
	})

	log.Println(aJson.ToString())

	bJson := aJson.Find("name", "1")
	if bJson == nil {
		log.Fatalln("find failed")
	}

	log.Println(bJson.ToString())
}

func TestAppend(t *testing.T) {
	aJson := New().Put(Array{
		Object{
			"name":  "1",
			"skill": "apple",
		},
		Object{
			"name":  "2",
			"skill": "banana",
		},
	})

	bJson := New().Put(Array{
		Object{
			"name":  "3",
			"skill": "apple",
		},
		Object{
			"name":  "4",
			"skill": "banana",
		},
	})

	log.Println(aJson.Append(bJson).ToString())
}

func TestCloneNil(t *testing.T) {
	aJson := New().Put(Array{
		Object{
			"name":  nil,
			"skill": "apple",
		},
		Object{
			"name":  "2",
			"skill": "banana",
		},
	})

	bJson := aJson.Clone()

	log.Println(bJson.ToString())
}

func TestPutArray(t *testing.T) {
	mJson := New()
	intArray := []int{1, 2, 3, 4, 5}
	intVal := 3
	strArray := []string{"yes", "this", "is right"}
	strVal := "abc"

	mJson.Put("int_array", NewArray(intArray))
	mJson.Put("int_val", intVal)
	mJson.Put("str_array", NewArray(strArray))
	mJson.Put("str_val", strVal)

	fmt.Println(mJson.ToString())
	// must be like {"int_array":[1,2,3,4,5],"int_val":3,"str_array":["yes","this","is right"],"str_val":"abc"}

	expected := `{"int_array":[1,2,3,4,5],"int_val":3,"str_array":["yes","this","is right"],"str_val":"abc"}`
	if result := mJson.ToString(); result != expected {
		t.Errorf("Expected %s, but got %s", expected, result)
	}
}

func TestSyntax12A(t *testing.T) {
	mJson := New()
	mJson.Put(Array{
		1, 2, 3, 4, 5, 6,
	})

	fmt.Println(mJson.ToString()) // must be [1,2,3,4,5,6]

	expected := "[1,2,3,4,5,6]"
	if result := mJson.ToString(); result != expected {
		t.Errorf("Expected %s, but got %s", expected, result)
	}
}

func TestSyntax12B(t *testing.T) {
	mJson := New()
	mJson.Put(Object{
		"name":  "Hery Victor",
		"idade": 32,
	})

	fmt.Println(mJson.ToString()) // must be {"idade":32,"name":"Hery Victor"}

	expected := `{"idade":32,"name":"Hery Victor"}`
	if result := mJson.ToString(); result != expected {
		t.Errorf("Expected %s, but got %s", expected, result)
	}
}

func TestSyntax13(t *testing.T) {
	mJson := New()
	mJson.Put(
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
	fmt.Println(mJson.ToString()) // must be like [{"idade":28,"name":"Ricardo Longa","skills":["Golang","Android"]},{"idade":32,"name":"Hery Victor","skills":["Golang","Java"]}]

	expected := `[{"idade":28,"name":"Ricardo Longa","skills":["Golang","Android"]},{"idade":32,"name":"Hery Victor","skills":["Golang","Java"]}]`
	if result := mJson.ToString(); result != expected {
		t.Errorf("Expected %s, but got %s", expected, result)
	}
}

func TestSyntax14A(t *testing.T) {
	mJson := New()

	mJson.Put(Array{
		1, 2, 3, 4, 5, 6,
	})

	mJson.Put(Array{
		7, 8, 9,
	})

	fmt.Println(mJson.ToString()) // must be [1,2,3,4,5,6,7,8,9]

	expected := "[1,2,3,4,5,6,7,8,9]" // 예상 출력값 설정
	if result := mJson.ToString(); result != expected {
		t.Errorf("Expected %s, but got %s", expected, result)
	}
}

func TestSyntax14B(t *testing.T) {
	// Object
	mJson := New()

	mJson.Put(Object{
		"name": "Hery Victor",
	})

	mJson.Put(Object{
		"idade": 28,
	})

	mJson.Put(Object{
		"name": "Ricardo Longa", // overwrite existing value
	})

	fmt.Println(mJson.ToString()) // must be like {"idade":28,"name":"Ricardo Longa"}

	expected := `{"idade":28,"name":"Ricardo Longa"}` // 예상 출력값 설정
	if result := mJson.ToString(); result != expected {
		t.Errorf("Expected %s, but got %s", expected, result)
	}
}

func TestSyntax15(t *testing.T) {
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
	fmt.Println(mJson.ToString()) // must be like [{"idade":28,"name":"Ricardo Longa","skills":["Golang","Android"]},{"idade":32,"name":"Hery Victor","skills":["Golang","Java"]}]

	expected := `[{"idade":28,"name":"Ricardo Longa","skills":["Golang","Android"]},{"idade":32,"name":"Hery Victor","skills":["Golang","Java"]}]`
	if result := mJson.ToString(); result != expected {
		t.Errorf("Expected %s, but got %s", expected, result)
	}
}

func TestSyntax16Simple(t *testing.T) {
	jsonDoc := `[
    {
        "name":"Ricardo Longa",
        "idade":28,
        "skills":[
            "Golang","Android"
        ]
    }
]`

	mJson := New().Parse(jsonDoc)

	// must be like [{"name":"Ricardo Longa","idade":28,"skills":["Golang","Android"]}]
	fmt.Println(mJson.ToString())

	// must be like {"name":"Ricardo Longa","idade":28,"skills":["Golang","Android"]}
	fmt.Println(mJson.String(0))

	aJson, _ := mJson.Object(0)

	fmt.Println(aJson.Int("idade"))    // 28
	fmt.Println(aJson.String("idade")) // 28

	fmt.Println(aJson.Int("name"))    // 0
	fmt.Println(aJson.String("name")) // Ricardo Longa

	fmt.Println(aJson.String("skills")) // ["Golang","Android"]
}

func TestSyntax16(t *testing.T) {
	jsonDoc := `[{
		"name":"Ricardo Longa",
		"idade":28,
		"skills":[
			"Golang","Android"
		]
	}]`

	mJson := New().Parse(jsonDoc)

	expectedJsonString := `[{"idade":28,"name":"Ricardo Longa","skills":["Golang","Android"]}]`
	expectedStringAtIndex0 := `{"idade":28,"name":"Ricardo Longa","skills":["Golang","Android"]}`

	if result := mJson.ToString(); result != expectedJsonString {
		t.Errorf("Expected %s, but got %s", expectedJsonString, result)
	}

	if result := mJson.String(0); result != expectedStringAtIndex0 {
		t.Errorf("Expected %s, but got %s", expectedStringAtIndex0, result)
	}

	aJson, _ := mJson.Object(0)

	if result := aJson.Int("idade"); result != 28 {
		t.Errorf("Expected 28, but got %d", result)
	}

	if result := aJson.String("idade"); result != "28" {
		t.Errorf("Expected '28', but got '%s'", result)
	}

	if result := aJson.Int("name"); result != 0 {
		t.Errorf("Expected 0, but got %d", result)
	}

	if result := aJson.String("name"); result != "Ricardo Longa" {
		t.Errorf("Expected 'Ricardo Longa', but got '%s'", result)
	}

	if result := aJson.String("skills"); result != `["Golang","Android"]` {
		t.Errorf(`Expected '["Golang","Android"]', but got '%s'`, result)
	}
}

func TestSyntax21(t *testing.T) {
	jsonDoc := `{
		"name":"Ricardo Longa",
		"idade":28,
		"skills":[
			"Golang","Android"
		]
	}`

	mJson := New().Parse(jsonDoc)
	aJson, _ := mJson.Array("skills")

	fmt.Println(mJson.HasKey("skills")) // true
	fmt.Println(aJson.HasKey(1))        // true
	fmt.Println(aJson.HasKey("addr"))   // false

	expectedHasKeySkill := true
	if result := mJson.HasKey("skills"); result != expectedHasKeySkill {
		t.Errorf("Expected %v, but got %v", expectedHasKeySkill, result)
	}

	expectedHasKey1 := true
	if result := aJson.HasKey(1); result != expectedHasKey1 {
		t.Errorf("Expected %v, but got %v", expectedHasKey1, result)
	}

	expectedHasKeyAddr := false
	if result := aJson.HasKey("addr"); result != expectedHasKeyAddr {
		t.Errorf("Expected %v, but got %v", expectedHasKeyAddr, result)
	}
}

func TestSyntax22(t *testing.T) {
	jsonDoc := `[{
		"name":"Ricardo Longa",
		"idade":28,
		"skills":["Golang","Android"]
	},
	{
		"name":"Hery Victor",
		"idade":32,
		"skills":["Golang", "Java"]
	}]`

	mJson := New().Parse(jsonDoc)

	_ = mJson.UpdatePath(`[1]["name"]`, Object{
		"family": "Victor",
		"first":  "Hery",
	})

	// must be like
	// [{"idade":28,"name":"Ricardo Longa","skills":["Golang","Android"]},
	// {"idade":32,"name":{"family":"Victor","first":"Hery"},"skills":["Golang","Java"]}]

	expected := `[{"idade":28,"name":"Ricardo Longa","skills":["Golang","Android"]},{"idade":32,"name":{"family":"Victor","first":"Hery"},"skills":["Golang","Java"]}]`

	if result := mJson.ToString(); result != expected {
		t.Errorf("Expected %s, but got %s", expected, result)
	}
}

func TestSyntax23(t *testing.T) {
	jsonDoc := `[{
		"name":"Ricardo Longa",
		"idade":28,
		"skills":["Golang","Android"]
	},
	{
		"name":"Hery Victor",
		"idade":32,
		"skills":["Golang", "Java"]
	}]`

	mJson := New().Parse(jsonDoc)

	_ = mJson.RemovePath(`[1]["skills"]`)

	expected := `[{"idade":28,"name":"Ricardo Longa","skills":["Golang","Android"]},{"idade":32,"name":"Hery Victor"}]`

	if result := mJson.ToString(); result != expected {
		t.Errorf("Expected %s, but got %s", expected, result)
	}
}

func TestSyntax24A(t *testing.T) {
	jsonDoc := `[{
		"name":"Ricardo Longa",
		"idade":28,
		"skills":["Golang","Android"]
	},
	{
		"name":"Hery Victor",
		"idade":32,
		"skills":["Golang","Java"]
	}]`

	aJson := New().Parse(jsonDoc)

	bJson, _ := aJson.Object(1) // bJson이 aJson과 같은 *djson.DO를 공유

	bJson.Put(Object{"hobbies": Array{"game"}})   // 객체에 배열 추가
	bJson.UpdatePath(`["hobbies"][1]`, "running") // 배열에 값 추가
	bJson.UpdatePath(`["hobbies"][0]`, "art")     // 배열 값 변경

	expected := `[{"idade":28,"name":"Ricardo Longa","skills":["Golang","Android"]},{"hobbies":["art"],"idade":32,"name":"Hery Victor","skills":["Golang","Java"]}]`

	if result := aJson.ToString(); result != expected {
		t.Errorf("Expected %s, but got %s", expected, result)
	}
}

func TestSyntax24B(t *testing.T) {
	jsonDoc := `{
		"name":"Ricardo Longa",
		"idade":28,
		"skills":["Golang","Android"]
	}`

	mJson := New().Parse(jsonDoc)

	fmt.Println(mJson.Int("idade", 10))   // must be 28
	fmt.Println(mJson.Int("name", 10))    // must be 10 because `name` field cannot be convert to integer
	fmt.Println(mJson.Int("hobbies", 10)) // must be 10 because no such field `hobbies`

	expectedIntIdade := int64(28)
	if result := mJson.Int("idade", 10); result != expectedIntIdade {
		t.Errorf("Expected %d, but got %d", expectedIntIdade, result)
	}

	expectedIntName := int64(10)
	if result := mJson.Int("name", 10); result != expectedIntName {
		t.Errorf("Expected %d, but got %d", expectedIntName, result)
	}

	expectedIntHobbies := int64(10)
	if result := mJson.Int("hobbies", 10); result != expectedIntHobbies {
		t.Errorf("Expected %d, but got %d", expectedIntHobbies, result)
	}
}

// - DJSON supports null package for sqlboiler (`github.com/volatiletech/null`)
func TestSyntax25(t *testing.T) {
	type User struct {
		Id    string      `json:"id"`
		Name  string      `json:"name"`
		Email null.String `json:"email"`
	}

	var user User

	mJson := New().Put(
		Object{
			"id":    "id-1234",
			"name":  "Ricardo Longa",
			"email": "longa@test.com",
		},
	)

	mJson.ToFields(&user, "id", "email") // only (`id`, `email`) tag

	fmt.Println(user) // must be {id-1234  {longa@test.com true}}

	// 추가 테스트 케이스
	t.Run("TestNameField", func(t *testing.T) {
		var user2 User
		// name 필드만 추출
		mJson.ToFields(&user2, "name")
		expectedUser := User{
			Id:    "",
			Name:  "Ricardo Longa",
			Email: null.String{},
		}

		// 예상 결과와 비교
		if user2 != expectedUser {
			t.Errorf("Expected %v, but got %v", expectedUser, user2)
		}
	})
}

func TestSyntax26(t *testing.T) {
	type User struct {
		Id      string      `json:"id"`
		Name    string      `json:"name"`
		Email   null.String `json:"email"`
		Address null.String `json:"address"`
	}

	var user = User{
		Id:   "id-1234",
		Name: "Ricardo Longa",
		Email: null.String{
			String: "longa@test.com",
			Valid:  true,
		},
		Address: null.String{
			String: "Area51",
			Valid:  false,
		},
	}

	mJson := New()
	mJson.FromFields(user) // no tag

	fmt.Println(mJson.ToString())

	expectedJSON := `{"address":"","email":"longa@test.com","id":"id-1234","name":"Ricardo Longa"}`
	if result := mJson.ToString(); result != expectedJSON {
		t.Errorf("Expected %s, but got %s", expectedJSON, result)
	}
}

func TestSyntax27(t *testing.T) {
	type Name struct {
		First  string `json:"first"`
		Family string `json:"family"`
	}

	user := make(map[string]interface{})

	user["id"] = "id-1234"
	user["name"] = Name{
		First:  "Ricardo",
		Family: "Longa",
	}

	user["email"] = null.String{
		String: "longa@test.com",
		Valid:  true,
	}

	mJson := New()
	mJson.FromFields(user, "name.first", "email") // tag has depth concept

	// must be like {"email":"longa@test.com","name":{"first":"Ricardo"}}
	fmt.Println(mJson.ToString())

	expectedJSON := `{"email":"longa@test.com","name":{"first":"Ricardo"}}`
	if result := mJson.ToString(); result != expectedJSON {
		t.Errorf("Expected %s, but got %s", expectedJSON, result)
	}
}

func TestOmitEmpty(t *testing.T) {
	type Name struct {
		First  string      `json:"first,omitempty"`
		Middle null.String `json:"middle,omitempty"`
		Family string      `json:"family"`
	}

	type User struct {
		Name    interface{} `json:"name,omitempty"`
		Address string      `json:"address"`
	}

	userName := Name{
		First:  "Ricardo",
		Family: "Longa",
		Middle: null.String{
			String: "Legend",
			Valid:  false,
		},
	}

	//mJson := New()
	//mJson.FromFields(userName) // tag has depth concept

	// must be like {"email":"longa@test.com","name":{"first":"Ricardo"}}
	//fmt.Println(mJson.ToString())

	user := User{
		Name:    userName,
		Address: "South Korea",
	}

	bJson := New()
	bJson.FromFields(user) // tag has depth concept

	fmt.Println(bJson.ToString())

}

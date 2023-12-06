# DJSON: Another JSON Library for who hates json.Unmarshal

## DJSON v2
- Compared to v1, the names of the functions are simple in v2.
- Various useful functions for handling djson structure are included in DJSON v2.
- github.com/lokks307/go-util will not serve DJSON v2 any more.

## 1. Basic Syntax for v1

### 1.1. New JSON

```go
mJson := djson.New()
```

### 1.2. Assign values to JSON

```go
// Array
mJson := djson.New()
mJson.Put(djson.Array{
    1,2,3,4,5,6,
})

fmt.Println(mJson.ToString()) // must be [1,2,3,4,5,6]
```

```go
// Object
mJson := djson.New()
mJson.Put(djson.Object{
    "name":  "Hery Victor",
    "idade": 32,
})

fmt.Println(mJson.ToString()) // must be {"idade":32,"name":"Hery Victor"}
```

### 1.2. Assign Arrays to JSON
```go
mJson := djson.New()
intArray := []int{1, 2, 3, 4, 5}
intVal := 3
strArray := []string{"yes", "this", "is right"}
strVal := "abc"

mJson.Put("int_array", djson.NewArray(intArray))
mJson.Put("int_val", intVal)
mJson.Put("str_array", djson.NewArray(strArray))
mJson.Put("str_val", strVal)

fmt.Println(mJson.ToString())
// must be like {"int_array":[1,2,3,4,5],"int_val":3,"str_array":["yes","this","is right"],"str_val":"abc"}
```


### 1.3. Assign table values to JSON
```go
mJson := djson.New()
mJson.Put(
    djson.Array{
        djson.Object{
            "name":  "Ricardo Longa",
            "idade": 28,
            "skills": djson.Array{
                "Golang",
                "Android",
            },
        },
        djson.Object{
            "name":  "Hery Victor",
            "idade": 32,
            "skills": djson.Array{
                "Golang",
                "Java",
            },
        },
    },
)

fmt.Println(mJson.ToString()) // must be like [{"idade":28,"name":"Ricardo Longa","skills":["Golang","Android"]},{"idade":32,"name":"Hery Victor","skills":["Golang","Java"]}]
```

### 1.4. Append values to existing JSON

```go
// Array
mJson := djson.New()

mJson.Put(djson.Array{
    1,2,3,4,5,6,
})

mJson.Put(djson.Array{
    7,8,9,
})

fmt.Println(mJson.ToString()) // must be [1,2,3,4,5,6,7,8,9]
```

```go
// Object
mJson := djson.New()

mJson.Put(djson.Object{
        "idade": 28,
    })

mJson.Put(djson.Object{
    "name":"Hery Victor",
})

mJson.Put(djson.Object{
    "name":"Ricardo Longa", // overwrite existing value
})

fmt.Println(mJson.ToString()) // must be like {"idade":28,"name":"Ricardo Longa"}
```

### 1.5. Parse existing JSON string

```go
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

mJson := djson.New().Parse(jsonDoc)

fmt.Println(mJson.ToString()) // must be like [{"idade":28,"name":"Ricardo Longa","skills":["Golang","Android"]},{"idade":32,"name":"Hery Victor","skills":["Golang","Java"]}]
```

### 1.6. Get values
```go
jsonDoc := `[
    {
        "name":"Ricardo Longa",
        "idade":28,
        "skills":[
            "Golang","Android"
        ]
    }
]`

mJson := djson.New().Parse(jsonDoc)

// must be like [{"idade":28,"name":"Ricardo Longa","skills":["Golang","Android"]}]
fmt.Println(mJson.ToString()) 

// must be like {"idade":28,"name":"Ricardo Longa","skills":["Golang","Android"]}
fmt.Println(mJson.String(0)) 

aJson, _ := mJson.Object(0)

fmt.Println(aJson.Int("idade")) // 28
fmt.Println(aJson.String("idade")) // 28

fmt.Println(aJson.Int("name")) // 0
fmt.Println(aJson.String("name")) // Ricardo Longa

fmt.Println(aJson.String("skills")) // ["Golang","Android"]
```

## 2. Advanced Syntax

### 2.1. Check if JSON has key

```go
jsonDoc := `{
    "name":"Ricardo Longa",
    "idade":28,
    "skills":[
        "Golang","Android"
    ]
}`

mJson := djson.New().Parse(jsonDoc)
aJson, _ := mJson.Array("skills")

fmt.Println(mJson.HasKey("skills")) // true
fmt.Println(aJson.Haskey(1)) // true
fmt.Println(aJson.HasKey("addr"))   // false

```

### 2.2. Update value via path

```go
jsonDoc := `[{
    "name":"Ricardo Longa",
    "idade":28,
    "skills":[ "Golang","Android" ]
},
{
    "name":"Hery Victor",
    "idade":32,
    "skills":[ "Golang", "Java" ]
}]`

mJson := djson.New().Parse(jsonDoc)

_ = mJson.UpdatePath(`[1]["name"]`, djson.Object{
    "family": "Victor",
    "first":  "Hery",
})

// must be like
// [{"idade":28,"name":"Ricardo Longa","skills":["Golang","Android"]},
// {"idade":32,"name":{"family":"Victor","first":"Hery"},"skills":["Golang","Java"]}]
fmt.Println(mJson.ToString()) 
```

### 2.3. Remove value via path

```go
jsonDoc := `[{
    "name":"Ricardo Longa",
    "idade":28,
    "skills":[ "Golang","Android" ]
},
{
    "name":"Hery Victor",
    "idade":32,
    "skills":[ "Golang", "Java" ]
}]`

mJson := djson.New().Parse(jsonDoc)

_ = mJson.RemovePath(`[1]["skills"]`)

// must be like
// [{"idade":28,"name":"Ricardo Longa","skills":["Golang","Android"]},{"idade":32,"name":"Hery Victor"}]
fmt.Println(mJson.ToString()) 
```

### 2.4. Manipluate value via sharing

```go
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

aJson := djson.New().Parse(jsonDoc)

bJson, _ := aJson.Object(1) // now, bJson shares *djson.DO with aJson

bJson.Put(djson.Object{"hobbies": djson.Array{"game"}}) // append Array to Object
bJson.UpdatePath(`["hobbies"][1]`, "running") // append value to Array
bJson.UpdatePath(`["hobbies"][0]`, "art") // replace value

// must be like
// [{"idade":28,"name":"Ricardo Longa","skills":["Golang","Android"]},
// {"hobbies":["art","running"],"idade":32,"name":"Hery Victor","skills":["Golang","Java"]}]
fmt.Println(aJson.ToString())
```


### 2.4. Default Value
```go
jsonDoc := `{
    "name":"Ricardo Longa",
    "idade":28,
    "skills":["Golang","Android"]
}`

mJson := djson.New().Parse(jsonDoc)

fmt.Println(mJson.Int("idade", 10)) // must be 28 which is int64
fmt.Println(mJson.Int("name", 10)) // must be 10 because `name` field cannot be convert to integer
fmt.Println(mJson.Int("hobbies", 10)) // must be 10 because no such field `hobbies`
```

### 2.5. To Structure
- DJSON supports null package for sqlboiler (`github.com/volatiletech/null`)
```go
type User struct {
    Id    string      `json:"id"`
    Name  string      `json:"name"`
    Email null.String `json:"email"`
}

var user User

mJson := djson.New().Put(
    Object{
        "id":    "id-1234",
        "name":  "Ricardo Longa",
        "email": "longa@test.com",
    },
)

mJson.ToFields(&user, "id", "email") // only (`id`, `email`) tag

fmt.Println(user) // must be {id-1234  {longa@test.com true}}
```

### 2.6. From Structure
```go
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

mJson := djson.New()
mJson.FromFields(user) // no tag

// must be like {"email":"longa@test.com","id":"id-1234","name":"Ricardo Longa"}
fmt.Println(mJson.ToString()) 

```

### 2.7. From Map and Structure
```go
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

mJson := djson.New()
mJson.FromFields(user, "name.first", "email") // tag has depth concept

// must be like {"email":"longa@test.com","name":{"first":"Ricardo"}}
fmt.Println(mJson.ToString())
```
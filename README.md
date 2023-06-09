# DJSON: Another JSON Library for who hates json.Unmarshal

## DJSON v2
- Compared to v1, the names of the functions are simple in v2.
- Various useful functions for handling djson structure are included in DJSON v2.
- github.com/lokks307/go-util will not serve DJSON v2 any more.

## 1. Basic Syntax for v1

### 1.1. New JSON

```go
mJson := djson.NewDJSON()
```

### 1.2. Assign values to JSON

```go
// Array
mJson := djson.NewDJSON()
mJson.Put(djson.Array{
    1,2,3,4,5,6
})
```

```go
// Object
mJson := djson.NewDJSON()
mJson.Put(djson.Object{
    "name":  "Hery Victor",
    "idade": 32,
})
```

### 1.3. Assign table values to JSON
```go
mJson := djson.NewDJSON()
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
```

### 1.4. Append values to existing JSON

```go
// Array
mJson := djson.NewDJSON()

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
mJson := djson.NewDJSON()

mJson.Put(djson.Object{
    "name":"Hery Victor",
})

mJson.Put(djson.Object{
    "idade": 28,
})

mJson.Put(djson.Object{
    "name":"Ricardo Longa", // overwrite existing value
})

fmt.Println(mJson.ToString()) // must be like {"name":"Ricardo Longa","idade":28}
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

mJson := NewDJSON().Parse(jsonDoc)
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

mJson := NewDJSON().Parse(jsonDoc)

// must be like [{"name":"Ricardo Longa","idade":28,"skills":["Golang","Android"]}]
fmt.Println(mJson.ToString()) 

// must be like {"name":"Ricardo Longa","idade":28,"skills":["Golang","Android"]}
fmt.Println(mJson.GetAsString(0)) 

aJson, _ := mJson.GetAsObject(0)

fmt.Println(aJson.GetAsInt("idade")) // 28
fmt.Println(aJson.GetAsString("idade")) // 28

fmt.Println(aJson.GetAsInt("name")) // 0
fmt.Println(aJson.GetAsString("name")) // Ricardo Longa

fmt.Println(aJson.GetAsString("skills")) // ["Golang","Android"]
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

mJson := NewDJSON().Parse(jsonDoc)
aJson, _ := mJson.GetAsArray("skills")

fmt.Println(mJson.HasKey("skill")) // true
fmt.Println(aJson.Haskey(1)) // true

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

mJson := NewDJSON().Parse(jsonDoc)

_ = mJson.UpdatePath(`[1]["name"]`, djson.Object{
    "first":  "Hery",
    "family": "Victor",
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

mJson := NewDJSON().Parse(jsonDoc)

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

aJson := NewDJSON().Parse(jsonDoc)

bJson, _ := aJson.GetAsObject(1) // now, bJson shares *djson.DO with aJson

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

mJson := NewDJSON().Parse(jsonDoc)

fmt.Println(mJson.GetAsInt("idade", 10)) // must be 28
fmt.Println(mJson.GetAsInt("name", 10)) // must be 10 because `name` field cannot be convert to integer
fmt.Println(mJson.GetAsInt("hobbies", 10)) // must be 10 because no such field `hobbies`
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

mJson := NewDJSON().Put(
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

mJson := NewDJSON()
mJson.FromFields(user) // no tag

// must be like {"email":"longa@test.com","id":"id-1234","name":"Ricardo Longa"}
fmt.Println(mJson.ToString()) 

}
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

mJson := NewDJSON()
mJson.FromFields(user, "name.first", "email") // tag has depth concept

// must be like {"email":"longa@test.com","name":{"first":"Ricardo"}}
fmt.Println(mJson.ToString())
```
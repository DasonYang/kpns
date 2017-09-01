package mongo

import (
    "fmt"
    "testing"
    // "reflect"
    // "gopkg.in/mgo.v2"
    // "gopkg.in/mgo.v2/bson"
)

// func TestMongo(t *testing.T) {
//     type Person struct {
//         Name string
//         Phone string
//     }

//     db := New()

//     db.Database = "test"
//     db.Collection = "people"

//     // c := db.Session.DB("test").C("people")

//     // result := Person{}
//     // err := c.Find(bson.M{"name": "Ale"}).One(&result)
//     // if err != nil {
//     //     log.Fatal(err)
//     // }

//     result := db.Read()

//     log.Println("Phone:", result)
// }

// func TestCount(t *testing.T) {
//     session, err := mgo.Dial("localhost")

//     if err != nil {
//         panic(err)
//     }

//     defer session.Close()

//     session.SetMode(mgo.Monotonic, true)

//     c := session.DB("tpns").C("allow")

//     // fmt.Printf("bson = %v\n", map[string]interface{}{"value.note":map[string]interface{}{"$regex":{Pattern:"SZ141", Options:"i"}}})

//     // fmt.Printf("bson = %v\n", reflect.TypeOf(bson.RegEx{"SZ141", "i"}))

//     count, err := c.Find(bson.M{"value.note": bson.M{"$regex":bson.RegEx{"SZ141", "i"}}}).Count()
//     // count, err := c.Find(nil).Count()

//     fmt.Printf("Count = %v\n", count)
// }

func TestQuery(t *testing.T) {
    {

        fmt.Println("\n********************* Test Normal Field **********************\n")
        var query = make(map[string]interface{})

        query["key"] = "com.uncord"
        query["value.appkey"] = "AIzaSyBnI5dDd_ZxoV-4zSl033jjxQZUArzxdVo"
        query["value.lasttime"] = 1425983387.86092

        fmt.Println(makeQuery(query))
    }

    {
        fmt.Println("\n********************* Test Dict Array **********************\n")
        var query = make(map[string]interface{})
        var or_array []map[string]interface{}
        var or_query_1 = make(map[string]interface{})
        var or_query_2 = make(map[string]interface{})

        or_query_1["key"] = "com.uncord"
        or_array = append(or_array, or_query_1)
        or_query_2["value.appkey"] = "AIzaSyBnI5dDd_ZxoV-4zSl033jjxQZUArzxdVo"
        or_array = append(or_array, or_query_2)


        query["$or"] = or_array

        fmt.Println(makeQuery(query))
    }

    {
        fmt.Println("\n********************* Test Exists **********************\n")
        var query = make(map[string]interface{})


        query["$exists"] = true

        fmt.Println(makeQuery(query))
    }

    {
        fmt.Println("\n********************* Test In **********************\n")
        var query = make(map[string]interface{})


        query["$in"] = []interface{}{"First", "Second", 3, 4}

        fmt.Println(makeQuery(query))
    }

    {
        fmt.Println("\n********************* Test Regex **********************\n")
        var query = make(map[string]interface{})


        query["$regex"] = "abcde"

        fmt.Println(makeQuery(query))
    }
}
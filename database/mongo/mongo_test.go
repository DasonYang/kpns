package mongo

import (
    "fmt"
    "testing"
    // "reflect"
    "gopkg.in/mgo.v2"
    "gopkg.in/mgo.v2/bson"
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

func TestCount(t *testing.T) {
    session, err := mgo.Dial("localhost")

    if err != nil {
        panic(err)
    }

    defer session.Close()

    session.SetMode(mgo.Monotonic, true)

    c := session.DB("tpns").C("allow")

    // fmt.Printf("bson = %v\n", map[string]interface{}{"value.note":map[string]interface{}{"$regex":{Pattern:"SZ141", Options:"i"}}})

    // fmt.Printf("bson = %v\n", reflect.TypeOf(bson.RegEx{"SZ141", "i"}))

    count, err := c.Find(bson.M{"value.note": bson.M{"$regex":bson.RegEx{"SZ141", "i"}}}).Count()
    // count, err := c.Find(nil).Count()

    fmt.Printf("Count = %v\n", count)
}
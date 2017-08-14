package database

import(
    "fmt"
    "log"
    "testing"

    "gopkg.in/mgo.v2"
    "gopkg.in/mgo.v2/bson"
)

// func TestMongo(t *testing.T) {

//     type Person struct {
//         Name string
//         Phone string
//     }

//     session, err := mgo.Dial("localhost")

//     if err != nil {
//         panic(err)
//     }

//     defer session.Close()

//     session.SetMode(mgo.Monotonic, true)

//     c := session.DB("test").C("people")

//     err = c.Insert(&Person{"Ale", "+55 53 8116 9639"},
//                    &Person{"Cla", "+55 53 8402 8510"})
//     if err != nil {
//         log.Fatal(err)
//     }

//     result := Person{}
//     err = c.Find(bson.M{"name": "Ale"}).One(&result)
//     if err != nil {
//         log.Fatal(err)
//     }

//     fmt.Println("Phone:", result.Phone)
// }


package mongo

import (
    "log"
    "testing"

    // "gopkg.in/mgo.v2/bson"
)

func TestMongo(t *testing.T) {
    type Person struct {
        Name string
        Phone string
    }

    db := New()

    db.Database = "test"
    db.Collection = "people"

    // c := db.Session.DB("test").C("people")

    // result := Person{}
    // err := c.Find(bson.M{"name": "Ale"}).One(&result)
    // if err != nil {
    //     log.Fatal(err)
    // }

    result := db.Read()

    log.Println("Phone:", result)
}
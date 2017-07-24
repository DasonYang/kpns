package mongo

import (
    "log"
    // "kpns/database"

    "gopkg.in/mgo.v2"
    "gopkg.in/mgo.v2/bson"
)

type DatabaseClient struct {
    Database    string
    Collection  string
    Session     *mgo.Session
}

func New() *DatabaseClient {
    session, err := mgo.Dial("localhost")

    if err != nil {
        panic(err)
    }

    session.SetMode(mgo.Monotonic, true)

    return &DatabaseClient{Session : session}
}

func (db *DatabaseClient) Read() map[string]interface{} {
    var result map[string]interface{}

    c := db.Session.DB(db.Database).C(db.Collection)

    err := c.Find(bson.M{"name": "Ale"}).One(&result)
    if err != nil {
        log.Fatal(err)
    }

    return result
}
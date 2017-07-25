package mongo

import (
    "log"
    // "kpns/database"

    "gopkg.in/mgo.v2"
    "gopkg.in/mgo.v2/bson"
)

func New() *DatabaseClient {
    session, err := mgo.Dial("localhost")

    if err != nil {
        panic(err)
    }

    session.SetMode(mgo.Monotonic, true)

    return &DatabaseClient{Session : session}
}

type DatabaseClient struct {
    Session     *mgo.Session
}

func (client *DatabaseClient) Read(db string, collection string) map[string]interface{} {
    var result map[string]interface{}

    c := client.Session.DB(db).C(collection)

    err := c.Find(bson.M{"name": "Ale"}).One(&result)
    if err != nil {
        log.Fatal(err)
    }

    return result
}

func (client *DatabaseClient) Write(db string, collection string, data map[string]interface{}) error {
    return nil
}

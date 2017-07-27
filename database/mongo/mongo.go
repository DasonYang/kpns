package mongo

import (
    "log"
    "fmt"
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

func (client *DatabaseClient) Read(db string, collection string, query map[string]interface{}) map[string]interface{} {
    var result map[string]interface{}

    c := client.Session.DB(db).C(collection)

    err := c.Find(query).One(&result)
    if err != nil {
        log.Println(err)
    }

    return result
}

func (client *DatabaseClient) Write(db string, collection string, data map[string]interface{}) error {
    log.Printf("Insert data into mongo : db = %v, collection = %v\n", db, collection)
    c := client.Session.DB(db).C(collection)
    key := fmt.Sprintf("%v", data["key"])

    /*
    type ChangeInfo struct {
        Updated    int         // Number of existing documents updated
        Removed    int         // Number of documents removed
        UpsertedId interface{} // Upserted _id field, when not explicitly provided
    }
    */

    info, err := c.Upsert(bson.M{"key": key}, data)


    log.Printf("key = %v, updated = %v, removed = %v, upsertedId = %v\n", key, info.Updated, info.Removed, info.UpsertedId)

    return err
}

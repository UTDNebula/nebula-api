package log

// import (
// 	"context"
// 	"encoding/json"

// 	"github.com/UTDNebula/nebula-api/api/configs"
// 	"go.mongodb.org/mongo-driver/mongo"
// )

// var logCollection = configs.GetCollection(configs.DB, "logs")

// func (coll *mongo.Collection) Write(p []byte) (n int, err error) {
// 	var log_obj map[string]interface{}
// 	if err := json.Unmarshal(p, &log_obj); err != nil {
// 		return 0, err
// 	}

// 	res, insert_err := coll.InsertOne(context.TODO(), log_obj)
// 	if insert_err != nil {
// 		return 0, insert_err
// 	}
// }

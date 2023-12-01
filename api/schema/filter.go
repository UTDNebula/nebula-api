package schema

import (
	"github.com/gin-gonic/gin"
	gs "github.com/gorilla/schema"
	"go.mongodb.org/mongo-driver/bson"
	"net/url"
)

var decoder = makeDecoder()
var encoder = makeEncoder()

func makeDecoder() *gs.Decoder {
	dec := gs.NewDecoder()
	dec.IgnoreUnknownKeys(true)
	return dec
}

func makeEncoder() *gs.Encoder {
	enc := gs.NewEncoder()
	return enc
}

func FilterQuery[F any](c *gin.Context) (bson.M, error) {
	src := c.Request.URL.Query()
	dst := make(url.Values)
	filter := new(F)
	// decode in bson dst
	if err := decoder.Decode(filter, src); err != nil {
		return nil, err
	}
	if err := encoder.Encode(filter, dst); err != nil {
		return nil, err
	}
	query := make(bson.M)
	// merge dst into bson.M
	for k, v := range src {
		if dst.Has(k) {
			query[k] = v[0]
		}
	}

	return query, nil
}

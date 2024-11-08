package schema

import (
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
)

func FilterQuery[F any](c *gin.Context) (bson.M, error) {

	// Placeholder until filtering is fixed
	query := bson.M{}
	for k := range c.Request.URL.Query() {
		query[k] = c.Query(k)
	}

	/*
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
	*/

	return query, nil
}

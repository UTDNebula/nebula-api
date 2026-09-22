package loaders

import (
	"context"
	"net/http"
	"time"

	"github.com/UTDNebula/nebula-api/graphql/graph/model"
	"github.com/UTDNebula/nebula-api/shared/configs"

	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/vikstrous/dataloadgen"
)

type ctxKey string

const (
	loadersKey = ctxKey("dataloaders")
)

type Loaders struct {
	SectionLoader *dataloadgen.Loader[primitive.ObjectID, *model.Section]
}

// NewLoaders instantiates data loaders for the middleware
func NewLoaders() *Loaders {
	// Define the dataloaders
	sectionReader := &sectionReader{collection: configs.GetCollection("sections"),}

	return &Loaders{
		SectionLoader: dataloadgen.NewLoader(sectionReader.getSections, dataloadgen.WithWait(time.Millisecond)),
	}
}

// Middleware injects data loaders into the context
func Middleware(next http.Handler) http.Handler {
	// Return a middleware that injects the loader to the request context
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		loader := NewLoaders()
		r = r.WithContext(context.WithValue(r.Context(), loadersKey, loader))
		next.ServeHTTP(w, r)
	})
}

// For returns the dataloader for a given context
func For(ctx context.Context) *Loaders {
	return ctx.Value(loadersKey).(*Loaders)
}

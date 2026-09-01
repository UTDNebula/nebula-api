package configs

import (
	"context"
	"database/sql"
	"strconv"
	"sync"
	"time"

	"log"

	"github.com/XSAM/otelsql"
	"github.com/gin-gonic/gin"
	_ "github.com/jackc/pgx/v5/stdlib"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.opentelemetry.io/contrib/instrumentation/go.mongodb.org/mongo-driver/mongo/otelmongo"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.37.0"

	sentryotlp "github.com/getsentry/sentry-go/otel/otlp"
)

type DBSingleton struct {
	client *mongo.Client
}

var dbInstance *DBSingleton
var once sync.Once

func ConnectDB() *mongo.Client {
	once.Do(func() {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		opts := options.Client().
			ApplyURI(GetEnvMongoURI()).
			SetMonitor(otelmongo.NewMonitor()) // To trace Mongo operation on Sentry

		client, err := mongo.Connect(ctx, opts)
		if err != nil {
			log.Fatalf("Unable to create MongoDB client")
		}

		// ping the database
		err = client.Ping(ctx, nil)
		if err != nil {
			log.Fatalf("Unable to ping Mongo database")
		}

		log.Printf("Connected to MongoDB")

		dbInstance = &DBSingleton{
			client: client,
		}
	})

	return dbInstance.client
}

// getting database collections
func GetCollection(collectionName string) *mongo.Collection {
	client := ConnectDB()
	collection := client.Database("combinedDB").Collection(collectionName)
	return collection
}

// Returns *options.FindOptions with a limit and offset applied.
// Produces an error if user-provided offset isn't able to be parsed.
func GetOptionLimit(query *bson.M, c *gin.Context) (*options.FindOptions, error) {
	delete(*query, "offset") // removes offset (if present) in query --offset is not field in collections

	// parses offset if included in the query
	var offset int64
	var err error

	var limit int64 = GetEnvLimit()

	if c.Query("offset") == "" {
		offset = 0 // default value for offset
	} else {
		offset, err = strconv.ParseInt(c.Query("offset"), 10, 64)
		if err != nil {
			return options.Find().SetSkip(0).SetLimit(limit), err // default value for offset
		}
	}

	return options.Find().SetSkip(offset).SetLimit(limit), err
}

// Returns the offsets and limit for pagination stage for aggregate endpoints pipeline
func GetAggregateLimit(query *bson.M, c *gin.Context) (map[string]bson.D, error) {
	// Parses offsets if included in the query
	paginateMap := map[string]bson.D{
		"former_offset": {{Key: "$skip", Value: 0}}, // Init the default value of offset
		"latter_offset": {{Key: "$skip", Value: 0}},
		"limit":         {{Key: "$limit", Value: GetEnvLimit()}},
	}
	var err error

	for field := range paginateMap {
		if field != "limit" && c.Query(field) != "" {
			delete(*query, field)

			offset, err := strconv.ParseInt(c.Query(field), 10, 64)
			if err != nil {
				// Return default value of offset
				return paginateMap, err
			}
			paginateMap[field] = bson.D{{Key: "$skip", Value: offset}}
		}
	}

	return paginateMap, err
}

var clubsDbInstance *sql.DB
var clubOnce sync.Once

func ConnectClubsDB() *sql.DB {
	clubOnce.Do(func() {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		db, err := otelsql.Open("pgx", GetClubsDBUri(), otelsql.WithAttributes(
			semconv.DBSystemNamePostgreSQL,
		))
		if err != nil {
			log.Panic("Unable to connect to clubs database.")
		}

		// ping the database
		err = db.PingContext(ctx)
		if err != nil {
			log.Panic("Unable to ping Clubs database")
		}

		log.Printf("Connected to Clubs DB")

		clubsDbInstance = db
	})

	return clubsDbInstance
}

// InitOtelTracer initializes the opentelementry tracer that is exported to Sentry
func InitOtelTracer(ctx context.Context, sentryDsn string, sentryEnv string) *sdktrace.TracerProvider {
	traceExporter, err := sentryotlp.NewTraceExporter(ctx, sentryDsn)
	if err != nil {
		log.Fatalf("Unable to create Sentry trace exporter %v\n", err)
	}

	resource, err := resource.New(
		ctx,
		resource.WithAttributes(
			semconv.ServiceName("nebula-api"),
			semconv.DeploymentEnvironmentName(sentryEnv),
		),
	)
	if err != nil {
		log.Fatalf("Unable to create Otel resource %v\n", err)
	}

	var sampler sdktrace.Sampler
	switch sentryEnv {
	case "development":
		sampler = sdktrace.AlwaysSample()
	case "production":
		sampler = sdktrace.TraceIDRatioBased(0.10)
	default:
		sampler = sdktrace.TraceIDRatioBased(0.10)
	}

	tracerProvider := sdktrace.NewTracerProvider(
		sdktrace.WithSampler(sdktrace.ParentBased(sampler)),
		sdktrace.WithBatcher(traceExporter),
		sdktrace.WithResource(resource),
	)
	otel.SetTracerProvider(tracerProvider)

	otel.SetTextMapPropagator(
		propagation.NewCompositeTextMapPropagator(
			propagation.TraceContext{},
			propagation.Baggage{},
		),
	)

	otel.SetErrorHandler(otel.ErrorHandlerFunc(func(err error) {
		log.Printf("OTEL ERROR: %s", err.Error())
	}))

	log.Printf("Initialized Otel tracer")

	return tracerProvider
}

// ShutdownOtelTracer shuts down the otel tracer
func ShutdownOtelTracer(ctx context.Context, provider *sdktrace.TracerProvider) {
	shutdownCtx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	if err := provider.Shutdown(shutdownCtx); err != nil {
		log.Printf("failed to shutdown tracer: %v", err)
	}
}

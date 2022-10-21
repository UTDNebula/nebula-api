package main

import (
	"encoding/json"
	"flag"
	"github.com/UTDNebula/nebula-api/scraper/ndt/requirements"
	"io"
	"os"
	"strconv"
	"time"
	// "fmt"
	// "context"
	// "go.mongodb.org/mongo-driver/bson"
	// "go.mongodb.org/mongo-driver/mongo"
	// "go.mongodb.org/mongo-driver/mongo/options"
	// "go.mongodb.org/mongo-driver/mongo/readpref"
)

func main() {
	subtype := flag.String("t", "", "The subtype of the degree")
	school := flag.String("s", "", "The school that the degree falls under")
	name := flag.String("n", "", "The name of the degree")
	year := flag.String("y", strconv.Itoa(time.Now().Year()), "The year of the degree")
	abbr := flag.String("a", "", "The abbreviation of the degree")
	minCredHours := flag.Int("m", 0, "The minimum credit hours for the degree")
	cURI := flag.String("u", "", "The catalog URI of the degree")
	reqFile := flag.String("r", "", "The path to the JSON file containing the requirements")
	outFile := flag.String("o", "", "The path to the JSON file to create")

	flag.Parse()

	deg := requirements.Degree{
		*subtype,
		*school,
		*name,
		*year,
		*abbr,
		*minCredHours,
		*cURI,
		nil,
	}

	deg.Requirements = requirements.Parse(*reqFile, &deg)

	var writer io.Writer
	if *outFile == "" {
		// No output file specified, write to stdout
		writer = os.Stdout
	} else {
		// Output file specified, create or overwrite existing
		f, err := os.OpenFile(*outFile, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0777)
		defer f.Close()
		if err != nil {
			panic(err)
		}
		writer = f
	}
	enc := json.NewEncoder(writer)
	enc.SetIndent("", "\t")
	enc.Encode(deg)
}

// func mongoTest() {
// // Set up 10-second timeout context
// ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
// defer cancel()
// // Connect to DB, init client
// client, err := mongo.Connect(ctx, options.Client().ApplyURI("DB_URI"))
// // Defer client disconnect
// defer func() {
// if err = client.Disconnect(ctx); err != nil {
// panic(err)
// }
// }()

// // Get courses database
// courses := client.Database("combinedDB").Collection("courses")
// // Find ACCT 2302
// filter := bson.D{{"course_number", "2302"}, {"subject_prefix", "ACCT"}}
// cursor, err := courses.Find(ctx, filter)
// if err != nil {
// panic(err)
// }
// // Unmarshall results as bson.D types (key/value slices)
// var results []bson.D
// if err = cursor.All(ctx, &results); err != nil {
// panic(err)
// }

// // Iterate over results
// for _, result := range results {
// // Get map form of result
// resMap := result.Map()

// // Re-marshall prereqs, unmarshall as a CollectionRequirement to collectionReq
// doc, err := bson.Marshal(resMap["prerequisites"])
// if err != nil {
// panic(err)
// }
// var collectionReq requirements.CollectionRequirement
// err = bson.Unmarshal(doc, &collectionReq)
// if err != nil {
// panic(err)
// }

// // Iterate over and print all course requirements
// for _, option := range collectionReq.Options {
// // Re-marshall the interface{} option, unmarshall as a CourseRequirement to courseReq
// doc, err = bson.Marshal(option)
// if err != nil {
// panic(err)
// }
// var courseReq requirements.CourseRequirement
// err = bson.Unmarshal(doc, &courseReq)
// if err != nil {
// panic(err)
// }
// // Print courseReq
// fmt.Println(courseReq)
// }
// }

// Collection requirement for testing purposes
// cr := requirements.NewCollectionRequirement("TestCollection", 2, []interface{}{
// requirements.NewCourseRequirement("ref", "C+"),
// requirements.NewCoreRequirement("16", 6),
// requirements.NewHoursRequirement(12, []*requirements.CourseRequirement{
// requirements.NewCourseRequirement("ref", "B-"),
// requirements.NewCourseRequirement("ref2", "C+"),
// }),
// })

// Uncomment to test inserting cr
// _, err = courses.InsertOne(context.TODO(), cr)
// if err != nil {
// panic(err)
// }

// Uncomment to test json-encoding cr
// enc := json.NewEncoder(os.Stdout)
// enc.SetIndent("", "\t")
// enc.Encode(cr)
// }

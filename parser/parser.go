// go run parser.go [file_path] [semster]
// example: go run parser.go "Fall 2019.csv" 19F

package main

import (
	"context"
	"encoding/csv"
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strconv"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"

	"parser/model"
)

func getCollection(client *mongo.Client, database string, collection string) (returnCollection *mongo.Collection) {
	return (client.Database(database).Collection(collection))
}

func DBConnect(URI string) (client *mongo.Client) {
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(URI))
	if err != nil {
		log.Panicf(err.Error())
	}
	if err := client.Ping(context.TODO(), readpref.Primary()); err != nil {
		log.Panicf(err.Error())
	}
	return client
}

func csvToClassesSlice(csvFile *os.File, logFile *os.File) (classes []model.Class) {
	reader := csv.NewReader(csvFile)
	records, err := reader.ReadAll() // records is [][]strings
	if err != nil {
		log.Panicf(err.Error())
	}
	// look for the subject column and w column
	subjectCol := -1
	catalogNumberCol := -1
	sectionCol := -1
	wCol := -1
	aPlusCol := -1
	for j := 0; j < len(records[0]); j++ {
		switch {
		case records[0][j] == "Subject":
			subjectCol = j
		case records[0][j] == "Catalog Number" || records[0][j] == "Catalog Nbr":
			catalogNumberCol = j
		case records[0][j] == "Section":
			sectionCol = j
		case records[0][j] == "W" || records[0][j] == "Total W" || records[0][j] == "W Total":
			wCol = j
		case records[0][j] == "A+":
			aPlusCol = j
		}
		if wCol == -1 || subjectCol == -1 || catalogNumberCol == -1 || sectionCol == -1 || aPlusCol == -1 {
			continue
		} else {
			break
		}
	}
	if wCol == -1 {
		logFile.WriteString("could not find W column")
		log.Panicf("could not find W column")
	}
	if sectionCol == -1 {
		logFile.WriteString("could not find Section column")
		log.Panicf("could not find Section column")
	}
	if subjectCol == -1 {
		logFile.WriteString("could not find Subject column")
		log.Panicf("could not find Subject column")
	}
	if aPlusCol == -1 {
		logFile.WriteString("could not find A+ column")
		log.Panicf("could not find A+ column")
	}

	for i := 1; i < len(records); i++ {
		// convert grade distribution from string to int
		var tempSlice bson.A
		for j := 0; j < 13; j++ {
			var tempInt int
			fmt.Sscan(records[i][aPlusCol+j], &tempInt)
			tempSlice = append(tempSlice, tempInt)
		}
		// add w number to the grade_distribution slice
		var tempInt int
		if wCol != -1 {
			fmt.Sscan(records[i][wCol], &tempInt)
		}
		tempSlice = append(tempSlice, tempInt)
		// add new class to classes slice
		classes = append(classes,
			model.Class{
				Subject:           records[i][subjectCol],
				CatalogNumber:     records[i][catalogNumberCol],
				Section:           records[i][sectionCol],
				GradeDistribution: tempSlice,
			})
	}
	return classes
}

// inserts grades into mongodb
func insertGrades(sectionsCollection *mongo.Collection, coursesCollection *mongo.Collection, classes []model.Class, academicSession string, logFile *os.File) {
	for i := 0; i < len(classes); i++ {
		var courseSearchResult model.Course
		err := coursesCollection.FindOne(context.TODO(), bson.D{{"course_number", classes[i].CatalogNumber}, {"subject_prefix", classes[i].Subject}}).Decode(&courseSearchResult)
		// if class is not in courses section
		if err != nil {
			// log that class could not be found
			logFile.WriteString("could not find course " + classes[i].Subject + " " + classes[i].CatalogNumber + ": " + err.Error() + "\n")
			fmt.Println("could not find course " + classes[i].Subject + " " + classes[i].CatalogNumber + ": " + err.Error())
			// fmt.Println(err)
			continue
		}
		var data []bson.M
		match :=
			bson.D{
				{"$match",
					bson.D{
						{"course_number", classes[i].CatalogNumber},
						{"subject_prefix", classes[i].Subject},
					},
				},
			}
		lookup :=
			bson.D{
				{"$lookup",
					bson.D{
						{"from", "sections"},
						{"localField", "sections"},
						{"foreignField", "_id"},
						{"as", "sections"},
					},
				},
			}
		unwind := bson.D{{"$unwind", bson.D{{"path", "$sections"}}}}
		matchSectionNo :=
			bson.D{
				{"$match",
					bson.D{
						{"sections.section_number", classes[i].Section},
						{"sections.academic_session.name", academicSession},
					},
				},
			}
		project := bson.D{{"$project", bson.D{{"section", "$sections"}}}}
		set :=
			bson.D{
				{"$set",
					bson.D{
						{"section.grade_distribution", classes[i].GradeDistribution},
					},
				},
			}
		cursor, err := coursesCollection.Aggregate(context.TODO(), mongo.Pipeline{match, lookup, unwind, matchSectionNo, project, set})
		if err != nil {
			logFile.WriteString(classes[i].Subject + " " + classes[i].CatalogNumber + "." + classes[i].Section + ": " + err.Error() + "\n")
			fmt.Println(classes[i].Subject + " " + classes[i].CatalogNumber + "." + classes[i].Section + ": " + err.Error())
			continue
		}
		err = cursor.All(context.TODO(), &data) // put cursor data into []primitive.M
		if err != nil {
			logFile.WriteString(classes[i].Subject + " " + classes[i].CatalogNumber + "." + classes[i].Section + ": " + err.Error() + "\n")
			fmt.Println(classes[i].Subject + " " + classes[i].CatalogNumber + "." + classes[i].Section + ": " + err.Error())
			continue
		}
		// if more than 1 result, log error and continue
		if len(data) != 1 {
			logFile.WriteString(classes[i].Subject + " " + classes[i].CatalogNumber + "." + classes[i].Section + ": recieved results " + strconv.Itoa(len(data)) + " from aggregation, expected 1\n")
			fmt.Println(classes[i].Subject + " " + classes[i].CatalogNumber + "." + classes[i].Section + ": recieved results " + strconv.Itoa(len(data)) + " from aggregation, expected 1")
			continue
		}

		section, ok := data[0]["section"].(primitive.M)
		if !ok { // if section assertion not work
			logFile.WriteString(classes[i].Subject + " " + classes[i].CatalogNumber + "." + classes[i].Section + "unable to get section from db")
			fmt.Println(classes[i].Subject + " " + classes[i].CatalogNumber + "." + classes[i].Section + "unable to get section from db")
			continue
		}
		sectionID, ok := section["_id"].(primitive.ObjectID)
		if !ok { // if sections id assertion not work
			logFile.WriteString(classes[i].Subject + " " + classes[i].CatalogNumber + "." + classes[i].Section + "unable to get section id from db")
			fmt.Println(classes[i].Subject + " " + classes[i].CatalogNumber + "." + classes[i].Section + "unable to get section id from db")
		}
		var gradeDistribution bson.A = data[0]["section"].(primitive.M)["grade_distribution"].(bson.A)
		_, err = sectionsCollection.UpdateByID(context.TODO(), sectionID, bson.D{{"$set", bson.D{{"grade_distribution", gradeDistribution}}}})
		if err != nil {
			fmt.Println(classes[i].Subject + " " + classes[i].CatalogNumber + ": " + err.Error())
			logFile.WriteString(classes[i].Subject + " " + classes[i].CatalogNumber + "." + classes[i].Section + ": " + err.Error() + "\n")
			continue
		}
		fmt.Println("added " + classes[i].Subject + " " + classes[i].CatalogNumber + "." + classes[i].Section)
	}
}

func EnvMongoURI() string {

	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file: %v", err)
	}

	return os.Getenv("MONGODB_URI")
}

func main() {
	URI := EnvMongoURI()
	fileFlag := flag.String("file", "", "csv grade file to be parsed")
	semesterFlag := flag.String("semester", "", "semester of the grades, ex: 18U, 19F")
	flag.Parse()
	csvPath := *fileFlag
	academicSession := *semesterFlag
	csvFile, err := os.Open(csvPath)
	if err != nil {
		fmt.Println("could not open file " + csvPath)
		log.Panicf(err.Error())
	}
	defer csvFile.Close()

	// create logs directory
	if _, err := os.Stat("logs"); err != nil {
		os.Mkdir("logs", os.ModePerm)
	}
	// create log file [name of csv].log in logs directory
	logFileName := filepath.Base(csvPath)
	logFile, err := os.Create("logs/" + logFileName + ".log")
	if err != nil {
		fmt.Println("could not create log file")
		log.Panicf(err.Error())
	}
	defer logFile.Close()

	// put class data from csv into classes slice
	classes := csvToClassesSlice(csvFile, logFile)
	client := DBConnect(URI)

	// insert grades into mongodb
	insertGrades(getCollection(client, "combinedDB", "sections"), getCollection(client, "combinedDB", "courses"), classes, academicSession, logFile)
}

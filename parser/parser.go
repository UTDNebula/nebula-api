// go run parser.go [file_path] [semster]
// example: go run parser.go "Fall 2019.csv" 19F

package main

import (
	"context"
	"encoding/csv"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"

	"parser/configs"
	"parser/model"
)

type Class struct {
	subject           string
	catalogNumber     string
	section           string
	gradeDistribution bson.A
}

func getCollection(client *mongo.Client, database string, collection string) (returnCollection *mongo.Collection) {
	return (client.Database(database).Collection(collection))
}

func DBConnect(URI string) (client *mongo.Client) {
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(URI))
	if err != nil {
		fmt.Println(err)
		panic(err)
	}
	if err := client.Ping(context.TODO(), readpref.Primary()); err != nil {
		fmt.Println(err)
		panic(err)
	}
	return client
}

func csvToClassesSlice(csvFile *os.File, logFile *os.File) (classes []Class) {
	reader := csv.NewReader(csvFile)
	records, err := reader.ReadAll() // records is [][]strings
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	// look for the subject column and w column
	subjectCol := -1
	wCol := -1
	for j := 0; j < len(records[0]); j++ {
		if strings.Compare(records[0][j], "Subject") == 0 {
			subjectCol = j
		}
		if strings.Compare(records[0][j], "W") == 0 || strings.Compare(records[0][j], "Total W") == 0 || strings.Compare(records[0][j], "W Total") == 0 {
			wCol = j
		}
		if wCol == -1 || subjectCol == -1 {
			continue
		} else {
			break
		}
	}
	if wCol == -1 {
		logFile.WriteString("could not find W column")
	}
	catalogNumberCol := subjectCol + 1
	sectionCol := subjectCol + 2

	for i := 1; i < len(records); i++ {
		// convert grade distribution from string to int
		var tempSlice bson.A
		for j := 0; j < 13; j++ {
			var tempInt int
			fmt.Sscan(records[i][3+subjectCol+j], &tempInt)
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
			Class{
				subject:           records[i][subjectCol],
				catalogNumber:     records[i][catalogNumberCol],
				section:           records[i][sectionCol],
				gradeDistribution: tempSlice,
			})
	}
	return classes
}

// inserts grades into mongodb
func insertGrades(sectionsCollection *mongo.Collection, coursesCollection *mongo.Collection, classes []Class, academicSession string, logFile *os.File) {
	for i := 0; i < len(classes); i++ {
		var courseSearchResult model.Course
		err := coursesCollection.FindOne(context.TODO(), bson.D{{"course_number", classes[i].catalogNumber}, {"subject_prefix", classes[i].subject}}).Decode(&courseSearchResult)
		// if class is not in courses section
		if err != nil {
			// log that class could not be found
			logFile.WriteString("could not find course " + classes[i].subject + " " + classes[i].catalogNumber + ": " + err.Error() + "\n")
			fmt.Println("could not find course " + classes[i].subject + " " + classes[i].catalogNumber + ": " + err.Error())
			// fmt.Println(err)
			continue
		}
		var data []bson.M
		match :=
			bson.D{
				{"$match",
					bson.D{
						{"course_number", classes[i].catalogNumber},
						{"subject_prefix", classes[i].subject},
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
						{"sections.section_number", classes[i].section},
						{"sections.academic_session.name", academicSession},
					},
				},
			}
		project := bson.D{{"$project", bson.D{{"section", "$sections"}}}}
		set :=
			bson.D{
				{"$set",
					bson.D{
						{"section.grade_distribution", classes[i].gradeDistribution},
					},
				},
			}
		cursor, err := coursesCollection.Aggregate(context.TODO(), mongo.Pipeline{match, lookup, unwind, matchSectionNo, project, set})
		if err != nil {
			logFile.WriteString(classes[i].subject + " " + classes[i].catalogNumber + "." + classes[i].section + ": " + err.Error() + "\n")
			fmt.Println(classes[i].subject + " " + classes[i].catalogNumber + "." + classes[i].section + ": " + err.Error())
			continue
		}
		err = cursor.All(context.TODO(), &data) // put cursor data into []primitive.M
		if err != nil {
			logFile.WriteString(classes[i].subject + " " + classes[i].catalogNumber + "." + classes[i].section + ": " + err.Error() + "\n")
			fmt.Println(classes[i].subject + " " + classes[i].catalogNumber + "." + classes[i].section + ": " + err.Error())
			continue
		}
		// if more than 1 result, log error and continue
		if len(data) != 1 {
			logFile.WriteString(classes[i].subject + " " + classes[i].catalogNumber + "." + classes[i].section + ": recieved results " + strconv.Itoa(len(data)) + " from aggregation, expected 1\n")
			fmt.Println(classes[i].subject + " " + classes[i].catalogNumber + "." + classes[i].section + ": recieved results " + strconv.Itoa(len(data)) + " from aggregation, expected 1")
			continue
		}
		sectionID := (data[0]["section"].(primitive.M))["_id"].(primitive.ObjectID)
		var gradeDistribution bson.A = data[0]["section"].(primitive.M)["grade_distribution"].(bson.A)
		_, err = sectionsCollection.UpdateByID(context.TODO(), sectionID, bson.D{{"$set", bson.D{{"grade_distribution", gradeDistribution}}}})
		if err != nil {
			fmt.Println(classes[i].subject + " " + classes[i].catalogNumber + ": " + err.Error())
			logFile.WriteString(classes[i].subject + " " + classes[i].catalogNumber + "." + classes[i].section + ": " + err.Error() + "\n")
			continue
		}
		fmt.Println("added " + classes[i].subject + " " + classes[i].catalogNumber + "." + classes[i].section)
	}
}

func main() {
	URI := configs.EnvMongoURI()
	fileFlag := flag.String("file", "", "csv grade file to be parsed")
	semesterFlag := flag.String("semester", "", "semester of the grades, ex: 18U, 19F")
	flag.Parse()
	csvPath := *fileFlag
	academicSession := *semesterFlag
	csvFile, err := os.Open(csvPath)
	if err != nil {
		fmt.Println("could not open file " + csvPath)
		fmt.Println(err)
		os.Exit(1)
	}

	// create logs directory
	if _, err := os.Stat("logs"); err != nil {
		os.Mkdir("logs", os.ModePerm)
	}
	// create log file [name of csv].log in logs directory
	logFileName := filepath.Base(csvPath)
	logFile, err := os.Create("logs/" + logFileName + ".log")
	if err != nil {
		fmt.Println("could not create log file")
		fmt.Println(err)
		os.Exit(1)
	}

	// put class data from csv into classes slice
	classes := csvToClassesSlice(csvFile, logFile)
	client := DBConnect(URI)

	// insert grades into mongodb
	insertGrades(getCollection(client, "combinedDB", "sections"), getCollection(client, "combinedDB", "courses"), classes, academicSession, logFile)
}

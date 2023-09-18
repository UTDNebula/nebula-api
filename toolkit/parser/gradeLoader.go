package parser

import (
	"encoding/csv"
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
)

func loadGrades(csvDir string) map[string]map[string][]int {

	// MAP[SEMESTER] -> MAP[SUBJECT + NUMBER + SECTION] -> GRADE DISTRIBUTION
	gradeMap := make(map[string]map[string][]int)

	if csvDir == "" {
		fmt.Printf("No grade data CSV directory specified. Grade data will not be included.\n")
		return gradeMap
	}

	dirPtr, err := os.Open(csvDir)
	if err != nil {
		panic(err)
	}
	defer dirPtr.Close()

	csvFiles, err := dirPtr.ReadDir(-1)
	if err != nil {
		panic(err)
	}

	for _, csvEntry := range csvFiles {

		if csvEntry.IsDir() {
			continue
		}

		csvPath := fmt.Sprintf("%s/%s", csvDir, csvEntry.Name())

		csvFile, err := os.Open(csvPath)
		if err != nil {
			panic(err)
		}
		defer csvFile.Close()

		// Create logs directory
		if _, err := os.Stat("logs"); err != nil {
			os.Mkdir("logs", os.ModePerm)
		}

		// Create log file [name of csv].log in logs directory
		basePath := filepath.Base(csvPath)
		csvName := strings.TrimSuffix(basePath, filepath.Ext(basePath))
		logFile, err := os.Create("logs/" + csvName + ".log")

		if err != nil {
			panic(errors.New("Could not create CSV log file."))
		}
		defer logFile.Close()

		// Put data from csv into map
		gradeMap[csvName] = csvToMap(csvFile, logFile)
	}

	return gradeMap
}

func csvToMap(csvFile *os.File, logFile *os.File) map[string][]int {
	reader := csv.NewReader(csvFile)
	records, err := reader.ReadAll() // records is [][]strings
	if err != nil {
		log.Panicf("Error parsing %s: %s", csvFile.Name(), err.Error())
	}
	// look for the subject column and w column
	subjectCol := -1
	catalogNumberCol := -1
	sectionCol := -1
	wCol := -1
	aPlusCol := -1

	headerRow := records[0]

	for j := 0; j < len(headerRow); j++ {
		switch {
		case headerRow[j] == "Subject":
			subjectCol = j
		case headerRow[j] == "Catalog Number" || headerRow[j] == "Catalog Nbr":
			catalogNumberCol = j
		case headerRow[j] == "Section":
			sectionCol = j
		case headerRow[j] == "W" || headerRow[j] == "Total W" || headerRow[j] == "W Total":
			wCol = j
		case headerRow[j] == "A+":
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
		//log.Panicf("could not find W column")
	}
	if sectionCol == -1 {
		logFile.WriteString("could not find Section column")
		log.Panicf("could not find Section column")
	}
	if subjectCol == -1 {
		logFile.WriteString("could not find Subject column")
		log.Panicf("could not find Subject column")
	}
	if catalogNumberCol == -1 {
		logFile.WriteString("could not find catalog # column")
		log.Panicf("could not find catalog # column")
	}
	if aPlusCol == -1 {
		logFile.WriteString("could not find A+ column")
		log.Panicf("could not find A+ column")
	}

	distroMap := make(map[string][]int)

	for _, record := range records {
		// convert grade distribution from string to int
		intSlice := make([]int, 0, 13)
		var tempInt int

		for j := 0; j < 13; j++ {
			fmt.Sscan(record[aPlusCol+j], &tempInt)
			intSlice = append(intSlice, tempInt)
		}
		// add w number to the grade_distribution slice
		if wCol != -1 {
			fmt.Sscan(record[wCol], &tempInt)
		}
		intSlice = append(intSlice, tempInt)

		// add new grade distribution to map, keyed by SUBJECT + NUMBER + SECTION
		distroKey := record[subjectCol] + record[catalogNumberCol] + record[sectionCol]
		distroMap[distroKey] = intSlice
	}
	return distroMap
}

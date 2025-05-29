package main

import (
	"fmt"
	"flag"
	"utilsDB/utils"
)

func main() {
	oldFile := flag.String("old", "", "Path to the old database file")
	newFile := flag.String("new", "", "Path to the new database file")
	flag.Parse()

	var err error
	var oldReader, newReader utils.DBReader
	var oldDB, newDB *utils.Recipes
	if *oldFile == "" || *newFile == "" {
		fmt.Println("Please provide both --old and --new file paths.")
		return
	}

	_, oldReader, err = utils.NewDBReader(*oldFile)
	if err == nil {
		oldDB, err = oldReader.Read(*oldFile)
		if err == nil {
			_, newReader, err = utils.NewDBReader(*newFile)
			if err == nil {
				newDB, err = newReader.Read(*newFile)
			}
		}
	}

	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	utils.CompareDatabases(oldDB, newDB)
}

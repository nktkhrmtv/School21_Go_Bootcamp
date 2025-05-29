package main

import (
	"flag"
	"fmt"
	"utilsDB/utils"
)

func main() {
	filename := flag.String("f", "", "Path to the database file")
	flag.Parse()

	var err error
	var fileType string
	var reader utils.DBReader
	var recipes *utils.Recipes

	fileType, reader, err = utils.NewDBReader(*filename)
	if err == nil {
		recipes, err = reader.Read(*filename)
		if err == nil {
			err = utils.ConvertAndPrint(recipes, fileType)
		}
	}

	if err != nil {
		fmt.Println("Error:", err)
	}
}

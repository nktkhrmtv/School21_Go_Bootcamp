package utils

import (
	"encoding/json"
	"encoding/xml"
	"errors"	
	"fmt"
	"io"
	"os"
	"strings"
)

type Ingredient struct {
	Name  string  `xml:"itemname" json:"ingredient_name"`
	Count string `xml:"itemcount" json:"ingredient_count"`
	Unit  string  `xml:"itemunit,omitempty" json:"ingredient_unit,omitempty"`
}
type Cake struct {
	Name        string       `xml:"name" json:"name"`
	StoveTime   string       `xml:"stovetime" json:"time"`
	Ingredients []Ingredient `xml:"ingredients>item" json:"ingredients"`
}
type Recipes struct {
	Cakes []Cake `xml:"cake" json:"cake"`
}


type DBReader interface {
	Read(filename string) (*Recipes, error)
}

type XMLReader struct{}
func (x XMLReader) Read(filename string) (*Recipes, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	data, err := io.ReadAll(file)
	if err != nil {
		return nil, err
	}

	var recipes Recipes
	err = xml.Unmarshal(data, &recipes)
	if err != nil {
		return nil, err
	}

	return &recipes, nil
}

type JSONReader struct{}
func (j JSONReader) Read(filename string) (*Recipes, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close() 
	data, err := io.ReadAll(file)
	if err != nil {
		return nil, err
	}

	var recipes Recipes
	err = json.Unmarshal(data, &recipes)
	if err != nil {
		return nil, err
	}

	return &recipes, nil
}


func NewDBReader(filename string) (string, DBReader, error) {
	fileType := ""
	if strings.HasSuffix(filename, ".xml") {
		fileType = "xml"
		return fileType, XMLReader{}, nil
	} else if strings.HasSuffix(filename, ".json") {
		fileType = "json"
		return fileType, JSONReader{}, nil
	}
	return fileType, nil, errors.New("unsupported file format")
}

func ConvertAndPrint(recipes *Recipes, inputFormat string) error {
	var output []byte
	var err error

	if inputFormat == "xml" {
		output, err = json.MarshalIndent(recipes, "", "    ")
		if err != nil {
			return err
		}
		fmt.Println(string(output))
	} else if inputFormat == "json" {
		output, err = xml.MarshalIndent(recipes, "", "    ")
		if err != nil {
			return err
		}
		fmt.Println(string(output))
	} else {
		return errors.New("unsupported input format")
	}

	return nil
}


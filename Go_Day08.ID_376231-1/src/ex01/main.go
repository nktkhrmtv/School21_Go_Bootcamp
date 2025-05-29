package main

import (
	"fmt"
	"reflect"
	"strings"
)

type UnknownPlant struct {
	FlowerType string
	LeafType   string
	Color      int `color_scheme:"rgb"`
}

type AnotherUnknownPlant struct {
	FlowerColor int
	LeafType    string
	Height      int `unit:"inches"`
}

func describePlant(plant interface{}) {
	val := reflect.ValueOf(plant) 
	typ := reflect.TypeOf(plant)   

	for i := 0; i < typ.NumField(); i++ {
		field := typ.Field(i) 
		value := val.Field(i)

		fieldName := field.Name

		tagName := strings.Split(string(field.Tag), ":")[0]
		tag :=  field.Tag.Get(string(tagName)) 
		if tag != "" {
			fieldName = fmt.Sprintf("%s(%s=%s)", fieldName, strings.Split(string(field.Tag), ":")[0], tag)
		}

		fmt.Printf("%s:%v\n", fieldName, value.Interface())
	}
}

func main() {
	plant1 := UnknownPlant{
		FlowerType: "rose",
		LeafType:   "oval",
		Color:      255,
	}

	plant2 := AnotherUnknownPlant{
		FlowerColor: 555,
		LeafType:    "lanceolate",
		Height:      15,
	}

	fmt.Println("Plant 1:")
	describePlant(plant1)

	fmt.Println("\nPlant 2:")
	describePlant(plant2)
}
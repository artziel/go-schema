package main

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/artziel/go-schema"
)

type Sample struct {
	ID     uint   `json:"id" schema:"name:id,required"`
	Status string `json:"status" schema:"name:status,required,restrictTo:'IS  ENABLED' DISABLED LOCKED"`
}

func main() {

	val := Sample{
		ID:     1,
		Status: "IS ENABLED",
	}

	result, err := schema.Validate(&val)
	if err != nil {
		fmt.Printf("%s\n%s\n", err.Error(), strings.Repeat("-", 60))
		formatted, _ := json.MarshalIndent(result, " ", "    ")
		fmt.Printf("%s\n", formatted)
	} else {
		fmt.Println("No Errors Found!")
	}

}

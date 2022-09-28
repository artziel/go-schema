package main

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/artziel/go-schema"
)

type Sample struct {
	ID       uint   `json:"id" schema:"name:id,required"`
	Username string `json:"username" schema:"name:username,require:Password"`
	Password string `json:"password" schema:"name:password"`
	Status   string `json:"status" schema:"name:status,restrictTo:ENABLED DISABLED LOCKED,require:Username"`
}

func main() {

	val := Sample{
		ID:       1,
		Username: "artziel",
		Status:   "ENABLED",
		Password: "123",
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

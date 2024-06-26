package main

import (
	"encoding/json"
	"fmt"
	"os"
)

type From struct {
	Name    string
	Address string
}

type Data struct {
	From      From
	Addresses []string
	Subject   string
	Body      string
}

func getData() (Data, error) {
	b, err := os.ReadFile("data.json")
	if err != nil {
		return Data{}, fmt.Errorf("Unable to read data file: %w", err)
	}

	data := Data{}

	err = json.Unmarshal(b, &data)
	if err != nil {
		return Data{}, fmt.Errorf("Unable to unmarshal data: %w", err)
	}

	return data, nil
}

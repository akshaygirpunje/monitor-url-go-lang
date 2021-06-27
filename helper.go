package main

import (
	"encoding/json"
	"fmt"
	"math"
	"os"
)

// Configuration for URLs Conf file
type Configuration struct {
	URLs []string
}

// GetStatus - returns 1 if site is Up, 0 if Down
func GetStatus(code int) (int, error) {

	if code == int(math.NaN()) {
		return -1, fmt.Errorf("Not a valid status code")
	}

	if code == 200 {
		return 1, nil
	}
	return 0, nil
}

// LoadConfig - Reads a JSON configuration file and outputs a Configuration struct
func LoadConfig(path string) (Configuration, error) {
	file, err := os.Open(path)
	if err != nil {
		fmt.Println("File Input Error: ", err)
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	configuration := Configuration{}
	derr := decoder.Decode(&configuration)

	return configuration, derr
}

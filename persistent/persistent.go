package persistent

import (
	"encoding/json"
	"os"
	"time"
)

type RequestData struct {
	Timestamp   time.Time `json:"timestamp"`
	Count       [60]int   `json:"count"`
	Acummulated [60]int   `json:"acummulated"`
}

func (data *RequestData) GetNumRequests() int {
	numRequests := 0

	diffTime := time.Since(data.Timestamp).Seconds()
	if diffTime < 60 {
		numRequests = data.Acummulated[data.Timestamp.Second()]
	}
	return numRequests
}

func (data *RequestData) ReadFromFile(path string) error {

	file, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	err = json.Unmarshal([]byte(file), data)
	if err != nil {
		return err
	}

	return nil
}

func (data *RequestData) WriteToFile(path string) error {

	file, err := json.MarshalIndent(data, "", " ")
	if err != nil {
		return err
	}
	err = os.WriteFile(path, file, 0644)
	if err != nil {
		return err
	}
	return nil
}

package handlers

import (
	"encoding/json"
	"os"
	"strings"
)

func ReadCassette(cassette string) (string, error) {
	data, err := os.ReadFile(cassette + ".cassette")
	if err != nil {
		return "", err
	}
	return string(data), nil
}

func WriteCassette(cassette string, response string) error {
	data := []byte(response)
	return os.WriteFile(cassette+".cassette", data, 0644)
}

type CassetteData struct {
	Prompt   string
	Response string
}

func ReadCassetteData(cassette string) (CassetteData, error) {
	data, err := os.ReadFile(cassette + ".cassette")
	if err != nil {
		return CassetteData{}, err
	}
	var cassetteData CassetteData
	err = json.Unmarshal(data, &cassetteData)
	if err != nil {
		return CassetteData{}, err
	}
	return cassetteData, nil
}

func WriteCassetteData(cassette string, prompt string, response string) error {
	cassetteData := CassetteData{
		Prompt:   prompt,
		Response: response,
	}
	data, err := json.Marshal(cassetteData)
	if err != nil {
		return err
	}
	return os.WriteFile(cassette+".cassette", data, 0644)
}

func CassetteName(testName string) string {
	name := strings.ToLower(testName)
	name = strings.ReplaceAll(name, " ", "_")
	return name
}

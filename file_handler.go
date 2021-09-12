package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/reconquest/karma-go"
	"github.com/reconquest/pkg/log"
)

const (
	VMI_EXTENTION = ".vmi"
	VMB_EXTENTION = ".vmb"
)

var (
	PATH_TO_INPUT_DIR  = "empty"
	PATH_TO_RESULT_DIR = "empty"
)

type VMI_FILE struct {
	ID            string    `json:"id"`
	DisplayName   string    `json:"displayName"`
	Group         string    `json:"group"`
	Region        string    `json:"region"`
	Min           string    `json:"min"`
	Max           string    `json:"max"`
	NumDeltas     string    "numDeltas"
	IsPoseControl string    "isPoseControl"
	Formula       []Formula `json:"formulas"`
}

type Formula struct {
	TargetType string `json:"targetType"`
	Target     string `json:"target"`
	Multiplier string `json:"multiplier"`
}

func getFileNamesInDir() ([]string, error) {
	var result []string
	path := PATH_TO_INPUT_DIR
	files, err := ioutil.ReadDir(path)
	if err != nil {
		return nil, karma.Format(
			err,
			"unable to read directory by path: %s",
			path,
		)
	}

	for _, file := range files {
		if strings.Contains(file.Name(), VMI_EXTENTION) {
			result = append(result, file.Name())
		}
	}

	return result, nil
}

func handleVMIFile(fileName string) error {
	firstFilePath := fmt.Sprintf("%s/%s", PATH_TO_INPUT_DIR, fileName)
	log.Infof(nil, "open file by path:%s", firstFilePath)

	file, err := os.Open(firstFilePath)
	if err != nil {
		log.Fatal(err)
	}

	defer func() {
		err = file.Close()
		if err != nil {
			log.Fatal(err)
		}
	}()

	log.Info("read file bytes")
	fileBytes, err := ioutil.ReadAll(file)
	if err != nil {
		return karma.Format(
			err,
			"unable to read file by path: %s",
			firstFilePath,
		)
	}

	log.Info("marshall json file bytes")
	var data VMI_FILE
	err = json.Unmarshal(fileBytes, &data)
	if err != nil {
		return karma.Format(
			err,
			"unable to unmarshal file by path: %s",
			firstFilePath,
		)
	}

	// edit necessary file strings
	displayName := data.DisplayName
	if data.ID != displayName {
		data.ID = displayName
	}

	// Convert golang object back to byte
	editedFileData, err := json.Marshal(data)
	if err != nil {
		return karma.Format(
			err,
			"unable to marshall data to file by path: %s",
			firstFilePath,
		)
	}

	// creating directory
	var pathToResultDir string
	if data.Group != "" {
		pathToResultDir = PATH_TO_RESULT_DIR + data.Group
	} else {
		pathToResultDir = PATH_TO_RESULT_DIR + data.Region
	}

	log.Infof(nil, "creating directory by path: %s", pathToResultDir)
	if _, err := os.Stat(pathToResultDir); os.IsNotExist(err) {
		err := os.MkdirAll(pathToResultDir, 0755)
		if err != nil {
			return karma.Format(
				err,
				"unable to create directory by path: %s",
				pathToResultDir,
			)
		}
	}

	// copy file to result directory

	pathToResultFileVMI := pathToResultDir + "/" + fileName
	log.Infof(nil, "copy .vmi file to directory by path: %s", pathToResultFileVMI)
	err = ioutil.WriteFile(pathToResultFileVMI, editedFileData, 0644)
	if err != nil {
		return karma.Format(
			err,
			"unable to write file by path: %s",
			pathToResultFileVMI,
		)
	}
	// copy .vmb file to result directory

	pathToCurrentFileVMB := strings.Replace(firstFilePath, ".vmi", ".vmb", -1)
	log.Infof(nil, "read .vmb file by path: %s", pathToCurrentFileVMB)
	fileVMB, err := ioutil.ReadFile(pathToCurrentFileVMB)
	if err != nil {
		return karma.Format(
			err,
			"unable to read file by path: %s",
			pathToCurrentFileVMB,
		)
	}

	pathToResultFileVMB := pathToResultDir + "/" + strings.Replace(fileName, ".vmi", ".vmb", -1)
	log.Infof(nil, "copy .vmb file by path: %s", pathToResultFileVMB)
	err = ioutil.WriteFile(pathToResultFileVMB, fileVMB, 0644)
	if err != nil {
		return karma.Format(
			err,
			"unable to write file by path: %s",
			pathToResultFileVMB,
		)
	}

	return nil
}

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
	PATH_TO_DIR   = "./Morphs_main/Morphs/Male_Genitalia"
	ROOT_DIR      = "./"
	RESULT_DIR    = "./Morphs_result/"
	VMI_EXTENTION = ".vmi"
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
	path := PATH_TO_DIR
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
	firstFilePath := fmt.Sprintf("%s/%s", PATH_TO_DIR, fileName)
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

	fileBytes, err := ioutil.ReadAll(file)
	if err != nil {
		return karma.Format(
			err,
			"unable to read file by path: %s",
			firstFilePath,
		)
	}

	var data VMI_FILE
	err = json.Unmarshal(fileBytes, &data)
	if err != nil {
		return karma.Format(
			err,
			"unable to unmarshal file by path: %s",
			firstFilePath,
		)
	}
	log.Warning("data before", data)

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

	// Write back to file
	// err = ioutil.WriteFile(fileName, editedFileData, 0644)
	// if err != nil {
	// 	return karma.Format(
	// 		err,
	// 		"unable to write file by path: %s",
	// 		firstFilePath,
	// 	)
	// }

	log.Warning("data after", data)

	// creating directory
	pathToResultDir := RESULT_DIR + data.Group
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

	// copy files to result directory
	pathToResultFileVMI := pathToResultDir + "/" + fileName
	err = ioutil.WriteFile(pathToResultFileVMI, editedFileData, 0644)
	if err != nil {
		return karma.Format(
			err,
			"unable to write file by path: %s",
			pathToResultFileVMI,
		)
	}
	// copy .vmb files to result directory
	pathToCurrentFileVMB := strings.Replace(firstFilePath, ".vmi", ".vmb", -1)
	fileVMB, err := ioutil.ReadFile(pathToCurrentFileVMB)
	if err != nil {
		return karma.Format(
			err,
			"unable to read file by path: %s",
			pathToCurrentFileVMB,
		)
	}

	pathToResultFileVMB := pathToResultDir + "/" + strings.Replace(fileName, ".vmi", ".vmb", -1)
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

// pathToResultFileVMB := pathToResultDir + "/" + strings.Replace(fileName, ".vmi", ".vmb", -1)
// err = ioutil.WriteFile(pathToResultFileVMB, editedFileData, 0644)
// if err != nil {
// 	return karma.Format(
// 		err,
// 		"unable to write file by path: %s",
// 		pathToResultFileVMB,
// 	)
// }

// lines := strings.Split(string(fileBytes), "\n")
// idContentLine := "\"id\" : " + "\""
// for i, line := range lines {
// 	log.Warning("line ", line)
// 	if strings.Contains(line, idContentLine) {
// 		lines[i] = idContentLine + " : " + "\"" + data.DisplayName + "\","
// 	}
// }

// output := strings.Join(lines, "\n")
// err = ioutil.WriteFile(fileName, []byte(output), 0644)
// if err != nil {
// 	return karma.Format(
// 		err,
// 		"unable to update file by path: %s",
// 		path,
// 	)
// }

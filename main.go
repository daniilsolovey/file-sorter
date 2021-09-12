package main

import (
	"github.com/reconquest/pkg/log"
)

func main() {
	listWithFileNames, err := getFileNamesInDir()
	if err != nil {
		log.Fatal(err)
	}

	log.Warning("len(listWithFileNames) ", len(listWithFileNames))
	for _, fileName := range listWithFileNames {
		err = handleVMIFile(fileName)
		if err != nil {
			log.Fatal(err)
		}

	}

}

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
	// log.Warning("listWithFileNames ", listWithFileNames)
	err = handleVMIFile("Base Up-Down.vmi")
	if err != nil {
		log.Fatal(err)
	}

}

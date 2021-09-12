package main

import (
	"github.com/daniilsolovey/file-sorter/internal/config"
	"github.com/docopt/docopt-go"
	"github.com/reconquest/karma-go"
	"github.com/reconquest/pkg/log"
)

var version = "[manual build]"

var usage = `file-sorter

Sort .vmi and .vmb files in specified directory.

Usage:
  file-sorter [options]

Options:
  -c --config <path>                Read specified config file. [default: config.yaml]
  --debug                           Enable debug messages.
  -v --version                      Print version.
  -h --help                         Show this help.

`

func main() {
	args, err := docopt.ParseArgs(
		usage,
		nil,
		"file-sorter "+version,
	)
	if err != nil {
		log.Fatal(err)
	}

	log.Infof(
		karma.Describe("version", version),
		"file-sorter started",
	)

	if args["--debug"].(bool) {
		log.SetLevel(log.LevelDebug)
	}

	log.Infof(nil, "loading configuration file: %q", args["--config"].(string))

	config, err := config.Load(args["--config"].(string))
	if err != nil {
		log.Fatal(err)
	}

	PATH_TO_INPUT_DIR = config.PathToInputDir
	PATH_TO_RESULT_DIR = config.PathToResultDir

	listWithFileNames, err := getFileNamesInDir()
	if err != nil {
		log.Fatal(err)
	}

	var filesWithError []map[string]string
	fileWithError := make(map[string]string)
	for _, fileName := range listWithFileNames {
		err = handleVMIFile(fileName)
		if err != nil {
			fileWithError[err.Error()] = fileName
			filesWithError = append(filesWithError, fileWithError)
			log.Error(err)
		}
	}

	if len(filesWithError) != 0 {
		log.Warning("list of errors during the program operations:")
		for _, fileWithError := range filesWithError {
			for key, value := range fileWithError {
				log.Warningf(nil, "error in file: %s\n error: %s", value, key)
			}
		}
	} else {
		log.Info("program operations was finished without errors")
	}
}

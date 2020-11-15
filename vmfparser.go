package main

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"sort"
	"strings"

	"github.com/sqweek/dialog"
)

func main() {
	createDirectoryIfItDoesNotExist("results")
	createDirectoryIfItDoesNotExist("sources")

	sourceFiles := walkSourceFiles()

	if len(sourceFiles) == 0 {
		printError("Please, put .vmf file(s) in \"sources\" folder.")
		return
	}

	for _, file := range sourceFiles {
		handleFile(file)
	}
}

func createDirectoryIfItDoesNotExist(directory string) {
	if _, err := os.Stat(directory); os.IsNotExist(err) {
		os.Mkdir(directory, os.ModePerm)
	}
}

func walkSourceFiles() []string {
	var sources []string

	err := filepath.Walk("sources", func(path string, info os.FileInfo, err error) error {
		if info.IsDir() || filepath.Ext(path) != ".vmf" {
			return nil
		}

		sources = append(sources, path)

		return nil
	})

	if err != nil {
		panic(err)
	}

	return sources
}

func handleFile(path string) {
	fileName := getFileNameForPath(path)
	file, err := os.Open(path)

	if err != nil {
		panic(err)
	}

	defer file.Close()

	// Splits on newlines by default.
	scanner := bufio.NewScanner(file)

	line := 1

	var materials []string
	var models []string

	for scanner.Scan() {
		text := scanner.Text()

		if strings.Contains(text, "material") {
			materialRow := strings.Split(text, " ")
			materialPath := strings.Trim(string(materialRow[1]), "\"")

			if !arrayContainsString(materials, materialPath) {
				materials = append(materials, materialPath)
			}
		}

		if strings.Contains(text, "model") {
			modelRow := strings.Split(text, " ")
			modelPath := strings.Trim(string(modelRow[1]), "\"")

			if !arrayContainsString(models, modelPath) {
				models = append(models, modelPath)
			}
		}

		line++
	}

	sort.Strings(materials)
	sort.Strings(models)

	writeLinesToFile(materials, "results/"+fileName+"_materials.txt")
	writeLinesToFile(models, "results/"+fileName+"_models.txt")

	if err := scanner.Err(); err != nil {
		printError("Error reading vmf: " + err.Error())
	}
}

func arrayContainsString(arr []string, str string) bool {
	for _, a := range arr {
		if a == str {
			return true
		}
	}

	return false
}

func getFileNameForPath(path string) string {
	path = strings.ReplaceAll(path, "\\", "/")
	substrings := strings.Split(path, "/")

	return substrings[len(substrings)-1]
}

func writeLinesToFile(lines []string, path string) error {
	file, err := os.Create(path)

	if err != nil {
		return err
	}

	defer file.Close()

	w := bufio.NewWriter(file)

	for _, line := range lines {
		fmt.Fprintln(w, line)
	}

	return w.Flush()
}

func printError(message string) {
	// For Windows and macOS show dialog message
	if runtime.GOOS == "windows" || runtime.GOOS == "darwin" {
		dialog.Message(message).Title("Error").Error()
	} else {
		fmt.Println(message)
	}
}

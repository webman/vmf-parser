package main

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

func main() {
	sourceFiles := walkSourceFiles()

	if _, err := os.Stat("results"); os.IsNotExist(err) {
		os.Mkdir("results", os.ModePerm)
	}

	for _, file := range sourceFiles {
		handleFile(file)
	}
}

func walkSourceFiles() []string {
	if _, err := os.Stat("sources"); os.IsNotExist(err) {
		fmt.Println("Please, create a \"sources\" folder and put .vmf files there.")
		return nil
	}

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

	var materialsArray []string

	for scanner.Scan() {
		text := scanner.Text()

		if strings.Contains(text, "material") {
			materialRow := strings.Split(text, " ")
			materialPath := strings.Trim(string(materialRow[1]), "\"")

			if !arrayContainsString(materialsArray, materialPath) {
				materialsArray = append(materialsArray, materialPath)
			}
		}

		line++
	}

	sort.Strings(materialsArray)
	writeLinesToFile(materialsArray, "results/"+fileName+"_materials.txt")

	if err := scanner.Err(); err != nil {
		fmt.Println("Error reading vmf: ", err)
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

package main

import (
	"fmt"
	"os"

	modifier "go-reloaded/modif"
)

func main() {
	if len(os.Args) != 3 {
		fmt.Fprintf(os.Stderr, "Использование: %s <входной_файл> <выходной_файл>\n", os.Args[0])
		os.Exit(1)
	}

	inputFile := os.Args[1]
	outputFile := os.Args[2]

	data, err := os.ReadFile(inputFile)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Ошибка чтения файла: %v\n", err)
		os.Exit(1)
	}

	result := modifier.Process(string(data))

	err = os.WriteFile(outputFile, []byte(result), 0o644)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Ошибка записи файла: %v\n", err)
		os.Exit(1)
	}
}

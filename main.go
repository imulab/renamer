package main

import (
	"bufio"
	"os"
)
import "fmt"

type fileAction func(*os.File) error
type lineAction func(string) error

func main() {
	args := os.Args[1:]
	if len(args) < 2 {
		panic("need two args, first is source file names, second is target file names")
	}

	sourceArg, targetArg := args[0], args[1]
	fmt.Printf("got source arg [%s] and target arg[%s]\n", sourceArg, targetArg)

	sourceCount, targetCount := 0, 0
	sourceFileNames, targetFileNames := make([]string, 0), make([]string, 0)
	if err := doWithFile(sourceArg, func(f *os.File) error {
		return eachLineOfFile(f, func(line string) error {
			sourceCount++
			if err := makeSureFileExists(line); err != nil {
				return err
			}
			sourceFileNames = append(sourceFileNames, line)
			return nil
		})
	}); err != nil {
		panic(err)
	}

	if err := doWithFile(targetArg, func(f *os.File) error {
		return eachLineOfFile(f, func(line string) error {
			targetCount++
			targetFileNames = append(targetFileNames, line)
			return nil
		})
	}); err != nil {
		panic(err)
	}

	if sourceCount != targetCount {
		panic(fmt.Errorf("source arg and target arg must provide same number of files, " +
			"but source count is %d, target count is %d", sourceCount, targetCount))
	}

	printConfirmation(sourceFileNames, targetFileNames)

	if yes := awaitConfirmation(); !yes {
		fmt.Println("user cancelled...")
		os.Exit(0)
	}

	fmt.Println("doing work...")
	if err := renameFiles(sourceFileNames, targetFileNames); err != nil {
		panic(err)
	}

	fmt.Println("done!")
}

func doWithFile(file string, action fileAction) error {
	f, err := os.OpenFile(file, os.O_RDONLY, os.ModePerm)
	if err != nil {
		return err
	}

	if err := action(f); err != nil {
		return err
	}

	defer f.Close()
	return nil
}

func eachLineOfFile(f *os.File, action lineAction) error {
	sc := bufio.NewScanner(f)
	for sc.Scan() {
		if err := action(sc.Text()); err != nil {
			return err
		}
	}
	if err := sc.Err(); err != nil {
		return err
	}
	return nil
}

func makeSureFileExists(fileName string) error {
	if _, err := os.Stat(fileName); os.IsNotExist(err) {
		return err
	}
	return nil
}

func printConfirmation(sourceNames, targetNames []string) {
	fmt.Println("Going to make the following changes: ")
	for i := 0; i < len(sourceNames); i++ {
		fmt.Printf("%s\n\t-> %s\n", sourceNames[i], targetNames[i])
	}
}

func awaitConfirmation() bool {
	fmt.Print("(y/N) > ")
	var ans string
	if _, err := fmt.Scanln(&ans); err != nil {
		return false
	}
	return ans == "y"
}

func renameFiles(sourceNames, targetNames []string) error {
	for i := 0; i < len(sourceNames); i++ {
		if err := os.Rename(sourceNames[i], targetNames[i]); err != nil {
			return err
		}
		fmt.Printf("Renamed %s\n\tto %s\n", sourceNames[i], targetNames[i])
	}
	return nil
}
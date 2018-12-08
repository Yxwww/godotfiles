package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/user"
	"path/filepath"
	"syscall"
)

func getHomeDirectory() string {
	usr, err := user.Current()
	if err != nil {
		log.Fatal(err)
	}
	return usr.HomeDir
}

func concat(a string, b string) string {
	return a + b
}

func handleSymlinkError(err error, file string) error {
	if err != nil {
		fmt.Printf("Symlink error installing %+v: %+v. ", file, err)
		if os.IsExist(err) {
			fmt.Printf("-force to overwrite\n")
			return err
		}
	}
	return err
}

func symlink(file string, target string) error {
	err := syscall.Link(file, target)
	handleSymlinkError(err, file)
	if err != nil {
		return err
	}
	fmt.Printf("Symlink %+v => %+v", file, target)
	return err
}

func nth(str string, index int) rune {
	return []rune(str)[index]
}

func isStringStartWithDotButNotDot(str string) bool {
	return string(nth(str, 0)) == "." && str != "."
}

func mapToHomeDir(name string) string {
	return concat(getHomeDirectory(), "/"+name)
}

func install(path string, file os.FileInfo, force bool) {
	if file.IsDir() {
		fmt.Printf("Ignore dir file %+v to install, path: %+v \n", file.Name(), path)
		return
	}
	target := mapToHomeDir(file.Name())
	if force {
		syscall.Unlink(target)
		symlink(path, target)
		return
	}
	symlink(path, target)
}

func generateWalkFunc(force bool) filepath.WalkFunc {
	return func(path string, info os.FileInfo, err error) error {
		if err != nil {
			fmt.Printf("walk error: %+v \n", err)
			return err
		}

		if info.IsDir() && isStringStartWithDotButNotDot(info.Name()) {
			fmt.Printf("skipping a dir without errors: %+v \n", info.Name())
			return filepath.SkipDir
		}

		install(path, info, force)
		fmt.Printf("\n")
		return nil
	}
}

func main() {
	force := flag.Bool("force", false, "a bool")
	wordPtr := flag.String("src", "dotfiles/", "a string")

	flag.Parse()

	filepath.Walk(*wordPtr, generateWalkFunc(*force))
	fmt.Printf("Done.\n")
}

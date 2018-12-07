package main

import (
	"fmt"
	"log"
	"os"
	"os/user"
	"path/filepath"
)

func getHomeDirectory() string {
	usr, err := user.Current()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(usr.HomeDir)
	return usr.HomeDir
}

func concat(a string, b string) string {
	return a + b
}

func symlink(file string, target string) error {
	err := os.Symlink(file, target)
	if err != nil {
		fmt.Println("symlink status:", err)
	}
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

func install(path string, file os.FileInfo) {
	if file.IsDir() {
		fmt.Printf("Ignore %+v to install ", file.Name())
		return
	}
	symlink(path, mapToHomeDir(file.Name()))
}

func walkFunc(path string, info os.FileInfo, err error) error {
	if err != nil {
		fmt.Printf("errors: %+v \n", err)
		return err
	}

	if info.IsDir() && isStringStartWithDotButNotDot(info.Name()) {
		fmt.Printf("skipping a dir without errors: %+v \n", info.Name())
		return filepath.SkipDir
	}

	// if string(nth(info.Name(), 0)) == "." {
	// 	fmt.Printf("skipping a hidden file: %+v \n", info.Name())
	// 	return filepath.SkipDir
	// }

	fmt.Printf("visited file or dir: %q\n", info.Name())
	install(path, info)
	return nil
}

func main() {
	// home := getHomeDirectory()
	// file := concat(home, "/chi.txt")
	// target := concat(home, "/.go/chi-link.txt")
	// err := os.Symlink(file, target)
	filepath.Walk("dotfiles/", walkFunc)
	// fmt.Println("symlink status:", err)
}

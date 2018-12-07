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
	fmt.Println(usr.HomeDir)
	return usr.HomeDir
}

func concat(a string, b string) string {
	return a + b
}

func symlink(file string, target string) error {
	err := syscall.Link(file, target)
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

func install(path string, file os.FileInfo, force bool) {
	if file.IsDir() {
		fmt.Printf("Ignore %+v to install, path: %+v \n", file.Name(), path)
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
		fmt.Println("hit file or dir: ", info.Name())
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
		install(path, info, force)
		return nil
	}
}

func main() {
	force := flag.Bool("force", false, "a bool")

	flag.Parse()

	filepath.Walk("dotfiles/", generateWalkFunc(*force))
}

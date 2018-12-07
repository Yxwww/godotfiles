package main
import (
  "fmt"
  "os/user"
  "os"
  "log"
 )

func getHomeDirectory() string {
  usr, err := user.Current()
  if err != nil {
    log.Fatal( err )
  }
  fmt.Println( usr.HomeDir )
  return usr.HomeDir;
}


func concat(a string, b string) string {
  return a + b
}

func main() {
  home := getHomeDirectory()
  file := concat(home, "/chi.txt")
  target := concat(home, "/.go/chi-link.txt")
  err := os.Symlink(file, target)
  // fmt.Println("symlink status:", err)
}



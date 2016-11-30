package main

import (
	"github.com/bluemonk3y/WinProcessMapper2/src/mapper"
)

/**
  ---> the pain the pain - seriously go?
 1. Update the package
   go get -u github.com/bluemonk3y/WinProcessMapper2/src/mapper
 2. build: go build -o mapper.exe -v mapperMain
 3. Run it> mapper.exe
 */
func main() {

	mapper.MainGo("etc/Handle.exe")
}
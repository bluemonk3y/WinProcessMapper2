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
 4. Cross compilation targets https://golang.org/doc/install/source#environment
 */
func main() {

	var handleExe = "etc/Handle.exe"
	// curl -i -XPOST http://localhost:8086/query --data-urlencode "q=CREATE DATABASE mydb"
	var influxURL = "http://192.168.99.100:32768"

	mapper.MainGo(handleExe, influxURL, "", "")
}
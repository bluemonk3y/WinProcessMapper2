
# Build & Compile
* Update the package
   go get -u github.com/bluemonk3y/WinProcessMapper2/src/mapper
*  build: go build -o mapper.exe -v mapperMain
* Run it> mapper.exe
* Cross compilation targets https://golang.org/doc/install/source#environment

See MapperMain.go & env.bat

# Dockerise it
* Need to expose docker params so it can be layers up on the fly

See Dockerfile
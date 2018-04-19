DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"
cd $DIR

GOPATH="$(pwd)" go generate main 
GOOS=windows GOPATH="$(pwd)" go build -o binaries/infosec_windows.exe main
GOOS=darwin GOPATH="$(pwd)" go build -o binaries/infosec_macos main
GOOS=linux GOPATH="$(pwd)" go build -o binaries/infosec_linux main

# GU PRINT

GU Print is a command line utility that allows you to print to Gonzaga printers
just like with guprint.gonzaga.edu.

## Use

From your terminal:
```
$ gu print myfile
```
You will then be asked to select a printer, number of copies, and the printing will begin! The job will be automatically released by default, so no waiting is necessary.

## To Install

1. [Install golang](https://golang.org/dl/). This will install go to `/Users/myusername/go` for mac or `c:\Go` for windows.
2. Download the project from your command line.
```
$ go get https://github.com/QuantamHD/gu.git
```
3. Navigate to the gu directory.
```
$ cd go/src/github.com/quantamhd/gu
```
4. Build and install the src files.
```
$ go build
$ go install
```
5. Open your profile at `$HOME/.profile`. For bash users this might be `$HOME/.bash_profile`. Add the following lines and save changes:
```
export GOPATH=$HOME/go
export PATH=$PATH:$GOROOT/bin:$GOPATH/bin
```
6. Restart your terminal.
7. You now should be able to run this utility by calling :
```
$ gu print myfile
```

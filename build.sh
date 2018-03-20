GOOS=darwin GOARCH=amd64 go build -ldflags "-w -s"
tar zcvf httpproxy-mac64.tar.gz httpproxy

GOOS=linux GOARCH=amd64 go build -ldflags "-w -s"
tar zcvf httpproxy-linux64.tar.gz httpproxy

GOOS=windows GOARCH=amd64 go build -ldflags "-w -s"
tar zcvf httpproxy-win64.tar.gz httpproxy

GOOS=windows GOARCH=386 go build -ldflags "-w -s"
tar zcvf httpproxy-win32.tar.gz httpproxy

GOOS=linux GOARCH=arm go build -ldflags "-w -s"
tar zcvf httpproxy-linux_arm.tar.gz httpproxy


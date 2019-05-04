
linux:
	CGO_ENABLED=0 GOOS=linux  GOARCH=amd64 go build

win:
	CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build

all :
	go build


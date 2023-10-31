
build:
	CGO_ENABLED=0 go build -o SystemMonitor .

build-windows:
	CGO_ENABLED=0 GOOS=windows go build -o SystemMonitor.exe .
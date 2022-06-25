all: samples-dal

.PHONY: samples-dal
samples-dal:
	@GO111MODULE=on go build -o ./tmp/sample-dal *.go

#Makefile
OUTDIR=bin
TOOLDIR=bin/tools
build:
	[ -d $(OUTDIR) ] || mkdir -p $(OUTDIR)
	[ -d $(TOOLDIR) ] || mkdir -p $(TOOLDIR)
	go build -race -o bin cmd/overseer/ovs.go
	go build -race -o bin cmd/ovswork/ovswork.go
	go build -o bin/tools/ cmd/ovstools/ovsres/ovsres.go
	go build -o bin/tools/ cmd/ovstools/ovstask/ovstask.go
	go build -o bin/chkprg chkprg/main.go
tests:
	go test `go list ./... | grep -v ./proto` -covermode=count --coverprofile='coverage.out'
sprof:
	go tool pprof -http localhost:9090 schedulerprofile.prof
wprof:
	go tool pprof -http localhost:9090 workerprofile.prof
cover:
	go tool cover  -html coverage.out  
protoc:
	 protoc --proto_path=./proto/  --go_out=plugins=grpc:./ proto/wservices.proto 
	 protoc --proto_path=./proto/  --go_out=plugins=grpc:./ proto/actions.proto 
	 protoc --proto_path=./proto/  --go_out=plugins=grpc:./ proto/services.proto

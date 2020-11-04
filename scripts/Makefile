#Makefile
build:
	go build -race -o bin cmd/overseer/ovs.go
	go build -race -o bin cmd/ovswork/ovswork.go
	go build -o bin/tools/ cmd/ovstools/ovsres/ovsres.go
	go build -o bin/tools/ cmd/ovstools/ovstask/ovstask.go
	go build -o bin/chkprg chkprg/main.go
tests:
	go test ./... --coverpkg=./... --coverprofile='goscheduler.out'
sprof:
	go tool pprof -http localhost:9090 schedulerprofile.prof
wprof:
	go tool pprof -http localhost:9090 workerprofile.prof
cover:
	go tool cover  -html goscheduler.out  
protoc:
	 protoc --proto_path=./proto/  --go_out=plugins=grpc:./ proto/wservices.proto 
	 protoc --proto_path=./proto/  --go_out=plugins=grpc:./ proto/actions.proto 
	 protoc  --go_out=plugins=grpc:./ proto/services.proto
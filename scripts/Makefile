#Makefile
OUTDIR=bin
TOOLDIR=bin/tools
TESTDIR=data/tests
SYSDIR=data/sysout
LOGDIR=logs
build:
	[ -d $(OUTDIR) ] || mkdir -p $(OUTDIR)
	[ -d $(TOOLDIR) ] || mkdir -p $(TOOLDIR)
	[ -d $(SYSDIR) ] || mkdir -p $(SYSDIR)
	[ -d $(LOGDIR) ] || mkdir -p $(LOGDIR)
	go build -race -o bin cmd/overseer/ovs.go
	go build -race -o bin cmd/ovswork/ovswork.go
	go build -race -o bin cmd/ovsgate/ovsgate.go
	go build -o bin/tools/ cmd/ovstools/ovscli.go
	go build -o bin/chkprg chkprg/main.go
tests:
	[ -d $(TESTDIR) ] || mkdir -p $(TESTDIR)
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
	 protoc -I ./proto --grpc-gateway_out ./proto/services --grpc-gateway_opt paths=source_relative --grpc-gateway_opt logtostderr=true services.proto
	 protoc -I ./proto --openapiv2_out ./www/swagger --openapiv2_opt logtostderr=true services.proto
dockerbuild:
	docker build -t overseer/mongodb -f docker/mongodb/Dockerfile .
start:
	docker run -d -l ovsdbmongo --name mongo_ovs -p23017:27017 -v ${PWD}/docker/mongodb/etc:/etc/mongo \
		-e MONGO_INITDB_ROOT_USERNAME=admin \
		-e MONGO_INITDB_ROOT_PASSWORD=notsecure \
		-e MONGO_INITDB_DATABASE=overseer \
		overseer/mongodb --config /etc/mongo/config.yml
stop:
	docker stop mongo_ovs
remove:
	docker rm -f $$( docker ps -qaf "label=ovsdbmongo")

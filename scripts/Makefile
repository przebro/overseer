#Makefile
OUTDIR=bin
TOOLDIR=bin/tools
TESTDIR=data/tests
SYSDIR=data/sysout
POOLDIR=def/.pool
POOLDIRT=def_test/.pool
LOGDIR=logs
build:
	[ -d $(OUTDIR) ] || mkdir -p $(OUTDIR)
	[ -d $(TOOLDIR) ] || mkdir -p $(TOOLDIR)
	[ -d $(SYSDIR) ] || mkdir -p $(SYSDIR)
	[ -d $(LOGDIR) ] || mkdir -p $(LOGDIR)
	[ -d $(POOLDIR) ] || mkdir -p $(POOLDIR)
	go build -race -o bin cmd/overseer/ovs.go
	go build -race -o bin cmd/ovswork/ovswork.go
	go build -race -o bin cmd/ovsgate/ovsgate.go
	go build -o bin/tools/ cmd/ovstools/ovscli.go
	go build -o bin/chkprg chkprg/main.go
arch:
	GOOS=linux GOARCH=amd64 go build -o bin/ovs-linux-amd64 cmd/overseer/ovs.go
	GOOS=linux GOARCH=amd64 go build -o bin/ovsgate-linux-amd64 cmd/ovsgate/ovsgate.go
	GOOS=linux GOARCH=amd64 go build -o bin/ovswork-linux-amd64 cmd/ovswork/ovswork.go
tests:
	[ -d $(TESTDIR) ] || mkdir -p $(TESTDIR)
	[ -d $(POOLDIRT) ] || mkdir -p $(POOLDIRT)
	go test `go list ./... | grep -v ./proto` -covermode=count --coverprofile='coverage.out'
sprof:
	go tool pprof -http localhost:9090 schedulerprofile.prof
wprof:
	go tool pprof -http localhost:9090 workerprofile.prof
cover:
	go tool cover  -html coverage.out  
protoc:
	 protoc --proto_path=./proto/ --go_out=paths=source_relative:./proto/wservices proto/wservices.proto  --go-grpc_out=./ proto/wservices.proto 
	 protoc --proto_path=./proto/ --go_out=paths=source_relative:./proto/actions proto/actions.proto   --go-grpc_out=./ proto/actions.proto 
	 protoc --proto_path=./proto/ --go_out=paths=source_relative:./proto/services proto/services.proto  --go-grpc_out=./ proto/services.proto
	 protoc -I ./proto --grpc-gateway_out ./proto/services --grpc-gateway_opt paths=source_relative --grpc-gateway_opt logtostderr=true services.proto
	 protoc -I ./proto --openapiv2_out ./www/swagger --openapiv2_opt logtostderr=true services.proto
dockerbuild:
	docker build -t overseer/mongodb -f docker/mongodb/Dockerfile .
	docker build -t ovswork:latest . -f Dockerfile.worker
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


$targetname = $args[0]

$targets = @(
    "build","tests","cover","protoc"
) 

if (!$targets.Contains($targetname)){
    Write-Host "Invalid target name:" $targetname
    exit(1)
}

if ((Test-Path -Path ".\bin" -PathType Container) -eq $False){
    mkdir "bin"
    mkdir "bin/tools"
}
if ((Test-Path -Path ".\data" -PathType Container) -eq $False){
    mkdir "data"
    mkdir "data/sysout"
}
if ((Test-Path -Path ".\logs" -PathType Container) -eq $False){
    mkdir "logs"
}

if ($targetname -eq "build"){
    go build -race -o bin/ cmd/overseer/ovs.go
	go build -race -o bin/ cmd/ovswork/ovswork.go
    go build -o bin/tools/ cmd/ovstools/ovscli.go
	go build -o bin/chkprg chkprg/main.go
}
if ($targetname -eq "tests"){
    go test ./... --coverpkg=./... --coverprofile='overseer.out'
}

if ($targetname -eq "cover"){
    go tool cover  -html overseer.out  
}
if ($targetname -eq "protoc"){
    protoc --proto_path=./proto/  --go_out=plugins=grpc:./ proto/wservices.proto 
    protoc --proto_path=./proto/  --go_out=plugins=grpc:./ proto/actions.proto 
    protoc  --go_out=plugins=grpc:./ proto/services.proto
}

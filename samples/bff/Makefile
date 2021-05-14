kratosThirdPartyDir:=github.com/go-kratos/kratos@v1.0.0/third_party

moduleName:=$(shell head -n 1 go.mod | awk '{print $$2}')
GOMODCACHE:=$(shell go env GOMODCACHE)

# protoc
protocImport:=-I$(GOMODCACHE)/../proto -I$(GOMODCACHE)/$(kratosThirdPartyDir) -I./:$(GOMODCACHE) -Iapi
protocOut:=--gofast_out=plugins=grpc,paths=source_relative:. --grpc-gateway_out=./

f:=api/*.proto


run:
	bin/app -env dev
proto:
	#protoc $(protocImport) $(protocOut) $(f)
	#protoc $(protocImport)--go_out=paths=source_relative:.   --go-grpc_out=paths=source_relative:.  --grpc-gateway_out=paths=source_relative:. --openapiv2_out=.  api/*.proto
	#protoc $(protocImport) --gotag_out=paths=source_relative:.  api/*.proto

	protoc $(protocImport) --go_out=paths=source_relative:.   --go-grpc_out=paths=source_relative:.  --grpc-gateway_out=paths=source_relative:. --openapiv2_out=.  api/*.proto
	protoc $(protocImport) --gotag_out=paths=source_relative:.  api/*.proto

build:
	go build -o ./bin/app $(moduleName)/cmd/app
wire:
	wire gen $(moduleName)/cmd/app

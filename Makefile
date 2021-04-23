GOPATH:=$(shell go env GOPATH)
VERSION=$(shell git describe --tags --always)
PROTO_FILES=$(shell find . -name *.proto)
KRATOS_VERSION=$(shell go mod graph |grep go-kratos/kratos/v2 |head -n 1 |awk -F '@' '{print $$2}')
KRATOS=$(GOPATH)/pkg/mod/github.com/go-kratos/kratos/v2@$(KRATOS_VERSION)
VALIDATE_VERSION=$(shell ls $(GOPATH)/pkg/mod/github.com/envoyproxy/|grep protoc-gen-validate|head -n 1)
VALIDATE=$(GOPATH)/pkg/mod/github.com/envoyproxy/$(VALIDATE_VERSION)
SWAG_VERSION=$(shell ls $(GOPATH)/pkg/mod/github.com/grpc-ecosystem/ |grep grpc-gateway|head -n 1)
SWAG=$(GOPATH)/pkg/mod/github.com/grpc-ecosystem/$(SWAG_VERSION)/third_party/googleapis/google/api
SWG_PROTO_FILES=$(shell find . -name *.proto|grep api|grep -v error);
MAIN_PATH=$(shell find . -name main.go|sed "s/main.go//")


.PHONY: init
# init env
init:
	go get -u github.com/go-kratos/kratos/cmd/kratos/v2
	go get -u github.com/go-kratos/kratos/cmd/protoc-gen-go-http/v2
	go get -u github.com/go-kratos/kratos/cmd/protoc-gen-go-errors/v2
	go get -u google.golang.org/protobuf/cmd/protoc-gen-go
	go get -u google.golang.org/grpc/cmd/protoc-gen-go-grpc
	go get -u github.com/envoyproxy/protoc-gen-validate
	go get -u github.com/grpc-ecosystem/grpc-gateway/protoc-gen-swagger
	go get -u entgo.io/ent/cmd/ent
	go get -u github.com/jteeuwen/go-bindata/...
	go get -u github.com/google/wire/cmd/wire

.PHONY: proto
# generate rpc code
proto:
	protoc --proto_path=. \
           --proto_path=$(KRATOS)/api \
           --proto_path=$(KRATOS)/third_party \
           --proto_path=$(GOPATH)/src \
           --proto_path=$(VALIDATE) \
           --validate_out="lang=go",paths=source_relative:. \
           --go_out=paths=source_relative:. \
           --go-grpc_out=paths=source_relative:. \
           --go-errors_out=paths=source_relative:. \
           $(PROTO_FILES)

.PHONY: http
# generate http code
http:
	protoc --proto_path=. \
           --proto_path=$(KRATOS)/api \
           --proto_path=$(KRATOS)/third_party \
           --proto_path=$(GOPATH)/src \
           --proto_path=$(VALIDATE) \
           --validate_out="lang=go",paths=source_relative:. \
           --go_out=paths=source_relative:. \
           --go-http_out=paths=source_relative:. \
           --go-errors_out=paths=source_relative:. \
           $(PROTO_FILES)

.PHONY: swag
# generate swagger
swag:
	rm -rf docs;
	protoc --proto_path=. \
		   --proto_path=$(KRATOS)/api \
		   --proto_path=$(KRATOS)/third_party \
		   --proto_path=$(GOPATH)/src \
		   --proto_path=$(VALIDATE) \
		   --proto_path=$(SWAG) \
		   --validate_out="lang=go",paths=source_relative:. \
		   --go_out=paths=source_relative:. \
		   --swagger_out=allow_delete_body=true,logtostderr=true,allow_merge=true:. \
		   --go-grpc_out=paths=source_relative:. \
		   --go-http_out=paths=source_relative:. \
		   --go-errors_out=paths=source_relative:. \
		   $(SWG_PROTO_FILES)
	mkdir docs;
	mv *.swagger.json docs;
#todo
#https://github.com/go-swagger/go-swagger
#swagger serve apidocs.swagger.json -Fswagger --no-open --port 8000

.PHONY: ent
# generate ent
ent:
	ent generate ./internal/data/ent/schema --target ./internal/data/ent

.PHONY: generate
# generate client code
generate:
	go generate ./...

.PHONY: all
# generate client code
all:
	make proto;
	make http;
	make swag;
	make ent;
	make generate

.PHONY: run
# run program
run:
	cd $(MAIN_PATH) && go run .

.PHONY: build
# build
build:
	mkdir -p bin/ && go build -ldflags "-X main.Version=$(VERSION)" -o ./bin/ ./...

.PHONY: test
# test
test:
	go test -v ./... -cover

.PHONY: clean
# clean generate code
clean:
	find . -name *.*.go -o -name *.swagger.json |xargs rm -rf;
	cd internal/data/ent && ls |grep -v schema|xargs rm -rf

# show help
help:
	@echo ''
	@echo 'Usage:'
	@echo ' make [target]'
	@echo ''
	@echo 'Targets:'
	@awk '/^[a-zA-Z\-\_0-9]+:/ { \
	helpMessage = match(lastLine, /^# (.*)/); \
		if (helpMessage) { \
			helpCommand = substr($$1, 0, index($$1, ":")-1); \
			helpMessage = substr(lastLine, RSTART + 2, RLENGTH); \
			printf "\033[36m%-22s\033[0m %s\n", helpCommand,helpMessage; \
		} \
	} \
	{ lastLine = $$0 }' $(MAKEFILE_LIST)

.DEFAULT_GOAL := help

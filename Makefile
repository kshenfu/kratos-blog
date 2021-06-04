GOPATH:=$(shell go env GOPATH)
VERSION=$(shell git describe --tags --always)
PROTO_FILES=$(shell find . -name *.proto)
API_PROTO_FILES=$(shell find api -name *.proto)
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
	go get -u entgo.io/ent/cmd/ent
	go get -u github.com/google/wire/cmd/wire
	go get -u github.com/jteeuwen/go-bindata/...
	go get -u github.com/go-kratos/kratos/cmd/kratos/v2
	go get -u github.com/go-kratos/kratos/cmd/protoc-gen-go-http/v2
	go get -u google.golang.org/protobuf/cmd/protoc-gen-go
	go get -u google.golang.org/grpc/cmd/protoc-gen-go-grpc
	go get -u github.com/envoyproxy/protoc-gen-validate
	go get -u github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2

.PHONY: ent
# generate ent
ent:
	ent generate ./internal/data/ent/schema --target ./internal/data/ent

.PHONY: proto
# generate internal proto
proto:
	protoc --proto_path=. \
           --proto_path=$(KRATOS)/third_party \
           --go_out=paths=source_relative:. \
           $(PROTO_FILES)

.PHONY: grpc
# generate api grpc code
grpc:
	protoc --proto_path=. \
           --proto_path=$(KRATOS)/third_party \
           --go-grpc_out=paths=source_relative:. \
           $(API_PROTO_FILES)

.PHONY: http
# generate api http code
http:
	protoc --proto_path=. \
           --proto_path=$(KRATOS)/third_party \
           --go-http_out=paths=source_relative:. \
           $(API_PROTO_FILES)

.PHONY: validate
# generate api validate code
validate:
	protoc --proto_path=. \
           --proto_path=$(KRATOS)/third_party \
           --validate_out="lang=go",paths=source_relative:. \
           $(API_PROTO_FILES)

.PHONY: swagger
# generate api swagger file
swagger:
	protoc --proto_path=. \
		--proto_path=$(KRATOS)/third_party \
		--openapiv2_out . \
		--openapiv2_opt logtostderr=true \
		--openapiv2_opt json_names_for_fields=false \
		$(API_PROTO_FILES)

.PHONY: all
# generate all
all: clean ent proto grpc http validate swagger

.PHONY: clean
# clean generate code
clean:
	find . -name *.*.go -o -name *.swagger.json |xargs rm -rf;
	cd internal/data/ent && ls |grep -v schema|xargs rm -rf

.PHONY: debug
# go run program
debug:
	cd $(MAIN_PATH) && go run .

.PHONY: build
# build
build:
	mkdir -p bin/ && go build -ldflags "-X main.Version=$(VERSION)" -o ./bin/ ./...

.PHONY: run
# run
run: build
	./bin/blog -conf ./configs/config.yaml

.PHONY: test
# test
test:
	go test -v ./... -cover

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

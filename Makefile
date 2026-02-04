default: build

.PHONY: build
build:
	go mod download
	go mod tidy
	go build -o terraform-provider-opnsense26

.PHONY: install
install: build
	mkdir -p ~/.terraform.d/plugins/registry.terraform.io/rgcosta7/opnsense26/0.1.0/linux_amd64
	cp terraform-provider-opnsense26 ~/.terraform.d/plugins/registry.terraform.io/rgcosta7/opnsense26/0.1.0/linux_amd64/

.PHONY: test
test:
	go test -v ./...

.PHONY: testacc
testacc:
	TF_ACC=1 go test -v ./... -timeout 120m

.PHONY: fmt
fmt:
	go fmt ./...
	terraform fmt -recursive ./examples/

.PHONY: lint
lint:
	golangci-lint run

.PHONY: docs
docs:
	go generate ./...

.PHONY: clean
clean:
	rm -f terraform-provider-opnsense26
	rm -rf dist/

.PHONY: mod
mod:
	go mod tidy
	go mod download

REGISTRY ?= quay.io/kubevirt
IMAGE_TAG ?= latest
IMAGE_GIT_TAG ?= $(shell git describe --abbrev=8 --tags)

COMPONENTS = $(sort \
			 $(subst /,-,\
			 $(patsubst cmd/%/,%,\
			 $(dir \
			 $(shell find cmd/ -type f -name '*.go')))))

BIN_DIR = $(CURDIR)/build/_output/bin/
export GOROOT=$(BIN_DIR)/go/
export GOBIN = $(GOROOT)/bin/
export PROTOCROOT = $(BIN_DIR)/protoc/
export PROTOCBIN = $(PROTOCROOT)/bin/
export PATH := $(GOBIN):$(PATH):$(BIN_DIR)
GOPATH = $(CURDIR)/.gopath
ORG_PATH = github.com/k8snetworkplumbingwg
PACKAGE = ovs-cni
REPO_PATH = $(ORG_PATH)/$(PACKAGE)
BASE = $(GOPATH)/src/$(REPO_PATH)
PKGS = $(or $(PKG),$(shell cd $(BASE) && env GOPATH=$(GOPATH) $(GO) list ./... | grep -v "$(PACKAGE)/vendor/" | grep -v "$(PACKAGE)/tests/cluster" | grep -v "$(PACKAGE)/tests/node"))
V = 0
Q = $(if $(filter 1,$V),,@)

all: lint build

GO := $(GOBIN)/go
PROTOC := $(PROTOCBIN)/protoc

$(GO):
	   hack/install-go.sh $(BIN_DIR)

$(PROTOC):
		hack/install-protoc.sh $(BIN_DIR)/protoc

$(BASE): ; $(info  setting GOPATH...)
	@mkdir -p $(dir $@)
	@ln -sf $(CURDIR) $@

GOLINT = $(GOBIN)/golint
$(GOBIN)/golint: | $(BASE) ; $(info  building golint...)
	$Q go install -mod=mod golang.org/x/lint/golint@v0.0.0-20210508222113-6edffad5e616

gen-grpc: $(GO) $(PROTOC)
		$(GO) install -mod=mod google.golang.org/protobuf/cmd/protoc-gen-go@v1.28
		$(GO) install -mod=mod google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.2
		$(PROTOC) --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative pkg/L2IPUInfraManager/l2_ipu_infra_mgr.proto

build: gen-grpc format $(patsubst %, build-%, $(COMPONENTS))

lint: | gen-grpc $(GO) $(BASE) $(GOLINT) ; $(info  running golint...) @ ## Run golint
	$Q cd $(BASE) && ret=0 && for pkg in $(PKGS); do \
		test -z "$$($(GOLINT) $$pkg | tee /dev/stderr)" || ret=1 ; \
	 done ; exit $$ret

build-%: $(GO)
	hack/version.sh > ./cmd/$(subst -,/,$*)/.version
	cd cmd/$(subst -,/,$*) && $(GO) fmt && $(GO) vet && GOOS=linux GOARCH=amd64 CGO_ENABLED=0 GO111MODULE=on $(GO) build -tags no_openssl -mod vendor

format: $(GO)
	$(GO) fmt ./pkg/... ./cmd/...
	$(GO) vet ./pkg/... ./cmd/...

docker-build: $(patsubst %, docker-build-%, $(COMPONENTS))

docker-build-%: build-%
	docker build -t ${REGISTRY}/ovs-cni-$*:${IMAGE_TAG} ./cmd/$(subst -,/,$*)

docker-push: $(patsubst %, docker-push-%, $(COMPONENTS))

docker-push-%:
	docker push ${REGISTRY}/ovs-cni-$*:${IMAGE_TAG}
	docker tag ${REGISTRY}/ovs-cni-$*:${IMAGE_TAG} ${REGISTRY}/ovs-cni-$*:${IMAGE_GIT_TAG}
	docker push ${REGISTRY}/ovs-cni-$*:${IMAGE_GIT_TAG}

.PHONY: build format docker-build docker-push lint gen-grpc
DOCKERHUB_USER := kaizoa
GITHUB_USER := kai-zoa

NAME       := d4aws
REPOSITORY := github.com/$(GITHUB_USER)/d4aws
REGISTRY   := $(DOCKERHUB_USER)/$(NAME)
VERSION    := 0.0.1
REVISION   := $(shell git rev-parse --short HEAD)
CODES      := $(shell find . -type f -name '*.go')
GO_VERSION := 1.8
LDFLAGS    := -ldflags="-s -w -X \"main.Version=$(VERSION)\""
BUILD_OPTS :=
RUN_ARGS   :=
#LDFLAGS    := -ldflags="-s -w -X \"main.Version=$(VERSION)\" -X \"main.Revision=$(REVISION)\" -extldflags \"-static\""

.PHONY: dep docker run release

dep:
ifeq ($(shell type dep 2> /dev/null),)
	go get github.com/golang/dep
endif
	dep ensure

$(NAME): dep
	docker run \
 -v $(GOPATH):/go \
 -w /go/src/$(REPOSITORY) \
 -it \
 -e GOOS=linux \
 -e GOARCH=amd64 \
 --entrypoint go \
 golang:$(GO_VERSION)-alpine \
 build $(BUILD_OPTS) $(LDFLAGS) -o $(NAME) ./main.go

docker: $(NAME)
	docker build -t $(NAME):$(VERSION) .
	docker tag $(NAME):$(VERSION) $(REGISTRY):$(VERSION)
	docker tag $(NAME):$(VERSION) $(REGISTRY):latest

run:
	docker run -v $(HOME)/.aws:/home/d4aws/.aws $(NAME):$(VERSION) $(RUN_ARGS)

release: docker
	docker push $(REGISTRY):$(VERSION)
	docker push $(REGISTRY):latest

.PHONY: clean
clean:
	rm -rf $(NAME)
	dep remove

.PHONY: test
test:
	go test -cover -v ./service/...


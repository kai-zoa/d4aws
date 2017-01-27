NAME      := d4aws
GITUB_REPOS   := github.com/kai-zoa/d4aws
DOCKERHUB_REPOS := kaizoa/d4aws
VERSION   := 0.0.1
REVISION  := $(shell git rev-parse --short HEAD)
CODES     := $(shell find . -type f -name '*.go')
LDFLAGS   := -ldflags="-s -w -X \"main.Version=$(VERSION)\" -X \"main.Revision=$(REVISION)\" -extldflags \"-static\""

.PHONY: _dep
_dep:
ifeq ($(shell type dep 2> /dev/null),)
	go get github.com/golang/dep
endif

.PHONY: dep
dep: _dep
	dep ensure

$(NAME): dep
	docker run \
 -v `pwd`:/go/src/$(GITUB_REPOS) \
 -w /go/src/$(GITUB_REPOS) \
 -it \
 golang:onbuild \
 sh -c 'CGO_ENABLED=0 go build -a -tags netgo -installsuffix netgo $(LDFLAGS) -o $(NAME) ./main.go'

.PHONY: docker-build
docker-build: $(NAME)
	docker build -t $(NAME):$(VERSION) .

docker-tag: docker-build
	docker tag $(NAME):$(VERSION) $(DOCKERHUB_REPOS):$(VERSION)
	docker tag $(NAME):$(VERSION) $(DOCKERHUB_REPOS):latest

docker-push: docker-tag
	docker push $(DOCKERHUB_REPOS):$(VERSION)
	docker push $(DOCKERHUB_REPOS):latest

.PHONY: clean
clean:
	rm -rf $(NAME)
	dep remove

.PHONY: test
test:
	go test -cover -v ./service/...

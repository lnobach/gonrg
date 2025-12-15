.PHONY : test unittest clean


GOARCH = amd64
GONRG_VERSION = $(shell awk -v FS="gonrg=" 'NF>1{print $$2}' VERSIONS)
GO_LDFLAGS := "\
	-X 'github.com/lnobach/gonrg/version.GonrgVersion=$(GONRG_VERSION)' \
	-extldflags '-static' -w -s"
GOENV := GOARCH=$(GOARCH) GOOS=linux
GO_FILES := $(wildcard **/*.go)
ALLSRC_FILES := go.mod go.sum VERSIONS $(GO_FILES)

all: gonrg

gonrg: $(ALLSRC_FILES)
	$(GOENV) go build -v -ldflags $(GO_LDFLAGS) -o gonrg ./cmd/gonrg/.

gonrg-mock: $(ALLSRC_FILES)
	echo Foo $(ALLSRC_FILES)
	$(GOENV) go build -v -ldflags $(GO_LDFLAGS) -tags gonrgmocks -o gonrg-mock ./cmd/gonrg/.

test: unittest

unittest:
	$(GOENV) go test ./... -coverprofile coverage.out
	$(GOENV) go tool cover -func=coverage.out

clean:
	rm -f gonrg gonrg-mock

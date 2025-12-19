.PHONY : test unittest clean

PLATFORMS := linux-amd64 windows-amd64.exe darwin-amd64 darwin-arm64 linux-arm-6 linux-arm64-8 linux-mips

releasebins := $(addprefix release/gonrg-, $(PLATFORMS))
temp = $(subst -, ,$@)
os = $(word 2, $(temp))
arch = $(word 1, $(subst .,  ,$(word 3, $(temp))))
version = $(word 4, $(temp))

GONRG_VERSION = $(shell awk -v FS="gonrg=" 'NF>1{print $$2}' VERSIONS)
GO_LDFLAGS := "\
	-X 'github.com/lnobach/gonrg/version.GonrgVersion=$(GONRG_VERSION)' \
	-extldflags '-static' -w -s"
GOENV := 
GO_FILES := $(wildcard **/*.go)
ALLSRC_FILES := go.mod go.sum VERSIONS $(GO_FILES)

all: gonrg

gonrg: $(ALLSRC_FILES)
	$(GOENV) go build -v -ldflags $(GO_LDFLAGS) -o gonrg ./cmd/gonrg/.

gonrg-mock: $(ALLSRC_FILES)
	$(GOENV) go build -v -ldflags $(GO_LDFLAGS) -tags gonrgmocks -o gonrg-mock ./cmd/gonrg/.

release: $(releasebins)

$(releasebins): $(ALLSRC_FILES)
	GOOS=$(os) GOARCH=$(arch) GOARM=$(version) go build -v -ldflags $(GO_LDFLAGS) -o '$@' ./cmd/gonrg/.

test: unittest

unittest:
	$(GOENV) go test ./... -coverprofile coverage.out
	$(GOENV) go tool cover -func=coverage.out

clean:
	rm -f gonrg gonrg-mock release

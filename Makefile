.PHONY : test unittest clean

PLATFORMS := linux-amd64 windows-amd64.exe darwin-amd64 darwin-arm64 linux-arm linux-arm64 linux-mips linux-mipsle
releasebins := $(addprefix release/gonrg-, $(PLATFORMS))

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

release/gonrg-linux-amd64: $(ALLSRC_FILES)
	GOOS=linux GOARCH=amd64 go build -v -ldflags $(GO_LDFLAGS) -o '$@' ./cmd/gonrg/.

release/gonrg-windows-amd64.exe: $(ALLSRC_FILES)
	GOOS=windows GOARCH=amd64 go build -v -ldflags $(GO_LDFLAGS) -o '$@' ./cmd/gonrg/.

release/gonrg-darwin-amd64: $(ALLSRC_FILES)
	GOOS=darwin GOARCH=amd64 go build -v -ldflags $(GO_LDFLAGS) -o '$@' ./cmd/gonrg/.

release/gonrg-darwin-arm64: $(ALLSRC_FILES)
	GOOS=darwin GOARCH=arm64 GOARM=8 go build -v -ldflags $(GO_LDFLAGS) -o '$@' ./cmd/gonrg/.

release/gonrg-linux-arm: $(ALLSRC_FILES)
	GOOS=linux GOARCH=arm GOARM=6 go build -v -ldflags $(GO_LDFLAGS) -o '$@' ./cmd/gonrg/.

release/gonrg-linux-arm64: $(ALLSRC_FILES)
	GOOS=linux GOARCH=arm64 GOARM=8 go build -v -ldflags $(GO_LDFLAGS) -o '$@' ./cmd/gonrg/.

release/gonrg-linux-mips: $(ALLSRC_FILES)
	GOOS=linux GOARCH=mips GOMIPS=softfloat go build -v -ldflags $(GO_LDFLAGS) -o '$@' ./cmd/gonrg/.

release/gonrg-linux-mipsle: $(ALLSRC_FILES)
	GOOS=linux GOARCH=mipsle GOMIPS=softfloat go build -v -ldflags $(GO_LDFLAGS) -o '$@' ./cmd/gonrg/.

releasegz: $(releasebins)
	gzip -v9f $^

test: unittest

unittest:
	$(GOENV) go test ./... -coverprofile coverage.out
	$(GOENV) go tool cover -func=coverage.out

clean:
	rm -f gonrg gonrg-mock release

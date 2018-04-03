VERSION ?= 0.1.0

BUILD_DIR ?= ./out
APPS := systembolaget-beer-releases

all: clean deps vet fmt test build

ci: vet deps cover

deps:
	dep ensure

fmt:
	go fmt ./...

vet:
	go vet ./...

test:
	go test -v ./...

cover:
	go test -v ./... -covermode=count -coverprofile=coverage.out

build:
	for target in $(APPS); do \
		$(BUILD_ENV_FLAGS) go build -o $(BUILD_DIR)/$$target ./cmd/$$target; \
	done

clean:
	rm -rf $(BUILD_DIR)

dockerbuild:
	./scripts/dockerbuild.sh $(VERSION)

.PHONY: build

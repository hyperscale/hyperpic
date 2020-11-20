BUILD_DIR ?= build
COMMIT = $(shell git rev-parse HEAD)
VERSION ?= $(shell git describe --match 'v[0-9]*' --dirty='-dev' --always)
ORG := hyperscale
PROJECT := hyperpic
REPOPATH ?= github.com/$(ORG)/$(PROJECT)
VERSION_PACKAGE = $(REPOPATH)/pkg/$(PROJECT)/version
IMAGE ?= $(ORG)/$(PROJECT)

GO_LDFLAGS :="
GO_LDFLAGS += -X $(VERSION_PACKAGE).Version=$(VERSION)
GO_LDFLAGS += -X $(VERSION_PACKAGE).BuildAt=$(shell date +'%Y-%m-%dT%H:%M:%SZ')
GO_LDFLAGS += -X $(VERSION_PACKAGE).Revision=$(COMMIT)
GO_LDFLAGS +="

GO_FILES := $(shell find . -type f -name '*.go' -not -path "./vendor/*")

HYPERPIC_AUTH_SECRET ?= c8da8ded-f9a2-429c-8811-9b2a07de8ede

ifeq ($(shell uname -s), Darwin)
	CGO_CFLAGS_ALLOW ?= "-Xpreprocessor"
	PKG_CONFIG_PATH ?= "/usr/local/opt/libffi/lib/pkgconfig:/usr/local/opt/vips/lib/pkgconfig"
endif

.PHONY: release
release:
	@echo "Release v$(version)"
	@git pull
	@git checkout master
	@git pull
	@git checkout develop
	@git flow release start $(version)
	@echo "$(version)" > .version
	@sed -e "s/version: .*/version: \"v$(version)\"/g" docs/swagger.yaml > docs/swagger.yaml.new && rm -rf docs/swagger.yaml && mv docs/swagger.yaml.new docs/swagger.yaml
	@git add .version docs/swagger.yaml
	@git commit -m "feat(project): update version file" .version docs/swagger.yaml
	@git flow release finish $(version)
	@git push
	@git push --tags
	@git checkout master
	@git push
	@git checkout develop
	@echo "Release v$(version) finished."

.PHONY: all
all: deps build test

.PHONY: deps
deps:
	@go mod download

.PHONY: clean
clean:
	@go clean -i ./...

create-build-dir:
	@mkdir -p $(BUILD_DIR)

$(BUILD_DIR)/coverage.out: $(GO_FILES) create-build-dir go.mod go.sum
	@CGO_CFLAGS_ALLOW=-Xpreprocessor go test -race -cover -coverprofile $(BUILD_DIR)/coverage.out.tmp ./...
	@cat $(BUILD_DIR)/coverage.out.tmp | grep -v '.pb.go' | grep -v 'mock_' > $(BUILD_DIR)/coverage.out
	@rm $(BUILD_DIR)/coverage.out.tmp

.PHONY: ci-test
ci-test:
	@go test -race -cover -coverprofile ./coverage.out.tmp -v ./... | go2xunit -fail -output tests.xml
	@cat ./coverage.out.tmp | grep -v '.pb.go' | grep -v 'mock_' > ./coverage.out
	@rm ./coverage.out.tmp
	@echo ""
	@go tool cover -func ./coverage.out

.PHONY: lint
lint:
	@golangci-lint run ./...

.PHONY: test
test: $(BUILD_DIR)/coverage.out

.PHONY: coverage
coverage: $(BUILD_DIR)/coverage.out
	@echo ""
	@go tool cover -func ./$(BUILD_DIR)/coverage.out

.PHONY: coverage-html
coverage-html: $(BUILD_DIR)/coverage.out
	@go tool cover -html ./$(BUILD_DIR)/coverage.out

generate:
	@go generate ./...

cmd/hyperpic/app/asset/bindata.go: docs/index.html docs/swagger.yaml
	@echo "Bin data..."
	@go-bindata -pkg asset -o cmd/hyperpic/app/asset/bindata.go docs/

${BUILD_DIR}/hyperpic: $(GO_FILES) cmd/hyperpic/app/asset/bindata.go go.mod go.sum
	@echo "Building $@..."
	@CGO_CFLAGS_ALLOW=-Xpreprocessor go generate ./cmd/$(subst ${BUILD_DIR}/,,$@)/
	@CGO_CFLAGS_ALLOW=-Xpreprocessor go build -ldflags $(GO_LDFLAGS) -o $@ ./cmd/$(subst ${BUILD_DIR}/,,$@)/

.PHONY: run-hyperpic
run-hyperpic: ${BUILD_DIR}/hyperpic
	@echo "Running $<..."
	@./$< --config=./cmd/hyperpic/config.yml

.PHONY: run
run: run-hyperpic

.PHONY: run-docker
run-docker: docker
	@docker run -e "HYPERPIC_AUTH_SECRET=$(HYPERPIC_AUTH_SECRET)" -p 8574:8080 -v $(shell pwd)/var/lib/hyperpic:/var/lib/hyperpic --rm $(IMAGE)

.PHONY: build
build: ${BUILD_DIR}/hyperpic

.PHONY: docker
docker:
	@docker build -f cmd/hyperpic/Dockerfile --rm -t $(IMAGE) .

.PHONY: publish
publish: docker
	@docker tag $(IMAGE) $(IMAGE):latest
	@docker push $(IMAGE)

.PHONY: publish-dev
publish-dev: docker
	@docker tag $(IMAGE) $(IMAGE):dev
	@docker push $(IMAGE):dev

heroku: docker
	@echo "Deploy Hyperpic on Heroku..."
	@#heroku container:push web --app=hyperpic
	@docker tag $(IMAGE) registry.heroku.com/hyperpic/web
	@docker push registry.heroku.com/hyperpic/web
	@heroku container:release web --app=hyperpic

upload-demo:
	@curl -F 'image=@_resources/demo/kayaks.jpg' -H "Authorization: Bearer $(HYPERPIC_AUTH_SECRET)" https://hyperpic.herokuapp.com/kayaks.jpg
	@curl -F 'image=@_resources/demo/smartcrop.jpg' -H "Authorization: Bearer $(HYPERPIC_AUTH_SECRET)" https://hyperpic.herokuapp.com/smartcrop.jpg


upload-demo-dev:
	@curl -F 'image=@_resources/demo/kayaks.jpg' -H "Authorization: Bearer foo" http://localhost:8574/kayaks.jpg
	@curl -F 'image=@_resources/demo/smartcrop.jpg' -H "Authorization: Bearer foo" http://localhost:8574/smartcrop.jpg

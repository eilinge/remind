
.PHONE: test, build, clean

app := remind

commit := $(shell git rev-parse HEAD)
commit_flag := -X main.Commit=$(commit)
build_time := $(shell date "+%Y-%m-%d.%H:%M:%S")
build_flag := -X main.Build=$(build_time)

docker_registry := registry.cn-shanghai.aliyuncs.com/eilingeloveduzi

# go mod vendor

run: export GOFLAGS=-mod=vendor
run:
	go run main.go

test: export GOFLAGS=-mod=vendor
test:
	go test -cover -race ./...

build: export CGO_ENABLED=0
build:
	go build -mod=vendor -v -ldflags '$(commit_flag) $(build_flag)' -o out/$(app)

deploy_test: export GOOS=linux
deploy_test: build
	ansible-playbook -i scripts/hosts scripts/deploy.yaml

clean:
	rm -rf out

.DEFAULT_GOAL=provision
.PHONY: build test get run tags

ensure_vendor:
	mkdir -pv vendor

clean:
	rm -rf ./vendor
	go clean .

vet:
	go vet .

generate:
	go generate -x

get:
	go get -u -v ./...

build: get vet generate
	go build .

test: build
	go test -v .

explore:
	go run main.go --level info explore

provision:
	go run main.go --level info provision --s3Bucket ${S3_BUCKET}

delete:
	go run main.go --level info delete

describe:
	go run main.go --level info describe --s3Bucket $(S3_BUCKET) --out ./graph.html

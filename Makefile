.DEFAULT_GOAL=provision
.PHONY: build test get run tags

clean:
	rm -rf ./vendor
	go clean .

format:
	go fmt .
	gofmt -s -w ./transforms/

vet:
	go vet .

generate:
	go generate -x

update:
	rm -rf $(GOPATH)/src/github.com/mweagle/Sparta
	mkdir -pv $(GOPATH)/src/github.com/mweagle/Sparta
	cp $(GOPATH)/src/Sparta/*.go $(GOPATH)/src/github.com/mweagle/Sparta
	cp -R $(GOPATH)/src/Sparta/aws $(GOPATH)/src/github.com/mweagle/Sparta
	cp -R $(GOPATH)/src/Sparta/resources $(GOPATH)/src/github.com/mweagle/Sparta
	rm -rf $(GOPATH)/src/github.com/mweagle/CloudFormationResources
	mkdir -pv $(GOPATH)/src/github.com/mweagle/CloudFormationResources
	cp $(GOPATH)/src/CloudFormationResources/*.go $(GOPATH)/src/github.com/mweagle/CloudFormationResources


get: clean
	rm -rf $(GOPATH)/src/github.com/aws/aws-sdk-go
	git clone --depth=1 https://github.com/aws/aws-sdk-go $(GOPATH)/src/github.com/aws/aws-sdk-go

	rm -rf $(GOPATH)/src/github.com/go-ini/ini
	git clone --depth=1 https://github.com/go-ini/ini $(GOPATH)/src/github.com/go-ini/ini

	rm -rf $(GOPATH)/src/github.com/jmespath/go-jmespath
	git clone --depth=1 https://github.com/jmespath/go-jmespath $(GOPATH)/src/github.com/jmespath/go-jmespath

	rm -rf $(GOPATH)/src/github.com/Sirupsen/logrus
	git clone --depth=1 https://github.com/Sirupsen/logrus $(GOPATH)/src/github.com/Sirupsen/logrus

	rm -rf $(GOPATH)/src/github.com/voxelbrain/goptions
	git clone --depth=1 https://github.com/voxelbrain/goptions $(GOPATH)/src/github.com/voxelbrain/goptions

	rm -rf $(GOPATH)/src/github.com/mjibson/esc
	git clone --depth=1 https://github.com/mjibson/esc $(GOPATH)/src/github.com/mjibson/esc

	rm -rf $(GOPATH)/src/github.com/crewjam/go-cloudformation
	git clone --depth=1 https://github.com/crewjam/go-cloudformation $(GOPATH)/src/github.com/crewjam/go-cloudformation

	rm -rf $(GOPATH)/src/github.com/mweagle/cloudformationresources
	git clone --depth=1 https://github.com/mweagle/cloudformationresources $(GOPATH)/src/github.com/mweagle/cloudformationresources

build: get format vet generate
	GO15VENDOREXPERIMENT=1 go build .

test: build
	GO15VENDOREXPERIMENT=1 go test ./test/...

tags:
	gotags -tag-relative=true -R=true -sort=true -f="tags" -fields=+l .

explore:
	go run main.go --level info explore

provision:
	go run main.go --level info provision --s3Bucket $(S3_BUCKET)

delete:
	go run main.go --level info delete

describe:
	go run main.go --level info describe --out ./graph.html

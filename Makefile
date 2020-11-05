TEST?=$$(go list ./... | grep -v 'vendor')
HOSTNAME=brinkmanlab.ca
NAMESPACE=brinkmanlab
NAME=galaxy
BINARY=terraform-provider-${NAME}
VERSION=0.2.3
OS_ARCH=linux_amd64

.ONESHELL:

.PHONY: default
default: install

.PHONY: build
build:
	go build -o ./bin/registry.terraform.io/${NAMESPACE}/${NAME}/${VERSION}/${OS_ARCH}/${BINARY}_${VERSION}

.PHONY: doc
doc:
	go run ./docgen

.PHONY: release
release:
	GOOS=darwin GOARCH=amd64 go build -o ./bin/${BINARY}_${VERSION}_darwin_amd64
	GOOS=freebsd GOARCH=386 go build -o ./bin/${BINARY}_${VERSION}_freebsd_386
	GOOS=freebsd GOARCH=amd64 go build -o ./bin/${BINARY}_${VERSION}_freebsd_amd64
	GOOS=freebsd GOARCH=arm go build -o ./bin/${BINARY}_${VERSION}_freebsd_arm
	GOOS=linux GOARCH=386 go build -o ./bin/${BINARY}_${VERSION}_linux_386
	GOOS=linux GOARCH=amd64 go build -o ./bin/${BINARY}_${VERSION}_linux_amd64
	GOOS=linux GOARCH=arm go build -o ./bin/${BINARY}_${VERSION}_linux_arm
	GOOS=openbsd GOARCH=386 go build -o ./bin/${BINARY}_${VERSION}_openbsd_386
	GOOS=openbsd GOARCH=amd64 go build -o ./bin/${BINARY}_${VERSION}_openbsd_amd64
	GOOS=solaris GOARCH=amd64 go build -o ./bin/${BINARY}_${VERSION}_solaris_amd64
	GOOS=windows GOARCH=386 go build -o ./bin/${BINARY}_${VERSION}_windows_386
	GOOS=windows GOARCH=amd64 go build -o ./bin/${BINARY}_${VERSION}_windows_amd64

.PHONY: install
install: build
	mkdir -p ~/.terraform.d/plugins/${HOSTNAME}/${NAMESPACE}/${NAME}/${VERSION}/${OS_ARCH}
	mv ${BINARY} ~/.terraform.d/plugins/${HOSTNAME}/${NAMESPACE}/${NAME}/${VERSION}/${OS_ARCH}

.PHONY: test
test: 
	go test -i $(TEST) || exit 1                                                   
	echo $(TEST) | xargs -t -n4 go test $(TESTARGS) -timeout=30s -parallel=4                    

.PHONY: testacc
testacc:
	function tearDown {
		docker kill $${TEST_BENCH}
	}
	trap tearDown EXIT
	TEST_BENCH=$$(docker run --rm -d -p 8080:80 -e GALAXY_CONFIG_OVERRIDE_ALLOW_USER_DELETION=true quay.io/bgruening/galaxy:19.09)
	until curl -sS --fail -o /dev/null "http://localhost:8080/api/version"; do sleep 1; done
	GALAXY_HOST=http://localhost:8080 GALAXY_API_KEY=admin TF_ACC=1 TF_LOG=DEBUG go test $(TEST) -v $(TESTARGS) -timeout 120m

TEST?=$$(go list ./... | grep -v 'vendor')
HOSTNAME=company.io
NAMESPACE=namespace
NAME=provider
BINARY=terraform-provider-${NAME}

OS_ARCH=darwin_amd64
VERSION=0.2
# provider source = company.io/namespace/provider
default: install

build:
	go build -o ${BINARY}
	chmod +x ${BINARY}

install:build
	mkdir -p ~/.terraform.d/plugins/${HOSTNAME}/${NAMESPACE}/${NAME}/${VERSION}/${OS_ARCH}
	mv ${BINARY} ~/.terraform.d/plugins/${HOSTNAME}/${NAMESPACE}/${NAME}/${VERSION}/${OS_ARCH}

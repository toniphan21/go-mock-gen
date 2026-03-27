GIT_COMMIT := $(shell git rev-parse --short HEAD)
RELEASE := $(shell sed -n 's/const Version = "\(.*\)"/\1/p' version.go)

SERVER := root@nhatp.com
WEB_ROOT := /var/www
BASE_PATH := /go/mock-gen
WEB_PATH := $(WEB_ROOT)$(BASE_PATH)

build-pkl:
	rm -rf ./.out
	pkl project package ./pkl --output-path=./.out --env-var=RELEASE=$(RELEASE)

upload-pkl:
	scp ./.out/pkl@$(RELEASE)* $(SERVER):$(WEB_PATH)

release-pkl: build-pkl upload-pkl

version:
	echo $(RELEASE)
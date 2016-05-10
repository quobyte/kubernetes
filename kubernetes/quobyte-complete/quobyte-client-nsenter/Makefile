.PHONY:	clean build push

IMAGE = quobyte-client-nsenter
TAG = $(shell head -n1 Dockerfile | cut -d: -f2)
ORG = quay.io/quobyte

all: push

clean:
	rm -rf nsenter

build: clean
	docker run --rm -v $$PWD:/target jpetazzo/nsenter
	docker build -t $(ORG)/$(IMAGE):$(TAG) .

push: build
	docker push $(ORG)/$(IMAGE):$(TAG)

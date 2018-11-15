
.PHONY: all source receiver

all: source receiver

SOURCE_IMAGE_NAME=$(if $(ENV_SOURCE_IMAGE_NAME),$(ENV_SOURCE_IMAGE_NAME),docker.io/rootfs/rgwpubsub-knative-source)
RECEIVER_IMAGE_NAME=$(if $(ENV_RECEIVER_IMAGE_NAME),$(ENV_RECEIVER_IMAGE_NAME),docker.io/rootfs/rgwpubsub-knative-receiver)

source:
	CGO_ENABLED=0 GOOS=linux go build -a -ldflags '-extldflags "-static"' -o  _output/rgwpubsub ./container-source

receiver:
	CGO_ENABLED=0 GOOS=linux go build -a -ldflags '-extldflags "-static"' -o  _output/receiver ./receiver

image: source receiver
	cp _output/rgwpubsub  container/source
	docker build -t $(SOURCE_IMAGE_NAME) container/source
	cp _output/receiver  container/receiver
	docker build -t $(RECEIVER_IMAGE_NAME) container/receiver

push-image: image
	docker push $(SOURCE_IMAGE_NAME)
	docker push $(RECEIVER_IMAGE_NAME)

clean:
	go clean -r -x
	rm -f _output
	rm -f container/source/rgwpubsub
	rm -f container/receiver/receiver


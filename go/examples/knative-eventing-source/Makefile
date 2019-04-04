TARGETS= source receiver vision resnet

.PHONY: all $(TARGETS)

REPO=rootfs

all: $(TARGETS)

source:
	CGO_ENABLED=0 GOOS=linux go build -ldflags '-extldflags "-static"' -o  _output/source ./container-source

receiver:
	CGO_ENABLED=0 GOOS=linux go build -ldflags '-extldflags "-static"' -o  _output/receiver ./receiver

vision:
	CGO_ENABLED=0 GOOS=linux go build -ldflags '-extldflags "-static"' -o  _output/vision ./vision

resnet:
	CGO_ENABLED=0 GOOS=linux go build -ldflags '-extldflags "-static"' -o  _output/resnet ./resnet-grpc


image: $(TARGETS)
	for t in $(TARGETS); do \
			cp _output/$$t  container/$$t; \
			docker build -t docker.io/$(REPO)/rgwpubsub-knative-$$t container/$$t; \
			rm container/$$t/$$t; \
	done

push-image: image
	for t in $(TARGETS); do \
		docker push docker.io/$(REPO)/rgwpubsub-knative-$$t;\
	done

clean:
	go clean -r -x
	rm -rf _output
	for t in $(TARGETS); do \
		rm -f container/$$t/$$t; \
	done



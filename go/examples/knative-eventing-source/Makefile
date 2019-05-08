TARGETS=source receiver vision resnet http-to-knative

.PHONY: all clean image push-image $(TARGETS)

REPO=rootfs

all: $(TARGETS)

source:
	gofmt -s -w ./container-source && \
	CGO_ENABLED=0 GOOS=linux go build -ldflags '-extldflags "-static"' -o  _output/$@ ./container-source

receiver:
	gofmt -s -w $@ && \
	CGO_ENABLED=0 GOOS=linux go build -ldflags '-extldflags "-static"' -o  _output/$@ ./$@

vision:
	gofmt -s -w $@ && \
	CGO_ENABLED=0 GOOS=linux go build -ldflags '-extldflags "-static"' -o  _output/$@ ./$@

resnet:
	gofmt -s -w ./resnet-grpc && \
	CGO_ENABLED=0 GOOS=linux go build -ldflags '-extldflags "-static"' -o  _output/resnet ./resnet-grpc

http-to-knative:
	gofmt -s -w $@ && \
	CGO_ENABLED=0 GOOS=linux go build -ldflags '-extldflags "-static"' -o  _output/$@ ./$@

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



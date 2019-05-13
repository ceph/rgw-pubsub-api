TARGETS=receiver http-to-aws test-receiver

.PHONY: create-test-receiver create-receiver all clean $(TARGETS)

all: $(TARGETS)

create-test-receiver:
	cd test-receiver && \
	gofmt -s -w main.go && \
	CGO_ENABLED=0 GOOS=linux go build -ldflags '-extldflags "-static"' -o  main main.go && \
	zip main.zip main && \
	aws lambda create-function \
  	--region ${AWS_REGION} \
  	--function-name test-receiver \
  	--memory 128 \
  	--role arn:aws:iam::${AWS_ACCOUNT}:role/service-role/${AWS_IAM_USER} \
  	--runtime go1.x \
  	--zip-file fileb://main.zip \
  	--handler main

create-receiver:
	cd receiver && \
	go generate && \
	gofmt -s -w *.go && \
	CGO_ENABLED=0 GOOS=linux go build -ldflags '-extldflags "-static"' -o  main *.go && \
	zip main.zip main && \
	aws lambda create-function \
  	--region ${AWS_REGION} \
  	--function-name receiver \
  	--memory 128 \
  	--role arn:aws:iam::${AWS_ACCOUNT}:role/service-role/${AWS_IAM_USER} \
  	--runtime go1.x \
  	--zip-file fileb://main.zip \
  	--handler main

test-receiver:
	cd $@ && \
	gofmt -s -w main.go && \
	CGO_ENABLED=0 GOOS=linux go build -ldflags '-extldflags "-static"' -o  main main.go && \
	zip main.zip main && \
	aws lambda update-function-code \
  	--region ${AWS_REGION} \
  	--function-name $@ \
  	--zip-file fileb://main.zip 

receiver:
	cd $@ && \
	go generate && \
	gofmt -s -w *.go && \
	CGO_ENABLED=0 GOOS=linux go build -ldflags '-extldflags "-static"' -o  main *.go && \
	zip main.zip main && \
	aws lambda update-function-code \
  	--region ${AWS_REGION} \
  	--function-name $@ \
  	--zip-file fileb://main.zip 

http-to-aws:
	gofmt -s -w $@ && \
	CGO_ENABLED=0 GOOS=linux go build -ldflags '-extldflags "-static"' -o  _output/$@ ./$@

clean:
	go clean -r -x
	rm -f receiver/main
	rm -f receiver/main.zip
	rm -f test-receiver/main
	rm -f test-receiver/main.zip
	rm -rf _output

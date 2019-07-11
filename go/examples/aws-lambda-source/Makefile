TARGETS=receiver http-to-aws

.PHONY: create-receiver all clean $(TARGETS)

all: $(TARGETS)

create-receiver:
	cd receiver && \
	gofmt -s -w *.go && \
	CGO_ENABLED=0 GOOS=linux go build -ldflags '-extldflags "-static"' -o  main *.go && \
	zip main.zip main && \
	aws lambda create-function \
  	--region ${AWS_DEFAULT_REGION} \
	--environment Variables="{S3_ACCESS_KEY_ID=${S3_ACCESS_KEY_ID},S3_SECRET_ACCESS_KEY=${S3_SECRET_ACCESS_KEY},S3_HOSTNAME=${S3_HOSTNAME},S3_REGION=${S3_REGION},GOOGLE_VISION_API_KEY=${GOOGLE_VISION_API_KEY}}" \
  	--function-name receiver \
  	--memory 128 \
  	--role arn:aws:iam::${AWS_ACCOUNT}:role/service-role/${AWS_IAM_USER} \
  	--runtime go1.x \
  	--zip-file fileb://main.zip \
  	--handler main

receiver:
	cd $@ && \
	gofmt -s -w *.go && \
	CGO_ENABLED=0 GOOS=linux go build -ldflags '-extldflags "-static"' -o  main *.go && \
	zip main.zip main && \
	aws lambda update-function-code \
  	--region ${AWS_DEFAULT_REGION} \
  	--function-name $@ \
  	--zip-file fileb://main.zip 

http-to-aws:
	gofmt -s -w $@ && \
	CGO_ENABLED=0 GOOS=linux go build -ldflags '-extldflags "-static"' -o  _output/$@ ./$@

clean:
	go clean -r -x
	rm -f receiver/main
	rm -f receiver/main.zip
	rm -rf _output

delete-receiver:
	aws lambda delete-function \
  		--region ${AWS_DEFAULT_REGION} \
		--function-name receiver


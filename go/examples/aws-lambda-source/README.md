# Build
Install ```aws``` command according to: https://docs.aws.amazon.com/cli/latest/userguide/install-linux.html
Set the following environment variables for the **build** process:
```AWS_REGION``` - the region where the lambda function should be created (e.g. "us-east-1")
```AWS_ACCOUNT``` - your AWS account ID
```AWS_IAM_USER``` - an AWS user allowed to create and update lambda functions
```AWS_ACCESS_KEY_ID``` - access ID for the AWS
```AWS_SECRET_ACCESS_KEY``` - secret key for the AWS
```AWS_DEFAULT_REGION``` - same as ```AWS_REGION```
```S3_ACCESS_KEY_ID``` - access ID for the radosgw
```S3_SECRET_ACCESS_KEY``` - secret key for the radosgw
```S3_HOSTNAME``` - URL for the main radosgw (e.g. "http://my-rgw:8000")
> Note that the hostname should be accessible from the lambda function. So, don't use something like "localhost" there
```GOOGLE_VISION_API_KEY``` - Google Vision API key

Run: ```make create-receiver``` to create the lambda function. This should be done once, and updates to the lambda function code would happen when running: ```make``` or ```make receiver```.
Run: ```make``` to build binaries and update lambda function (after it was created)

# Run
The following environment variable should be changed to the radosgw hostname of the pubsub zone: ```S3_HOSTNAME```
Execute the event converter: ```_output/http-to-aws -username tester -v 10 -zonegroup a -subscriptionname catsub -port 9001```
> Note that the port must match the endpoint definition configured on the subscription.

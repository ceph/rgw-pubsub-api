# Run RGW Pubsub go client

## Build

```bash
go build rgw_ps_cli.go
```

## Setup test environment

```console
export S3_ACCESS_KEY_ID=<s3 key id>
export S3_SECRET_ACCESS_KEY=<s3 secret>
export S3_HOSTNAME=<rgw pubsub endpoint URL>
```
## Run test

```bash
./rgw_ps_cli -topicname="foobar" -subscription="sub1" -zonegroup=a -username="rgwtest" -cleanup=true -v 10
```

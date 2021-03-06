# Build
First install missing packages:
```bash
go get google.golang.org/grpc
go get github.com/golang/protobuf/proto
```
Then build all binaries using: ```make```

# Setup

Follow [Knative installation instructions](https://www.knative.dev/docs/install/)

## Install Knative Serving
```bash
curl -L https://github.com/knative/serving/releases/download/v0.4.0/serving.yaml \
  | sed 's/LoadBalancer/NodePort/' \
  | kubectl apply --filen
```

## Install Knative Eventing
```bash
kubectl apply --filename https://github.com/knative/eventing/releases/download/v0.4.0/release.yaml
kubectl apply --filename https://github.com/knative/eventing-sources/releases/download/v0.4.0/release.yaml
```

## Install RGW Event source
Edit `sources_v1alpha1_containersources_rgwpubsub.yaml` and `service-entry.yaml` to reflect local settings.

```bash
cd deploy

kubectl apply -f rgwpubsub-ns.yaml
kubectl apply -f channel.yaml
kubectl apply -f subscription.yaml
kubectl apply -f service-entry.yaml
kubectl apply -f sources_v1alpha1_containersources_rgwpubsub.yaml
```

# Health checking

## Service is up and ready
```console
# kubectl get revision -n rgwpubsub
NAME                  SERVICE NAME                  READY   REASON
rgwpubsub-svc-00001   rgwpubsub-svc-00001-service   True
```
## RGW PubSub messages are received and posted
```console
# kubectl logs -lsource=containersource-rgwpubsub -c source -n rgwpubsub
2018/11/19 15:25:11 Target is: "http://rgw-ps-channel-channel.default.svc.cluster.local/"
2018/11/19 15:25:16 Posting to "http://rgw-ps-channel-channel.default.svc.cluster.local/"
2018/11/19 15:25:17 response Status: 202 Accepted
```
## Events are received on Knative serving function
```console
# kubectl logs -lserving.knative.dev/service=rgwpubsub-svc -c user-container | more
2018/11/19 15:25:08 Ready and listening on port 8080
2018/11/19 15:25:17 [2018-11-19T15:25:16Z] application/json rgwpubsub. Object: "test3"  Bucket: "buck"
```

# Advanced Applications

More advanced applications can go into the service functions. As illustrated in the [Google vision app](vision) and [ResNet app](resnet-grpc), the serving
functions can retrieve the RGW objects, send images to inference services, and get their annotation/classes.

## Google Vision Serving Function

This serving function uses Google Vision service to annotate an image.

First, edit [service-entry.yaml](deploy/google-vision-svc/service-entry.yaml) and [subscription.yaml](deploy/google-vision-svc/subscription.yaml)
to reflect local RGW settings and your Google Vision API Key.

Then run the following:

```bash
kubectl apply -f deploy/google-vision-svc/service-entry.yaml
kubectl apply -f deploy/google-vision-svc/subscription.yaml
```

Then upload an cat image into RGW:

```console
# wget https://r.hswstatic.com/w_907/gif/tesla-cat.jpg
# ./s3 put buck/cat1.jpg --in-file=./tesla-cat.jpg
# ./s3 put buck/cat2.jpg --in-file=./tesla-cat.jpg
```

Checking the serving container's log:
```console
# kubectl logs -lserving.knative.dev/service=rgwpubsub-svc -c user-container 
2018/11/29 16:22:49 Ready and listening on port 8080
2018/11/29 16:23:42 [2018-11-29T16:23:41Z] application/json rgwpubsub. Object: "cat1.jpg"  Bucket: "buck"
2018/11/29 16:23:43 label: cat, Score: 0.993347
2018/11/29 16:25:01 [2018-11-29T16:25:01Z] application/json rgwpubsub. Object: "cat2.jpg"  Bucket: "buck"
2018/11/29 16:25:02 label: cat, Score: 0.993347
```

The cat is identified!

## Tensorflow ResNet Serving Function

This serving function uses ResNet with a  pre-trained ImageNet model to classify an image.

First, edit [service-entry.yaml](deploy/resnet-grpc/service-entry.yaml) and [subscription-resnet.yaml](deploy/resnet-grpc/subscription-resnet.yaml)
to reflect local RGW settings and your Tensorflow Serving endpoint.

Then run the following:

```bash
kubectl apply -f deploy/resnet-grpc/service-entry.yaml
kubectl apply -f deploy/resnet-grpc/subscription-grpc.yaml
```

Then upload cat and dog images into RGW:

```console
# wget https://r.hswstatic.com/w_907/gif/tesla-cat.jpg
# wget https://upload.wikimedia.org/wikipedia/commons/d/d9/Collage_of_Nine_Dogs.jpg
#./s3 put buck/dogs.jpg --in-file=./Collage_of_Nine_Dogs.jpg
#./s3 put buck/telsa-cat.jpg --in-file=./tesla-cat.jpg
```

Checking the serving container's log:
```console
# kubectl logs -lserving.knative.dev/service=rgwpubsub-svc -c user-container
2018/11/29 18:52:23 Ready and listening on port 8080
2018/11/29 18:54:31 [2018-11-29T18:54:31Z] application/json rgwpubsub. Object: "dogs.jpg"  Bucket: "buck"
2018/11/29 18:54:32 classes: [162]
2018/11/29 18:57:54 [2018-11-29T18:57:51Z] application/json rgwpubsub. Object: "telsa-cat.jpg"  Bucket: "buck"
2018/11/29 18:57:54 classes: [286]
```

Note, refer to ImageNet classes, classes 162 is 'beagle', class 286 is 'cougar, puma, catamount, mountain lion, painter, panther, Felis concolor'.
The classifier is close enough!

# Splitting Traffic between Google Vision and ResNet functions

This is a more advanced Knative serving configuration: splitting traffic between two functions.

## Why to split traffic?

There are a plenty of use cases for this, e.g. you start with your own inference service but as business grows, a Cloud burst is
the natural choice as your in-house inference service fail to scale economically.

## How to split traffic in Knative
Splitting traffic between multiple services requires a more advanced Knative configuration. We need to use `Route` rather than `Service`
in our `Subscription`, as illustrated below:
```yaml
apiVersion: eventing.knative.dev/v1alpha1
kind: Subscription
metadata:
  name: rgw-ps-subscription
spec:
  channel:
    apiVersion: eventing.knative.dev/v1alpha1
    kind: Channel
    name: rgw-ps-channel
  subscriber:
    ref:
      apiVersion: serving.knative.dev/v1alpha1
      kind: Route
      name: rgwpubsub-route
```

In Knative Serving, a route is an aggregation of configurations and ratios of how much traffic is dispatched into each configration.
Our traffic splitting example evenly divides requests into Google Vision and ResNet:

```yaml
apiVersion: serving.knative.dev/v1alpha1
kind: Route
metadata:
  name: rgwpubsub-route
spec:
 traffic:
 - configurationName: google-vision-configuration
   percent: 50
 - configurationName: resnet-configuration
   percent: 50
```

A Knative Serving `Configuration` is very similar to `Service`, as illustrated in [route.yaml](deploy/split-traffic/route.yaml)

## Test traffic split

First, edit [route.yaml](deploy/split-traffic/route.yaml) to reflect your RGW credentials and endpoints and Google Vision API key.

Then run the following:

```console
kubectl apply -f deploy/split-traffic/route.yaml
```

Once the serving Pods are up and running, upload some images to RGW:

```console
# wget https://r.hswstatic.com/w_907/gif/tesla-cat.jpg
# for i in $(seq 1 10); do ./s3 put buck/cat-${i}.jpg --in-file=./tesla-cat.jpg; done
```

From the logs of the ResNet Serving function:

```console
# kubectl logs -lserving.knative.dev/configuration=resnet-configuration -c user-container
2018/12/03 15:20:49 Ready and listening on port 8080
2018/12/03 15:24:10 [2018-12-03T15:24:10Z] application/json rgwpubsub. Object: "cat-7.jpg"  Bucket: "buck"
2018/12/03 15:24:11 classes: [286]
2018/12/03 15:24:12 [2018-12-03T15:24:12Z] application/json rgwpubsub. Object: "cat-2.jpg"  Bucket: "buck"
2018/12/03 15:24:12 classes: [286]
2018/12/03 15:24:12 [2018-12-03T15:24:12Z] application/json rgwpubsub. Object: "cat-3.jpg"  Bucket: "buck"
2018/12/03 15:24:12 classes: [286]
2018/12/03 15:24:13 [2018-12-03T15:24:13Z] application/json rgwpubsub. Object: "cat-4.jpg"  Bucket: "buck"
2018/12/03 15:24:13 classes: [286]
2018/12/03 15:24:16 [2018-12-03T15:24:16Z] application/json rgwpubsub. Object: "cat-9.jpg"  Bucket: "buck"
2018/12/03 15:24:16 classes: [286]
2018/12/03 15:24:16 [2018-12-03T15:24:16Z] application/json rgwpubsub. Object: "cat-10.jpg"  Bucket: "buck"
2018/12/03 15:24:16 classes: [286]
```

And from the logs of the Google Vision Serving fucntion:
```console
# kubectl logs -lserving.knative.dev/configuration=google-vision-configuration -c user-container
2018/12/03 15:20:48 Ready and listening on port 8080
2018/12/03 15:24:11 [2018-12-03T15:24:11Z] application/json rgwpubsub. Object: "cat-1.jpg"  Bucket: "buck"
2018/12/03 15:24:11 label: cat, Score: 0.993347
2018/12/03 15:24:11 [2018-12-03T15:24:11Z] application/json rgwpubsub. Object: "cat-5.jpg"  Bucket: "buck"
2018/12/03 15:24:12 label: cat, Score: 0.993347
2018/12/03 15:24:12 [2018-12-03T15:24:12Z] application/json rgwpubsub. Object: "cat-6.jpg"  Bucket: "buck"
2018/12/03 15:24:13 label: cat, Score: 0.993347
2018/12/03 15:24:15 [2018-12-03T15:24:15Z] application/json rgwpubsub. Object: "cat-8.jpg"  Bucket: "buck"
2018/12/03 15:24:16 label: cat, Score: 0.993347
```

The traffic indeed ran into the two serving functions.

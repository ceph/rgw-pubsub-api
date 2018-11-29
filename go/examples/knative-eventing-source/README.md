# Setup

Follow [Knative installation instructions](https://github.com/knative/docs)

## Install Knative serving
## Install Knative eventing
## Install Knative event source CRD and Controller


# Install RGW Event source

Edit `sources_v1alpha1_containersources_rgwpubsub.yaml` and `service-entry.yaml` to reflect local settings.

```bash
cd deploy
kubectl apply -f channel.yaml
kubectl apply -f subscription.yaml
kubectl apply -f service-entry.yaml
kubectl apply -f sources_v1alpha1_containersources_rgwpubsub.yaml
```

# Health checking

## Service is up and ready
```console
# kubectl get revision
NAME                  SERVICE NAME                  READY   REASON
rgwpubsub-svc-00001   rgwpubsub-svc-00001-service   True
```
## RGW PubSub messages are received and posted
```console
#kubectl logs -lsource=containersource-rgwpubsub -c source
2018/11/19 15:25:11 Target is: "http://rgw-ps-channel-channel.default.svc.cluster.local/"
2018/11/19 15:25:16 Posting to "http://rgw-ps-channel-channel.default.svc.cluster.local/"
2018/11/19 15:25:17 response Status: 202 Accepted
```
## Events are received on Knative serving function
```console
# kubectl logs -lserving.knative.dev/service=rgwpubsub-svc -c user-container |more
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

This serving function uses ResNet witha  pre-trained ImageNet model to classify an image.

First, edit [service-entry.yaml](deploy/resnet-grpc/service-entry.yaml) and [subscription-resnet.yaml](deploy/resnet-grpc/subscription-grpc.yaml)
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

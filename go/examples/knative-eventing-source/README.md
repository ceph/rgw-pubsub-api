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


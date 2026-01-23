root@node1 ~ # helm install redis-cluster oci://registry-1.docker.io/bitnamicharts/redis-cluster -n prod
Pulled: registry-1.docker.io/bitnamicharts/redis-cluster:11.2.1
Digest: sha256:f761a6f26a5a54649e40dae85322664262dca565963e7c763740e5cf9818a0ce
NAME: redis-cluster
LAST DEPLOYED: Sat Dec 28 04:23:14 2024
NAMESPACE: prod
STATUS: deployed
REVISION: 1
TEST SUITE: None
NOTES:
CHART NAME: redis-cluster
CHART VERSION: 11.2.1
APP VERSION: 7.4.1

Did you know there are enterprise versions of the Bitnami catalog? For enhanced secure software supply chain features, unlimited pulls from Docker, LTS support, or application customization, see Bitnami Premium or Tanzu Application Catalog. See https://www.arrow.com/globalecs/na/vendors/bitnami for more information.** Please be patient while the chart is being deployed **


To get your password run:
export REDIS_PASSWORD=$(kubectl get secret --namespace "prod" redis-cluster -o jsonpath="{.data.redis-password}" | base64 -d)

You have deployed a Redis&reg; Cluster accessible only from within you Kubernetes Cluster.INFO: The Job to create the cluster will be created.To connect to your Redis&reg; cluster:

1. Run a Redis&reg; pod that you can use as a client:
   kubectl run --namespace prod redis-cluster-client --rm --tty -i --restart='Never' \
   --env REDIS_PASSWORD=$REDIS_PASSWORD \
   --image docker.io/bitnami/redis-cluster:7.4.1-debian-12-r3 -- bash

2. Connect using the Redis&reg; CLI:

redis-cli -c -h redis-cluster -a $REDIS_PASSWORD



WARNING: There are "resources" sections in the chart not set. Using "resourcesPreset" is not recommended for production. For production installations, please set the following values according to your workload needs:
- redis.resources
- updateJob.resources
  +info https://kubernetes.io/docs/concepts/configuration/manage-resources-containers/
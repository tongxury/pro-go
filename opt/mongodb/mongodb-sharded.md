
helm uninstall mongodb-sharded   -n prod
helm install mongodb-sharded  oci://registry-1.docker.io/bitnamicharts/mongodb-sharded -n prod


root@node1 ~ # helm install mongodb-sharded  oci://registry-1.docker.io/bitnamicharts/mongodb-sharded -n prod
Pulled: registry-1.docker.io/bitnamicharts/mongodb-sharded:9.1.0
Digest: sha256:6e16b0a9f40da3a5ffeecb587992e5afb7806df52da22ff8ab7ff8644f23d72d
NAME: mongodb-sharded
LAST DEPLOYED: Fri Dec 27 09:02:20 2024
NAMESPACE: prod
STATUS: deployed
REVISION: 1
TEST SUITE: None
NOTES:
CHART NAME: mongodb-sharded
CHART VERSION: 9.1.0
APP VERSION: 8.0.4

Did you know there are enterprise versions of the Bitnami catalog? For enhanced secure software supply chain features, unlimited pulls from Docker, LTS support, or application customization, see Bitnami Premium or Tanzu Application Catalog. See https://www.arrow.com/globalecs/na/vendors/bitnami for more information.

** Please be patient while the chart is being deployed **

The MongoDB&reg; Sharded cluster can be accessed via the Mongos instances in port 27017 on the following DNS name from within your cluster:

    mongodb-sharded.prod.svc.cluster.local

To get the root password run:

    export MONGODB_ROOT_PASSWORD=$(kubectl get secret --namespace prod mongodb-sharded -o jsonpath="{.data.mongodb-root-password}" | base64 -d)

To connect to your database run the following command:

    kubectl run --namespace prod mongodb-sharded-client --rm --tty -i --restart='Never' --image docker.io/bitnami/mongodb-sharded:8.0.4-debian-12-r0 --command -- mongosh admin --host mongodb-sharded --authenticationDatabase admin -u root -p $MONGODB_ROOT_PASSWORD

To connect to your database from outside the cluster execute the following commands:

    kubectl port-forward --namespace prod svc/mongodb-sharded 27017:27017 &
    mongosh --host 127.0.0.1 --authenticationDatabase admin -p $MONGODB_ROOT_PASSWORD

WARNING: There are "resources" sections in the chart not set. Using "resourcesPreset" is not recommended for production. For production installations, please set the following values according to your workload needs:
- configsvr.resources
- mongos.resources
- shardsvr.dataNode.resources
  +info https://kubernetes.io/docs/concepts/configuration/manage-resources-containers/
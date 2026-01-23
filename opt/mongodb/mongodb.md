helm install mongodb oci://registry-1.docker.io/bitnamicharts/mongodb -n prod
helm uninstall mongodb -n prod

Pulled: registry-1.docker.io/bitnamicharts/mongodb:16.3.1
Digest: sha256:325ffcd6acbefa8181eeccac8689dcb3d294b4d17af5edd3f4d217924141ed79
NAME: mongodb
LAST DEPLOYED: Thu Nov 28 10:17:29 2024
NAMESPACE: prod
STATUS: deployed
REVISION: 1
TEST SUITE: None
NOTES:
CHART NAME: mongodb
CHART VERSION: 16.3.1
APP VERSION: 8.0.3

** Please be patient while the chart is being deployed **

MongoDB&reg; can be accessed on the following DNS name(s) and ports from within your cluster:

    mongodb.prod.svc.cluster.local

To get the root password run:

    export MONGODB_ROOT_PASSWORD=$(kubectl get secret --namespace prod mongodb -o jsonpath="{.data.mongodb-root-password}" | base64 -d)

To connect to your database, create a MongoDB&reg; client container:

    kubectl run --namespace prod mongodb-client --rm --tty -i --restart='Never' --env="MONGODB_ROOT_PASSWORD=$MONGODB_ROOT_PASSWORD" --image docker.io/bitnami/mongodb:8.0.3-debian-12-r0 --command -- bash

Then, run the following command:
mongosh admin --host "mongodb" --authenticationDatabase admin -u $MONGODB_ROOT_USER -p $MONGODB_ROOT_PASSWORD

To connect to your database from outside the cluster execute the following commands:

    kubectl port-forward --namespace prod svc/mongodb 27017:27017 &
    mongosh --host 127.0.0.1 --authenticationDatabase admin -p $MONGODB_ROOT_PASSWORD

WARNING: There are "resources" sections in the chart not set. Using "resourcesPreset" is not recommended for production. For production installations, please set the following values according to your workload needs:
- arbiter.resources
- resources
  +info https://kubernetes.io/docs/concepts/configuration/manage-resources-containers/
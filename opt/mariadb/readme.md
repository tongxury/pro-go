helm install mariadb  oci://registry-1.docker.io/bitnamicharts/mariadb

Pulled: registry-1.docker.io/bitnamicharts/mariadb:20.1.0
Digest: sha256:88c1c428670a3fd1f4eec54ecd83dbdb1af1acd35ef79fd387a73eeee86981c2
NAME: mariadb
LAST DEPLOYED: Thu Nov 28 11:16:09 2024
NAMESPACE: default
STATUS: deployed
REVISION: 1
TEST SUITE: None
NOTES:
CHART NAME: mariadb
CHART VERSION: 20.1.0
APP VERSION: 11.4.4

** Please be patient while the chart is being deployed **

Tip:

Watch the deployment status using the command: kubectl get pods -w --namespace default -l app.kubernetes.io/instance=mariadb

Services:

echo Primary: mariadb.default.svc.cluster.local:3306

Administrator credentials:

Username: root
Password : $(kubectl get secret --namespace default mariadb -o jsonpath="{.data.mariadb-root-password}" | base64 -d)

To connect to your database:

1. Run a pod that you can use as a client:

   kubectl run mariadb-client --rm --tty -i --restart='Never' --image  docker.io/bitnami/mariadb:11.4.4-debian-12-r0 --namespace default --command -- bash

2. To connect to primary service (read/write):

   mysql -h mariadb.default.svc.cluster.local -uroot -p my_database

To upgrade this helm chart:

1. Obtain the password as described on the 'Administrator credentials' section and set the 'auth.rootPassword' parameter as shown below:

   ROOT_PASSWORD=$(kubectl get secret --namespace default mariadb -o jsonpath="{.data.mariadb-root-password}" | base64 -d)
   helm upgrade --namespace default mariadb oci://registry-1.docker.io/bitnamicharts/mariadb --set auth.rootPassword=$ROOT_PASSWORD

WARNING: There are "resources" sections in the chart not set. Using "resourcesPreset" is not recommended for production. For production installations, please set the following values according to your workload needs:
- primary.resources
- secondary.resources
  +info https://kubernetes.io/docs/concepts/configuration/manage-resources-containers/
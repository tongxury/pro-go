helm install kafka oci://registry-1.docker.io/bitnamicharts/kafka --set image.tag=3.6.1-debian-11-r24 -n prod
helm install kafka -f values.yaml oci://registry-1.docker.io/bitnamicharts/kafka  -n prod
helm install kafka oci://registry-1.docker.io/bitnamicharts/kafka -n prod

helm uninstall kafka -n prod 

kubectl run kafka-client --restart='Never' --image docker.io/bitnami/kafka:3.6.1-debian-11-r24 --namespace prod --command -- sleep infinity
kubectl cp --namespace prod client.properties kafka-client:/tmp/client.properties
kubectl exec --tty -i kafka-client --namespace prod -- bash

kafka-console-producer.sh --producer.config /tmp/client.properties --broker-list kafka5-controller-0.kafka5-controller-headless.prod.svc.cluster.local:9092,kafka5-controller-1.kafka5-controller-headless.prod.svc.cluster.local:9092,kafka5-controller-2.kafka5-controller-headless.prod.svc.cluster.local:9092 --topic sing
kafka-console-consumer.sh --consumer.config /tmp/client.properties --bootstrap-server kafka5.prod.svc.cluster.local:9092 --topic log_received --from-beginning

kafka.prod.svc.cluster.local

kubectl -n prod port-forward --address 0.0.0.0 service/kafka 9092:9092
kubectl -n prod port-forward --address 0.0.0.0 service/kafka-headless 9092:9092

[//]: # (kafka-headless.prod:9092)
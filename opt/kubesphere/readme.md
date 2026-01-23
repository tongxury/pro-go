

rm -rf /etc/etcd.env kubekey /etc/ssl/etcd 
./kk create cluster -f config-sample.yaml  --with-local-storage
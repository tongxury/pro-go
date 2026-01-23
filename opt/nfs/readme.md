apt install nfs-kernel-server
!!!! 客户端机器也要下载 nft相关包

mkdir -p /data/nfs-storage
chmod 777 /data/nfs-storage

vim /etc/exports 加入下面这句话
/data/nfs-storage *(rw,insecure,sync,no_subtree_check,no_root_squash)

systemctl restart nfs-kernel-server
systemctl enable nfs-kernel-server

mkdir /tmp/testnfs6 \
&& mount -t nfs 13.212.187.169:/srv/nfs4/homes /tmp/testnfs6 \
&& echo "hello nfs" >> /tmp/testnfs6/test1.txt \
&& cat /tmp/testnfs6/test1.txt

helm repo add nfs-subdir-external-provisioner https://kubernetes-sigs.github.io/nfs-subdir-external-provisioner/

helm install nfs-subdir-external-provisioner nfs-subdir-external-provisioner/nfs-subdir-external-provisioner \
--set nfs.server=18.140.5.226 \
--set nfs.path=/data/nfs-storage

helm install redis-cluster oci://registry-1.docker.io/bitnamicharts/redis-cluster -n prod
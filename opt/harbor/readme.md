

helm install harbor harbor/harbor \
--set expose.type=domain \
--set expose.host=harbor.example.com \
--set persistence.storageClass=your-storage-class \
--set persistence.enabled=true \
--set harborAdminPassword=your-admin-password 
[//]: # ()
[//]: # (helm repo add ingress-nginx https://kubernetes.github.io/ingress-nginx)

[//]: # ()
[//]: # (helm install nginx ingress-nginx/ingress-nginx -n prod)

[//]: # ()
[//]: # (helm uninstall nginx -n prod)



好的！使用私有镜像仓库是最可靠的方式。我看到您的 Makefile 中使用的是 yoozy-cn-shanghai.cr.volces.com。以下是完整的步骤：

步骤 1: 在本地（能访问外网的机器）拉取镜像
bash
# 拉取 ingress-nginx controller 镜像
docker pull registry.k8s.io/ingress-nginx/controller:v1.9.4

# 如果上面拉不到，用镜像加速
docker pull docker.m.daocloud.io/ingress-nginx/controller:v1.9.4
步骤 2: 重新打标签并推送到您的私有仓库
bash
# 假设您的私有仓库是 yoozy-cn-shanghai.cr.volces.com
PRIVATE_REGISTRY="yoozy-cn-shanghai.cr.volces.com/k8s"

# 重新打标签
docker tag registry.k8s.io/ingress-nginx/controller:v1.9.4 ${PRIVATE_REGISTRY}/ingress-nginx-controller:v1.9.4

# 登录私有仓库（如果需要）
docker login yoozy-cn-shanghai.cr.volces.com

# 推送到私有仓库
docker push ${PRIVATE_REGISTRY}/ingress-nginx-controller:v1.9.4
步骤 3: 在 K8s 集群中使用私有仓库镜像安装
bash
helm install ingress-nginx ingress-nginx/ingress-nginx \
-n prod \
--set controller.image.registry=yoozy-cn-shanghai.cr.volces.com \
--set controller.image.image=k8s/ingress-nginx-controller \
--set controller.image.tag=v1.9.4 \
--set controller.admissionWebhooks.enabled=false \
--set defaultBackend.enabled=false
完整脚本（一键完成）
创建一个脚本 push-ingress-images.sh：

bash
#!/bin/bash

# 配置您的私有仓库地址
PRIVATE_REGISTRY="yoozy-cn-shanghai.cr.volces.com/k8s"

# 镜像列表
IMAGES=(
"registry.k8s.io/ingress-nginx/controller:v1.9.4"
"registry.k8s.io/ingress-nginx/kube-webhook-certgen:v20231011-8b53cabe0"
)

# 如果官方镜像拉不到，使用镜像加速
MIRROR_IMAGES=(
"docker.m.daocloud.io/ingress-nginx/controller:v1.9.4"
"docker.m.daocloud.io/ingress-nginx/kube-webhook-certgen:v20231011-8b53cabe0"
)

echo "开始拉取镜像..."

for i in "${!IMAGES[@]}"; do
OFFICIAL_IMAGE="${IMAGES[$i]}"
MIRROR_IMAGE="${MIRROR_IMAGES[$i]}"

# 提取镜像名称和标签
IMAGE_NAME=$(echo $OFFICIAL_IMAGE | awk -F'/' '{print $NF}' | awk -F':' '{print $1}')
IMAGE_TAG=$(echo $OFFICIAL_IMAGE | awk -F':' '{print $NF}')

echo "处理镜像: $IMAGE_NAME:$IMAGE_TAG"

# 尝试拉取官方镜像，失败则使用镜像加速
if ! docker pull $OFFICIAL_IMAGE 2>/dev/null; then
echo "官方镜像拉取失败，使用镜像加速..."
docker pull $MIRROR_IMAGE
docker tag $MIRROR_IMAGE $OFFICIAL_IMAGE
fi

# 重新打标签
PRIVATE_IMAGE="${PRIVATE_REGISTRY}/${IMAGE_NAME}:${IMAGE_TAG}"
docker tag $OFFICIAL_IMAGE $PRIVATE_IMAGE

# 推送到私有仓库
echo "推送到私有仓库: $PRIVATE_IMAGE"
docker push $PRIVATE_IMAGE

echo "完成: $PRIVATE_IMAGE"
echo "---"
done

echo "所有镜像已推送到私有仓库！"
步骤 4: 执行脚本并安装
bash
# 在本地执行
chmod +x push-ingress-images.sh
./push-ingress-images.sh

# 在 K8s 集群上安装
helm install ingress-nginx ingress-nginx/ingress-nginx \
-n prod \
--set controller.image.registry=yoozy-cn-shanghai.cr.volces.com \
--set controller.image.image=k8s/ingress-nginx-controller \
--set controller.image.tag=v1.9.4 \
--set controller.admissionWebhooks.enabled=false \
--set defaultBackend.enabled=false
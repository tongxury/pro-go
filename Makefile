VERSION=$(shell git describe --tags --always)
COMMIT_MESSAGE = $(shell git log --pretty=format:"%s" $(VERSION) -1)

INTERNAL_PROTO_FILES=$(shell find internal -name *.proto)
API_PROTO_FILES=$(shell find api -name *.proto)

#REGISTRY=registry.eu-central-1.aliyuncs.com
#REGISTRY_REPO=registry.eu-central-1.aliyuncs.com/proreg
REGISTRY_REPO=usernx
VERSION := $(shell date +%Y%m%d%H%M%S)

.PHONY: init
# init env
init:
	# brew install protobuf
	go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
	go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
	go install github.com/go-kratos/kratos/cmd/kratos/v2@latest
	go install github.com/go-kratos/kratos/cmd/protoc-gen-go-http/v2@latest
	go install github.com/google/gnostic/cmd/protoc-gen-openapi@latest
	go install github.com/google/wire/cmd/wire@latest
	go install github.com/envoyproxy/protoc-gen-validate@latest
	go install github.com/favadi/protoc-go-inject-tag@latest

.PHONY: api
# generate api proto
api:
	protoc --proto_path=./api \
           --proto_path=./pkg/publicprotos \
 	       --go_out=paths=source_relative:./api \
 	       --validate_out=paths=source_relative,lang=go:./api \
 	       --go-http_out=paths=source_relative:./api \
 	       --go-grpc_out=paths=source_relative:./api \
	       --openapi_out=fq_schema_naming=true,default_response=false:. \
	       $(API_PROTO_FILES)
	protoc-go-inject-tag -input=./api/**/**.pb.go

.PHONY: generate
# generate
generate: api
	#go mod tidy
	#go get github.com/google/wire/cmd/wire
	go generate ./...

.PHONY: build
# build
build:
	mkdir -p bin && CGO_ENABLED=0 GOOS=linux go build -ldflags "-X main.Version=$(VERSION)" -o ./bin ./app/${SRV}/...

build-deploy:
	@docker buildx build -t ${REGISTRY_REPO}/${SUB_ID}:${VER} --build-arg SRV=${SRV} .
	@docker push ${REGISTRY_REPO}/${SUB_ID}:${VER}

#    @aws ecr get-login-password --region cn-north-1 | docker login --username AWS --password-stdin 501487856280.dkr.ecr.cn-north-1.amazonaws.com.cn
	@docker tag ${REGISTRY_REPO}/${SUB_ID}:${VER} 501487856280.dkr.ecr.cn-north-1.amazonaws.com.cn/veogo/${SUB_ID}:${VER}
	@docker push 501487856280.dkr.ecr.cn-north-1.amazonaws.com.cn/veogo/${SUB_ID}:${VER}

##    docker login crpi-yju5wi02oeo074jw.cn-hangzhou.personal.cr.aliyuncs.com -u nick1153926706 -p veogo.ai
#	@docker build -t crpi-yju5wi02oeo074jw.cn-hangzhou.personal.cr.aliyuncs.com/${REGISTRY_REPO}/${SUB_ID}:${VER} --build-arg SRV=${SRV} .
#	@docker push crpi-yju5wi02oeo074jw.cn-hangzhou.personal.cr.aliyuncs.com/${REGISTRY_REPO}/${SUB_ID}:${VER}

release_dex-tgbot-server:
	@make build-deploy SRV=dex-tgbot SUB_ID=p3 VER=1.0.0112
	@curl -X PUT \
         -H "Content-Type: application/yaml" \
         -H "Cookie: KuboardUsername=admin; KuboardAccessKey=8sy2yedcjpm6.sthcttwjn7dnp72ty2m6xbrr8z7873m6" \
         -d '{"kind":"deployments","namespace":"prod","name":"dex-tgbot-server"}' \
         "http://3.0.101.23:18060/kuboard-api/cluster/trade-master/kind/CICDApi/admin/resource/restartWorkload"

release_dex-admin:
	@make build-deploy SRV=dex-admin SUB_ID=dm VER=latest
	@curl -X PUT \
         -H "Content-Type: application/yaml" \
         -H "Cookie: KuboardUsername=admin; KuboardAccessKey=t8xmfrrantn6.w3bwm73wn4f5rmk8s7ya2zjfbsax8j6c" \
         -d '{"kind":"deployments","namespace":"prod","name":"dex-admin"}' \
         "http://13.212.187.169:18060/kuboard-api/cluster/pro/kind/CICDApi/admin/resource/restartWorkload"


release_tracker:
	@make build-deploy SRV=tracker SUB_ID=p11 VER=latest


release_market:
	@make build-deploy SRV=dex-market SUB_ID=mk VER=latest
	@curl -X PUT \
         -H "Content-Type: application/yaml" \
         -H "Cookie: KuboardUsername=admin; KuboardAccessKey=t8xmfrrantn6.w3bwm73wn4f5rmk8s7ya2zjfbsax8j6c" \
         -d '{"kind":"deployments","namespace":"prod","name":"dex-market"}' \
         "http://13.212.187.169:18060/kuboard-api/cluster/trade-master/kind/CICDApi/admin/resource/restartWorkload"





release_dex-gateway:
	@make build-deploy SRV=dex-gateway SUB_ID=dgw VER=latest
	@curl -X PUT \
         -H "Content-Type: application/yaml" \
         -H "Cookie: KuboardUsername=admin; KuboardAccessKey=t8xmfrrantn6.w3bwm73wn4f5rmk8s7ya2zjfbsax8j6c" \
         -d '{"kind":"deployments","namespace":"prod","name":"dex-gateway"}' \
         "http://13.212.187.169:18060/kuboard-api/cluster/trade-master/kind/CICDApi/admin/resource/restartWorkload"

release_user:
	@docker buildx build -t yoozy-cn-shanghai.cr.volces.com/ve/user:${VERSION} --build-arg SRV=user -f Dockerfile-arm64 . --push
	@curl -X PUT \
         -H "content-type: application/json" \
         -H "Cookie: KuboardUsername=admin; KuboardAccessKey=4342pex2bf55.kkx8hf6z665b85fsf38ichrnfrc238kz" \
         -d '{"kind":"deployments","namespace":"prod","name":"user","images":{"yoozy-cn-shanghai.cr.volces.com/ve/user":"yoozy-cn-shanghai.cr.volces.com/ve/user:${VERSION}"}}' \
         "http://118.196.63.209:18060/kuboard-api/cluster/qiniu/kind/CICDApi/admin/resource/updateImageTag"
#	@docker push yoozy-cn-shanghai.cr.volces.com/ve/user:${VERSION}
#	@make build-deploy SRV=user SUB_ID=user VER=$(VERSION)
#	@curl -X PUT \
#         -H "content-type: application/json" \
#         -H "Cookie: KuboardUsername=admin; KuboardAccessKey=eb3ie8csecpa.3dzkhe8x7jhpddz2fcwkx6xycbwfppk4" \
#         -d '{"kind":"deployments","namespace":"prod","name":"user","images":{"501487856280.dkr.ecr.cn-north-1.amazonaws.com.cn/veogo/user":"501487856280.dkr.ecr.cn-north-1.amazonaws.com.cn/veogo/user:$(VERSION)"}}' \
#         "http://118.196.63.209:18060/kuboard-api/cluster/veogo-cn/kind/CICDApi/admin/resource/updateImageTag"
#	@curl -X PUT \
#		-H "content-type: application/json" \
#		-H "Cookie: KuboardUsername=admin; KuboardAccessKey=s76zhyaezhcp.wszzs7wnsp4kmxkh2zzwzi6y2mjh34jr" \
#		-d '{"kind":"deployments","namespace":"prod","name":"user","images":{"usernx/user":"usernx/user:$(VERSION)"}}' \
#		"http://13.250.123.59:18060/kuboard-api/cluster/eks-ap/kind/CICDApi/admin/resource/updateImageTag"
#


release_payment:
	@docker buildx build --platform linux/arm64 -t yoozy-cn-shanghai.cr.volces.com/ve/payment:${VERSION} --build-arg SRV=payment -f Dockerfile-arm64 . --push

#	@make build-deploy SRV=payment SUB_ID=payment VER=$(VERSION)
#	@curl -X PUT \
#         -H "content-type: application/json" \
#         -H "Cookie: KuboardUsername=admin; KuboardAccessKey=eb3ie8csecpa.3dzkhe8x7jhpddz2fcwkx6xycbwfppk4" \
#         -d '{"kind":"deployments","namespace":"prod","name":"payment","images":{"501487856280.dkr.ecr.cn-north-1.amazonaws.com.cn/veogo/payment":"501487856280.dkr.ecr.cn-north-1.amazonaws.com.cn/veogo/payment:$(VERSION)"}}' \
#         "http://118.196.63.209:18060/kuboard-api/cluster/veogo-cn/kind/CICDApi/admin/resource/updateImageTag"
#	@curl -X PUT \
#         -H "content-type: application/json" \
#         -H "Cookie: KuboardUsername=admin; KuboardAccessKey=5t8cfdibws6j.wi436t7c33jyszrj5ahtzctsk4zfpjtf" \
#         -d '{"kind":"deployments","namespace":"prod","name":"payment","images":{"usernx/payment":"usernx/payment:$(VERSION)"}}' \
#         "http://13.250.123.59:18060/kuboard-api/cluster/eks-ap/kind/CICDApi/admin/resource/updateImageTag"


release_bi:
	@make build-deploy SRV=bi SUB_ID=bi VER=$(VERSION)
	@curl -X PUT \
         -H "content-type: application/json" \
         -H "Cookie: KuboardUsername=admin; KuboardAccessKey=eb3ie8csecpa.3dzkhe8x7jhpddz2fcwkx6xycbwfppk4" \
         -d '{"kind":"deployments","namespace":"prod","name":"bi","images":{"501487856280.dkr.ecr.cn-north-1.amazonaws.com.cn/veogo/bi":"501487856280.dkr.ecr.cn-north-1.amazonaws.com.cn/veogo/bi:$(VERSION)"}}' \
         "http://118.196.63.209:18060/kuboard-api/cluster/veogo-cn/kind/CICDApi/admin/resource/updateImageTag"
	@curl -X PUT \
		 -H "content-type: application/json" \
		 -H "Cookie: KuboardUsername=admin; KuboardAccessKey=s76zhyaezhcp.wszzs7wnsp4kmxkh2zzwzi6y2mjh34jr" \
		 -d '{"kind":"deployments","namespace":"prod","name":"bi","images":{"usernx/bi":"usernx/bi:$(VERSION)"}}' \
		 "http://13.250.123.59:18060/kuboard-api/cluster/eks-ap/kind/CICDApi/admin/resource/updateImageTag"

release_aiagent:
	@docker buildx build --platform linux/arm64 -t yoozy-cn-shanghai.cr.volces.com/ve/aiagent:${VERSION} --build-arg SRV=aiagent -f Dockerfile-arm64 . --push
	@curl -X PUT \
         -H "content-type: application/json" \
         -H "Cookie: KuboardUsername=admin; KuboardAccessKey=4342pex2bf55.kkx8hf6z665b85fsf38ichrnfrc238kz" \
         -d '{"kind":"deployments","namespace":"prod","name":"aiagent","images":{"yoozy-cn-shanghai.cr.volces.com/ve/aiagent":"yoozy-cn-shanghai.cr.volces.com/ve/aiagent:${VERSION}"}}' \
         "http://118.196.63.209:18060/kuboard-api/cluster/qiniu/kind/CICDApi/admin/resource/updateImageTag"


usercenter_my_release:
	@docker buildx build  -t usernx/usercenter:${VERSION} --build-arg SRV=usercenter -f Dockerfile-arm64 . --push
	@curl -X PUT \
    -H "content-type: application/json" \
    -H "Cookie: KuboardUsername=admin; KuboardAccessKey=4342pex2bf55.kkx8hf6z665b85fsf38ichrnfrc238kz" \
    -d '{"kind":"deployments","namespace":"prod","name":"usercenter","images":{"usernx/usercenter":"usernx/usercenter:${VERSION}"}}' \
    "http://118.196.63.209:18060/kuboard-api/cluster/my-server/kind/CICDApi/admin/resource/updateImageTag"

usercenter_qiniu_release:
	@docker buildx build --platform=linux/arm64  -t yoozy-cn-shanghai.cr.volces.com/ve/usercenter:${VERSION} --build-arg SRV=usercenter -f Dockerfile-arm64 . --push
	@curl -X PUT \
         -H "content-type: application/json" \
         -H "Cookie: KuboardUsername=admin; KuboardAccessKey=4342pex2bf55.kkx8hf6z665b85fsf38ichrnfrc238kz" \
         -d '{"kind":"deployments","namespace":"prod","name":"usercenter","images":{"yoozy-cn-shanghai.cr.volces.com/ve/usercenter":"yoozy-cn-shanghai.cr.volces.com/ve/usercenter:${VERSION}"}}' \
         "http://118.196.63.209:18060/kuboard-api/cluster/qiniu/kind/CICDApi/admin/resource/updateImageTag"

voiceagent_qiniu_release:
	@docker buildx build --platform linux/arm64 -t yoozy-cn-shanghai.cr.volces.com/ve/voiceagent:${VERSION} --build-arg SRV=voiceagent -f Dockerfile-arm64 . --push
	@curl -X PUT \
         -H "content-type: application/json" \
         -H "Cookie: KuboardUsername=admin; KuboardAccessKey=4342pex2bf55.kkx8hf6z665b85fsf38ichrnfrc238kz" \
         -d '{"kind":"deployments","namespace":"prod","name":"voiceagent","images":{"yoozy-cn-shanghai.cr.volces.com/ve/voiceagent":"yoozy-cn-shanghai.cr.volces.com/ve/voiceagent:${VERSION}"}}' \
         "http://118.196.63.209:18060/kuboard-api/cluster/qiniu/kind/CICDApi/admin/resource/updateImageTag"

#	@docker push yoozy-cn-shanghai.cr.volces.com/ve/aiagent:${VERSION}
#	@make build-deploy SRV=aiagent SUB_ID=p12 VER=$(VERSION)
#	@curl -X PUT \
#         -H "content-type: application/json" \
#         -H "Cookie: KuboardUsername=admin; KuboardAccessKey=eb3ie8csecpa.3dzkhe8x7jhpddz2fcwkx6xycbwfppk4" \
#         -d '{"kind":"deployments","namespace":"prod","name":"aiagent","images":{"501487856280.dkr.ecr.cn-north-1.amazonaws.com.cn/veogo/p12":"501487856280.dkr.ecr.cn-north-1.amazonaws.com.cn/veogo/p12:$(VERSION)"}}' \
#         "http://118.196.63.209:18060/kuboard-api/cluster/veogo-cn/kind/CICDApi/admin/resource/updateImageTag"
#
#	@curl -X PUT \
#         -H "content-type: application/json" \
#         -H "Cookie: KuboardUsername=admin; KuboardAccessKey=s76zhyaezhcp.wszzs7wnsp4kmxkh2zzwzi6y2mjh34jr" \
#         -d '{"kind":"deployments","namespace":"prod","name":"aiagent","images":{"usernx/p12":"usernx/p12:$(VERSION)"}}' \
#         "http://13.250.123.59:18060/kuboard-api/cluster/eks-ap/kind/CICDApi/admin/resource/updateImageTag"
#
#	@curl -X PUT \
#         -H "content-type: application/json" \
#         -H "Cookie: KuboardUsername=admin; KuboardAccessKey=s76zhyaezhcp.wszzs7wnsp4kmxkh2zzwzi6y2mjh34jr" \
#         -d '{"kind":"deployments","namespace":"prod","name":"aiagent-cronjob","images":{"usernx/p12":"usernx/p12:$(VERSION)"}}' \
#         "http://13.250.123.59:18060/kuboard-api/cluster/eks-ap/kind/CICDApi/admin/resource/updateImageTag"

release_monday:
	@docker buildx build -t yoozy-cn-shanghai.cr.volces.com/yoozy/monday:${VERSION} --build-arg SRV=monday -f Dockerfile-chrome .
	@docker push yoozy-cn-shanghai.cr.volces.com/yoozy/monday:${VERSION}
	@curl -X PUT \
         -H "content-type: application/json" \
         -H "Cookie: KuboardUsername=admin; KuboardAccessKey=eb3ie8csecpa.3dzkhe8x7jhpddz2fcwkx6xycbwfppk4" \
         -d '{"kind":"deployments","namespace":"prod","name":"monday","images":{"yoozy-cn-shanghai.cr.volces.com/yoozy/monday":"yoozy-cn-shanghai.cr.volces.com/yoozy/monday:${VERSION}"}}' \
         "http://118.196.63.209:18060/kuboard-api/cluster/yoozy-cn/kind/CICDApi/admin/resource/updateImageTag"
	@docker rmi -f yoozy-cn-shanghai.cr.volces.com/yoozy/monday:${VERSION}


release_monday_admin:
	@docker buildx build -t yoozy-cn-shanghai.cr.volces.com/yoozy/monday-admin:${VERSION} --build-arg SRV=monday-admin -f Dockerfile-chrome .
	@docker push yoozy-cn-shanghai.cr.volces.com/yoozy/monday-admin:${VERSION}
	@curl -X PUT \
         -H "content-type: application/json" \
         -H "Cookie: KuboardUsername=admin; KuboardAccessKey=eb3ie8csecpa.3dzkhe8x7jhpddz2fcwkx6xycbwfppk4" \
         -d '{"kind":"deployments","namespace":"prod","name":"monday-admin","images":{"yoozy-cn-shanghai.cr.volces.com/yoozy/monday-admin":"yoozy-cn-shanghai.cr.volces.com/yoozy/monday-admin:${VERSION}"}}' \
         "http://118.196.63.209:18060/kuboard-api/cluster/yoozy-cn/kind/CICDApi/admin/resource/updateImageTag"
	@docker rmi -f yoozy-cn-shanghai.cr.volces.com/yoozy/monday-admin:${VERSION}


env:
	@docker buildx build --platform linux/amd64,linux/arm64,linux/arm/v7 -t ${REGISTRY_REPO}/goenv:v1.0.0 -f Dockerfile-env .
	@docker push ${REGISTRY_REPO}/goenv:v1.0.0

goenv:
	@docker buildx build --platform=linux/amd64 -t yoozy-cn-shanghai.cr.volces.com/yoozy/goenv:v3.0.5 -f Dockerfile-goenv . --push



release_proj_user:
	@docker buildx build -t yoozy-cn-shanghai.cr.volces.com/yoozy/user:${VERSION} --build-arg SRV=user -f Dockerfile-yoozy2 . --push
	@curl -X PUT \
         -H "content-type: application/json" \
         -H "Cookie: KuboardUsername=admin; KuboardAccessKey=eb3ie8csecpa.3dzkhe8x7jhpddz2fcwkx6xycbwfppk4" \
         -d '{"kind":"deployments","namespace":"prod","name":"user","images":{"yoozy-cn-shanghai.cr.volces.com/yoozy/user":"yoozy-cn-shanghai.cr.volces.com/yoozy/user:${VERSION}"}}' \
         "http://118.196.63.209:18060/kuboard-api/cluster/yoozy-cn/kind/CICDApi/admin/resource/updateImageTag"
	@docker rmi -f yoozy-cn-shanghai.cr.volces.com/yoozy/user:${VERSION}


release_proj_admin:
	@docker buildx build --platform linux/amd64 -t yoozy-cn-shanghai.cr.volces.com/yoozy/proj-admin:${VERSION} --build-arg SRV=proj-admin -f Dockerfile-yoozy2 .  --push
	@curl -X PUT \
         -H "content-type: application/json" \
         -H "Cookie: KuboardUsername=admin; KuboardAccessKey=4342pex2bf55.kkx8hf6z665b85fsf38ichrnfrc238kz" \
         -d '{"kind":"deployments","namespace":"prod","name":"admin","images":{"yoozy-cn-shanghai.cr.volces.com/yoozy/proj-admin":"yoozy-cn-shanghai.cr.volces.com/yoozy/proj-admin:${VERSION}"}}' \
         "http://118.196.63.209:18060/kuboard-api/cluster/vke-shanghai/kind/CICDApi/admin/resource/updateImageTag"
	@docker rmi -f yoozy-cn-shanghai.cr.volces.com/yoozy/proj-admin:${VERSION}


test_proj_admin:
	@docker buildx build --platform=linux/amd64  -t yoozy-cn-shanghai.cr.volces.com/yoozy/proj-admin:test --build-arg SRV=proj-admin -f Dockerfile-yoozy2 . --push
	@curl -X PUT \
         -H "Content-Type: application/yaml" \
         -H "Cookie: KuboardUsername=admin; KuboardAccessKey=4342pex2bf55.kkx8hf6z665b85fsf38ichrnfrc238kz" \
         -d '{"kind":"deployments","namespace":"beta","name":"admin"}' \
         "http://118.196.63.209:18060/kuboard-api/cluster/vke-shanghai/kind/CICDApi/admin/resource/restartWorkload"

release_proj:
	@docker buildx build --platform=linux/amd64  -t yoozy-cn-shanghai.cr.volces.com/yoozy/proj:${VERSION} --build-arg SRV=proj -f Dockerfile-chrome . --push
	@curl -X PUT \
         -H "content-type: application/json" \
         -H "Cookie: KuboardUsername=admin; KuboardAccessKey=4342pex2bf55.kkx8hf6z665b85fsf38ichrnfrc238kz" \
         -d '{"kind":"deployments","namespace":"prod","name":"proj","images":{"yoozy-cn-shanghai.cr.volces.com/yoozy/proj":"yoozy-cn-shanghai.cr.volces.com/yoozy/proj:${VERSION}"}}' \
         "http://118.196.63.209:18060/kuboard-api/cluster/vke-shanghai/kind/CICDApi/admin/resource/updateImageTag"
	@docker rmi -f yoozy-cn-shanghai.cr.volces.com/yoozy/proj:${VERSION}


release_proj_pro:
	@docker buildx build --platform=linux/amd64 -t yoozy-cn-shanghai.cr.volces.com/yoozy/proj-pro:${VERSION} --build-arg SRV=proj-pro -f Dockerfile-yoozy2 . --push
	@curl -X PUT \
         -H "content-type: application/json" \
         -H "Cookie: KuboardUsername=admin; KuboardAccessKey=4342pex2bf55.kkx8hf6z665b85fsf38ichrnfrc238kz" \
         -d '{"kind":"deployments","namespace":"prod","name":"proj-pro","images":{"yoozy-cn-shanghai.cr.volces.com/yoozy/proj-pro":"yoozy-cn-shanghai.cr.volces.com/yoozy/proj-pro:${VERSION}"}}' \
         "http://118.196.63.209:18060/kuboard-api/cluster/vke-shanghai/kind/CICDApi/admin/resource/updateImageTag"
	@docker rmi -f yoozy-cn-shanghai.cr.volces.com/yoozy/proj-pro:${VERSION}


release_proj_pro_chromedp:
	@docker buildx build --platform=linux/amd64 -t yoozy-cn-shanghai.cr.volces.com/yoozy/proj-pro-chromedp:${VERSION} --build-arg SRV=proj-pro -f Dockerfile-chromedp . --push
#	@curl -X PUT \
#         -H "content-type: application/json" \
#         -H "Cookie: KuboardUsername=admin; KuboardAccessKey=4342pex2bf55.kkx8hf6z665b85fsf38ichrnfrc238kz" \
#         -d '{"kind":"deployments","namespace":"prod","name":"proj-pro","images":{"yoozy-cn-shanghai.cr.volces.com/yoozy/proj-pro":"yoozy-cn-shanghai.cr.volces.com/yoozy/proj-pro:${VERSION}"}}' \
#         "http://118.196.63.209:18060/kuboard-api/cluster/vke-shanghai/kind/CICDApi/admin/resource/updateImageTag"
#	@docker rmi -f yoozy-cn-shanghai.cr.volces.com/yoozy/proj-pro:${VERSION}


test_proj_pro:
	@docker buildx build --platform=linux/amd64  -t yoozy-cn-shanghai.cr.volces.com/yoozy/proj-pro:test --build-arg SRV=proj-pro -f Dockerfile-yoozy2 . --push
	@curl -X PUT \
         -H "Content-Type: application/yaml" \
         -H "Cookie: KuboardUsername=admin; KuboardAccessKey=4342pex2bf55.kkx8hf6z665b85fsf38ichrnfrc238kz" \
         -d '{"kind":"deployments","namespace":"beta","name":"proj-pro"}' \
         "http://118.196.63.209:18060/kuboard-api/cluster/vke-shanghai/kind/CICDApi/admin/resource/restartWorkload"

release_usercenter:
	@docker buildx build --platform=linux/amd64  -t yoozy-cn-shanghai.cr.volces.com/yoozy/usercenter:${VERSION} --build-arg SRV=usercenter -f Dockerfile-yoozy2 . --push
	@curl -X PUT \
         -H "content-type: application/json" \
         -H "Cookie: KuboardUsername=admin; KuboardAccessKey=4342pex2bf55.kkx8hf6z665b85fsf38ichrnfrc238kz" \
         -d '{"kind":"deployments","namespace":"prod","name":"usercenter","images":{"yoozy-cn-shanghai.cr.volces.com/yoozy/usercenter":"yoozy-cn-shanghai.cr.volces.com/yoozy/usercenter:${VERSION}"}}' \
         "http://118.196.63.209:18060/kuboard-api/cluster/vke-shanghai/kind/CICDApi/admin/resource/updateImageTag"




test_usercenter:
	@docker buildx build --platform=linux/amd64  -t yoozy-cn-shanghai.cr.volces.com/yoozy/usercenter:test --build-arg SRV=usercenter -f Dockerfile-yoozy2 . --push
	@curl -X PUT \
         -H "Content-Type: application/yaml" \
         -H "Cookie: KuboardUsername=admin; KuboardAccessKey=4342pex2bf55.kkx8hf6z665b85fsf38ichrnfrc238kz" \
         -d '{"kind":"deployments","namespace":"beta","name":"usercenter"}' \
         "http://118.196.63.209:18060/kuboard-api/cluster/vke-shanghai/kind/CICDApi/admin/resource/restartWorkload"

.DEFAULT_GOAL := build

#helm install ingress-nginx ingress-nginx/ingress-nginx \
#--namespace prod \
#--set controller.image.repository=registry.aliyuncs.com/google_containers/nginx-ingress-controller \
#--set controller.image.tag=v1.9.4 \
#--set controller.admissionWebhooks.patch.image.repository=registry.aliyuncs.com/google_containers/kube-webhook-certgen \
#--set controller.admissionWebhooks.patch.image.tag=v1.3.0
#
#
#chart=oci://hub.kubesphere.com.cn/kse/ks-core
#version=1.2.3
#helm upgrade --install -n kubesphere-system --create-namespace ks-core $chart --debug --wait --version $version --reset-values --set extension.imageRegistry=swr.cn-north-9.myhuaweicloud.com/ks

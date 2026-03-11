VERSION=latest

SERVER_NAME=game
SERVER_TYPE=rpc

## 测试环境配置

# 镜像发布地址
DOCKER_REPO_TEST=crpi-52ct675t9ypuyp0h.cn-hangzhou.personal.cr.aliyuncs.com/dicetales/${SERVER_NAME}-${SERVER_TYPE}
# 测试版本
VERSION_TEST=$(VERSION)
# 镜像名称
APP_NAME_TEST=dicetales-${SERVER_NAME}-${SERVER_TYPE}-test
# Dockerfile编译文件
DOCKER_FILE_TEST=./deploy/docker/Dockerfile_${SERVER_NAME}_${SERVER_TYPE}

# 测试环境的编译发布
build-test:
	GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o bin/${SERVER_NAME}-${SERVER_TYPE} ./apps/${SERVER_NAME}/${SERVER_TYPE}/${SERVER_NAME}.go
	docker build . -f ${DOCKER_FILE_TEST} --no-cache -t ${APP_NAME_TEST}

# 镜像的测试标签
tag-test:
	@echo 'create tag ${VERSION_TEST}'
	docker tag ${APP_NAME_TEST} ${DOCKER_REPO_TEST}:${VERSION_TEST}

publish-test:
	@echo 'publish ${VERSION_TEST} to ${DOCKER_REPO_TEST}'
	docker push $(DOCKER_REPO_TEST):${VERSION_TEST}

release-test: build-test tag-test publish-test

build-local-test: 
	GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o bin/${SERVER_NAME}-${SERVER_TYPE} ./apps/${SERVER_NAME}/${SERVER_TYPE}/${SERVER_NAME}.go
	docker cp bin/${SERVER_NAME}-${SERVER_TYPE} ${SERVER_NAME}-${SERVER_TYPE}-test:/${SERVER_NAME}/bin/
	docker cp ./apps/${SERVER_NAME}/${SERVER_TYPE}/etc/${SERVER_NAME}.yaml ${SERVER_NAME}-${SERVER_TYPE}-test:/${SERVER_NAME}/conf/
	docker restart ${SERVER_NAME}-${SERVER_TYPE}-test

fresh-local-test: build-test
	docker stop ${SERVER_NAME}-${SERVER_TYPE}-test 2>/dev/null || true
	docker rm ${SERVER_NAME}-${SERVER_TYPE}-test 2>/dev/null || true
	docker run -d --name ${SERVER_NAME}-${SERVER_TYPE}-test --network dicetale -p 10002:10002 ${APP_NAME_TEST}

local-test:
	@if docker ps -a --format '{{.Names}}' | grep -q "^${SERVER_NAME}-${SERVER_TYPE}-test$$"; then \
		echo "Container exists. Updating binary and restarting..."; \
		$(MAKE) -f deploy/make/game-rpc.mk build-local-test; \
	else \
		echo "Container not found. Building and starting fresh container..."; \
		$(MAKE) -f deploy/make/game-rpc.mk fresh-local-test; \
	fi

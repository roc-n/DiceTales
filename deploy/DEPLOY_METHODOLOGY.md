# go-zero 微服务部署与测试工具链沉淀

本文档记录了在本项目中，如何为一个新创建的 go-zero rpc/api 模块（例如 `game`）快速构建 Docker 和 Make 工具链，以便支持一键本地编译、容器构建与快速重启。

## 1. 核心架构逻辑

对于每一个微服务模块（以 `SERVER_NAME=game`, `SERVER_TYPE=rpc` 为例），我们需要三部分支持文件：
1. **Dockerfile**: 存放在 `deploy/docker/Dockerfile_game_rpc`
2. **Make 模板**: 存放在 `deploy/make/game-rpc.mk`
3. **主 Makefile 统一入口**: 位于项目根目录 `Makefile`

### Docker 容器的运行约定
- **编译产物路径**: 挂载到容器内的 `/game/bin/game-rpc`
- **配置文件路径**: 挂载到容器内的 `/game/conf/game.yaml`
- **基础镜像**: 使用 `alpine` 加载东八区时区。

---

## 2. 标准化操作步骤

假设我们要为全新的模块 `[模块名]`（如 `post`）搭建这套系统，请按照以下步骤执行：

### Step 1: 确定端口规划
检查现有的 yaml 配置文件并进行端口分配，防止端口冲突。
- `user-rpc`: 10000
- `social-rpc`: 10001
- `game-rpc`: 10002
- `[新模块]`: 10003...

### Step 2: 编写 Dockerfile
在 `deploy/docker/` 下创建 `Dockerfile_[模块名]_[类型]`（例如 `Dockerfile_post_rpc`）。
文件内容基本固定，只需要将 `ARG SERVER_NAME` 的值改为对应的模块名：
```dockerfile
FROM alpine:3.22.2

# 注入时区
RUN echo -e "https://mirrors.aliyun.com/alpine/v3.22/main\nhttps://mirrors.aliyun.com/alpine/v3.15/community" > /etc/apk/repositories && \
    apk update &&\
    apk --no-cache add tzdata && \
    cp /usr/share/zoneinfo/Asia/Shanghai /etc/localtime && \
    echo "Asia/Shanghai" > /etc/timezone
ENV TZ=Asia/Shanghai

ARG SERVER_NAME=[模块名]
ARG SERVER_TYPE=rpc

ENV RUN_BIN=bin/${SERVER_NAME}-${SERVER_TYPE}
ENV RUN_CONF=/${SERVER_NAME}/conf/${SERVER_NAME}.yaml

RUN mkdir /$SERVER_NAME && mkdir /$SERVER_NAME/bin && mkdir /$SERVER_NAME/conf
COPY ./bin/$SERVER_NAME-$SERVER_TYPE /$SERVER_NAME/bin/
COPY ./apps/$SERVER_NAME/$SERVER_TYPE/etc/$SERVER_NAME.yaml /$SERVER_NAME/conf/
RUN chmod +x /$SERVER_NAME/bin/$SERVER_NAME-$SERVER_TYPE

WORKDIR /$SERVER_NAME
ENTRYPOINT ["/[模块名]/bin/[模块名]-rpc"]
CMD ["-f", "/[模块名]/conf/[模块名].yaml"]
```

### Step 3: 编写专属 Make 模板
在 `deploy/make/` 目录下创建 `[模块名]-[类型].mk`。
复制已有的 `.mk` 模板，**重点修改以下两处**：
1. 头部的 `SERVER_NAME=[模块名]`
2. `fresh-local-test` 指令中的宿主机端口映射，如 `-p 10003:10003`

### Step 4: 更新主 Makefile 聚合指令
在根目录 `Makefile` 中追加该模块的专属快捷命令：
```makefile
[模块名]-rpc-gen:
	@make -f apps/[模块名]/rpc/gen.mk [模块名]-rpc-gen

[模块名]-rpc-local-test:
	@make -f deploy/make/[模块名]-rpc.mk local-test
```
注意这里必须使用 **Tab** 作为缩进。

---

## 3. 工具链的优势与运行机制
当我们在根目录敲下 `make game-rpc-local-test` 时的执行流程：
1. `make` 读取部署脚本，判断对应容器 (如 `game-rpc-test`) 是否存在。
2. **首次运行 (fresh-local-test)**: 用 `go build` 静态编译到 `bin/` 目录 -> 用 Dockerfile 将 `bin` 与 `yaml` 打包成镜像 -> 运行守护进程容器并挂载入自定义网络 (`dicetale`)。
3. **后续调测 (build-local-test)**: 不再重新打包整个镜像。仅在外部执行 `go build`，然后通过 `docker cp` 热替换容器内的二进制文件和配置文件，最后 `docker restart`。这极大提升了本地微服务 Vibe Coding 时改代码看效果的效率。

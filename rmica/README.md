# Simple Container Runtime Drop-in Replacement

这是一个简单的 runc drop-in replacement，实现了基本的容器生命周期管理功能。它可以作为 Docker 的替代运行时使用。

## 功能

- create: 创建容器
- start: 启动容器
- run: 创建并启动容器
- kill: 终止容器
- delete: 删除容器
- ps: 查看容器进程
- exec: 在容器中执行命令
- list: 列出容器
- state: 查看容器状态

## 构建

```bash
go build -o rmica
```

## 使用

### 作为独立运行时

```bash
# 创建容器
./rmica create <container-id>

# 启动容器
./rmica start <container-id>

# 创建并启动容器
./rmica run <container-id>

# 终止容器
./rmica kill <container-id> [signal]

# 删除容器
./rmica delete <container-id>

# 查看容器进程
./rmica ps <container-id>

# 在容器中执行命令
./rmica exec <container-id> <command>

# 列出容器
./rmica list

# 查看容器状态
./rmica state <container-id>
```

### 作为 Docker 运行时

1. 将编译好的 rmica 二进制文件复制到系统路径：
```bash
sudo cp rmica /usr/local/bin/
```

2. 配置 Docker 使用 rmica 作为运行时：
```bash
sudo mkdir -p /etc/docker
sudo tee /etc/docker/daemon.json <<-'EOF'
{
  "runtimes": {
    "rmica": {
      "path": "/usr/local/bin/rmica",
      "runtimeArgs": []
    }
  }
}
EOF
```

3. 重启 Docker 服务：
```bash
sudo systemctl restart docker
```

4. 使用 rmica 运行时运行容器：
```bash
docker run --runtime=rmica <image>
```

## 注意

这是一个简单的实现，目前只实现了基本的容器状态管理。要完全支持作为 Docker 运行时，还需要实现：

1. 完整的 OCI Runtime Specification 支持
2. 容器进程管理
3. 文件系统隔离
4. 网络隔离
5. 资源限制（cgroups）
6. 安全特性（capabilities, seccomp, etc.） 

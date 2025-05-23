# MicaRuntime

MicaRuntime is a simple container runtime implementation for containerd that demonstrates the shim v2 API integration. It currently only prints "hello" messages for each function call without implementing actual container functionality.

## Building

To build the MicaRuntime shim:

```bash
cd core/runtime/v2/micaruntime
go build -o containerd-shim-mica-v1 ./cmd
```

## Installation

1. Copy the built binary to a location in your PATH:
```bash
sudo cp containerd-shim-mica-v1 /usr/local/bin/
```

2. Configure containerd to use MicaRuntime by adding the following to your containerd config.toml:

```toml
[plugins."io.containerd.grpc.v1.cri".containerd.runtimes]
  [plugins."io.containerd.grpc.v1.cri".containerd.runtimes.mica]
    runtime_type = "io.containerd.mica.v1"
```

## Usage

### With containerd

You can use MicaRuntime with containerd by specifying the runtime when creating a container:

```bash
ctr run --runtime io.containerd.mica.v1 docker.io/library/alpine:latest test
```

### With Docker

To use MicaRuntime with Docker:

1. Configure Docker to use containerd as the runtime:
```bash
sudo mkdir -p /etc/docker
sudo tee /etc/docker/daemon.json <<-'EOF'
{
  "runtimes": {
    "micarun": {
      "path": "/usr/local/bin/containerd-shim-mica-v1",
      "runtimeArgs": []
    }
  }
}
EOF
```

2. Restart Docker:
```bash
sudo systemctl restart docker
```

3. Run a container using MicaRuntime:
```bash
docker run --runtime=micarun alpine echo "Hello from MicaRuntime"
```

## Notes

- This is a minimal implementation that only prints "hello" messages for each function call
- No actual container functionality is implemented
- The runtime is meant to demonstrate the shim v2 API integration with containerd 
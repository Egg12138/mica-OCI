# <a name="linuxRuntime" />Linux 运行时

## <a name="runtimeLinuxFileDescriptors" />文件描述符

默认情况下，运行时只为应用程序保持 `stdin`、`stdout` 和 `stderr` 文件描述符打开。
运行时可以向应用程序传递额外的文件描述符以支持[套接字激活][socket-activated-containers]等功能。
即使某些文件描述符是打开的，它们也可能被重定向到 `/dev/null`。

## <a name="runtimeLinuxDevSymbolicLinks" />Dev 符号链接

在创建容器时（[生命周期](runtime.md#lifecycle)中的第2步），如果在处理[`mounts`](config.md#mounts)后源文件存在，运行时必须创建以下符号链接：

|    源文件        | 目标文件    |
| --------------- | ----------- |
| /proc/self/fd   | /dev/fd     |
| /proc/self/fd/0 | /dev/stdin  |
| /proc/self/fd/1 | /dev/stdout |
| /proc/self/fd/2 | /dev/stderr |

[socket-activated-containers]: https://0pointer.de/blog/projects/socket-activated-containers.html 
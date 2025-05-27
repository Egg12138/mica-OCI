# <a name="ZOSContainerConfiguration" />z/OS 容器配置 (z/OS Container Configuration)

本文档描述了[容器配置](config.md)中[平台特定配置](config.md#platform-specific-configuration)的 z/OS 特定部分的模式。
z/OS 容器规范使用 z/OS UNIX 内核功能（如命名空间和文件系统监狱）来实现规范。

期望 z/OS 环境的应用程序很可能期望这些文件路径被正确设置。

以下文件系统应该在每个容器的文件系统中可用：

| 路径     | 类型   |
| -------- | ------ |
| /proc    | [proc][] |

## <a name="configZOSNamespaces" />命名空间 (Namespaces)

命名空间将全局系统资源包装在一个抽象中，使命名空间内的进程看起来拥有自己的全局资源隔离实例。
对全局资源的更改对作为命名空间成员的其他进程可见，但对其他进程不可见。
有关更多信息，请参见 https://www.ibm.com/docs/zos/latest?topic=planning-namespaces-zos-unix。

命名空间在 `namespaces` 根字段内指定为条目数组。
可以指定以下参数来设置命名空间：

* **`type`** *(string, 必需)* - 命名空间类型。应该支持以下命名空间类型：
    * **`pid`** 容器内的进程只能看到同一容器内或同一 pid 命名空间内的其他进程。
    * **`mount`** 容器将有一个隔离的挂载表。
    * **`ipc`** 容器内的进程只能通过系统级 IPC 与同一容器内的其他进程通信。
    * **`uts`** 容器将能够拥有自己的主机名和域名。
* **`path`** *(string, 可选)* - 命名空间文件。
    此值必须是[运行时挂载命名空间](glossary.md#runtime-namespace)中的绝对路径。
    运行时必须将容器进程放在与该 `path` 关联的命名空间中。
    如果 `path` 与类型为 `type` 的命名空间不关联，运行时必须[生成错误](runtime.md#errors)。

    如果未指定 `path`，运行时必须创建类型为 `type` 的新[容器命名空间](glossary.md#container-namespace)。

如果 `namespaces` 数组中未指定命名空间类型，容器必须继承该类型的[运行时命名空间](glossary.md#runtime-namespace)。
如果 `namespaces` 字段包含具有相同 `type` 的重复命名空间，运行时必须[生成错误](runtime.md#errors)。

### 示例

```json
"namespaces": [
    {
        "type": "pid",
        "path": "/proc/1234/ns/pid"
    },
    {
        "type": "mount"
    },
    {
        "type": "ipc"
    },
    {
        "type": "uts"
    }
]
``` 
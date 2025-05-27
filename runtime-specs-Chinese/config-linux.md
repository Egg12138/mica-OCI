# <a name="linuxContainerConfiguration" />Linux 容器配置

本文档描述了[容器配置](config.md)中[Linux特定部分](config.md#platform-specific-configuration)的架构。
Linux容器规范使用各种内核特性，如命名空间、cgroups、capabilities、LSM和文件系统隔离来实现规范。

## <a name="configLinuxDefaultFilesystems" />默认文件系统

Linux ABI 包括系统调用和几个特殊的文件路径。
期望在Linux环境中运行的应用程序很可能需要这些文件路径被正确设置。

以下文件系统应该在每个容器的文件系统中可用：

| 路径     | 类型   |
| -------- | ------ |
| /proc    | [proc][] |
| /sys     | [sysfs][]  |
| /dev/pts | [devpts][] |
| /dev/shm | [tmpfs][]  |

## <a name="configLinuxNamespaces" />命名空间

命名空间将全局系统资源包装在一个抽象层中，使得命名空间内的进程看起来拥有自己的全局资源隔离实例。
对全局资源的更改对作为命名空间成员的其他进程可见，但对其他进程不可见。
更多信息，请参见 [namespaces(7)][namespaces.7_2] 手册页。

命名空间在 `namespaces` 根字段内以数组形式指定。
可以指定以下参数来设置命名空间：

* **`type`** *(string, REQUIRED)* - 命名空间类型。应该支持以下命名空间类型：
    * **`pid`** 容器内的进程只能看到同一容器内或同一pid命名空间内的其他进程。
    * **`network`** 容器将拥有自己的网络栈。
    * **`mount`** 容器将拥有隔离的挂载表。
    * **`ipc`** 容器内的进程只能通过系统级IPC与同一容器内的其他进程通信。
    * **`uts`** 容器将能够拥有自己的主机名和域名。
    * **`user`** 容器将能够将主机上的用户和组ID重新映射到容器内的本地用户和组。
    * **`cgroup`** 容器将拥有cgroup层次结构的隔离视图。
    * **`time`** 容器将能够拥有自己的时钟。
* **`path`** *(string, OPTIONAL)* - 命名空间文件。
    该值必须是[runtime挂载命名空间](glossary.md#runtime-namespace)中的绝对路径。
    运行时必须将容器进程放在与该`path`关联的命名空间中。
    如果`path`与`type`类型的命名空间不关联，运行时必须[生成错误](runtime.md#errors)。

    如果未指定`path`，运行时必须创建类型为`type`的新[容器命名空间](glossary.md#container-namespace)。

如果`namespaces`数组中未指定命名空间类型，容器必须继承该类型的[runtime命名空间](glossary.md#runtime-namespace)。
如果`namespaces`字段包含具有相同`type`的重复命名空间，运行时必须[生成错误](runtime.md#errors)。

### 示例

```json
"namespaces": [
    {
        "type": "pid",
        "path": "/proc/1234/ns/pid"
    },
    {
        "type": "network",
        "path": "/var/run/netns/neta"
    },
    {
        "type": "mount"
    },
    {
        "type": "ipc"
    },
    {
        "type": "uts"
    },
    {
        "type": "user"
    },
    {
        "type": "cgroup"
    },
    {
        "type": "time"
    }
]
```

## <a name="configLinuxUserNamespaceMappings" />用户命名空间映射

**`uidMappings`** (对象数组，OPTIONAL) 描述了从主机到容器的用户命名空间uid映射。
**`gidMappings`** (对象数组，OPTIONAL) 描述了从主机到容器的用户命名空间gid映射。

每个条目具有以下结构：

* **`containerID`** *(uint32, REQUIRED)* - 是容器中的起始uid/gid。
* **`hostID`** *(uint32, REQUIRED)* - 是主机上要映射到*containerID*的起始uid/gid。
* **`size`** *(uint32, REQUIRED)* - 是要映射的id数量。

运行时不应该修改引用文件系统的所有权来实现映射。
注意，映射条目的数量可能受到[内核][user-namespaces]的限制。

### 示例

```json
"uidMappings": [
    {
        "containerID": 0,
        "hostID": 1000,
        "size": 32000
    }
],
"gidMappings": [
    {
        "containerID": 0,
        "hostID": 1000,
        "size": 32000
    }
]
```

## <a name="configLinuxTimeOffset" />时间命名空间的偏移量

**`timeOffsets`** (对象，OPTIONAL) 设置时间命名空间的偏移量。更多信息
请参见 [time_namespaces][time_namespaces.7]。

时钟的名称是条目的键。
条目值是具有以下属性的对象：

* **`secs`** *(int64, OPTIONAL)* - 是容器中时钟的偏移量（以秒为单位）。
* **`nanosecs`** *(uint32, OPTIONAL)* - 是容器中时钟的偏移量（以纳秒为单位）。

## <a name="configLinuxDevices" />设备

**`devices`** (对象数组，OPTIONAL) 列出了容器中必须可用的设备。
运行时可以以任何方式提供这些设备（通过[`mknod`][mknod.2]，从运行时挂载命名空间绑定挂载，使用符号链接等）。

每个条目具有以下结构：

* **`type`** *(string, REQUIRED)* - 设备类型：`c`、`b`、`u`或`p`。
    更多信息请参见 [mknod(1)][mknod.1]。
* **`path`** *(string, REQUIRED)* - 容器内设备的完整路径。
    如果`path`处已存在与请求的设备不匹配的[文件][]，运行时必须生成错误。
    路径可以是容器文件系统中的任何位置，特别是在`/dev`之外。
* **`major, minor`** *(int64, 除非`type`是`p`否则REQUIRED)* - 设备的[主设备号，次设备号][devices]。
* **`fileMode`** *(uint32, OPTIONAL)* - 设备的文件模式。
    您也可以通过cgroups[控制对设备的访问](#configLinuxDeviceAllowedlist)。
* **`uid`** *(uint32, OPTIONAL)* - [容器命名空间](glossary.md#container-namespace)中设备所有者的id。
* **`gid`** *(uint32, OPTIONAL)* - [容器命名空间](glossary.md#container-namespace)中设备组的id。

相同的`type`、`major`和`minor`不应该用于多个设备。

容器可能无法访问任何未在**`devices`**数组中明确引用或未列为[默认设备](#configLinuxDefaultDevices)一部分的设备节点。
理由：基于虚拟机的运行时需要能够调整节点设备，访问未调整的设备节点可能具有未定义的行为。

### 示例

```json
"devices": [
    {
        "path": "/dev/fuse",
        "type": "c",
        "major": 10,
        "minor": 229,
        "fileMode": 438,
        "uid": 0,
        "gid": 0
    },
    {
        "path": "/dev/sda",
        "type": "b",
        "major": 8,
        "minor": 0,
        "fileMode": 432,
        "uid": 0,
        "gid": 0
    }
]
```

### <a name="configLinuxDefaultDevices" />默认设备

除了通过此设置配置的任何设备外，运行时还必须提供：

* [`/dev/null`][null.4]
* [`/dev/zero`][zero.4]
* [`/dev/full`][full.4]
* [`/dev/random`][random.4]
* [`/dev/urandom`][random.4]
* [`/dev/tty`][tty.4]
* 如果在配置中启用了[`terminal`](config.md#process)，则通过将伪终端pty绑定挂载到`/dev/console`来设置`/dev/console`。
* [`/dev/ptmx`][pts.4]。
  容器`/dev/pts/ptmx`的[绑定挂载或符号链接][devpts]。

## <a name="configLinuxNetworkDevices" />网络设备

Linux网络设备是发送和接收数据包的实体。它们
不以文件形式存在于`/dev`目录中。相反，它们由
Linux内核中的[`net_device`][net_device]数据结构表示。网络
设备只能属于一个网络命名空间，并使用一组与常规文件操作不同的操作。
网络设备可以分为**物理**或**虚拟**：

* **物理网络设备**对应于硬件接口，如
    以太网卡（例如，`eth0`、`enp0s3`）。它们直接与
    物理网络硬件关联。
* **虚拟网络设备**是软件定义的接口，如回环
    设备（`lo`）、虚拟以太网对（`veth`）、网桥（`br0`）、VLAN和
    MACVLAN。它们由内核创建和管理，不对应于
    物理硬件。

此架构仅关注将主机网络命名空间中按名称识别的现有网络设备
移动到容器网络命名空间中。它不涵盖网络设备创建或网络配置的复杂性，
如IP地址分配、路由和DNS设置。

**`netDevices`** (对象，OPTIONAL) - 必须在容器中可用的一组网络设备。
运行时负责移动这些设备；底层机制由实现定义。

网络设备的名称是条目的键。条目值是具有
以下属性的对象：

* **`name`** *(string, OPTIONAL)* - 容器命名空间内网络设备的名称。
    如果未指定，则使用主机名称。

运行时必须检查将网络接口移动到容器
命名空间是否可能。如果容器命名空间中已存在具有指定名称的网络设备，
运行时必须[生成错误](runtime.md#errors)，除非用户通过在新名称后附加
`%d`提供了模板。在这种情况下，运行时必须允许移动，并且
内核将为容器网络命名空间内的接口生成唯一名称。

运行时必须保留现有网络接口属性，包括所有
具有全局作用域（RT_SCOPE_UNIVERSE值）的任何族的永久IP地址（IFA_F_PERMANENT标志），
如[`RFC 3549 Section 2.3.3.2`][rfc3549]中定义。这确保只有用于持久外部通信的
地址被转移。

运行时必须在将网络设备移动到网络命名空间后将其状态设置为"up"，
以允许容器通过该设备发送和接收网络流量。

### 命名空间生命周期和容器终止

运行时不得主动管理容器网络命名空间*内*的接口生命周期和配置。
这是因为网络接口本质上与网络命名空间本身绑定，因此它们的生命周期由
网络命名空间的所有者管理。通常，这种所有权和管理由更高级别的容器运行时
编排器处理，而不是直接在容器内运行的进程。

[proc]: https://www.kernel.org/doc/html/latest/filesystems/proc.html
[sysfs]: https://www.kernel.org/doc/html/latest/filesystems/sysfs.html
[devpts]: https://www.kernel.org/doc/html/latest/filesystems/devpts.html
[tmpfs]: https://www.kernel.org/doc/html/latest/filesystems/tmpfs.html
[namespaces.7_2]: https://man7.org/linux/man-pages/man7/namespaces.7.html
[user-namespaces]: https://www.kernel.org/doc/html/latest/admin-guide/namespaces/compatibility-list.html
[time_namespaces.7]: https://man7.org/linux/man-pages/man7/time_namespaces.7.html
[mknod.2]: https://man7.org/linux/man-pages/man2/mknod.2.html
[mknod.1]: https://man7.org/linux/man-pages/man1/mknod.1.html
[devices]: https://www.kernel.org/doc/html/latest/admin-guide/devices.html
[null.4]: https://man7.org/linux/man-pages/man4/null.4.html
[zero.4]: https://man7.org/linux/man-pages/man4/zero.4.html
[full.4]: https://man7.org/linux/man-pages/man4/full.4.html
[random.4]: https://man7.org/linux/man-pages/man4/random.4.html
[tty.4]: https://man7.org/linux/man-pages/man4/tty.4.html
[pts.4]: https://man7.org/linux/man-pages/man4/pts.4.html
[devpts]: https://www.kernel.org/doc/html/latest/filesystems/devpts.html
[net_device]: https://www.kernel.org/doc/html/latest/networking/netdevices.html
[rfc3549]: https://www.rfc-editor.org/rfc/rfc3549.html#section-2.3.3.2 
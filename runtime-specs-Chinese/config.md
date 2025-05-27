# 配置

此配置文件包含实现针对容器的[标准操作](runtime.md#operations)所需的元数据。
这包括要运行的进程、要注入的环境变量、要使用的沙箱功能等。

规范模式在本文档中定义，但在[`schema/config-schema.json`](schema/config-schema.json)中有JSON Schema，在[`specs-go/config.go`](specs-go/config.go)中有Go绑定。
[平台](spec.md#platforms)特定的配置模式在下面链接的[平台特定文档](#platform-specific-configuration)中定义。
对于仅在部分[平台](spec.md#platforms)上定义的属性，Go属性有一个`platform`标签列出这些协议（例如`platform:"linux,solaris"`）。

以下是配置格式中定义的每个字段的详细说明，并指定了有效值。
平台特定字段已标识为特定字段。
对于所有平台特定的配置值，下面[平台特定配置](#platform-specific-configuration)部分中定义的范围适用。

## 规范版本

* **`ociVersion`** (字符串，必需) 必须是[SemVer v2.0.0][semver-v2.0.0]格式，并指定bundle所遵循的Open Container Initiative运行时规范的版本。
    Open Container Initiative运行时规范遵循语义化版本控制，并在主要版本内保持向前和向后兼容性。
    例如，如果配置符合本规范的1.1版本，则它与支持本规范1.1或更高版本的任何运行时兼容，但与仅支持1.0而不支持1.1的运行时不兼容。

### 示例

```json
"ociVersion": "0.1.0"
```

## 根文件系统

**`root`** (对象，可选) 指定容器的根文件系统。
在Windows上，对于Windows Server容器，此字段是必需的。
对于[Hyper-V容器](config-windows.md#hyperv)，不得设置此字段。

在所有其他平台上，此字段是必需的。

* **`path`** (字符串，必需) 指定容器根文件系统的路径。
    * 在Windows上，`path`必须是[卷GUID路径][naming-a-volume]。
    * 在POSIX平台上，`path`是绝对路径或相对于bundle的路径。
        例如，对于位于`/to/bundle`的bundle和位于`/to/bundle/rootfs`的根文件系统，`path`值可以是`/to/bundle/rootfs`或`rootfs`。
        该值应该是常规的`rootfs`。

    在字段声明的路径处必须存在一个目录。

* **`readonly`** (布尔值，可选) 如果为true，则容器内的根文件系统必须是只读的，默认为false。
    * 在Windows上，此字段必须省略或为false。

### 示例（POSIX平台）

```json
"root": {
    "path": "rootfs",
    "readonly": true
}
```

## 挂载点

**`mounts`** (对象数组，可选) 指定除[`root`](#root)之外的额外挂载点。
运行时必须按列出的顺序挂载条目。
对于Linux，参数如[mount(2)][mount.2]系统调用手册页中所述。
对于Solaris，挂载条目对应于[zonecfg(1M)][zonecfg.1m]手册页中的'fs'资源。

* **`destination`** (字符串，必需) 挂载点的目标：容器内的路径。
    * Linux：此值应该是绝对路径。
      为了与旧工具和配置兼容，它可以是相对路径，在这种情况下，它必须被解释为相对于"/"。
      相对路径已**弃用**。
    * Windows：此值必须是绝对路径。
      一个挂载目标不得嵌套在另一个挂载中（例如，c:\\foo和c:\\foo\\bar）。
    * Solaris：此值必须是绝对路径。
      对应于[zonecfg(1M)][zonecfg.1m]中fs资源的"dir"。
    * 对于所有其他平台：此值必须是绝对路径。
* **`source`** (字符串，可选) 设备名称，但也可以是绑定挂载的文件或目录名称，或虚拟设备。
    绑定挂载的路径值要么是绝对的，要么是相对于bundle的。
    如果选项中有`bind`或`rbind`，则挂载是绑定挂载。
    * Windows：容器主机文件系统上的本地目录。不支持UNC路径和映射驱动器。
    * Solaris：对应于[zonecfg(1M)][zonecfg.1m]中fs资源的"special"。
* **`options`** (字符串数组，可选) 要使用的文件系统的挂载选项。
    * Linux：参见下面的[Linux挂载选项](#configLinuxMountOptions)。
    * Solaris：对应于[zonecfg(1M)][zonecfg.1m]中fs资源的"options"。
    * Windows：当给出`ro`时，运行时必须支持`ro`，以只读方式挂载文件系统。

### Linux挂载选项

运行时必须/应该/可以实现以下Linux选项字符串：

 选项名称      | 要求      | 描述
--------------|-----------|-----------------------------------------------------
 `async`      | 必须      | [^1]
 `atime`      | 必须      | [^1]
 `bind`       | 必须      | 绑定挂载 [^2]
 `defaults`   | 必须      | [^1]
 `dev`        | 必须      | [^1]
 `diratime`   | 必须      | [^1]
 `dirsync`    | 必须      | [^1]
 `exec`       | 必须      | [^1]
 `iversion`   | 必须      | [^1]
 `lazytime`   | 必须      | [^1]
 `loud`       | 必须      | [^1]
 `mand`       | 可以      | [^1]（在kernel 5.15，util-linux 2.38中弃用）
 `noatime`    | 必须      | [^1]
 `nodev`      | 必须      | [^1]
 `nodiratime` | 必须      | [^1]
 `noexec`     | 必须      | [^1]
 `noiversion` | 必须      | [^1]
 `nolazytime` | 必须      | [^1]
 `nomand`     | 可以      | [^1]
 `norelatime` | 必须      | [^1]
 `nostrictatime` | 必须   | [^1]
 `nosuid`     | 必须      | [^1]
 `nosymfollow`| 应该      | [^1]（在kernel 5.10，util-linux 2.38中引入）
 `private`    | 必须      | 绑定挂载传播 [^2]
 `ratime`     | 应该      | 递归`atime` [^3]
 `rbind`      | 必须      | 递归绑定挂载 [^2]
 `rdev`       | 应该      | 递归`dev` [^3]
 `rdiratime`  | 应该      | 递归`diratime` [^3]
 `relatime`   | 必须      | [^1]
 `remount`    | 必须      | [^1]
 `rexec`      | 应该      | 递归`dev` [^3]
 `rnoatime`   | 应该      | 递归`noatime` [^3]
 `rnodiratime`| 应该      | 递归`nodiratime` [^3]
 `rnoexec`    | 应该      | 递归`noexec` [^3]
 `rnorelatime`| 应该      | 递归`norelatime` [^3]
 `rnostrictatime` | 应该 | 递归`nostrictatime` [^3]
 `rnosuid`    | 应该      | 递归`nosuid` [^3]
 `rnosymfollow` | 应该    | 递归`nosymfollow` [^3]
 `ro`         | 必须      | [^1]
 `rprivate`   | 必须      | 绑定挂载传播 [^2]
 `rrelatime`  | 应该      | 递归`relatime` [^3]
 `rro`        | 应该      | 递归`ro` [^3]
 `rrw`        | 应该      | 递归`rw` [^3]
 `rshared`    | 必须      | 绑定挂载传播 [^2]
 `rslave`     | 必须      | 绑定挂载传播 [^2]
 `rstrictatime` | 应该    | 递归`strictatime` [^3]
 `rsuid`      | 应该      | 递归`suid` [^3]
 `rsymfollow` | 应该      | 递归`symfollow` [^3]
 `runbindable`| 必须      | 绑定挂载传播 [^2]
 `rw`         | 必须      | [^1]
 `shared`     | 必须      | [^1]
 `silent`     | 必须      | [^1]
 `slave`      | 必须      | 绑定挂载传播 [^2]
 `strictatime`| 必须      | [^1]
 `suid`       | 必须      | [^1]
 `symfollow`  | 应该      | `nosymfollow`的反向
 `sync`       | 必须      | [^1]
 `tmpcopyup`  | 可以      | 将内容复制到tmpfs
 `idmap`      | 应该      | 表示挂载必须应用idmapping。此选项不应传递给底层的[`mount(2)`][mount.2]调用。如果为挂载指定了`uidMappings`或`gidMappings`，运行时必须使用这些值进行挂载的映射。如果未指定，运行时可以使用容器的用户命名空间映射，否则必须[返回错误](runtime.md#errors)。如果未指定`uidMappings`和`gidMappings`且容器不使用用户命名空间，必须[返回错误](runtime.md#errors)。这应该使用[`mount_setattr(MOUNT_ATTR_IDMAP)`][mount_setattr.2]实现，自Linux 5.12起可用。
 `ridmap`     | 应该      | 表示挂载必须应用idmapping，并且映射是递归应用的[^3]。此选项不应传递给底层的[`mount(2)`][mount.2]调用。如果为挂载指定了`uidMappings`或`gidMappings`，运行时必须使用这些值进行挂载的映射。如果未指定，运行时可以使用容器的用户命名空间映射，否则必须[返回错误](runtime.md#errors)。如果未指定`uidMappings`和`gidMappings`且容器不使用用户命名空间，必须[返回错误](runtime.md#errors)。这应该使用[`mount_setattr(MOUNT_ATTR_IDMAP)`][mount_setattr.2]实现，自Linux 5.12起可用。

[^1]: 对应于[`mount(8)`（文件系统无关）][mount.8-filesystem-independent]。
[^2]: 对应于[绑定挂载和共享子树][mount-bind]。
[^3]: 这些`AT_RECURSIVE`选项需要kernel 5.12或更高版本。参见[`mount_setattr(2)`][mount_setattr.2]

"必须"选项对应于[`mount(8)`][mount.8]。

运行时也可以实现上表中未列出的自定义选项字符串。
如果自定义选项字符串已被[`mount(8)`][mount.8]识别，运行时应该遵循[`mount(8)`][mount.8]的行为。

运行时应该将未知选项视为[文件系统特定的选项][mount.8-filesystem-specific]，并将它们作为逗号分隔的字符串传递给[`mount(2)`][mount.2]的第五个（`const void *data`）参数。

### 示例（Linux）

```json
"mounts": [
    {
        "destination": "/tmp",
        "type": "tmpfs",
        "source": "tmpfs",
        "options": ["nosuid","strictatime","mode=755","size=65536k"]
    },
    {
        "destination": "/data",
        "type": "none",
        "source": "/volumes/testing",
        "options": ["rbind","rw"]
    }
]
```

## 进程

**`process`** (对象，可选) 指定容器进程。
当调用[`start`](runtime.md#start)时，此属性是必需的。

* **`terminal`** (布尔值，可选) 指定是否将终端附加到进程，默认为false。
    例如，如果在Linux上设置为true，则为进程分配一个伪终端对，并将伪终端pty复制到进程的[标准流][stdin.3]。
* **`consoleSize`** (对象，可选) 指定终端的控制台大小（以字符为单位）。
    如果`terminal`为`false`或未设置，运行时必须忽略`consoleSize`。
    * **`height`** (无符号整数，必需)
    * **`width`** (无符号整数，必需)
* **`cwd`** (字符串，必需) 是将为可执行文件设置的工作目录。
    此值必须是绝对路径。
* **`env`** (字符串数组，可选) 与[IEEE Std 1003.1-2008的`environ`][ieee-1003.1-2008-xbd-c8.1]具有相同的语义。 
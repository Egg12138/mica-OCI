# 功能结构

[运行时](glossary.md#runtime)可以向[运行时调用者](glossary.md#runtime-caller)提供关于其实现功能的JSON结构。
这个JSON结构被称为["功能结构"](glossary.md#features-structure)。

功能结构与主机操作系统中功能的实际可用性无关。
因此，功能结构的内容应该在运行时的编译时确定，而不是在执行时。

除了`ociVersionMin`和`ociVersionMax`之外，功能结构中的所有属性都可以不存在或具有`null`值。
`null`值不得与空值（如`0`、`false`、`""`、`[]`和`{}`）混淆。

## 规范版本

* **`ociVersionMin`** (字符串，必需) Open Container Initiative运行时规范的最低识别版本。
  运行时必须接受此值作为[`config.json`的`ociVersion`属性](config.md#specification-version)。

* **`ociVersionMax`** (字符串，必需) Open Container Initiative运行时规范的最高识别版本。
  运行时必须接受此值作为[`config.json`的`ociVersion`属性](config.md#specification-version)。
  该值不得小于`ociVersionMin`属性的值。
  功能结构不得包含未在此版本的Open Container Initiative运行时规范中定义的属性。

### 示例
```json
{
  "ociVersionMin": "1.0.0",
  "ociVersionMax": "1.1.0"
}
```

## 钩子
* **`hooks`** (字符串数组，可选) 识别的[钩子](config.md#posix-platform-hooks)名称。
  运行时必须支持此数组中的元素作为[`config.json`的`hooks`属性](config.md#posix-platform-hooks)。

### 示例
```json
"hooks": [
  "prestart",
  "createRuntime",
  "createContainer",
  "startContainer",
  "poststart",
  "poststop"
]
```

## 挂载选项

* **`mountOptions`** (字符串数组，可选) 识别的挂载选项名称，包括可能不被主机操作系统支持的选项。
  运行时必须将此数组中的元素识别为[`config.json`中`mounts`对象的`options`](config.md#mounts)。
  * Linux：此数组不应包含作为`const void *data`传递给[mount(2)][mount.2]系统调用的文件系统特定挂载选项。

### 示例

```json
"mountOptions": [
  "acl",
  "async",
  "atime",
  "bind",
  "defaults",
  "dev",
  "diratime",
  "dirsync",
  "exec",
  "iversion",
  "lazytime",
  "loud",
  "mand",
  "noacl",
  "noatime",
  "nodev",
  "nodiratime",
  "noexec",
  "noiversion",
  "nolazytime",
  "nomand",
  "norelatime",
  "nostrictatime",
  "nosuid",
  "nosymfollow",
  "private",
  "ratime",
  "rbind",
  "rdev",
  "rdiratime",
  "relatime",
  "remount",
  "rexec",
  "rnoatime",
  "rnodev",
  "rnodiratime",
  "rnoexec",
  "rnorelatime",
  "rnostrictatime",
  "rnosuid",
  "rnosymfollow",
  "ro",
  "rprivate",
  "rrelatime",
  "rro",
  "rrw",
  "rshared",
  "rslave",
  "rstrictatime",
  "rsuid",
  "rsymfollow",
  "runbindable",
  "rw",
  "shared",
  "silent",
  "slave",
  "strictatime",
  "suid",
  "symfollow",
  "sync",
  "tmpcopyup",
  "unbindable"
]
```

## 平台特定功能

* **`linux`** (对象，可选) [Linux特定功能](features-linux.md)。
  如果运行时支持`linux`平台，可以设置此项。

## 注解

**`annotations`** (对象，可选) 包含运行时的任意元数据。
此信息可以是结构化的或非结构化的。
注解必须是键值映射，遵循与[`config.json`的`annotations`属性](config.md#annotations)的键和值相同的约定。
但是，注解不需要包含[`config.json`的`annotations`属性](config.md#annotations)的可能值。
当前版本的规范没有提供枚举[`config.json`的`annotations`属性](config.md#annotations)的可能值的方法。

### 示例
```json
"annotations": {
  "org.opencontainers.runc.checkpoint.enabled": "true",
  "org.opencontainers.runc.version": "1.1.0"
}
```

## `config.json`中的不安全注解

**`potentiallyUnsafeConfigAnnotations`** (字符串数组，可选) 包含[`config.json`的`annotations`属性](config.md#annotations)的值，
这些值可能会潜在地改变运行时的行为。

以"."结尾的值被解释为注解的前缀。

### 示例
```json
"potentiallyUnsafeConfigAnnotations": [
  "com.example.foo.bar",
  "org.systemd.property."
]
```

上面的示例匹配`com.example.foo.bar`、`org.systemd.property.ExecStartPre`等。
该示例不匹配`com.example.foo.bar.baz`。

# 完整示例

以下是完整示例供参考。

```json
{
  "ociVersionMin": "1.0.0",
  "ociVersionMax": "1.1.0-rc.2",
  "hooks": [
    "prestart",
    "createRuntime",
    "createContainer",
    "startContainer",
    "poststart",
    "poststop"
  ],
  "mountOptions": [
    "async",
    "atime",
    "bind",
    "defaults",
    "dev",
    "diratime",
    "dirsync",
    "exec",
    "iversion",
    "lazytime",
    "loud",
    "mand",
    "noatime",
    "nodev",
    "nodiratime",
    "noexec",
    "noiversion",
    "nolazytime",
    "nomand",
    "norelatime",
    "nostrictatime",
    "nosuid",
    "nosymfollow",
    "private",
    "ratime",
    "rbind",
    "rdev",
    "rdiratime",
    "relatime",
    "remount",
    "rexec",
    "rnoatime",
    "rnodev",
    "rnodiratime",
    "rnoexec",
    "rnorelatime",
    "rnostrictatime",
    "rnosuid",
    "rnosymfollow",
    "ro",
    "rprivate",
    "rrelatime",
    "rro",
    "rrw",
    "rshared",
    "rslave",
    "rstrictatime",
    "rsuid",
    "rsymfollow",
    "runbindable",
    "rw",
    "shared",
    "silent",
    "slave",
    "strictatime",
    "suid",
    "symfollow",
    "sync",
    "tmpcopyup",
    "unbindable"
  ],
  "linux": {
    "namespaces": [
      "cgroup",
      "ipc",
      "mount",
      "network",
      "pid",
      "user",
      "uts"
    ],
    "capabilities": [
      "CAP_CHOWN",
      "CAP_DAC_OVERRIDE",
      "CAP_DAC_READ_SEARCH",
      "CAP_FOWNER",
      "CAP_FSETID",
      "CAP_KILL",
      "CAP_SETGID",
      "CAP_SETUID",
      "CAP_SETPCAP",
      "CAP_LINUX_IMMUTABLE",
      "CAP_NET_BIND_SERVICE",
      "CAP_NET_BROADCAST",
      "CAP_NET_ADMIN",
      "CAP_NET_RAW",
      "CAP_IPC_LOCK",
      "CAP_IPC_OWNER",
      "CAP_SYS_MODULE",
      "CAP_SYS_RAWIO",
      "CAP_SYS_CHROOT",
      "CAP_SYS_PTRACE",
      "CAP_SYS_PACCT",
      "CAP_SYS_ADMIN",
      "CAP_SYS_BOOT",
      "CAP_SYS_NICE",
      "CAP_SYS_RESOURCE",
      "CAP_SYS_TIME",
      "CAP_SYS_TTY_CONFIG",
      "CAP_MKNOD",
      "CAP_LEASE",
      "CAP_AUDIT_WRITE",
      "CAP_AUDIT_CONTROL",
      "CAP_SETFCAP",
      "CAP_MAC_OVERRIDE",
      "CAP_MAC_ADMIN",
      "CAP_SYSLOG",
      "CAP_WAKE_ALARM",
      "CAP_BLOCK_SUSPEND",
      "CAP_AUDIT_READ"
    ]
  }
}
``` 
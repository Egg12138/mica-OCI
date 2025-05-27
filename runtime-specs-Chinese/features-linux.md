# <a name="linuxFeatures" />Linux 功能结构

本文档描述了[功能结构](features.md)的[Linux特定部分](features.md#platform-specific-features)。

## <a name="linuxFeaturesNamespaces" />命名空间

* **`namespaces`** (字符串数组，OPTIONAL) 命名空间的识别名称，包括主机操作系统可能不支持的名称空间。
  运行时必须将此数组中的元素识别为[`config.json`中`linux.namespaces`对象的`type`](config-linux.md#namespaces)。

### 示例

```json
"namespaces": [
  "cgroup",
  "ipc",
  "mount",
  "network",
  "pid",
  "user",
  "uts"
]
```

## <a name="linuxFeaturesCapabilities" />Capabilities

* **`capabilities`** (字符串数组，OPTIONAL) capabilities的识别名称，包括主机操作系统可能不支持的功能。
  运行时必须将此数组中的元素识别为[`config.json`的`process.capabilities`对象](config.md#linux-process)。

### 示例

```json
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
  "CAP_AUDIT_READ",
  "CAP_PERFMON",
  "CAP_BPF",
  "CAP_CHECKPOINT_RESTORE"
]
```

## <a name="linuxFeaturesCgroup" />Cgroup

**`cgroup`** (对象，OPTIONAL) 表示运行时对cgroup管理器的实现状态。
与主机操作系统的cgroup版本无关。

* **`v1`** (布尔值，OPTIONAL) 表示运行时是否支持cgroup v1。
* **`v2`** (布尔值，OPTIONAL) 表示运行时是否支持cgroup v2。
* **`systemd`** (布尔值，OPTIONAL) 表示运行时是否支持系统范围的systemd cgroup管理器。
* **`systemdUser`** (布尔值，OPTIONAL) 表示运行时是否支持用户范围的systemd cgroup管理器。
* **`rdma`** (布尔值，OPTIONAL) 表示运行时是否支持RDMA cgroup控制器。

### 示例

```json
"cgroup": {
  "v1": true,
  "v2": true,
  "systemd": true,
  "systemdUser": true,
  "rdma": false
}
```

## <a name="linuxFeaturesSeccomp" />Seccomp

**`seccomp`** (对象，OPTIONAL) 表示运行时对seccomp的实现状态。
与主机操作系统的内核版本无关。

* **`enabled`** (布尔值，OPTIONAL) 表示运行时是否支持seccomp。
* **`actions`** (字符串数组，OPTIONAL) seccomp操作的识别名称。
  运行时必须将此数组中的元素识别为[`config.json`中`linux.seccomp`对象的`syscalls[].action`属性](config-linux.md#seccomp)。
* **`operators`** (字符串数组，OPTIONAL) seccomp操作符的识别名称。
  运行时必须将此数组中的元素识别为[`config.json`中`linux.seccomp`对象的`syscalls[].args[].op`属性](config-linux.md#seccomp)。
* **`archs`** (字符串数组，OPTIONAL) seccomp架构的识别名称。
  运行时必须将此数组中的元素识别为[`config.json`中`linux.seccomp`对象的`architectures`属性](config-linux.md#seccomp)。
* **`knownFlags`** (字符串数组，OPTIONAL) seccomp标志的识别名称。
  运行时必须将此数组中的元素识别为[`config.json`中`linux.seccomp`对象的`flags`属性](config-linux.md#seccomp)。
* **`supportedFlags`** (字符串数组，OPTIONAL) seccomp标志的识别和支持名称。
  由于某些标志不被当前内核和/或libseccomp支持，此列表可能是`knownFlags`的子集。
  运行时必须识别并支持此数组中的元素在[`config.json`中`linux.seccomp`对象的`flags`属性](config-linux.md#seccomp)中的使用。

### 示例

```json
"seccomp": {
  "enabled": true,
  "actions": [
    "SCMP_ACT_ALLOW",
    "SCMP_ACT_ERRNO",
    "SCMP_ACT_KILL",
    "SCMP_ACT_LOG",
    "SCMP_ACT_NOTIFY",
    "SCMP_ACT_TRACE",
    "SCMP_ACT_TRAP"
  ],
  "operators": [
    "SCMP_CMP_EQ",
    "SCMP_CMP_GE",
    "SCMP_CMP_GT",
    "SCMP_CMP_LE",
    "SCMP_CMP_LT",
    "SCMP_CMP_MASKED_EQ",
    "SCMP_CMP_NE"
  ],
  "archs": [
    "SCMP_ARCH_AARCH64",
    "SCMP_ARCH_ARM",
    "SCMP_ARCH_MIPS",
    "SCMP_ARCH_MIPS64",
    "SCMP_ARCH_MIPS64N32",
    "SCMP_ARCH_MIPSEL",
    "SCMP_ARCH_MIPSEL64",
    "SCMP_ARCH_MIPSEL64N32",
    "SCMP_ARCH_PPC",
    "SCMP_ARCH_PPC64",
    "SCMP_ARCH_PPC64LE",
    "SCMP_ARCH_S390",
    "SCMP_ARCH_S390X",
    "SCMP_ARCH_X32",
    "SCMP_ARCH_X86",
    "SCMP_ARCH_X86_64"
  ],
  "knownFlags": [
    "SECCOMP_FILTER_FLAG_LOG"
  ],
  "supportedFlags": [
    "SECCOMP_FILTER_FLAG_LOG"
  ]
}
```

## <a name="linuxFeaturesApparmor" />AppArmor

**`apparmor`** (对象，OPTIONAL) 表示运行时对AppArmor的实现状态。
与主机操作系统上AppArmor的可用性无关。

* **`enabled`** (布尔值，OPTIONAL) 表示运行时是否支持AppArmor。

### 示例

```json
"apparmor": {
  "enabled": true
}
```

## <a name="linuxFeaturesApparmor" />SELinux

**`selinux`** (对象，OPTIONAL) 表示运行时对SELinux的实现状态。
与主机操作系统上SELinux的可用性无关。

* **`enabled`** (布尔值，OPTIONAL) 表示运行时是否支持SELinux。

### 示例

```json
"selinux": {
  "enabled": true
}
```

## <a name="linuxFeaturesIntelRdt" />Intel RDT

**`intelRdt`** (对象，OPTIONAL) 表示运行时对Intel RDT的实现状态。
与主机操作系统上Intel RDT的可用性无关。

* **`enabled`** (布尔值，OPTIONAL) 表示运行时是否支持Intel RDT。
* **`schemata`** (布尔值，OPTIONAL) 表示是否支持
  ([`config.json`中`linux.intelRdt`的`schemata`字段](config-linux.md#intelrdt))。

### 示例

```json
"intelRdt": {
  "enabled": true,
  "schemata": true
}
```

## <a name="linuxFeaturesMountExtensions" />挂载扩展

**`mountExtensions`** (对象，OPTIONAL) 表示运行时是否支持某些挂载功能，与主机操作系统上这些功能的可用性无关。

* **`idmap`** (对象，OPTIONAL) 表示运行时是否支持使用挂载的`uidMappings`和`gidMappings`属性的idmap挂载。
  * **`enabled`** (布尔值，OPTIONAL) 表示如果提供了挂载的`uidMappings`和`gidMappings`属性，运行时是否解析并尝试使用它们。
    注意，运行时可能对id-mapped挂载支持有部分实现（例如只允许具有与容器用户命名空间匹配的映射的挂载，或只允许id-mapped绑定挂载）。
    在这种情况下，运行时仍必须将此值设置为`true`，以表明运行时识别`uidMappings`和`gidMappings`属性。

### 示例

```json
"mountExtensions": {
  "idmap":{
    "enabled": true
  }
}
```

## <a name="linuxFeaturesNetDevices" />网络设备

**`netDevices`** (对象，OPTIONAL) 表示运行时对Linux网络设备的实现状态。

* **`enabled`** (布尔值，OPTIONAL) 表示运行时是否支持将Linux网络设备移动到容器的网络命名空间的能力。

### 示例

```json
"netDevices": {
  "enabled": true
}
``` 
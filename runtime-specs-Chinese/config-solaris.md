# <a name="solarisApplicationContainerConfiguration" />Solaris 应用程序容器配置 (Solaris Application Container Configuration)

Solaris 应用程序容器可以使用以下属性进行配置，除了 milestone 之外，以下所有属性都映射到 [zonecfg(1M)][zonecfg.1m_2] 手册页中指定的属性。

## <a name="configSolarisMilestone" />milestone (milestone)
在容器内启动所需进程之前应该进入"online"状态的 SMF(Service Management Facility) FMRI。

**`milestone`** *(string, OPTIONAL)*

### 示例
```json
"milestone": "svc:/milestone/container:default"
```

## <a name="configSolarisLimitpriv" />limitpriv
此容器中任何进程可以获得的最大权限集。
该属性应该包含一个逗号分隔的权限集规范，如相应 Solaris 版本的 [priv_str_to_set(3C)][priv-str-to-set.3c] 手册页中所述。

**`limitpriv`** *(string, OPTIONAL)*

### 示例
```json
"limitpriv": "default"
```

## <a name="configSolarisMaxShmMemory" />maxShmMemory
允许此应用程序容器使用的最大共享内存量。
可以为这些数字中的每一个应用一个比例（K、M、G、T）（例如，1M 表示一兆字节）。
映射到 [zonecfg(1M)][zonecfg.1m_2] 手册页中的 `max-shm-memory`。

**`maxShmMemory`** *(string, OPTIONAL)*

### 示例
```json
"maxShmMemory": "512m"
```

## <a name="configSolarisCappedCpu" />cappedCPU
设置容器可以使用的 CPU 时间量的限制。
使用的单位转换为容器中所有用户线程可以使用的单个 CPU 的百分比，表示为分数（例如，.75）或混合数（整数和分数，例如，1.25）。
ncpu 值为 1 表示 CPU 的 100%，值为 1.25 表示 125%，.75 表示 75%，依此类推。
当有上限的容器内的项目有自己的上限时，最小值优先。
cappedCPU 映射到 [zonecfg(1M)][zonecfg.1m_2] 手册页中的 `capped-cpu`。

* **`ncpus`** *(string, OPTIONAL)*

### 示例
```json
"cappedCPU": {
    "ncpus": "8"
}
```

## <a name="configSolarisCappedMemory" />cappedMemory
此应用程序容器可以使用的内存的物理和交换上限。
可以为这些数字中的每一个应用一个比例（K、M、G、T）（例如，1M 表示一兆字节）。
cappedMemory 映射到 [zonecfg(1M)][zonecfg.1m_2] 手册页中的 `capped-memory`。

* **`physical`** *(string, OPTIONAL)*
* **`swap`** *(string, OPTIONAL)*

### 示例
```json
"cappedMemory": {
    "physical": "512m",
    "swap": "512m"
}
```

## <a name="configSolarisNetwork" />网络

### <a name="configSolarisAutomaticNetwork" />自动网络 (anet)
anet 被指定为一个数组，用于为 Solaris 应用程序容器设置网络。
anet 资源表示自动为应用程序容器创建网络资源。
区域管理守护进程 zoneadmd 是管理容器虚拟平台的主要进程。
守护进程的职责之一是创建和拆除容器的网络。
有关守护进程的更多信息，请参见 [zoneadmd(1M)][zoneadmd.1m] 手册页。
当启动这样的容器时，会自动为容器创建一个临时的 VNIC（虚拟网卡）。
当容器被拆除时，VNIC 会被删除。
以下属性可用于设置自动网络。
有关属性的其他信息，请查看相应 Solaris 版本的 [zonecfg(1M)][zonecfg.1m_2] 手册页。

* **`linkname`** *(string, OPTIONAL)* 为自动创建的 VNIC 数据链路指定一个名称。
* **`lowerLink`** *(string, OPTIONAL)* 指定将在其上创建 VNIC 的链路。
映射到 [zonecfg(1M)][zonecfg.1m_2] 手册页中的 `lower-link`。
* **`allowedAddress`** *(string, OPTIONAL)* 容器可以使用的 IP 地址集可能通过指定 `allowedAddress` 属性来限制。
    如果未指定 `allowedAddress`，则它们可以使用网络资源关联的物理接口上的任何 IP 地址。
    否则，当指定 `allowedAddress` 时，容器不能使用不在物理地址的 `allowedAddress` 列表中的 IP 地址。
    映射到 [zonecfg(1M)][zonecfg.1m_2] 手册页中的 `allowed-address`。
* **`configureAllowedAddress`** *(string, OPTIONAL)* 如果 `configureAllowedAddress` 设置为 true，则每次容器启动时都会自动在接口上配置 `allowedAddress` 指定的地址。
    当它设置为 false 时，`allowedAddress` 将不会在容器启动时配置。
    映射到 [zonecfg(1M)][zonecfg.1m_2] 手册页中的 `configure-allowed-address`。
* **`defrouter`** *(string, OPTIONAL)* OPTIONAL 默认路由器的值。
* **`macAddress`** *(string, OPTIONAL)* 根据指定的值或关键字设置 VNIC 的 MAC 地址。
    如果不是关键字，则将其解释为单播 MAC 地址。
    有关支持的关键字列表，请参阅相应 Solaris 版本的 [zonecfg(1M)][zonecfg.1m_2] 手册页。
    映射到 [zonecfg(1M)][zonecfg.1m_2] 手册页中的 `mac-address`。
* **`linkProtection`** *(string, OPTIONAL)* 使用逗号分隔的值启用一种或多种类型的链路保护。
    有关相应 Solaris 版本中支持的值，请参见 dladm(8) 中的保护属性。
    映射到 [zonecfg(1M)][zonecfg.1m_2] 手册页中的 `link-protection`。

#### 示例
```json
"anet": [
    {
        "allowedAddress": "172.17.0.2/16",
        "configureAllowedAddress": "true",
        "defrouter": "172.17.0.1/16",
        "linkProtection": "mac-nospoof, ip-nospoof",
        "linkname": "net0",
        "lowerLink": "net2",
        "macAddress": "02:42:f8:52:c7:16"
    }
]
```

[priv-str-to-set.3c]: https://docs.oracle.com/cd/E86824_01/html/E54766/priv-str-to-set-3c.html
[zoneadmd.1m]: https://docs.oracle.com/cd/E86824_01/html/E54764/zoneadmd-1m.html
[zonecfg.1m_2]: https://docs.oracle.com/cd/E86824_01/html/E54764/zonecfg-1m.html 
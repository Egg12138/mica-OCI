# <a name="windowsSpecificContainerConfiguration" />Windows 特定容器配置 (Windows-specific Container Configuration)

本文档描述了[容器配置](config.md)中[平台特定配置](config.md#platform-specific-configuration)的 Windows 特定部分的模式。
Windows 容器规范使用 Windows 主机计算服务(HCS)提供的 API 来实现规范。

## <a name="configWindowsLayerFolders" />层文件夹 (LayerFolders)

**`layerFolders`** (字符串数组，必需) 指定容器镜像所依赖的层文件夹列表。列表从最顶层到基础层排序，最后一个条目是 scratch。
`layerFolders` 必须至少包含一个条目。

### 示例

```json
"windows": {
    "layerFolders": [
        "C:\\Layers\\layer2",
        "C:\\Layers\\layer1",
        "C:\\Layers\\layer-base",
        "C:\\scratch",
    ]
}
```

## <a name="configWindowsDevices" />设备 (Devices)

**`devices`** (对象数组，可选) 列出必须在容器中可用的设备。

每个条目具有以下结构：

* **`id`** *(string, 必需)* - 指定运行时必须在容器中提供的设备。
* **`idType`** *(string, 必需)* - 告诉运行时如何解释 `id`。目前，Windows 仅支持值 `class`，它将 `id` 标识为[设备接口类 GUID][interfaceGUID]。

[interfaceGUID]: https://docs.microsoft.com/en-us/windows-hardware/drivers/install/overview-of-device-interface-classes

### 示例

```json
"windows": {
    "devices": [
        {
            "id": "24E552D7-6523-47F7-A647-D3465BF1F5CA",
            "idType": "class"
        },
        {
            "id": "5175d334-c371-4806-b3ba-71fd53c9258d",
            "idType": "class"
        }
    ]
}
```

## <a name="configWindowsResources" />资源 (Resources)

您可以通过 Windows 配置的可选 `resources` 字段来配置容器的资源限制。

### <a name="configWindowsMemory" />内存

`memory` 是容器内存使用的可选配置。

可以指定以下参数：

* **`limit`** *(uint64, 可选)* - 设置内存使用限制（以字节为单位）。

#### 示例

```json
"windows": {
    "resources": {
        "memory": {
            "limit": 2097152
        }
    }
}
```

### <a name="configWindowsCpu" />CPU

`cpu` 是容器 CPU 使用的可选配置。

可以指定以下参数（互斥）：

* **`count`** *(uint64, 可选)* - 指定容器可用的 CPU 数量。它表示容器中配置的处理器 `count` 相对于主机可用处理器的比例。该比例最终决定了容器中的线程在每个调度间隔期间可以使用的处理器周期部分，以每 10,000 个周期的周期数表示。
* **`shares`** *(uint16, 可选)* - 限制容器相对于处理器上其他工作负载的处理器时间份额。处理器 `shares`（平台级别的 `weight`）是 0 到 10,000 之间的值。
* **`maximum`** *(uint16, 可选)* - 确定容器中的线程在每个调度间隔期间可以使用的处理器周期部分，以每 10,000 个周期的周期数表示。将处理器 `maximum` 设置为百分比乘以 100。
* **`affinity`** *(对象数组, 可选)* - 指定要为此容器关联的 CPU 集合。

  每个条目具有以下结构：

  参考：https://learn.microsoft.com/en-us/windows-hardware/drivers/ddi/miniport/ns-miniport-_group_affinity

  * **`mask`** *(uint64, 必需)* - 指定相对于此 CPU 组的 CPU 掩码。
  * **`group`** *(uint32, 必需)* - 指定此掩码引用的处理器组，由 GetLogicalProcessorInformationEx 返回。

参考：https://docs.microsoft.com/en-us/virtualization/api/hcs/schemareference#Container_Processor

#### 示例

```json
"windows": {
    "resources": {
        "cpu": {
            "maximum": 5000
        }
    }
}
```

### <a name="configWindowsStorage" />存储

`storage` 是容器存储使用的可选配置。

可以指定以下参数：

* **`iops`** *(uint64, 可选)* - 指定容器系统驱动器的每秒最大 IO 操作数。
* **`bps`** *(uint64, 可选)* - 指定容器系统驱动器的每秒最大字节数。
* **`sandboxSize`** *(uint64, 可选)* - 指定系统驱动器的最小大小（以字节为单位）。

#### 示例

```json
"windows": {
    "resources": {
        "storage": {
            "iops": 50
        }
    }
}
```

## <a name="configWindowsNetwork" />网络 (Network)

您可以通过 Windows 配置的可选 `network` 字段来配置容器的网络选项。

可以指定以下参数：

* **`endpointList`** *(字符串数组, 可选)* - 容器应该连接的 HNS（主机网络服务）端点列表。
* **`allowUnqualifiedDNSQuery`** *(bool, 可选)* - 指定是否允许不合格的 DNS 名称解析。
* **`DNSSearchList`** *(字符串数组, 可选)* - 用于名称解析的 DNS 后缀的逗号分隔列表。
* **`networkSharedContainerName`** *(string, 可选)* - 我们将与之共享网络堆栈的容器的名称（ID）。
* **`networkNamespace`** *(string, 可选)* - 将用于容器的网络命名空间的名称（ID）。如果指定了网络命名空间，则不得指定其他参数。

### 示例

```json
"windows": {
    "network": {
        "endpointList": [
            "7a010682-17e0-4455-a838-02e5d9655fe6"
        ],
        "allowUnqualifiedDNSQuery": true,
        "DNSSearchList": [
            "a.com",
            "b.com"
        ],
        "networkSharedContainerName": "containerName",
        "networkNamespace": "168f3daf-efc6-4377-b20a-2c86764ba892"
    }
}
```

## <a name="configWindowsCredentialSpec" />凭据规范 (Credential Spec)

您可以通过 Windows 配置的可选 `credentialSpec` 字段来配置容器的组托管服务帐户(gMSA)。
`credentialSpec` 是一个 JSON 对象，其属性由实现定义。
有关 gMSA 的更多信息，请参阅[Windows 容器的 Active Directory 服务帐户][gMSAOverview]。
有关生成 gMSA 的工具的更多信息，请参阅[部署概述][gMSATooling]。

[gMSAOverview]: https://aka.ms/windowscontainers/manage-serviceaccounts
[gMSATooling]: https://aka.ms/windowscontainers/credentialspec-tools

## <a name="configWindowsServicing" />服务 (Servicing)

当容器终止时，主机计算服务会指示是否有 Windows 更新服务操作待处理。
您可以通过 Windows 配置的可选 `servicing` 字段来指示容器应该以应用待处理服务操作的模式启动。

### 示例

```json
"windows": {
    "servicing": true
}
```

## <a name="configWindowsIgnoreFlushesDuringBoot" />启动期间忽略刷新 (IgnoreFlushesDuringBoot)

您可以通过 Windows 配置的可选 `ignoreFlushesDuringBoot` 字段来指示容器应该以在容器启动期间不执行磁盘刷新的模式启动。

### 示例

```json
"windows": {
    "ignoreFlushesDuringBoot": true
}
```

## <a name="configWindowsHyperV" />Hyper-V (HyperV)

`hyperv` 是 Windows 配置的可选字段。
如果存在，容器必须使用 Hyper-V 隔离运行。
如果省略，容器必须作为 Windows Server 容器运行。

可以指定以下参数：

* **`utilityVMPath`** *(string, 可选)* - 指定用于实用程序 VM 的镜像路径。
    如果使用不包含实用程序 VM 镜像的基础镜像，则需要指定此路径。
    如果未提供，运行时将从最底层开始向上搜索容器文件系统层，直到找到 "UtilityVM"，并默认使用该路径。

### 示例

```json
"windows": {
    "hyperv": {
        "utilityVMPath": "C:\\path\\to\\utilityvm"
    }
}
``` 
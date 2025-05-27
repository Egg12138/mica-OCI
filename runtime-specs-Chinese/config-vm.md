# <a name="VirtualMachineSpecificContainerConfiguration" />虚拟机特定容器配置 (Virtual-machine-specific Container Configuration)

本节描述了[容器配置](config.md)中[平台特定配置](config.md#platform-specific-configuration)的虚拟机特定部分的模式。
虚拟机容器规范为管理程序、内核和镜像提供了额外的配置。

## <a name="HypervisorObject" />管理程序对象 (Hypervisor Object)

**`hypervisor`** (对象，可选) 指定管理容器虚拟机的管理程序的详细信息。
* **`path`** (string, 必需) 管理容器虚拟机的管理程序二进制文件的路径。
    此值必须是[运行时挂载命名空间](glossary.md#runtime-namespace)中的绝对路径。
* **`parameters`** (字符串数组，可选) 指定要传递给管理程序的参数数组。

### 示例

```json
    "hypervisor": {
        "path": "/path/to/vmm",
        "parameters": ["opts1=foo", "opts2=bar"]
    }
```

## <a name="KernelObject" />内核对象 (Kernel Object)

**`kernel`** (对象，必需) 指定用于启动容器虚拟机的内核的详细信息。
* **`path`** (string, 必需) 用于启动容器虚拟机的内核的路径。
    此值必须是[运行时挂载命名空间](glossary.md#runtime-namespace)中的绝对路径。
* **`parameters`** (字符串数组，可选) 指定要传递给内核的参数数组。
* **`initrd`** (string, 可选) 容器虚拟机使用的初始 ramdisk 的路径。
    此值必须是[运行时挂载命名空间](glossary.md#runtime-namespace)中的绝对路径。

### 示例

```json
    "kernel": {
        "path": "/path/to/vmlinuz",
        "parameters": ["foo=bar", "hello world"],
        "initrd": "/path/to/initrd.img"
    }
```

## <a name="ImageObject" />镜像对象 (Image Object)

**`image`** (对象，可选) 指定包含容器虚拟机根文件系统的镜像的详细信息。
* **`path`** (string, 必需) 容器虚拟机根镜像的路径。
    此值必须是[运行时挂载命名空间](glossary.md#runtime-namespace)中的绝对路径。
* **`format`** (string, 必需) 容器虚拟机根镜像的格式。常用的支持格式有：
    * **`raw`** [原始磁盘镜像格式][raw-image-format]。未设置的 `format` 值将默认为该格式。
    * **`qcow2`** [QEMU 镜像格式][qcow2-image-format]。
    * **`vdi`** [VirtualBox 1.1 兼容镜像格式][vdi-image-format]。
    * **`vmdk`** [VMware 兼容镜像格式][vmdk-image-format]。
    * **`vhd`** [虚拟硬盘镜像格式][vhd-image-format]。

此镜像包含虚拟机 **`kernel`** 将启动到的根文件系统，不要与容器根文件系统本身混淆。后者由[根配置](config.md#Root-Configuration)部分中的 **`path`** 指定，将在虚拟机内部由基于虚拟机的运行时选择的位置挂载。

### 示例

```json
    "image": {
        "path": "/path/to/vm/rootfs.img",
        "format": "raw"
    }
```

[raw-image-format]: https://en.wikipedia.org/wiki/IMG_(file_format)
[qcow2-image-format]: https://git.qemu.org/?p=qemu.git;a=blob_plain;f=docs/interop/qcow2.txt;hb=HEAD
[vdi-image-format]: https://forensicswiki.org/wiki/Virtual_Disk_Image_(VDI)
[vmdk-image-format]: http://www.vmware.com/app/vmdk/?src=vmdk
[vhd-image-format]: https://github.com/libyal/libvhdi/blob/master/documentation/Virtual%20Hard%20Disk%20(VHD)%20image%20format.asciidoc 
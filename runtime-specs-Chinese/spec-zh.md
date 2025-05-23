# <a name="openContainerInitiativeRuntimeSpecification" />Open Container Initiative Runtime Specification

[Open Container Initiative][oci] 为操作系统进程和应用程序容器开发标准规范。

# <a name="ociRuntimeSpecAbstract" />Abstract

Open Container Initiative 运行时规范旨在规定容器的配置、执行环境和生命周期。

容器的配置被指定为支持平台的 `config.json`，详细说明了启用容器创建的字段。
执行环境的规范确保在容器内运行的应用程序在运行时之间具有一致的环境，并为容器的生命周期定义了通用操作。

# <a name="ociRuntimeSpecPlatforms" />Platforms

本规范定义的平台包括：

* `linux`: [runtime.md](runtime.md), [config.md](config.md), [features.md](features.md), [config-linux.md](config-linux.md), [runtime-linux.md](runtime-linux.md), 和 [features-linux.md](features-linux.md)。
* `solaris`: [runtime.md](runtime.md), [config.md](config.md), [features.md](features.md), 和 [config-solaris.md](config-solaris.md)。
* `windows`: [runtime.md](runtime.md), [config.md](config.md), [features.md](features.md), 和 [config-windows.md](config-windows.md)。
* `vm`: [runtime.md](runtime.md), [config.md](config.md), [features.md](features.md), 和 [config-vm.md](config-vm.md)。
* `zos`: [runtime.md](runtime.md), [config.md](config.md), [features.md](features.md), 和 [config-zos.md](config-zos.md)。

# <a name="ociRuntimeSpecTOC" />Table of Contents

- [介绍](spec.md)
    - [符号约定](#notational-conventions)
    - [容器原则](principles.md)
- [文件系统包](bundle.md)
- [运行时和生命周期](runtime.md)
    - [Linux特定的运行时和生命周期](runtime-linux.md)
- [配置](config.md)
    - [Linux特定的配置](config-linux.md)
    - [Solaris特定的配置](config-solaris.md)
    - [Windows特定的配置](config-windows.md)
    - [虚拟机特定的配置](config-vm.md)
    - [z/OS特定的配置](config-zos.md)
- [功能结构](features.md)
    - [Linux特定的功能结构](features-linux.md)
- [术语表](glossary.md)

# <a name="ociRuntimeSpecNotationalConventions" />Notational Conventions

关键词"MUST"、"MUST NOT"、"REQUIRED"、"SHALL"、"SHALL NOT"、"SHOULD"、"SHOULD NOT"、"RECOMMENDED"、"NOT RECOMMENDED"、"MAY"和"OPTIONAL"的解释如[RFC 2119][rfc2119]中所述。

关键词"unspecified"、"undefined"和"implementation-defined"的解释如[C99标准原理][c99-unspecified]中所述。

如果实现未能满足其实现的[平台](#platforms)的一个或多个MUST、REQUIRED或SHALL要求，则该实现在给定的CPU架构上不符合规范。
如果实现满足其实现的[平台](#platforms)的所有MUST、REQUIRED和SHALL要求，则该实现在给定的CPU架构上符合规范。


[c99-unspecified]: https://www.open-std.org/jtc1/sc22/wg14/www/C99RationaleV5.10.pdf#page=18
[oci]: https://opencontainers.org
[rfc2119]: https://www.rfc-editor.org/rfc/rfc2119.html 
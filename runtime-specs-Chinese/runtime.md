# 运行时和生命周期

## 容器的范围

使用运行时创建容器的实体必须能够使用本规范中定义的操作来操作该容器。其他使用相同或不同运行时实例的实体是否能看到该容器不在本规范的范围内。

## 状态

容器的状态包含以下属性：

* **`ociVersion`** (字符串，必需) 是容器状态所遵循的Open Container Initiative运行时规范的版本。
* **`id`** (字符串，必需) 是容器的ID。
    该ID在主机上的所有容器中必须是唯一的。
    不要求在不同主机之间保持唯一。
* **`status`** (字符串，必需) 是容器的运行时状态。
    值可以是以下之一：

    * `creating`：容器正在创建中（生命周期中的第2步）
    * `created`：运行时已完成[创建操作](#create)（生命周期第2步之后），且容器进程既未退出也未执行用户指定的程序
    * `running`：容器进程已执行用户指定的程序但尚未退出（生命周期第8步之后）
    * `stopped`：容器进程已退出（生命周期第10步）

    运行时可以定义其他值，但这些值必须用于表示上述未定义的新运行时状态。
* **`pid`** (整数，在Linux上当`status`为`created`或`running`时必需，在其他平台上可选) 是容器进程的ID。
  对于在运行时命名空间中执行的钩子，它是运行时看到的pid。
  对于在容器命名空间中执行的钩子，它是容器看到的pid。
* **`bundle`** (字符串，必需) 是容器bundle目录的绝对路径。
    提供此路径是为了让使用者能够在主机上找到容器的配置和根文件系统。
* **`annotations`** (映射，可选) 包含与容器关联的注解列表。
    如果未提供注解，则此属性可以不存在或为空映射。

状态可以包含其他属性。

当以JSON格式序列化时，格式必须符合JSON Schema [`schema/state-schema.json`](schema/state-schema.json)。

有关检索容器状态的信息，请参见[查询状态](#query-state)。

### 示例

```json
{
    "ociVersion": "0.2.0",
    "id": "oci-container1",
    "status": "running",
    "pid": 4422,
    "bundle": "/containers/redis",
    "annotations": {
        "myKey": "myValue"
    }
}
```

## 生命周期

生命周期描述了从容器创建到停止存在期间发生的事件时间线。

1. 调用符合OCI规范的运行时的[`create`](runtime.md#create)命令，提供bundle位置的引用和唯一标识符。
2. 必须根据[`config.json`](config.md)中的配置创建容器的运行时环境。
    如果运行时无法创建[`config.json`](config.md)中指定的环境，它必须[生成错误](#errors)。
    虽然必须创建[`config.json`](config.md)中请求的资源，但此时不得运行用户指定的程序（来自[`process`](config.md#process)）。
    在此步骤之后对[`config.json`](config.md)的任何更新都不得影响容器。
3. 运行时必须调用[`prestart`钩子](config.md#prestart)。
    如果任何`prestart`钩子失败，运行时必须[生成错误](#errors)，停止容器，并在第12步继续生命周期。
4. 运行时必须调用[`createRuntime`钩子](config.md#createRuntime-hooks)。
    如果任何`createRuntime`钩子失败，运行时必须[生成错误](#errors)，停止容器，并在第12步继续生命周期。
5. 运行时必须调用[`createContainer`钩子](config.md#createContainer-hooks)。
    如果任何`createContainer`钩子失败，运行时必须[生成错误](#errors)，停止容器，并在第12步继续生命周期。
6. 使用容器的唯一标识符调用运行时的[`start`](runtime.md#start)命令。
7. 运行时必须调用[`startContainer`钩子](config.md#startContainer-hooks)。
    如果任何`startContainer`钩子失败，运行时必须[生成错误](#errors)，停止容器，并在第12步继续生命周期。
8. 运行时必须运行用户指定的程序，如[`process`](config.md#process)所指定。
9. 运行时必须调用[`poststart`钩子](config.md#poststart)。
    如果任何`poststart`钩子失败，运行时必须[记录警告](#warnings)，但其余钩子和生命周期继续执行，就像钩子成功了一样。
10. 容器进程退出。
    这可能由于出错、退出、崩溃或调用运行时的[`kill`](runtime.md#kill)操作而发生。
11. 使用容器的唯一标识符调用运行时的[`delete`](runtime.md#delete)命令。
12. 必须通过撤销创建阶段（第2步）执行的步骤来销毁容器。
13. 运行时必须调用[`poststop`钩子](config.md#poststop)。
    如果任何`poststop`钩子失败，运行时必须[记录警告](#warnings)，但其余钩子和生命周期继续执行，就像钩子成功了一样。

## 错误

在指定操作生成错误的情况下，本规范不强制要求如何甚至是否向实现用户返回或暴露该错误。
除非另有说明，生成错误必须使环境状态保持为操作从未尝试过的状态 - 除了可能的琐碎辅助更改，如日志记录。

## 警告

在指定操作记录警告的情况下，本规范不强制要求如何甚至是否向实现用户返回或暴露该警告。
除非另有说明，记录警告不会改变操作的流程；它必须继续执行，就像没有记录警告一样。

## 操作

除非另有说明，运行时必须支持以下操作。

注意：这些操作不指定任何命令行API，参数是通用操作的输入。

### 查询状态

`state <container-id>`

如果未提供容器的ID，此操作必须[生成错误](#errors)。
尝试查询不存在的容器必须[生成错误](#errors)。
此操作必须返回[状态](#state)部分中指定的容器状态。

### 创建

`create <container-id> <path-to-bundle>`

如果未提供bundle的路径和要关联的容器ID，此操作必须[生成错误](#errors)。
如果提供的ID在运行时的范围内不是唯一的，或以任何其他方式无效，实现必须[生成错误](#errors)，并且不得创建新容器。
此操作必须创建新容器。

必须应用[`config.json`](config.md)中配置的所有属性，除了[`process`](config.md#process)。
直到由[`start`](#start)操作触发时，才应用[`process.args`](config.md#process)。
此操作可以应用其余的`process`属性。
如果运行时无法应用[配置](config.md)中指定的属性，它必须[生成错误](#errors)，并且不得创建新容器。

运行时可以在创建容器之前（[第2步](#lifecycle)）根据本规范验证`config.json`，可以是通用的，也可以针对本地系统功能进行验证。
对预创建验证感兴趣的[运行时调用者](glossary.md#runtime-caller)可以在调用创建操作之前运行[bundle验证工具](implementations.md#testing--tools)。

在此操作之后对[`config.json`](config.md)文件的任何更改都不会影响容器。

### 启动
`start <container-id>`

如果未提供容器ID，此操作必须[生成错误](#errors)。
尝试`start`不是[`created`](#state)状态的容器不得对容器产生任何影响，并且必须[生成错误](#errors)。
此操作必须运行[`process`](config.md#process)指定的用户指定程序。
如果未设置`process`，此操作必须生成错误。

### 终止
`kill <container-id> <signal>`

如果未提供容器ID，此操作必须[生成错误](#errors)。
尝试向既不是[`created`也不是`running`](#state)状态的容器发送信号不得对容器产生任何影响，并且必须[生成错误](#errors)。
此操作必须向容器进程发送指定的信号。

### 删除
`delete <container-id>`

如果未提供容器ID，此操作必须[生成错误](#errors)。
尝试`delete`不是[`stopped`](#state)状态的容器不得对容器产生任何影响，并且必须[生成错误](#errors)。
删除容器必须删除在`create`步骤期间创建的资源。
注意，与容器关联但不由该容器创建的资源不得被删除。
一旦容器被删除，其ID可以被后续容器使用。

## 钩子
本规范中指定的许多操作都有"钩子"，允许在每个操作之前或之后执行其他操作。
有关更多信息，请参见[运行时配置中的钩子](./config.md#posix-platform-hooks)。 
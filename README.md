# 设计

实现一个 mica runtime 作为 low-level 的“容器”运行时，相当于一个 runc

然而我们并不真正地管理容器，而是通过 mica hub 来进行对RTOS的管理工作

## mica runtime based-in Shim V2

## mica runtime as a runc drop-in replacement

## mica runtime works on both shim V2 and drop-in replacement

理论可行，但是LGTM. 我们会在 mica-both 中这样尝试

## mica runtime for isulad-shim

# 目录说明

* shimv2-c: 基于C的、对接containerd shimv2 的mica runtime 实现(便于与mica对接)，其中, shim 的实现不一定基于C。
* shimv2: 基于go的、对接containerd shimv2 的mica runtime 实现, 我们目前都在 containerd源码树下实现，软连接过来，很快就会去掉这个方式
* runmica: 基于go的mica runtime 实现, 直接作为 runc drop-in replacement, 同样暂时都是软连接runc的源码树，很快会
* runmicars: 基于rust的 mica runtime实现, 直接作为 runc drop-in replacement, 以youki 为参考实现

## TODO

# 调研

* runc
* runv
* kata

# Roadmap

* 网络
* 完善的生命周期
* 原子化服务



# 功能验证 (前、中端)

## OCI 接口验证

* docker <cmd> --runtime=rmica <args>

## rmica 与 micad 的对接模拟验证

### 需要相关组件

* rmica 二进制， 作为runc drop-in replacement
  * mica module (or communication module)

## 模拟启动一个非标准OCI-spec的服务

### 需要相关组件

* docker 用来启动请求
* rmica 二进制， 作为runc drop-in replacement
  * 添加一个 manual_test option， 手动发送一系列的测试语句;
* pseudo-mica daemon, 一个micad server模拟器
  * 监听 /run/micad.sock
  * 把来自rmica的请求转发给 mcs_task
  * 
* mcs_task, 一个二进制，模拟 micad 的控制对象
  * 


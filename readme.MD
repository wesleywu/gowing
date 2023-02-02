# Gowing 为 Golang 开发插上翅膀

Golang开箱即用的万能代码生成器，支持生成Restful API、微服务、管理中后台、React 前端UI

## 开发计划
### 1. CLI

#### 1.1 项目生成CLI
    gw init web|rpc_provider|rpc_consumer
        Restful Web API
        RPC Provider
        Web RPC Consumer

#### 1.2  代码生成CLI
    gw gen table_def|crud_web|crud_rpc|web_rpc
        Table Definition Yaml
        CRUD Service & Restful Web API
        CRUD Service & RPC Provider
        Restful Web API & RPC Consumer

### 2. 参考实现代码
#### 2.1 CRUD Service
- 【完成】查询接口支持Query By Example
- 【完成】Query By Example 支持 in/between/like/isnull查询，支持not修饰
- 【完成】Query By Example支持自定义Condition查询
- 【P1】支持关联查询
- 【P1】支持虚拟字段
- 【P1】支持级联查询
- 【P2】Query By Example支持Or查询（通过Condition扩展）
- 【P1】支持规则化缓存中间层，中间层不依赖缓存实现
- 【P1】缓存中间层缺省支持Redis协议
#### 2.2 Restful Web API
- 【完成】Web API Route绑定CRUD Service
- 支持完整输入校验
    - 非空
    - 整形
    - 浮点型
    - 字符串长度
    - 日期时间
    - bit和枚举类型
- 支持错误处理
    - 对于底层error，返回 Status 500 Internal Server Error
    - 对于输入参数校验失败，返回 Status 400 Bad Request
#### 2.3 React Web UI
- 支持Arco-Design控件
    - 自动补全 AutoComplete
    - 级联选择 Cascader
    - 复选框 Checkbox
    - 日期选择器 DatePicker
    - 表单 Form
    - 【完成】输入框 Input
    - 标签输入框 InputTag
    - 数字输入框 InputNumber
    - 提及 Mentions
    - 单选框 Radio
    - 评分 Rate
    - 选择器 Select
    - 滑动输入条 Slider
    - 开关 Switch
    - 时间选择器 TimePicker
    - 数据穿梭框 Transfer
    - 树选择 TreeSelect
    - 上传 Upload
      2.4 RPC Provider
- 【完成】支持多个Service注册到同一个Provider服务
- 【完成】字段支持google.protobuf.Any类型
    - 支持Any字段的wrap和unwrap工具方法
- 【完成】字段支持Condition复合条件类型
- 【完成】支持dubbo-go协议
- 【完成】支持nacos注册中心
- 【完成】支持对protobuf定义的message添加自定义tag，例如v（数据校验）, json（通常用于去掉缺省存在的omitempty）等
- 【完成】支持输入结构体的validation
- 【完成】支持为传入参数 Struct 设定 gmeta.Meta （可用于在orm操作时忽略 nil 值的字段）
- 【完成】支持与Pixiu API网关的集成
- 【完成】支持 rpc call 的声明式透明缓存（对业务代码无侵入）
- 【完成】支持链路跟踪
- 【P1】更完善的错误处理
- 【P1】支持ConfigCenter
- 【P1】支持Stream形式的调用
- 【P2】支持etcd注册中心
- 【P2】健康检测
- 【P2】支持国际化
- 【P3】支持go-zero
#### 2.5 RPC Consumer
- 【P1】启动自动连接RPC Provider
- 【P2】断线重连
- 【P2】健康检测

### 3. 前后端分离的中后台管理系统 (前端用React实现)

### 4. 代码生成模板

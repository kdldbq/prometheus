

# 改造
## rule 规则-从DB加载
参考《深入浅出 Prometheus：原理、应用、源码与拓展详情》11.6 源码改造-MySQL规则存储
对应配置项，在 prometheus.yml 的 global 段中

## 使用
1.启动时，自动会从mysql中加载配置

2.触发重加载，还是调用 prometheus的 reload 接口
- 规则是否修改的判断逻辑：同 prometheus 原有逻辑，即规则的 key = rule_fn;rule_name
- 规则的修改，会引发 rule 重新计算，即 firing 状态的 rule 会恢复 pending 并重新开始计算

# Note
**prometheus 源码对配置对象的处理不是很一致**
- 在 config/ 下有 config 对象，但在 main.go 中另定义了 cfg 对象，又未完整引用 config 对象。
- 没有一个全局配置对象，代码中要想应用配置项，基本需要在其模块中定义局部配置变量，然后在 main.go 中对其赋值

# 相关
## 辅助
有两个相关项目：
* [mysql4Prom](https://github.com/huangwei2013/mysql4prom)：用于解析 prometheus 规则文件，并导入DB
* [mysql4PromUI](https://github.com/huangwei2013/mysql4promUI)：用于 DB 对应的简单管理界面

# 注意事项
* 新版Prometheus需go1.16（主要是 k8s.io/client_go 使用了 io/fs 等新特性）
* 若要调试，dvl需自行编译
  * dvl 源码中（MaxSupportedVersionOfGoMinor）限定仅支持到 go1.14，手动修改到 go1.16 编译也可正常使用

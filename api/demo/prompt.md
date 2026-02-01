## api定义
每个项目可能有多种资源, 每种资源有自己的service。

要求：
1. request 相关的message 放在service 文件中，其他的资源message 要用独立的文件
2. 不要api定义上默认加上类似v1这样的版本控制
3. 取名字要通用一些， 便于业务切换，比如 VoiceMessage 直接叫 Message , VoiceSession 直接叫Session。
4. 路径命名规则：/api/{service_name_short}/资源名复数 (例如 foo 缩写为 fo, manage 缩写为 mg)。
5. 定义 extra.go 文件，是为了补充生成的go struct 的方法，
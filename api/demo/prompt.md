## api定义
每个项目可能有多种资源, 每种资源有自己的service。

要求：
1. request 相关的message 放在service 文件中，其他的资源message 要用独立的文件
2. 不要api定义上默认加上类似v1这样的版本控制
3. 取名字要通用一些， 便于业务切换，比如 VoiceMessage 直接叫 Message , VoiceSession 直接叫Session。
4. 路径命名规则：/api/{service_name_short}/资源名复数 (例如 foo 缩写为 fo, manage 缩写为 mg)。
5. 定义 extra.go 文件，是为了补充生成的go struct 的方法，


## 实现
- 项目结构要参考 app/demo的目录结构。要先完整复制一份app/demo,然后在这基础上改, 去掉无用的代码。
1. mongodb的泛型要定义成指针。示例: *mgz.Core[*demo.Foo]
2. 数据层 核心要用 `*mgz.Core` 。涉及到的资源message 要实现对应的方法 示例：
```
    type FooCollection struct {    
        *mgz.Core[*demo.Foo]
    }

    func NewFooCollection(db *mongo.Database) *FooCollection {
        return &FooCollection{
            Core: mgz.NewCore[*demo.Foo](db, "foos"),
        }
    }

需要实现的方法：
    func (m *Foo) GetID() string {
        if m = nil {
            return ""
        }
        return m.XId
    }

    func (m *Foo) SetID(id string) {
        m.XId = id
    }
    
```
3. 在service 文件下编写对应的实现。 每个方法要独立一个文件，文件名格式是 {service名}_{方法名}。示例： voiceagent_CreatePersona.go

4.根据  api 中的解释和定义， 优化所有的接口的实现 ，加上必要的逻辑 
5.直接访问数据库 不要这个biz层
# 前端

1. 如果涉及到 LinearGradient  记得将样式写在 style中 不然可能会不生效



请按照我的工作流帮我编写代码：
## 数据格式定义。

api定义要分 service和 普通资源 message
 1. request 相关的message  放在service 文件中

 3. 不要api定义上默认加上类似v1这样的版本控制
 4. 普通资源 message 需要加上userId

 取名字要通用一些， 便于业务切换，比如 VoiceMessage 直接叫 Message , VoiceSession 直接叫Session


## 实现
- 项目结构要参考 app/usercenter的目录结构。要先完整复制一份app/usercenter,然后在这基础上改
1. mongodb的泛型要定义成指针。示例: *mgz.CoreCollectionV3[*voiceagent.Persona]
2. 数据层 核心要用 `*mgz.Core` 。涉及到的资源message 要实现对应的方法 示例：
```
    type UserCollection struct {
        *mgz.Core[*ucpb.User]
    }

    func NewUserCollection(db *mongo.Database) *UserCollection {
        return &UserCollection{
            Core: mgz.NewCore[*ucpb.User](db, "users"),
        }
    }

需要实现的方法：
    func (m *User) GetID() string {
        return m.XId
    }

    func (m *User) SetID(id string) {
        m.XId = id
    }
    
```
3. 在service 文件下编写对应的实现。 每个方法要独立一个文件，文件名格式是 {service名}_{方法名}。示例： voiceagent_CreatePersona.go

4.根据  api 中的解释和定义， 优化所有的接口的实现 ，加上必要的逻辑 
5.直接访问数据库 不要这个biz层



这是开发者文档  请仔细分析 。给出api定义中所有接口的 必要性和详细的逻辑， 输出输出分别是什么，详细解释输入输出的每个字段的含义和参考， 注意事项是什么，参考的文档的那个点 。以注释的形式 写在api的定义处。直接以注释的形式写在代码中

这是开发者文档  请仔细分析 。api定义中所有的message及其所有字段端的必要性 和含义，用法，已经参考文档的哪个点。写在api的定义处。直接以注释的形式写在代码中


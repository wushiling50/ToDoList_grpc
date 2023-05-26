# ToDoList_grpc
基于gin+gorm+grpc+etcd+mysql的备忘录demo

#在线api文档地址：https://apifox.com/apidoc/shared-bbd81ba6-7d76-4e96-afd5-0cf6939df8ed

#项目主要依赖
> -gin
> -gorm 
> -grpc
> -etcd
> -sentinel
> -logrus
> -viper
> -jwt-go
> -protobuf

#项目结构
ToDoList-grpc/
├── api-gateway //网关部分
|   |---cmd //启动函数
|   |---discovery //etcd服务发现
|   |---inner //业务逻辑
|   |---middleware //中间件
|   |---routes //路由组
├── conf //配置文件
├── pkg //工具包
├── service //存放pb文件及其生成的服务文件
├── pkg //工具包
├── task //任务部分
|   |---cmd //启动函数
|   |---inner //业务逻辑
|   |---login //etcd服务注册
├── user //用户部分
|   |---cmd //启动函数
|   |---inner //业务逻辑
|   |---login //etcd服务注册









# 项目说明
    微服务框架
## 目录结构
    ├── core  核心文件                  
            ├── breaker  熔断
            ├── conf  配置文件加载           
            ├── discov  服务发现
            ├── logsj  日志
            ├── rollingwindow 滑动时间窗口
            └── trace  链路追踪
    
    ├── httpsj  http服务端
            ├── handler  http中间件
            ├── init.go
            └── server.go
    
    ├── srpc rpc服务端+rpc客户端
            ├── clientinterceptor  grpc客户端拦截器
            ├── serverinterceptors grpc服务端拦截器
            ├── codes
            ├── client.go
            ├── init.go
            ├── server.go
            └── types.go

    └── utils  工具包
            ├── ginfunc gin相关工具包
            ├── logfunc 日志相关工具包
            ├── mathfunc 数字相关工具包
            ├── stringfunc 字符串相关工具包
            ├── sysfunc 系统相关工具包
            └── timefunc 时间相关工具包

# mirai-go
基于mirai http-api的 golang框架  
处于刚刚起步的开发阶段

- [x] 群组消息解析
- [x] 用户消息解析
- [x] 发送好友消息
- [x] 发送群组信息
- [x] 消息回复回调框架
- [x] 测试补全benchmark 现阶段event center和朴素的sync.RWMap的时间 

对EventCenter以及平凡Map实现的事件回调中心的BenchMark

| type                             | ops | ns/op          |
| -------------------------------- | --- | -------------- |
| BenchmarkEventCenter-8           | 99  | 11680378 ns/op |
| BenchmarkPlain-8                 | 70  | 15864608 ns/op |
| BenchmarkEventCenterSequential-8 | 172 | 6921466 ns/op  |
| BenchmarkPlainSequential-8       | 72  | 16749177 ns/op |
| BenchmarkEventCenterParallel-8   | 94  | 12190873 ns/op |
| BenchmarkPlainParallel-8         | 33  | 30933554 ns/op |
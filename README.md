# proto2json 
 * 解析protobuf String()生成的字符串为json格式
# 使用方法
 * go run main.go -h
```
  -d    active debug log
  -h    print this
  -t string
        set template value
```
 * go run main.go -t 'b:1 b:2 b:3 d:<d:5 d:6 > c:<c:4 >'
```
   {"b":[1,2,3],"c":{"c":4},"d":{"d":[5,6]}}
```
# 已知问题
 * 当数组只有一个元素时会被解析为对象

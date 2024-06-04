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

 * go run main.go -t 'b:1 b:2 b:3 d:<d:5 d:6 > c:<c:4 >' -a='c' 强制把c解析为数组
```
   {"b":[1,2,3],"c":[{"c":[4]}],"d":{"d":[5,6]}}
```

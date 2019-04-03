# commandstomap
convert commands to map

# 工具名

commandstomap

#功能

将命令行参数转换成数组和健值对，可以支持换行、转义等等

# DEMO

```go
func Main() {
	cmd := `-aaaaa 4 \
	-bbbbb 8 \
	-env ""   \
	-config "hellworl" `
	a := StringToArray(cmd)
	for i, d := range a {
		fmt.Println(i, ":", d)
	}
	m, _ := StringToMap(cmd)
	fmt.Println(m["env"])

}
```


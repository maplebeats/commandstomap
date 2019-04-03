package commandstomap

import (
	"fmt"
	"strings"
)

//StringToArray 将shell命令转换成Array
//-a 123 -b hello ---> ["-a","123","-b","hello"]
func StringToArray(ext string) []string {
	extRaw := []byte(strings.TrimSpace(ext))
	var array []string
	var tmp []byte
	//sOff 单引号开始,doff 双引号开始, escape \进行转换
	var sOff, dOff, escape bool
	for offset := 0; offset < len(extRaw); offset++ {
		switch extRaw[offset] {
		case ' ':
			if sOff || dOff {
				tmp = append(tmp, extRaw[offset])
				continue
			}
			if tmp != nil {
				array = append(array, string(tmp))
				tmp = nil
			}
		case '"':
			//单双引号互相引用
			if dOff || escape {
				tmp = append(tmp, extRaw[offset])
				escape = false
				continue
			}
			//开关
			sOff = !sOff
			if !sOff && tmp == nil {
				//如果关闭之后数据还为空，增加一个空字段进去
				array = append(array, string(tmp))
			}
		case '\'':
			if sOff || escape {
				tmp = append(tmp, extRaw[offset])
				escape = false
				continue
			}
			dOff = !dOff
			if !dOff && tmp == nil {
				array = append(array, string(tmp))
			}
			/*
				case '-':
					//如果不是开头-,不会忽略
					if (sOff || dOff) || (tmp != nil) {
						tmp = append(tmp, extRaw[offset])
						continue
					}
			*/
		case '\t':
			if sOff || dOff {
				tmp = append(tmp, extRaw[offset])
				continue
			}
		case '\\':
			//开始转义
			if escape {
				tmp = append(tmp, extRaw[offset])
				escape = !escape
				continue
			}
			escape = !escape //true
		case '\r':
			fallthrough
		case '\n':
			if sOff || dOff {
				tmp = append(tmp, extRaw[offset])
				continue
			}
			if escape {
				//换行结束，重置
				escape = !escape
				continue
			}
			if tmp != nil {
				array = append(array, string(tmp))
				tmp = nil
			}
		default:
			tmp = append(tmp, extRaw[offset])
			if escape {
				//退出转义
				escape = !escape
			}
		}
		if offset == (len(extRaw) - 1) {
			array = append(array, string(tmp))
		}
	}
	return array
}

//StringToMap 将shell命令转换成Map
//-a 123 -b hello ---> {"a":"123","b":"hello"}
func StringToMap(ext string) (map[string]string, error) {
	array := StringToArray(ext)
	if len(array)%2 != 0 {
		return nil, fmt.Errorf("array length error:%d", len(array))
	}
	m := make(map[string]string, len(array)/2)
	for i := range array {
		if i%2 == 0 {
			k := strings.TrimLeft(array[i], "-")
			m[k] = array[i+1]
		}
	}
	return m, nil
}

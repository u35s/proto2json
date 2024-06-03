package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"strconv"
	"strings"
)

const (
	EleTypeObjectStart = iota
	EleTypeObjectEnd
	// EleTypeObjectStart
	EleTypeObjectKey
	EleTypeObjectValueStr
	EleTypeObjectValue
)

type LogWriter struct {
	Enable bool
}

func (w *LogWriter) Write(p []byte) (n int, err error) {
	if w.Enable {
		fmt.Printf("%s", p)
		return len(p), nil
	}
	return 0, nil
}

type Element struct {
	value     interface{}
	lastValue interface{}
	lastKey   string
	arr       []interface{}
	dict      map[string]interface{}
}

func NewEle() *Element {
	return &Element{
		arr:  make([]interface{}, 0),
		dict: make(map[string]interface{}),
	}
}

func (e *Element) Set(child interface{}) {
	e.value = child
}

func (e *Element) Add(k string, child *Element) {
	if e.lastKey != k {
		if len(e.arr) > 0 {
		        e.arr = append(e.arr, e.lastValue)
			e.dict[e.lastKey] = e.arr
			e.arr = make([]interface{}, 0)
		} else if e.lastKey != "" {
			e.dict[e.lastKey] = e.lastValue
		}
	}
	if e.lastKey == k {
		e.arr = append(e.arr, e.lastValue)
	}
	e.lastKey = k
	e.lastValue = child
}

func (e *Element) MarshalJSON() ([]byte, error) {
	if e.value != nil {
		return json.Marshal(e.value)
	}
	if len(e.arr) > 0 {
		e.arr = append(e.arr, e.lastValue)
		e.dict[e.lastKey] = e.arr
	} else if e.lastKey != "" {
		e.dict[e.lastKey] = e.lastValue
	}
	return json.Marshal(e.dict)
}

func parseValue(s string) interface{} {
	// 去除字符串前后的空白字符
	s1 := strings.TrimSpace(s)

	// 尝试转换为整数
	if a, err := strconv.Atoi(s1); err == nil {
		return a
	}

	// 尝试转换为浮点数
	if b, err := strconv.ParseFloat(s1, 64); err == nil {
		return b
	}

	// 如果都失败，则不是数字
	return s1
}

func parseEle(tpl string, start int, parent *Element) int {
	for {
		key, tp, pos := parseType(tpl, start)
		if tp == EleTypeObjectStart {
			v := NewEle()
			parent.Add(key, v)
			start = parseEle(tpl, pos+2, v)
			start++
		}

		if tp == EleTypeObjectEnd {
			return pos + 1
		}
		if tp == EleTypeObjectKey {
			val1, tp1, pos1 := parseType(tpl, pos+1)
			v := NewEle()
			if tp1 == EleTypeObjectValue {
				v.Set(parseValue(val1))
				parent.Add(key, v)
				start = pos1 + 1
			} else if tp1 == EleTypeObjectValueStr {
				v.Set(val1)
				parent.Add(key, v)
				start = pos1 + 1
			} else {
				log.Printf("obj value parse error,pos:%v,tp:%v,val:%v", pos1, tp1, val1)
				return -1
			}
		}
		if start >= len(tpl) {
			break
		}
	}
	return start
}

func parseType(tpl string, start int) (string, int, int) {
	isStr := tpl[start] == '"'
	for i := start; i < len(tpl); i++ {
		if isStr {
			isStr2 := tpl[i] == '"'
			if isStr && i != start && isStr2 {
				tp := EleTypeObjectValueStr
				log.Printf("parse string,start:%v,end:%v,str:%v", start, i, tpl[start+1:i-1])
				return tpl[start+1 : i-1], tp, i + 1
			}
			continue
		}
		if tpl[i] == '>' {
			tp := EleTypeObjectEnd
			return tpl[start:i], tp, i
		}
		if tpl[i] == ':' {
			tp := EleTypeObjectStart
			if tpl[i+1] != '<' {
				tp = EleTypeObjectKey
			}
			return tpl[start:i], tp, i
		}

		if tpl[i] == ' ' {
			tp := EleTypeObjectValue
			log.Printf("parse value,start:%v,end:%v,str:%v", start, i, tpl[start:i])
			return tpl[start:i], tp, i
		}
		if len(tpl)-1 == i {
			tp := EleTypeObjectValue
			log.Printf("parse value2,start:%v,end:%v,str:%v", start, i, tpl[start:i+1])
			return tpl[start : i+1], tp, i
		}

	}
	return "", -1, -1
}

func main() {
	var res = NewEle()
	var template = "result:<top_entries:<id:EI_PINK col:2 index:1 start:1 len:1 > top_entries:<id:EI_SCATTER col:3 index:1 start:1 len:1 > top_entries:<id:EI_WILD col:4 index:1 start:1 len:1 > top_entries:<id:EI_RED col:5 index:1 start:1 len:1 > entries:<id:EI_Q col:1 index:2 start:2 len:1 > entries:<id:EI_SCATTER col:1 index:3 start:3 len:1 > entries:<id:EI_J col:1 index:4 start:4 len:1 > entries:<id:EI_J col:1 index:5 start:5 len:1 > entries:<id:EI_10 col:1 index:6 start:6 len:1 > entries:<id:EI_10 col:2 index:2 start:2 len:3 border:Silver > entries:<id:EI_PINK col:2 index:3 start:5 len:1 > entries:<id:EI_Q col:2 index:4 start:6 len:1 > entries:<id:EI_PURPLE col:3 index:2 start:2 len:4 border:Silver > entries:<id:EI_BLUE col:3 index:3 start:6 len:1 > entries:<id:EI_J col:4 index:2 start:2 len:3 > entries:<id:EI_BLUE col:4 index:3 start:5 len:2 > entries:<id:EI_SCATTER col:5 index:2 start:2 len:2 > entries:<id:EI_10 col:5 index:3 start:4 len:1 > entries:<id:EI_Q col:5 index:4 start:5 len:1 > entries:<id:EI_10 col:5 index:5 start:6 len:1 > entries:<id:EI_A col:6 index:2 start:2 len:1 > entries:<id:EI_J col:6 index:3 start:3 len:1 > entries:<id:EI_10 col:6 index:4 start:4 len:1 > entries:<id:EI_K col:6 index:5 start:5 len:1 > entries:<id:EI_GREEN col:6 index:6 start:6 len:1 > cat_rate:1 >"
	var logWriter = &LogWriter{}

	templateP := flag.String("t", template, "-t:set template value")
	debugLogP := flag.Bool("d", false, "-d:active debug log")
	flag.Parse()

	template = *templateP
	logWriter.Enable = *debugLogP

	log.SetOutput(logWriter)
	parseEle(template, 0, res)
	bts, err := json.Marshal(res)
	if err == nil {
		fmt.Println(string(bts))
	}
}

package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
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
		e.flushLastValue()
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
	e.flushLastValue()
	return json.Marshal(e.dict)
}

func (e *Element) flushLastValue() {
	if e.lastKey == "" {
		return
	}
	if len(e.arr) > 0 {
		e.arr = append(e.arr, e.lastValue)
		e.dict[e.lastKey] = e.arr
		e.arr = make([]interface{}, 0)
	} else if e.lastKey != "" {
		if _, ok := arrayKeyMap[e.lastKey]; ok {
			e.arr = append(e.arr, e.lastValue)
			e.dict[e.lastKey] = e.arr
			e.arr = make([]interface{}, 0)
		} else {
			e.dict[e.lastKey] = e.lastValue
		}
	}
	e.lastKey = ""
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
			log.Printf("parse key,start:%v,end:%v,str:%v", start, i, tpl[start:i])
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

type Args struct {
	Str         string
	StrFile     string
	Debug       bool
	WriteResult bool
	Indent      bool
	Help        bool
	ArrayKey    string
}

func (a *Args) Parse() {
	flag.StringVar(&a.Str, "s", "", "to parse value")
	flag.StringVar(&a.StrFile, "f", "proto.txt", "to parse file")
	flag.StringVar(&a.ArrayKey, "a", "", "array key,split by ,")
	flag.BoolVar(&a.WriteResult, "w", true, "write result to file")
	flag.BoolVar(&a.Indent, "t", true, "json text format")
	flag.BoolVar(&a.Debug, "d", false, "active debug log")
	flag.BoolVar(&a.Help, "h", false, "print this")
	flag.Parse()
	if a.Help {
		flag.Usage()
		os.Exit(0)
	}
	if a.Str == "" {
		bts, err := os.ReadFile(a.StrFile)
		if err == nil && len(bts) > 0 {
			a.Str = string(bts)
		} else if err != nil {
			if a.StrFile != "proto.txt" {
				fmt.Printf("file %v not exists", a.StrFile)
				os.Exit(0)
			}
			flag.Usage()
			os.Exit(0)
		}
	}

	if a.ArrayKey != "" {
		strs := strings.Split(a.ArrayKey, ",")
		for i := 0; i < len(strs); i++ {
			arrayKeyMap[strs[i]] = true
		}
		fmt.Printf("array key %v\n", arrayKeyMap)
	}
}

var arrayKeyMap = make(map[string]bool)

func main() {
	var res = NewEle()
	var logWriter = &LogWriter{}

	var args = &Args{}
	args.Parse()
	logWriter.Enable = args.Debug

	log.SetOutput(logWriter)
	parseEle(args.Str, 0, res)
	var (
		bts []byte
		err error
	)
	if args.Indent {
		bts, err = json.MarshalIndent(res, "", "\t")
	} else {
		bts, err = json.Marshal(res)
	}
	if err == nil {
		fmt.Println(string(bts))
		if args.WriteResult {
			os.WriteFile("json_"+args.StrFile, bts, 0644)
		}
	}
}

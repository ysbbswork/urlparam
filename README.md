# urlparam  provides url.Value to Struct and Struct to url.Value, and support protobuf struct.
# urlparam 快速将url参数和struct直接相互转换

urlparam 提供url参数自动解析到结构体，实现url.Value to Struct和 Struct to url.Value
支持结构体标签，利用结构体的反射标签可以自定义参数名称的映射关系，
默认利用json标签，可以使此库可直接用pb.go生成的结构体，便于不同服务之间使用pb管理url params 的协议
支持修改标签类型，解析标签名可修改导出变量URLParamTag来改变
没有定义标签则默认使用结构体的字段名称作为url参数名称，并且 - 标签作为不解析标记，和json标签语法兼容

# 使用示例

以UserReportParams作为示例
```
	// UserReportParams 作为示例
	type UserReportParams struct {
	BusiType   int32  `json:"busi_type,omitempty"` // 同一个标签里面有option参数，例如PB生成的结构
	Uid        string `json:"uid"`                 // 一个字段一个标签 string类型
	AdposType  int32  `json:"adpos_type"`
	ActionType int32  // 无标签
	FeedsIndex uint32 `json:"feeds_index"` // uint32 类型
	AdposId    string `json:"-"`           // 屏蔽标签
	}
```

## 结构体转为url参数
```
	// 结构体转为URL参数
	userData := &UserReportParams{
		BusiType:   1,
		Uid:        "2222",
		AdposType:  555,
		ActionType: -666,
		FeedsIndex: 7777,
		AdposId:    "aaaaa",
	}
	values, errs := urlparam.Encode(userData)
	if errs != nil {
		fmt.Printf("errors %v", errs)
		return
	}
	fmt.Printf("%v\n", values.Encode())
```

输出结果
```
ActionType=-666&adpos_type=555&busi_type=1&feeds_index=7777&uid=2222
```
可以看到结构体值被序列化成了url参数
1.AdposType没有定义标签，则使用字段名作为参数名
2.AdposId 写了屏蔽标签，即使有值也会被编码到url参数里

## url参数转为结构体

```
	// url参数解析到结构体
	urlStr := "AdposId=aaaa&ActionType=-666&adpos_type=555&busi_type=1&feeds_index=7777&uid=2222"
	values, _ := url.ParseQuery(urlStr)
	res := UserReportParams{}
	err := urlparam.Decode(values, &res)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Printf("%+v\n", res)
```

输出结果：
```
{BusiType:1 Uid:2222 AdposType:555 ActionType:-666 FeedsIndex:7777 AdposId:}
```
可以看到按结构体标签对应的url 参数被解析到结构体中，AdposId 写了屏蔽标签，即使有值也不会被解析
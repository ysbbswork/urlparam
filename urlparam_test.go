package urlparam

import (
	"net/url"
	"reflect"
	"testing"
)

// UserReportParams 作为示例
type UserReportParams struct {
	BusiType   int32  `json:"busi_type,omitempty"` // 同一个标签里面有option参数，例如PB生成的结构
	UID        string `json:"uid"`                 // 一个字段一个标签 string类型
	AdposType  int32  `json:"adpos_type"`
	ActionType int32  // 无标签
	FeedsIndex uint32 `json:"feeds_index"` // uint32 类型
	AdposID    string `json:"-"`           // 屏蔽标签
}

func TestEncode(t *testing.T) {
	type args struct {
		s interface{}
	}
	want := url.Values{}
	want.Set("busi_type", "1")
	want.Set("uid", "2222")
	want.Set("adpos_type", "555")
	want.Set("ActionType", "-666")
	want.Set("feeds_index", "7777")
	tests := []struct {
		name    string
		args    args
		want    url.Values
		wantErr bool
	}{
		{
			name: "Encode",
			args: args{&UserReportParams{
				BusiType:   1,
				UID:        "2222",
				AdposType:  555,
				ActionType: -666,
				FeedsIndex: 7777,
				AdposID:    "aaaaa",
			}},
			want: want,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Encode(tt.args.s)
			if (err != nil) != tt.wantErr {
				t.Errorf("Encode() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Encode() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDecode(t *testing.T) {
	urlStr := "ActionType=-666&adpos_type=555&busi_type=1&feeds_index=7777&uid=2222"
	values, _ := url.ParseQuery(urlStr)
	got := &UserReportParams{}
	want := &UserReportParams{
		BusiType:   1,
		UID:        "2222",
		AdposType:  555,
		ActionType: -666,
		FeedsIndex: 7777,
	}
	type args struct {
		urlValues url.Values
		s         interface{}
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Decode",
			args: args{
				urlValues: values,
				s:         got,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := Decode(tt.args.urlValues, tt.args.s); (err != nil) != tt.wantErr {
				t.Errorf("Decode() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !reflect.DeepEqual(got, want) {
				t.Errorf("Encode() = %v, want %v", got, want)
			}

		})
	}
}
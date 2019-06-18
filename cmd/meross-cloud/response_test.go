package main

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strconv"
	"testing"
	"time"
)

func MustParse(s string) time.Time {
	t, err := time.ParseInLocation(time.RFC3339, s, serverTimezone)
	if err != nil {
		panic(err)
	}
	return t
}

func Test_ResponseUnmarshal(t *testing.T) {
	tt := []struct {
		in     string
		expOut Response
		expErr string
	}{
		{
			`{"apiStatus":0,"sysStatus":0,"data":{"token":"abc123df","key":"0000","userid":"1111","email":"me@email.address.com"},"info":"Success","timeStamp":1558237198}`,
			Response{
				APIStatus: 0,
				SysStatus: 0,
				Data:      json.RawMessage(`{"token":"abc123df","key":"0000","userid":"1111","email":"me@email.address.com"}`),
				Info:      "Success",
				Timestamp: ResponseTime(time.Unix(1558237198, 0)),
			},
			"<nil>",
		},
		{
			`{"apiStatus":1004,"sysStatus":0,"data":null,"info":"A","timeStamp":"2019-05-19 11:18:52"}`,
			Response{
				APIStatus: 1004,
				SysStatus: 0,
				Data:      json.RawMessage(`null`),
				Info:      "A",
				Timestamp: ResponseTime(MustParse("2019-05-19T11:18:52+08:00")),
			},
			"<nil>",
		},
	}

	for i, tc := range tt {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			var out Response
			err := json.Unmarshal([]byte(tc.in), &out)
			if fmt.Sprint(err) != tc.expErr {
				t.Errorf("expected %q got %q", tc.expErr, err)
			}
			if !reflect.DeepEqual(out, tc.expOut) {
				t.Errorf("diff \n expected %#v \n got      %#v", tc.expOut, out)
			}
		})
	}
}

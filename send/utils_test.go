/*
 * Copyright 2018 Xiaomi, Inc.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */
package send

import (
	"os"
	"reflect"
	"testing"

	"github.com/open-falcon/gpu-mon/cfg"
)

func Test_decimal(t *testing.T) {
	type args struct {
		value float64
	}
	tests := []struct {
		name string
		args args
		want float64
	}{
		{"test1", args{3.1786}, 3.18},
		{"test2", args{3.00}, 3.00},
		{"test3", args{34}, 34.00},
		{"test4", args{3.2}, 3.20},
		{"test5", args{3.299}, 3.30},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := decimal(tt.args.value); got != tt.want {
				t.Errorf("decimal() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_checkMetricData(t *testing.T) {
	type args struct {
		data interface{}
	}
	var (
		value1 uint   = 10
		value2 uint64 = 32
		value3        = 3.14
		value4        = -3
		value5        = "0"
	)

	tests := []struct {
		name      string
		args      args
		wantRes   interface{}
		wantIsNaN bool
	}{

		{"test1", args{&value1}, uint(10), false},
		{"test2", args{&value2}, uint64(32), false},
		{"test3", args{&value3}, 3.14, false},
		{"test4", args{&value4}, -3, false},
		{"test5", args{&value5}, -1, true},
		{"test6", args{nil}, -1, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotRes, gotIsNaN := checkMetricData(tt.args.data)
			if !reflect.DeepEqual(gotRes, tt.wantRes) {
				t.Errorf("checkMetricData() gotRes = %v, want %v", gotRes, tt.wantRes)
			}
			if gotIsNaN != tt.wantIsNaN {
				t.Errorf("checkMetricData() gotIsNaN = %v, want %v", gotIsNaN, tt.wantIsNaN)
			}
		})
	}
}

func Test_getFalconAgent(t *testing.T) {

	tests := []struct {
		name             string
		wantFalconAgent  string
		inputFalconAgent string
	}{
		{"test1", "http://127.0.0.1:1988/v1/push", ""},
		{"test2", "testAgent", "testAgent"},
	}
	for _, tt := range tests {
		conf := cfg.Config()
		conf.Falcon.Agent = tt.inputFalconAgent
		t.Run(tt.name, func(t *testing.T) {
			if gotFalconAgent := getFalconAgent(); gotFalconAgent != tt.wantFalconAgent {
				t.Errorf("getFalconAgent() = %v, want %v", gotFalconAgent, tt.wantFalconAgent)
			}
		})
	}
}

func Test_getEndPoint(t *testing.T) {
	hostname, _ := os.Hostname()
	tests := []struct {
		name          string
		wantEndPoint  string
		inputEndPoint string
	}{

		{"test1", hostname, ""},
		{"test2", "testhost", "testhost"},
	}
	for _, tt := range tests {
		conf := cfg.Config()
		conf.Metric.EndPoint = tt.inputEndPoint
		t.Run(tt.name, func(t *testing.T) {
			if gotEndPoint := getEndPoint(); gotEndPoint != tt.wantEndPoint {
				t.Errorf("getEndPoint() = %v, want %v", gotEndPoint, tt.wantEndPoint)
			}
		})
	}
}

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
package common

import (
	"reflect"
	"testing"
)

var wantConfig1 = Conf{
	Falcon:       falconConf{""},
	Metric:       metricConf{nil, ""},
	Log:          logConf{"", ""},
	MetricFilter: map[string]struct{}{},
	IsCrontab:    false,
}

var wantConfig2 = Conf{
	Falcon: falconConf{"http://127.0.0.1:1988/v1/push"},
	Metric: metricConf{
		[]string{
			"FanSpeed",
			"Tx",
			"Rx"},
		"testEndPoint"},
	Log: logConf{"Warn", "logs"},
	MetricFilter: map[string]struct{}{
		"FanSpeed": struct{}{},
		"Tx":       struct{}{},
		"Rx":       struct{}{},
	},
	IsCrontab: false,
}

func TestInitConfig(t *testing.T) {
	type args struct {
		configPath string
		isCrontab  bool
	}
	tests := []struct {
		name          string
		args          args
		wantErr       bool //true means err!=nil
		wantIsCrontab bool
	}{
		{"t_confFileNotExist", args{"./testData/notExit.json", false}, true, false},
		{"t_IsCrontab", args{"./testData/test.cfg.json", false}, false, false},
		{"t_IsCrontab", args{"./testData/test.cfg.json", true}, false, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := InitCommon(tt.args.configPath, tt.args.isCrontab); (err != nil) != tt.wantErr {
				t.Errorf("InitConfig() error = %v, wantErr %v", err, tt.wantErr)
			}
			if globalConf.IsCrontab != tt.wantIsCrontab {
				t.Errorf("InitConfig() \n got isCrontab = %v, want isCrontab %v",
					globalConf.IsCrontab, tt.wantIsCrontab)
			}

		})
	}
}

func Test_readConfigFile(t *testing.T) {
	type args struct {
		configPath string
	}
	tests := []struct {
		name     string
		args     args
		wantErr  bool
		wantConf Conf
	}{
		{"test1", args{"./testData/notExit.json"}, true, wantConfig1},
		{"test2", args{"./testData/test.cfg.json"}, false, wantConfig2},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := readConfigFile(tt.args.configPath)
			if (err != nil) != tt.wantErr {
				t.Errorf("readConfigFile() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !tt.wantErr {
				initIgnoreMetrics()
				globalConf.IsCrontab = false //默认为false
				if !(reflect.DeepEqual(globalConf.MetricFilter, tt.wantConf.MetricFilter)) {
					t.Errorf("initIgnoreMetrics() FAILED:\nget ignoreMetric = %v \nwant ignoreMetric = %v",
						globalConf.MetricFilter, tt.wantConf.MetricFilter)
				}
				if !(reflect.DeepEqual(globalConf, tt.wantConf)) {
					t.Errorf("readConfigFile() \ngot globalconf = %v\nwant config = %v", globalConf, tt.wantConf)
				}
			}
		})
	}
}

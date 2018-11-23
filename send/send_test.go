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
	"reflect"
	"testing"
	"time"

	"github.com/open-falcon/gpu-mon/common"
	"github.com/open-falcon/gpu-mon/fetch"

	"github.com/sirupsen/logrus"
)

func init() {
	common.Logger.SetLevel(logrus.PanicLevel)
}
func Test_buildMetaData(t *testing.T) {
	conf := common.Config()
	conf.Metric.EndPoint = "testEndPoint"
	metricdata := 3
	wantMetaData := MetaData{
		"testName",
		"testEndPoint",
		time.Now().Unix(),
		60,
		metricdata,
		"GAUGE",
		"GpuId=0",
	}

	type args struct {
		metricName  string
		metricValue interface{}
		tags        string
		step        int
	}
	tests := []struct {
		name         string
		args         args
		wantMetaData MetaData
	}{
		{"test1", args{"testName", &metricdata, "GpuId=0", 60}, wantMetaData},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotMetaData := buildMetaData(tt.args.metricName, tt.args.metricValue, tt.args.tags, tt.args.step); !reflect.DeepEqual(gotMetaData, tt.wantMetaData) {
				t.Errorf("buildMetaData() = %v, want %v", gotMetaData, tt.wantMetaData)
			}
		})
	}
}

func TestBuildMetaDatas(t *testing.T) {
	type args struct {
		rawDataList []fetch.RawData
	}

	uintValue := uint(49)
	uint64Value := uint64(12345)
	float64Value := float64(3.1415)
	intValue := int(0)
	inputValues := fetch.MetricValues{
		GPUUtils:      &uintValue,
		Rx:            &uint64Value,
		PowerUsed:     &float64Value,
		DcgmSupported: &intValue,
		MemUtils:      nil,
	}

	rawData := fetch.RawData{
		GpuID:  0,
		Values: inputValues,
	}
	inputData := []fetch.RawData{
		rawData,
	}
	timestamp := time.Now().Unix()

	wantMetaDataMap := map[string]MetaData{
		"GPUUtils":      MetaData{"GPUUtils", "testEndPoint", timestamp, 60, uintValue, "GAUGE", "GpuId=0"},
		"Rx":            MetaData{"Rx", "testEndPoint", timestamp, 60, uint64Value, "GAUGE", "GpuId=0"},
		"PowerUsed":     MetaData{"PowerUsed", "testEndPoint", timestamp, 60, float64Value, "GAUGE", "GpuId=0"},
		"DcgmSupported": MetaData{"DcgmSupported", "testEndPoint", timestamp, 60, intValue, "GAUGE", "GpuId=0"},
		"MemUtils":      MetaData{"MemUtils", "testEndPoint", timestamp, 60, -1, "GAUGE", "GpuId=0"},
	}

	tests := []struct {
		name          string
		args          args
		wantMetaDatas map[string]MetaData
	}{
		{"test", args{inputData}, wantMetaDataMap},
	}

	for _, tt := range tests {
		config := common.Config()
		config.Metric.EndPoint = "testEndPoint"

		t.Run(tt.name, func(t *testing.T) {
			gotMetaDataList := BuildMetaDatas(tt.args.rawDataList)
			for _, gotmetaData := range gotMetaDataList {
				if wantMetaData, ok := tt.wantMetaDatas[gotmetaData.Metric]; ok {
					if gotmetaData.TAGS != "type=ave" && gotmetaData.TAGS != "type=sum" && !reflect.DeepEqual(gotmetaData, wantMetaData) {
						t.Errorf("BuildMetaDatas()\nget %v\nwant %v", gotmetaData, wantMetaData)
					}
				}
			}
		})
	}
}

func TestData(t *testing.T) {
	inputMetaDataList := []MetaData{
		MetaData{
			Metric:      "test_gpu_1",
			Endpoint:    "192.168.0.1:1988",
			Timestamp:   time.Now().Unix(),
			Step:        60,
			Value:       "ok",
			CounterType: "GAUGE",
			TAGS:        "type=test",
		},
		MetaData{
			Metric:      "test_endpoint2",
			Endpoint:    "192.168.0.1:1988",
			Timestamp:   time.Now().Unix(),
			Step:        60,
			Value:       1,
			CounterType: "COUNTER",
			TAGS:        "type=test",
		},
	}

	type args struct {
		metaDataList []MetaData
	}

	tests := []struct {
		name    string
		args    args
		agent   string
		wantErr bool
	}{
		{"test1", args{inputMetaDataList}, "nil", true},
		{"test2", args{inputMetaDataList}, "", false},
		{"test3", args{inputMetaDataList}, "http://127.0.0.1:1988/v1/push", false},
	}
	for _, tt := range tests {
		conf := common.Config()
		conf.Falcon.Agent = tt.agent
		t.Run(tt.name, func(t *testing.T) {
			if err := pushAgent(tt.args.metaDataList); (err != nil) != tt.wantErr {
				t.Errorf("pushAgent() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

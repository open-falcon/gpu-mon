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
	"fmt"
	"os"
	"strconv"

	"github.com/open-falcon/gpu-mon/common"
)

// 保留float数的两位小数
func decimal(value float64) float64 {
	value, _ = strconv.ParseFloat(fmt.Sprintf("%.2f", value), 64)
	return value
}

// result = i * k / j
func dvInt(i, j, k int) float64 {
	if j == 0 {
		return -1
	}
	res := float64(i*k) / float64(j)
	return decimal(res)
}

// GetHostName 获取主机名
func getHostName() (host string, err error) {
	host, err = os.Hostname()
	if err != nil {
		return
	}
	return
}

// 检查MetricData，使用isNaN标识异常数据，避免使用interface{} 判断数据的类型断言
func checkMetricData(data interface{}) (res interface{}, isNaN bool) {
	switch data.(type) {
	case *uint:
		if data.(*uint) == nil || *data.(*uint) >= uintThreshold {
			return -1, true
		}
		res = *data.(*uint)
	case *uint64:
		if data.(*uint64) == nil || *data.(*uint64) >= uint64Threshold {
			return -1, true
		}
		res = *data.(*uint64)
	case *float64:
		if data.(*float64) == nil || *data.(*float64) >= uint64Threshold {
			return -1, true
		}
		res = *data.(*float64)
	case *int:
		if data.(*int) == nil {
			return -1, true
		}
		res = *data.(*int)
	default:
		return -1, true
	}
	return res, false
}

func isIgnore(metricName string) bool {
	_, ok := common.Config().MetricFilter[metricName]
	return ok
}

func getEndPoint() (endPoint string) {
	globalConfig := common.Config()
	if globalConfig.Metric.EndPoint != "" {
		endPoint = globalConfig.Metric.EndPoint
	} else {
		endPoint, _ = getHostName() // 默认主机名，可以配置
	}
	return endPoint
}

func getFalconAgent() (falconAgent string) {
	url := common.Config().Falcon.Agent
	if url != "" {
		falconAgent = url
	} else {
		falconAgent = "http://127.0.0.1:1988/v1/push"
	}
	return falconAgent
}

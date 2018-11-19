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

var falconAgent string

// MetaData 发送到falcon的数据结构
type MetaData struct {
	Metric      string      `json:"metric"`      //
	Endpoint    string      `json:"endpoint"`    //hostname
	Timestamp   int64       `json:"timestamp"`   // s
	Step        int         `json:"step"`        // interval
	Value       interface{} `json:"value"`       // number or string
	CounterType string      `json:"counterType"` // GAUGE  原值, COUNTER 差值(ps)
	TAGS        string      `json:"tags"`        // port=3306,k=v
}

const (
	gpuid int = iota
	sum
	ave
	tagCount
)

// borrow from NVIDIA dcgm
const (
	uintThreshold   = 0x7ffffff0         // 2147483632
	uint64Threshold = 0x7ffffffffffffff0 // 9223372036854775792
)

var sumMetric = map[string]int{
	"FBUsed":   -1,
	"GPUUtils": -1,
	"MemUtils": -1,
}

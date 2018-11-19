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
package fetch

import (
	"fmt"

	"github.com/NVIDIA/gpu-monitoring-tools/bindings/go/dcgm"
	"github.com/sirupsen/logrus"

	"github.com/open-falcon/gpu-mon/cfg"
)

//Data 获取监控数据
func Data() ([]RawData, error) {
	var rawDataList []RawData

	// 初始化 DCGM
	err := dcgm.Init(dcgm.Standalone, "localhost", "0")
	if err != nil {
		err = fmt.Errorf("initialization of DCGM failed | %s", err)
		cfg.CommonLogger.Error(err)
		return nil, err
	}
	cfg.CommonLogger.Info("Successful initialization of Dcgm")

	defer func() {
		if err := dcgm.Shutdown(); err != nil {
			cfg.CommonLogger.Error(err) // Calls os.Exit(1) after loggings
		}
	}()

	// 获取所有GPU的数量
	count, err := dcgm.GetAllDeviceCount()
	if err != nil {
		return nil, fmt.Errorf("failed to get GPU device number | %s", err)
	}
	cfg.CommonLogger.Info("successful get all gpu device")

	// 获取支持DCGM的gpu设备，并构建DCGM map用于遍历
	gpus, err := dcgm.GetSupportedDevices()
	if err != nil {
		err = fmt.Errorf("failed to get GPU devices supporting DCGM | %s", err)
		return nil, err
	}
	cfg.CommonLogger.Info("successful get GPU devices supporting dcgm")

	deviceSupportDCGM := make(map[uint]uint)
	for _, gpuID := range gpus {
		deviceSupportDCGM[gpuID] = gpuID
		cfg.CommonLogger.WithFields(logrus.Fields{
			"gpuID": gpuID,
		}).Info("Set deviceSupportDCGM map: dcgm supported device")
	}

	// 获取GPU设备信息
	rawDataList = make([]RawData, count)
	for i := uint(0); i < count; i++ {
		var dcgmSupported = -1  //默认不支持DCGM
		var values MetricValues //
		if gpuID, ok := deviceSupportDCGM[i]; ok {
			//gpu支持DCGM
			dcgmSupported = 1
			cfg.CommonLogger.WithFields(logrus.Fields{
				"gpuId": gpuID,
			}).Info("device support dcgm")
			values, err = fetchValues(gpuID)
			if err != nil {
				cfg.CommonLogger.Error(err)
			}
			cfg.CommonLogger.WithFields(logrus.Fields{
				"dcgmSupportedGPU": gpuID,
			}).Info("get DCGM supported Device Info")
			values.DcgmSupported = &dcgmSupported
		} else {
			cfg.CommonLogger.WithFields(logrus.Fields{
				"gpuId": i,
			}).Error("device do not support dcgm")
			values.DcgmSupported = &dcgmSupported
		}

		rawData := RawData{
			GpuID:  i,
			Values: values,
		}
		rawDataList[i] = rawData
		cfg.CommonLogger.WithFields(logrus.Fields{
			"gpuId": i,
		}).Info("successful fetch gpu info")
	}
	return rawDataList, nil
}

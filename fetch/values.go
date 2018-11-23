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

/*
#cgo LDFLAGS: -ldl -Wl,--unresolved-symbols=ignore-in-object-files
#cgo CFLAGS: -I /usr/include

#include "dcgm_agent.h"
#include "dcgm_structs.h"
*/
import "C"
import (
	"fmt"
	"github.com/open-falcon/gpu-mon/common"
	"math/rand"
	"unsafe"

	// 需要使用dcgm库的一些私有函数及变量
	_ "github.com/NVIDIA/gpu-monitoring-tools/bindings/go/dcgm"
)

//go:linkname fieldGroupCreate github.com/open-falcon/gpu-mon/vendor/github.com/NVIDIA/gpu-monitoring-tools/bindings/go/dcgm.fieldGroupCreate
//go:linkname fieldGroupDestroy github.com/open-falcon/gpu-mon/vendor/github.com/NVIDIA/gpu-monitoring-tools/bindings/go/dcgm.fieldGroupDestroy
//go:linkname watchFields github.com/open-falcon/gpu-mon/vendor/github.com/NVIDIA/gpu-monitoring-tools/bindings/go/dcgm.watchFields
//go:linkname destroyGroup github.com/open-falcon/gpu-mon/vendor/github.com/NVIDIA/gpu-monitoring-tools/bindings/go/dcgm.destroyGroup
//go:linkname errorString github.com/open-falcon/gpu-mon/vendor/github.com/NVIDIA/gpu-monitoring-tools/bindings/go/dcgm.errorString
//go:linkname uintPtrUnsafe github.com/open-falcon/gpu-mon/vendor/github.com/NVIDIA/gpu-monitoring-tools/bindings/go/dcgm.uintPtrUnsafe
//go:linkname handle github.com/open-falcon/gpu-mon/vendor/github.com/NVIDIA/gpu-monitoring-tools/bindings/go/dcgm.handle
//go:linkname uint64PtrUnsafe github.com/open-falcon/gpu-mon/vendor/github.com/NVIDIA/gpu-monitoring-tools/bindings/go/dcgm.uint64PtrUnsafe
//go:linkname dblToFloatUnsafe github.com/open-falcon/gpu-mon/vendor/github.com/NVIDIA/gpu-monitoring-tools/bindings/go/dcgm.dblToFloatUnsafe

func fieldGroupCreate(fieldsGroupName string, fields []C.ushort, count int) (fieldsId fieldHandle, err error)
func fieldGroupDestroy(fieldsGroup fieldHandle) (err error)
func watchFields(gpuId uint, fieldsGroup fieldHandle, groupName string) (groupId groupHandle, err error)
func destroyGroup(groupId groupHandle) (err error)
func errorString(result C.dcgmReturn_t) error
func uintPtrUnsafe(p unsafe.Pointer) *uint
func uint64PtrUnsafe(p unsafe.Pointer) *uint64
func dblToFloatUnsafe(p unsafe.Pointer) *float64

var handle dcgmHandle

type dcgmHandle struct{ handle C.dcgmHandle_t }
type fieldHandle struct{ handle C.dcgmFieldGrp_t }
type groupHandle struct{ handle C.dcgmGpuGrp_t }

func fetchValues(gpuID uint) (values MetricValues, err error) {
	var deviceFields []C.ushort = make([]C.ushort, fieldsCount)

	deviceFields[gpuUtils] = C.DCGM_FI_DEV_GPU_UTIL
	deviceFields[memUtils] = C.DCGM_FI_DEV_MEM_COPY_UTIL
	deviceFields[encoder] = C.DCGM_FI_DEV_ENC_UTIL
	deviceFields[decoder] = C.DCGM_FI_DEV_DEC_UTIL
	deviceFields[singleBitError] = C.DCGM_FI_DEV_ECC_SBE_AGG_TOTAL
	deviceFields[doubleBitError] = C.DCGM_FI_DEV_ECC_DBE_AGG_TOTAL
	deviceFields[smClock] = C.DCGM_FI_DEV_SM_CLOCK
	deviceFields[memClock] = C.DCGM_FI_DEV_MEM_CLOCK
	deviceFields[bAR1Used] = C.DCGM_FI_DEV_BAR1_USED
	deviceFields[rx] = C.DCGM_FI_DEV_PCIE_RX_THROUGHPUT
	deviceFields[tx] = C.DCGM_FI_DEV_PCIE_TX_THROUGHPUT
	deviceFields[replays] = C.DCGM_FI_DEV_PCIE_REPLAY_COUNTER
	deviceFields[performance] = C.DCGM_FI_DEV_PSTATE
	deviceFields[fanSpeed] = C.DCGM_FI_DEV_FAN_SPEED
	deviceFields[powerUsed] = C.DCGM_FI_DEV_POWER_USAGE
	deviceFields[fBUsed] = C.DCGM_FI_DEV_FB_USED
	deviceFields[deviceTemperature] = C.DCGM_FI_DEV_GPU_TEMP
	deviceFields[memTemperature] = C.DCGM_FI_DEV_MEMORY_TEMP
	deviceFields[slowdownTemperature] = C.DCGM_FI_DEV_SLOWDOWN_TEMP
	deviceFields[shutdownTemperature] = C.DCGM_FI_DEV_SHUTDOWN_TEMP
	deviceFields[powerCurrentLimit] = C.DCGM_FI_DEV_POWER_MGMT_LIMIT
	deviceFields[powerMinManLimit] = C.DCGM_FI_DEV_POWER_MGMT_LIMIT_MIN
	deviceFields[powerMaxManLimit] = C.DCGM_FI_DEV_POWER_MGMT_LIMIT_MAX
	deviceFields[powerDefaultManLimit] = C.DCGM_FI_DEV_POWER_MGMT_LIMIT_DEF
	deviceFields[powerEnforcedLimit] = C.DCGM_FI_DEV_ENFORCED_POWER_LIMIT
	deviceFields[powerViolationTime] = C.DCGM_FI_DEV_POWER_VIOLATION
	deviceFields[fBtotal] = C.DCGM_FI_DEV_FB_TOTAL
	deviceFields[fBfree] = C.DCGM_FI_DEV_FB_FREE
	deviceFields[memAppClock] = C.DCGM_FI_DEV_APP_MEM_CLOCK
	deviceFields[sMAppClock] = C.DCGM_FI_DEV_APP_SM_CLOCK
	deviceFields[videoEnClock] = C.DCGM_FI_DEV_VIDEO_CLOCK
	deviceFields[rPSingleError] = C.DCGM_FI_DEV_RETIRED_SBE
	deviceFields[rPDoubleError] = C.DCGM_FI_DEV_RETIRED_DBE
	deviceFields[packagePend] = C.DCGM_FI_DEV_RETIRED_PENDING
	deviceFields[sBErrors] = C.DCGM_FI_DEV_ECC_SBE_VOL_TOTAL
	deviceFields[dBErrors] = C.DCGM_FI_DEV_ECC_DBE_VOL_TOTAL
	deviceFields[memSBAErrors] = C.DCGM_FI_DEV_ECC_SBE_AGG_DEV
	deviceFields[memDBAErrors] = C.DCGM_FI_DEV_ECC_DBE_AGG_DEV
	deviceFields[deviceMemSBErrors] = C.DCGM_FI_DEV_ECC_SBE_VOL_DEV
	deviceFields[deviceMemDBErrors] = C.DCGM_FI_DEV_ECC_DBE_VOL_DEV
	deviceFields[registerSBErrors] = C.DCGM_FI_DEV_ECC_SBE_VOL_REG
	deviceFields[registerDBErrors] = C.DCGM_FI_DEV_ECC_DBE_VOL_REG

	fieldsName := fmt.Sprintf("devStatusFields%d", rand.Uint64())
	fieldsID, err := fieldGroupCreate(fieldsName, deviceFields, fieldsCount)
	if err != nil {
		return values, fmt.Errorf("failed to create dcgm group | %v", err)
	}

	common.Logger.Info("successfully create dcgm group")

	groupName := fmt.Sprintf("devStatus%d", rand.Uint64())
	groupID, err := watchFields(gpuID, fieldsID, groupName)
	if err != nil {
		err = fmt.Errorf("watch fields failed | %v", err)
		destroyErr := fieldGroupDestroy(fieldsID)
		if destroyErr != nil {
			destroyErr = fmt.Errorf("destroy group fields failed | %v", err)
			common.Logger.Error(destroyErr)
		}
		return values, err
	}

	orivalues := make([]C.dcgmFieldValue_t, fieldsCount)
	result := C.dcgmGetLatestValuesForFields(handle.handle, C.int(gpuID),
		&deviceFields[0], C.uint(fieldsCount), &orivalues[0])

	if err = errorString(result); err != nil {
		_ = fieldGroupDestroy(fieldsID)
		_ = destroyGroup(groupID)
		return values, fmt.Errorf("get device status failed | %v", err)
	}

	// 返回metricValue
	values = MetricValues{
		GPUUtils:             uintPtrUnsafe(unsafe.Pointer(&orivalues[gpuUtils].value)),
		MemUtils:             uintPtrUnsafe(unsafe.Pointer(&orivalues[memUtils].value)),
		Encoder:              uintPtrUnsafe(unsafe.Pointer(&orivalues[encoder].value)),
		Decoder:              uintPtrUnsafe(unsafe.Pointer(&orivalues[decoder].value)),
		SingleBitError:       uintPtrUnsafe(unsafe.Pointer(&orivalues[singleBitError].value)),
		DoubleBitError:       uintPtrUnsafe(unsafe.Pointer(&orivalues[doubleBitError].value)),
		SmClock:              uintPtrUnsafe(unsafe.Pointer(&orivalues[smClock].value)),
		MemClock:             uintPtrUnsafe(unsafe.Pointer(&orivalues[memClock].value)),
		BAR1Used:             uintPtrUnsafe(unsafe.Pointer(&orivalues[bAR1Used].value)),
		Rx:                   uint64PtrUnsafe(unsafe.Pointer(&orivalues[rx].value)),
		Tx:                   uint64PtrUnsafe(unsafe.Pointer(&orivalues[tx].value)),
		Replays:              uint64PtrUnsafe(unsafe.Pointer(&orivalues[replays].value)),
		Performance:          uintPtrUnsafe(unsafe.Pointer(&orivalues[performance].value)),
		FanSpeed:             uintPtrUnsafe(unsafe.Pointer(&orivalues[fanSpeed].value)),
		PowerUsed:            dblToFloatUnsafe(unsafe.Pointer(&orivalues[powerUsed].value)),
		FBUsed:               uintPtrUnsafe(unsafe.Pointer(&orivalues[fBUsed].value)),
		DeviceTemperature:    uintPtrUnsafe(unsafe.Pointer(&orivalues[deviceTemperature].value)),
		MemTemperature:       uintPtrUnsafe(unsafe.Pointer(&orivalues[memTemperature].value)),
		SlowdownTemperature:  uintPtrUnsafe(unsafe.Pointer(&orivalues[slowdownTemperature].value)),
		ShutdownTemperature:  uintPtrUnsafe(unsafe.Pointer(&orivalues[shutdownTemperature].value)),
		PowerCurrentLimit:    dblToFloatUnsafe(unsafe.Pointer(&orivalues[powerCurrentLimit].value)),
		PowerMinManLimit:     dblToFloatUnsafe(unsafe.Pointer(&orivalues[powerMinManLimit].value)),
		PowerMaxManLimit:     dblToFloatUnsafe(unsafe.Pointer(&orivalues[powerMaxManLimit].value)),
		PowerDefaultManLimit: dblToFloatUnsafe(unsafe.Pointer(&orivalues[powerDefaultManLimit].value)),
		PowerEnforcedLimit:   dblToFloatUnsafe(unsafe.Pointer(&orivalues[powerEnforcedLimit].value)),
		PowerViolationTime:   dblToFloatUnsafe(unsafe.Pointer(&orivalues[powerViolationTime].value)),
		FBtotal:              uintPtrUnsafe(unsafe.Pointer(&orivalues[fBtotal].value)),
		FBfree:               uintPtrUnsafe(unsafe.Pointer(&orivalues[fBfree].value)),
		MemAppClock:          uintPtrUnsafe(unsafe.Pointer(&orivalues[memAppClock].value)),
		SMAppClock:           uintPtrUnsafe(unsafe.Pointer(&orivalues[sMAppClock].value)),
		VideoEnClock:         uintPtrUnsafe(unsafe.Pointer(&orivalues[videoEnClock].value)),
		RPSingleError:        uintPtrUnsafe(unsafe.Pointer(&orivalues[rPSingleError].value)),
		RPDoubleError:        uintPtrUnsafe(unsafe.Pointer(&orivalues[rPDoubleError].value)),
		PackagePend:          uintPtrUnsafe(unsafe.Pointer(&orivalues[packagePend].value)),
		SBErrors:             uintPtrUnsafe(unsafe.Pointer(&orivalues[sBErrors].value)),
		DBErrors:             uintPtrUnsafe(unsafe.Pointer(&orivalues[dBErrors].value)),
		MemSBAErrors:         uintPtrUnsafe(unsafe.Pointer(&orivalues[memSBAErrors].value)),
		MemDBAErrors:         uintPtrUnsafe(unsafe.Pointer(&orivalues[memDBAErrors].value)),
		DeviceMemSBErrors:    uintPtrUnsafe(unsafe.Pointer(&orivalues[deviceMemSBErrors].value)),
		DeviceMemDBErrors:    uintPtrUnsafe(unsafe.Pointer(&orivalues[deviceMemDBErrors].value)),
		RegisterSBErrors:     uintPtrUnsafe(unsafe.Pointer(&orivalues[registerSBErrors].value)),
		RegisterDBErrors:     uintPtrUnsafe(unsafe.Pointer(&orivalues[registerDBErrors].value)),
	}
	return values, err
}

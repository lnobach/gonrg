package sml_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/lnobach/gonrg/d0"
	"github.com/lnobach/gonrg/sml"
	"github.com/stretchr/testify/assert"
)

func TestParsing(t *testing.T) {

	testMeter(t, "DrNeuhaus_SMARTY_ix-130", "1DNT4294975536", 7)
	testMeter(t, "DZG_DVS-7420.2V.G2_mtr2", "1DZG0060694611", 5)
	testMeter(t, "dzg_dwsb20_2th_3byte", "1DZG0040051478", 5)
	testMeter(t, "EMH_eHZ-GW8E2A500AK2", "02280816", 6)
	testMeter(t, "EMH_eHZ-HW8E2A5L0EK2P_2", "EMH10491535441", 7)
	testMeter(t, "HOLLEY_DTZ541-BDBA_with_PIN", "1HLY8590814182", 4)
	testMeter(t, "ISKRA_MT175_D1A52-V22-K0t", "ISK60656331162", 13)
	testMeter(t, "ISKRA_MT175_eHZ", "1ISK0067362659", 10)
	testMeter(t, "ISKRA_MT631-D1A52-K0z-H01_with_PIN", "1ISK0075126084", 5)
	testMeter(t, "ISKRA_MT691_eHZ-MS2020", "1ISK0070409925", 4)
	testMeter(t, "EMH_eHZ-HW8E2A5L0EK2P_2", "EMH10491535441", 7)
	testMeter(t, "ITRON_OpenWay-3.HZ_with_PIN", "1ITR0055046208", 4)

}

func testMeter(t *testing.T, binfile, expectedDeviceID string, obisEntries int) {

	mt := time.Unix(12345678, 0)
	raw, err := sml.RawDataFromFile(fmt.Sprintf("../testdata/sml/%s.bin", binfile))
	assert.NoError(t, err)
	assert.NoError(t, err)
	result, err := d0.ParseOBISList(&d0.ParseConfig{}, raw, mt)
	assert.NoError(t, err)
	if result != nil {
		assert.Equal(t, expectedDeviceID, result.DeviceID)
		assert.Len(t, result.List, obisEntries)
	}

}

package modem_alcatel_mw40v

import (
	"fmt"
	_ "io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

func runTestServer(t *testing.T, expectedUrl string, expectedMethod string, testFile string) *httptest.Server {
	// HTTP Server for test
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, testFile)

		if fmt.Sprintf("%s", r.URL) != expectedUrl {
			t.Logf("Expected URL: %s, got: %s\n", expectedUrl, r.URL)
			t.Fail()
		}
		if r.Method != expectedMethod {
			t.Logf("Expected Method: %s, got: %s\n", expectedMethod, r.Method)
			t.Fail()
		}
	}))
	return ts
}

func TestGetSystemInfo(t *testing.T) {
	expectedResults := SystemInfo{
		SoftwareVersion: "MW40_E6_02.00_05",
		HardwareVersion: "MW40-V-V1.0",
		MacAddress:      "c4:43:13:c5:12:34",
		IMEI:            "123456789012345",
		IMSI:            "987654321098765",
		ICCID:           "0123456789012345678p",
	}

	expectedUrl := "/jrd/webapi?api=GetSystemInfo"
	expectedMethod := "POST"

	ts := runTestServer(t, expectedUrl, expectedMethod, "testdata/getSystemInfo.json")
	defer ts.Close()

	modem := New(ts.URL)

	err, systemInfo := modem.GetSystemInfo()
	if err != nil {
		t.Logf("[TestGetSystemStatus] Error: %s", err.Error())
		t.Fail()
	}

	if expectedResults.SoftwareVersion != systemInfo.SoftwareVersion {
		t.Logf("Expected software version: [%+v], got: %s", expectedResults.SoftwareVersion, systemInfo.SoftwareVersion)
		t.Fail()
	}
	if expectedResults.HardwareVersion != systemInfo.HardwareVersion {
		t.Logf("Expected hardware version: %s, got: %s", expectedResults.HardwareVersion, systemInfo.HardwareVersion)
		t.Fail()
	}
	if expectedResults.MacAddress != systemInfo.MacAddress {
		t.Logf("Expected mac address: %s, got: %s", expectedResults.MacAddress, systemInfo.MacAddress)
		t.Fail()
	}
	if expectedResults.IMEI != systemInfo.IMEI {
		t.Logf("Expected IMEI: %s, got: %s", expectedResults.IMEI, systemInfo.IMEI)
		t.Fail()
	}
	if expectedResults.IMSI != systemInfo.IMSI {
		t.Logf("Expected IMSI: %s, got: %s", expectedResults.IMSI, systemInfo.IMSI)
		t.Fail()
	}
	if expectedResults.ICCID != systemInfo.ICCID {
		t.Logf("Expected ICCID: %s, got: %s", expectedResults.ICCID, systemInfo.ICCID)
		t.Fail()
	}
}

func TestGetSystemStatus(t *testing.T) {
	expectedResults := SystemStatus{
		BatteryCapacity:   100,
		BatteryLevel:      4,
		Roaming:           1,
		DomesticRoaming:   1,
		SignalStrength:    0,
		CurrentConnection: 5,
		TotalConnection:   6,
	}
	expectedUrl := "/jrd/webapi?api=GetSystemStatus"
	expectedMethod := "POST"

	ts := runTestServer(t, expectedUrl, expectedMethod, "testdata/getSystemStatus.json")
	defer ts.Close()

	modem := New(ts.URL)

	err, systemStatus := modem.GetSystemStatus()
	if err != nil {
		t.Logf("[TestGetSystemStatus] Error: %s", err.Error())
		t.Fail()
	}

	if expectedResults.BatteryCapacity != systemStatus.BatteryCapacity {
		t.Logf("Expected battery capacity: %f, got: %f", expectedResults.BatteryCapacity, systemStatus.BatteryCapacity)
		t.Fail()
	}
	if expectedResults.BatteryLevel != systemStatus.BatteryLevel {
		t.Logf("Expected battery level: %f, got: %f", expectedResults.BatteryLevel, systemStatus.BatteryLevel)
		t.Fail()
	}
	if expectedResults.Roaming != systemStatus.Roaming {
		t.Logf("Expected roaming: %f, got: %f", expectedResults.Roaming, systemStatus.Roaming)
		t.Fail()
	}
	if expectedResults.DomesticRoaming != systemStatus.DomesticRoaming {
		t.Logf("Expected domestic roaming: %f, got: %f", expectedResults.DomesticRoaming, systemStatus.DomesticRoaming)
		t.Fail()
	}
	if expectedResults.SignalStrength != systemStatus.SignalStrength {
		t.Logf("Expected signal strength: %f, got: %f", expectedResults.SignalStrength, systemStatus.SignalStrength)
		t.Fail()
	}
	if expectedResults.CurrentConnection != systemStatus.CurrentConnection {
		t.Logf("Expected current connection: %f, got: %f", expectedResults.CurrentConnection, systemStatus.CurrentConnection)
		t.Fail()
	}
	if expectedResults.TotalConnection != systemStatus.TotalConnection {
		t.Logf("Expected total connection: %f, got: %f", expectedResults.TotalConnection, systemStatus.TotalConnection)
		t.Fail()
	}
}

func TestGetConnectionState(t *testing.T) {
	expectedResults := ConnectionState{
		ConnectionStatus: 2,
		ConProfileError:  0,
		IPv4Address:      "89.180.91.116",
		IPv6Address:      "0::0",
		SpeedDownload:    67,
		SpeedUpload:      83,
		DownloadRate:     100000000,
		UploadRate:       50000000,
		ConnectionTime:   1676,
		UploadBytes:      1176744,
		DownloadBytes:    2656630,
	}

	expectedUrl := "/jrd/webapi?api=GetConnectionState"
	expectedMethod := "POST"

	ts := runTestServer(t, expectedUrl, expectedMethod, "testdata/getConnectionState.json")
	defer ts.Close()

	modem := New(ts.URL)

	err, connectionState := modem.GetConnectionState()
	if err != nil {
		t.Logf("[TestGetConnectionState] Error: %s", err.Error())
		t.Fail()
	}

	if expectedResults.IPv4Address != connectionState.IPv4Address {
		t.Logf("Expected IPv4Address: %s, got: %s", expectedResults.IPv4Address, connectionState.IPv4Address)
		t.Fail()
	}
	if expectedResults.IPv6Address != connectionState.IPv6Address {
		t.Logf("Expected IPv6Address: %s, got: %s", expectedResults.IPv6Address, connectionState.IPv6Address)
		t.Fail()
	}
	if expectedResults.ConnectionStatus != connectionState.ConnectionStatus {
		t.Logf("Expected ConnectionStatus: %f, got: %f", expectedResults.ConnectionStatus, connectionState.ConnectionStatus)
		t.Fail()
	}
	if expectedResults.ConProfileError != connectionState.ConProfileError {
		t.Logf("Expected ConProfileError: %f, got: %f", expectedResults.ConProfileError, connectionState.ConProfileError)
		t.Fail()
	}
	if expectedResults.SpeedDownload != connectionState.SpeedDownload {
		t.Logf("Expected SpeedDownload: %f, got: %f", expectedResults.SpeedDownload, connectionState.SpeedDownload)
		t.Fail()
	}
	if expectedResults.SpeedUpload != connectionState.SpeedUpload {
		t.Logf("Expected SpeedUpload: %f, got: %f", expectedResults.SpeedUpload, connectionState.SpeedUpload)
		t.Fail()
	}
	if expectedResults.DownloadRate != connectionState.DownloadRate {
		t.Logf("Expected DownloadRate: %f, got: %f", expectedResults.DownloadRate, connectionState.DownloadRate)
		t.Fail()
	}
	if expectedResults.UploadRate != connectionState.UploadRate {
		t.Logf("Expected UploadRate: %f, got: %f", expectedResults.UploadRate, connectionState.UploadRate)
		t.Fail()
	}
	if expectedResults.ConnectionTime != connectionState.ConnectionTime {
		t.Logf("Expected ConnectionTime: %f, got: %f", expectedResults.ConnectionTime, connectionState.ConnectionTime)
		t.Fail()
	}
	if expectedResults.UploadBytes != connectionState.UploadBytes {
		t.Logf("Expected UploadBytes: %f, got: %f", expectedResults.UploadBytes, connectionState.UploadBytes)
		t.Fail()
	}
	if expectedResults.DownloadBytes != connectionState.DownloadBytes {
		t.Logf("Expected DownloadBytes: %f, got: %f", expectedResults.DownloadBytes, connectionState.DownloadBytes)
		t.Fail()
	}
}

func TestGetSMSStorageState(t *testing.T) {
	expectedResult := SMSStorageState{
		UnreadReport:   0,
		LeftCount:      100,
		MaxCount:       100,
		TUseCount:      0,
		UnreadSMSCount: 0,
	}

	expectedUrl := "/jrd/webapi?api=GetSMSStorageState"
	expectedMethod := "POST"

	ts := runTestServer(t, expectedUrl, expectedMethod, "testdata/getSMSStorageState.json")
	defer ts.Close()

	modem := New(ts.URL)

	err, smsStorageState := modem.GetSMSStorageState()
	if err != nil {
		t.Logf("[TestGetConnectionState] Error: %s", err.Error())
		t.Fail()
		return
	}

	if expectedResult.UnreadSMSCount != smsStorageState.UnreadSMSCount {
		t.Logf("Expected UnreadSMSCount: %f, got: %f", expectedResult.UnreadSMSCount, smsStorageState.UnreadSMSCount)
		t.Fail()
	}
	if expectedResult.TUseCount != smsStorageState.TUseCount {
		t.Logf("Expected TUseCount: %f, got: %f", expectedResult.TUseCount, smsStorageState.TUseCount)
		t.Fail()
	}
	if expectedResult.MaxCount != smsStorageState.MaxCount {
		t.Logf("Expected MaxCount: %f, got: %f", expectedResult.MaxCount, smsStorageState.MaxCount)
		t.Fail()
	}
	if expectedResult.LeftCount != smsStorageState.LeftCount {
		t.Logf("Expected LeftCount: %f, got: %f", expectedResult.LeftCount, smsStorageState.LeftCount)
		t.Fail()
	}
	if expectedResult.UnreadReport != smsStorageState.UnreadReport {
		t.Logf("Expected UnreadReport: %f, got: %f", expectedResult.UnreadReport, smsStorageState.UnreadReport)
		t.Fail()
	}
}

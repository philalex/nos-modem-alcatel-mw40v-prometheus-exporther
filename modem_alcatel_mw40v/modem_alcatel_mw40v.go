package modem_alcatel_mw40v

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	log "github.com/sirupsen/logrus"
)

// HTTP_TIMEOUT in second
const HTTP_TIMEOUT = 10

type Modem struct {
	Url string
}

type AuthResponse struct {
	Result struct {
		Token float64
	}
	Id string
}

// system status
type SystemStatusResult struct {
	Result SystemStatus `json:"result"`
	Id     string       `json:"id"`
}

type SystemStatus struct {
	BatteryCapacity   float64 `json:"bat_cap"`
	BatteryLevel      float64 `json:"bat_level"`
	Roaming           float64 `json:"Roaming"`
	DomesticRoaming   float64 `json:"Domestic_Roaming"`
	SignalStrength    float64 `json:"SignalStrength"`
	CurrentConnection float64 `json:"curr_num"`
	TotalConnection   float64 `json:"TotalConnNum"`
}

// system info
type SystemInfoResult struct {
	Result SystemInfo `json:"result"`
	Id     string     `json:"id"`
}

type SystemInfo struct {
	SoftwareVersion string `json:"SwVersion"`
	HardwareVersion string `json:"HwVersion"`
	MacAddress      string `json:"MacAddress"`
	IMEI            string `json:"IMEI"`
	IMSI            string `json:"IMSI"`
	ICCID           string `json:"ICCID"`
}

// Connection state
type ConnectionStateResult struct {
	Result ConnectionState `json:"result"`
	Id     string          `json:"id"`
}

type ConnectionState struct {
	ConnectionStatus float64 `json:"ConnectionStatus"`
	ConProfileError  float64 `json:"Conprofileerror"`
	IPv4Address      string  `json:"IPv4Adrress"`
	IPv6Address      string  `json:"IPv6Adrress"`
	SpeedDownload    float64 `json:"Speed_Dl"`
	SpeedUpload      float64 `json:"Speed_Ul"`
	DownloadRate     float64 `json:"DlRate"`
	UploadRate       float64 `json:"UlRate"`
	ConnectionTime   float64 `json:"ConnectionTime"`
	UploadBytes      float64 `json:"UlBytes"`
	DownloadBytes    float64 `json:"DlBytes"`
}

// SMS storage state
type SMSStorageStateResult struct {
	Result SMSStorageState `json:"result"`
	Id     string          `json:"id"`
}

type SMSStorageState struct {
	UnreadReport   float64 `json:"UnreadReport"`
	LeftCount      float64 `json:"LeftCount"`
	MaxCount       float64 `json:"MaxCount"`
	TUseCount      float64 `json:"TUseCount"`
	UnreadSMSCount float64 `json:"UnreadSMSCount"`
}

func New(url string) *Modem {
	tmpUrl := url
	if strings.HasSuffix(tmpUrl, "/") == false {
		tmpUrl = tmpUrl + "/"
	}

	return &Modem{Url: tmpUrl}
}

// GetSystemInfo get modem identification Software & hardware version, mac address, IMEI, IMSI and ICCID
func (modem *Modem) GetSystemInfo() (*SystemInfo, error) {
	var systemInfoResult SystemInfoResult

	var jsonStr = []byte(`{"jsonrpc":"2.0","method":"GetSystemInfo","params":null,"id":"13.1"}`)

	body, err := modem.postRequest("jrd/webapi?api=GetSystemInfo", jsonStr)
	if err != nil {
		return nil, err
	}
	log.Debugf("[GetSystemInfo] Body: %+s\n", string(body))

	err = json.Unmarshal(body, &systemInfoResult)
	if err != nil {
		return nil, err
	}

	// Remove \n suffix
	systemInfoResult.Result.SoftwareVersion = strings.TrimSuffix(systemInfoResult.Result.SoftwareVersion, "\n")
	systemInfoResult.Result.MacAddress = strings.TrimSuffix(systemInfoResult.Result.MacAddress, "\n")

	return &systemInfoResult.Result, nil
}

// GetSystemStatus get modem status: battery capacity, battery level, roaming, domestic roaming, signal strength, number device(s) connected, total device(s) connected
func (modem *Modem) GetSystemStatus() (*SystemStatus, error) {
	var systemStatusResult SystemStatusResult
	var jsonStr = []byte(`{"jsonrpc":"2.0","method":"GetSystemStatus","params":null,"id":"13.4"}`)
	body, err := modem.postRequest("jrd/webapi?api=GetSystemStatus", jsonStr)
	if err != nil {
		return nil, err
	}
	log.Debugf("[GetSystemStatus] Body: %+s\n", string(body))

	err = json.Unmarshal(body, &systemStatusResult)
	if err != nil {
		return nil, err
	}

	return &systemStatusResult.Result, nil
}

// GetConnectionState
func (modem *Modem) GetConnectionState() (*ConnectionState, error) {
	var connectionStateResult ConnectionStateResult
	var jsonStr = []byte(`{"jsonrpc":"2.0","method":"GetConnectionState","params":null,"id":"3.1"}`)
	body, err := modem.postRequest("jrd/webapi?api=GetConnectionState", jsonStr)
	if err != nil {
		return nil, err
	}
	log.Debugf("[GetConnectionState] Body: %+s\n", string(body))

	err = json.Unmarshal(body, &connectionStateResult)
	if err != nil {
		return nil, err
	}

	return &connectionStateResult.Result, nil
}

// GetSMSStorageState
func (modem *Modem) GetSMSStorageState() (*SMSStorageState, error) {
	var smsStorageStateResult SMSStorageStateResult
	var jsonStr = []byte(`{"jsonrpc":"2.0","method":"GetSMSStorageState","params":null,"id":"6.4"}`)
	body, err := modem.postRequest("jrd/webapi?api=GetSMSStorageState", jsonStr)
	if err != nil {
		return nil, err
	}
	log.Debugf("[GetSMSStorageState] Body: %+s\n", string(body))

	err = json.Unmarshal(body, &smsStorageStateResult)
	if err != nil {
		return nil, err
	}

	return &smsStorageStateResult.Result, nil
}

// postRequest
func (modem *Modem) postRequest(url string, jsonStr []byte) ([]byte, error) {
	var emptyByte []byte

	requestUrl := modem.Url + url

	req, err := http.NewRequest("POST", requestUrl, bytes.NewBuffer(jsonStr))

	client := &http.Client{
		Timeout: HTTP_TIMEOUT * time.Second,
	}

	resp, err := client.Do(req)
	if err != nil {
		return emptyByte, err
	}
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)

	return body, nil
}

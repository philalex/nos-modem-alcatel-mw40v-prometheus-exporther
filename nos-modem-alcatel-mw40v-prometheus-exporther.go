package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	log "github.com/sirupsen/logrus"

	"nos-modem-alcatel-mw40v-prometheus-exporther/modem_alcatel_mw40v"
)

var BUILD_DATE = "Undefined"
var GIT_BRANCH = "Undefined"
var GIT_HASH = "Undefined"

var (
	batteryCapacityGauge = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "battery_capacity_percent",
			Help: "Battery capacity",
		},
		[]string{"IMEI", "IMSI", "MacAddress"},
	)
	batteryLevelGauge = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "battery_level",
			Help: "Battery level",
		},
		[]string{"IMEI", "IMSI", "MacAddress"},
	)
	currentConnectionGauge = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "current_connection_count",
			Help: "Current connection(s)",
		},
		[]string{"IMEI", "IMSI", "MacAddress"},
	)
	totalConnectionGauge = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "total_connection_count",
			Help: "total connection(s)",
		},
		[]string{"IMEI", "IMSI", "MacAddress"},
	)
	connectionStatusGauge = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "connection_status",
			Help: "Connection status",
		},
		[]string{"IMEI", "IMSI", "MacAddress"},
	)
	speedDownloadGauge = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "speed_download",
			Help: "Max speed download",
		},
		[]string{"IMEI", "IMSI", "MacAddress"},
	)
	speedUploadGauge = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "speed_upload",
			Help: "Max speed upload",
		},
		[]string{"IMEI", "IMSI", "MacAddress"},
	)
	downloadRateGauge = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "download_rate",
			Help: "Download rate",
		},
		[]string{"IMEI", "IMSI", "MacAddress"},
	)
	uploadRateGauge = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "upload_rate",
			Help: "Upload rate",
		},
		[]string{"IMEI", "IMSI", "MacAddress"},
	)
	downloadBytesGauge = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "download_bytes",
			Help: "Download bytes",
		},
		[]string{"IMEI", "IMSI", "MacAddress"},
	)
	uploadBytesGauge = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "upload_bytes",
			Help: "Upload bytes",
		},
		[]string{"IMEI", "IMSI", "MacAddress"},
	)
	// SMS storage state
	unreadSMSCountGauge = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "unread_sms_count",
			Help: "Unread SMS",
		},
		[]string{"IMEI", "IMSI", "MacAddress"},
	)
)

func init() {
	// System status
	prometheus.MustRegister(batteryCapacityGauge)
	prometheus.MustRegister(batteryLevelGauge)
	prometheus.MustRegister(currentConnectionGauge)
	prometheus.MustRegister(totalConnectionGauge)
	// Connection state
	prometheus.MustRegister(connectionStatusGauge)
	prometheus.MustRegister(speedDownloadGauge)
	prometheus.MustRegister(speedUploadGauge)
	prometheus.MustRegister(downloadRateGauge)
	prometheus.MustRegister(uploadRateGauge)
	prometheus.MustRegister(downloadBytesGauge)
	prometheus.MustRegister(uploadBytesGauge)
	// SMS
	prometheus.MustRegister(unreadSMSCountGauge)
}

func main() {
	var cmdlineVersion = flag.Bool("v", false, "Version")
	flag.Parse()

	if *cmdlineVersion {
		fmt.Printf("Git hash: %s\n", GIT_HASH)
		fmt.Printf("Git branch: %s\n", GIT_BRANCH)
		fmt.Printf("Build date: %s\n", BUILD_DATE)
		os.Exit(0)
	}

	// Set log level
	logLevel := os.Getenv("LOG_LEVEL")
	if strings.TrimSpace(logLevel) == "" {
		logLevel = "info"
	}
	log.SetOutput(os.Stdout)
	loglevel, err := log.ParseLevel(logLevel)
	if err != nil {
		log.Fatal(err)
	}
	log.SetLevel(loglevel)

	modemUrl := os.Getenv("MODEM_URL")
	if strings.TrimSpace(modemUrl) == "" {
		modemUrl = "http://192.168.1.1"
	}

	strUpdateInterval := os.Getenv("UPDATE_INTERVAL")
	if strings.TrimSpace(strUpdateInterval) == "" {
		strUpdateInterval = "10s"
	}

	updateInterval, err := time.ParseDuration(strUpdateInterval)
	if err != nil {
		log.Fatal(err)
	}

	ticker := time.NewTicker(updateInterval)
	defer ticker.Stop()

	done := make(chan bool)

	modem := modem_alcatel_mw40v.New(modemUrl)
	systemInfo, err := modem.GetSystemInfo()
	if err != nil {
		log.Fatal(err)
	}

	scraper := func() error {
		systemStatus, err := modem.GetSystemStatus()
		if err != nil {
			return err
		}
		batteryCapacityGauge.With(prometheus.Labels{"IMEI": systemInfo.IMEI, "IMSI": systemInfo.IMSI, "MacAddress": systemInfo.MacAddress}).Set(systemStatus.BatteryCapacity)
		batteryLevelGauge.With(prometheus.Labels{"IMEI": systemInfo.IMEI, "IMSI": systemInfo.IMSI, "MacAddress": systemInfo.MacAddress}).Set(systemStatus.BatteryLevel)
		currentConnectionGauge.With(prometheus.Labels{"IMEI": systemInfo.IMEI, "IMSI": systemInfo.IMSI, "MacAddress": systemInfo.MacAddress}).Set(systemStatus.CurrentConnection)
		totalConnectionGauge.With(prometheus.Labels{"IMEI": systemInfo.IMEI, "IMSI": systemInfo.IMSI, "MacAddress": systemInfo.MacAddress}).Set(systemStatus.TotalConnection)

		connectionState, err := modem.GetConnectionState()
		if err != nil {
			return err
		}
		connectionStatusGauge.With(prometheus.Labels{"IMEI": systemInfo.IMEI, "IMSI": systemInfo.IMSI, "MacAddress": systemInfo.MacAddress}).Set(connectionState.ConnectionStatus)
		speedDownloadGauge.With(prometheus.Labels{"IMEI": systemInfo.IMEI, "IMSI": systemInfo.IMSI, "MacAddress": systemInfo.MacAddress}).Set(connectionState.SpeedDownload)
		speedUploadGauge.With(prometheus.Labels{"IMEI": systemInfo.IMEI, "IMSI": systemInfo.IMSI, "MacAddress": systemInfo.MacAddress}).Set(connectionState.SpeedUpload)
		downloadRateGauge.With(prometheus.Labels{"IMEI": systemInfo.IMEI, "IMSI": systemInfo.IMSI, "MacAddress": systemInfo.MacAddress}).Set(connectionState.DownloadRate)
		uploadRateGauge.With(prometheus.Labels{"IMEI": systemInfo.IMEI, "IMSI": systemInfo.IMSI, "MacAddress": systemInfo.MacAddress}).Set(connectionState.UploadRate)
		downloadBytesGauge.With(prometheus.Labels{"IMEI": systemInfo.IMEI, "IMSI": systemInfo.IMSI, "MacAddress": systemInfo.MacAddress}).Set(connectionState.DownloadBytes)
		uploadBytesGauge.With(prometheus.Labels{"IMEI": systemInfo.IMEI, "IMSI": systemInfo.IMSI, "MacAddress": systemInfo.MacAddress}).Set(connectionState.UploadBytes)

		smsStorageState, err := modem.GetSMSStorageState()
		if err != nil {
			return err
		}
		unreadSMSCountGauge.With(prometheus.Labels{"IMEI": systemInfo.IMEI, "IMSI": systemInfo.IMSI, "MacAddress": systemInfo.MacAddress}).Set(smsStorageState.UnreadSMSCount)
		return nil
	}

	// first run
	err = scraper()
	if err != nil {
		log.Fatal(err)
	}

	go func() {
		for {
			select {
			case <-done:
				return

			case t := <-ticker.C:
				log.Debug("Current time: ", t)
				err = scraper()
				if err != nil {
					log.Error(err)
				}
			}
		}
	}()

	http.Handle("/metrics", prometheus.Handler())
	http.ListenAndServe(":8080", nil)

	done <- true
}

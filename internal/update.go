package magellan

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/OpenCHAMI/magellan/internal/util"
)

type UpdateParams struct {
	CollectParams
	FirmwarePath     string
	FirmwareVersion  string
	Component        string
	TransferProtocol string
}

// UpdateFirmwareRemote() uses 'gofish' to update the firmware of a BMC node.
// The function expects the firmware URL, firmware version, and component flags to be
// set from the CLI to perform a firmware update.
func UpdateFirmwareRemote(q *UpdateParams) error {
	url := baseRedfishUrl(&q.CollectParams) + "/redfish/v1/UpdateService/Actions/SimpleUpdate"
	headers := map[string]string{
		"Content-Type":  "application/json",
		"cache-control": "no-cache",
	}
	b := map[string]any{
		"UpdateComponent":  q.Component, // BMC, BIOS
		"TransferProtocol": q.TransferProtocol,
		"ImageURI":         q.FirmwarePath,
	}
	data, err := json.Marshal(b)
	if err != nil {
		return fmt.Errorf("failed to marshal data: %v", err)
	}
	res, body, err := util.MakeRequest(nil, url, "POST", data, headers)
	if err != nil {
		return fmt.Errorf("something went wrong: %v", err)
	} else if res == nil {
		return fmt.Errorf("no response returned (url: %s)", url)
	}
	if len(body) > 0 {
		fmt.Printf("%d: %v\n", res.StatusCode, string(body))
	}
	return nil
}

func GetUpdateStatus(q *UpdateParams) error {
	url := baseRedfishUrl(&q.CollectParams) + "/redfish/v1/UpdateService"
	res, body, err := util.MakeRequest(nil, url, "GET", nil, nil)
	if err != nil {
		return fmt.Errorf("something went wrong: %v", err)
	} else if res == nil {
		return fmt.Errorf("no response returned (url: %s)", url)
	} else if res.StatusCode != http.StatusOK {
		return fmt.Errorf("returned status code %d", res.StatusCode)
	}
	if len(body) > 0 {
		fmt.Printf("%v\n", string(body))
	}
	return nil
}

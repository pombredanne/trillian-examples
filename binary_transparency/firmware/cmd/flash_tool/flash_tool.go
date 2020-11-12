// flash_tool is a util to flash firmware update packages created by the publisher tool onto devices.
//
// Currently, the only device is a dummy device, which simply sorts the firmware+metadata on local disk.
//
// Usage:
//   go run ./cmd/flash_tool/ --logtostderr --dummy_storage_dir=/path/to/dir --update_file=/path/to/update.json
//
// The first time you use this tool there will be no prior firmware metadata
// stored on the device and the tool will fail.  In this case, use the --force
// flag to apply the update anyway thereby creating the metadata.
// Subsequent invocations should then work without needing the --force flag.
package main

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"net/url"
	"os"

	"github.com/golang/glog"
	"github.com/google/trillian-examples/binary_transparency/firmware/api"
	"github.com/google/trillian-examples/binary_transparency/firmware/cmd/flash_tool/devices"
	"github.com/google/trillian-examples/binary_transparency/firmware/internal/client"
)

var (
	logURL     = flag.String("log_url", "http://localhost:8000", "Base URL of the log HTTP API")
	updateFile = flag.String("update_file", "", "File path to read the update package from")
	force      = flag.Bool("force", false, "Ignore errors and force update")
)

func fatal(msg string) {
	if !*force {
		glog.Exit(msg)
	}
	glog.Warning(msg)
}

func main() {
	flag.Parse()

	logURL, err := url.Parse(*logURL)
	if err != nil {
		glog.Exitf("log_url is invalid: %v", err)
	}
	c := &client.Client{LogURL: logURL}

	up, err := readUpdateFileFromFlags()
	if err != nil {
		glog.Exitf("Failed to read update package file: %q", err)
	}

	var dev devices.Device
	dev, err = devices.NewDummyDeviceFromFlags()
	if err != nil {
		switch t := err.(type) {
		case devices.ErrNeedsInit:
			fatal(fmt.Sprintf("Device needs to be force initialised: %q", err))
		default:
			fatal(fmt.Sprintf("Failed to open dummy device: %q", t))
		}
	}

	if err := checkManifestHash(up); err != nil {
		fatal(fmt.Sprintf("Manifest/firmware is corrupt: %q", err))
	}

	if err := checkInclusion(c, up); err != nil {
		fatal(fmt.Sprintf("Could not verify manifest inclusion under update checkpoint: %q", err))
	}

	if err := checkConsistency(c, up, dev); err != nil {
		fatal(fmt.Sprintf("Could not verify consistency between manifest checkpoint and device checkpoint: %q", err))
	}

	glog.Info("Update verified, about to apply to device...")

	if err := dev.ApplyUpdate(up); err != nil {
		glog.Exitf("Failed to apply update to device: %q", err)
	}

	glog.Info("Update applied.")

}

func readUpdateFileFromFlags() (api.UpdatePackage, error) {
	if len(*updateFile) == 0 {
		return api.UpdatePackage{}, errors.New("must specify update_file")
	}

	f, err := os.OpenFile(*updateFile, os.O_RDONLY, os.ModePerm)
	if err != nil {
		glog.Exitf("Failed to open update package file %q: %q", *updateFile, err)
	}
	defer f.Close()

	var up api.UpdatePackage
	if err := json.NewDecoder(f).Decode(&up); err != nil {
		glog.Exitf("Failed to parse update package file: %q", err)
	}
	return up, nil
}

func checkManifestHash(up api.UpdatePackage) error {
	// TODO(al): implement this
	return nil
}

func checkInclusion(c *client.Client, up api.UpdatePackage) error {
	// TODO(al): implement this
	return nil
}

func checkConsistency(c *client.Client, up api.UpdatePackage, dev devices.Device) error {
	// TODO(al): implement this
	return nil
}

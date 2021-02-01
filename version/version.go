package version

import (
	"fmt"
	"sync"

	"encoding/base64"
	"github.com/golang/glog"
)

var (
	Version   = ""
	GitSHA    = "Not Provided"
	BuildTime = "Not Provided"
	Mods      = "Not Provided"
)

var modesOnce sync.Once
var modsInfo = ""

// GetVersionInfo function
func GetVersionInfo() string {
	return fmt.Sprintf("version %v(%v), build time: %v", Version, GitSHA, BuildTime)
}

// GetMods .
func GetMods() string {
	modesOnce.Do(func() {
		if Mods != "" {
			modsInfoData, err := base64.StdEncoding.DecodeString(Mods)
			if err != nil {
				glog.Errorf("decode failed: %v, %v", err, Mods)
				return
			}
			modsInfo = string(modsInfoData)
		}
	})
	return modsInfo
}

// Print function
func Print() {
	glog.Info(GetVersionInfo())
}

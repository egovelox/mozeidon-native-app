package models

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/google/uuid"
)

/*
 * Each file located in this profile directory is a native-app profile.
 * It contains information about a running native-app,
 * notably the ipc-name where this native-app can be reached ( by the mozeidon CLI )
 */
func GetProfileDirectory() (string, error) {
	configDir, err := os.UserConfigDir()
	if err != nil {
		return "", err
	}
	profilesDir := filepath.Join(configDir, "mozeidon_profiles")
	err = os.MkdirAll(profilesDir, 0755)
	return profilesDir, err
}

type NativeAppProfile struct {
	IpcName             string `json:"ipcName"`
	FileName            string `json:"fileName"`
	BrowserName         string `json:"browserName"`
	BrowserEngine       string `json:"browserEngine"`
	BrowserVersion      string `json:"browserVersion"`
	ProfileId           string `json:"profileId"`
	ProfileRank         int    `json:"profileRank"`
	ProfileName         string `json:"profileName"`
	ProfileAlias        string `json:"profileAlias"`
	ProfileCommandAlias string `json:"profileCommandAlias"`
	InstanceId          string `json:"instanceId"`
	UserAgent           string `json:"userAgent"`
	Pid                 int    `json:"pid"`
	RegisteredAt        string `json:"registeredAt"`
}

func GetNativeAppProfile(response *RegistrationInfoResponse) *NativeAppProfile {
	instanceId := uuid.New().String()
	pid := os.Getpid()

	return &NativeAppProfile{
		// e.g mozeidon_native_app_12345678_67891011
		IpcName: fmt.Sprintf("mozeidon_native_app_%d_%s", pid, response.Data.ProfileId[:8]),
		// e.g 12345_67891011.json
		FileName: fmt.Sprintf(
			"%d_%s.json",
			pid,
			response.Data.ProfileId[:8],
		),
		BrowserName:         response.Data.BrowserName,
		BrowserEngine:       response.Data.BrowserEngine,
		BrowserVersion:      response.Data.BrowserVersion,
		ProfileId:           response.Data.ProfileId,
		ProfileRank:         response.Data.ProfileRank,
		ProfileName:         response.Data.ProfileName,
		ProfileAlias:        response.Data.ProfileAlias,
		ProfileCommandAlias: response.Data.ProfileCommandAlias,
		InstanceId:          instanceId,
		UserAgent:           response.Data.UserAgent,
		Pid:                 pid,
		RegisteredAt:        response.Data.RegisteredAt,
	}
}

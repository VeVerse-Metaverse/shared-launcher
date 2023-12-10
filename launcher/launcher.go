package launcher

import (
	sm "dev.hackerman.me/artheon/veverse-shared/model"
)

type Launcher interface {
	// GetLauncherMetadata returns the launcher metadata from the cache if it exists, otherwise returns nil
	GetLauncherMetadata() *sm.LauncherV2
	// FetchLauncherMetadata fetches the launcher metadata from the API and stores it in the cache, returns an error if it fails
	FetchLauncherMetadata() (*sm.LauncherV2, error)
}

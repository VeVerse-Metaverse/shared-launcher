package tests

import sm "dev.hackerman.me/artheon/veverse-shared/model"

type testLauncher struct {
}

func (m *testLauncher) FetchLauncherMetadata() (*sm.LauncherV2, error) {
	panic("not implemented")
}

func (m *testLauncher) GetLauncherMetadata() *sm.LauncherV2 {
	return &sm.LauncherV2{
		Entity: sm.Entity{
			Identifier:  sm.Identifier{},
			Timestamps:  sm.Timestamps{},
			EntityType:  "",
			Public:      false,
			Views:       0,
			Owner:       nil,
			Accessibles: nil,
			Files:       nil,
			Links:       nil,
			Properties:  nil,
			Likables:    nil,
			Comments:    nil,
			Liked:       nil,
			Likes:       nil,
			Dislikes:    nil,
		},
		Name:        "",
		Description: "",
		Releases:    nil,
		Apps:        nil,
	}
}

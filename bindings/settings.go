package bindings

import (
	"context"
	"encoding/json"
	"os"
	"path/filepath"
	"time"

	"github.com/pkg/errors"
	"github.com/satisfactorymodding/SatisfactoryModManager/utils"
	"github.com/spf13/viper"
	wailsRuntime "github.com/wailsapp/wails/v2/pkg/runtime"
)

type SettingsData struct {
	FavouriteMods    []string `json:"favouriteMods"`
	AppHeight        int      `json:"appHeight"`
	ExpandedAppWidth int      `json:"expandedAppWidth"`
}

type Settings struct {
	ctx  context.Context
	Data SettingsData
}

var SettingsFileName = "settings.json"

func MakeSettings() (*Settings, error) {
	s := &Settings{}

	if err := s.load(); err != nil {
		return nil, errors.Wrap(err, "failed to load settings")
	}

	return s, nil
}

func (s *Settings) load() error {
	settingsFilePath := filepath.Join(viper.GetString("local-dir"), SettingsFileName)

	_, err := os.Stat(settingsFilePath)
	if err != nil {
		if !os.IsNotExist(err) {
			return errors.Wrap(err, "failed to stat settings file")
		}

		s.Data = SettingsData{
			FavouriteMods:    []string{},
			AppHeight:        utils.UnexpandedMinHeight,
			ExpandedAppWidth: utils.UnexpandedMinWidth,
		}
		err = s.save()
		if err != nil {
			return errors.Wrap(err, "failed to save default settings")
		}
	}

	settingsFile, err := os.ReadFile(filepath.Join(viper.GetString("local-dir"), SettingsFileName))
	if err != nil {
		return errors.Wrap(err, "failed to read settings")
	}

	if err := json.Unmarshal(settingsFile, &s.Data); err != nil {
		return errors.Wrap(err, "failed to unmarshal settings")
	}

	return nil
}

func (s *Settings) save() error {
	settingsFile, err := json.MarshalIndent(s.Data, "", "  ")
	if err != nil {
		return errors.Wrap(err, "failed to marshal settings")
	}
	err = os.WriteFile(filepath.Join(viper.GetString("local-dir"), SettingsFileName), settingsFile, 0755)
	if err != nil {
		return errors.Wrap(err, "failed to write settings")
	}

	return nil
}

func (s *Settings) startup(ctx context.Context) {
	s.ctx = ctx
	sizeTicker := time.NewTicker(100 * time.Millisecond)
	go func() {
		for range sizeTicker.C {
			w, h := wailsRuntime.WindowGetSize(s.ctx)
			if BindingsInstance.App.isExpanded {
				if w != s.Data.ExpandedAppWidth {
					s.Data.ExpandedAppWidth = w
					s.save()
				}
			}
			if h != s.Data.AppHeight {
				s.Data.AppHeight = h
				s.save()
			}
		}
	}()
}

func (s *Settings) FavouriteMod(modReference string) bool {
	idx := -1
	for i, mod := range s.Data.FavouriteMods {
		if mod == modReference {
			idx = i
			break
		}
	}
	if idx != -1 {
		return false
	}
	s.Data.FavouriteMods = append(s.Data.FavouriteMods, modReference)
	s.save()
	s.emitFavouriteMods()
	return true
}

func (s *Settings) UnFavouriteMod(modReference string) bool {
	idx := -1
	for i, mod := range s.Data.FavouriteMods {
		if mod == modReference {
			idx = i
			break
		}
	}
	if idx == -1 {
		return false
	}
	s.Data.FavouriteMods = append(s.Data.FavouriteMods[:idx], s.Data.FavouriteMods[idx+1:]...)
	s.save()
	s.emitFavouriteMods()
	return true
}

func (s *Settings) GetFavouriteMods() []string {
	return s.Data.FavouriteMods
}

func (s *Settings) emitFavouriteMods() {
	wailsRuntime.EventsEmit(s.ctx, "favouriteMods", s.Data.FavouriteMods)
}
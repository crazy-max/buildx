package localstate

import (
	"encoding/json"
	"os"
	"path/filepath"

	"github.com/pkg/errors"
)

func (ls *LocalState) Migrate() error {
	currentVersion := ls.readVersion()
	if currentVersion == version {
		return nil
	}
	migrations := map[int]func(*LocalState) error{
		2: (*LocalState).migration2,
	}
	for v := currentVersion + 1; v <= version; v++ {
		migration, found := migrations[v]
		if !found {
			return errors.Errorf("localstate migration v%d not found", v)
		}
		if err := migration(ls); err != nil {
			return errors.Wrapf(err, "localstate migration v%d failed", v)
		}
	}
	return ls.writeVersion(version)
}

func (ls *LocalState) migration2() error {
	return filepath.Walk(ls.GroupDir(), func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}
		dt, err := os.ReadFile(path)
		if err != nil {
			return err
		}
		var ndt map[string]interface{}
		if err := json.Unmarshal(dt, &ndt); err != nil {
			return err
		}
		delete(ndt, "Definition")
		delete(ndt, "Inputs")
		mdt, err := json.Marshal(ndt)
		if err != nil {
			return err
		}
		if err := os.WriteFile(path, mdt, 0600); err != nil {
			return err
		}
		return nil
	})
}

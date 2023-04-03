package localstate

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/docker/docker/pkg/ioutils"
	"github.com/gofrs/flock"
)

const refsDir = "refs"

type Ref struct {
	Cmd            string `json:"command"`
	LocalPath      string `json:"localPath"`
	DockerfilePath string `json:"dockerfilePath"`
}

type LocalState struct {
	root string
}

func New(root string) (*LocalState, error) {
	if root == "" {
		return nil, fmt.Errorf("root dir empty")
	}
	if err := os.MkdirAll(filepath.Join(root, refsDir), 0700); err != nil {
		return nil, err
	}
	return &LocalState{
		root: root,
	}, nil
}

type Txn struct {
	ls *LocalState
}

func (ls *LocalState) Txn() (*Txn, func(), error) {
	l := flock.New(filepath.Join(ls.root, "localstate.lock"))
	if err := l.Lock(); err != nil {
		return nil, nil, err
	}
	return &Txn{
			ls: ls,
		}, func() {
			l.Close()
		}, nil
}

func (t *Txn) GetRef(builderName, nodeName, id string) (*Ref, error) {
	if err := t.ls.validate(builderName, nodeName, id); err != nil {
		return nil, err
	}
	dt, err := os.ReadFile(filepath.Join(t.ls.root, refsDir, builderName, nodeName, id))
	if err != nil {
		return nil, err
	}
	var ref Ref
	if err := json.Unmarshal(dt, &ref); err != nil {
		return nil, err
	}
	return &ref, nil
}

func (t *Txn) SetRef(builderName, nodeName, id string, ref Ref) error {
	if err := t.ls.validate(builderName, nodeName, id); err != nil {
		return err
	}
	refDir := filepath.Join(t.ls.root, refsDir, builderName, nodeName)
	if err := os.MkdirAll(refDir, 0700); err != nil {
		return err
	}
	dt, err := json.Marshal(ref)
	if err != nil {
		return err
	}
	return ioutils.AtomicWriteFile(filepath.Join(refDir, id), dt, 0600)
}

func (t *Txn) Remove(builderName, nodeName, id string) error {
	if err := t.ls.validate(builderName, nodeName, id); err != nil {
		return err
	}
	return os.Remove(filepath.Join(t.ls.root, refsDir, builderName, nodeName, id))
}

func (t *Txn) RemoveForBuilder(builderName string) error {
	if builderName == "" {
		return fmt.Errorf("builder name empty")
	}
	return os.RemoveAll(filepath.Join(t.ls.root, refsDir, builderName))
}

func (t *Txn) RemoveForNode(builderName string, nodeName string) error {
	if builderName == "" {
		return fmt.Errorf("builder name empty")
	}
	if nodeName == "" {
		return fmt.Errorf("node name empty")
	}
	return os.RemoveAll(filepath.Join(t.ls.root, refsDir, builderName, nodeName))
}

func (ls *LocalState) validate(builderName, nodeName, id string) error {
	if builderName == "" {
		return fmt.Errorf("builder name empty")
	}
	if nodeName == "" {
		return fmt.Errorf("node name empty")
	}
	if id == "" {
		return fmt.Errorf("ref ID empty")
	}
	return nil
}

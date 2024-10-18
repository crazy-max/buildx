package confutil

import (
	"crypto/rand"
	"encoding/hex"
	"os"
	"path/filepath"
	"sync"
)

var nodeIdentifierMu sync.Mutex

func TryNodeIdentifier(cfg *Config) (out string) {
	nodeIdentifierMu.Lock()
	defer nodeIdentifierMu.Unlock()
	sessionFilename := ".buildNodeID"
	sessionFilepath := filepath.Join(cfg.Dir(), sessionFilename)
	if _, err := os.Lstat(sessionFilepath); err != nil {
		if os.IsNotExist(err) { // create a new file with stored randomness
			b := make([]byte, 8)
			if _, err := rand.Read(b); err != nil {
				return out
			}
			if err := cfg.WriteFile(sessionFilename, []byte(hex.EncodeToString(b)), 0600); err != nil {
				return out
			}
		}
	}

	dt, err := os.ReadFile(sessionFilepath)
	if err == nil {
		return string(dt)
	}
	return
}

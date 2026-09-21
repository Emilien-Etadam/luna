package crypto

import (
	"os"
	"testing"

	"luna-backend/config"
)

func TestSymmetricKeyFileExists(t *testing.T) {
	dir := t.TempDir()
	cfg := &config.CommonConfig{
		Env: &config.Environmental{DATA_PATH: dir},
	}

	exists, err := SymmetricKeyFileExists(cfg, "passwordPepper")
	if err != nil {
		t.Fatal(err)
	}
	if exists {
		t.Fatal("expected missing key file")
	}

	keysDir := cfg.Env.GetKeysPath()
	if err := os.MkdirAll(keysDir, 0o770); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(keysDir+"/passwordPepper.key", []byte("x"), 0o600); err != nil {
		t.Fatal(err)
	}

	exists, err = SymmetricKeyFileExists(cfg, "passwordPepper")
	if err != nil {
		t.Fatal(err)
	}
	if !exists {
		t.Fatal("expected key file to exist")
	}
}

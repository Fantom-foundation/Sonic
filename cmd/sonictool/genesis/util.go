package genesis

import (
	"os"
	"path"
)

func setGenesisProcessing(chaindataDir string) {
	f, _ := os.Create(path.Join(chaindataDir, "unfinished"))
	if f != nil {
		_ = f.Close()
	}
}

func setGenesisComplete(chaindataDir string) {
	_ = os.Remove(path.Join(chaindataDir, "unfinished"))
}

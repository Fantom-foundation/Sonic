package integration

import (
	"path"

	"github.com/Fantom-foundation/go-opera/utils"
)

func isInterrupted(chaindataDir string) bool {
	return utils.FileExists(path.Join(chaindataDir, "unfinished"))
}

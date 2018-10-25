package destroy

import (
	"encoding/json"
	"io/ioutil"
	"path/filepath"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"

	"github.com/openshift/installer/pkg/asset/cluster"
	"github.com/openshift/installer/pkg/types"
)

// Destroyer allows multiple implementations of destroy
// for different platforms.
type Destroyer interface {
	Run() error
}

// NewFunc is an interface for creating platform-specific destroyers.
type NewFunc func(logger logrus.FieldLogger, metadata *types.ClusterMetadata) (Destroyer, error)

// Registry maps ClusterMetadata.Platform() to per-platform Destroyer creators.
var Registry = make(map[string]NewFunc)

// New returns a Destroyer based on `metadata.json` in `rootDir`.
func New(logger logrus.FieldLogger, rootDir string) (Destroyer, error) {
	path := filepath.Join(rootDir, cluster.MetadataFilename)
	raw, err := ioutil.ReadFile(filepath.Join(rootDir, cluster.MetadataFilename))
	if err != nil {
		return nil, errors.Wrapf(err, "failed to read %s file", cluster.MetadataFilename)
	}

	var cmetadata *types.ClusterMetadata
	if err := json.Unmarshal(raw, &cmetadata); err != nil {
		return nil, errors.Wrapf(err, "failed to Unmarshal data from %s file to types.ClusterMetadata", cluster.MetadataFilename)
	}

	platform := cmetadata.Platform()
	if platform == "" {
		return nil, errors.Errorf("no platform configured in %q", path)
	}

	creator, ok := Registry[platform]
	if !ok {
		return nil, errors.Errorf("no destroyers registered for %q", platform)
	}
	return creator(logger, cmetadata)
}

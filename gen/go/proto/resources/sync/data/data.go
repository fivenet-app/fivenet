package syncdata

import (
	"fmt"

	"github.com/cespare/xxhash/v2"
	"github.com/fivenet-app/fivenet/v2026/pkg/utils/protoutils"
)

func (x *DataUser) GetJSONAndHash() ([]byte, uint64, error) {
	out, err := protoutils.MarshalToJSON(x)
	if err != nil {
		return nil, 0, fmt.Errorf(
			"failed to marshal user data to JSON for user %d (%s). %w",
			x.GetUserId(),
			x.GetIdentifier(),
			err,
		)
	}

	hash := xxhash.Sum64(out)

	return out, hash, nil
}

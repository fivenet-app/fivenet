package centrum

import (
	"database/sql/driver"

	users "github.com/galexrt/fivenet/gen/go/proto/resources/users"
	"github.com/galexrt/fivenet/pkg/utils/protoutils"
	"google.golang.org/protobuf/encoding/protojson"
)

func (x *Disponents) Merge(in *Disponents) *Disponents {
	if len(in.Disponents) == 0 {
		x.Disponents = []*users.UserShort{}
	} else {
		x.Disponents = in.Disponents
	}

	return x
}

func (x *UserUnitMapping) Merge(in *UserUnitMapping) *UserUnitMapping {
	if x.UnitId != in.UnitId {
		x.UnitId = in.UnitId
	}

	if x.UserId != in.UserId {
		x.UserId = in.UserId
	}

	return x
}

func (x *PredefinedStatus) Scan(value any) error {
	switch t := value.(type) {
	case string:
		return protojson.Unmarshal([]byte(t), x)
	case []byte:
		return protojson.Unmarshal(t, x)
	}
	return nil
}

// Scan implements driver.Valuer for protobuf PredefinedStatus.
func (x *PredefinedStatus) Value() (driver.Value, error) {
	if x == nil {
		return nil, nil
	}

	out, err := protoutils.Marshal(x)
	return string(out), err
}

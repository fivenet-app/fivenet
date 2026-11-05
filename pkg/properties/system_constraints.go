package properties

import (
	"fmt"

	propertypb "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/properties"
	"google.golang.org/protobuf/encoding/protojson"
)

// MustParseValueConstraintsJSON decodes code-generated system-definition
// constraints. It panics for invalid generator input, which is a startup-time
// programming error rather than a recoverable instance configuration error.
func MustParseValueConstraintsJSON(value string) *propertypb.PropertyValueConstraints {
	constraints := &propertypb.PropertyValueConstraints{}
	if err := protojson.Unmarshal([]byte(value), constraints); err != nil {
		panic(fmt.Sprintf("invalid generated property value constraints: %v", err))
	}
	return constraints
}

package idutil_test

import (
	"testing"

	"github.com/maqdev/go-be-template/util/idutil"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

func Test_New(t *testing.T) {
	req := require.New(t)
	id := idutil.NewID()
	req.Len(id, 22)

	decoded := idutil.MustDecodeShortUUID(id)
	req.NotEqual(uuid.Nil, decoded)

	encoded := idutil.ShortenUUID(decoded)
	req.Equal(id, encoded)
}

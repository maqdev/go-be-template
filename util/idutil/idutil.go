package idutil

import (
	"encoding/base64"

	"github.com/google/uuid"
)

/*
	IDs generated are based on UUIDs that aren't time-based and encoded with URL safe base64 alphabet,
	We don't use canonical UUID format with dashes to preserve space as those are mostly stored and passed as strings.
	google's uuid.New generates cryptographically random UUID
*/

func NewID() string {
	return ShortenUUID(uuid.New())
}

func ShortenUUID(u uuid.UUID) string {
	return base64.RawURLEncoding.EncodeToString(u[:])
}

func DecodeShortUUID(short string) (uuid.UUID, error) {
	bytes, err := base64.RawURLEncoding.DecodeString(short)
	if err != nil {
		return uuid.Nil, err
	}
	return uuid.FromBytes(bytes)
}

func MustDecodeShortUUID(short string) uuid.UUID {
	r, err := DecodeShortUUID(short)
	if err != nil {
		panic(err)
	}
	return r
}

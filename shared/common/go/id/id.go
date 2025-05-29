package id

/*
Package for generating unique identifiers (IDs) based on a base64 representation of the bytes of a UUID.
The purpose of this package is to create short (22 characters), URL-safe IDs that can be used in various applications.
Also includes functions for converting these IDs to and from classic UUIDs.
*/

import (
	"encoding/base64"
	"errors"
	"github.com/google/uuid"
	"strconv"
)

func New() (string, error) {
	newUUID, err := uuid.NewRandom()
	if err != nil {
		return "", err
	}

	idString, err := FromUUID(newUUID)
	if err != nil {
		return "", err
	}

	return idString, nil
}

func FromUUID(uuid uuid.UUID) (string, error) {
	bytes := uuid[:]
	idString := base64.RawURLEncoding.EncodeToString(bytes)

	return idString, nil
}

func ToUUID(idString string) (uuid.UUID, error) {
	bytes, err := base64.RawURLEncoding.DecodeString(idString)
	if err != nil {
		return uuid.UUID{}, err
	}

	if len(bytes) != 16 {
		return uuid.UUID{}, errors.New("invalid ID length: expected 16 bytes, got " + strconv.Itoa(len(bytes)))
	}

	return uuid.FromBytes(bytes)
}

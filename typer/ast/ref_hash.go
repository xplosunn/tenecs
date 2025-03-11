package ast

import (
	"crypto/sha1"
	"encoding/base64"
	"errors"
	"github.com/fsamin/go-dump"
)

type RefHashes map[Ref]string

func DetermineRefHashes(program Program) (RefHashes, error) {
	result := RefHashes{}
	for ref, exp := range program.Declarations {
		_, ok := result[ref]
		if ok {
			return result, errors.New("tried rehashing same ref")
		}
		hashed, err := hash(exp)
		if err != nil {
			return result, err
		}
		result[ref] = hashed
	}
	for ref, varType := range program.StructFunctions {
		_, ok := result[ref]
		if ok {
			return result, errors.New("tried rehashing same ref")
		}
		hashed, err := hash(varType)
		if err != nil {
			return result, err
		}
		result[ref] = hashed
	}
	return result, nil
}

func hash(thingToHash any) (string, error) {
	str, err := dump.Sdump(thingToHash)
	if err != nil {
		return "", err
	}
	hasher := sha1.New()
	_, err = hasher.Write([]byte(str))
	if err != nil {
		return "", err
	}
	sha := base64.URLEncoding.EncodeToString(hasher.Sum(nil))
	return sha, nil
}

package checkfile

import (
	"fmt"

	"github.com/HardDie/ghashdeep/internal/entities/hash"
)

type Object struct {
	Name string
	Hash hash.Hash
}

func NewObjectFromString(line string, hashLen int) (Object, error) {
	hashString := line[:hashLen]
	name := line[hashLen+2:]
	h, err := hash.FromString(hashString)
	if err != nil {
		return Object{}, fmt.Errorf("hash.FromString: %w", err)
	}
	return Object{
		Name: name,
		Hash: h,
	}, nil
}

func (o Object) String() string {
	return o.Hash.String() + "  " + o.Name
}

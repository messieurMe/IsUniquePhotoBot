package util

import (
	"bytes"
	"image"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"math/bits"

	goimagehash "github.com/corona10/goimagehash"
)

type HashHelper struct {
	similarityThreshold int
}

func NewHashHelper(similarityThreshold int) *HashHelper {
	return &HashHelper{
		similarityThreshold: similarityThreshold,
	}
}

func (_ *HashHelper) computePHash(data []byte) (uint64, error) {
	img, _, err := image.Decode(bytes.NewReader(data))
	if err != nil {
		return 0, err
	}
	h, err := goimagehash.PerceptionHash(img)
	if err != nil {
		return 0, err
	}
	return h.GetHash(), nil
}

func (this *HashHelper) AreSimilar(a, b uint64) bool {
	return bits.OnesCount64(a^b) < this.similarityThreshold
}

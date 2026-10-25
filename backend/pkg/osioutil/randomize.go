package osioutil

import (
	"math/rand/v2"
	"time"
)

type RandomizeStruct struct {
	randSource *rand.Rand
}

func NewRandomize() RandomizeStruct {
	time.Sleep(1 * time.Millisecond)
	pvSeed1 := uint64(time.Now().Add(-25 * 365 * 24 * time.Hour).UnixNano())
	pvSeed2 := uint64(time.Now().Add(22 * 365 * 24 * time.Hour).UnixNano())

	return RandomizeStruct{
		randSource: rand.New(rand.NewPCG(pvSeed1, pvSeed2)),
	}
}

func (parRand *RandomizeStruct) GerRandomInt(parMin, parMax int) int {
	if parMax-parMin < 0 {
		return 0
	}

	return parRand.randSource.IntN(parMax-parMin) + parMin
}

func (parRand *RandomizeStruct) GerRandomByte(parMin, parMax byte) byte {

	pvMin := uint(parMin)
	pvMax := uint(parMax) + 1

	if int(pvMax)-int(pvMin) < 0 {
		return 0
	}

	return byte(parRand.randSource.UintN(pvMax-pvMin) + pvMin)
}

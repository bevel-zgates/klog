package record

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestSumUpTimes(t *testing.T) {
	r := NewRecord(Ɀ_Date_(2020, 1, 1))
	r.AddDuration(NewDuration(1, 0))
	r.AddDuration(NewDuration(2, 0))
	assert.Equal(t, NewDuration(3, 0), Total(r))
}

func TestSumUpZeroIfNoTimesAvailable(t *testing.T) {
	r := NewRecord(Ɀ_Date_(2020, 1, 1))
	assert.Equal(t, NewDuration(0, 0), Total(r))
}

func TestSumUpRanges(t *testing.T) {
	range1 := Ɀ_Range_(Ɀ_Time_(9, 7), Ɀ_Time_(12, 59))
	range2 := Ɀ_Range_(Ɀ_Time_(13, 49), Ɀ_Time_(17, 12))
	r := NewRecord(Ɀ_Date_(2020, 1, 1))
	r.AddRange(range1)
	r.AddRange(range2)
	assert.Equal(t, NewDuration(7, 15), Total(r))
}

func TestSumUpTimesAndRanges(t *testing.T) {
	range1 := Ɀ_Range_(Ɀ_Time_(8, 0), Ɀ_Time_(12, 0))
	r := NewRecord(Ɀ_Date_(2020, 1, 1))
	r.AddDuration(NewDuration(1, 33))
	r.AddRange(range1)
	assert.Equal(t, NewDuration(5, 33), Total(r))
}

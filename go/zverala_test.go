package main

import (
	"fmt"
	"math/rand"
	"os"
	"sort"
	"testing"
	"time"

	ktime "zverala/klvanistic_time"

	"github.com/stretchr/testify/require"
)

func TestComposeKyear(t *testing.T) {
	for _, item := range []struct {
		// Input
		start time.Time
		end   time.Time
		// Expected results
		DoubleYear       int
		DoubleYearDigits []int
		Direction        ktime.DirType
		Length           int
	}{
		{
			start: time.Date(2013, 12, 21, 11, 0, 0, 0, time.UTC),
			end:   time.Date(2014, 12, 21, 11, 0, 0, 0, time.UTC),

			DoubleYear:       27073,
			DoubleYearDigits: []int{2, 7, 0, 7, 3},
			Direction:        ktime.OUT,
			Length:           365,
		},
		{
			start: time.Date(2014, 12, 21, 11, 0, 0, 0, time.UTC),
			end:   time.Date(2015, 12, 22, 11, 0, 0, 0, time.UTC),

			DoubleYear:       27073,
			DoubleYearDigits: []int{2, 7, 0, 7, 3},
			Direction:        ktime.IN,
			Length:           366,
		},
		{
			start: time.Date(2015, 12, 22, 11, 0, 0, 0, time.UTC),
			end:   time.Date(2016, 12, 21, 11, 0, 0, 0, time.UTC),

			DoubleYear:       27074,
			DoubleYearDigits: []int{2, 7, 0, 7, 4},
			Direction:        ktime.OUT,
			Length:           365,
		},
	} {
		t.Run(fmt.Sprintf("year_%d", item.start.Year()), func(t *testing.T) {
			kyear := ktime.ComputeKyear(item.start, item.end)
			require.Equal(t, item.DoubleYear, kyear.Doubleyear)
			require.Equal(t, item.Direction, kyear.Direction)
			require.Equal(t, item.Length, kyear.Length)
			require.True(t, item.start.Equal(kyear.NormalYearStart))
			require.Len(t, kyear.DoubleyearDigits, len(item.DoubleYearDigits))
			for i, dig := range kyear.DoubleyearDigits {
				require.Equal(t, item.DoubleYearDigits[i], dig)
			}
		})
	}
}

func TestGetOrderedSteps(t *testing.T) {
	var sins, cosins []float64
	for i := 0; i < NUM_CREATURES; i++ {
		sins = append(sins, rand.Float64())
		cosins = append(cosins, rand.Float64())
	}

	sort.Slice(sins, func(i, j int) bool { return sins[i] < sins[j] })
	sort.Slice(cosins, func(i, j int) bool { return cosins[i] < cosins[j] })

	results1 := getOrderedSteps(sins, cosins)

	for i := range sins {
		sins[i] += 10
	}
	for i := range cosins {
		cosins[i] += 10
	}

	results2 := getOrderedSteps(sins, cosins)

	for i, res := range results1 {
		if i == 0 {
			continue
		}
		// Because of how floats work, we can't compare these numbers directly.
		// Best we can test is whether they are very similar.
		require.Less(t, res-results2[i], 0.0005)
	}
}

func TestOrderedCreatures(t *testing.T) {
	var days []int
	for i := 0; i < NUM_CREATURES-1; i++ {
		// Reserve 0 days for chimera
		days = append(days, i+1)
	}

	t.Run("ktime.OUT", func(t *testing.T) {
		creatures := getCreaturesInOrder(ktime.OUT, 0, days)
		for i, creature := range creatures {
			expectedName := Creatures[i]
			require.Equal(t, expectedName, creature.name)
			require.Equal(t, i, creature.days)
		}

		t.Run("dragon_days", func(t *testing.T) {
			creatures = addDragonDays(true, ktime.OUT, creatures)
			for i, creatureIndex := range DragonsAfterCreatureIndex {
				creature := creatures[creatureIndex+1+i]
				require.Equal(t, Dragons[i].name, creature.name)
				require.Equal(t, 1, Dragons[i].days)
			}
		})

	})

	t.Run("ktime.IN", func(t *testing.T) {
		creatures := getCreaturesInOrder(ktime.IN, 0, days)
		for i, creature := range creatures {
			expectName := Creatures[NUM_CREATURES-1-i]
			require.Equal(t, expectName, creature.name)
			expectDays := i + 1
			if i == NUM_CREATURES-1 {
				expectDays = 0 // Chimera has 0 days
			}
			require.Equal(t, expectDays, creature.days)
		}

		t.Run("dragon_days", func(t *testing.T) {
			creatures = addDragonDays(true, ktime.IN, creatures)
			for i, creatureIndex := range DragonsAfterCreatureIndex {
				// I appreciate this is quite difficult to reason about.
				// Basically we're doing this:
				//
				//  - Get the total Length of the new creatures array ->
				//    NUM_CREATURES + NUM_DRAGONS
				//
				//  - For each index in the DragonsAfterCreatureIndex, count
				//    backwards in the creatures array. Arrays' indeces only
				//    reach len(array)-1 -> -1
				//
				//  - The dragon indeces are actually indeces of creatures
				//    followed by dragons. To get dragons -> another -1
				//
				//  - Each index must take into account the dragons previously
				//    considered -> -i
				creature := creatures[NUM_CREATURES+NUM_DRAGONS-i-creatureIndex-2]
				require.Equal(t, Dragons[i].name, creature.name)
				require.Equal(t, 1, Dragons[i].days)
			}
		})
	})
}

func TestDragonYear(t *testing.T) {
	for _, item := range []struct {
		kYear        kYear
		isDragonYear bool
	}{
		{
			kYear: ktime.ComputeKyear(time.Date(2013, 12, 21, 0, 0, 0, 0, time.UTC),
				time.Date(2014, 12, 21, 0, 0, 0, 0, time.UTC)),
			isDragonYear: false,
		},
		{
			kYear: ktime.ComputeKyear(time.Date(2012, 12, 21, 0, 0, 0, 0, time.UTC),
				time.Date(2013, 12, 21, 0, 0, 0, 0, time.UTC)),
			isDragonYear: true,
		},
		{
			kYear: ktime.ComputeKyear(time.Date(2048, 12, 22, 0, 0, 0, 0, time.UTC),
				time.Date(2049, 12, 21, 0, 0, 0, 0, time.UTC)),
			isDragonYear: true,
		},
	} {
		t.Run(fmt.Sprintf("checking_kyear_%s", item.kYear.ToReadableString()), func(t *testing.T) {
			is := ktime.IsDragonYear(item.kYear)
			require.Equal(t, item.isDragonYear, is)
		})
	}
}

func TestFileManipulation(t *testing.T) {
	year := 2019
	sol1In := time.Date(year, 12, 21, 0, 0, 0, 0, time.Local)
	sol2In := time.Date(year, 12, 22, 0, 0, 0, 0, time.Local)
	sol3In := time.Date(year, 12, 21, 0, 0, 0, 0, time.Local)

	tmpFile, err := os.CreateTemp("/tmp/", "zverala_test")
	require.NoError(t, err)
	file = tmpFile.Name()
	defer os.RemoveAll("/tmp/" + tmpFile.Name())

	dYear := ktime.DoubleYear{
		OutKyear: ktime.ComputeKyear(sol1In, sol2In),
		InKyear:  ktime.ComputeKyear(sol2In, sol3In),
		EndTime:  sol3In,
	}
	writeYearToFile(dYear)

	t.Run(fmt.Sprintf("get_cached_year_%d", year), func(t *testing.T) {
		sol1, sol2, sol3, found := cachedYear(year)

		if !found {
			var content []byte
			tmpFile.Read(content)
			fmt.Printf("Content of file:\n%s", string(content))
		}

		require.True(t, found)
		require.True(t, sol1.Day() == sol1In.Day())
		require.True(t, sol2.Day() == sol2In.Day())
		require.True(t, sol3.Day() == sol3In.Day())
	})
}

func TestCalculateABC(t *testing.T) {
	t.Run("calculate_a", func(t *testing.T) {
		for _, item := range []struct {
			input  int
			output int
		}{
			{0, 1},
			{1, 2},
			{2, 3},
			{8, 9},
			{9, 1},
			{10, 2},
			{12, 4},
		} {
			t.Run(fmt.Sprintf("input_is_%d", item.input), func(t *testing.T) {
				out := calculate_a(item.input)
				require.Equal(t, item.output, out)
			})
		}
	})

	t.Run("calculate_b", func(t *testing.T) {
		for _, item := range []struct {
			digits     []int
			outwardOut int
			inwardOut  int
		}{
			{
				digits:     []int{2, 2, 2, 2},
				outwardOut: 8,
				inwardOut:  5,
			},
			{
				digits:     []int{9, 0, 1, 9},
				outwardOut: 4,
				inwardOut:  4,
			},
		} {
			t.Run(fmt.Sprintf("input_year_is_%v", item.digits), func(t *testing.T) {
				out := calculate_b(item.digits, true)
				require.Equal(t, item.outwardOut, out)
				out = calculate_b(item.digits, false)
				require.Equal(t, item.inwardOut, out)
			})
		}
	})

	t.Run("calculate_c", func(t *testing.T) {
		for _, item := range []struct {
			DoubleYear int
			digits     []int
			output     int
		}{
			{
				DoubleYear: 27076,
				digits:     []int{2, 7, 0, 7, 6},
				output:     7,
			},
		} {
			t.Run(fmt.Sprintf("year_%d", item.DoubleYear), func(t *testing.T) {
				out := calculate_c(item.DoubleYear, item.digits)
				require.Equal(t, item.output, out)
			})
		}
	})
}

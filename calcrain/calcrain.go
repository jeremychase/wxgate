package calcrain

import (
	"errors"
	"fmt"
	"sort"
	"time"
)

type record struct {
	Time   time.Time
	Amount float64
}

type Data struct {
	Rain    []record
	PrevMax record
}

// Append amount to in-memory dataset.
func (d *Data) Append(amount float64, t time.Time) {
	r := record{
		Amount: amount,
		Time:   t,
	}

	// Keep track of when the day changes and update PrevMax.
	if len(d.Rain) > 1 {
		if d.Rain[len(d.Rain)-1].Time.Local().Weekday() != r.Time.Local().Weekday() {
			d.PrevMax = d.Rain[len(d.Rain)-1]
		}
	}

	d.Rain = append(d.Rain, r)
}

// RainLast24Hours returns rainfall over the trailing 24 hours based on
// recorded data. This prunes old data.
func (d *Data) RainLast24Hours(amount float64, t time.Time, threshold uint) (float64, bool, error) {
	prev, pruned, err := d.prevNow(t, threshold)
	if err != nil {
		return 0.0, pruned, err
	}

	return d.PrevMax.Amount - prev + amount, pruned, err
}

// prevNow returns cumulative rain for day previous to `t`. `threshold`
// represents the maximum age allowed.
func (d *Data) prevNow(t time.Time, threshold uint) (amount float64, pruned bool, err error) {
	// x is the first data point less than 24 hours from t
	x := sort.Search(len(d.Rain), func(i int) bool { return t.Sub(d.Rain[i].Time) < time.Hour*24 })

	if x < len(d.Rain) && x > 0 {
		if t.Sub(d.Rain[x-1].Time) < time.Hour*24+time.Minute*time.Duration(threshold) {
			amount = d.Rain[x-1].Amount
		} else {
			err = errors.New("threshold exceeded")
		}
	} else {
		err = errors.New("insufficient data")
	}

	// prune old data
	if x > len(d.Rain)/4 {
		fmt.Printf("Pruning at:\t%v\t%v\n", x, len(d.Rain))
		pruned = true
		d.Rain = append([]record(nil), d.Rain[x-x/10:]...) // keep some for future calculations
	}

	return
}

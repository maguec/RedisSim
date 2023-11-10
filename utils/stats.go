package utils

import (
	"fmt"

	"github.com/jamiealquiza/tachymeter"
	"github.com/maguec/metermaid"
)

func ShowStats(tach *tachymeter.Tachymeter, mm *metermaid.Metermaid) {

	results := tach.Calc()
	fmt.Println("------------------ Latency ------------------")
	fmt.Printf(
		"Max:\t\t%s\nMin:\t\t%s\nP95:\t\t%s\nP99:\t\t%s\nP99.9:\t\t%s\n\n",
		results.Time.Max,
		results.Time.Min,
		results.Time.P95,
		results.Time.P99,
		results.Time.P999,
	)
	fmt.Println("-------------- Latency Histogram ------------")
	fmt.Println("")
	fmt.Println(results.Histogram.String(10))
	rates := mm.Calc()
	fmt.Println("-------------------- Rate -------------------")
	fmt.Printf(
		"MaxRate:\t%.1f/s\nMinRate:\t%.1f/s\nP95Rate:\t%.1f/s\nP99Rate:\t%.1f/s\nP99.9Rate:\t%.1f/s\n",
		rates.MaxRate, rates.MinRate, rates.P95Rate, rates.P99Rate, rates.P999Rate,
	)

}

// Package float helps to interact with float type data
package float

import "math"

// Return round float into specific precision
// input 12.4312 will return 12.43
func Round[F float64 | float32](val F, precision uint) F {
	ratio := math.Pow(10, float64(precision))
	result := math.Round(float64(val)*ratio) / ratio

	return F(result)
}

// Package poisson provides functions for the Poisson probability distribution.
// The Poisson distribution models the number of independent events in a fixed
// interval at a known average rate — in football, this is goals per match.
package poisson

import "math"

// PMF (Probability Mass Function) returns the probability of exactly k events
// given expected average rate lambda.
// Formula: P(X=k) = (e^-λ × λ^k) / k!
func PMF(lambda float64, k int) float64 {
	if lambda <= 0 || k < 0 {
		return 0
	}
	return math.Exp(-lambda) * math.Pow(lambda, float64(k)) / factorial(k)
}

// factorial computes n! iteratively (not recursively) for performance and safety.
func factorial(n int) float64 {
	if n == 0 {
		return 1.0
	}
	result := 1.0
	for i := 2; i <= n; i++ {
		result *= float64(i)
	}
	return result
}

// SimulateGoals maps a uniform random value in [0,1) to a Poisson-distributed
// integer using the inverse CDF (quantile function) method.
// It partitions [0,1) into segments proportional to each outcome's probability,
// then returns whichever segment randVal falls into.
func SimulateGoals(lambda, randVal float64) int {
	cumulative := 0.0
	for k := 0; k <= 20; k++ {
		cumulative += PMF(lambda, k)
		if randVal < cumulative {
			return k
		}
	}
	return 20 // hard cap; unreachable under any realistic lambda
}
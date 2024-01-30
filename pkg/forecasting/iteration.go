package forecasting

// create a 'forward projection' iteration which connects to
// an online learning iteration, reads in its parameters and
// projects forward in time starting from the beginning of the
// historic window, continuing through the present point in time
// and ending up some point into the future

// iterations planned to be used with this iterator:
// - probabilistic reweighting iteration (evolves the mean and
//	covariance states forward in time)
// - any general stochadex simulation (combines with simulation
//  inference)

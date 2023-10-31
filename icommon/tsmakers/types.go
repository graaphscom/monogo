package tsmakers

type Maker func(src, dest string, resultCh chan<- makeResult)

type makeResult struct {
	success *makeDetails
	err     *makeError
}

type makeError struct {
	details makeDetails
	err     error
}

type makeDetails struct {
	divisionPath string
	iconName     string
}

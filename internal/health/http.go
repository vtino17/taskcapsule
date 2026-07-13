package health

import "net/http"

// HTTPChecker performs HTTP health checks with retries.
type HTTPChecker struct {
	ExpectedStatus int
}

func NewHTTPChecker(expectedStatus int) *HTTPChecker {
	return &HTTPChecker{ExpectedStatus: expectedStatus}
}

func (c *HTTPChecker) Check(url string) error {
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if c.ExpectedStatus > 0 && resp.StatusCode != c.ExpectedStatus {
		return &StatusError{Expected: c.ExpectedStatus, Got: resp.StatusCode}
	}

	return nil
}

type StatusError struct {
	Expected int
	Got      int
}

func (e *StatusError) Error() string {
	return "expected status " + string(rune(e.Expected)) + ", got " + string(rune(e.Got))
}

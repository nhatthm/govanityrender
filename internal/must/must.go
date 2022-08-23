package must

import "fmt"

// NoError panics only if the error is not nil.
func NoError(err error) {
	if err != nil {
		panic(fmt.Errorf("unexpected error: %w", err))
	}
}

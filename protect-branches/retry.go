/*
   Utility function to retry.
*/
package main

import (
	"fmt"
	"time"
)

func retry(attempts int, f func() error) error {
	var err error

	for i := 0; i < attempts; i++ {
		err = f()
		if err == nil {
			return nil
		}

		time.Sleep(30 * time.Millisecond)
	}
	return fmt.Errorf("after %d attempts, last error: %s", attempts, err)
}

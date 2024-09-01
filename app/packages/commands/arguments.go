package commands

import (
	"fmt"
	"strconv"
	"time"
)

func setExpiry(value string, key string) error {
	ex, err := strconv.Atoi(value)
	if err != nil {
		return fmt.Errorf("invalid value %s for expiry", value)
	}

	go func() {
		time.Sleep(time.Duration(ex) * time.Millisecond)
		delete(dRedis, key)
		fmt.Printf("deleted %s\n", key)
	}()

	return nil
}

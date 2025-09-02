package main

import (
	"fmt"
	"time"

	"collections/httpx/hooks"
)

func main() {
	bkfj := hooks.NewBackoffWithJitter(2*time.Second, 10*time.Minute, hooks.WithoutJitter)
	for attempt := range 30 {
		fmt.Println(bkfj.NextWaitDuration(nil, nil, attempt+1))
	}
}

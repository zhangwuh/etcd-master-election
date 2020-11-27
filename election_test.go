package etcd_master_election

import (
	"context"
	"fmt"
	"testing"
	"time"
)

//check the leadership for every 5 seconds, then master will resign the leadership and having the next master selection when
func TestCampaign(t *testing.T) {
	path := "/foo"
	done := make(chan bool)
	go func() {
		for {
			campaign(path, "s1")
		}
	}()

	go func() {
		for {
			campaign(path, "s2")
		}
	}()
	<-done
}

func campaign(path string, value string) {
	fmt.Println(fmt.Sprintf("%s start to compaign", value))
	ctx, cancel := context.WithCancel(context.Background())
	ticker := time.NewTicker(5 * time.Second)
	ch, _ := AcquireLeadership(ctx, []string{"127.0.0.1:2379"}, path, "s1", 10)
	var isMaster = false
	for {
		select {
		case <-ticker.C:
			if isMaster {
				fmt.Println(fmt.Sprintf("%s resigned his leadership", value))
				cancel()
				return
			}
		case v := <-ch:
			isMaster = v
			fmt.Println(fmt.Sprintf("status changed: %s is master:%v", value, v))
		}
	}
}

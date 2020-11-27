package etcd_master_election

import (
	"context"
	"fmt"

	v3 "github.com/coreos/etcd/clientv3"
	"go.etcd.io/etcd/clientv3/concurrency"
)

func AcquireLeadership(ctx context.Context,etcdNodes []string, electionPath string, val string, ttl int) (<-chan bool, error) {
	ch := make(chan bool)
	cli, err := v3.New(v3.Config{Endpoints: etcdNodes})
	if err != nil {
		return nil, err
	}
	go func() {
		defer cli.Close()
		defer close(ch)
		for {
			session, err := concurrency.NewSession(cli, concurrency.WithTTL(ttl))
			if err != nil {
				continue
			}
			e := concurrency.NewElection(session, electionPath)
			if err = e.Campaign(ctx, val); err != nil {
				fmt.Println(fmt.Sprintf("err when campain:%s", err.Error()))
				continue
			}
			ch <- true

			select {
			case <-session.Done():
				ch <- false
			case <-ctx.Done():
				return
			}
		}
	}()
	return ch, nil
}

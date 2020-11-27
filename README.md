# etcd-master-election
Implement Master election with etcd in a master/slave application cluster, if the master is unavailable, new master will be selected from other standby candidates.

## Usage:
	ctx, cancel := context.WithCancel(context.Background())
	leaseTtl := 10 //seconds
	ch, _ := AcquireLeadership(ctx, []string{"127.0.0.1:2379"}, "/election-path", "app", leaseTtl)
	for {
		select {
		case isMaster := <- ch:
			if isMaster:
				fmt.Println("i'm the master now")
			//do master's work
		}
	}
	

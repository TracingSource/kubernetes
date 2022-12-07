// Package namespace is a clientv3 wrapper that translates all keys to begin
// with a given prefix.
//
// First, create a client:
//
//	cli, err := clientv3.New(clientv3.Config{Endpoints: []string{"localhost:2379"}})
//	if err != nil {
//		// handle error!
//	}
//
// Next, override the client interfaces:
//
//	unprefixedKV := cli.KV
//	cli.KV = namespace.NewKV(cli.KV, "my-prefix/")
//	cli.Watcher = namespace.NewWatcher(cli.Watcher, "my-prefix/")
//	cli.Lease = namespace.NewLease(cli.Lease, "my-prefix/")
//
// Now calls using 'cli' will namespace / prefix all keys with "my-prefix/":
//
//	cli.Put(context.TODO(), "abc", "123")
//	resp, _ := unprefixedKV.Get(context.TODO(), "my-prefix/abc")
//	fmt.Printf("%s\n", resp.Kvs[0].Value)
//	// Output: 123
//	unprefixedKV.Put(context.TODO(), "my-prefix/abc", "456")
//	resp, _ = cli.Get(context.TODO(), "abc")
//	fmt.Printf("%s\n", resp.Kvs[0].Value)
//	// Output: 456
//
package namespace

# Network Status

`network-status` is a simple tool to check the status of your network connection. Originally
it was meant to be used on Unify Security Gateway router (hence support of `mips` architecture),
but practically can be used on any host that runs Linux (or WSL2 for that matter).

## Usage

You can of course build the binary using `go build` or `go run` commands, but there is a helpful
script aptly named `build.sh` that will build binaries for multiple architecture for you, and apply
`upx` compression to the binary. You can run it like this:

```bash

# Let's build the binary first (or you can go run it directly)
$ ./build.sh

# Now we can run the binary. It requires `sudo` because of direct socket access
$ sudo ./bin/network-status-amd64 -count 10 -max-latency 100ms google.com yahoo.com 1.1.1.1

# If the exit code is 0, then everything is fine
$ echo $?
```

You can now `scp` the `mips` binary to your router and run it there. Please check the 
`network-test-script.sh` for example of how to run it in the environment with multiple interfaces.

In my configuration, I have two interfaces: `eth0` and `pppoe0` which is used as a failover. I copied 
the `mips64` binary and `network-test-script.sh` to `/config/scripts` directory on the router. Then
I configured the loadbalancing as follows (you will need to `sudo su` before doing any of this):

```bash
# Let's check the script can even run
$ chmod +x /config/scripts/network-test-script.sh
$ /config/scripts/network-test-script.sh
...
...
$ echo $?
0
```

We see that script runs and produces desired output. To be absolutely sure, modify the required
latency to check that script will indeed fail (say, set it to 5ms in the script header). Now let's 
configure load balancing:

```bash
$ configure
$ set load-balance group wan_failover interface eth0 route-test \
  type script /config/scripts/network-test.sh
$ commit; save
```

Now if you check load balancing watchdog, it should show something like that:

```bash
root@USG-3P:/config/scripts# show load-balance watchdog
Group wan_failover
  eth0
  status: Running 
  pings: 5
  fails: 0
  run fails: 0/3
  route drops: 1
  test script : /config/scripts/network-test.sh - OK
  last route drop   : Thu Apr  6 12:03:10 2023
  last route recover: Thu Apr  6 12:04:22 2023

  pppoe1
  status: Running 
  failover-only mode
  pings: 10751
  fails: 1
  run fails: 0/3
  route drops: 0
  ping gateway: ping.ubnt.com - REACHABLE
```

Note how `pppoe1` shows ping gateway, whereas `eth0` uses test script for that. Running 
`show load-balance status` will show you currently have `eth0` as a primary interface.

Now change the latency in the script header to a small number (say 5ms) and wait for a minute or two.
`watchdog` command will now show:

```
Group wan_failover
  eth0
  status: Waiting on recovery (0/3)
  pings: 5
  fails: 3
  run fails: 3/3
  route drops: 1
  test script : /config/scripts/network-test.sh - FAIL
  last route drop   : Thu Apr  6 12:03:10 2023

  pppoe1
  status: Running 
  failover-only mode
  pings: 10739
  fails: 1
  run fails: 0/3
  route drops: 0
  ping gateway: ping.ubnt.com - REACHABLE
```

Changing the latency value back will bring the interface back up (provided your primary network 
connection holds).
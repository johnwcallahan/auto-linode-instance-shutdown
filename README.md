# Auto Linode Instance shutdown

This little program shuts down your Linodes that do not have tags labeled "persist" or "secure".

**WARNING:** This *will* power off your Linodes.

Usage:

1. Build: `$ go build`
2. Set token: `$ export LINODE_TOKEN=your_token_here`
3. Run: `$ ./auto-linode-instance-shutdown`


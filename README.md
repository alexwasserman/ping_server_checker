# ping_server_checker

Pings hosts

Written using ChatGPT in 15 mins so I have a cronjob that can ping a list of hosts. 

When run using `chronic` and the `-quiet` flag, it'll only email the job failure when something is down.

```bash
# ./ping_server_checker --help
Usage of ./ping_server_checker:
  -file string
        Path to file containing server hostnames, one per line
  -quiet
        Suppress non-error output
```

Example use:
```bash
#  ./ping_server_checker
Enter server hostnames, one per line (Press Ctrl+D when done):
google.com
Pinging server: google.com
Server google.com is responsive
All servers are responsive.
```

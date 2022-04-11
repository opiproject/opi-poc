# Software Firewall PoC

The Proof of Concept is made up of three types of containers: traffic sources,
traffic targets, and proxies.  The initial PoC will be:
- Servers hosting traffic sources sending traffic to
- Servers with IPU or DPU cards, with traffic targets on the servers
- Proxies running on IPU or DPU cards

The source, target and proxies are all containers so then can all run on a
laptop.

# Healthy traffic generation
Evaluate using iperf first

# Malicious traffic generation
Since we're using OWASP CRS (see below), we will first evaluate [Zed Attack
Proxy](https://github.com/zaproxy/zaproxy).

# The proxy
The proxy will be open source NGINX with the [SpiderLabs
Modsecurity](https://github.com/SpiderLabs/ModSecurity-nginx) module running
with the open source  [OWASP Core Rule Set
(CRS)](https://github.com/coreruleset/coreruleset).

# The target
No specific target has been chosen, the plan was to evaluate members of [this
list](https://ultimateqa.com/dummy-automation-websites/).  Entries under
consideration will be open source and fully runnable locally, and relatively
simple to containerize if needed.

# Containers
Each of the above will be running in a container.  The long term plan is to have
multiple types of each container, so more than one way to generate healthy
traffic, many types of malicious traffic, multiple sites to automate against,
etc.

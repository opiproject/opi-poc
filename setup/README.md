# What is this?

This is a set of ansible playbooks for helping set up an environment to run the
tests.

## How do you use this?

Pick the OS type that most closely resembles your test environment.

For example, to run on Fedora, RHEL, or similar:

```bash
cd redhat
```

Create an inventory file like this, replacing `mytesthost.domain` with the host
name or IP address to configure and `mytestuser` with the user to remotely log
in as (if different from the local user):

`mytesthost.domain ansible_user=mytestuser`

The ansible_user should have sudo privileges.  Preferably, set up password-less
ssh and password-less sudo for that user, though you can instead instruct
ansible to ask for passwords.

Then run the playbook:

```bash
ansible-playbook -i inventory setup.yml
```

### Example inventory file

The playbooks are set up to run on `hosts: all` so it will run on each host in
your inventory file simultaneously.  You could instead change the playbook to
run it on a group of hosts (e.g., `opi-hosts`) or on individual hosts.  See the
ansible documentation for details.

```text
[opi-hosts]
myhost1.mydomain ansible_user=test
192.168.11.12 ansible_user=foo
```

## What about Debian?

Similarly to the redhat case above, there is an `ubuntu` playbook that may also
work with Debian.

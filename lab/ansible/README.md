# What is this?

This is an ansible playbook for helping set up an OPI LAB environment.

## How do you use this?

The ansible_user should have sudo privileges.  Preferably, set up password-less
ssh and password-less sudo for that user, though you can instead instruct
ansible to ask for passwords.

This is an example how to set password-less:

```bash
 ssh-copy-id root@172.22.X.X
```

Install dependencies:

```bash
python3 -m pip install ansible-pylibssh
ansible-galaxy collection install -r requirements.yml
```

Then run the playbook:

```bash
ansible-playbook -i inventory -e ansible_password='TAKE-FROM-1-PASSWORD' setup.yml
```

## Example log

```bash
root@dh1:~/opi-poc/lab/ansible# ansible-playbook -i inventory setup.yml

PLAY [Set up for OPI LAB environment] *******************************************************************************************************************************************************

TASK [Gathering Facts] **********************************************************************************************************************************************************************
ok: [dh2]
ok: [dh4]
ok: [dh1]
fatal: [dh3]: UNREACHABLE! => {"changed": false, "msg": "Failed to connect to the host via ssh: ssh: connect to host 172.22.1.3 port 22: No route to host", "unreachable": true}
fatal: [amd]: UNREACHABLE! => {"changed": false, "msg": "Failed to connect to the host via ssh: ssh: connect to host 172.22.3.1 port 22: No route to host", "unreachable": true}
fatal: [mrv]: UNREACHABLE! => {"changed": false, "msg": "Failed to connect to the host via ssh: ssh: connect to host 172.22.3.3 port 22: No route to host", "unreachable": true}
ok: [mev]
ok: [bf2]

TASK [Testing] ******************************************************************************************************************************************************************************
changed: [dh2]
changed: [dh4]
changed: [dh1]
changed: [bf2]
changed: [mev]

PLAY RECAP **********************************************************************************************************************************************************************************
amd                        : ok=0    changed=0    unreachable=1    failed=0    skipped=0    rescued=0    ignored=0
bf2                        : ok=2    changed=1    unreachable=0    failed=0    skipped=0    rescued=0    ignored=0
dh1                        : ok=2    changed=1    unreachable=0    failed=0    skipped=0    rescued=0    ignored=0
dh2                        : ok=2    changed=1    unreachable=0    failed=0    skipped=0    rescued=0    ignored=0
dh3                        : ok=0    changed=0    unreachable=1    failed=0    skipped=0    rescued=0    ignored=0
dh4                        : ok=2    changed=1    unreachable=0    failed=0    skipped=0    rescued=0    ignored=0
mev                        : ok=2    changed=1    unreachable=0    failed=0    skipped=0    rescued=0    ignored=0
mrv                        : ok=0    changed=0    unreachable=1    failed=0    skipped=0    rescued=0    ignored=0
```

Demo for Running NGINX on Intel IPU with Red Hat OpenShift
==========================================================

This is for the [opilab][1], which installs an OCP cluster and an IPU host
with RHEL and microshift.

The setup is similar to Red Hat Summit 2025, see the [demo video][2] and the
[blog article about the setup][3].

Demo Setup
----------

In the [opilab][1], we configure an Intel IPU inside host `dh4`. The IPU runs
RHEL with microshift and [dpu-operator][4].

The host `tgen1` runs 3 virtual machines `opicluster-master-{1,2,3}` that are
OCP 4.19 master nodes. The host `dh4` is a worker node in the OCP cluster. The
OCP cluster also has dpu-operator installed.

DPU Operator is installed via [Operatorhub][5].

The cluster was installed and configured manually and using console.redhat.com.
See also the [blog article about the setup][3].

[1]: https://github.com/opiproject/lab
[2]: https://www.linkedin.com/posts/ppindell_update-revealed-at-red-hat-summit-last-year-activity-7333338464651792384-e1yW?utm_source=share&utm_medium=member_ios&rcm=ACoAAAARCJUBkVc4lx8FfmYUF_zaKRbmhvyimTk
[3]: https://access.redhat.com/articles/7120276
[4]: https://github.com/openshift/dpu-operator/
[5]: https://catalog.redhat.com/en/search?gs=&q=dpu&searchType=all

Usage
-----

First connect to VPN to [access the OPI lab][6].

[6]: https://github.com/opiproject/lab?tab=readme-ov-file#access-the-lab

Run `./demo.sh -h` for available commands or inspect the sources of the demo
script.

# Platform/PoC/Reference Architecture Meeting Minutes

## 10/5/2021

Talk about goals: Define a PoC

PoC ideas:

* Nginx already runs on BF2 (but slowly)
* Nginx on Phantom Lake?
* Mt Evans, Oak Springs Canyon complex, likely not MVP

Actions:

* Set up meeting with Marvell: Octeon <- Kris
* Talk to Intel Phantom Lake people about Diamond Bluff <- Tim

Definition of MVP:

* OS - free
* Free app in container: eg Squid Proxy
* Ability to swap out container
* Freely downloadable
* Blog or other publishing

## 11/1/2021

Agenda:

* Refresh on previous discussion
* Follow up on card availability: Tim and Kris
  * F5 has Intel Phantom Lake cards incoming
    * Nginx dpdk support provided by vendor (Silicom?)
  * Intel Big Springs Canyon suggested, available now, Xeon, with IPDK
    * FPGA works out of the box
    * Runtime programming all open source, P4 compiler backend closed
    * ACTION: Tom to investigate hardware
    * ACTION: Tim to verify nginx dpdk support can be made available
  * Intel Mt Evans early access possible as well
    * Requires 5.10 and 5.14 for the two cores
* Get input from new members
* Discuss direction and next steps
  * ACTION: Steve to create github repo for PoC working group
  * OS Selection discussion
    * RHEL with dev sub for initial PoC
    * Figure out OS strategy later
* Meeting cadence and date/time
  * Weekly Wed 1-3pm slot

## 11/12/2021

* Dan: Building cluster
  * Big Springs Canyon on Broadwell with 2x25GbE
  * Mt Evans on Sapphire Rapids with 2x100GbE
  * What OS to deploy?  Start with RHEL 8.5, but Red Hat to investigate OCP
    options.

## 12/1/2021

* Introductions to Kyle Mestery and Maxime Coquelin
* Brief review of the goals of the working group and prior discussions

## 12/15/2021

Merging with the dev platform group

* OS requirements on device
  * Looked at Ubuntu, Fedora, CentOS Stream, and RHEL
  * CentOS Stream 9 with help from the community sounding possible
  * [CentOS SIG](https://wiki.centos.org/SpecialInterestGroup)
* Canceling 12/22/2021 meeting

## 1/5/2022

Light discussion today

* Brief discussion about the Dell presentation, touching on SONiC and whether
  that applies to DB or not (not general purpose OS)
* Kyle mentioned attempting to get IPDK environment working, will present the
  process next week.

## 1/12/2022

* CentOS SIG discussion
* [Meeting recording](https://drive.google.com/file/d/1gPSXwK7xHvFcfWUL3vay7zEaMkjBJQPU/view?usp=sharing)

## 1/19/2022

* CentOS SIG discussion continuation

## 1/26/2022

Dan Daly and Tim Worsley attended

* Dan update on lab: all parts are in
  * Putting CentOS 8.2 on Mt Evans
  * Putting CentOS 7.x? on Big Springs Canyon
* Tim update on PoC workload
  * Simple learning switch based on nginx
  * Pure software, no acceleration
* How do we generate test load in the Intel lab?

For next session:

* Expecting to not have made progress with hardware, so focus on the test
* Create document describing scenarios, test network, traffic generation,
  traffic targets
* Goal is to be able to limit test the solution, which F5 hasn’t done yet
* Also look into open source tools or languages for describing the test, so
  diamond bluff doesn’t have to invent something.  Example would be the open
  ddos (proposed) standard

## 2/1/2022

Attendees: Tim Worsley, Tom Rix, Mark Sanders, Steven Royer

* Discussion on virtual development environment:
  * Similar to the ipdk virtual environment
  * Need volunteers to define and implement virtual development environment
    * Tom Rix volunteered!
  * Prefer to start with arm64
  * Start with RHEL for developer (free)
* OS Image requirements
  * DPDK
  * docker/podman
* App requirements
  * Looking to start with open source nginx in a container
  * Use pure CPU cores initially
  * Tim Worsley has agreed to provide this once a platform exists
* Deliveries:
  * General agreement that DPUs aren’t build systems
  * Go into github: recipes, etc
  * Artifactory???  For container images: Prefer quay.io.

## 2/9/2022

Attendees: Tim Worsley, Steven Royer

* PoC idea
  * Open source nginx with modules like modsec in a container
  * Generate data streams where at least one is “malicious”
    * Container for generating traffic
    * Container as traffic target

## 2/16/2022

Attendees: Kyle Mestery, Maxime Coquelin, Mark Sanders, Michal Kalderon, Laura
JH, Steven Royer

* Review
* Michal Kalderon to look at providing access to Marvell hardware
* Artifacts: look at integrating with github actions + github docker registry

## 2/23/2022

Attendees: Kyle Mestery, Maxime Coquelin, Mark Sanders, Tom Rix, Dan Daly,
Steven Royer, Michal Kalderon

* Dev environment:
  * Provide VM image with some OS (CentOS Stream 9, consider future Debian 11)
    * Tom to talk to Steve about this
* Steve to find out if Red Hat could host some vendor hardware for Diamond Bluff
* Michal to see if something like
  [RH and NVIDIA](https://www.redhat.com/en/blog/optimizing-server-utilization-datacenters-offloading-network-functions-nvidia-bluefield-2-dpus)
  can be done on Marvell hardware

## 3/2/2022

Attendees: Tim Worsely, Michal Kalderon, Laura JH, Steven Royer

* Dev environment follow up: concerns around will adapter vendors support CentOS
  Stream 9 and Debian 11 (drivers etc…)
  * Marvell: yes, drivers are upstream
  * Intel: ??? no attendees
* Hosting vendor hardware by Red Hat:
  * Answer is no, the plan should be to work this through the foundation when
    that is set up.
* Action:
  * Tim: put up PR for current PoC plan

## 3/9/2022

Attendees: Steven Royer, Ted Streete, Michal Kalderon, Mark Sanders, Dan Daly

* Ted leaving the subgroup to focus on other areas
* Tim Worsely out sick
* Kyle on vacation
* Actions:
  * Steve: start OS discussion in slack
  * Steve: update working group description
  * Steve: migrate minutes to github
  * Dan: create document about hardware hosting requirements/options
* Dan: hardware update
  * Big Springs Canyon and Mt Evans in place
  * Access control in place
  * Targeting two weeks for schedule

## 3/23/2022

Attendees: Steven Royer, Gene Bagwell, Shafiq Abedin, Tim Worsley, Michal
Kalderon, Lionel P

* Follow-ups:
  * Steve: Started OS discussion [in slack](https://opi-project.slack.com/archives/C033E418VCK/p1648048161645019)
  * Steve: Update working group description: in progress
  * Steve: Migrate minutes to [github](https://github.com/Diamond-Bluff/dbluff-poc/pull/1)
  * Tim: Still working on providing firewall PoC documentation
* What's the difference between this group and the use cases group?
  * Uses cases group to define use cases
  * PoC group to implement them
* New actions:
  * Tim to provide Dockerfiles for firewall PoC application and tests
  * Tim to publish containers to dockerhub or equivalent
  * Michal to investigate defining VM emulating DPU
  * Steve to add people to the github org for access to private repos
* Task list to be managed via github issues in the poc repo
* Suggestion to move to meeting every other week instead of weekly

## 3/30/2022

Attendees: Steven Royer, Venkat Pullela, Lionel P, Laura JH, Michal Kalderon,
Shafiq Abedin, Mark Sanders

* Meeting cadence
  * Moving to every other week meetings
  * Meeting April 6
* [Review](https://github.com/Diamond-Bluff/dbluff-poc/pull/2)

## 4/6/2022

Attendees: Steven Royer, Venkat Pullela, Shafiq Abedin, Michal Kaderon, Mark
Sanders, Harry Quackenboss

* Venkat volunteers to be involved in setting up open infrastructure, open
  testing that keeps high performance in mind for PoCs
* What additional PoCs should we be thinking about?
  * Layering of function: network, storage, ai, etc implies an ordering
  * Once pure software firewall is done, look at accelerating it
  * Look into storage PoC
    * Shafiq volunteers to work on a storage PoC
* There was a brief discussion where interest was expressed on hearing more
  about the IPDK container environment to use as a starting point for us
  * Ask Kyle for IPDK container presentation at next meeting
* Reminder meetings are now every other week, next meeting is 4/20/2022

## 4/20/2022

Attendees: Steven Royer, Kyle Mestery, Mark Sanders, Dan Daly, Tim Worsley,
Venkat Pullela

* Tim shared his near term plans on the software firewall:
  * PRs for containers over the next weeks
* Kyle shared demo of the new p4-eBPF IPDK container environment
  * [IPDK P4-eBPF](https://github.com/ipdk-io/ipdk/tree/main/build/networking_ebpf)
  * Should be able to clone and replicate easily
  * Should be able to build a version of Tim's software firewall on top of this
    container environment
  * p4-DPDK is currently limited in ways that it is inconvenient to use
    * Dan to work with the team to improve this

## 5/4/2022

Attendees: Steven Royer, Venkat Pullela, Kyle Mestery, Shafiq Abedin

* Work with Mark Sanders and the API group on basic 2 node PoC
  * Can API group build on top of the
    [networking](https://github.com/opiproject/opi-poc/tree/main/networking)
    PoC?
* Can still layer Tim's nginx firewall on top of existing PoC

## 6/15/2022

Attendees: Yuval Caduri, Venkat Pullela, Tim Worsley, Dan Daly, Steven Royer

* Refresh and welcome to Yuval!
* Dan brought up concerns that it's too early to invest much in PoC before the
  other subgroups have produced much
  * General agreements, but there is still value in laying the framework.
* Firewall PoC
  * Tim to update existing networking PoC to include his firewall bits

## 6/29/2022

Attendees: Steven Royer, Dan Daly, Kyle Mestery, Tim Worsley, Anh Thu, Venkat
Pullela

* Talk about OS again
  * Device OS development can occur on any OS, but there needs to be a
    "supported" option, e.g. Red Hat
  * Hope is that applications can be developed in containers with base image of
    their choice.  TBD on how libraries/APIs work across specific base image and
    device OS types.
* New organization discussion
  * Brief discussion follow-on to main meeting topic related to expanding the
    role of this group to include hosting vendor code.  TBD
* Firewall
  * Tim to replace nginx in existing PoC
* Dan: Testing discussion
  * To post slides around CI/CD: slack and github
  * Venkat to bring Keysight people to talk about test pipeline

## 7/13/2022

Attendees: Steven Royer, Kyle Mestery, Boris Glimcher, Timothy Worsley, Dan Daly, Renato Recio, Venkat Pullela, Anh Thu

* Boris presented the current state of the developer platform
* Dan to work with IPDK group for how to merge storage test
  * Need IPDK/SPDK to publish SPDK container image

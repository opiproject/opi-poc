# OPI Lab Goals

The OPI members may contribute some hardware
to be co-managed between the member and the lab staff.
Having hardware in the OPI lab helps
to validate compatibility and to avoid regressions.

## Usage

The OPI lab should be used for:

Phase 0

- define partitioning
- budgeting
- hardware duplication only when needed and if needed

- automatic testing
  - Phase 1:
    - Functional validation and regression
      - scheduled functional validation
      - patch validation (CI)
      - No Access except the lab owner(s) responsible to maintain the CI/CD
    - Conformance
  - Phase 2:
    - Interoperability
      - Reserve IPU/DPUs from multiple vendors to run the test
      - Test portability, compatibility  
- manual
  - Phase 3:
    - Drive adoption  

      - Open Access All IPU/DPUs accessible to all OPI members
      - Demos, Conferences - OPI Lab preferred
      - Learning, Training - Vendor Lab preferred
      - Development lab for SW vendors

  - Phase 4:
  - certification

The performance data generated in the lab shall not be used
for competitive marketing purposes.
Conformance data or certifications from OPI labs can be used for marketing purposes.

## Properties

In order to serve the OPI community, these properties should be targeted:

- reliability: always running, without false positive result
- neutrality: same automatic tests and configurations on all HW/SW platforms
- security: access is restricted to avoid unexpected changes

## Hosting considerations

- 2 years minimum engagement
- 1 year billing intervals

## For consideration

- ease to replicate the testbed
- ease of use
- separation (physical or by time bound reservation) of HW/SW used for automated regression from manual
- clean start before each reservation (restore to factory defaults functionality)

## OPI Lab Requirements

## Access

- restricted access to the automated setup (only maintainers get access)
- reservation based access to the manual setup

## Backup

- All critical servers must be backed up.
- Restoration of critical servers must be done within a couple of hours.
- Hosting provider must be trained to do backup/restore.

## Capacity

- Wait queue for automatic tests should be less than 1 hour.
- Sufficient excess capacity so that manual workloads does not affect automatic tests too much.
- System must support firmware updates.

## Process

- End users given 2 weeks notice for major lab works
- End users given 1 week notice for regular lab maintenance works
- End users notified as soon as it is known immediate emergency lab work need to be done
- Hosting provider must be notified 2 days in advance of a change on a machine
  in order to avoid false negative results due to an upgrade in process.
- Any environment change (firmware, OS update or new software)
  and test changes (new or updated use case)
  must be validated first before being used in automatic tests.
- Lab inventory management system in place
- Lab tech updates inventory as needed
- Lab manager assigned
- Process in place to forecast equipment needs
- Questionnaire in place to be handed out to lab users to help with equipment forecast
- Lab requirements reviewed by OPI TSC
- Sufficient budget in place to fund initial capital costs
- Sufficient budget in place to fund ongoing operations
- Ticketing system in place to request lab services

## Documentation

- Wiki showing servers, assigned owners, emergency console and SSH access info
- Documentation in place for browsing, requesting, and releasing lab resources
  - e.g. how to request for new physical servers with power cycler
- Documentation in place for lab resources and interconnection
  - example: <https://opendcim.org/screenshots.html>

## SLA

- SSH access at 99% availability.
- Automatic tests running at 95% availability (includes maintenance downtime).

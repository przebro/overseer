# goscheduler
[![Build Status](https://travis-ci.com/przebro/goscheduler.svg?token=BuDzHpxjhcjeKFWW17aH&branch=main)](https://travis-ci.com/przebro/goscheduler)
## About
Goscheduler is a workflow manager and a task scheduler. Unlike cron-like schedulers,running tasks in goscheduler are controlled by resources: tickets and flags. this feature makes it easy to create robust and flexible workflows. Gosheduler is inspired by Control-M from BMC.

Note that currently, goscheduler is in a demo stage; therefore, some parts of a project will change.

### Features
* Scheduling options: Manual, Daily, Day of a Week, Selected months, Specific Date
* Time criteria: Active from time, Active to Time
* Hold / Free task
* Confirm task
* Task types: Dummy, OS
* Global variables, Task variables
### TODO
* Cyclic tasks
* Another task types: FTP, Messages, Services, Azure jobs, ...
* Extensive task history
* Users & User roles
* ...
### Installation
```
  git clone https://github.com/przebro/goscheduler
  cd goscheduler
  git checkout develop
  make -f scripts/Makefile build
```
### Sample workflow 
Inside def/samples directory there are sample definitions that consist together for an example workflow.
After successful build binaries can be found inside the bin catalog.
* ovs - scheduler, by default it starts listening on 127.0.0.1:7053
* ovswork - worker, by default it starts listening on 127.0.0.1:7055
* chkprg - sample program
* tools/ovsres - tool, helps manage resources
* tools/ovstask - tool, helps manage tasks

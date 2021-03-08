# overseer
[![Build Status](https://travis-ci.com/przebro/overseer.svg?token=BuDzHpxjhcjeKFWW17aH&branch=develop)](https://travis-ci.com/przebro/overseer)
[![codecov](https://codecov.io/gh/przebro/overseer/branch/develop/graph/badge.svg?token=GGT2W1ARNU)](https://codecov.io/gh/przebro/overseer)
[![Go Report Card](https://goreportcard.com/badge/github.com/przebro/overseer)](https://goreportcard.com/report/github.com/przebro/overseer)

- [About](#about)
- [Features](#features)
- [TODO](#todo)
- [Installation](#installation)
- [Getting Started](#getting-started)

### About
Overseer is a workflow manager and a task scheduler. In overseer tasks are controlled by resources: tickets and flags. This feature makes it easy to create robust and flexible workflows. Overseer is inspired by Control-M from BMC.

Note that currently, overseer is in a demo stage; therefore, some parts of a project will change.

### Features
* Scheduling options: Manual, Daily, Day of a Week, Selected months, Specific Date
* Time criteria: Active from time, Active to Time
* Hold / Free task
* Confirm task
* Task types: Dummy, OS
* Global variables, Task variables
### TODO
* Cyclic tasks
* Flags
* Management of an internal state of components, quiesce mode...
* Another task types: FTP, Messages, Services, Azure jobs, ...
* ...
### Installation
```
  git clone https://github.com/przebro/overseer
  cd overseer
  git checkout develop
  make -f scripts/Makefile build
```
### Getting Started 
Inside def/samples directory there are sample definitions that consist together for an example workflow.
After successful build binaries can be found inside the bin catalog.
* ovs - scheduler, by default it starts listening on 127.0.0.1:7053
* ovswork - worker, by default it starts listening on 127.0.0.1:7055
* chkprg - sample program
* tools/ovscli - tool, helps manage resources and tasks. There is also an Electron-based client:(https://github.com/przebro/overseergui)


**Task definition**
Currently tasks definitions are stored in catalog, by default it is `def` catalog located in a root directory of a project.
The basic task definition is:
```
{
    "type" : "dummy",
    "name" :"minimal",
    "group" : "samples",
    "schedule" :{"type" : "manual"}
}
```
**type**:Defines a kind of a task. Right now, there are two kinds of the task:\
- **dummy** - an empty task that does not call any specific program or service but can have scheduling criteria and can add/remove tickets for other tasks.
- **os** - executes scripts or programs on a worker.
**name** and **group** are unique identifiers of a task. Group represents a subfolde of a root of the definition catalog and name represents the name of a json file.\

**schedule**:Basically, there are two kinds of a task scheduling, manual and time-based. Manual means that task will be not taken into account by daily ordering process, any other type will be checked against specific, time-based criteria:\
- **daily**: The task will be ordered everday.\
- **weekday**: The task will be ordered on a specific day of a week where 1 means Monday and 7 means Sunday.\
- **dayofmonth**: The task will be ordered on a specific day of a month.\
- **fromend**: The task will be ordered on a specific day from the end of a month, where 1 means the end of the month, 2 a day before the end, and so on.
It is relative value so, fromend=2 in July will be resolved to 30 of July and in February it will be resolved to 27 of February or 28 if it is a leap year.\
- **exact**: The task will be ordered exactly on a spcefic day: '2020-05-11','2021-04-07'...\

Each of those values can also be restricted by specifying months, in which the task can be ordered.The above example definition can be changed to:
```
{
    "type" : "dummy",
    "name" :"minimal",
    "description" :"sample dummy task definition",
    "group" : "samples",
    "schedule" :{
    "type" : "daily",
    "months" :[1,3,7,12]
    }
}
```
This task will be ordered every day of January, March, July, and December.\
Despite calendar-based criteria, every task can be restricted to run in a specific period. this can be achieved by setting two properties of schedule:
from and to. For instance:
```
    "schedule" :{
    "type" : "daily",
    "from" : "11:15",
    "to" : "13:00",
    "months" :[1,3,7,12]
    }
```
These values restrict the time window when the task can run to hours between 11:15 and 13:00. A task can be capped at the bottom or the end so, setting only one of the mentioned values will make the time window half-open or half-closed.
**tickets**
Tickets are an essential element of tasks processing. They allow to define dependencies between tasks, use of tickets gives a posibility to build compound 
workflows. A task can require a ticket to run and, after complete their work, it can create a ticket for another task. More than one ticket can be expected by the task before starting the job. In that situation, a relation between input tickets can be defined. If relation is set to 'OR' then one of tickets is required
and if relation is set to 'AND' then all of tickets are requiered.\
Here is an example:\
```
"inticket" : [
        {"name" : "IN-SAMPLE01A","odate" : "ODATE" },
        {"name" : "IN-SAMPLE02A","odate" : "ODATE" },
],
"relation" :"AND"
```
A task with the above declaration can run only if these two tickets exist: IN-SAMPLE01A and IN-SAMPLE02A. Additionally, there is an order date parameter therefore tickets must have the same order date as a task. The "odate" parameter is resolved during the ordering process of a task and it is related to the task's calendar. This field can contain values like: "ODATE", "", "PREV", "NEXT","+nnn", "-nnn" where "nnn" is a number between 001 and 999. As mentioned, the calculated value depends on the task's schedule:\
```
    "schedule" :{
    "type" : "daily",
    }
...
    "schedule" :{
    "type" : "weekday",
    "values" :[3,5,7]
    }
```
If these definitions are ordered on the Friday 8 January 2021, the "ODATE" value will always resolve to the current date for both definitions but, "PREV" and "NEXT" will be different. For the first definition, "NEXT" and "PREV" will be resolved to Thursday and Saturday, but for the second definition "PREV" will resolve to Wednesday(3), and "NEXT" will resolve to Sunday(7).\
For this definition:
```
    "schedule" :{
    "type" : "fromend",
    "values" :[1]
    }
```
The task is ordered on the last day of a month, If the order date is 31 March, the "NEXT" will resolve to 30 April, and "PREV" will resolve to 28 February or,
29 February if it is a leap year.
**Issuing tickets**
There are two ways how the ticket can be added. Manually or, each task definition can contain a section that defines tickets that will be issued after successful completion of a task.
```
"outticket" :[
        {"name" : "IN-SAMPLE01A","odate" : "ODATE" ,"action":"REM"},
        {"name" : "IN-SAMPLE02A","odate" : "ODATE" ,"action":"ADD"},
        {"name" : "IN-SAMPLE02B","odate" : "ODATE" ,"action":"ADD"},
    ]
```
If a task with that definition ends, the ticket "IN-SAMPLE01A" will be removed and, the ticket "IN-SAMPLE02A" and "IN-SAMPLE02B" will be added with order date as a "ODATE". The Same rules as for an "inticket" definition applies to "odate" field in an "outticket" definition so, "NEXT", "PREV" and, +nnn will resolve accordingly to the task's schedule definition.
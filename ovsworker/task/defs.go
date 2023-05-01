package task

import (
	//indirect call init() in drivers packages
	_ "github.com/przebro/overseer/ovsworker/jobs/drv/aws"
	_ "github.com/przebro/overseer/ovsworker/jobs/drv/os"
)

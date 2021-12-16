package task

import (
	//indirect call init() in drivers packages
	_ "overseer/ovsworker/jobs/drv/aws"
	_ "overseer/ovsworker/jobs/drv/dummy"
	_ "overseer/ovsworker/jobs/drv/os"
)

package overseer

import (
	"goscheduler/common/logger"
	"goscheduler/datastore"
	"goscheduler/overseer/config"
	"goscheduler/overseer/internal/events"
	"goscheduler/overseer/internal/pool"
	"goscheduler/overseer/internal/resources"
	"goscheduler/overseer/internal/taskdef"
	"goscheduler/overseer/internal/work"
	"goscheduler/overseer/services"
	"goscheduler/overseer/services/handlers"
	"goscheduler/overseer/services/middleware"
	"path/filepath"
)

//Overseer - main  component
type Overseer struct {
	conf       config.OverseerConfiguration
	resources  resources.ResourceManager
	taskdef    taskdef.TaskDefinitionManager
	taskpool   *pool.ActiveTaskPool
	wrunner    work.WorkerManager
	logger     logger.AppLogger
	ovsGrpcSrv services.OvsGrpcServer
}

//ServerAction - wrapper for a Overseer
type ServerAction interface {
	Start() error
}

//Start - starts server
func (s *Overseer) Start() error {

	var defPath string
	var err error

	s.logger.Info("Starting events dispatcher")
	evDispatcher := events.NewDispatcher()
	s.logger.Info("Starting resources manager")

	if filepath.IsAbs(defPath) {
		defPath = filepath.Join(s.conf.RootDirectory, s.conf.DefinitionDirectory)
	} else {
		defPath = s.conf.DefinitionDirectory
	}

	s.logger.Info("definitions:", defPath)

	s.logger.Info("Start Datastore provider")
	dataProvider, err := datastore.NewDataProvider(s.conf.GetStoreProviderConfiguration())
	if err != nil {
		return err
	}

	s.resources, err = resources.NewManager(evDispatcher, s.logger, s.conf.GetResourceConfiguration(), dataProvider)
	if err != nil {
		s.logger.Error(err)
		return err
	}
	s.logger.Info("Starting definition manager")
	s.taskdef, err = taskdef.NewManager(defPath)

	if err != nil {
		s.logger.Error(err)
		return err
	}

	s.logger.Info("Start timer")
	timer := overseerTimer{}
	timer.tickerFunc(evDispatcher, s.conf.TimeInterval)

	s.logger.Info("Start work runner")
	s.wrunner = work.NewWorkerManager(evDispatcher, s.conf.Workers)
	s.wrunner.Run()

	s.logger.Info("Start taskpool")
	s.taskpool = pool.NewTaskPool(evDispatcher, s.conf.GetActivePoolConfiguration())
	taskManager := pool.NewActiveTaskPoolManager(evDispatcher, s.taskdef, s.taskpool)
	daily := pool.NewDailyExecutor(evDispatcher, taskManager, s.taskpool)

	if s.conf.GetActivePoolConfiguration().ForceNewDayProc {
		s.logger.Info("Forcing New Day Procedure")
		daily.DailyProcedure()
	}

	s.logger.Info("Start grpc")

	tcv, err := services.NewTokenCreatorVerifier(s.conf.GetSecurityConfiguration())
	if err != nil {
		s.logger.Error(err)
		return err
	}

	aservice, err := services.NewAuthenticateService(s.conf.GetSecurityConfiguration(), tcv, dataProvider)
	if err != nil {
		s.logger.Error(err)
		return err
	}

	authhandler, err := handlers.NewServiceAuthorizeHandler(s.conf.GetSecurityConfiguration(), tcv, dataProvider)

	if err != nil {
		s.logger.Error(err)
		return err
	}

	middleware.RegisterHandler(authhandler)

	rservice := services.NewResourceService(s.resources)
	dservice := services.NewDefinistionService(s.taskdef)
	tservice := services.NewTaskService(taskManager, s.taskpool)

	s.ovsGrpcSrv = services.NewOvsGrpcServer(evDispatcher, rservice, dservice, tservice, aservice)

	err = s.ovsGrpcSrv.Listen(s.conf.Host, s.conf.Port)

	return err

}

//NewInstance - creates new instance of a Overseer
func NewInstance(config config.OverseerConfiguration) (*Overseer, error) {

	ov := new(Overseer)
	ov.conf = config
	logDirectory := config.GetLogConfiguration().LogDirectory
	level := config.GetLogConfiguration().LogLevel
	ov.logger = logger.NewLogger(logDirectory, level)

	return ov, nil

}

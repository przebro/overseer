package overseer

import (
	"errors"
	"overseer/common/logger"
	"overseer/datastore"
	"overseer/overseer/auth"
	"overseer/overseer/config"
	"overseer/overseer/internal/events"
	"overseer/overseer/internal/journal"
	"overseer/overseer/internal/pool"
	"overseer/overseer/internal/resources"
	"overseer/overseer/internal/taskdef"
	"overseer/overseer/internal/work"
	"overseer/overseer/services"
	"overseer/overseer/services/handlers"
	"overseer/overseer/services/middleware"
	"path/filepath"
)

//Overseer - main  component
type Overseer struct {
	conf       config.OverseerConfiguration
	resources  resources.ResourceManager
	taskdef    taskdef.TaskDefinitionManager
	taskpool   *pool.ActiveTaskPool
	taskman    *pool.ActiveTaskPoolManager
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

	if !filepath.IsAbs(s.conf.DefinitionDirectory) {
		defPath = filepath.Join(s.conf.Server.RootDirectory, s.conf.DefinitionDirectory)
	} else {
		defPath = s.conf.DefinitionDirectory
	}

	s.logger.Info("definitions:", defPath)

	s.logger.Info("Start Datastore provider")
	dataProvider, err := datastore.NewDataProvider(s.conf.GetStoreProviderConfiguration())
	if err != nil {
		s.logger.Error(err)
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
	s.wrunner = work.NewWorkerManager(evDispatcher, s.conf.GetWorkerManagerConfiguration())
	s.wrunner.Run()

	tJournal, err := journal.NewTaskJournal(s.conf.GetJournalConfiguration(), evDispatcher, dataProvider)
	if err != nil {
		s.logger.Error(err)
		return err
	}

	s.logger.Info("Start taskpool")
	if s.taskpool, err = pool.NewTaskPool(evDispatcher, s.conf.GetActivePoolConfiguration(), dataProvider); err != nil {
		s.logger.Error(err)
		return err
	}

	s.logger.Info("Start taskpool manager")
	if s.taskman, err = pool.NewActiveTaskPoolManager(evDispatcher, s.taskdef, s.taskpool, dataProvider); err != nil {
		s.logger.Error(err)
		return err
	}

	daily := pool.NewDailyExecutor(evDispatcher, s.taskman, s.taskpool)

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

	um, err := auth.NewUserManager(s.conf.GetSecurityConfiguration(), dataProvider)
	if err != nil {
		s.logger.Error(err)
		return err
	}

	rm, err := auth.NewRoleManager(s.conf.GetSecurityConfiguration(), dataProvider)
	if err != nil {
		s.logger.Error(err)
		return err
	}

	am, err := auth.NewRoleAssociationManager(s.conf.GetSecurityConfiguration(), dataProvider)
	if err != nil {
		s.logger.Error(err)
		return err
	}

	rservice := services.NewResourceService(s.resources)
	dservice := services.NewDefinistionService(s.taskdef)
	tservice := services.NewTaskService(s.taskman, s.taskpool, tJournal)
	admservice := services.NewAdministrationService(um, rm, am)
	statservice := services.NewStatusService()

	s.ovsGrpcSrv = services.NewOvsGrpcServer(evDispatcher,
		rservice,
		dservice,
		tservice,
		aservice,
		admservice,
		statservice,
		s.conf.GetServerConfiguration(),
	)

	if s.ovsGrpcSrv == nil {
		return errors.New("fatal error, unable to start grpc server")
	}

	err = s.ovsGrpcSrv.Start()

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

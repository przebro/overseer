package overseer

import (
	"overseer/common/core"
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
	conf             config.OverseerConfiguration
	logger           logger.AppLogger
	srvComponent     core.OverseerComponent
	poolComponent    core.ComponentQuiescer
	journalComponent core.OverseerComponent
	resComponent     core.OverseerComponent
	wmanager         work.WorkerManager
	dispatcher       events.Dispatcher
}

//New - creates a new instance of Overseer
func New(config config.OverseerConfiguration, quiesce bool) (core.RunnableComponent, error) {

	var defPath string
	var err error
	var lg logger.AppLogger
	var pl *pool.ActiveTaskPool
	var pm *pool.ActiveTaskPoolManager
	var dm taskdef.TaskDefinitionManager
	var rm resources.ResourceManager
	var jn journal.TaskJournal
	var gs *services.OvsGrpcServer
	var ds events.Dispatcher

	lg = logger.NewLogger(config.GetLogConfiguration().LogDirectory, config.GetLogConfiguration().LogLevel)

	ds = events.NewDispatcher()

	dataProvider, err := datastore.NewDataProvider(config.GetStoreProviderConfiguration())
	if err != nil {
		lg.Error(err)
		return nil, err
	}

	if !filepath.IsAbs(config.DefinitionDirectory) {
		defPath = filepath.Join(config.Server.RootDirectory, config.DefinitionDirectory)
	} else {
		defPath = config.DefinitionDirectory
	}

	if rm, err = resources.NewManager(ds, lg, config.GetResourceConfiguration(), dataProvider); err != nil {
		return nil, err
	}

	if dm, err = taskdef.NewManager(defPath); err != nil {
		return nil, err
	}

	if pl, err = pool.NewTaskPool(ds, config.PoolConfiguration, dataProvider, !quiesce); err != nil {
		return nil, err
	}

	if pm, err = pool.NewActiveTaskPoolManager(ds, dm, pl, dataProvider); err != nil {
		return nil, err
	}

	if jn, err = journal.NewTaskJournal(config.GetJournalConfiguration(), ds, dataProvider); err != nil {
		return nil, err
	}

	daily := pool.NewDailyExecutor(ds, pm, pl)

	if config.GetActivePoolConfiguration().ForceNewDayProc {
		daily.DailyProcedure()
	}

	wrunner := work.NewWorkerManager(ds, config.GetWorkerManagerConfiguration())

	if gs, err = createServiceServer(config.GetServerConfiguration(), ds, rm, dm, pm, pl, jn, dataProvider, config.GetSecurityConfiguration()); err != nil {
		return nil, err
	}

	ovs := &Overseer{
		logger:           lg,
		srvComponent:     gs,
		poolComponent:    pl,
		journalComponent: jn,
		resComponent:     rm,
		dispatcher:       ds,
		conf:             config,
		wmanager:         wrunner,
	}
	return ovs, nil
}

func createServiceServer(config config.ServerConfiguration,
	disp events.Dispatcher,
	rm resources.ResourceManager,
	dm taskdef.TaskDefinitionManager,
	tm *pool.ActiveTaskPoolManager,
	pv *pool.ActiveTaskPool,
	jrnl journal.TaskJournal,
	provider *datastore.Provider,
	sec config.SecurityConfiguration,
) (*services.OvsGrpcServer, error) {

	tcv, err := services.NewTokenCreatorVerifier(sec)
	if err != nil {
		return nil, err
	}

	aservice, err := services.NewAuthenticateService(sec, tcv, provider)
	if err != nil {
		return nil, err
	}

	authhandler, err := handlers.NewServiceAuthorizeHandler(sec, tcv, provider)

	if err != nil {
		return nil, err
	}

	middleware.RegisterHandler(authhandler)

	um, err := auth.NewUserManager(sec, provider)
	if err != nil {
		return nil, err
	}

	rmn, err := auth.NewRoleManager(sec, provider)
	if err != nil {
		return nil, err
	}

	am, err := auth.NewRoleAssociationManager(sec, provider)
	if err != nil {
		return nil, err
	}

	rservice := services.NewResourceService(rm)
	dservice := services.NewDefinistionService(dm)
	tservice := services.NewTaskService(tm, pv, jrnl)
	admservice := services.NewAdministrationService(um, rmn, am, pv)
	statservice := services.NewStatusService()

	grpcsrv := services.NewOvsGrpcServer(disp,
		rservice,
		dservice,
		tservice,
		aservice,
		admservice,
		statservice,
		config,
	)

	return grpcsrv, nil
}

//Start - starts service
func (s *Overseer) Start() error {

	s.logger.Info("starting worker manager")
	s.wmanager.Run()

	timer := overseerTimer{}
	timer.tickerFunc(s.dispatcher, s.conf.TimeInterval)

	s.logger.Info("starting task journal")
	s.journalComponent.Start()
	s.logger.Info("starting pool")
	s.poolComponent.Start()
	s.logger.Info("starting resource manager")
	s.resComponent.Start()
	s.logger.Info("starting grpc server")
	s.srvComponent.Start()

	return nil
}

//Shutdown - stops service
func (s *Overseer) Shutdown() error {

	s.logger.Info("Shutdown grpc")
	s.srvComponent.Shutdown()
	s.logger.Info("Shutdown pool")
	s.poolComponent.Shutdown()
	s.logger.Info("Shutdown journal")
	s.journalComponent.Shutdown()
	s.logger.Info("Shutdown resource")
	s.resComponent.Shutdown()

	return nil
}

//ServiceName - returns the name of a service
func (s *Overseer) ServiceName() string {
	return s.conf.Server.ServiceName
}

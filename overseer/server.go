package overseer

import (
	"path/filepath"

	"github.com/przebro/overseer/common/core"
	"github.com/przebro/overseer/datastore"
	"github.com/przebro/overseer/overseer/auth"
	"github.com/przebro/overseer/overseer/config"
	"github.com/przebro/overseer/overseer/internal/journal"
	"github.com/przebro/overseer/overseer/internal/pool"
	"github.com/przebro/overseer/overseer/internal/proc"
	"github.com/przebro/overseer/overseer/internal/resources"
	"github.com/przebro/overseer/overseer/internal/taskdef"
	"github.com/przebro/overseer/overseer/internal/work"
	"github.com/przebro/overseer/overseer/services"
	"github.com/przebro/overseer/overseer/services/handlers"
	"github.com/przebro/overseer/overseer/services/middleware"
	"github.com/rs/zerolog/log"
)

// Overseer - main  component
type Overseer struct {
	conf             config.OverseerConfiguration
	srvComponent     core.OverseerComponent
	poolComponent    core.ComponentQuiescer
	journalComponent core.OverseerComponent
	resComponent     core.OverseerComponent
	wmanager         core.OverseerComponent
	timer            core.OverseerComponent
}

// New - creates a new instance of Overseer
func New(config config.OverseerConfiguration, quiesce bool) (core.RunnableComponent, error) {

	var defPath string
	var err error
	var pl *pool.ActiveTaskPool
	var pm *pool.ActiveTaskPoolManager
	var dm taskdef.TaskDefinitionManager
	var rm *resources.ResourceManagerImpl
	var jn *journal.TaskLogJournal
	var gs *services.OvsGrpcServer

	dataProvider, err := datastore.NewDataProvider(config.StoreConfiguration)
	if err != nil {
		log.Err(err).Str("connection_string", config.StoreConfiguration.ConnectionString).Msg("Error creating data provider")
		return nil, err
	}

	if !filepath.IsAbs(config.DefinitionDirectory) {
		defPath = filepath.Join(config.Server.RootDirectory, config.DefinitionDirectory)
	} else {
		defPath = config.DefinitionDirectory
	}

	wrunner := work.NewWorkerWorkManager(config.WorkerManager, config.Server.Security)

	if rm, err = resources.NewManager(config.Resources, dataProvider); err != nil {
		log.Err(err).Msg("Error creating resource manager")
		return nil, err
	}

	if dm, err = taskdef.NewManager(defPath); err != nil {
		log.Err(err).Msg("Error creating task definition manager")
		return nil, err
	}

	if jn, err = journal.NewTaskJournal(config.Journal, dataProvider); err != nil {
		log.Err(err).Msg("Error creating journal")
		return nil, err
	}

	if pl, err = pool.NewTaskPool(config.PoolConfiguration, dataProvider, !quiesce, dm, wrunner, rm, jn); err != nil {
		log.Err(err).Msg("Error creating task pool")
		return nil, err
	}

	if pm, err = pool.NewActiveTaskPoolManager(dm, pl, dataProvider); err != nil {
		log.Err(err).Msg("Error creating task pool manager")
		return nil, err
	}

	daily := proc.NewDailyExecutor(pm, pl)

	if config.PoolConfiguration.ForceNewDayProc {
		daily.DailyProcedure()
	}

	if gs, err = createServiceServer(config.Server, rm, dm, pm, pl, jn, dataProvider, config.Security); err != nil {
		return nil, err
	}

	tm := newTimer(config.TimeInterval)
	tm.AddReceiver(pl)
	tm.AddReceiver(daily)

	ovs := &Overseer{
		srvComponent:     gs,
		poolComponent:    pl,
		journalComponent: jn,
		resComponent:     rm,
		conf:             config,
		wmanager:         wrunner,
		timer:            tm,
	}
	return ovs, nil
}

func createServiceServer(config config.ServerConfiguration,
	rm resources.ResourceManager,
	dm taskdef.TaskDefinitionManager,
	tm *pool.ActiveTaskPoolManager,
	pv *pool.ActiveTaskPool,
	jrnl journal.TaskJournal,
	provider *datastore.Provider,
	sec config.SecurityConfiguration,
) (*services.OvsGrpcServer, error) {

	log := log.With().Str("component", "grpc-service").Logger()
	tcv, err := services.NewTokenCreatorVerifier(sec)
	if err != nil {
		return nil, err
	}

	authman, err := auth.NewAuthenticationManager(provider)
	if err != nil {
		return nil, err
	}

	aservice, err := services.NewAuthenticateService(sec, tcv, authman)
	if err != nil {
		return nil, err
	}

	authhandler, err := handlers.NewServiceAuthorizeHandler(sec, tcv, provider)

	if err != nil {
		return nil, err
	}

	loghandler := handlers.NewServiceLoggerHandler(&log)

	middleware.RegisterHandler(authhandler)
	middleware.RegisterStreamHandler(authhandler)
	middleware.RegisterHandler(loghandler)
	middleware.RegisterStreamHandler(loghandler)

	auth.FirstRun(provider)

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

	grpcsrv := services.NewOvsGrpcServer(
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

// Start - starts service
func (s *Overseer) Start() error {

	lg := log.With().Str("component", "service").Logger()

	lg.Info().Msg("starting worker manager")
	s.wmanager.Start()
	s.timer.Start()

	lg.Info().Msg("starting task journal")
	s.journalComponent.Start()
	lg.Info().Msg("starting pool")
	s.poolComponent.Start()
	lg.Info().Msg("starting resource manager")
	s.resComponent.Start()
	lg.Info().Msg("starting grpc server")
	s.srvComponent.Start()

	return nil
}

// Shutdown - stops service
func (s *Overseer) Shutdown() error {

	lg := log.With().Str("component", "service").Logger()

	lg.Info().Msg("Shutdown grpc")
	s.srvComponent.Shutdown()
	lg.Info().Msg("Shutdown pool")
	s.poolComponent.Shutdown()
	lg.Info().Msg("Shutdown journal")
	s.journalComponent.Shutdown()
	lg.Info().Msg("Shutdown resource")
	s.resComponent.Shutdown()
	lg.Info().Msg("Shutdown worker manager")
	s.wmanager.Shutdown()

	return nil
}

// ServiceName - returns the name of a service
func (s *Overseer) ServiceName() string {
	return s.conf.Server.ServiceName
}

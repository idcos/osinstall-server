package osinstallserver

import (
	"config"
	"config/jsonconf"
	"github.com/AlexanderChen1989/go-json-rest/rest"
	"logger"
	"model"
	"model/mysqlrepo"
	"net/http"
)

type OsInstallServer struct {
	conf    *config.Config
	log     logger.Logger
	repo    model.Repo
	handler http.Handler
}

func NewServer(confPath string, setup PipelineSetupFunc) (*OsInstallServer, error) {
	conf, err := jsonconf.New(confPath).Load()
	if err != nil {
		return nil, err
	}
	log := logger.NewLogrusLogger(conf)
	repo, err := mysqlrepo.NewRepo(conf, log)
	if err != nil {
		return nil, err
	}

	api := rest.NewAPI()

	api.Use(setup(conf, log, repo)...)

	// routes a global
	router, err := rest.MakeRouter(routes...)

	api.SetApp(router)

	server := &OsInstallServer{
		conf:    conf,
		log:     log,
		repo:    repo,
		handler: api.MakeHandler(),
	}

	return server, nil
}

func (server *OsInstallServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	server.handler.ServeHTTP(w, r)
}

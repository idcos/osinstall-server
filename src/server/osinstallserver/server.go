package osinstallserver

import (
	"config"
	"logger"
	"model"
	"model/mysqlrepo"
	"net/http"

	"github.com/AlexanderChen1989/go-json-rest/rest"
)

type OsInstallServer struct {
	Conf    *config.Config
	Log     logger.Logger
	Repo    model.Repo
	handler http.Handler
}

// NewServer 实例化http server
func NewServer(log logger.Logger, conf *config.Config, setup PipelineSetupFunc) (*OsInstallServer, error) {
	repo, err := mysqlrepo.NewRepo(conf, log)
	if err != nil {
		log.Error(err)
		return nil, err
	}

	api := rest.NewAPI()

	api.Use(setup(conf, log, repo)...)

	// routes a global
	router, err := rest.MakeRouter(routes...)
	if err != nil {
		log.Error(err)
		return nil, err
	}

	api.SetApp(router)

	return &OsInstallServer{
		Conf:    conf,
		Log:     log,
		Repo:    repo,
		handler: api.MakeHandler(),
	}, nil
}

func (server *OsInstallServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	server.handler.ServeHTTP(w, r)
}

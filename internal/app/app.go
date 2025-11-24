package app

import (
	"github.com/gin-gonic/gin"

	"github.com/kgugunava/avito-tech-internship/internal/adapters/postgres"
	"github.com/kgugunava/avito-tech-internship/internal/api"
	"github.com/kgugunava/avito-tech-internship/internal/api/handlers"
	"github.com/kgugunava/avito-tech-internship/internal/config"
	"github.com/kgugunava/avito-tech-internship/internal/service"
)

type App struct {
	Cfg config.Config
	Router *gin.Engine
	DB *postgres.Postgres
}

func NewApp() *App {
	app := &App{
		Cfg: config.NewConfig(),
	}
	app.Cfg.InitConfig()

	db := postgres.NewPostgres()

    if err := db.ConnectToPostgresMainDatabase(app.Cfg); err != nil {
        panic(err)
    }

	if err := db.CreateDatabase(app.Cfg); err != nil {
        panic(err)
    }
    
    if err := db.ConnectToDatabase(app.Cfg); err != nil {
        panic(err)
    }
    
    if err := db.CreateDatabaseTables(); err != nil {
        panic(err)
    }
    
    app.DB = &db

	pullRequestRepository := postgres.NewPullRequestRepository(app.DB.Pool)
	teamRepository := postgres.NewTeamRepository(app.DB.Pool)
	userRepository := postgres.NewUserRepository(app.DB.Pool)

	pullRequestService := service.NewPullRequestService(pullRequestRepository)
	teamService := service.NewTeamService(teamRepository)
	userService := service.NewUserService(userRepository)

	apiPullRequests := handlers.NewPullRequestAPI(pullRequestService)
	apiTeams := handlers.NewTeamsAPI(teamService)
	apiUsers := handlers.NewUserAPI(userService)

	apiHandleFunctions := api.ApiHandleFunctions{
		PullRequestsAPI: *apiPullRequests,
		TeamsAPI: *apiTeams,
		UsersAPI: *apiUsers,
	}
    
    app.Router = api.NewRouter(apiHandleFunctions)
    
    return app
}
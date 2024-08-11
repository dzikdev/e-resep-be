package application

import (
	"context"
	"fmt"
	"net/http"

	"e-resep-be/internal/config"

	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/sirupsen/logrus"
	"github.com/xendit/xendit-go/v6"
)

// App..
type App struct {
	Application *echo.Echo
	Context     context.Context
	Config      *config.Configuration
	Logger      *logrus.Logger
	DB          *pgxpool.Pool
	HTTPClient  *http.Client
	XenditSDK   *xendit.APIClient
}

// SetupApplication configuring dependencies app needed
func SetupApplication(ctx context.Context) (*App, error) {
	var err error

	app := &App{}
	app.Context = context.TODO()
	app.Config = config.NewConfig()

	// initialize custom log app with logrus
	logWithLogrus := logrus.New()
	logWithLogrus.Formatter = &logrus.JSONFormatter{}
	logWithLogrus.ReportCaller = true
	app.Logger = logWithLogrus

	// "postgres://username:password@localhost:5432/database_name?sslmode=disable"
	dbpool, err := pgxpool.Connect(context.Background(), fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=%s", app.Config.Database.Username, app.Config.Database.Password, app.Config.Database.Host, app.Config.Database.Port, app.Config.Database.Name, app.Config.Database.SslMode))
	if err != nil {
		app.Logger.Error("failed create pool connection Postgres", err)
		return app, err
	}

	// initialize http client
	app.HTTPClient = &http.Client{}

	// initialize db pool
	app.DB = dbpool

	// initialize echo framework with config
	app.Application = echo.New()
	app.Application.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{"GET", "POST", "PUT", "DELETE"},
		AllowHeaders: []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept, echo.HeaderAuthorization},
	}))
	app.Application.Use(middleware.RequestID())
	app.Application.Use(middleware.Logger())

	// initialize xendit sdk
	app.XenditSDK = xendit.NewClient(app.Config.Xendit.XenditAPIKey)

	app.Logger.Info("APP RUN SUCCESSFULLY ON PORT: ", app.Config.Server.AppPort)

	return app, nil
}

// Close method will close any instances before app terminated
func (a *App) Close(ctx context.Context) {
	a.Logger.Info("APP CLOSED SUCCESSFULLY")

	defer func(ctx context.Context) {
		// DB
		a.DB.Close()
	}(ctx)
}

package main

import (
	"database/sql"
	"fmt"
	"os"

	logger "github.com/codeinuit/test-billing-api/pkg/log"
	"github.com/codeinuit/test-billing-api/pkg/log/logrus"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
)

const (
	POSTGRES_HOST = "POSTGRES_HOST"
	POSTGRES_PORT = "POSTGRES_PORT"
	POSTGRES_PASS = "POSTGRES_PASS"
	POSTGRES_USER = "POSTGRES_USER"
	POSTGRES_DB   = "POSTGRES_DB"

	PORT = "PORT"
)

type Config struct {
	// PostgresDB configuration
	pqHost string
	pqPort string
	pqPass string
	pqUser string
	pqDb   string
}

type BillingAPI struct {
	db     *sql.DB
	log    logger.Logger
	engine *gin.Engine
}

// Run launch the API engine
func (api BillingAPI) Run() (err error) {
	if err = api.engine.Run(); err != nil {
		api.log.Error(err.Error())
	}

	return err
}

// getConfiguration retrieve configuration from ENV
// and returns a Config struct
func getConfiguration() Config {
	return Config{
		pqHost: os.Getenv(POSTGRES_HOST),
		pqPort: os.Getenv(POSTGRES_PORT),
		pqPass: os.Getenv(POSTGRES_PASS),
		pqUser: os.Getenv(POSTGRES_USER),
		pqDb:   os.Getenv(POSTGRES_DB),
	}
}

// pqConnect connects to DB using a given configuration and
// tests the connexion
func pqConnect(c Config) (*sql.DB, error) {
	connStr := fmt.Sprintf("postgresql://%s:%s@%s:%s?sslmode=disable", c.pqUser, c.pqPass, c.pqHost, c.pqPort)

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return db, err
	}

	return db, db.Ping()
}

func main() {
	e := gin.Default()
	l := logrus.NewLogrusLogger()
	c := getConfiguration()

	db, err := pqConnect(c)
	if err != nil {
		l.Error("could not connect to database : ", err.Error())
		return
	}
	l.Info("connexion to pqDb OK")
	defer func() {
		if err := db.Close(); err != nil {
			l.Error("could not close database properly : ", err.Error())
		}
	}()

	api := BillingAPI{db: db, engine: e}
	h := &handlers{log: l, db: &Database{log: l, conn: db}}

	// route declaration
	e.GET("/health", h.healthcheck)
	e.GET("/users", h.getUsers)
	e.POST("/invoice", h.postInvoice)
	e.POST("/transaction", h.postTransaction)

	api.Run()
}

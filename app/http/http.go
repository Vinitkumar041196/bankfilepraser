package http

import (
	"bank_file_analyser/accounts/service"
	"bank_file_analyser/app/http/controllers"
	"bank_file_analyser/config"
	_ "bank_file_analyser/docs"
	"bank_file_analyser/domain"
	"bank_file_analyser/fileparser"
	"context"
	"crypto/tls"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

type httpApp struct {
	Config            *config.AppConfig
	AccBalanceService domain.BalanceGeneratorService
}

func NewHttpApp(conf *config.AppConfig) domain.App {
	//initialize file parser
	parser := fileparser.NewCSVParser(rune(conf.FileColumnSeparator[0]), conf.FileHasHeader)

	//initialize accounts service
	accSrvc := service.NewBalanceGeneratorService(parser, conf.PayRefRegex, conf.DecimalPrecision)

	return &httpApp{Config: conf, AccBalanceService: accSrvc}
}

func (app *httpApp) NewRouter() *gin.Engine {
	router := gin.Default()
	// setup cors
	cors_config := cors.DefaultConfig()
	cors_config.AllowOrigins = []string{"*"}
	cors_config.AllowMethods = []string{"*"}
	cors_config.AllowHeaders = []string{"*"}
	router.Use(cors.New(cors_config))

	router.GET("/docs/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	router.GET("/v1/health_check", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "OK"})
	})

	httpHandler := controllers.NewHTTPHandler(app.AccBalanceService)
	v1 := router.Group("/v1")
	v1.POST("/process_statement", httpHandler.ProcessStatement)

	return router
}

func (app *httpApp) Run() {
	router := app.NewRouter()

	server := &http.Server{
		Addr:    app.Config.ServerAddress,
		Handler: router,
	}

	c := make(chan os.Signal, 1)
	defer close(c)

	signal.Notify(c, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-c
		server.Shutdown(context.Background())
	}()


	log.Printf("Server started at: %s", server.Addr)

	if app.Config.EnableTLS {
		server.ListenAndServeTLS(app.Config.SSLCertPath, app.Config.SSLKeyPath)
	} else {
		//Globally disabling SSL certificate check
		http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
		server.ListenAndServe()
	}

	log.Printf("Server stopped")
}

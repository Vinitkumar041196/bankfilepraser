package app

import (
	"bank_file_analyser/app/http/controllers"
	_ "bank_file_analyser/docs"
	"bank_file_analyser/domain"
	"context"
	"crypto/tls"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func NewRouter(accBalanceService domain.BalanceGeneratorService) *gin.Engine {
	//new router
	router := gin.Default()

	// setup cors middleware
	cors_config := cors.DefaultConfig()
	cors_config.AllowOrigins = []string{"*"}
	cors_config.AllowMethods = []string{"*"}
	cors_config.AllowHeaders = []string{"*"}
	router.Use(cors.New(cors_config))

	//swagger doc endpoint
	router.GET("/docs/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	//health check endpoint
	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "OK"})
	})

	//v1 endpoints
	v1 := router.Group("/v1")
	{
		balancesHandler := controllers.NewBalancesHandler(accBalanceService)
		v1.POST("/process_statement", balancesHandler.ProcessStatement)
	}

	return router
}

func RunHTTPApp(app *App) error {
	//create new router
	router := NewRouter(app.AccBalanceService)

	//set server address
	app.Config.ServerAddress = strings.TrimSpace(app.Config.ServerAddress)
	if app.Config.ServerAddress == "" {
		app.Config.ServerAddress = ":9001" //set default server address
	}

	//create new http server
	server := &http.Server{
		Addr:    app.Config.ServerAddress,
		Handler: router,
	}

	//channel to capture interrupt and shutdown
	c := make(chan os.Signal, 1)
	defer close(c)

	//listen for interrupts
	signal.Notify(c, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-c                                   //keep waiting for interrupt signal
		server.Shutdown(context.Background()) //stop the server
	}()

	log.Println("Server started at: ", server.Addr)

	//Globally disabling SSL certificate check
	http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
	
	//start http server 
	server.ListenAndServe()

	//prints this log on exit
	log.Printf("Server stopped")
	return nil
}

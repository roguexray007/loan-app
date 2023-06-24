package boot

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/roguexray007/loan-app/internal/app/loans"
	"github.com/roguexray007/loan-app/internal/app/loans/payments"
	users "github.com/roguexray007/loan-app/internal/app/users"
	"github.com/roguexray007/loan-app/internal/controllers"
	"github.com/roguexray007/loan-app/internal/provider"
	"github.com/roguexray007/loan-app/internal/routing/router"
)

var (
	loanService *loans.Service
	userService *users.Service
)

func registerDefaultHandlers() *gin.Engine {
	// initialize default router
	ginEngine := router.Initialize()
	return ginEngine
}

func RegisterApplicationHandler(engine *gin.Engine) {
	initializeServices()
	controllers.NewLoanController(loanService)
	controllers.NewUserController(userService)
	// Initialise application routes
	router.InitializeApplicationRoutes(engine)
}

func initializeServices() {
	loanService = newLoanService()
	userService = newUserService()
	createLoanPaymentCore()
}

func newLoanService() *loans.Service {
	loanCore := createLoanCore()
	service := loans.NewLoanService(loanCore)
	return service
}

func createLoanCore() loans.ILoanCore {
	loanRepo := loans.NewRepo(provider.GetDatabase(nil))
	loanCore := loans.NewCore(loanRepo)
	return loanCore
}

func createLoanPaymentCore() payments.ILoanPaymentCore {
	loanPaymentRepo := payments.NewRepo(provider.GetDatabase(nil))
	loanPaymentCore := payments.NewCore(loanPaymentRepo)
	return loanPaymentCore
}

func newUserService() *users.Service {
	userCore := createUserCore()
	service := users.NewUserService(userCore)
	return service
}

func createUserCore() users.IUserCore {
	userRepo := users.NewRepo(provider.GetDatabase(nil))
	userCore := users.NewCore(userRepo)
	return userCore
}

func serve(ctx context.Context, ginEngine *gin.Engine) {

	// Listen to port
	listener, err := net.Listen("tcp4", provider.GetConfig(ctx).App.Port)
	if err != nil {
		panic(err)
	}

	// Serve request - http.Serve
	httpServer := http.Server{
		Handler: ginEngine,
	}

	go func() {
		if err := httpServer.Serve(listener); err != nil && err != http.ErrServerClosed {
			fmt.Println("Failed to start http listener %m", map[string]interface{}{"error": err})
		}
	}()

	c := make(chan os.Signal, 1)

	// accept graceful shutdowns when quit via SIGINT (Ctrl+C) or SIGTERM.
	// SIGKILL, SIGQUIT will not be caught.
	signal.Notify(c, syscall.SIGTERM, syscall.SIGINT)

	// Block until signal is received.
	<-c

	shutDown(ctx, &httpServer)
}

func shutDown(ctx context.Context, httpServer *http.Server) {

	// send unhealthy status to the health check probe and let
	// it mark this pod OOR first before shutting the server down
	fmt.Println("marking server unhealthy")

	// wait for ShutdownDelay seconds
	time.Sleep(time.Duration(provider.GetConfig(ctx).App.ShutdownDelay) * time.Second)

	// Create a deadline to wait for.
	ctxWithTimeout, cancel := context.WithTimeout(
		ctx, time.Duration(
			provider.
				GetConfig(ctx).
				App.ShutdownTimeout)*time.Second)

	defer cancel()

	fmt.Println("shutting down application")
	err := httpServer.Shutdown(ctxWithTimeout)

	if err != nil {
		fmt.Println("failed to initiate shutdown %m", map[string]interface{}{"error": err})
	}
}

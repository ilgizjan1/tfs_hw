package main

import (
	"context"
	"errors"
	"fmt"
	grpchandler "homework10/internal/api/handlers/grpc"
	contracts "homework10/internal/api/handlers/grpc/contracts/langs/go"
	"homework10/internal/api/handlers/grpc/interceptors"
	"homework10/internal/api/handlers/httpgin"
	"homework10/internal/api/handlers/httpgin/middlewares"
	localrepo "homework10/internal/repository/local-repo"

	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc"
	"homework10/internal/service"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

const (
	grpcPortNum = ":50054"
	httpPortNum = ":9000"
)

func main() {
	adService := service.NewAdService(localrepo.NewAdRepo())
	userService := service.NewUserService(localrepo.NewUserRepo())

	grpcListener, err := net.Listen("tcp", grpcPortNum)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	grpcUserMiddleware := interceptors.NewGRPCUserIdentityMiddleware(userService)

	grpcServer := grpc.NewServer(
		grpc.ChainUnaryInterceptor(
			interceptors.LoggingInterceptor,
			interceptors.RecoverInterceptor,
			grpcUserMiddleware.GRPCUserMiddleware,
		),
	)

	grpcAdHandler := grpchandler.NewAdHandler(adService)
	contracts.RegisterAdServiceServer(grpcServer, grpcAdHandler)

	grpcUserHandler := grpchandler.NewUserHandler(userService)
	contracts.RegisterUserServiceServer(grpcServer, grpcUserHandler)

	userMiddleware := middlewares.NewUserIdentityMiddleware(userService)

	httpAdHandler := httpgin.NewAdHandler(adService, userMiddleware)
	httpUserHandler := httpgin.NewUserHandler(userService)
	httpRouter := httpgin.MakeRoutes(httpgin.ApiV1, httpAdHandler, httpUserHandler)

	httpServer := &http.Server{Addr: httpPortNum, Handler: httpRouter}

	eg, ctx := errgroup.WithContext(context.Background())

	sigQuit := make(chan os.Signal, 1)
	signal.Ignore(syscall.SIGHUP, syscall.SIGPIPE)
	signal.Notify(sigQuit, syscall.SIGINT, syscall.SIGTERM)

	eg.Go(func() error {
		select {
		case s := <-sigQuit:
			log.Printf("captured signal: %v\n", s)
			return fmt.Errorf("captured signal: %v", s)
		case <-ctx.Done():
			return nil
		}
	})

	// run grpc server
	eg.Go(func() error {
		log.Printf("starting grpc server, listening on %s\n", grpcPortNum)
		defer log.Printf("close grpc server listening on %s\n", grpcPortNum)

		errCh := make(chan error)

		defer func() {
			grpcServer.GracefulStop()

			_ = grpcListener.Close()

			close(errCh)
		}()

		go func() {
			if err := grpcServer.Serve(grpcListener); err != nil {
				errCh <- err
			}
		}()

		select {
		case <-ctx.Done():
			return ctx.Err()
		case err := <-errCh:
			return fmt.Errorf("grpc server can't listen and serve requests: %w", err)
		}
	})

	// run http server
	eg.Go(func() error {
		log.Printf("starting http server, listening on %s\n", httpServer.Addr)
		defer log.Printf("close http server listening on %s\n", httpServer.Addr)

		errCh := make(chan error)

		defer func() {
			shCtx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
			defer cancel()

			if err := httpServer.Shutdown(shCtx); err != nil {
				log.Printf("can't close http server listening on %s: %s", httpServer.Addr, err.Error())
			}

			close(errCh)
		}()

		go func() {
			if err := httpServer.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
				errCh <- err
			}
		}()

		select {
		case <-ctx.Done():
			return ctx.Err()
		case err := <-errCh:
			return fmt.Errorf("http server can't listen and serve requests: %w", err)
		}
	})

	if err := eg.Wait(); err != nil {
		log.Printf("gracefully shutting down the servers: %s\n", err.Error())
	}

	log.Println("servers were successfully shutdown")
}

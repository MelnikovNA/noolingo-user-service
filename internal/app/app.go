package app

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/MelnikovNA/noolingo-user-service/internal/domain"
	grpcserver "github.com/MelnikovNA/noolingo-user-service/internal/transport/grpc/server"
	"github.com/ilyakaznacheev/cleanenv"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"golang.org/x/sync/errgroup"
)

var (
	ErrSignalReceived = errors.New("signal received")
)

func Run(config string) error {
	cfg := &domain.Config{}

	if err := cleanenv.ReadConfig(config, cfg); err != nil {
		return err
	}

	log := logrus.New()
	log.Infof("Hello app!%#v", cfg)

	eg, egctx := errgroup.WithContext(context.Background())
	g := grpcserver.New(cfg.GrpcServer.Host, cfg.GrpcServer.Port, log)

	_ = egctx
	_ = g

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)

	eg.Go(func() error {
		log.Infof("Grpc listener has started on %s:%s", cfg.GrpcServer.Host, cfg.GrpcServer.Port)
		return g.Serve()
	})

	eg.Go(func() error {
		log.Infof("starting signal handler")
		select {
		case q := <-quit:
			log.Infof("%s signal received, stopping gracefully...", q.String())

			return errors.Wrapf(ErrSignalReceived, "%v", q.String())
		case <-egctx.Done():
			return nil
		}

	})

	// stop servers
	eg.Go(func() error {
		<-egctx.Done()
		g.Stop()
		log.Info("grpc listener has closed")
		return nil
	})

	err := eg.Wait()
	return err
}

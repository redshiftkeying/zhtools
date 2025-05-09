package main

import (
	"context"
	"fmt"
	"github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/jetstream"
	"github.com/redshiftkeying/zhtools/internal/config"
	"log/slog"
	"os"
	"time"
)

type UserService struct {
	js     jetstream.JetStream
	nc     *nats.Conn
	logger *slog.Logger
}

func NewNatsService() (*UserService, error) {
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	slog.SetDefault(logger)

	nc, err := nats.Connect(nats.DefaultURL)
	if err != nil {
		logger.Error(err.Error())
		return nil, err
	}

	js, err := jetstream.New(nc)
	if err != nil {
		nc.Close()
		logger.Error(err.Error())
		return nil, err
	}

	// Ensure stream exists
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err = js.Stream(ctx, config.StreamName)
	if err != nil {
		logger.Error(err.Error(), slog.String("Stream name", config.StreamName))
		logger.Info("attempting to create default stream")
		defaultStreamCfg := jetstream.StreamConfig{
			Name:        config.StreamName,
			Description: "Stores user information",
			Subjects:    []string{fmt.Sprintf("%s.*", config.SubjectPrefix)}, // Subject hierarchy
			MaxAge:      24 * time.Hour,                                      // Example retention policy
			Storage:     jetstream.FileStorage,                               // also allows memory storage
		}
		_, err = js.CreateStream(ctx, defaultStreamCfg)
		if err != nil {
			logger.Error(err.Error())
			nc.Close()
			return nil, fmt.Errorf("failed to create stream '%s': %w", config.StreamName, err)
		}
		logger.Info("created default stream", slog.String("Stream name", config.StreamName))
	} else {
		logger.Info("found existing stream", slog.String("Stream name", config.StreamName))
	}

	return &UserService{
		logger: logger,
		js:     js,
		nc:     nc,
	}, nil
}

func (s *UserService) Close() {
	if s.nc != nil {
		s.nc.Close()
	}
}

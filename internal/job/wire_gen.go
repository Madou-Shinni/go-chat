// Code generated by Wire. DO NOT EDIT.

//go:generate go run github.com/google/wire/cmd/wire
//go:build !wireinject
// +build !wireinject

package main

import (
	"context"
	"github.com/google/wire"
	"go-chat/config"
	"go-chat/internal/cache"
	"go-chat/internal/dao"
	"go-chat/internal/job/internal/command"
	cron2 "go-chat/internal/job/internal/command/cron"
	other2 "go-chat/internal/job/internal/command/other"
	"go-chat/internal/job/internal/command/queue"
	"go-chat/internal/job/internal/handle/cron"
	"go-chat/internal/job/internal/handle/other"
	queue2 "go-chat/internal/job/internal/handle/queue"
	"go-chat/internal/pkg/client"
	"go-chat/internal/pkg/filesystem"
	"go-chat/internal/provider"
)

// Injectors from wire.go:

func Initialize(ctx context.Context, conf *config.Config) *Provider {
	client := provider.NewRedisClient(ctx, conf)
	sidServer := cache.NewSid(client)
	clearWsCacheHandle := cron.NewClearWsCacheHandle(sidServer)
	db := provider.NewMySQLClient(conf)
	filesystemFilesystem := filesystem.NewFilesystem(conf)
	clearArticleHandle := cron.NewClearArticle(db, filesystemFilesystem)
	clearTmpFileHandle := cron.NewClearTmpFile(db, filesystemFilesystem)
	clearExpireServerHandle := cron.NewClearExpireServer(sidServer)
	handles := &cron2.Handles{
		ClearWsCacheHandle:      clearWsCacheHandle,
		ClearArticleHandle:      clearArticleHandle,
		ClearTmpFileHandle:      clearTmpFileHandle,
		ClearExpireServerHandle: clearExpireServerHandle,
	}
	cronCommand := cron2.NewCrontabCommand(handles)
	subcommands := &queue.Subcommands{}
	queueCommand := queue.NewQueueCommand(subcommands)
	exampleHandle := other.NewExampleHandle()
	exampleCommand := other2.NewExampleCommand(exampleHandle)
	otherSubcommands := &other2.Subcommands{
		ExampleCommand: exampleCommand,
	}
	otherCommand := other2.NewOtherCommand(otherSubcommands)
	commands := &command.Commands{
		CrontabCommand: cronCommand,
		QueueCommand:   queueCommand,
		OtherCommand:   otherCommand,
	}
	mainProvider := &Provider{
		Config:   conf,
		Commands: commands,
	}
	return mainProvider
}

// wire.go:

var providerSet = wire.NewSet(provider.NewMySQLClient, provider.NewRedisClient, provider.NewHttpClient, client.NewHttpClient, filesystem.NewFilesystem, cache.NewSid, dao.NewBaseDao, cron2.NewCrontabCommand, cron.NewClearTmpFile, cron.NewClearArticle, cron.NewClearWsCacheHandle, cron.NewClearExpireServer, wire.Struct(new(cron2.Handles), "*"), queue.NewQueueCommand, wire.Struct(new(queue.Subcommands), "*"), queue2.NewEmailHandle, other2.NewOtherCommand, other2.NewExampleCommand, wire.Struct(new(other2.Subcommands), "*"), other.NewExampleHandle, wire.Struct(new(command.Commands), "*"), wire.Struct(new(Provider), "*"))

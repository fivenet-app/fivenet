package cmd

import (
	"context"
	"log"

	"github.com/fivenet-app/fivenet/v2025/pkg/utils/instance"
	"github.com/kardianos/service"
	"go.uber.org/fx"
)

var svcConfig = &service.Config{
	Name:        "FiveNetDBSync",
	DisplayName: "FiveNet DBSync",
	Description: "The DBSync for FiveNet, used to synchronize your gameservers data with your FiveNet instance.",
}

type DBSyncCmd struct {
	RunCmd struct{} `cmd:"" default:"1" help:"Run the DBSync service (default if not subcommand is specified)" name:"run"`

	Start StartCmd `cmd:"" help:"Start the DBSync service via your OS's service manager"`
	Stop  StopCmd  `cmd:"" help:"Stop the DBSync service via your OS's service manager"`

	Install   InstallCmd   `cmd:"" help:"Install the DBSync service to your OS's service manager"`
	Uninstall UninstallCmd `cmd:"" help:"Uninstall the DBSync service from your OS's service manager"`
}

func getService() service.Service {
	fxOpts := getFxBaseOpts(Cli.StartTimeout, false)
	fxOpts = append(fxOpts, FxDBSyncOpts()...)

	app := fx.New(fxOpts...)

	prg := &dbSyncProgram{
		app: app,
	}
	s, err := service.New(prg, svcConfig)
	if err != nil {
		log.Fatal(err)
	}

	return s
}

func (c *DBSyncCmd) Run(_ *Context) error {
	instance.SetComponent("dbsync")

	s := getService()
	if err := s.Run(); err != nil {
		return err
	}

	return nil
}

type StartCmd struct{}

func (c *StartCmd) Run(_ *Context) error {
	log.Println("Starting FiveNet DBSync service")

	s := getService()
	if err := s.Start(); err != nil {
		return err
	}

	log.Println("Started FiveNet DBSync service")

	return nil
}

type StopCmd struct{}

func (c *StopCmd) Run(_ *Context) error {
	log.Println("Stopping FiveNet DBSync service")

	s := getService()
	if err := s.Stop(); err != nil {
		return err
	}

	log.Println("Stopped FiveNet DBSync service")

	return nil
}

type InstallCmd struct{}

func (c *InstallCmd) Run(_ *Context) error {
	log.Println("Installing FiveNet DBSync service")

	s := getService()
	if err := s.Install(); err != nil {
		log.Fatalf("Failed to install service. %v", err)
	}

	log.Println(
		"Service installed successfully. You can now start the service with 'fivenet dbsync start' or via your OS's service manager.",
	)
	return nil
}

type UninstallCmd struct{}

func (c *UninstallCmd) Run(_ *Context) error {
	log.Println("Uninstalling FiveNet DBSync service")

	s := getService()
	if err := s.Uninstall(); err != nil {
		log.Fatalf("Failed to uninstall service. %v", err)
	}

	log.Println("Uninstalled FiveNet DBSync service")

	return nil
}

type dbSyncProgram struct {
	app *fx.App
}

func (p *dbSyncProgram) Start(s service.Service) error {
	go p.app.Run()

	return nil
}

func (p *dbSyncProgram) Stop(s service.Service) error {
	ctx, cancel := context.WithTimeout(context.Background(), stopTimeout)
	defer cancel()

	return p.app.Stop(ctx)
}

package cmd

import (
	"context"
	"errors"
	"log"
	"os"
	"path/filepath"

	"github.com/fivenet-app/fivenet/v2026/pkg/server/admin"
	"github.com/fivenet-app/fivenet/v2026/pkg/utils/instance"
	"github.com/kardianos/service"
	"go.uber.org/fx"
)

var svcConfig = &service.Config{
	Name:        "FiveNetDBSync",
	DisplayName: "FiveNet DBSync",
	Description: "The DBSync for FiveNet, used to synchronize your gameservers data with your FiveNet instance.",
	Arguments:   []string{"dbsync", "run"},
}

type DBSyncCmd struct {
	RunCmd RunCmd `cmd:"" default:"1" help:"Run the DBSync service (default if not subcommand is specified)" name:"run"`

	Start   StartCmd   `cmd:"" help:"Start the DBSync service via your OS's service manager"`
	Restart RestartCmd `cmd:"" help:"Restart the DBSync service via your OS's service manager"`
	Status  StatusCmd  `cmd:"" help:"Get the status of the DBSync service via your OS's service manager"`
	Stop    StopCmd    `cmd:"" help:"Stop the DBSync service via your OS's service manager"`

	Install   InstallCmd   `cmd:"" help:"Install the DBSync service to your OS's service manager"`
	Uninstall UninstallCmd `cmd:"" help:"Uninstall the DBSync service from your OS's service manager"`
}

func getService() service.Service {
	fxOpts := getFxBaseOpts(Cli.StartTimeout, false, false)
	fxOpts = append(fxOpts, FxDBSyncOpts()...)
	fxOpts = append(fxOpts, fx.Invoke(func(admin.AdminServer) {}))

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

type RunCmd struct{}

func (c *RunCmd) Run(ctx *Context) error {
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

type RestartCmd struct{}

func (c *RestartCmd) Run(_ *Context) error {
	log.Println("Restarting FiveNet DBSync service")

	s := getService()
	if err := s.Restart(); err != nil {
		return err
	}

	log.Println("Restarted FiveNet DBSync service")

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

type StatusCmd struct{}

func (c *StatusCmd) Run(_ *Context) error {
	s := getService()
	status, err := s.Status()
	if err != nil {
		return err
	}

	statusMap := map[service.Status]string{
		service.StatusRunning: "Running",
		service.StatusStopped: "Stopped",
	}

	statusStr, exists := statusMap[status]
	if !exists {
		statusStr = "Unknown"
	}

	log.Println("FiveNet DBSync service status:", statusStr)
	return nil
}

type InstallCmd struct{}

func (c *InstallCmd) Run(_ *Context) error {
	log.Println("Installing FiveNet DBSync service")

	// Check default config file location/name
	c.checkIfConfigInWd("dbsync.yaml")

	s := getService()
	if err := s.Install(); err != nil {
		log.Fatalf("Failed to install service. %v", err)
	}

	log.Println(
		"Service installed successfully. You can now start the service with 'fivenet dbsync start' or via your OS's service manager.",
	)
	return nil
}

func (c *InstallCmd) checkIfConfigInWd(cfg string) {
	wd, err := os.Getwd()
	if err != nil {
		log.Fatalf("Failed to get current working directory. %v", err)
	}

	if _, err := os.Stat(filepath.Join(wd, cfg)); err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return
		}
		log.Fatalf("Failed to stat dbsync config file in working directory. %v", err)
	}

	log.Println(
		`WARNING!
The FiveNetDBSync service requires the dbsync config file to be in the /etc/fivenet directory.
You must copy the dbsync config file to the /etc/fivenet directory yourself for the service to find the file successfully.`,
	)
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
	ctx, cancel := context.WithTimeout(context.Background(), defaultStopTimeout)
	defer cancel()

	return p.app.Stop(ctx)
}

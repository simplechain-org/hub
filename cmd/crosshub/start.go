package main

import (
	"fmt"
	"github.com/simplechain-org/crosshub/repo"
	"github.com/simplechain-org/crosshub/swarm"
	"github.com/simplechain-org/go-simplechain/log"
	"github.com/urfave/cli"
)

func startCMD() cli.Command {
	return cli.Command{
		Name:   "start",
		Usage:  "Start a long-running start process",
		Action: start,
	}
}

func start(ctx *cli.Context) error {
	ch := make(chan bool)
	log.Info("start")
	repoRoot, err := repo.PathRootWithDefault(ctx.GlobalString("repo"))
	if err != nil {
		log.Error("PathRootWithDefault","err",err)
		return fmt.Errorf("get repo path: %w", err)
	}

	repo, err := repo.Load(repoRoot)
	if err != nil {
		log.Error("repo.Load","err",err)
		return fmt.Errorf("repo load: %w", err)
	}


	if s,err := swarm.New(repo); err != nil {
		log.Error("swarm.New","err",err)
		return err
	} else {
		if err := s.Start(); err != nil {
			log.Error("s.Start","err",err)
			return err
		}
	}

	<- ch
	return nil
}
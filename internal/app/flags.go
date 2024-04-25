package app

import (
	"flag"
	"fmt"
	"os"

	"github.com/noolingo/user-service/internal/domain"
	"github.com/ilyakaznacheev/cleanenv"
	"gopkg.in/yaml.v3"
)

func parseFlags(cfg *domain.Config) {
	dumpConfig := flag.Bool("dc", false,
		"dump current running config (with env and defaults)")
	descEnv := flag.Bool("de", false, "get descriptions of all environment variables")
	flag.Parse()
	if *dumpConfig {
		b, _ := yaml.Marshal(cfg)
		fmt.Print(string(b))
		os.Exit(0)
	}
	if *descEnv {
		b, err := cleanenv.GetDescription(cfg, nil)
		if err != nil {
			panic(err)
		}
		fmt.Print(b)
		os.Exit(0)
	}
}

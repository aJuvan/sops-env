package sops

import (
	"errors"
	"os"
	"os/exec"

	"github.com/aJuvan/sops-env/config"
	"gopkg.in/yaml.v2"
)

type EnvData struct {
	Env  *map[string]string
}

type output struct {
	Data []byte
}

func (o *output) Write(data []byte) (int, error) {
	o.Data = append(o.Data, data...);
	return len(data), nil;
}

func Sops(conf config.Config) EnvData {
	if conf.RecurseParents {
		err := findFile(conf.File);
		if err != nil {
			config.Log(&conf, config.LogLevelError, "Error while searching for the file:", err);
			os.Exit(1);
		}
	}

	args := append([]string{"-d", conf.File}, conf.SopsExtraArgs...);
	out, err := sopsExec(conf.SopsBinary, args);
	if err != nil {
		config.Log(&conf, config.LogLevelError, "An error occured when running sops:", err);
		os.Exit(1);
	}

	envData := sopsParse(out);
	
	return envData;
}

func findFile(file string) (error) {
	for {
		_, err := os.Stat(file);
		exists := !errors.Is(err, os.ErrNotExist);
		if exists {
			return nil;
		}

		wd, err := os.Getwd();
		if err != nil {
			return err;
		}
		if wd == "/" {
			return os.ErrNotExist;
		}
		os.Chdir("..");
	}
}

func sopsExec(binary string, args[]string) ([]byte, error) {
	var out output;

	process := exec.Command(binary, args...);
	process.Stdin  = os.Stdin
	process.Stdout = &out
	process.Stderr = os.Stderr
	err := process.Run();
	if err != nil {
		return nil, err;
	}

	return out.Data, nil;
}

func sopsParse(data []byte) EnvData {
	var envData EnvData;
	yaml.Unmarshal(data, &envData);
	return envData;
}

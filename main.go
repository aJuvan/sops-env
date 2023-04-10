package main

import (
	"fmt"
	"strings"

	"github.com/aJuvan/sops-env/config"
	"github.com/aJuvan/sops-env/sops"
)

func main() {
	conf := config.GetConfig();
	config.Log(&conf, config.LogLevelDebug, "Parsed config:", conf);
	
	envData := sops.Sops(conf);
	config.Log(&conf, config.LogLevelDebug, "Parsed file:", envData);

	if envData.Env != nil {
		for key, value := range *envData.Env {
			fmt.Println(parseEnv(key, value));
		}
	}
}

func parseEnv(key string, value string) string {
	value = strings.Replace(value, "\\", "\\\\", -1);
	value = strings.Replace(value, "\n", "\\n",  -1);
	value = strings.Replace(value, "\t", "\\\t", -1);
	value = strings.Replace(value, "\"", "\\\"", -1);
	value = strings.Replace(value, " ",  "\\ ",  -1);

	return "export " + key + "=" + value;
}

package config

import (
	"flag"
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

const prefix = "SOPS_ENV__";
const (
	LogLevelDebug = iota + 1
	LogLevelInfo
	LogLevelWarning
	LogLevelError
)
var LogLevels = map[string]int{
	"debug":   LogLevelDebug,
	"info":    LogLevelInfo,
	"warning": LogLevelWarning,
	"error":   LogLevelError,
}
var LogLevelsRev = map[int]string{
	LogLevelDebug:   "debug",
	LogLevelInfo:    "info",
	LogLevelWarning: "warning",
	LogLevelError:   "error",
}

type Config struct {
	File           string
	RecurseParents bool

	SopsBinary     string
	SopsExtraArgs  []string

	LogLevel       int
}

func GetConfig() Config {
	var config Config;

	var logLevel string;

  godotenv.Load();
	flag.Usage = printUsage;

	setFlagString(&config.File, "file", "f", "FILE", "", "Input file for decryption");
	setFlagBool(&config.RecurseParents, "recurse-parents", "p", "RECURSE_PARENTS", "", "Recurse parents for the file");
	setFlagString(&config.SopsBinary, "sops-binary", "b", "SOPS_BINARY", "sops", "Sops binary location (default: sops)");
	setFlagString(&logLevel, "log-level", "l", "LOG_LEVEL", "warning", "Logging level");

	flag.Parse();
	
	if config.File == "" {
		fmt.Fprintln(os.Stderr, "No file specified!");
		printUsage();
	}

	logLevelNum, ok := LogLevels[logLevel]
	if ok == false {
		fmt.Fprintln(os.Stderr, "Invalid log level '" + logLevel + "'!");
		printUsage();
	}

	config.LogLevel = logLevelNum;
	config.SopsExtraArgs = flag.Args();

	return config;
}

func printUsage() {
	fmt.Fprintln(os.Stderr, "Usage: ./" + os.Args[0] + " [flags] <file> [-- [extra arguments]]");
	fmt.Fprintln(os.Stderr, "");
	fmt.Fprintln(os.Stderr, "Flags:");
	fmt.Fprintln(os.Stderr, "\t" + "--file | -f               File to decrypt");
	fmt.Fprintln(os.Stderr, "\t" + "--recurse-parents | -p    File to decrypt");
	fmt.Fprintln(os.Stderr, "\t" + "--sops-binary | -b        SOPS binary location (default: sops)");
	fmt.Fprintln(os.Stderr, "\t" + "--log-level | -l          Logging level (default: warning) [debug, log, warning, error]");
	fmt.Fprintln(os.Stderr, "");
	fmt.Fprintln(os.Stderr, "Environment:");
	fmt.Fprintln(os.Stderr, "\t" + prefix + "FILE           File to decrypt");
	fmt.Fprintln(os.Stderr, "\t" + prefix + "SOPS_BINARY    SOPS binary location");
	fmt.Fprintln(os.Stderr, "\t" + prefix + "LOG_LEVEL      Logging level");
	os.Exit(1);
}

func setFlagString(destination *string, argkey string, argshort string, envkey string, envdefault string, usage string) {
	env := os.Getenv(prefix + envkey);
	if env == "" {
		env = envdefault;
	}

	flag.StringVar(destination, argkey, env, usage);
	flag.StringVar(destination, argshort, env, usage);
}

func setFlagBool(destination *bool, argkey string, argshort string, envkey string, envdefault string, usage string) {
	env := os.Getenv(prefix + envkey);
	if env == "" {
		env = envdefault;
	}
	val := env != "";

	flag.BoolVar(destination, argkey, val, usage);
	flag.BoolVar(destination, argshort, val, usage);
}

func Log(config *Config, LogLevel int, log ...any) {
	if config.LogLevel <= LogLevel {
		all := "[" + LogLevelsRev[LogLevel] + "] ";
		for _, l := range log {
			all += fmt.Sprint(l) + " ";
		}
		fmt.Fprintln(os.Stderr, all);
	}
}

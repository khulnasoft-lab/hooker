package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/signal"
	"runtime"
	"syscall"

	"github.com/khulnasoft-lab/hooker/v2/controller"

	"github.com/khulnasoft-lab/hooker/v2/dbservice"
	"github.com/khulnasoft-lab/hooker/v2/router"
	"github.com/khulnasoft-lab/hooker/v2/runner"
	"github.com/khulnasoft-lab/hooker/v2/utils"
	"github.com/khulnasoft-lab/hooker/v2/webserver"
	"github.com/spf13/cobra"
)

const (
	URL       = "0.0.0.0:8082"
	TLS       = "0.0.0.0:8445"
	URL_USAGE = "The socket to bind to, specified using host:port."
	TLS_USAGE = "The TLS socket to bind to, specified using host:port."
	CFG_FILE  = "/config/cfg.yaml"
	CFG_USAGE = "The alert configuration file."
)

var (
	url            = ""
	tls            = ""
	cfgfile        = ""
	controllerMode = false

	controllerURL          = ""
	controllerCARootPath   = ""
	controllerTLSCertPath  = ""
	controllerTLSKeyPath   = ""
	controllerSeedFilePath = ""
	runnerSeedFilePath     = ""

	runnerName        = ""
	runnerCARootPath  = ""
	runnerTLSCertPath = ""
	runnerTLSKeyPath  = ""
)

var rootCmd = &cobra.Command{
	Use:   "webhooksrv",
	Short: fmt.Sprintf("Khulnasoft Container Security Webhook server\n"),
	Long:  fmt.Sprintf("Khulnasoft Container Security Webhook server\n"),
}

func init() {
	rootCmd.Flags().StringVar(&url, "url", URL, URL_USAGE)
	rootCmd.Flags().StringVar(&tls, "tls", TLS, TLS_USAGE)
	rootCmd.Flags().StringVar(&cfgfile, "cfgfile", CFG_FILE, CFG_USAGE)

	rootCmd.Flags().BoolVar(&controllerMode, "controller-mode", false, "run hooker in controller mode")
	rootCmd.Flags().StringVar(&controllerURL, "controller-url", "", "hooker controller URL")
	rootCmd.Flags().StringVar(&controllerCARootPath, "controller-ca-root", "", "hooker controller ca root file")
	rootCmd.Flags().StringVar(&controllerTLSCertPath, "controller-tls-cert", "", "hooker controller TLS cert file")
	rootCmd.Flags().StringVar(&controllerTLSKeyPath, "controller-tls-key", "", "hooker controller TLS key file")
	rootCmd.Flags().StringVar(&controllerSeedFilePath, "controller-seed-file", "", "hooker controller AuthN seed file")

	rootCmd.Flags().StringVar(&runnerName, "runner-name", "", "hooker runner name")
	rootCmd.Flags().StringVar(&runnerCARootPath, "runner-ca-root", "", "hooker runner ca root file")
	rootCmd.Flags().StringVar(&runnerTLSCertPath, "runner-tls-cert", "", "hooker runner tls cert file")
	rootCmd.Flags().StringVar(&runnerTLSKeyPath, "runner-tls-key", "", "hooker runner tls key file")
	rootCmd.Flags().StringVar(&runnerSeedFilePath, "runner-seed-file", "", "hooker runner AuthN seed file")
}

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	utils.InitDebug()

	rootCmd.Run = func(cmd *cobra.Command, args []string) {
		rtr := router.Instance()

		if runnerName != "" {
			if controllerMode {
				log.Fatal("hooker cannot run as a controller when running in runner mode")
			}

			f, err := ioutil.TempFile("", "temp-hooker-config-*") // TODO: Find a better way
			if err != nil {
				log.Fatal("Unable to create temp file for runner config on disk: ", err)
			}

			rnr := runner.Runner{
				ControllerURL:      controllerURL,
				RunnerSeedFilePath: runnerSeedFilePath,
				RunnerCARootPath:   runnerCARootPath,
				RunnerTLSKeyPath:   runnerTLSKeyPath,
				RunnerTLSCertPath:  runnerTLSCertPath,
				RunnerName:         runnerName,
			}
			if err := rnr.Setup(rtr, f); err != nil {
				log.Fatal("Failed to launch runner: ", err)
			}
			defer func() { os.Remove(f.Name()) }()

			cfgfile = f.Name()
		}

		if controllerMode {
			if runnerName != "" {
				log.Fatal("hooker cannot run as a runner when running in controller mode")
			}

			ctr := controller.Controller{
				ControllerURL:          controllerURL,
				ControllerSeedFilePath: controllerSeedFilePath,
				ControllerCAFile:       controllerCARootPath,
				ControllerTLSKeyPath:   controllerTLSKeyPath,
				ControllerTLSCertPath:  controllerTLSCertPath,
				RunnerName:             runnerName,
			}
			if err := ctr.Setup(rtr); err != nil {
				log.Fatal("Failed to launch controller: ", err)
			}
		}

		if os.Getenv("KHULNASOFTALERT_URL") != "" {
			url = os.Getenv("KHULNASOFTALERT_URL")
		}

		if os.Getenv("HOOKER_HTTP") != "" {
			url = os.Getenv("HOOKER_HTTP")
		}

		if os.Getenv("KHULNASOFTALERT_TLS") != "" {
			tls = os.Getenv("KHULNASOFTALERT_TLS")
		}

		if os.Getenv("HOOKER_HTTPS") != "" {
			tls = os.Getenv("HOOKER_HTTPS")
		}

		if os.Getenv("KHULNASOFTALERT_CFG") != "" {
			cfgfile = os.Getenv("KHULNASOFTALERT_CFG")
		}

		if os.Getenv("HOOKER_CFG") != "" {
			cfgfile = os.Getenv("HOOKER_CFG")
		}

		if os.Getenv("PATH_TO_DB") != "" {
			dbservice.SetNewDbPathFromEnv()
		}

		err := rtr.Start(cfgfile)
		if err != nil {
			log.Printf("Can't start alert manager %v", err)
			return
		}

		defer rtr.Terminate()

		go webserver.Instance().Start(url, tls)
		defer webserver.Instance().Terminate()

		Daemonize()
	}
	err := rootCmd.Execute()
	if err != nil {
		log.Printf("Can't start command %v", err)
		return
	}
}

func Daemonize() {
	sigs := make(chan os.Signal, 1)
	done := make(chan bool, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		sig := <-sigs
		log.Println(sig)
		done <- true
	}()

	<-done
}

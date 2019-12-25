package main

import (
	client_native "github.com/haproxytech/client-native"
	"github.com/haproxytech/client-native/configuration"
	"github.com/haproxytech/client-native/runtime"
	log "github.com/sirupsen/logrus"
	"os"
)

func init() {
	log.SetOutput(os.Stdout)
}

func main() {
	//cmd.Execute()

	confClient := &configuration.Client{}
	confParams := configuration.ClientParams{
		ConfigurationFile:      "/etc/haproxy/haproxy.cfg",
		Haproxy:                "/usr/sbin/haproxy",
		UseValidation:          true,
		PersistentTransactions: true,
		TransactionDir:         "/tmp/haproxy",
	}

	err := confClient.Init(confParams)
	if err != nil {
		log.Println("Error setting up configuration client, using default one")
		confClient, err = configuration.DefaultClient()
		if err != nil {
			log.Println("Error setting up default configuration client, exiting...")
			os.Exit(1)
		}
	}

	runtimeClient := &runtime.Client{}
	_, globalConf, err := confClient.GetGlobalConfiguration("")
	if err == nil {
		socketList := make([]string, 0, 1)
		runtimeAPIs := globalConf.RuntimeApis

		if len(runtimeAPIs) != 0 {
			for _, r := range runtimeAPIs {
				socketList = append(socketList, *r.Address)
			}
			if err := runtimeClient.Init(socketList, "", 0); err != nil {
				log.Println("Error setting up runtime client, not using one")
			}
		} else {
			log.Println("Runtime API not configured, not using it")
			runtimeClient = nil
		}
	} else {
		log.Println("Cannot read runtime API configuration, not using it")
		runtimeClient = nil
	}

	client := &client_native.HAProxyClient{}
	client.Init(confClient, runtimeClient)

	ver, _ := client.Configuration.GetVersion("")
	trx, _ := client.Configuration.StartTransaction(ver)

}

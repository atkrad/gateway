package main

import (
	"fmt"
	client_native "github.com/haproxytech/client-native"
	"github.com/haproxytech/client-native/configuration"
	"github.com/haproxytech/client-native/runtime"
	"github.com/haproxytech/models"
	"github.com/hashicorp/hcl/v2/hclsimple"
	"github.com/imdario/mergo"
	log "github.com/sirupsen/logrus"
	"os"
)

func init() {
	log.SetOutput(os.Stdout)
}

type HaProxyConfiguration struct {
	LogLevel string `hcl:"log_level"`
	EntryPoints []EntryPoint `hcl:"entrypoint,block"`
}

type EntryPoint struct {
	Name string `hcl:"name,label"`
	Bind []string `hcl:"bind"`
}

func main() {
	var hapConf HaProxyConfiguration
	err := hclsimple.DecodeFile("./config.hcl", nil, &hapConf)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	fmt.Printf("%#v \n", hapConf)


	//cmd.Execute()

	confClient := &configuration.Client{}
	confParams := configuration.ClientParams{
		ConfigurationFile:      "/etc/haproxy/haproxy.cfg",
		Haproxy:                "/usr/sbin/haproxy",
		UseValidation:          true,
		PersistentTransactions: true,
		TransactionDir:         "/tmp/haproxy",
	}

	err = confClient.Init(confParams)
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

	transactions, err := client.Configuration.GetTransactions("success")
	fmt.Printf("%T\n", transactions)
	for _, transaction := range *transactions {
		println(transaction.ID)
		println(transaction.Status)
		println(transaction.Version)
		println("======================")
		client.Configuration.CommitTransaction(transaction.ID)
	}

	ver, _ := client.Configuration.GetVersion("")
	trx, _ := client.Configuration.StartTransaction(ver)

	ver, backend, err := client.Configuration.GetBackend("static", trx.ID)
	//fmt.Printf("%#v \n", backend)

	backend.Mode = "tcp"
	//fmt.Printf("%#v \n", backend)
	mergo.Merge(&backend, models.Backend{Mode: "tcp"}, mergo.WithOverride)

	//fmt.Printf("%#v \n", backend)
	client.Configuration.EditBackend("static", backend, trx.ID, 1)
	client.Configuration.CommitTransaction(trx.ID)

}

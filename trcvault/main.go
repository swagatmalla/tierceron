package main

import (
	"log"
	"os"
	"tierceron/trcvault/factory"
	vscutils "tierceron/trcvault/util"
	eUtils "tierceron/utils"
	"tierceron/utils/mlock"

	tclib "VaultConfig.TenantConfig/lib"
	tcutil "VaultConfig.TenantConfig/util"
	"github.com/hashicorp/vault/api"
	"github.com/hashicorp/vault/sdk/plugin"
)

func main() {
	// mLockErr := unix.Mlockall(unix.MCL_CURRENT | unix.MCL_FUTURE)
	// if mLockErr != nil {
	// 	fmt.Println(mLockErr)
	// 	os.Exit(-1)
	// }
	eUtils.InitHeadless(true)
	f, logErr := os.OpenFile("trcvault.log", os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
	logger := log.New(f, "[trcvault]", log.LstdFlags)
	eUtils.CheckError(&eUtils.DriverConfig{Insecure: true, Log: logger, ExitOnFailure: true}, logErr, true)

	tclib.SetLogger(logger.Writer())
	factory.Init(logger)
	mlock.Mlock(logger)

	apiClientMeta := api.PluginAPIClientMeta{}
	flags := apiClientMeta.FlagSet()

	args := os.Args
	vaultHost, lvherr := vscutils.GetLocalVaultHost(false, logger)
	if lvherr != nil {
		logger.Println("Host lookup failure.")
		os.Exit(-1)
	}

	if vaultHost == tcutil.GetLocalVaultAddr() {
		logger.Println("Running in developer mode with self signed certs.")
		args = append(args, "--tls-skip-verify=true")
	} else {
		// TODO: this may not be needed...
		//	args = append(args, fmt.Sprintf("--ca-cert=", caPEM))
	}

	argErr := flags.Parse(args[1:])
	if argErr != nil {
		logger.Fatal(argErr)
	}
	logger.Print("Warming up...")

	tlsConfig := apiClientMeta.GetTLSConfig()
	tlsProviderFunc := api.VaultPluginTLSProvider(tlsConfig)

	err := plugin.Serve(&plugin.ServeOpts{
		BackendFactoryFunc: factory.TrcFactory,
		TLSProviderFunc:    tlsProviderFunc,
	})
	if err != nil {
		logger.Fatal("Plugin shutting down")
	}
}
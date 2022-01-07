package util

import (
	"tierceron/trcconfig/utils"

	"tierceron/vaulthelper/kv"

	sys "tierceron/vaulthelper/system"
)

//Properties stores all configuration properties for a project.
type Properties struct {
	mod          *kv.Modifier
	authMod      *kv.Modifier
	AuthEndpoint string
	cds          *utils.ConfigDataStore
}

/*
//NewProperties stores all configuration properties for a project.
func NewProperties(tokenNamePtr *string, authTokenNamePtr *string, appRoleIDPtr *string, secretIDPtr *string, addrPtr *string, env string, authenv string, project string, service string) (*Properties, error) {
	properties := Properties{}
	if len(*tokenNamePtr) > 0 {
		if len(*appRoleIDPtr) == 0 || len(*secretIDPtr) == 0 {
			eUtils.CheckError(fmt.Errorf("Need both public and secret app role to retrieve token from vault"), true)
		}
		v, err := sys.NewVault(false, *addrPtr, env, false, false, false)
		eUtils.CheckError(err, true)

		mod, err := kv.NewModifier(false, token, *addrPtr, env, nil)
		eUtils.CheckError(err, true)
		mod.Env = "bamboo"

		tokenHandle, err := mod.ReadValue("super-secrets/tokens", *authTokenNamePtr)
		properties.mod, err = kv.NewModifier(false, tokenHandle, *addrPtr, env, nil)
		eUtils.CheckError(err, true)
		properties.mod.Env = env

		properties.cds = new(utils.ConfigDataStore)
		var commonPaths []string
		properties.cds.Init(properties.mod, true, true, project, commonPaths, service)
	}

	return &properties, nil
}
*/
func NewProperties(v *sys.Vault, mod *kv.Modifier, env string, project string, service string) (*Properties, error) {
	properties := Properties{}
	properties.mod = mod
	properties.mod.Env = env

	properties.cds = new(utils.ConfigDataStore)
	var commonPaths []string
	properties.cds.Init(properties.mod, true, true, project, commonPaths, service)

	return &properties, nil
}

//GetValue gets an invididual configuration value for a service from the data store.
func (p *Properties) GetValue(service string, keyPath []string, key string) (string, error) {
	return p.cds.GetValue(service, keyPath, key)
}

//GetConfigValue gets an invididual configuration value for a service from the data store.
func (p *Properties) GetConfigValue(service string, config string, key string) (string, bool) {
	return p.cds.GetConfigValue(service, config, key)
}

//GetConfigValues gets an invididual configuration value for a service from the data store.
func (p *Properties) GetConfigValues(service string, config string) (map[string]interface{}, bool) {
	return p.cds.GetConfigValues(service, config)
}

func ResolveTokenName(env string) string {
	tokenNamePtr := ""
	switch env {
	case "local":
		tokenNamePtr = "config_token_local"
		break
	case "dev":
		tokenNamePtr = "config_token_dev"
		break
	case "QA":
		tokenNamePtr = "config_token_QA"
		break
	case "RQA":
		tokenNamePtr = "config_token_RQA"
		break
	case "staging":
		tokenNamePtr = "config_token_staging"
		break
	default:
		tokenNamePtr = "config_token_local"
		break
	}
	return tokenNamePtr
}

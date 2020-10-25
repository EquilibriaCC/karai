package configuration

import (
	"encoding/json"
	"io/ioutil"

	//"io/ioutil"
	"log"
	"os"
)

func InitConfig() Config {
	var con Config
	r, err := os.Open("./config.json")
	if err != nil {
		r, err = os.Open("./configuration/config.json")
		if err != nil {
			log.Panic("Could not find configuation file", err.Error())
		}
	}
	err = nil
	decoder := json.NewDecoder(r)
	err = decoder.Decode(&con)
	if err != nil {
		log.Panic("Failed to duplicate config file")
	}
	copyConfigFile, err := json.MarshalIndent(&con, "", "    ")
	if err != nil {
		log.Panic("Failed to duplicate config file")
	}
	_ = ioutil.WriteFile("./build/config.json", copyConfigFile, 0644)

	return con
}

func (c Config) GetAppName() string {
	return c.AppName
}

func (c Config) GetAppDev() string {
	return c.AppDev
}

func (c Config) GetConfigDir() string {
	return c.ConfigDir
}

func (c Config) Getp2pConfigDir() string {
	return c.GetConfigDir() + c.P2pConfigDir
}

func (c Config) Getp2pWhitelistDir() string {
	return c.Getp2pConfigDir() + c.P2pWhitelistDir
}

func (c Config) Getp2pBlacklistDir() string {
	return c.Getp2pConfigDir() + c.P2pBlacklistDir
}

func (c Config) GetcertPathDir() string {
	return c.Getp2pConfigDir() + c.CertPathDir
}

func (c Config) GetcertPathSelfDir() string {
	return c.GetcertPathDir() + c.CertPathSelfDir
}

func (c Config) GetcertPathRemote() string {
	return c.GetcertPathDir() + c.CertPathRemote
}

func (c Config) GetWantsClean() bool {
	return c.WantsClean
}

func (c Config) GetDBName() string {
	return c.DbName
}

func (c Config) GetDBUser() string {
	return c.DbUser
}

func (c Config) GetDBSSL() string {
	return c.DbSSL
}

func (c Config) GetTableName() string {
	return c.TableName
}

func (c Config) GetListenPort() int {
	return c.Lport
}

func (c Config) SetkaraiAPIPort(port int) {
	c.Lport = port
}

func (c Config) SetwantsClean(wants_clean bool) {
	c.WantsClean = wants_clean
}

func (c Config) SetconfigDir(dir string) {
	c.ConfigDir = dir
}

func (c Config) Setlport(port int) {
	c.Lport = port
}

func (c Config) SettableName(name string) {
	c.TableName = name
}

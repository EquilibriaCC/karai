package configuration

type Config struct {
	AppName           string `json:"app_name"`
	AppDev            string `json:"app_dev"`
	AppLicense        string `json:"app_license"`
	AppRepository     string `json:"app_repository"`
	AppURL            string `json:"app_url"`
	ConfigDir         string `json:"config_dir"`
	P2pConfigDir      string `json:"p2p_config_dir"`
	P2pWhitelistDir   string `json:"p2p_whitelist_dir"`
	P2pBlacklistDir   string `json:"p2p_blacklist_dir"`
	CertPathDir       string `json:"cert_path_dir"`
	CertPathSelfDir   string `json:"cert_path_self_dir"`
	CertPathRemote    string `json:"cert_path_remote"`
	PubKeyFilePath    string `json:"pubkey_filepath"`
	PrivKeyFilePath   string `json:"privkey_filepath"`
	SignedKeyFilePath string `json:"signedkey_filepath"`
	SelfCertFilePath  string `json:"self_cert_filepath"`
	PeeridPath        string `json:"peerid_path"`
	ChannelName       string `json:"channel_name"`
	ChannelDesc       string `json:"channel_desc"`
	ChannelCont       string `json:"channel_cont"`
	NodeKey           string `json:"node_key"`
	DbUser            string `json:"db_user"`
	DbPassword        string `json:"db_password"`
	DbName            string `json:"db_name"`
	DbPort            string `json:"db_port"`
	DbSSL             string `json:"db_SSL"`
	KaraiAPIPort      int    `json:"karai_api_port"`
	TableName         string `json:"table_name"`
	Lport             int    `json:"listen_port"`
	WantsClean        bool   `json:"wants_clean"`
}
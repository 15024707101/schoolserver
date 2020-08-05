package conf

import (
	"time"
)

type (
	MysqlParams struct {
		AllowAllFiles           bool
		AllowCleartextPasswords bool
		AllowNativePasswords    bool
		AllowOldPasswords       bool
		Charset                 string
		Collation               string
		ClientFoundRows         bool
		ColumnsWithAlias        bool
		InterpolateParams       bool
		Loc                     string
		MaxAllowedPacket        uint32
		MultiStatements         bool
		ParseTime               bool
		ReadTimeout             string
		RejectReadOnly          bool
		Timeout                 string
		Tls                     bool
		WriteTimeout            string

		ParamsStr string
	}

	MysqlConfigBase struct {
		Host       string
		Port       uint32
		Username   string
		Password   string
		Protocol   string
		UnixDomain string
		Dbname     string

		MaxOpenConns    int
		MaxIdleConns    int
		ConnMaxLifetime time.Duration

		Params MysqlParams
	}

	//----------------------------------------------------------------------------------
	// redis配置选项, 拷贝自redis.v6.Options.go部分选项
	RedisOptions struct {
		// The network type, either tcp or unix.
		// Default is tcp.
		Network string
		// host:port address.
		Addr string

		Password string
		// Database to be selected after connecting to the server.
		DB int

		// Maximum number of retries before giving up.
		// Default is to not retry failed commands.
		MaxRetries int
		// Minimum backoff between each retry.
		// Default is 8 milliseconds; -1 disables backoff.
		MinRetryBackoff string
		// Maximum backoff between each retry.
		// Default is 512 milliseconds; -1 disables backoff.
		MaxRetryBackoff string

		// Dial timeout for establishing new connections.
		// Default is 5 seconds.
		DialTimeout string
		// Timeout for socket reads. If reached, commands will fail
		// with a timeout instead of blocking.
		// Default is 3 seconds.
		ReadTimeout string
		// Timeout for socket writes. If reached, commands will fail
		// with a timeout instead of blocking.
		// Default is ReadTimeout.
		WriteTimeout string

		// Maximum number of socket connections.
		// Default is 10 connections per every CPU as reported by runtime.NumCPU.
		PoolSize int
		// Amount of time client waits for connection if all connections
		// are busy before returning an error.
		// Default is ReadTimeout + 1 second.
		PoolTimeout string
		// Amount of time after which client closes idle connections.
		// Should be less than server's timeout.
		// Default is 5 minutes.
		IdleTimeout string
		// Frequency of idle checks.
		// Default is 1 minute.
		// When minus value is set, then idle check is disabled.
		IdleCheckFrequency string

		// Enables read only queries on slave nodes.
		readOnly bool
	}

	RedisConfig struct {
		Redis RedisOptions
	}

	AllConfig struct {
		Mysql       MysqlConfigBase
		MysqlMsg    MysqlConfigBase `yaml:"mysqlmsg"`
		MysqlSchool MysqlConfigBase `yaml:"mysqlschool"`
		Redis       RedisOptions
		Es          []EsConfig
		Sys         SysConfig
	}

	EsConfig struct {
		EsHost string
		EsPort uint
	}
	SysConfig struct {
		MaxMessageNum          int    `yaml:"max_message_num"`
		MaxMessageReceiverNum  int    `yaml:"max_messager_num"`
		MaxMessageReceiverUNum int    `yaml:"max_messageru_num"`
		LogLevel               string `yaml:"log_level"`
		LogFileName            string `yaml:"log_file_name"`
		BackendAddr            string `yaml:"backend_addr"`
		DbLogLevel             string `yaml:"db_log_level"`
		DbLogFilePath          string `yaml:"db_log_file_path"`
		ServiceMode            string `yaml:"service_mode"`
		ServiceName            string `yaml:"service_name"`
		OtherProvince          string `yaml:"other_province"`
		ParamsOutTime          int64  `yaml:"params_out_time"`
		JwtKey                 string `yaml:"jwt_key"`
	}
)

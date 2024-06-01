module github.com/NubeIO/rxlib

go 1.21.6

toolchain go1.21.9

//replace github.com/NubeIO/mqttwrapper => /home/aidan/code/go/mqttwrapper

require (
	github.com/NubeIO/engineering-units v0.0.2
	github.com/NubeIO/mqttwrapper v0.0.9
	github.com/NubeIO/scheduler v0.0.1
	github.com/NubeIO/schema v0.0.1
	github.com/expr-lang/expr v1.16.5
	github.com/gin-gonic/gin v1.10.0
	github.com/go-gota/gota v0.12.0
	github.com/go-resty/resty/v2 v2.11.0
	github.com/golang/protobuf v1.5.3
	github.com/google/uuid v1.6.0
	github.com/grpc-ecosystem/grpc-gateway/v2 v2.19.1
	github.com/jackc/pgx/v4 v4.18.3
	github.com/jackpal/gateway v1.0.14
	github.com/patrickmn/go-cache v2.1.0+incompatible
	github.com/rs/zerolog v1.32.0
	github.com/sashabaranov/go-openai v1.21.0
	github.com/shirou/gopsutil v3.21.11+incompatible
	github.com/sirupsen/logrus v1.4.2
	github.com/spf13/cobra v1.8.0
	github.com/spf13/viper v1.18.2
	github.com/texttheater/golang-levenshtein/levenshtein v0.0.0-20200805054039-cae8b0eaed6c
	github.com/thlib/go-timezone-local v0.0.0-20210907160436-ef149e42d28e
	github.com/tidwall/gjson v1.17.1
	google.golang.org/genproto/googleapis/api v0.0.0-20240221002015-b0ce06bbee7c
	google.golang.org/grpc v1.62.0
	google.golang.org/protobuf v1.34.1
	zgo.at/tz v0.0.0-20230117232324-bf333631bec4
)

require (
	github.com/bytedance/sonic v1.11.6 // indirect
	github.com/bytedance/sonic/loader v0.1.1 // indirect
	github.com/cloudwego/base64x v0.1.4 // indirect
	github.com/cloudwego/iasm v0.2.0 // indirect
	github.com/davecgh/go-spew v1.1.2-0.20180830191138-d8f796af33cc // indirect
	github.com/eclipse/paho.golang v0.21.0 // indirect
	github.com/eclipse/paho.mqtt.golang v1.4.3 // indirect
	github.com/fsnotify/fsnotify v1.7.0 // indirect
	github.com/gabriel-vasile/mimetype v1.4.3 // indirect
	github.com/gin-contrib/sse v0.1.0 // indirect
	github.com/go-co-op/gocron/v2 v2.2.4 // indirect
	github.com/go-ole/go-ole v1.2.6 // indirect
	github.com/go-playground/locales v0.14.1 // indirect
	github.com/go-playground/universal-translator v0.18.1 // indirect
	github.com/go-playground/validator/v10 v10.20.0 // indirect
	github.com/goccy/go-json v0.10.2 // indirect
	github.com/gorilla/websocket v1.5.1 // indirect
	github.com/hashicorp/hcl v1.0.0 // indirect
	github.com/inconshreveable/mousetrap v1.1.0 // indirect
	github.com/jackc/chunkreader/v2 v2.0.1 // indirect
	github.com/jackc/pgconn v1.14.3 // indirect
	github.com/jackc/pgio v1.0.0 // indirect
	github.com/jackc/pgpassfile v1.0.0 // indirect
	github.com/jackc/pgproto3/v2 v2.3.3 // indirect
	github.com/jackc/pgservicefile v0.0.0-20221227161230-091c0ba34f0a // indirect
	github.com/jackc/pgtype v1.14.0 // indirect
	github.com/jackc/puddle v1.3.0 // indirect
	github.com/jonboulle/clockwork v0.4.0 // indirect
	github.com/json-iterator/go v1.1.12 // indirect
	github.com/klauspost/cpuid/v2 v2.2.7 // indirect
	github.com/konsorten/go-windows-terminal-sequences v1.0.2 // indirect
	github.com/leodido/go-urn v1.4.0 // indirect
	github.com/magiconair/properties v1.8.7 // indirect
	github.com/mattn/go-colorable v0.1.13 // indirect
	github.com/mattn/go-isatty v0.0.20 // indirect
	github.com/mitchellh/mapstructure v1.5.0 // indirect
	github.com/modern-go/concurrent v0.0.0-20180306012644-bacd9c7ef1dd // indirect
	github.com/modern-go/reflect2 v1.0.2 // indirect
	github.com/pelletier/go-toml/v2 v2.2.2 // indirect
	github.com/pmezard/go-difflib v1.0.1-0.20181226105442-5d4384ee4fb2 // indirect
	github.com/robfig/cron/v3 v3.0.1 // indirect
	github.com/rogpeppe/go-internal v1.12.0 // indirect
	github.com/sagikazarmark/locafero v0.4.0 // indirect
	github.com/sagikazarmark/slog-shim v0.1.0 // indirect
	github.com/sourcegraph/conc v0.3.0 // indirect
	github.com/spf13/afero v1.11.0 // indirect
	github.com/spf13/cast v1.6.0 // indirect
	github.com/spf13/pflag v1.0.5 // indirect
	github.com/stretchr/objx v0.5.2 // indirect
	github.com/stretchr/testify v1.9.0 // indirect
	github.com/subosito/gotenv v1.6.0 // indirect
	github.com/tidwall/match v1.1.1 // indirect
	github.com/tidwall/pretty v1.2.0 // indirect
	github.com/tklauser/go-sysconf v0.3.13 // indirect
	github.com/tklauser/numcpus v0.7.0 // indirect
	github.com/twitchyliquid64/golang-asm v0.15.1 // indirect
	github.com/ugorji/go/codec v1.2.12 // indirect
	github.com/yusufpapurcu/wmi v1.2.4 // indirect
	go.uber.org/atomic v1.9.0 // indirect
	go.uber.org/multierr v1.9.0 // indirect
	golang.org/x/arch v0.8.0 // indirect
	golang.org/x/crypto v0.23.0 // indirect
	golang.org/x/exp v0.0.0-20231219180239-dc181d75b848 // indirect
	golang.org/x/net v0.25.0 // indirect
	golang.org/x/sync v0.6.0 // indirect
	golang.org/x/sys v0.20.0 // indirect
	golang.org/x/text v0.15.0 // indirect
	gonum.org/v1/gonum v0.9.1 // indirect
	google.golang.org/genproto v0.0.0-20240213162025-012b6fc9bca9 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20240213162025-012b6fc9bca9 // indirect
	gopkg.in/ini.v1 v1.67.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)

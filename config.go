package simutils

import (
	"encoding/json"
	"errors"
	"flag"
	"log"
	"os"
	"path"
	"runtime"
	"strings"
	"time"

	"github.com/go-resty/resty/v2"
	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cast"
	"github.com/spf13/viper"
	"gorm.io/gorm"

	"github.com/alifakhimi/simple-utils-go/simrest"
)

// error block
var (
	ErrHttpServerNotFound = errors.New("http server not found")
	ErrClientNotFound     = errors.New("client not found")
	ErrDBConnNotFound     = errors.New("db connection not found")
)

type (
	Config struct {
		// Name is name of service
		Name string `json:"name,omitempty"`
		// DisplayName displays name of service in human readable
		DisplayName string `json:"display_name,omitempty"`
		// Version show version in banner
		Version string `json:"version,omitempty"`
		// Description
		Description string `json:"description,omitempty"`
		// Website
		Website string `json:"website,omitempty"`
		// Logger
		Logger Logger `json:"logger,omitempty"`
		// Address optionally specifies the TCP address for the server to listen on,
		HttpServers HttpServers `json:"http_servers,omitempty"`
		// Clients is a list of rest client
		Clients simrest.Clients `json:"clients,omitempty"`
		// Databases is a list of database connection
		Databases DBs `json:"databases,omitempty"`
		// Meta
		Meta any `json:"meta,omitempty"`
		// Banner will be displayed when the service starts
		Banners []*Banner `json:"banners,omitempty"`
		// viper is a config tools
		*viper.Viper
	}

	Logger struct {
		Formatter LoggerFormatter `json:"formatter,omitempty"`
		Level     string          `json:"level,omitempty"`
		Output    map[string]any  `json:"output,omitempty"`
	}

	LoggerFormatter struct {
		Use  string              `json:"use,omitempty"`
		JSON JSONLoggerFormatter `json:"json,omitempty"`
		Text TextLoggerFormatter `json:"text,omitempty"`
	}

	JSONLoggerFormatter struct {
		// TimestampFormat sets the format used for marshaling timestamps.
		// The format to use is the same than for time.Format or time.Parse from the standard
		// library.
		// The standard Library already provides a set of predefined format.
		TimestampFormat string `json:"timestamp_format,omitempty"`

		// DisableTimestamp allows disabling automatic timestamps in output
		DisableTimestamp bool `json:"disable_timestamp,omitempty"`

		// DisableHTMLEscape allows disabling html escaping in output
		DisableHTMLEscape bool `json:"disable_html_escape,omitempty"`

		// DataKey allows users to put all the log entry parameters into a nested dictionary at a given key.
		DataKey string `json:"data_key,omitempty"`

		// FieldMap allows users to customize the names of keys for default fields.
		// As an example:
		// formatter := &JSONFormatter{
		//   	FieldMap: FieldMap{
		// 		 FieldKeyTime:  "@timestamp",
		// 		 FieldKeyLevel: "@level",
		// 		 FieldKeyMsg:   "@message",
		// 		 FieldKeyFunc:  "@caller",
		//    },
		// }
		FieldMap FieldMap `json:"field_map,omitempty"`

		// CallerPrettyfier can be set by the user to modify the content
		// of the function and file keys in the json data when ReportCaller is
		// activated. If any of the returned value is the empty string the
		// corresponding key will be removed from json fields.
		CallerPrettyfier func(*runtime.Frame) (function string, file string)

		// PrettyPrint will indent all json logs
		PrettyPrint bool `json:"pretty_print,omitempty"`
	}

	// TextFormatter formats logs into text
	TextLoggerFormatter struct {
		// Set to true to bypass checking for a TTY before outputting colors.
		ForceColors bool `json:"force_colors,omitempty"`

		// Force disabling colors.
		DisableColors bool `json:"disable_colors,omitempty"`

		// Force quoting of all values
		ForceQuote bool `json:"force_quote,omitempty"`

		// DisableQuote disables quoting for all values.
		// DisableQuote will have a lower priority than ForceQuote.
		// If both of them are set to true, quote will be forced on all values.
		DisableQuote bool `json:"disable_quote,omitempty"`

		// Override coloring based on CLICOLOR and CLICOLOR_FORCE. - https://bixense.com/clicolors/
		EnvironmentOverrideColors bool `json:"environment_override_colors,omitempty"`

		// Disable timestamp logging. useful when output is redirected to logging
		// system that already adds timestamps.
		DisableTimestamp bool `json:"disable_timestamp,omitempty"`

		// Enable logging the full timestamp when a TTY is attached instead of just
		// the time passed since beginning of execution.
		FullTimestamp bool `json:"full_timestamp,omitempty"`

		// TimestampFormat to use for display when a full timestamp is printed.
		// The format to use is the same than for time.Format or time.Parse from the standard
		// library.
		// The standard Library already provides a set of predefined format.
		TimestampFormat string `json:"timestamp_format,omitempty"`

		// The fields are sorted by default for a consistent output. For applications
		// that log extremely frequently and don't use the JSON formatter this may not
		// be desired.
		DisableSorting bool `json:"disable_sorting,omitempty"`

		// The keys sorting function, when uninitialized it uses sort.Strings.
		SortingFunc func([]string)

		// Disables the truncation of the level text to 4 characters.
		DisableLevelTruncation bool `json:"disable_level_truncation,omitempty"`

		// PadLevelText Adds padding the level text so that all the levels output at the same length
		// PadLevelText is a superset of the DisableLevelTruncation option
		PadLevelText bool `json:"pad_level_text,omitempty"`

		// QuoteEmptyFields will wrap empty fields in quotes if true
		QuoteEmptyFields bool `json:"quote_empty_fields,omitempty"`

		// FieldMap allows users to customize the names of keys for default fields.
		// As an example:
		// formatter := &TextFormatter{
		//     FieldMap: FieldMap{
		//         FieldKeyTime:  "@timestamp",
		//         FieldKeyLevel: "@level",
		//         FieldKeyMsg:   "@message"}}
		FieldMap FieldMap `json:"field_map,omitempty"`

		// CallerPrettyfier can be set by the user to modify the content
		// of the function and file keys in the data when ReportCaller is
		// activated. If any of the returned value is the empty string the
		// corresponding key will be removed from fields.
		CallerPrettyfier func(*runtime.Frame) (function string, file string)
	}

	FieldMap map[string]string

	LoggerOutputFile struct {
		Type string `json:"type,omitempty"`
		Path string `json:"path,omitempty"`
	}

	Banner struct {
		Text  string `json:"text,omitempty"`
		Font  string `json:"font,omitempty"`
		Color string `json:"color,omitempty"`
	}
)

// Setup logger
func (l *Logger) UnmarshalJSON(data []byte) (err error) {
	var (
		loggerMap = map[string]any{}
		fieldMap  FieldMap
	)

	if err := json.Unmarshal(data, &loggerMap); err != nil {
		return err
	}

	// Log Level
	l.Level = cast.ToString(loggerMap["level"])

	if lvl, err := logrus.ParseLevel(l.Level); err != nil {
		logrus.SetLevel(logrus.DebugLevel)
	} else {
		logrus.SetLevel(lvl)
	}

	// Log Output
	l.Output = cast.ToStringMap(loggerMap["output"])
	logPath := strings.Replace(l.Output["path"].(string), "{{now}}", time.Now().Format(time.RFC3339), -1)

	switch l.Output["type"] {
	case "file":
		if f, err := CreateFile(logPath); err != nil {
			return err
		} else {
			logrus.SetOutput(f)
		}
	default:
		logrus.SetOutput(os.Stdout)
	}

	// Log Formatter
	if b, err := json.Marshal(loggerMap["formatter"]); err != nil {
		return err
	} else if err := json.Unmarshal(b, &l.Formatter); err != nil {
		return err
	}

	switch l.Formatter.Use {
	case "json":
		jsonFormatter := l.Formatter.JSON
		logrus.SetFormatter(&logrus.JSONFormatter{
			TimestampFormat:   jsonFormatter.TimestampFormat,
			DisableTimestamp:  jsonFormatter.DisableTimestamp,
			DisableHTMLEscape: jsonFormatter.DisableHTMLEscape,
			DataKey:           jsonFormatter.DataKey,
			PrettyPrint:       jsonFormatter.PrettyPrint,
		})
		fieldMap = jsonFormatter.FieldMap
	case "text":
		textFormatter := l.Formatter.Text
		logrus.SetFormatter(&logrus.TextFormatter{
			ForceColors:               textFormatter.ForceColors,
			DisableColors:             textFormatter.DisableColors,
			ForceQuote:                textFormatter.ForceQuote,
			DisableQuote:              textFormatter.DisableQuote,
			EnvironmentOverrideColors: textFormatter.EnvironmentOverrideColors,
			DisableTimestamp:          textFormatter.DisableTimestamp,
			FullTimestamp:             textFormatter.FullTimestamp,
			TimestampFormat:           textFormatter.TimestampFormat,
			DisableSorting:            textFormatter.DisableSorting,
			DisableLevelTruncation:    textFormatter.DisableLevelTruncation,
			PadLevelText:              textFormatter.PadLevelText,
			QuoteEmptyFields:          textFormatter.QuoteEmptyFields,
		})
		fieldMap = textFormatter.FieldMap
	default:
		logrus.SetFormatter(&logrus.TextFormatter{})
	}

	// Convert map[string]string to logrus.Fields (map[string]interface{})
	logrusFields := make(logrus.Fields)
	for key, value := range fieldMap {
		logrusFields[key] = value
	}
	if len(logrusFields) > 0 {
		// Use logrus with custom fields
		logrus.WithFields(logrusFields).Info("Logging with custom fields")
	}

	logrus.Infoln("logger configured successfully")

	return nil
}

func NewConfig(path ...string) *Config {
	c := Config{
		HttpServers: HttpServers{},
		Clients:     simrest.Clients{},
		Databases:   DBs{},
	}

	if len(path) > 0 && path[0] != "" {
		if err := ReadConfig(path[0], &c); err != nil {
			logrus.Panicln(err)
		}
	}

	return &c
}

func (conf *Config) GetHttpServer(name string) (h *HttpServer, err error) {
	if len(conf.HttpServers) == 0 {
		return nil, ErrHttpServerNotFound
	}

	if d, ok := conf.HttpServers[name]; !ok {
		return nil, ErrHttpServerNotFound
	} else {
		return d, nil
	}
}

func (conf *Config) GetHttpServerEcho(name string) (ech *echo.Echo, err error) {
	if d, err := conf.GetHttpServer(name); err != nil {
		return nil, err
	} else {
		return d.Echo(), nil
	}
}

func (conf *Config) GetClient(name string) (client *simrest.Client, err error) {
	if len(conf.Clients) == 0 {
		return nil, ErrClientNotFound
	}

	if c, ok := conf.Clients[name]; !ok {
		return nil, ErrClientNotFound
	} else {
		return c, nil
	}
}

func (conf *Config) GetRestyClient(name string) (client *resty.Client, err error) {
	if c, err := conf.GetClient(name); err != nil {
		return nil, err
	} else {
		return c.Client, nil
	}
}

func (conf *Config) GetDB(name string) (db *DBConnection, err error) {
	if len(conf.Databases) == 0 {
		return nil, ErrDBConnNotFound
	}

	if d, ok := conf.Databases[name]; !ok {
		return nil, ErrDBConnNotFound
	} else {
		return d, nil
	}
}

func (conf *Config) GetDBGorm(name string) (db *gorm.DB, err error) {
	if d, err := conf.GetDB(name); err != nil {
		return nil, err
	} else {
		return d.DB, nil
	}
}

func readConfig(path string) error {
	logrus.Infoln("using config file:", path)

	viper.SetConfigType("json")
	viper.SetConfigFile(path)

	if err := viper.ReadInConfig(); err != nil {
		return err
	}

	viper.AutomaticEnv()

	return nil
}

func ReadConfig(path string, conf any) (err error) {
	// Read config from path
	if err = readConfig(path); err != nil {
		return err
	}

	configMap := make(map[string]any)

	if err := viper.Unmarshal(&configMap); err != nil {
		return err
	}

	b, err := json.Marshal(configMap)
	if err != nil {
		return err
	}

	if err := json.Unmarshal(b, conf); err != nil {
		return err
	}

	if c, ok := conf.(*Config); ok {
		c.Viper = viper.GetViper()
	}

	return nil
}

func ReadConfigFromFlag(conf any) {
	var (
		configPath string
	)

	flag.StringVar(&configPath, "c", path.Join(CurrentDirectory(), "config.json"), "config path with json extension")
	flag.Parse()

	if err := ReadConfig(configPath, conf); err != nil {
		log.Fatal(err)
	}
}

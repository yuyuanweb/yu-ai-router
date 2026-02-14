package config

import (
	"bufio"
	"errors"
	"flag"
	"os"
	"strconv"
	"strings"
)

const (
	defaultServerPort     = "8123"
	defaultContextPath    = "/api"
	defaultRedisAddr      = "127.0.0.1:6379"
	defaultRedisDB        = 0
	defaultSessionMaxAge  = 2592000
	defaultSessionName    = "SESSION"
	defaultAIBaseURL      = "https://dashscope.aliyuncs.com/compatible-mode"
	defaultAIModel        = "qwen-plus"
	defaultAIAPIKey       = "YOUR_QWEN_API_KEY"
	defaultAICacheEnabled = true
	defaultAICacheTTLSeconds = 300
	defaultStripeSuccessURL = "http://localhost:5173/recharge/success"
	defaultStripeCancelURL  = "http://localhost:5173/recharge/cancel"
)

type Config struct {
	ServerPort    string
	ContextPath   string
	MySQLDSN      string
	RedisAddr     string
	RedisUsername string
	RedisPassword string
	RedisDB       int

	SessionName   string
	SessionSecret string
	SessionMaxAge int

	AIBaseURL string
	AIAPIKey  string
	AIModel   string
	AICacheEnabled bool
	AICacheTTLSeconds int

	StripeAPIKey       string
	StripeWebhookSecret string
	StripeSuccessURL   string
	StripeCancelURL    string
	PluginSerpAPIKey   string
}

func Load() (*Config, error) {
	baseEnv := currentEnvKeys()
	envName := detectEnvName()
	loadEnvFiles(envName, baseEnv)

	serverPortFlag := flag.String("server-port", "", "server port")
	contextPathFlag := flag.String("server-context-path", "", "server context path")
	envFlag := flag.String("env", "", "runtime env name, such as dev or prod")
	mysqlDSNFlag := flag.String("mysql-dsn", "", "mysql dsn")
	redisAddrFlag := flag.String("redis-addr", "", "redis address")
	redisUsernameFlag := flag.String("redis-username", "", "redis username")
	redisPasswordFlag := flag.String("redis-password", "", "redis password")
	redisDBFlag := flag.Int("redis-db", -1, "redis db index")
	sessionNameFlag := flag.String("session-name", "", "session cookie name")
	sessionSecretFlag := flag.String("session-secret", "", "session secret")
	sessionMaxAgeFlag := flag.Int("session-max-age-seconds", -1, "session max age in seconds")
	aiBaseURLFlag := flag.String("ai-base-url", "", "ai provider base url")
	aiAPIKeyFlag := flag.String("ai-api-key", "", "ai provider api key")
	aiModelFlag := flag.String("ai-model", "", "default ai model")
	stripeAPIKeyFlag := flag.String("stripe-api-key", "", "stripe api key")
	stripeWebhookSecretFlag := flag.String("stripe-webhook-secret", "", "stripe webhook secret")
	stripeSuccessURLFlag := flag.String("stripe-success-url", "", "stripe success callback url")
	stripeCancelURLFlag := flag.String("stripe-cancel-url", "", "stripe cancel callback url")
	pluginSerpAPIKeyFlag := flag.String("plugin-serpapi-api-key", "", "serpapi api key for web_search plugin")
	flag.Parse()
	if *envFlag != "" && *envFlag != envName {
		loadEnvFiles(*envFlag, baseEnv)
	}

	mysqlDSN := pickString(*mysqlDSNFlag, os.Getenv("MYSQL_DSN"), "")
	if mysqlDSN == "" {
		return nil, errors.New("MYSQL_DSN is required (or pass --mysql-dsn)")
	}
	sessionSecret := pickString(*sessionSecretFlag, os.Getenv("SESSION_SECRET"), "")
	if sessionSecret == "" {
		return nil, errors.New("SESSION_SECRET is required (or pass --session-secret)")
	}
	aiAPIKey := pickString(*aiAPIKeyFlag, pickString("", os.Getenv("AI_API_KEY"), os.Getenv("QWEN_API_KEY")), defaultAIAPIKey)

	return &Config{
		ServerPort:    pickString(*serverPortFlag, os.Getenv("SERVER_PORT"), defaultServerPort),
		ContextPath:   pickString(*contextPathFlag, os.Getenv("SERVER_CONTEXT_PATH"), defaultContextPath),
		MySQLDSN:      mysqlDSN,
		RedisAddr:     pickString(*redisAddrFlag, os.Getenv("REDIS_ADDR"), defaultRedisAddr),
		RedisUsername: pickString(*redisUsernameFlag, os.Getenv("REDIS_USERNAME"), ""),
		RedisPassword: pickString(*redisPasswordFlag, os.Getenv("REDIS_PASSWORD"), ""),
		RedisDB:       pickInt(*redisDBFlag, getIntEnvOrDefault("REDIS_DB", defaultRedisDB), defaultRedisDB),
		SessionName:   pickString(*sessionNameFlag, os.Getenv("SESSION_NAME"), defaultSessionName),
		SessionSecret: sessionSecret,
		SessionMaxAge: pickInt(*sessionMaxAgeFlag, getIntEnvOrDefault("SESSION_MAX_AGE_SECONDS", defaultSessionMaxAge), defaultSessionMaxAge),
		AIBaseURL:     pickString(*aiBaseURLFlag, os.Getenv("AI_BASE_URL"), defaultAIBaseURL),
		AIAPIKey:      aiAPIKey,
		AIModel:       pickString(*aiModelFlag, os.Getenv("AI_MODEL"), defaultAIModel),
		AICacheEnabled: getBoolEnvOrDefault("AI_CACHE_ENABLED", defaultAICacheEnabled),
		AICacheTTLSeconds: getIntEnvOrDefault("AI_CACHE_TTL_SECONDS", defaultAICacheTTLSeconds),
		StripeAPIKey:  pickString(*stripeAPIKeyFlag, os.Getenv("STRIPE_API_KEY"), ""),
		StripeWebhookSecret: pickString(*stripeWebhookSecretFlag, os.Getenv("STRIPE_WEBHOOK_SECRET"), ""),
		StripeSuccessURL: pickString(*stripeSuccessURLFlag, os.Getenv("STRIPE_SUCCESS_URL"), defaultStripeSuccessURL),
		StripeCancelURL:  pickString(*stripeCancelURLFlag, os.Getenv("STRIPE_CANCEL_URL"), defaultStripeCancelURL),
		PluginSerpAPIKey: pickString(*pluginSerpAPIKeyFlag, os.Getenv("PLUGIN_SERPAPI_API_KEY"), ""),
	}, nil
}

func detectEnvName() string {
	envName := strings.TrimSpace(os.Getenv("APP_ENV"))
	if envName == "" {
		return "dev"
	}
	return envName
}

func loadEnvFiles(envName string, baseEnv map[string]struct{}) {
	paths := []string{
		".env",
		".env." + envName,
	}
	for _, path := range paths {
		loadEnvFile(path, baseEnv)
	}
}

func loadEnvFile(path string, baseEnv map[string]struct{}) {
	file, err := os.Open(path)
	if err != nil {
		return
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			continue
		}
		key := strings.TrimSpace(parts[0])
		value := strings.TrimSpace(parts[1])
		value = strings.Trim(value, `"'`)
		if key == "" {
			continue
		}
		// Do not override real process-level env,
		// but allow later .env files to override earlier .env files.
		if _, exists := baseEnv[key]; exists {
			continue
		}
		_ = os.Setenv(key, value)
	}
}

func currentEnvKeys() map[string]struct{} {
	result := make(map[string]struct{})
	for _, item := range os.Environ() {
		parts := strings.SplitN(item, "=", 2)
		if len(parts) == 0 {
			continue
		}
		result[parts[0]] = struct{}{}
	}
	return result
}

func pickString(flagValue, envValue, defaultValue string) string {
	if flagValue != "" {
		return flagValue
	}
	if envValue != "" {
		return envValue
	}
	return defaultValue
}

func pickInt(flagValue, envValue, defaultValue int) int {
	if flagValue >= 0 {
		return flagValue
	}
	if envValue >= 0 {
		return envValue
	}
	return defaultValue
}

func getIntEnvOrDefault(key string, defaultValue int) int {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	parsed, err := strconv.Atoi(value)
	if err != nil {
		return defaultValue
	}
	return parsed
}

func getBoolEnvOrDefault(key string, defaultValue bool) bool {
	value := strings.TrimSpace(strings.ToLower(os.Getenv(key)))
	if value == "" {
		return defaultValue
	}
	if value == "1" || value == "true" || value == "yes" || value == "y" {
		return true
	}
	if value == "0" || value == "false" || value == "no" || value == "n" {
		return false
	}
	return defaultValue
}

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

package config

import (
	"fmt"
	"net"
	"os"
	"strconv"
	"strings"
	"time"

	"notex/pkg/types"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Server    ServerConfig        `yaml:"server" json:"server"`
	Database  DatabaseConfig      `yaml:"database" json:"database"`
	Email     EmailConfig         `yaml:"email" json:"email"`
	JWT       JWTConfig           `yaml:"jwt" json:"jwt"`
	RateLimit RateLimitConfig     `yaml:"rate_limit" json:"rate_limit"`
	Storage   types.StorageConfig `yaml:"storage" json:"storage"`
}

type ServerConfig struct {
	Port int    `yaml:"port" json:"port"`
	Host string `yaml:"host" json:"host"`
}

type DatabaseConfig struct {
	Host     string `yaml:"host" json:"host"`
	Port     int    `yaml:"port" json:"port"`
	User     string `yaml:"user" json:"user"`
	Password string `yaml:"password" json:"password"`
	DBName   string `yaml:"dbname" json:"dbname"`
}

type EmailConfig struct {
	Host         string `yaml:"host" json:"host"`
	Port         int    `yaml:"port" json:"port"`
	Username     string `yaml:"username" json:"username"`
	Password     string `yaml:"password" json:"password"`
	FromEmail    string `yaml:"from_email" json:"from_email"`
	FromName     string `yaml:"from_name" json:"from_name"`
	UseTLS       bool   `yaml:"use_tls" json:"use_tls"`
	TemplatesDir string `yaml:"templates_dir" json:"templates_dir"`
	LocalesDir   string `yaml:"locales_dir" json:"locales_dir"`
}

type JWTConfig struct {
	SecretKey          string `yaml:"secret_key" json:"secret_key"`
	AccessTokenExpiry  int    `yaml:"access_token_expiry" json:"access_token_expiry"`   // 小时
	RefreshTokenExpiry int    `yaml:"refresh_token_expiry" json:"refresh_token_expiry"` // 小时
}

// RateLimitItem 限流配置项
type RateLimitItem struct {
	Rate  float64       `yaml:"rate" json:"rate"`   // 每秒请求数
	Burst int           `yaml:"burst" json:"burst"` // 令牌桶容量
	TTL   time.Duration `yaml:"ttl" json:"ttl"`     // 限流器生存时间
}

// RateLimitConfig 限流配置
type RateLimitConfig struct {
	IP    RateLimitItem `yaml:"ip" json:"ip"`
	API   RateLimitItem `yaml:"api" json:"api"`
	Login RateLimitItem `yaml:"login" json:"login"`
}

var (
	DefaultConfig = Config{
		Server: ServerConfig{
			Port: 8080,
			Host: "localhost",
		},
		Database: DatabaseConfig{
			Host:     "localhost",
			Port:     5432,
			User:     "postgres",
			Password: "",
			DBName:   "notex",
		},
		Email: EmailConfig{
			Host:         "smtp.gmail.com",
			Port:         587,
			Username:     "your-email@gmail.com",
			Password:     "your-app-password",
			FromEmail:    "your-email@gmail.com",
			FromName:     "Notex System",
			UseTLS:       true,
			TemplatesDir: "templates/email",
			LocalesDir:   "locales",
		},
		JWT: JWTConfig{
			SecretKey:          "your-secret-key",
			AccessTokenExpiry:  24,  // 24小时
			RefreshTokenExpiry: 168, // 7天
		},
		RateLimit: RateLimitConfig{
			IP: RateLimitItem{
				Rate:  100, // 每秒100个请求
				Burst: 200, // 令牌桶容量200
				TTL:   time.Hour,
			},
			API: RateLimitItem{
				Rate:  1000, // 每秒1000个请求
				Burst: 2000, // 令牌桶容量2000
				TTL:   time.Hour,
			},
			Login: RateLimitItem{
				Rate:  5,  // 每秒5个请求
				Burst: 10, // 令牌桶容量10
				TTL:   time.Hour,
			},
		},
		Storage: types.StorageConfig{
			Type:          "local",
			MaxSize:       10 << 20, // 10MB
			AllowedTypes:  []string{"image/jpeg", "image/png", "image/gif", "image/webp"},
			ThumbnailSize: 300,
			Local: struct {
				UploadDir string `yaml:"upload_dir" json:"upload_dir"`
				URLPrefix string `yaml:"url_prefix" json:"url_prefix"`
			}{
				UploadDir: "uploads",
				URLPrefix: "/uploads/",
			},
		},
	}
	LoadedConfig Config
)

// Validate 验证配置是否有效
func (c *Config) Validate() error {
	// 验证服务器配置
	if err := c.Server.Validate(); err != nil {
		return fmt.Errorf("server config error: %v", err)
	}

	// 验证数据库配置
	if err := c.Database.Validate(); err != nil {
		return fmt.Errorf("database config error: %v", err)
	}

	// 验证邮件配置
	if err := c.Email.Validate(); err != nil {
		return fmt.Errorf("email config error: %v", err)
	}

	// 验证JWT配置
	if err := c.JWT.Validate(); err != nil {
		return fmt.Errorf("jwt config error: %v", err)
	}

	// 验证限流配置
	if err := c.RateLimit.Validate(); err != nil {
		return fmt.Errorf("rate limit config error: %v", err)
	}

	// 验证存储配置
	if err := c.Storage.Validate(); err != nil {
		return fmt.Errorf("storage config error: %v", err)
	}

	return nil
}

// Validate 验证服务器配置
func (c *ServerConfig) Validate() error {
	if c.Port <= 0 || c.Port > 65535 {
		return fmt.Errorf("invalid port number: %d", c.Port)
	}

	if c.Host == "" {
		return fmt.Errorf("host cannot be empty")
	}

	return nil
}

// Validate 验证数据库配置
func (c *DatabaseConfig) Validate() error {
	if c.Port <= 0 || c.Port > 65535 {
		return fmt.Errorf("invalid port number: %d", c.Port)
	}

	if c.Host == "" {
		return fmt.Errorf("host cannot be empty")
	}

	if c.User == "" {
		return fmt.Errorf("user cannot be empty")
	}

	if c.DBName == "" {
		return fmt.Errorf("database name cannot be empty")
	}

	return nil
}

// Validate 验证邮件配置
func (c *EmailConfig) Validate() error {
	if c.Port <= 0 || c.Port > 65535 {
		return fmt.Errorf("invalid port number: %d", c.Port)
	}

	if c.Host == "" {
		return fmt.Errorf("host cannot be empty")
	}

	if c.Username == "" {
		return fmt.Errorf("username cannot be empty")
	}

	if c.Password == "" {
		return fmt.Errorf("password cannot be empty")
	}

	if c.FromEmail == "" {
		return fmt.Errorf("from_email cannot be empty")
	}

	if !isValidEmail(c.FromEmail) {
		return fmt.Errorf("invalid from_email format: %s", c.FromEmail)
	}

	if c.FromName == "" {
		return fmt.Errorf("from_name cannot be empty")
	}

	return nil
}

// Validate 验证JWT配置
func (c *JWTConfig) Validate() error {
	if c.SecretKey == "" {
		return fmt.Errorf("secret_key cannot be empty")
	}

	if len(c.SecretKey) < 32 {
		return fmt.Errorf("secret_key should be at least 32 characters long")
	}

	if c.AccessTokenExpiry <= 0 {
		return fmt.Errorf("access_token_expiry should be positive")
	}

	if c.RefreshTokenExpiry <= 0 {
		return fmt.Errorf("refresh_token_expiry should be positive")
	}

	if c.RefreshTokenExpiry <= c.AccessTokenExpiry {
		return fmt.Errorf("refresh_token_expiry should be greater than access_token_expiry")
	}

	return nil
}

// Validate 验证限流配置
func (c *RateLimitConfig) Validate() error {
	if c.IP.Rate <= 0 {
		return fmt.Errorf("ip rate must be positive")
	}
	if c.IP.Burst <= 0 {
		return fmt.Errorf("ip burst must be positive")
	}
	if c.IP.TTL <= 0 {
		return fmt.Errorf("ip ttl must be positive")
	}

	if c.API.Rate <= 0 {
		return fmt.Errorf("api rate must be positive")
	}
	if c.API.Burst <= 0 {
		return fmt.Errorf("api burst must be positive")
	}
	if c.API.TTL <= 0 {
		return fmt.Errorf("api ttl must be positive")
	}

	if c.Login.Rate <= 0 {
		return fmt.Errorf("login rate must be positive")
	}
	if c.Login.Burst <= 0 {
		return fmt.Errorf("login burst must be positive")
	}
	if c.Login.TTL <= 0 {
		return fmt.Errorf("login ttl must be positive")
	}

	return nil
}

// isValidEmail 验证邮箱格式是否正确
func isValidEmail(email string) bool {
	parts := strings.Split(email, "@")
	if len(parts) != 2 {
		return false
	}

	if len(parts[0]) == 0 || len(parts[1]) == 0 {
		return false
	}

	if _, err := net.LookupMX(parts[1]); err != nil {
		return false
	}

	return true
}

// LoadConfig 加载配置文件
func LoadConfig(configPath string) error {
	// 使用默认配置
	LoadedConfig = DefaultConfig

	// 如果配置文件路径为空，直接返回
	if configPath == "" {
		return nil
	}

	// 读取配置文件
	data, err := os.ReadFile(configPath)
	if err != nil {
		return fmt.Errorf("failed to read config file: %v", err)
	}

	// 解析YAML
	if err := yaml.Unmarshal(data, &LoadedConfig); err != nil {
		return fmt.Errorf("failed to parse config file: %v", err)
	}

	// 验证配置
	if err := LoadedConfig.Validate(); err != nil {
		return fmt.Errorf("invalid config: %v", err)
	}

	// 打印加载的配置
	fmt.Printf("Loaded configuration:\n")
	fmt.Printf("Database: host=%s, port=%d, user=%s, dbname=%s\n",
		LoadedConfig.Database.Host,
		LoadedConfig.Database.Port,
		LoadedConfig.Database.User,
		LoadedConfig.Database.DBName)

	return nil
}

// overrideWithEnv 使用环境变量覆盖配置
func overrideWithEnv(cfg *Config) {
	// 数据库配置
	if dbHost := os.Getenv("DB_HOST"); dbHost != "" {
		cfg.Database.Host = dbHost
	}
	if dbPort := os.Getenv("DB_PORT"); dbPort != "" {
		if port, err := strconv.Atoi(dbPort); err == nil {
			cfg.Database.Port = port
		}
	}
	if dbUser := os.Getenv("DB_USER"); dbUser != "" {
		cfg.Database.User = dbUser
	}
	if dbPass := os.Getenv("DB_PASSWORD"); dbPass != "" {
		cfg.Database.Password = dbPass
	}
	if dbName := os.Getenv("DB_NAME"); dbName != "" {
		cfg.Database.DBName = dbName
	}

	// 邮件配置
	if emailHost := os.Getenv("EMAIL_HOST"); emailHost != "" {
		cfg.Email.Host = emailHost
	}
	if emailPort := os.Getenv("EMAIL_PORT"); emailPort != "" {
		if port, err := strconv.Atoi(emailPort); err == nil {
			cfg.Email.Port = port
		}
	}
	if emailUser := os.Getenv("EMAIL_USERNAME"); emailUser != "" {
		cfg.Email.Username = emailUser
	}
	if emailPass := os.Getenv("EMAIL_PASSWORD"); emailPass != "" {
		cfg.Email.Password = emailPass
	}

	// JWT配置
	if jwtSecret := os.Getenv("JWT_SECRET_KEY"); jwtSecret != "" {
		cfg.JWT.SecretKey = jwtSecret
	}

	// 限流配置
	if ipRate := os.Getenv("RATE_LIMIT_IP_RATE"); ipRate != "" {
		if rate, err := strconv.ParseFloat(ipRate, 64); err == nil {
			cfg.RateLimit.IP.Rate = rate
		}
	}
	if ipBurst := os.Getenv("RATE_LIMIT_IP_BURST"); ipBurst != "" {
		if burst, err := strconv.Atoi(ipBurst); err == nil {
			cfg.RateLimit.IP.Burst = burst
		}
	}
	if apiRate := os.Getenv("RATE_LIMIT_API_RATE"); apiRate != "" {
		if rate, err := strconv.ParseFloat(apiRate, 64); err == nil {
			cfg.RateLimit.API.Rate = rate
		}
	}
	if apiBurst := os.Getenv("RATE_LIMIT_API_BURST"); apiBurst != "" {
		if burst, err := strconv.Atoi(apiBurst); err == nil {
			cfg.RateLimit.API.Burst = burst
		}
	}
	if loginRate := os.Getenv("RATE_LIMIT_LOGIN_RATE"); loginRate != "" {
		if rate, err := strconv.ParseFloat(loginRate, 64); err == nil {
			cfg.RateLimit.Login.Rate = rate
		}
	}
	if loginBurst := os.Getenv("RATE_LIMIT_LOGIN_BURST"); loginBurst != "" {
		if burst, err := strconv.Atoi(loginBurst); err == nil {
			cfg.RateLimit.Login.Burst = burst
		}
	}

	// 存储配置
	if storageType := os.Getenv("STORAGE_TYPE"); storageType != "" {
		cfg.Storage.Type = storageType
	}
	if maxSize := os.Getenv("STORAGE_MAX_SIZE"); maxSize != "" {
		if size, err := strconv.ParseInt(maxSize, 10, 64); err == nil {
			cfg.Storage.MaxSize = size
		}
	}
	if allowedTypes := os.Getenv("STORAGE_ALLOWED_TYPES"); allowedTypes != "" {
		cfg.Storage.AllowedTypes = strings.Split(allowedTypes, ",")
	}
	if thumbnailSize := os.Getenv("STORAGE_THUMBNAIL_SIZE"); thumbnailSize != "" {
		if size, err := strconv.Atoi(thumbnailSize); err == nil {
			cfg.Storage.ThumbnailSize = size
		}
	}

	// 本地存储配置
	if localUploadDir := os.Getenv("STORAGE_LOCAL_UPLOAD_DIR"); localUploadDir != "" {
		cfg.Storage.Local.UploadDir = localUploadDir
	}
	if localURLPrefix := os.Getenv("STORAGE_LOCAL_URL_PREFIX"); localURLPrefix != "" {
		cfg.Storage.Local.URLPrefix = localURLPrefix
	}

	// OSS配置
	if ossEndpoint := os.Getenv("STORAGE_OSS_ENDPOINT"); ossEndpoint != "" {
		cfg.Storage.OSS.Endpoint = ossEndpoint
	}
	if ossAccessKey := os.Getenv("STORAGE_OSS_ACCESS_KEY"); ossAccessKey != "" {
		cfg.Storage.OSS.AccessKey = ossAccessKey
	}
	if ossAccessSecret := os.Getenv("STORAGE_OSS_ACCESS_SECRET"); ossAccessSecret != "" {
		cfg.Storage.OSS.AccessSecret = ossAccessSecret
	}
	if ossBucketName := os.Getenv("STORAGE_OSS_BUCKET_NAME"); ossBucketName != "" {
		cfg.Storage.OSS.BucketName = ossBucketName
	}
	if ossRegion := os.Getenv("STORAGE_OSS_REGION"); ossRegion != "" {
		cfg.Storage.OSS.Region = ossRegion
	}
	if ossRoleArn := os.Getenv("STORAGE_OSS_ROLE_ARN"); ossRoleArn != "" {
		cfg.Storage.OSS.RoleArn = ossRoleArn
	}
	if ossURLPrefix := os.Getenv("STORAGE_OSS_URL_PREFIX"); ossURLPrefix != "" {
		cfg.Storage.OSS.URLPrefix = ossURLPrefix
	}

	// COS配置
	if cosEndpoint := os.Getenv("STORAGE_COS_ENDPOINT"); cosEndpoint != "" {
		cfg.Storage.COS.Endpoint = cosEndpoint
	}
	if cosAccessKey := os.Getenv("STORAGE_COS_ACCESS_KEY"); cosAccessKey != "" {
		cfg.Storage.COS.AccessKey = cosAccessKey
	}
	if cosAccessSecret := os.Getenv("STORAGE_COS_ACCESS_SECRET"); cosAccessSecret != "" {
		cfg.Storage.COS.AccessSecret = cosAccessSecret
	}
	if cosBucketName := os.Getenv("STORAGE_COS_BUCKET_NAME"); cosBucketName != "" {
		cfg.Storage.COS.BucketName = cosBucketName
	}
	if cosRegion := os.Getenv("STORAGE_COS_REGION"); cosRegion != "" {
		cfg.Storage.COS.Region = cosRegion
	}
	if cosURLPrefix := os.Getenv("STORAGE_COS_URL_PREFIX"); cosURLPrefix != "" {
		cfg.Storage.COS.URLPrefix = cosURLPrefix
	}

	// MinIO配置
	if minioEndpoint := os.Getenv("STORAGE_MINIO_ENDPOINT"); minioEndpoint != "" {
		cfg.Storage.MinIO.Endpoint = minioEndpoint
	}
	if minioAccessKey := os.Getenv("STORAGE_MINIO_ACCESS_KEY"); minioAccessKey != "" {
		cfg.Storage.MinIO.AccessKey = minioAccessKey
	}
	if minioAccessSecret := os.Getenv("STORAGE_MINIO_ACCESS_SECRET"); minioAccessSecret != "" {
		cfg.Storage.MinIO.AccessSecret = minioAccessSecret
	}
	if minioBucketName := os.Getenv("STORAGE_MINIO_BUCKET_NAME"); minioBucketName != "" {
		cfg.Storage.MinIO.BucketName = minioBucketName
	}
	if minioRegion := os.Getenv("STORAGE_MINIO_REGION"); minioRegion != "" {
		cfg.Storage.MinIO.Region = minioRegion
	}
	if minioURLPrefix := os.Getenv("STORAGE_MINIO_URL_PREFIX"); minioURLPrefix != "" {
		cfg.Storage.MinIO.URLPrefix = minioURLPrefix
	}
	if minioUseSSL := os.Getenv("STORAGE_MINIO_USE_SSL"); minioUseSSL != "" {
		if useSSL, err := strconv.ParseBool(minioUseSSL); err == nil {
			cfg.Storage.MinIO.UseSSL = useSSL
		}
	}
}

// GetConfig 获取当前配置
func GetConfig() *Config {
	if LoadedConfig.Server.Port == 0 {
		return &DefaultConfig
	}
	return &LoadedConfig
}

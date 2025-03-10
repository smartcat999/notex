package config

import (
	"fmt"
	"net"
	"os"
	"strconv"
	"strings"
	"time"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Server      ServerConfig      `yaml:"server" json:"server"`
	Database    DatabaseConfig    `yaml:"database" json:"database"`
	Email       EmailConfig       `yaml:"email" json:"email"`
	JWT         JWTConfig         `yaml:"jwt" json:"jwt"`
	RateLimit   RateLimitConfig   `yaml:"rate_limit" json:"rate_limit"`
	FileStorage FileStorageConfig `yaml:"file_storage" json:"file_storage"`
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

// FileStorageConfig 文件存储配置
type FileStorageConfig struct {
	UploadDir     string   `yaml:"upload_dir" json:"upload_dir"`         // 上传目录
	AllowedTypes  []string `yaml:"allowed_types" json:"allowed_types"`   // 允许的文件类型
	MaxSize       int64    `yaml:"max_size" json:"max_size"`             // 最大文件大小（字节）
	URLPrefix     string   `yaml:"url_prefix" json:"url_prefix"`         // 文件访问URL前缀
	ThumbnailSize int      `yaml:"thumbnail_size" json:"thumbnail_size"` // 缩略图大小
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
		FileStorage: FileStorageConfig{
			UploadDir:     "uploads",
			AllowedTypes:  []string{"image/jpeg", "image/png", "image/gif", "image/webp"},
			MaxSize:       10 << 20, // 10MB
			URLPrefix:     "/uploads/",
			ThumbnailSize: 300,
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

	// 验证文件存储配置
	if err := c.FileStorage.Validate(); err != nil {
		return fmt.Errorf("file storage config error: %v", err)
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

// Validate 验证文件存储配置
func (c *FileStorageConfig) Validate() error {
	if c.UploadDir == "" {
		return fmt.Errorf("upload_dir cannot be empty")
	}
	if len(c.AllowedTypes) == 0 {
		return fmt.Errorf("allowed_types cannot be empty")
	}
	if c.MaxSize <= 0 {
		return fmt.Errorf("max_size must be positive")
	}
	if c.URLPrefix == "" {
		return fmt.Errorf("url_prefix cannot be empty")
	}
	if c.ThumbnailSize <= 0 {
		return fmt.Errorf("thumbnail_size must be positive")
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

	// 文件存储配置
	if uploadDir := os.Getenv("FILE_UPLOAD_DIR"); uploadDir != "" {
		cfg.FileStorage.UploadDir = uploadDir
	}
	if maxSize := os.Getenv("FILE_MAX_SIZE"); maxSize != "" {
		if size, err := strconv.ParseInt(maxSize, 10, 64); err == nil {
			cfg.FileStorage.MaxSize = size
		}
	}
	if urlPrefix := os.Getenv("FILE_URL_PREFIX"); urlPrefix != "" {
		cfg.FileStorage.URLPrefix = urlPrefix
	}
	if thumbnailSize := os.Getenv("FILE_THUMBNAIL_SIZE"); thumbnailSize != "" {
		if size, err := strconv.Atoi(thumbnailSize); err == nil {
			cfg.FileStorage.ThumbnailSize = size
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

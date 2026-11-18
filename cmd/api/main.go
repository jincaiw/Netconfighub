package main

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"flag"
	"fmt"
	"io/fs"
	"log"
	"net/http"
	"os"
	"os/signal"
	"path"
	"path/filepath"
	"strings"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/netconfighub/netconfighub/internal/baseline"
	"github.com/netconfighub/netconfighub/internal/collector"
	"github.com/netconfighub/netconfighub/internal/config"
	"github.com/netconfighub/netconfighub/internal/gitstore"
	"github.com/netconfighub/netconfighub/internal/handler"
	"github.com/netconfighub/netconfighub/internal/middleware"
	"github.com/netconfighub/netconfighub/internal/model"
	"github.com/netconfighub/netconfighub/internal/notifier"
	"github.com/netconfighub/netconfighub/internal/repository"
	"github.com/netconfighub/netconfighub/internal/sanitizer"
	"github.com/netconfighub/netconfighub/internal/scheduler"
	"github.com/netconfighub/netconfighub/internal/service"
	"github.com/netconfighub/netconfighub/web"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	_ "github.com/glebarez/sqlite"
	gormSqlite "github.com/glebarez/sqlite"
)

var version = "dev"

func main() {
	showVersion := flag.Bool("version", false, "print version and exit")
	flag.Parse()
	if *showVersion {
		fmt.Printf("NetConfigHub %s\n", version)
		return
	}

	configPath := "configs/config.yaml"
	allowMissingConfig := true
	if envPath := os.Getenv("NCH_CONFIG_PATH"); envPath != "" {
		configPath = envPath
		allowMissingConfig = false
	}

	cfg, err := config.LoadConfig(configPath, allowMissingConfig)
	if err != nil {
		log.Fatalf("加载配置文件失败: %v", err)
	}
	if err := config.ValidateProduction(cfg); err != nil {
		log.Fatalf("%v", err)
	}

	log.Printf("NetConfigHub %s 配置加载成功: server=%s, database.driver=%s", version, cfg.Server.Addr(), cfg.Database.Driver)

	db, err := initDatabase(cfg)
	if err != nil {
		log.Fatalf("初始化数据库失败: %v", err)
	}

	if err := autoMigrate(db); err != nil {
		log.Fatalf("数据库迁移失败: %v", err)
	}
	log.Println("数据库迁移完成")

	adminRepo := repository.NewAdminRepository(db)
	tokenRepo := repository.NewAPITokenRepository(db)
	deviceRepo := repository.NewDeviceRepository(db)
	groupRepo := repository.NewGroupRepository(db)
	backupTaskRepo := repository.NewBackupTaskRepository(db)
	baselineRepo := repository.NewBaselineRepository(db)
	deviationRepo := repository.NewDeviationRepository(db)
	alertRepo := repository.NewAlertRepository(db)
	hookRepo := repository.NewHookRepository(db)
	configRepo := repository.NewConfigRepository(db)
	auditLogRepo := repository.NewAuditLogRepository(db)

	gitStore, err := gitstore.NewStore(cfg.Git.RepoPath, cfg.Git.AuthorName, cfg.Git.AuthorEmail, cfg.Git.GPGSignKey, cfg.Git.GitCryptEnabled, cfg.Git.GitCryptKeyFile)
	if err != nil {
		log.Fatalf("初始化 Git 存储失败: %v", err)
	}

	baselineEngine := baseline.NewEngine(gitStore)

	c := collector.NewCollector()

	san := sanitizer.NewSanitizer(cfg.Sanitize)

	ntf := notifier.NewNotifier(cfg.Notify)

	alertService := service.NewAlertService(alertRepo)
	hookService := service.NewHookService(hookRepo)
	configService := service.NewConfigService(configRepo)
	auditService := service.NewAuditService(auditLogRepo)

	sched := scheduler.NewScheduler(scheduler.SchedulerDeps{
		Cfg:            &cfg.Scheduler,
		DeviceRepo:     deviceRepo,
		BackupTaskRepo: backupTaskRepo,
		BaselineRepo:   baselineRepo,
		DeviationRepo:  deviationRepo,
		Collector:      c,
		GitStore:       gitStore,
		BaselineEngine: baselineEngine,
		AlertService:   alertService,
		Notifier:       ntf,
		EncryptionKey:  cfg.Server.EncryptionKey,
		HookService:    hookService,
		AuditService:   auditService,
	})

	adminService := service.NewAdminService(adminRepo)
	tokenService := service.NewAPITokenService(tokenRepo)
	authService := service.NewAuthService(adminService, cfg.Server.JWTSecret, cfg.Server.JWTExpireHours)
	deviceService := service.NewDeviceService(deviceRepo, groupRepo, backupTaskRepo, gitStore, deviationRepo, baselineEngine, san, cfg.Server.EncryptionKey)
	groupService := service.NewGroupService(groupRepo, deviceRepo, backupTaskRepo)
	backupService := service.NewBackupService(backupTaskRepo, deviceRepo, gitStore, san)
	backupService.SetScheduler(sched)
	baselineService := service.NewBaselineService(baselineRepo, deviceRepo, deviationRepo)
	deviationService := service.NewDeviationService(deviationRepo)

	initDefaultAdmin(authService)

	if err := sched.Start(); err != nil {
		log.Fatalf("启动调度器失败: %v", err)
	}
	log.Println("调度器启动完成")

	if strings.EqualFold(strings.TrimSpace(os.Getenv("NCH_ENV")), "production") {
		gin.SetMode(gin.ReleaseMode)
	}
	engine := gin.Default()

	if err := engine.SetTrustedProxies(cfg.Server.TrustedProxies); err != nil {
		log.Fatalf("设置可信代理失败: %v", err)
	}
	engine.Use(securityHeaders(cfg.Server.TLS.Enabled))
	engine.Use(corsMiddleware(cfg.Server.CORS))

	rateLimiter := middleware.NewRateLimiter(5000, time.Minute)
	engine.Use(rateLimiter.Middleware())

	loginRateLimiter := middleware.NewRateLimiter(300, time.Minute)
	router := handler.NewRouter(handler.RouterDeps{
		AdminService:     adminService,
		TokenService:     tokenService,
		AuthService:      authService,
		DeviceService:    deviceService,
		GroupService:     groupService,
		BackupService:    backupService,
		BaselineService:  baselineService,
		DeviationService: deviationService,
		AlertService:     alertService,
		HookService:      hookService,
		ConfigService:    configService,
		AuditService:     auditService,
		TokenAuth:        middleware.APITokenAuth(tokenRepo),
		JWTSecret:        cfg.Server.JWTSecret,
		LoginRateLimiter: loginRateLimiter,
	})
	router.Setup(engine)

	gitStore.Start()
	log.Println("Git 存储启动完成")

	setupFrontend(engine)

	addr := cfg.Server.Addr()
	srv := &http.Server{
		Addr:    addr,
		Handler: engine,
	}

	go func() {
		if cfg.Server.TLS.Enabled {
			log.Printf("NetConfigHub 服务启动（HTTPS），监听地址: %s", addr)
			if err := srv.ListenAndServeTLS(cfg.Server.TLS.CertFile, cfg.Server.TLS.KeyFile); err != nil && err != http.ErrServerClosed {
				log.Fatalf("HTTPS 服务启动失败: %v", err)
			}
		} else {
			log.Printf("NetConfigHub 服务启动（HTTP），监听地址: %s", addr)
			if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
				log.Fatalf("HTTP 服务启动失败: %v", err)
			}
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("正在关闭服务...")

	sched.Stop()
	log.Println("调度器已停止")

	rateLimiter.Stop()
	loginRateLimiter.Stop()
	log.Println("限流器已停止")

	gitStore.Stop()
	log.Println("Git 存储已停止")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("服务关闭失败: %v", err)
	}

	log.Println("服务已关闭")
}

func setupFrontend(engine *gin.Engine) {
	subFS, err := fs.Sub(web.DistFS, "dist")
	if err != nil {
		log.Printf("警告: 无法加载前端静态文件: %v", err)
		return
	}

	engine.NoRoute(func(c *gin.Context) {
		requestPath := c.Request.URL.Path
		if strings.HasPrefix(requestPath, "/api") {
			c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
			return
		}

		cleanPath, ok := normalizeFrontendPath(requestPath)
		if !ok {
			c.JSON(http.StatusBadRequest, gin.H{"error": "bad request"})
			return
		}

		data, err := fs.ReadFile(subFS, cleanPath)
		if err != nil {
			if isPageRoute(cleanPath) {
				data, err = fs.ReadFile(subFS, "index.html")
				if err != nil {
					c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
					return
				}
				cleanPath = "index.html"
			} else {
				c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
				return
			}
		}

		contentType := getContentType(cleanPath)

		c.Writer.Header().Set("Content-Type", contentType)
		c.Writer.WriteHeader(http.StatusOK)
		c.Writer.Write(data)
	})
}

// securityHeaders adds baseline defensive response headers. The CSP keeps the
// existing single-origin Vue/TDesign assets working (including injected styles).
func securityHeaders(tlsEnabled bool) gin.HandlerFunc {
	return func(c *gin.Context) {
		h := c.Writer.Header()
		h.Set("X-Content-Type-Options", "nosniff")
		h.Set("X-Frame-Options", "DENY")
		h.Set("Referrer-Policy", "no-referrer")
		h.Set("Content-Security-Policy", "default-src 'self'; script-src 'self'; style-src 'self' 'unsafe-inline'; img-src 'self' data:; font-src 'self' data:; connect-src 'self'; object-src 'none'; base-uri 'self'; frame-ancestors 'none'")
		if tlsEnabled {
			h.Set("Strict-Transport-Security", "max-age=31536000; includeSubDomains")
		}
		c.Next()
	}
}

func corsMiddleware(corsCfg config.CORSConfig) gin.HandlerFunc {
	allowedOrigins := make(map[string]struct{}, len(corsCfg.AllowedOrigins))
	allowWildcard := false
	for _, origin := range corsCfg.AllowedOrigins {
		origin = strings.TrimSpace(origin)
		if origin == "" {
			continue
		}
		if origin == "*" {
			allowWildcard = true
		}
		allowedOrigins[origin] = struct{}{}
	}
	methods := strings.Join(corsCfg.AllowedMethods, ", ")
	if allowWildcard && !strings.EqualFold(strings.TrimSpace(os.Getenv("NCH_ENV")), "production") {
		log.Println("警告: CORS 允许 *（仅限非生产环境；生产环境请在配置中限制来源）")
	}
	headers := strings.Join(corsCfg.AllowedHeaders, ", ")
	maxAge := fmt.Sprintf("%d", corsCfg.MaxAgeSeconds)

	return func(c *gin.Context) {
		origin := c.Request.Header.Get("Origin")
		switch {
		case allowWildcard:
			c.Header("Access-Control-Allow-Origin", "*")
		case origin != "":
			if _, ok := allowedOrigins[origin]; ok {
				c.Header("Access-Control-Allow-Origin", origin)
				c.Header("Vary", "Origin")
			}
		}
		c.Header("Access-Control-Allow-Methods", methods)
		c.Header("Access-Control-Allow-Headers", headers)
		c.Header("Access-Control-Max-Age", maxAge)
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}
		c.Next()
	}
}

func normalizeFrontendPath(requestPath string) (string, bool) {
	if strings.ContainsRune(requestPath, '\x00') {
		return "", false
	}
	cleanPath := strings.TrimPrefix(path.Clean("/"+requestPath), "/")
	if cleanPath == "" || cleanPath == "." {
		cleanPath = "index.html"
	}
	if !fs.ValidPath(cleanPath) {
		return "", false
	}
	return cleanPath, true
}

func isPageRoute(path string) bool {
	if strings.HasPrefix(path, "@") {
		return false
	}
	ext := filepath.Ext(path)
	if ext == "" || ext == ".html" {
		return true
	}
	return false
}

func getContentType(path string) string {
	switch {
	case strings.HasSuffix(path, ".html"):
		return "text/html; charset=utf-8"
	case strings.HasSuffix(path, ".js"), strings.HasSuffix(path, ".mjs"):
		return "application/javascript; charset=utf-8"
	case strings.HasSuffix(path, ".css"):
		return "text/css; charset=utf-8"
	case strings.HasSuffix(path, ".svg"):
		return "image/svg+xml"
	case strings.HasSuffix(path, ".png"):
		return "image/png"
	case strings.HasSuffix(path, ".ico"):
		return "image/x-icon"
	case strings.HasSuffix(path, ".json"):
		return "application/json"
	case strings.HasSuffix(path, ".woff"):
		return "font/woff"
	case strings.HasSuffix(path, ".woff2"):
		return "font/woff2"
	case strings.HasSuffix(path, ".ttf"):
		return "font/ttf"
	default:
		return "application/octet-stream"
	}
}

func initDatabase(cfg *config.Config) (*gorm.DB, error) {
	var dialector gorm.Dialector

	switch cfg.Database.Driver {
	case "sqlite":
		dir := filepath.Dir(cfg.Database.SQLitePath)
		if err := os.MkdirAll(dir, 0755); err != nil {
			return nil, fmt.Errorf("创建数据库目录失败: %w", err)
		}
		dialector = gormSqlite.Open(cfg.Database.SQLitePath)
	case "mysql":
		if cfg.Database.MySQLDSN == "" {
			return nil, fmt.Errorf("MySQL 驱动需要配置 mysql_dsn")
		}
		dialector = mysql.Open(cfg.Database.MySQLDSN)
	default:
		return nil, fmt.Errorf("不支持的数据库驱动: %s", cfg.Database.Driver)
	}

	db, err := gorm.Open(dialector, &gorm.Config{
		Logger: logger.Default.LogMode(gormLogLevel(cfg.Log.Level)),
	})
	if err != nil {
		return nil, fmt.Errorf("连接数据库失败: %w", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("获取底层数据库连接失败: %w", err)
	}
	sqlDB.SetMaxOpenConns(cfg.Database.MaxOpenConns)
	sqlDB.SetMaxIdleConns(cfg.Database.MaxIdleConns)
	lifetime, err := time.ParseDuration(cfg.Database.ConnMaxLifetime)
	if err != nil || lifetime <= 0 {
		lifetime = 30 * time.Minute
	}
	sqlDB.SetConnMaxLifetime(lifetime)

	return db, nil
}

func gormLogLevel(level string) logger.LogLevel {
	switch strings.ToLower(strings.TrimSpace(level)) {
	case "debug":
		return logger.Info
	case "warn", "warning":
		return logger.Warn
	case "error":
		return logger.Error
	case "silent":
		return logger.Silent
	default:
		return logger.Warn
	}
}

func autoMigrate(db *gorm.DB) error {
	// Default to true so first boot and E2E still create tables. Production is
	// advised to set NCH_AUTO_MIGRATE=false and run a reviewed migration.
	if raw := strings.TrimSpace(os.Getenv("NCH_AUTO_MIGRATE")); strings.EqualFold(raw, "false") || raw == "0" {
		log.Println("数据库自动迁移已跳过（NCH_AUTO_MIGRATE=false）")
		return nil
	}
	log.Println("执行数据库自动迁移 ...")
	if err := db.AutoMigrate(
		&model.Admin{},
		&model.APIToken{},
		&model.Group{},
		&model.Device{},
		&model.BackupTask{},
		&model.Baseline{},
		&model.Deviation{},
		&model.Alert{},
		&model.Hook{},
		&model.SystemConfig{},
		&model.AuditLog{},
	); err != nil {
		return err
	}
	log.Println("数据库自动迁移完成")
	return nil
}

func initDefaultAdmin(authService service.AuthService) {
	username, password, generated := defaultAdminCredentials()
	admin, err := authService.InitAdmin(nil, username, password)
	if err != nil {
		return
	}
	log.Printf("========================================")
	log.Printf("已创建默认管理员账户")
	log.Printf("用户名: %s", username)
	if generated {
		log.Printf("密码: %s（随机生成，仅打印一次，请妥善保存）", password)
		log.Printf("警告: 未设置 NCH_ADMIN_PASSWORD，已使用随机口令；请立即修改！")
	} else {
		log.Printf("密码: 已通过 NCH_ADMIN_PASSWORD 设置")
	}
	log.Printf("管理员 ID: %d", admin.ID)
	log.Printf("========================================")
}

// defaultAdminCredentials returns the bootstrap admin credentials. When no
// password is supplied, a strong random password is generated and printed once;
// the literal "admin" is never used as a fallback.
func defaultAdminCredentials() (username, password string, generated bool) {
	username = strings.TrimSpace(os.Getenv("NCH_ADMIN_USERNAME"))
	password = os.Getenv("NCH_ADMIN_PASSWORD")
	if username == "" {
		username = "admin"
	}
	if password == "" {
		buf := make([]byte, 24)
		if _, err := rand.Read(buf); err != nil {
			password = fmt.Sprintf("nch-%d-%d", os.Getpid(), time.Now().UnixNano())
		} else {
			password = hex.EncodeToString(buf)
		}
		generated = true
	}
	return username, password, generated
}

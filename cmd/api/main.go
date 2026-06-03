package main

import (
	"context"
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

func main() {
	configPath := "configs/config.yaml"
	if envPath := os.Getenv("NCH_CONFIG_PATH"); envPath != "" {
		configPath = envPath
	}

	cfg, err := config.LoadConfig(configPath)
	if err != nil {
		log.Fatalf("加载配置文件失败: %v", err)
	}

	log.Printf("配置加载成功: server=%s, database.driver=%s", cfg.Server.Addr(), cfg.Database.Driver)

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

	sched := scheduler.NewScheduler(db, &cfg.Scheduler, deviceRepo, backupTaskRepo, baselineRepo, deviationRepo, c, gitStore, baselineEngine, alertService, ntf, cfg.Server.EncryptionKey, hookService, auditService)

	adminService := service.NewAdminService(adminRepo)
	tokenService := service.NewAPITokenService(tokenRepo)
	authService := service.NewAuthService(adminService, cfg.Server.JWTSecret, cfg.Server.JWTExpireHours)
	deviceService := service.NewDeviceService(deviceRepo, groupRepo, backupTaskRepo, gitStore, deviationRepo, baselineEngine, san, cfg.Server.EncryptionKey)
	groupService := service.NewGroupService(groupRepo, deviceRepo, backupTaskRepo)
	backupService := service.NewBackupService(backupTaskRepo, deviceRepo, gitStore, san)
	baselineService := service.NewBaselineService(baselineRepo, deviceRepo, deviationRepo)
	deviationService := service.NewDeviationService(deviationRepo)

	initDefaultAdmin(authService)

	if err := sched.Start(); err != nil {
		log.Fatalf("启动调度器失败: %v", err)
	}
	log.Println("调度器启动完成")

	engine := gin.Default()

	engine.Use(corsMiddleware(cfg.Server.CORS))

	rateLimiter := middleware.NewRateLimiter(5000, time.Minute)
	engine.Use(rateLimiter.Middleware())

	loginRateLimiter := middleware.NewRateLimiter(300, time.Minute)
	router := handler.NewRouter(adminService, tokenService, authService, deviceService, groupService, backupService, baselineService, deviationService, alertService, hookService, configService, auditService, sched, cfg.Server.JWTSecret, tokenRepo, loginRateLimiter)
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
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		return nil, fmt.Errorf("连接数据库失败: %w", err)
	}

	return db, nil
}

func autoMigrate(db *gorm.DB) error {
	return db.AutoMigrate(
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
	)
}

func initDefaultAdmin(authService service.AuthService) {
	username, password, usingDefault := defaultAdminCredentials()
	admin, err := authService.InitAdmin(nil, username, password)
	if err != nil {
		return
	}
	log.Printf("========================================")
	log.Printf("已创建默认管理员账户")
	log.Printf("用户名: %s", username)
	if usingDefault {
		log.Printf("密码: admin")
		log.Printf("警告: 正在使用内置默认管理员密码，请立即修改！")
	} else {
		log.Printf("密码: 已通过 NCH_ADMIN_PASSWORD 设置")
	}
	log.Printf("管理员 ID: %d", admin.ID)
	log.Printf("========================================")
}

func defaultAdminCredentials() (username, password string, usingDefault bool) {
	username = strings.TrimSpace(os.Getenv("NCH_ADMIN_USERNAME"))
	password = os.Getenv("NCH_ADMIN_PASSWORD")
	if username == "" {
		username = "admin"
	}
	if password == "" {
		password = "admin"
		usingDefault = true
	}
	return username, password, usingDefault
}

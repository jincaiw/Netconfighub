package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/netconfighub/netconfighub/internal/middleware"
	"github.com/netconfighub/netconfighub/internal/repository"
	"github.com/netconfighub/netconfighub/internal/service"
)

type BackupScheduler interface {
	SubmitDevice(deviceID uint) error
	SubmitGroup(groupID uint) error
}

type Router struct {
	adminService      service.AdminService
	tokenService      service.APITokenService
	authService       service.AuthService
	deviceService     service.DeviceService
	groupService      service.GroupService
	backupService     service.BackupService
	baselineService   service.BaselineService
	deviationService  service.DeviationService
	alertService      service.AlertService
	hookService       service.HookService
	configService     service.ConfigService
	auditService      service.AuditService
	scheduler         BackupScheduler
	jwtSecret         string
	tokenRepo         repository.APITokenRepository
	loginRateLimiter  *middleware.RateLimiter
}

func NewRouter(
	adminService service.AdminService,
	tokenService service.APITokenService,
	authService service.AuthService,
	deviceService service.DeviceService,
	groupService service.GroupService,
	backupService service.BackupService,
	baselineService service.BaselineService,
	deviationService service.DeviationService,
	alertService service.AlertService,
	hookService service.HookService,
	configService service.ConfigService,
	auditService service.AuditService,
	scheduler BackupScheduler,
	jwtSecret string,
	tokenRepo repository.APITokenRepository,
	loginRateLimiter *middleware.RateLimiter,
) *Router {
	return &Router{
		adminService:      adminService,
		tokenService:      tokenService,
		authService:       authService,
		deviceService:     deviceService,
		groupService:      groupService,
		backupService:     backupService,
		baselineService:   baselineService,
		deviationService:  deviationService,
		alertService:      alertService,
		hookService:       hookService,
		configService:     configService,
		auditService:      auditService,
		scheduler:         scheduler,
		jwtSecret:         jwtSecret,
		tokenRepo:         tokenRepo,
		loginRateLimiter:  loginRateLimiter,
	}
}

func (r *Router) Setup(engine *gin.Engine) {
	api := engine.Group("/api/v1")
	{
		api.GET("/health", r.healthCheck)

		auth := api.Group("/auth")
		{
			auth.POST("/login", r.loginRateLimiter.Middleware(), r.login)
			auth.POST("/logout", r.logout)
			auth.POST("/refresh", r.refreshToken)
			auth.PUT("/change-password", middleware.JWTAuth(r.jwtSecret), r.changePassword)
		}

		authenticated := api.Group("")
		authenticated.Use(middleware.JWTAuth(r.jwtSecret))
		{
			admins := authenticated.Group("/admins")
			{
				admins.GET("", r.listAdmins)
				admins.POST("", r.createAdmin)
				admins.GET("/:id", r.getAdmin)
				admins.PUT("/:id", r.updateAdmin)
				admins.DELETE("/:id", r.deleteAdmin)
			}

			tokens := authenticated.Group("/tokens")
			{
				tokens.GET("", r.listTokens)
				tokens.POST("", r.createToken)
				tokens.DELETE("/:id", r.deleteToken)
			}

			devices := authenticated.Group("/devices")
			{
				devices.GET("", r.listDevices)
				devices.POST("", r.createDevice)
				devices.GET("/:id", r.getDevice)
				devices.PUT("/:id", r.updateDevice)
				devices.DELETE("/:id", r.deleteDevice)
				devices.PUT("/:id/enable", r.enableDevice)
				devices.PUT("/:id/disable", r.disableDevice)
				devices.POST("/:id/backup", r.triggerBackup)
				devices.GET("/:id/config", r.getDeviceConfig)
				devices.GET("/:id/versions", r.getDeviceVersions)
				devices.GET("/:id/deviations", r.getDeviceDeviations)
				devices.GET("/:id/diff", r.getConfigDiff)
				devices.GET("/:id/diff-latest", r.getLatestDiff)
				devices.GET("/:id/diff-download", r.downloadDiff)
				devices.POST("/import", r.importDevices)
				devices.GET("/export", r.exportDevices)
			}

			groups := authenticated.Group("/groups")
			{
				groups.GET("", r.listGroups)
				groups.POST("", r.createGroup)
				groups.GET("/:id", r.getGroup)
				groups.PUT("/:id", r.updateGroup)
				groups.DELETE("/:id", r.deleteGroup)
				groups.POST("/:id/backup", r.triggerGroupBackup)
			}

			backups := authenticated.Group("/backups")
			{
				backups.GET("", r.listBackups)
				backups.GET("/:id", r.getBackup)
				backups.GET("/:id/config", r.getBackupConfig)
			}

			jobs := authenticated.Group("/jobs")
			{
				jobs.GET("", r.listBackups)
				jobs.GET("/:id", r.getBackup)
			}

			baselines := authenticated.Group("/baselines")
			{
				baselines.GET("", r.listBaselines)
				baselines.POST("", r.createBaseline)
				baselines.GET("/:id", r.getBaseline)
				baselines.PUT("/:id", r.updateBaseline)
				baselines.DELETE("/:id", r.deleteBaseline)
			}

			deviations := authenticated.Group("/deviations")
			{
				deviations.GET("", r.listDeviations)
				deviations.GET("/:id", r.getDeviation)
				deviations.GET("/export", r.exportDeviations)
			}

			alerts := authenticated.Group("/alerts")
			{
				alerts.GET("", r.listAlerts)
				alerts.GET("/unread-count", r.countUnreadAlerts)
				alerts.PUT("/:id/read", r.markAlertAsRead)
				alerts.PUT("/read-all", r.markAllAlertsAsRead)
				alerts.DELETE("/:id", r.deleteAlert)
			}

			hooks := authenticated.Group("/hooks")
			{
				hooks.GET("", r.listHooks)
				hooks.POST("", r.createHook)
				hooks.GET("/:id", r.getHook)
				hooks.PUT("/:id", r.updateHook)
				hooks.DELETE("/:id", r.deleteHook)
			}

			configs := authenticated.Group("/configs")
			{
				configs.GET("", r.listConfigs)
				configs.GET("/:key", r.getConfig)
				configs.PUT("", r.setConfig)
			}

			authenticated.GET("/failed-devices", r.listFailedDevices)
			authenticated.GET("/audit-logs", r.listAuditLogs)
		}

		tokenAuth := api.Group("/ext")
		tokenAuth.Use(middleware.APITokenAuth(r.tokenRepo))
		{
			tokenAuth.GET("/devices", r.listDevices)
			tokenAuth.POST("/devices/:id/backup", r.triggerBackup)
			tokenAuth.GET("/devices/:id/config", r.getDeviceConfig)
			tokenAuth.GET("/devices/:id/versions", r.getDeviceVersions)
			tokenAuth.GET("/devices/:id/deviations", r.getDeviceDeviations)
			tokenAuth.GET("/backups", r.listBackups)
			tokenAuth.GET("/backups/:id/config", r.getBackupConfig)
			tokenAuth.GET("/baselines", r.listBaselines)
			tokenAuth.GET("/baselines/:id", r.getBaseline)
			tokenAuth.GET("/deviations", r.listDeviations)
			tokenAuth.GET("/deviations/:id", r.getDeviation)
			tokenAuth.GET("/alerts", r.listAlerts)
			tokenAuth.GET("/alerts/unread-count", r.countUnreadAlerts)
			tokenAuth.GET("/failed-devices", r.listFailedDevices)
		}
	}
}

func (r *Router) healthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status":  "ok",
		"service": "NetConfigHub",
	})
}

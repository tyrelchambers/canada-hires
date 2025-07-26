package container

import (
	"canada-hires/controllers"
	"canada-hires/db"
	"canada-hires/middleware"
	"canada-hires/repos"
	"canada-hires/services"
	"net/http"

	"go.uber.org/dig"
)

// registerProviders registers all provider functions with the DI container
func registerProviders(c *dig.Container) error {
	// Database providers
	if err := c.Provide(NewDatabaseConfig); err != nil {
		return err
	}

	if err := c.Provide(NewDatabase); err != nil {
		return err
	}

	// Repository providers
	if err := c.Provide(NewUserRepository); err != nil {
		return err
	}
	if err := c.Provide(NewLoginTokenRepository); err != nil {
		return err
	}
	if err := c.Provide(NewBusinessRepository); err != nil {
		return err
	}
	if err := c.Provide(NewReportRepository); err != nil {
		return err
	}

	if err := c.Provide(NewSessionRepository); err != nil {
		return err
	}

	if err := c.Provide(NewLMIARepository); err != nil {
		return err
	}

	if err := c.Provide(NewJobBankRepository); err != nil {
		return err
	}

	// Service providers
	if err := c.Provide(NewEmailService); err != nil {
		return err
	}
	if err := c.Provide(NewAuthService); err != nil {
		return err
	}
	if err := c.Provide(NewBusinessService); err != nil {
		return err
	}
	if err := c.Provide(NewReportService); err != nil {
		return err
	}
	if err := c.Provide(NewUserService); err != nil {
		return err
	}

	if err := c.Provide(NewLMIAService); err != nil {
		return err
	}

	if err := c.Provide(NewCronService); err != nil {
		return err
	}

	if err := c.Provide(NewJobBankService); err != nil {
		return err
	}

	if err := c.Provide(NewJobBankBrowserService); err != nil {
		return err
	}

	if err := c.Provide(NewJobBankConcurrentService); err != nil {
		return err
	}

	if err := c.Provide(NewJobService); err != nil {
		return err
	}

	// Controller providers
	if err := c.Provide(NewAuthController); err != nil {
		return err
	}
	if err := c.Provide(NewBusinessController); err != nil {
		return err
	}
	if err := c.Provide(NewReportController); err != nil {
		return err
	}
	if err := c.Provide(NewUserController); err != nil {
		return err
	}

	if err := c.Provide(NewLMIAController); err != nil {
		return err
	}

	if err := c.Provide(NewJobController); err != nil {
		return err
	}

	// Middleware providers
	if err := c.Provide(NewAuthMiddleware); err != nil {
		return err
	}

	return nil
}

// NewDatabaseConfig creates database configuration from environment
func NewDatabaseConfig() db.Config {
	return db.NewConfigFromEnv()
}

// NewDatabase creates a new database connection
func NewDatabase(config db.Config) (db.Database, error) {
	return db.NewPostgresDB(config)
}

// NewUserRepository creates a new user repository
func NewUserRepository(database db.Database) repos.UserRepository {
	return repos.NewUserRepository(database.GetDB())
}

// NewLoginTokenRepository creates a new login token repository
func NewLoginTokenRepository(database db.Database) repos.LoginTokenRepository {
	return repos.NewLoginTokenRepository(database.GetDB())
}

// NewBusinessRepository creates a new business repository
func NewBusinessRepository(database db.Database) repos.BusinessRepository {
	return repos.NewBusinessRepository(database.GetDB())
}

// NewReportRepository creates a new report repository
func NewReportRepository(database db.Database) repos.ReportRepository {
	return repos.NewReportRepository(database.GetDB())
}

// NewSessionRepository creates a new session repository
func NewSessionRepository(database db.Database) repos.SessionRepository {
	return repos.NewSessionRepository(database.GetDB())
}

// NewEmailService creates a new email service
func NewEmailService() services.EmailService {
	return services.NewEmailService()
}

// NewAuthService creates a new auth service
func NewAuthService(userRepo repos.UserRepository, tokenRepo repos.LoginTokenRepository, sessionRepo repos.SessionRepository, emailService services.EmailService) services.AuthService {
	return services.NewAuthService(userRepo, tokenRepo, sessionRepo, emailService)
}

// NewBusinessService creates a new business service
func NewBusinessService(repo repos.BusinessRepository) services.BusinessService {
	return services.NewBusinessService(repo)
}

// NewReportService creates a new report service
func NewReportService(repo repos.ReportRepository) services.ReportService {
	return services.NewReportService(repo)
}

// NewAuthController creates a new auth controller
func NewAuthController(authService services.AuthService, userService services.UserService) controllers.AuthController {
	return controllers.NewAuthController(authService, userService)
}

// NewBusinessController creates a new business controller
func NewBusinessController(service services.BusinessService) controllers.BusinessController {
	return controllers.NewBusinessController(service)
}

// NewReportController creates a new report controller
func NewReportController(service services.ReportService) controllers.ReportController {
	return controllers.NewReportController(service)
}

func NewUserService(userRepo repos.UserRepository) services.UserService {
	return services.NewUserService(userRepo)
}

func NewUserController(userService services.UserService) controllers.UserController {
	return controllers.NewUserController(userService)
}

// NewAuthMiddleware creates a new auth middleware
func NewAuthMiddleware(authService services.AuthService, userService services.UserService) func(http.Handler) http.Handler {
	authMW := middleware.NewAuthMiddleware(authService, userService)
	return authMW.Middleware
}

// NewLMIARepository creates a new LMIA repository
func NewLMIARepository(database db.Database) repos.LMIARepository {
	return repos.NewLMIARepository(database.GetDB())
}

// NewLMIAService creates a new LMIA service
func NewLMIAService(repo repos.LMIARepository) services.LMIAService {
	return services.NewLMIAService(repo)
}

// NewCronService creates a new cron service
func NewCronService(lmiaService services.LMIAService, repo repos.LMIARepository) services.CronService {
	return services.NewCronService(lmiaService, repo)
}

// NewLMIAController creates a new LMIA controller
func NewLMIAController(lmiaService services.LMIAService, cronService services.CronService, repo repos.LMIARepository) *controllers.LMIAController {
	return controllers.NewLMIAController(lmiaService, cronService, repo)
}

// NewJobBankRepository creates a new Job Bank repository
func NewJobBankRepository(database db.Database) repos.JobBankRepository {
	return repos.NewJobBankRepository(database.GetDB())
}

// NewJobBankService creates a new Job Bank service
func NewJobBankService(repo repos.JobBankRepository) services.JobBankService {
	return services.NewJobBankService(repo)
}

// NewJobBankBrowserService creates a new Job Bank browser service  
func NewJobBankBrowserService(repo repos.JobBankRepository) services.JobBankBrowserService {
	return services.NewJobBankBrowserService(repo)
}

// NewJobBankConcurrentService creates a new Job Bank concurrent service
func NewJobBankConcurrentService(repo repos.JobBankRepository) services.JobBankConcurrentService {
	return services.NewJobBankConcurrentService(repo)
}

// NewJobService creates a new Job service
func NewJobService(repo repos.JobBankRepository) services.JobService {
	return services.NewJobService(repo)
}

// NewJobController creates a new Job controller
func NewJobController(repo repos.JobBankRepository, jobService services.JobService) *controllers.JobController {
	return controllers.NewJobController(repo, jobService)
}

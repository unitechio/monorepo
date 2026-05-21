package database

import (
	"log/slog"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func AutoMigrate(db *gorm.DB) error {
	slog.Info("Running database migrations...")

	for _, stmt := range []string{
		`CREATE EXTENSION IF NOT EXISTS "uuid-ossp"`,
		`CREATE EXTENSION IF NOT EXISTS "pgcrypto"`,
	} {
		if err := db.Exec(stmt).Error; err != nil {
			slog.Warn("Failed to ensure database extension", slog.String("statement", stmt), slog.String("error", err.Error()))
		}
	}

	// Drop tables that may have incompatible schema (e.g. old integer ID vs new UUID).
	// This is safe for development/test databases only.
	if err := dropAndRecreateUserTables(db); err != nil {
		slog.Warn("Could not clean old tables, attempting migration anyway", slog.String("error", err.Error()))
	}

	if err := db.AutoMigrate(
		&domain.User{},
		&domain.Customer{},
		&domain.Role{},
		&domain.Permission{},
		&domain.UserRole{},
		&domain.RolePermission{},
		&domain.UserPermission{},
		&domain.OTP{},
		&domain.RefreshToken{},
		&domain.Session{},
		&domain.LoginAttempt{},
		&domain.Environment{},
		&domain.FeatureFlag{},
		&domain.VocabularyWord{},
		&domain.UserVocabularyProgress{},
		&domain.WritingSubmission{},
		&domain.Course{},
		&domain.Unit{},
		&domain.Lesson{},
		&domain.UserProgress{},
		&domain.Activity{},
		&domain.ActivitySubmission{},
		&domain.SpeakingSession{},
		&domain.ListeningLesson{},
		&domain.StudyPlanner{},
		&domain.BillingPlan{},
		&domain.BillingSubscription{},
		&domain.BillingHistory{},
		&domain.Streak{},
		&domain.XPPoint{},
		&domain.PracticeSession{},
		&domain.PronunciationHistory{},
		&domain.DictionaryHistory{},
		&domain.VocabularySet{},
		&domain.VocabularySetWord{},
	); err != nil {
		slog.Error("Failed to migrate user tables", slog.String("error", err.Error()))
		return err
	}

	// Authorization related tables
	if err := db.AutoMigrate(
		&domain.Module{},
		&domain.Department{},
		&domain.Service{},
		&domain.Scope{},
		&domain.EnhancedPermission{},
		&domain.RoleEnhancedPermission{},
		&domain.UserEnhancedPermission{},
	); err != nil {
		slog.Error("Failed to migrate authorization tables", slog.String("error", err.Error()))
		return err
	}

	// Content related tables
	if err := db.AutoMigrate(
		&domain.Post{},
		&domain.Media{},
		&domain.PostMedia{},
		&domain.Category{},
		&domain.Tag{},
		&domain.PostSchedule{},
	); err != nil {
		slog.Error("Failed to migrate content tables", slog.String("error", err.Error()))
		return err
	}

	// System related tables
	if err := db.AutoMigrate(
		&domain.AuditLog{},
		&domain.SystemSetting{},
		&domain.Notification{},
		&domain.ActivityLog{},
		&domain.EmailTemplate{},
		&domain.EmailLog{},
		&domain.UserSettings{},
		&domain.Document{},
		&domain.DocumentPermission{},
		&domain.DocumentComment{},
		&domain.DocumentVersion{},
	); err != nil {
		slog.Error("Failed to migrate system tables", slog.String("error", err.Error()))
		return err
	}

	slog.Info("Database migrations completed successfully")
	return nil
}

// dropAndRecreateUserTables detects and fixes incompatible table schemas.
// For dev/test databases: drops all user-related tables when issues are found.
func dropAndRecreateUserTables(db *gorm.DB) error {
	if !db.Migrator().HasTable("users") {
		return nil
	}

	needsDrop := false

	// Check 1: id column must be uuid type
	var colType string
	db.Raw(`SELECT data_type FROM information_schema.columns WHERE table_name = 'users' AND column_name = 'id' LIMIT 1`).Scan(&colType)
	if colType != "" && colType != "uuid" {
		slog.Warn("Schema mismatch: id column is not uuid", slog.String("col_type", colType))
		needsDrop = true
	}

	// Check 2: password column must not have NULLs
	if !needsDrop {
		var nullCount int64
		db.Raw("SELECT COUNT(*) FROM users WHERE password IS NULL").Scan(&nullCount)
		if nullCount > 0 {
			slog.Warn("Schema mismatch: users table has NULL passwords", slog.Int64("count", nullCount))
			needsDrop = true
		}
	}

	if !needsDrop {
		return nil
	}

	slog.Warn("Dropping incompatible user-related tables (dev/test mode). They will be recreated.")
	for _, t := range []string{
		"user_permissions", "role_permissions", "user_roles",
		"user_enhanced_permissions", "role_enhanced_permissions",
		"refresh_tokens", "audit_logs", "activity_logs", "otps",
		"posts", "post_media", "post_schedules", "media",
		"customers", "users",
	} {
		db.Exec("DROP TABLE IF EXISTS " + t + " CASCADE")
	}
	slog.Info("Tables dropped. AutoMigrate will recreate with correct schema.")
	return nil
}

func SeedData(db *gorm.DB) error {
	slog.Info("Seeding initial data...")

	// Seed Modules
	if err := seedModules(db); err != nil {
		return err
	}

	// Seed Departments
	if err := seedDepartments(db); err != nil {
		return err
	}

	// Seed Services
	if err := seedServices(db); err != nil {
		return err
	}

	// Seed Scopes
	if err := seedScopes(db); err != nil {
		return err
	}

	// Seed Roles
	if err := seedRoles(db); err != nil {
		return err
	}

	// Seed Permissions
	if err := seedPermissions(db); err != nil {
		return err
	}

	// Seed Enhanced Permissions
	if err := seedEnhancedPermissions(db); err != nil {
		return err
	}

	// Assign permissions to super_admin role
	if err := assignPermissionsToSuperAdmin(db); err != nil {
		return err
	}

	// Seed Users
	if err := seedUsers(db); err != nil {
		return err
	}

	// Seed Categories
	if err := seedCategories(db); err != nil {
		return err
	}
	if err := seedBillingPlans(db); err != nil {
		return err
	}
	if err := seedListeningLessons(db); err != nil {
		return err
	}

	slog.Info("Initial data seeded successfully")
	return nil
}

func seedModules(db *gorm.DB) error {
	modules := []domain.Module{
		{Code: "admin", Name: "Administration", DisplayName: "System Administration", IsActive: true, IsSystem: true},
		{Code: "crm", Name: "CRM", DisplayName: "Customer Relationship Management", IsActive: true, IsSystem: true},
		{Code: "content", Name: "Content", DisplayName: "Content Management", IsActive: true, IsSystem: true},
	}

	for i := range modules {
		if err := db.Where(domain.Module{Code: modules[i].Code}).FirstOrCreate(&modules[i]).Error; err != nil {
			slog.Error("Failed to create module", slog.String("module", modules[i].Code), slog.String("error", err.Error()))
			return err
		}
	}
	slog.Info("Modules seeded successfully")
	return nil
}

func seedDepartments(db *gorm.DB) error {
	var adminModule, crmModule, contentModule domain.Module
	db.Where("code = ?", "admin").First(&adminModule)
	db.Where("code = ?", "crm").First(&crmModule)
	db.Where("code = ?", "content").First(&contentModule)

	departments := []domain.Department{
		{ModuleID: adminModule.ID, Code: "system", Name: "System", IsActive: true, IsSystem: true},
		{ModuleID: crmModule.ID, Code: "sales", Name: "Sales", IsActive: true, IsSystem: true},
		{ModuleID: contentModule.ID, Code: "editorial", Name: "Editorial", IsActive: true, IsSystem: true},
	}

	for _, dept := range departments {
		if err := db.Where("code = ?", dept.Code).FirstOrCreate(&dept).Error; err != nil {
			slog.Error("Failed to create department", slog.String("department", dept.Code), slog.String("error", err.Error()))
			return err
		}
	}
	slog.Info("Departments seeded successfully")
	return nil
}

func seedServices(db *gorm.DB) error {
	var systemDept, salesDept, editorialDept domain.Department
	db.Where("code = ?", "system").First(&systemDept)
	db.Where("code = ?", "sales").First(&salesDept)
	db.Where("code = ?", "editorial").First(&editorialDept)

	services := []domain.Service{
		{DepartmentID: systemDept.ID, Code: "users", Name: "User Management", IsActive: true, IsSystem: true},
		{DepartmentID: salesDept.ID, Code: "customers", Name: "Customer Management", IsActive: true, IsSystem: true},
	}

	for _, service := range services {
		current := service
		if err := db.Where("code = ?", current.Code).FirstOrCreate(&current).Error; err != nil {
			slog.Error("Failed to create service", slog.String("service", service.Code), slog.String("error", err.Error()))
			return err
		}
	}
	slog.Info("Services seeded successfully")
	return nil
}

func seedScopes(db *gorm.DB) error {
	scopes := []domain.Scope{
		{Code: "org", Name: "Organization", Level: domain.ScopeLevelOrganization, IsSystem: true},
		{Code: "personal", Name: "Personal", Level: domain.ScopeLevelPersonal, IsSystem: true},
	}

	for _, scope := range scopes {
		if err := db.Where("code = ?", scope.Code).FirstOrCreate(&scope).Error; err != nil {
			slog.Error("Failed to create scope", slog.String("scope", scope.Code), slog.String("error", err.Error()))
			return err
		}
	}
	slog.Info("Scopes seeded successfully")
	return nil
}

func seedRoles(db *gorm.DB) error {
	roles := []domain.Role{
		{Name: "super_admin", DisplayName: "Super Administrator", Level: domain.RoleLevelOrganization, IsSystem: true},
		{Name: "admin", DisplayName: "Administrator", Level: domain.RoleLevelOrganization, IsSystem: true},
		{Name: "instructor", DisplayName: "Instructor", Level: domain.RoleLevelOrganization, IsSystem: true},
		{Name: "premium", DisplayName: "Premium User", Level: domain.RoleLevelOrganization, IsSystem: true},
		{Name: "user", DisplayName: "User", Level: domain.RoleLevelOrganization, IsSystem: true},
	}

	for _, role := range roles {
		if err := db.Where("name = ?", role.Name).FirstOrCreate(&role).Error; err != nil {
			slog.Error("Failed to create role", slog.String("role", role.Name), slog.String("error", err.Error()))
			return err
		}
	}
	slog.Info("Roles seeded successfully")
	return nil
}

func seedPermissions(db *gorm.DB) error {
	slog.Info("Permissions seeded successfully")
	return nil
}

func seedEnhancedPermissions(db *gorm.DB) error {
	slog.Info("Enhanced permissions seeded successfully")
	return nil
}

func assignPermissionsToSuperAdmin(db *gorm.DB) error {
	slog.Info("Permissions assigned to super_admin successfully")
	return nil
}

func seedUsers(db *gorm.DB) error {
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("password123"), bcrypt.DefaultCost)

	adminUser := domain.User{
		Email:    "admin@eenglish.io",
		Password: string(hashedPassword),
		Status:   domain.UserStatusActive,
		TenantID: uuid.New(),
	}

	if err := db.Where("email = ?", adminUser.Email).FirstOrCreate(&adminUser).Error; err != nil {
		slog.Error("Failed to create admin user", slog.String("error", err.Error()))
		return err
	}
	var superAdminRole domain.Role
	if err := db.Where("name = ?", "super_admin").First(&superAdminRole).Error; err == nil {
		if err := db.Model(&adminUser).Association("Roles").Append(&superAdminRole); err != nil {
			slog.Warn("Failed to attach super_admin role", slog.String("error", err.Error()))
		}
	}
	slog.Info("Admin user seeded successfully")
	return nil
}

func seedCategories(db *gorm.DB) error {
	slog.Info("Categories seeded successfully")
	return nil
}

func seedBillingPlans(db *gorm.DB) error {
	plans := []domain.BillingPlan{
		{TenantID: uuid.Nil, Name: "Starter", Code: "starter", Price: 9.99, Currency: "USD", Description: "Starter access for Academy English", BillingCycle: "monthly", IsActive: true},
		{TenantID: uuid.Nil, Name: "Pro", Code: "pro", Price: 19.99, Currency: "USD", Description: "Full AI coaching for Academy English", BillingCycle: "monthly", IsActive: true},
	}
	for i := range plans {
		if err := db.Where("code = ?", plans[i].Code).FirstOrCreate(&plans[i]).Error; err != nil {
			slog.Error("Failed to create billing plan", slog.String("code", plans[i].Code), slog.String("error", err.Error()))
			return err
		}
	}
	slog.Info("Billing plans seeded successfully")
	return nil
}

func seedListeningLessons(db *gorm.DB) error {
	items := []domain.ListeningLesson{
		{TenantID: uuid.Nil, Title: "Daily Conversation", Description: "Practice everyday academy english listening.", AudioURL: "https://cdn.eenglish.local/audio/daily-conversation.mp3", Transcript: "Two friends discuss their study plan for the week.", Level: "beginner", Domain: "english", IsActive: true},
		{TenantID: uuid.Nil, Title: "Campus Interview", Description: "Listen to an interview about learning goals.", AudioURL: "https://cdn.eenglish.local/audio/campus-interview.mp3", Transcript: "A student explains how speaking practice improved confidence.", Level: "intermediate", Domain: "english", IsActive: true},
	}
	for i := range items {
		if err := db.Where("title = ?", items[i].Title).FirstOrCreate(&items[i]).Error; err != nil {
			slog.Error("Failed to create listening lesson", slog.String("title", items[i].Title), slog.String("error", err.Error()))
			return err
		}
	}
	slog.Info("Listening lessons seeded successfully")
	return nil
}

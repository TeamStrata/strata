package main

import (
	"context"
	"log"
	"time"

	_ "github.com/TeamStrata/strata/docs"
	"github.com/TeamStrata/strata/pkg/api"
	"github.com/TeamStrata/strata/pkg/database"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func main() {
	server := gin.Default()

	//cors config (redux)
	config := cors.Config{
		AllowOrigins:     []string{"http://localhost:5173"},
		AllowMethods:     []string{"GET", "PATCH", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}
	server.Use(cors.New(config))

	// Get database connection string
	conStr, err := database.GetConnectionString(database.DefaultPath)
	if err != nil {
		log.Fatalf("error loading .env file: %s", err.Error())
	}

	// Initialize database manager
	db, err := database.NewDbManager(conStr, context.Background())
	if err != nil {
		log.Fatalf("error initializing DB manager: %s", err.Error())
	}

	// Create the client database which will be used for custom queries
	cdb_str, cerr := db.GetSetting("cdb")
	if cerr != nil {
		log.Printf("error reading DB settings: %s", cerr.Error())
		return
	}

	cdb, cerr := database.NewDbManager(cdb_str, context.Background())
	if cerr != nil {
		log.Printf("error initializing client DB manager: %s", cerr.Error())
		return
	}

	// Map for active users and uuids
	// key: uuid
	// value: username
	activeUsers := make(map[string]api.UserSessionData)

	// Frontend
	server.Static("/assets", "ui/dist/assets")
	server.StaticFile("/", "ui/dist/index.html")
	server.StaticFile("/favicon.ico", "ui/dist/favicon.ico")

	// Auth endpoints
	server.POST("/login", api.LoginHandler(db, activeUsers))

	protected := server.Group("/api")
	{
		protected.Use(api.AuthHandler(activeUsers))

		// Auth routes
		protected.POST("/logout", api.LogoutHandler(activeUsers))
		protected.POST("/auth", api.AuthHandler(activeUsers))
		protected.GET("/isadmin/:uname", api.IsAdminHandler(activeUsers))

		// Dashboard Components
		/// Charts
		protected.GET("/charts", api.GetChartListHandler(db))
		protected.GET("/chart/:cid", api.GetChartHandler(db))
		protected.POST("/chart", api.CreateChartHandler(db))
		protected.DELETE("/chart/:cid", api.DeleteChartHandler(db))
		protected.PATCH("/chart/:cid", api.UpdateChartHandler(db))
		/// Chart Series
		protected.GET("/chart/:cid/series", api.GetChartSeriesListHandler(db))
		protected.GET("/chart/:cid/series/:sid", api.GetChartSingleSeriesHandler(db))
		protected.POST("/chart/:cid/series", api.AddChartSeriesHandler(db))
		protected.DELETE("/chart/:cid/series/:sid", api.DeleteChartSingleSeriesHandler(db))
		protected.PATCH("/chart/:cid/series/:sid", api.UpdateChartSeriesHandler(db))
		protected.DELETE("/chart/:cid/series", api.DeleteAllChartSeriesHandler(db))
		/// Dashboard itself
		protected.GET("/dashboards", api.GetDashboardListHandler(db))
		protected.GET("/dashboard/:did", api.GetDashboardHandler(db))
		protected.POST("/dashboard", api.CreateDashboardPageHandler(db))
		protected.DELETE("/dashboard/:did", api.DeleteDashboardHandler(db))
		protected.GET("/dashboard/:did/charts", api.ListDashboardChartsHandler(db))
		protected.PATCH("/dashboard/:did/chart/:cid", api.AppendChartToDashboardHandler(db))
		protected.DELETE("/dashboard/:did/chart/:cid", api.RemoveChartFromDashboardHandler(db))

		// Query Endpoints
		protected.GET("/queries", api.GetQueryList(db))
		protected.POST("/query/executeLiteral", api.ExecuteQueryLiteralHandler(db)) // Execute a saved query (custom or standard saved queries)
		protected.GET("/query/:qid", api.ReadQueryLiteralHandler(db))               // Return the query SQL string
		protected.GET("/query/:qid/execute", api.ExecuteQueryHandler(db))           // Execute a saved query (custom or standard saved queries)
		protected.POST("/query/:qid", api.SaveQueryHandler(db))
		protected.DELETE("/query/:qid", api.DeleteQueryHandler(db))
		protected.PATCH("/query/:qid", api.UpdateQueryHandler(db)) // Update a custom query in the database

		// Documenation Endpoints
		protected.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

		// Misc Endpoints
		protected.GET("/ping", api.PingHandler)

		admin := protected.Group("/admin")
		{
			admin.Use(api.AuthHandler(activeUsers))

			// Account creation
			admin.POST("/signup", api.SignUpHandler(db, activeUsers))

			// Users
			admin.GET("/users", api.GetUsersHandler(db))
			admin.GET("/user/:uid", api.GetUserHandler(db))
			admin.DELETE("/user/:uname", api.DeleteUserHandler(db, activeUsers))
			admin.PATCH("/user/:uid", api.UpdateUserHandler(db))

			// User roles
			admin.DELETE("/user/:uname/role/:rname", api.DeleteUserRoleHandler(db))
			admin.POST("/user/:uname/role/:rname", api.AddUserRoleHandler(db))

			// Roles
			admin.GET("/roles", api.GetRolesHandler(db))
			admin.POST("/role", api.AddRoleHandler(db))
			admin.PATCH("/role/:rid", api.UpdateRoleHandler(db))
			admin.DELETE("/role/:rid", api.DeleteRoleHandler(db))

			// Permissions
			admin.GET("/permissions", api.GetPermissionsHandler(db))
			admin.GET("/permissions/:scope", api.GetScopedPermissionsHandler(db))
			admin.POST("/dashboard/:did/role/:rid/permission/:pid", api.AddDashboardRolePermissionHandler(db))
			admin.DELETE("/dashboard/:did/role/:rid/permission/:pid", api.DeleteDashboardRolePermissionHandler(db))

			// Settings
			admin.GET("/settings", api.GetAllSettingsHandler(db))
			admin.GET("/settings/:key", api.GetSettingHandler(db))
			admin.PATCH("/settings/:key", api.UpdateSettingHandler(db, &cdb))
		}
	}

	err = server.Run(":8080")
	if err != nil {
		panic(err)
	}
}

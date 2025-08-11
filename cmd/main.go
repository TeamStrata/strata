package main

import (
	"log"

	_ "github.com/TeamStrata/strata/docs"
	"github.com/TeamStrata/strata/pkg/api"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func main() {
	server, err := api.InitServer()
	if err != nil {
		log.Printf("error initializing server: %s", err.Error())
		return
	}

	// Frontend assets
	server.Engine.Static("/assets", "ui/dist/assets")
	server.Engine.StaticFile("/", "ui/dist/index.html")
	server.Engine.StaticFile("/favicon.ico", "ui/dist/favicon.ico")

	// Documenation endpoints
	server.Engine.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// Auth endpoints
	server.Engine.POST("/login", server.LoginHandler)

	protected := server.Engine.Group("/api")
	{
		protected.Use(server.AuthHandler)

		// Auth routes
		protected.POST("/logout", server.LogoutHandler)
		protected.POST("/auth", server.AuthHandler)
		protected.GET("/isadmin/:uname", server.IsAdminHandler)

		// Dashboard Components
		// Charts
		protected.GET("/charts", server.GetChartListHandler)
		protected.GET("/chart/:cid", server.GetChartHandler)
		protected.POST("/chart", server.CreateChartHandler)
		protected.DELETE("/chart/:cid", server.DeleteChartHandler)
		protected.PATCH("/chart/:cid", server.UpdateChartHandler)
		// Chart Series
		protected.GET("/chart/:cid/series", server.GetChartSeriesListHandler)
		protected.GET("/chart/:cid/series/:sid", server.GetChartSingleSeriesHandler)
		protected.POST("/chart/:cid/series", server.AddChartSeriesHandler)
		protected.DELETE("/chart/:cid/series/:sid", server.DeleteChartSingleSeriesHandler)
		protected.PATCH("/chart/:cid/series/:sid", server.UpdateChartSeriesHandler)
		protected.DELETE("/chart/:cid/series", server.DeleteAllChartSeriesHandler)
		// Dashboard itself
		protected.GET("/dashboards", server.GetDashboardListHandler)
		protected.GET("/dashboard/:did", server.GetDashboardHandler)
		protected.POST("/dashboard", server.CreateDashboardPageHandler)
		protected.DELETE("/dashboard/:did", server.DeleteDashboardHandler)
		protected.GET("/dashboard/:did/charts", server.ListDashboardChartsHandler)
		protected.PATCH("/dashboard/:did/chart/:cid", server.AppendChartToDashboardHandler)
		protected.DELETE("/dashboard/:did/chart/:cid", server.RemoveChartFromDashboardHandler)
		protected.GET("/dashboard/:did/full", server.GetFullDashboardData)
		// Dashboard permissions
		protected.GET("/dashboard/:did/permissions", server.GetDashboardRolePermissionsHandler)
		protected.PATCH("/dashboard", server.UpdateDashboardHandler)

		// Query Endpoints
		protected.GET("/queries", server.GetQueryList)
		protected.POST("/query/executeLiteral", server.ExecuteQueryLiteralHandler) // Execute a saved query (custom or standard saved queries)
		protected.GET("/query/:qid", server.ReadQueryLiteralHandler)               // Return the query SQL string
		protected.GET("/query/:qid/execute", server.ExecuteQueryHandler)           // Execute a saved query (custom or standard saved queries)
		protected.POST("/query/:qid", server.SaveQueryHandler)
		protected.DELETE("/query/:qid", server.DeleteQueryHandler)
		protected.PATCH("/query/:qid", server.UpdateQueryHandler) // Update a custom query in the database

		// Misc Endpoints
		protected.GET("/ping", server.PingHandler)

		// AI
		protected.POST("/stratabot", server.StrataBotHandler)
		protected.GET("/cdb-schema", server.SchemaDump)

		admin := protected.Group("/admin")
		{
			admin.Use(server.AuthHandler)

			// Account creation
			admin.POST("/signup", server.SignUpHandler)

			// Users
			admin.GET("/users", server.GetUsersHandler)
			admin.GET("/user/:uid", server.GetUserHandler)
			admin.DELETE("/user/:uname", server.DeleteUserHandler)
			admin.PATCH("/user/:uid", server.UpdateUserHandler)

			// User roles
			admin.DELETE("/user/:uname/role/:rname", server.DeleteUserRoleHandler)
			admin.POST("/user/:uname/role/:rname", server.AddUserRoleHandler)

			// Roles
			admin.GET("/roles", server.GetRolesHandler)
			admin.POST("/role", server.AddRoleHandler)
			admin.PATCH("/role/:rid", server.UpdateRoleHandler)
			admin.DELETE("/role/:rid", server.DeleteRoleHandler)

			// Permissions
			admin.GET("/permissions", server.GetPermissionsHandler)
			admin.GET("/permissions/:scope", server.GetScopedPermissionsHandler)
			admin.POST("/dashboard/:did/role/:rid/permission/:pid", server.AddDashboardRolePermissionHandler)
			admin.DELETE("/dashboard/:did/role/:rid/permission/:pid", server.DeleteDashboardRolePermissionHandler)

			// Settings
			admin.GET("/settings", server.GetAllSettingsHandler)
			admin.GET("/settings/:key", server.GetSettingHandler)
			admin.PATCH("/settings/:key", server.UpdateSettingHandler)
		}
	}

	err = server.Engine.Run(":8080")
	if err != nil {
		panic(err)
	}
}

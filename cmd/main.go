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

	// Initialize map for active users and uuids
	activeUsers := make(map[string]string)

	server.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// Auth endpoints
	server.POST("/login", api.LoginHandler(db, activeUsers))
	server.POST("/signup", api.SignUpHandler(db, activeUsers))
	server.POST("/logout", api.LogoutHandler(activeUsers))
	server.POST("/auth", api.AuthHandler(activeUsers))

	// Frontend
	server.Static("/assets", "ui/dist/assets")
	server.StaticFile("/", "ui/dist/index.html")
	server.StaticFile("/favicon.ico", "ui/dist/favicon.ico")

	// Users
	server.GET("/users", api.GetUsersHandler(db))
	server.GET("/user/:uid", api.GetUserHandler(db))
	server.DELETE("/user/:uname", api.DeleteUserHandler(db, activeUsers))

	// User roles
	server.DELETE("/user/:uname/role/:rname", api.DeleteUserRoleHandler(db))
	server.PATCH("/user/:uname/role/:rname", api.AddUserRoleHandler(db))

	// Roles
	server.GET("/roles", api.GetRolesHandler(db))
	server.POST("/role", api.AddRoleHandler(db))
	server.PATCH("/role/:rid", api.UpdateRoleHandler(db))
	server.DELETE("/role/:rid", api.DeleteRoleHandler(db))

	// Dashboard Components
	/// Charts
	server.GET("/charts", api.GetChartListHandler(db))
	server.GET("/chart/:cid", api.GetChartHandler(db))
	server.POST("/chart", api.CreateChartHandler(db))
	server.DELETE("/chart/:cid", api.DeleteChartHandler(db))
	/// Chart Series
	server.GET("/chart/:cid/series", api.GetChartSeriesListHandler(db))
	server.GET("/chart/:cid/series/:sid", api.GetChartSingleSeriesHandler(db))
	server.POST("/chart/:cid/series", api.AddChartSeriesHandler(db))
	server.DELETE("/chart/:cid/series/:sid", api.DeleteChartSingleSeriesHandler(db))
	server.DELETE("/chart/:cid/series", api.DeleteAllChartSeriesHandler(db))
	/// Dashboard itself
	server.GET("/dashboards", api.GetDashboardListHandler(db))
	server.GET("/dashboard/:did", api.GetDashboardHandler(db))
	server.POST("/dashboard", api.CreateDashboardPageHandler(db))
	server.DELETE("/dashboard/:did", api.DeleteDashboardHandler(db))
	server.GET("/dashboard/:did/charts", api.ListDashboardChartsHandler(db))
	server.PATCH("/dashboard/:did/chart/:cid", api.AppendChartToDashboardHandler(db))
	server.DELETE("/dashboard/:did/chart/:cid", api.RemoveChartFromDashboardHandler(db))

	// Query Endpoints
	server.GET("/queries", api.GetQueryList(db))
	server.GET("/query/:qid", api.ReadQueryLiteralHandler(db))     // Return the query SQL string
	server.GET("/query/:qid/execute", api.ExecuteQueryHandler(db)) // Execute a saved query (custom or standard saved queries)
	server.POST("/query/:qid", api.SaveQueryHandler(db))
	server.DELETE("/query/:qid", api.DeleteQueryHandler(db))

	// Misc Endpoints
	server.GET("/ping", api.PingHandler)

	err = server.Run(":8080")
	if err != nil {
		panic(err)
	}
}

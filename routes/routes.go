package routes

import (
	"gin-M-TIX/config"
	"gin-M-TIX/controllers"
	bookingfacade "gin-M-TIX/patterns/facade"
	"gin-M-TIX/repositories"
	"gin-M-TIX/services"

	"github.com/gin-gonic/gin"
)

func SetupRouter(db *config.Database) *gin.Engine {
	router := gin.Default()

	// Serve static files for frontend
	router.Static("/ui", "./public")

	// Redirect root to frontend
	router.GET("/", func(c *gin.Context) {
		c.Redirect(302, "/ui/")
	})

	movieRepo := repositories.NewMovieRepository(db)
	studioRepo := repositories.NewStudioRepository(db)
	scheduleRepo := repositories.NewScheduleRepository(db)
	bookingRepo := repositories.NewBookingRepository(db)
	userRepo := repositories.NewUserRepository(db)

	pricingService := services.NewPricingService()
	bookingService := services.NewBookingService(bookingRepo, scheduleRepo, pricingService)
	paymentService := services.NewPaymentService(bookingRepo)
	bookingFacade := bookingfacade.NewBookingFacade(bookingService, paymentService)

	movieController := controllers.NewMovieController(movieRepo)
	studioController := controllers.NewStudioController(studioRepo)
	scheduleController := controllers.NewScheduleController(scheduleRepo, pricingService)
	bookingController := controllers.NewBookingController(bookingFacade)
	authController := controllers.NewAuthController(userRepo)

	router.POST("/register", authController.Register)
	router.POST("/login", authController.Login)

	router.GET("/movies", movieController.GetMovies)
	router.GET("/schedules", scheduleController.GetSchedules)
	router.GET("/schedules/:id/seats", scheduleController.GetScheduleSeats)
	router.GET("/studios", studioController.GetStudios)

	auth := router.Group("")
	auth.Use(authController.RequireAuth)
	auth.POST("/logout", authController.Logout)
	auth.GET("/users/me", authController.Me)
	auth.POST("/users/me/student-application", authController.SubmitStudentApplication)
	auth.GET("/users/me/bookings", bookingController.GetUserBookings)
	auth.POST("/bookings", authController.RequireNonAdmin, bookingController.CreateBooking)
	auth.GET("/bookings/:id", bookingController.GetBooking)
	auth.DELETE("/bookings/:id", bookingController.CancelBooking)
	auth.POST("/payments", authController.RequireNonAdmin, bookingController.Pay)

	admin := router.Group("")
	admin.Use(authController.RequireAdmin)
	admin.POST("/movies", movieController.CreateMovie)
	admin.PUT("/movies/:id", movieController.UpdateMovie)
	admin.DELETE("/movies/:id", movieController.DeleteMovie)
	admin.POST("/schedules", scheduleController.CreateSchedule)
	admin.PUT("/schedules/:id", scheduleController.UpdateSchedule)
	admin.DELETE("/schedules/:id", scheduleController.DeleteSchedule)
	admin.POST("/studios", studioController.CreateStudio)
	admin.PUT("/studios/:id", studioController.UpdateStudio)
	admin.DELETE("/studios/:id", studioController.DeleteStudio)
	admin.GET("/admin/student-applications", authController.ListStudentApplications)
	admin.GET("/admin/student-applications/:id/evidence", authController.StudentEvidence)
	admin.POST("/admin/student-applications/:id/resolve", authController.ResolveStudentApplication)

	return router
}

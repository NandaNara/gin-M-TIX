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

	movieRepo := repositories.NewMovieRepository(db)
	scheduleRepo := repositories.NewScheduleRepository(db)
	bookingRepo := repositories.NewBookingRepository(db)

	pricingService := services.NewPricingService()
	bookingService := services.NewBookingService(bookingRepo, scheduleRepo, pricingService)
	paymentService := services.NewPaymentService(bookingRepo)
	bookingFacade := bookingfacade.NewBookingFacade(bookingService, paymentService)

	movieController := controllers.NewMovieController(movieRepo)
	scheduleController := controllers.NewScheduleController(scheduleRepo, pricingService)
	bookingController := controllers.NewBookingController(bookingFacade)

	router.POST("/login", controllers.Login)

	router.GET("/movies", movieController.GetMovies)
	router.POST("/movies", movieController.CreateMovie)
	router.PUT("/movies/:id", movieController.UpdateMovie)
	router.DELETE("/movies/:id", movieController.DeleteMovie)

	router.GET("/schedules", scheduleController.GetSchedules)
	router.POST("/schedules", scheduleController.CreateSchedule)
	router.GET("/schedules/:id/seats", scheduleController.GetScheduleSeats)

	router.POST("/bookings", bookingController.CreateBooking)
	router.GET("/bookings/:id", bookingController.GetBooking)
	router.GET("/users/:id/bookings", bookingController.GetUserBookings)

	router.POST("/payments", bookingController.Pay)

	return router
}

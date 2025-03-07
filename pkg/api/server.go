package http

import (
	"log"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	handler "ahava/pkg/api/handler"
	"ahava/pkg/api/middleware"
	"ahava/pkg/routes"
)

type ServerHTTP struct {
	engine *gin.Engine
}

func NewServerHTTP(
	userHandler handler.UserHandler,
	adminHandler handler.AdminHandler,
	productHandler handler.ProductHandler,
	// otpHandler handler.OtpHandler,
	orderHandler handler.OrderHandler,
	cartHandler handler.CartHandler,
	// couponHandler handler.CouponHandler,
	paymentHandler handler.PaymentHandler,
	// offerhandler handler.OfferHandler,
	wishlistHandler handler.WishlistHandler,
	newsHandler handler.NewsHandler,
	uploadHandler handler.UploadHandler,
	db *gorm.DB,
) *ServerHTTP {

	gin.SetMode(gin.ReleaseMode)
	engine := gin.New()
	engine.Use(middleware.CORSMiddleware())
	engine.Use(middleware.DefaultStructuredLogger())
	// engine.Use(gin.Logger())
	go middleware.SaveRequestTransaction(db)

	engine.GET("/validate-token", adminHandler.ValidateRefreshTokenAndCreateNewAccess)

	routes.UserRoutes(engine.Group("/api"),
		userHandler,
		// otpHandler,
		productHandler,
		orderHandler,
		cartHandler,
		paymentHandler,
		wishlistHandler,
		newsHandler,
		// couponHandler,
	)
	routes.AdminRoutes(engine.Group("/admin"),
		adminHandler,
		productHandler,
		userHandler,
		uploadHandler,
		orderHandler,
		newsHandler,
		// couponHandler,
		// offerhandler,
	)

	return &ServerHTTP{engine: engine}
}

func (sh *ServerHTTP) Start() {
	err := sh.engine.Run(":8088")
	if err != nil {
		log.Fatal("gin engine couldn't start")
	}
}

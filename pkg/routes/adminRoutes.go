package routes

import (
	"ahava/pkg/api/handler"
	"ahava/pkg/api/middleware"

	"github.com/gin-gonic/gin"
)

func AdminRoutes(
	engine *gin.RouterGroup,
	adminHandler handler.AdminHandler,
	productHandler handler.ProductHandler,
	userHandler handler.UserHandler,
	// orderHandler handler.OrderHandler,
	// couponHandler handler.CouponHandler,
	// offerHandler handler.OfferHandler,
) {

	engine.POST("/login", adminHandler.Login)
	engine.Use(middleware.AdminAuthMiddleware)
	{
		usermanagement := engine.Group("/users")
		{
			usermanagement.GET("", adminHandler.GetUsers)
			usermanagement.PUT("/block/:id", adminHandler.BlockUser)
			usermanagement.PUT("/unblock/:id", adminHandler.UnBlockUser)
		}

		productmanagement := engine.Group("/product")
		{
			productmanagement.GET("", productHandler.ListAllProducts)
			productmanagement.GET("/detail", productHandler.GetProductDetails)
			productmanagement.POST("", productHandler.AddProduct)
			productmanagement.POST("/:product_id", productHandler.AddProductImages)
			productmanagement.DELETE("/:product_id", productHandler.DeleteProduct)
			productmanagement.PUT("/:product_id", productHandler.UpdateProduct)
			productmanagement.PUT("/:product_id/image", productHandler.UpdateProductImage)
		}

		// payment := engine.Group("/payment-method")
		// {
		// 	payment.POST("", adminHandler.NewPaymentMethod)
		// 	payment.GET("", adminHandler.ListPaymentMethods)
		// 	payment.DELETE("", adminHandler.DeletePaymentMethod)
		// }

		// orders := engine.Group("/orders")
		// {
		// 	orders.PUT("/status", orderHandler.EditOrderStatus)
		// 	orders.PUT("/payment-status", orderHandler.MakePaymentStatusAsPaid)
		// 	orders.GET("", orderHandler.AdminOrders)
		// 	orders.GET("/:id", orderHandler.GetIndividualOrderDetails)
		// }

		// coupons := engine.Group("/coupons")
		// {
		// 	coupons.GET("", couponHandler.GetAllCoupons)
		// 	coupons.POST("", couponHandler.CreateNewCoupon)
		// 	coupons.DELETE("", couponHandler.MakeCouponInvalid)
		// 	//reactivation of coupons
		// 	coupons.PUT("", couponHandler.ReActivateCoupon)
		// }

		// offers := engine.Group("/offer")
		// {
		// 	offers.GET("", offerHandler.GetOffers)
		// 	offers.POST("", offerHandler.AddNewOffer)
		// 	offers.DELETE("/:product_id", offerHandler.ExpireOffer)
		// 	offers.PUT("/:product_id", offerHandler.UpdateOffer)
		// }
	}

}

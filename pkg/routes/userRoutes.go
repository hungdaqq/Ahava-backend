package routes

import (
	"ahava/pkg/api/handler"
	"ahava/pkg/api/middleware"

	"github.com/gin-gonic/gin"
)

func UserRoutes(
	engine *gin.RouterGroup,
	userHandler *handler.UserHandler,
	// otpHandler *handler.OtpHandler,
	productHandler *handler.ProductHandler,
	orderHandler *handler.OrderHandler,
	cartHandler *handler.CartHandler,
	// paymentHandler *handler.PaymentHandler,
	wishlisthandler *handler.WishlistHandler,
	// categoryHandler *handler.CategoryHandler,
	// couponHandler *handler.CouponHandler
) {

	engine.POST("/signup", userHandler.UserSignUp)
	engine.POST("/login", userHandler.LoginHandler)
	// engine.GET("/forgot-password", userHandler.ForgotPasswordSend)
	// engine.POST("/forgot-password", userHandler.ForgotPasswordVerifyAndChange)

	// engine.POST("/otplogin", otpHandler.SendOTP)
	// engine.POST("/verifyotp", otpHandler.VerifyOTP)

	// payment := engine.Group("/payment")
	// {
	// 	payment.GET("/razorpay", paymentHandler.MakePaymentRazorPay)
	// 	payment.GET("/update_status", paymentHandler.VerifyPayment)
	// }

	engine.Use(middleware.UserAuthMiddleware)
	{

		// engine.GET("/banners", categoryHandler.GetBannersForUsers)

		home := engine.Group("/home")
		{
			home.GET("/products", productHandler.ListProductsForUser)
			home.GET("/products/details", productHandler.ShowProductDetails)
			// home.POST("/add-to-cart", cartHandler.AddToCart)
			// home.POST("/wishlist/add", wishlisthandler.AddToWishlist)
			home.POST("/search", productHandler.SearchProducts)
			home.GET("/search", productHandler.GetSearchHistory)

		}

		// categorymanagement := engine.Group("/category")
		// {
		// 	categorymanagement.GET("", categoryHandler.GetCategory)
		// 	categorymanagement.GET("/products", categoryHandler.GetProductDetailsInACategory)

		// }

		profile := engine.Group("/profile")
		{
			profile.GET("/detail", userHandler.GetUserDetails)
			profile.GET("/address", userHandler.GetAddresses)
			profile.POST("/address", userHandler.AddAddress)
			profile.PUT("/address", userHandler.UpdateAddress)
			profile.DELETE("/address", userHandler.DeleteAddress)

			// profile.GET("/reference-link", userHandler.GetMyReferenceLink)

			// orders := profile.Group("/orders")
			// {
			// 	orders.GET("", orderHandler.GetOrders)
			// 	orders.GET("/:id", orderHandler.GetIndividualOrderDetails)
			// 	orders.DELETE("", orderHandler.CancelOrder)
			// 	orders.PUT("/return", orderHandler.ReturnOrder)
			// }

			edit := profile.Group("/edit")
			{
				edit.PUT("", userHandler.EditProfile)
				edit.PUT("/password", userHandler.ChangePassword)

			}
		}

		cart := engine.Group("/cart")
		{
			cart.GET("", cartHandler.GetCart)
			cart.POST("", cartHandler.AddToCart)
			cart.DELETE("/:cart_id", cartHandler.RemoveFromCart)
			cart.PUT("/:cart_id/plus", cartHandler.UpdateQuantityAdd)
			cart.PUT("/:cart_id/minus", cartHandler.UpdateQuantityLess)
			cart.PUT("/:cart_id", cartHandler.UpdateQuantityLess)

			cart.POST("/check-out", cartHandler.CheckOut)
		}

		wishlist := engine.Group("/wishlist")
		{
			wishlist.POST("/:product_id", wishlisthandler.AddToWishlist)
			wishlist.GET("", wishlisthandler.GetWishList)
			wishlist.DELETE(":product_id", wishlisthandler.RemoveFromWishlist)
		}

		checkout := engine.Group("/order")
		{
			checkout.POST("", orderHandler.PlaceOrder)
		}

		payment := engine.Group("/payment")
		{
			payment.POST("/sepay", paymentHandler.CreateQR)
		}

		// engine.GET("/coupon", couponHandler.GetAllCoupons)

	}

}

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
	// orderHandler *handler.OrderHandler,
	cartHandler *handler.CartHandler,
	// paymentHandler *handler.PaymentHandler,
	// wishlisthandler *handler.WishlistHandler,
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

		}

		// categorymanagement := engine.Group("/category")
		// {
		// 	categorymanagement.GET("", categoryHandler.GetCategory)
		// 	categorymanagement.GET("/products", categoryHandler.GetProductDetailsInACategory)

		// }

		profile := engine.Group("/profile")
		{
			profile.GET("/detail", userHandler.GetUserDetails)
			// profile.GET("/address", userHandler.GetAddresses)
			// profile.POST("/address", userHandler.AddAddress)
			profile.GET("/reference-link", userHandler.GetMyReferenceLink)

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
			// cart.DELETE("", cartHandler.RemoveFromCart)
			cart.PUT("/quantity/plus", cartHandler.UpdateQuantityAdd)
			// cart.PUT("/updateQuantity/minus", cartHandler.UpdateQuantityLess)
			// hello
		}

		// wishlist := engine.Group("/wishlist")
		// {
		// 	wishlist.GET("/", wishlisthandler.GetWishList)
		// 	wishlist.DELETE("/remove", wishlisthandler.RemoveFromWishlist)
		// }

		// checkout := engine.Group("/check-out")
		// {
		// 	checkout.GET("", cartHandler.CheckOut)
		// 	checkout.POST("/order", orderHandler.OrderItemsFromCart)
		// }

		// engine.GET("/coupon", couponHandler.GetAllCoupons)

	}

}

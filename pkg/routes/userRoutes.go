package routes

import (
	"ahava/pkg/api/handler"
	"ahava/pkg/api/middleware"

	"github.com/gin-gonic/gin"
)

func UserRoutes(
	engine *gin.RouterGroup,
	userHandler handler.UserHandler,
	// otpHandler handler.OtpHandler,
	productHandler handler.ProductHandler,
	orderHandler handler.OrderHandler,
	cartHandler handler.CartHandler,
	paymentHandler handler.PaymentHandler,
	wishlisthandler handler.WishlistHandler,
	// couponHandler handler.CouponHandler
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

	payment := engine.Group("/payment")
	{
		payment.POST("/webhook", paymentHandler.Webhook)
	}

	home := engine.Group("/home")
	{
		home.POST("/search", productHandler.SearchProducts)
		// home.GET("/search", productHandler.GetSearchHistory)
	}

	product := engine.Group("/product")
	{
		product.GET("/detail", productHandler.GetProductDetails)
		product.GET("", productHandler.ListCategoryProducts)
		product.GET("/featured", productHandler.ListFeaturedProducts)
	}

	engine.Use(middleware.UserAuthMiddleware)
	{
		profile := engine.Group("/profile")
		{
			profile.GET("/detail", userHandler.GetUserDetails)
			address := profile.Group("/address")
			{
				address.GET("", userHandler.GetAddresses)
				address.POST("", userHandler.AddAddress)
				address.PUT("/:address_id", userHandler.UpdateAddress)
				address.DELETE("/:address_id", userHandler.DeleteAddress)
			}
			// profile.GET("/reference-link", userHandler.GetMyReferenceLink)
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
			cart.PUT("/:cart_id", cartHandler.UpdateQuantity)
			cart.POST("/check-out", cartHandler.CheckOut)
		}

		wishlist := engine.Group("/wishlist")
		{
			wishlist.POST("", wishlisthandler.AddToWishlist)
			wishlist.GET("", wishlisthandler.GetWishList)
			wishlist.DELETE("/:wishlist_id", wishlisthandler.RemoveFromWishlist)
		}

		order := engine.Group("/order")
		{
			order.GET("/detail", orderHandler.GetOrderDetails)
			order.POST("", orderHandler.PlaceOrder)
		}

		payment := engine.Group("/payment")
		{
			payment.POST("/qr", paymentHandler.CreateQR)
		}
		// engine.GET("/coupon", couponHandler.GetAllCoupons)
	}
}

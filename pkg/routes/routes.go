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
	// couponHandler handler.CouponHandler,
	// offerHandler handler.OfferHandler,
	uploadHandler handler.UploadHandler,
	orderHandler handler.OrderHandler,
	newsHandler handler.NewsHandler,
) {
	engine.POST("/login", adminHandler.Login)
	engine.Use(middleware.AdminAuthMiddleware)
	{
		filemanagement := engine.Group("/file")
		{
			filemanagement.POST("/upload", uploadHandler.FileUpload)
		}
		usermanagement := engine.Group("/user")
		{
			usermanagement.GET("", adminHandler.ListAllUsers)
			usermanagement.PUT("/block/:id", adminHandler.BlockUser)
			usermanagement.PUT("/unblock/:id", adminHandler.UnBlockUser)
		}

		productmanagement := engine.Group("/product")
		{
			productmanagement.GET("", productHandler.ListAllProducts)
			productmanagement.GET("/detail", productHandler.GetProductDetails)
			productmanagement.POST("", productHandler.AddProduct)
			productmanagement.DELETE("/:product_id", productHandler.DeleteProduct)
			productmanagement.PUT("/:product_id", productHandler.UpdateProduct)
		}
		ordermanagement := engine.Group("/order")
		{
			ordermanagement.GET("", orderHandler.ListAllOrders)
		}
		newsmanagement := engine.Group("/news")
		{
			newsmanagement.GET("", newsHandler.ListAllNews)
			newsmanagement.GET("/:news_id", newsHandler.GetNewsByID)
			newsmanagement.POST("", newsHandler.AddNews)
			newsmanagement.PUT("/:news_id", newsHandler.UpdateNews)
			newsmanagement.DELETE("/:news_id", newsHandler.DeleteNews)
		}
		// payment := engine.Group("/payment-method")
		// {
		// 	payment.POST("", adminHandler.NewPaymentMethod)
		// 	payment.GET("", adminHandler.ListPaymentMethods)
		// 	payment.DELETE("", adminHandler.DeletePaymentMethod)
		// }

		// coupons := engine.Group("/coupons")
		// {
		// 	coupons.GET("", couponHandler.GetAllCoupons)
		// 	coupons.POST("", couponHandler.CreateNewCoupon)
		// 	coupons.DELETE("", couponHandler.MakeCouponInvalid)
		// 	//reactivation of coupons
		// 	coupons.PUT("", couponHandler.ReActivateCoupon)
		// }
	}
}

func UserRoutes(
	engine *gin.RouterGroup,
	userHandler handler.UserHandler,
	// otpHandler handler.OtpHandler,
	productHandler handler.ProductHandler,
	orderHandler handler.OrderHandler,
	cartHandler handler.CartHandler,
	paymentHandler handler.PaymentHandler,
	wishlisthandler handler.WishlistHandler,
	newsHandler handler.NewsHandler,
	// couponHandler handler.CouponHandler
) {

	engine.POST("/signup", userHandler.Register)
	engine.POST("/login", userHandler.Login)
	// engine.GET("/forgot-password", userHandler.ForgotPasswordSend)
	// engine.POST("/forgot-password", userHandler.ForgotPasswordVerifyAndChange)

	// engine.POST("/otplogin", otpHandler.SendOTP)
	// engine.POST("/verifyotp", otpHandler.VerifyOTP)

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
	news := engine.Group("/news")
	{
		news.GET("", newsHandler.GetFeaturedNews)
		news.GET("/:news_id", newsHandler.GetNewsByID)
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
			cart.PUT("/:cart_id", cartHandler.UpdateQuantity)
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

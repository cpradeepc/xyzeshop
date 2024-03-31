package router

import (
	"net/http"
	"xyzeshop/constval"
	"xyzeshop/controller"
)

var healthCheckRoutes = Routes{
	Route{
		Name:        "Health check",
		Method:      http.MethodGet,
		Pattern:     constval.HealthCheckRoutes,
		HandlerFunc: controller.HealthCheck,
	},
}

var userRoutes = Routes{
	Route{
		Name:        "VerifyEmail",
		Method:      http.MethodPost,
		Pattern:     constval.VerifyEmailRoute,
		HandlerFunc: controller.VerifyEmail,
	},
	Route{
		Name:        "VerifyOtp",
		Method:      http.MethodPost,
		Pattern:     constval.VerifyOtpRoute,
		HandlerFunc: controller.VerifyOtp,
	},
	Route{
		Name:        "ResendEmail",
		Method:      http.MethodPost,
		Pattern:     constval.ResendEmailRoute,
		HandlerFunc: controller.VerifyEmail,
	},
	// register user
	Route{
		Name:        "Register User",
		Method:      http.MethodPost,
		Pattern:     constval.UserRegdRoute,
		HandlerFunc: controller.RegisterUser,
	},
	Route{
		Name:        "Signin User",
		Method:      http.MethodPost,
		Pattern:     constval.UserSigninRoute,
		HandlerFunc: controller.UserSignin,
	},
}

var productGlobalRoutes = Routes{
	Route{
		Name:        "List Product",
		Method:      http.MethodGet,
		Pattern:     constval.ListProductRoute,
		HandlerFunc: controller.ListProducts,
	},
	Route{
		Name:        "Search Product",
		Method:      http.MethodGet,
		Pattern:     constval.SearchProductRoute,
		HandlerFunc: controller.SearchProduct,
	},
}

var productRoutes = Routes{
	Route{
		Name:        "Register Product",
		Method:      http.MethodPost,
		Pattern:     constval.RegdProductRoute,
		HandlerFunc: controller.RegProduct,
	},
	Route{
		Name:        "Update Product",
		Method:      http.MethodPut,
		Pattern:     constval.UpdateProductRoute,
		HandlerFunc: controller.UpdateProduct,
	},
	Route{
		Name:        "Delete Product",
		Method:      http.MethodDelete,
		Pattern:     constval.DeleteProductRoute,
		HandlerFunc: controller.DelProduct,
	},
}

var userAuthRoutes = Routes{
	Route{
		Name:        "Add to cart",
		Method:      http.MethodPost,
		Pattern:     constval.AddToCartRoute,
		HandlerFunc: controller.AddCart,
	},
	Route{
		Name:        "Add Address",
		Method:      http.MethodPost,
		Pattern:     constval.AddAddrRoute,
		HandlerFunc: controller.AddAddressUser,
	},
	Route{
		Name:        "Get single user",
		Method:      http.MethodGet,
		Pattern:     constval.GetSingleUserRoute,
		HandlerFunc: controller.GetUser,
	},
	Route{
		Name:        "Update User",
		Method:      http.MethodPut,
		Pattern:     constval.UpdateUserRoute,
		HandlerFunc: controller.UpdateUser,
	},
	Route{
		Name:        "Checkout Order",
		Method:      http.MethodPut,
		Pattern:     constval.CheckoutRoute,
		HandlerFunc: controller.CheckoutProductById,
	},
}

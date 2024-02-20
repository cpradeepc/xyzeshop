package constval

const (
	APIversion string = "v1"
	BadReqMsg  string = "request not fullfilled"

	//schedular constants value
	HealthCheckRoutes string = "/health"
	MDBuri            string = "localhost:27017"
	DbName            string = "ecomshop"
	Sender            string = ""

	//email verification routes
	VerifyEmailRoute string = "/verify_email"
	VerifyOtpRoute   string = "/verify_otp"
	ResendEmailRoute string = "/resend_eamil"

	//users routes
	UserRegdRoute   string = "/user-regd"
	UserSigninRoute string = "/signin"

	//product routes
	RegdProductRoute   string = "/product_regd"
	ListProductRoute   string = "/list_products"
	SearchProductRoute string = "/src_products"
	UpdateProductRoute string = "/update_products"
	DeleteProductRoute string = "/del_products"
	AddToCartRoute     string = "/cart"
	AddAddrRoute       string = "/addr"
	GetSingleUserRoute string = "/user/:id"
	UpdateUserRoute    string = "/update_user"
	CheckoutRoute      string = "user/:id"
)

const (
	GuestUser string = "guestuser"
	AdminUser string = "admin"
)

const (
	//time slot for otp validation
	OtpValidation int = 60
)

// collection
const (
	VerificationsCollection string = "verification"
	UserCollection          string = "user"
	ProductCollection       string = "products"
	AddressCollection       string = "user_addresses"
	CartCollection          string = "user_cart"
)

// messages
const (
	AlreadyRegisterWithThisEmail string = "already register with this email"
	EmailIsNotVerified           string = "your email is not verified please verify your email"
	EmailValidationError         string = "wrong email passed"
	OtpValidationError           string = "wrong otp passed"
	OtpExpiredValidationError    string = "otp expired"
	AlreadyVerifiedError         string = "already verified"
	OptAlreadySentError          string = "otp already sent to email"
	NotRegisteredUser            string = "you are not register user"
	PasswordNotMatchedError      string = "password doesn't match"
	NotAuthorizedUserError       string = "you are not authorized to do this"
	NoProductAvaliable           string = "no product avaliable"
	UserDoesNotExists            string = "user not exists"
	AddressNotExists             string = "address not exists. please add one address"
)

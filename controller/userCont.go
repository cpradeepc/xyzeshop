package controller

import (
	"log"
	"net/http"
	"os"
	"time"
	"xyzeshop/auth"
	"xyzeshop/constval"
	"xyzeshop/dbs"
	"xyzeshop/helper"
	"xyzeshop/payloads"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
)

// email verify process
func VerifyEmail(c *gin.Context) {
	var req payloads.Verification
	err := c.BindJSON(&req)
	if err != nil {
		log.Println("error invocked during bind json :", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": true, "msg": err.Error()})
		return
	}

	if req.Email == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": true, "msg": constval.EmailValidationError})
		return
	}

	resp, _ := dbs.Mgr.GetSingleRecordByEmail(req.Email, constval.VerificationsCollection)
	if resp.OTP != 0 {
		sec := resp.CreatedAt + constval.OtpValidation

		//check otp is expired
		if sec < time.Now().Unix() {
			reqMail, errMail := helper.SendMailOtp(req)
			if errMail != nil {
				log.Println("error invoked during send mail :", errMail)
				c.JSON(http.StatusBadRequest, gin.H{"error": true, "msg": constval.EmailValidationError})
				return
			}
			//update email sending time
			reqMail.CreatedAt = time.Now().Unix()
			//update operation
			dbs.Mgr.UpdateVerification(reqMail, constval.VerificationsCollection)
			c.JSON(http.StatusOK, gin.H{"error": false, "msg": "success"})
			return
		} else {
			c.JSON(http.StatusBadRequest, gin.H{"error": true, "msg": constval.OptAlreadySentError})
			return
		}
	}
	reqMail, errMail := helper.SendMailOtp(req)
	if errMail != nil {
		log.Println("error invoked during send mail :", errMail)
		c.JSON(http.StatusBadRequest, gin.H{"error": true, "msg": constval.EmailValidationError})
		return
	}
	reqMail.CreatedAt = time.Now().Unix()
	//insertion of record
	dbs.Mgr.InsertOne(reqMail, constval.VerificationsCollection)
	c.JSON(http.StatusOK, gin.H{"error": false, "msg": "success"})
}

// verify otp end to end
func VerifyOtp(c *gin.Context) {
	var req payloads.Verification
	err := c.BindJSON(&req)
	if err != nil {
		log.Println("error was invoked during bind with reqest model :", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": true, "msg": err.Error()})
		return
	}

	//checking email
	if req.Email == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": true, "msg": constval.EmailValidationError})
		return
	}

	log.Println("mail :", req.Email, "otp :", req.OTP)
	//checking otp field not empty
	if req.OTP <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": true, "msg": constval.OtpValidationError})
		return
	}

	//checking record in verification collection
	resp, _ := dbs.Mgr.GetSingleRecordByEmail(req.Email, constval.VerificationsCollection)

	//if status or email is already verified
	if resp.Status {
		c.JSON(http.StatusBadRequest, gin.H{"error": true, "msg": constval.AlreadyVerifiedError})
		return
	}

	sec := resp.CreatedAt + constval.OtpValidation
	//otp got in request and from dbs can not matched
	if resp.OTP != req.OTP {
		c.JSON(http.StatusBadRequest, gin.H{"error": true, "msg": constval.OtpValidationError})
		return
	}

	//if otp expired
	if sec < time.Now().Unix() {
		c.JSON(http.StatusBadRequest, gin.H{"error": true, "msg": constval.OtpExpiredValidationError})
		return
	}

	req.Status = true
	req.CreatedAt = time.Now().Unix()
	err = dbs.Mgr.UpdateEmailVerifiedStatus(req, constval.VerificationsCollection)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": true, "msg": constval.OtpValidationError})
		return
	}
	c.JSON(http.StatusOK, gin.H{"error": false, "msg": "success"})
}

// user registration
func RegisterUser(c *gin.Context) {
	var guest payloads.GuestUser //default nil value
	var user payloads.RegdUser

	//bind data with models
	err := c.BindJSON(&guest)
	if err != nil {
		log.Println("error was invoked during bind with user json :", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": true, "msg": err.Error()})
		return
	}

	//check user  validation
	err = helper.CheckUserValidation(guest)
	if err != nil {
		log.Println("error was invoked during user validation :", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": true, "msg": err.Error()})
		return
	}

	//verify email
	resp, _ := dbs.Mgr.GetSingleRecordByEmail(guest.Email, constval.VerificationsCollection)
	if !resp.Status {
		c.JSON(http.StatusBadRequest, gin.H{"error": true, "msg": constval.EmailIsNotVerified})
		return
	}

	//check for duplicate user
	respUser := dbs.Mgr.GetSingleRecordByEmailForUser(guest.Email, constval.VerificationsCollection)
	if respUser.Email != "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": true, "msg": constval.AlreadyRegisterWithThisEmail})
		return
	}

	user.Email = guest.Email
	user.Name = guest.Name
	user.Phone = guest.Phone
	user.UserType = constval.GuestUser

	//encrypted hash psw generate
	hashPsw := helper.GenPassHash(guest.Password)
	user.Password = hashPsw
	user.CreatedAt = time.Now().Unix()
	user.UpdatedAt = time.Now().Unix()

	id, err := dbs.Mgr.InsertOne(user, constval.UserCollection)
	if err != nil {
		log.Println("error was invoked during insert user in db :", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": true, "msg": err.Error()})
		return
	}
	//jwt prepare
	jwtWrap := auth.JwtWrapper{
		SecretKey:      os.Getenv("JwtSecrets"),
		Issuer:         os.Getenv("JwtIssuer"),
		ExpirationTime: 48,
	}
	userID := id.(primitive.ObjectID)
	//generate token
	token, err := jwtWrap.GenerateToken(userID, user.Email, user.UserType)
	if err != nil {
		log.Println("error was invoked during generate token of user :", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": true, "msg": err.Error()})
		return
	}
	//user.Password = ""
	c.JSON(http.StatusOK, gin.H{
		"error": false,
		"msg":   "success",
		"data":  user,
		"token": token,
	})
}

// user login
func UserSignin(c *gin.Context) {
	var loginReq payloads.Login

	err := c.BindJSON(&loginReq)
	if err != nil {
		log.Println("error was invoked during bind json with model :", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": true, "msg": err.Error()})
		return
	}
	user := dbs.Mgr.GetSingleRecordByEmailForUser(loginReq.Email, constval.UserCollection)
	if user.Email == "" {
		log.Println("error was invoked during read mail :", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": true, "msg": constval.NotRegisteredUser})
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(loginReq.Password))
	if err != nil {
		log.Println("error was invoked during match password :", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": true, "msg": constval.PasswordNotMatchedError})
		return
	}

	jwtWrap := auth.JwtWrapper{
		SecretKey:      os.Getenv("JwtSecrets"),
		Issuer:         os.Getenv("JwtIssuer"),
		ExpirationTime: 48,
	}
	token, err := jwtWrap.GenerateToken(user.Id, user.Email, user.UserType)
	if err != nil {
		log.Println("error was invoked during generate token of user :", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": true, "msg": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"error": false, "token": token})

}

// add cart
func AddCart(c *gin.Context) {
	email, ok := c.Get("email")
	if !ok {
		log.Println("error was invoked during get email :", email)
		c.JSON(http.StatusBadRequest, gin.H{"error": true, "msg": constval.EmailIsNotVerified})
		return
	}

	userDbResp, _ := dbs.Mgr.GetSingleRecordByEmail(email.(string), constval.UserCollection)
	if userDbResp.Email == "" {
		log.Println("error was invoked due to empty mail :", userDbResp.Email)
		c.JSON(http.StatusBadRequest, gin.H{"error": true, "msg": constval.EmailIsNotVerified})
		return
	}

	userAdr, err := dbs.Mgr.GetSingleAddress(userDbResp.ID, constval.AddressCollection)
	if err != nil {
		log.Println("error was invoked address not parse :", userAdr.Address1)
		c.JSON(http.StatusBadRequest, gin.H{"error": true, "msg": err.Error()})
		return
	}
	if userAdr.Address1 == "" {
		log.Println("error was invoked address not exist :", userAdr.Address1)
		c.JSON(http.StatusBadRequest, gin.H{"error": true, "msg": constval.AddressNotExists})
		return
	}

	var cartdb payloads.Cart
	var cart payloads.CartUser
	err = c.BindJSON(&cart)
	if err != nil {
		log.Println("error was invoked  during bind json with cart :", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": true, "msg": err.Error()})
		return
	}
	projectId, _ := primitive.ObjectIDFromHex(cart.ProductID)
	userId, _ := primitive.ObjectIDFromHex(cart.UserID)
	cartdb.UserID = userId
	cartdb.ProductID = projectId
	cartdb.Checkout = false
	_, err = dbs.Mgr.InsertOne(cartdb, constval.CartCollection)
	if err != nil {
		log.Println("error was invoked  during insert data of cart :", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": true, "msg": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"error": false, "msg": "success"})
}

// add address of user
func AddAddressUser(c *gin.Context) {
	var adrReq payloads.RegdUserPostalAddress
	err := c.BindJSON(&adrReq)
	if err != nil {
		log.Println("error was invoked  during bind data with user address :", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": true, "msg": err.Error()})
		return
	}

	userId, err := primitive.ObjectIDFromHex(adrReq.UserId)
	if err != nil {
		log.Println("error was invoked  during parse user id :", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": true, "msg": err.Error()})
		return
	}

	var adrDb payloads.RegdUserAddress
	adrDb.Address1 = adrReq.Address1
	adrDb.UserId = userId
	adrDb.City = adrReq.City
	adrDb.Country = adrReq.Country

	_, err = dbs.Mgr.InsertOne(adrDb, constval.AddressCollection)
	if err != nil {
		log.Println("error was invoked  during insert address data :", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": true, "msg": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"error": false, "msg": "success"})

}

// get user only
func GetUser(c *gin.Context) {
	userIdStr := c.Param("id")
	userId, err := primitive.ObjectIDFromHex(userIdStr)
	if err != nil {
		log.Println("error was invoked  during read user id :", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": true, "msg": err.Error()})
		return
	}
	user := dbs.Mgr.GetSingleUserByUserId(userId, constval.UserCollection)
	if user.Email == "" {
		log.Println("error was invoked  empty email :", user.Email)
		c.JSON(http.StatusBadRequest, gin.H{"error": true, "msg": constval.UserDoesNotExists})
		return
	}

	//user.Password == ""
	c.JSON(http.StatusOK, gin.H{"error": false, "msg": "success", "data": user})
}

// update user
func UpdateUser(c *gin.Context) {
	var userUpdate payloads.RegdUserUpdate
	err := c.BindJSON(&userUpdate)
	if err != nil {
		log.Println("error was invoked  during bind with user data :", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": true, "msg": err.Error()})
		return
	}
	userId, err := primitive.ObjectIDFromHex(userUpdate.Id)
	if err != nil {
		log.Println("error was invoked  during parse user object id :", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": true, "msg": err.Error()})
		return
	}
	userResp := dbs.Mgr.GetSingleUserByUserId(userId, constval.UserCollection)
	if userResp.Email == "" {
		log.Println("error was invoked  due to empty user email id :", userResp.Email)
		c.JSON(http.StatusBadRequest, gin.H{"error": true, "msg": constval.UserDoesNotExists})
		return
	}

	var user payloads.RegdUser
	user.Id = userId
	user.Email = userResp.Email
	user.Password = userResp.Password
	user.Phone = userResp.Phone
	user.UserType = userResp.UserType
	user.UpdatedAt = time.Now().Unix()
	//user.CreatedAt = userResp.CreatedAt  //not in all cases except creation time

	if userUpdate.Email != "" {
		user.Email = userUpdate.Email
	}
	if userUpdate.Password != "" {
		user.Password = helper.GenPassHash(userUpdate.Password)
	}
	if userUpdate.Phone != "" {
		user.Phone = userUpdate.Phone
	}
	if userUpdate.Name != "" {
		user.Name = userUpdate.Name
	}

	err = dbs.Mgr.UpdateUser(user, constval.UserCollection)
	if err != nil {
		log.Println("error was invoked  during update user :", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": true, "msg": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"error": false, "msg": "success"})

}

package controller

import (
	"log"
	"net/http"
	"time"
	"xyzeshop/constval"
	"xyzeshop/dbs"
	"xyzeshop/helper"
	"xyzeshop/payloads"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// register product
func RegProduct(c *gin.Context) {
	userMail, ok := c.Get("email")
	if !ok {
		log.Println("error invocked during get email :", userMail)
		c.JSON(http.StatusInternalServerError, gin.H{"error": true, "msg": constval.NotRegisteredUser})
		return
	}

	user := dbs.Mgr.GetSingleRecordByEmailForUser(userMail.(string), constval.UserCollection)
	if user.UserType != constval.AdminUser {
		log.Println("error invocked due to user is not admin :", userMail)
		c.JSON(http.StatusInternalServerError, gin.H{"error": true, "msg": constval.NotAuthorizedUserError})
		return
	}

	var productReq payloads.ProductUser
	var product payloads.Product
	err := c.BindJSON(&productReq)
	if err != nil {
		log.Println("error invocked during  bind with product object :", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": true, "msg": err.Error()})
		return
	}

	err = helper.CheckProductValidation(productReq)
	if err != nil {
		log.Println("error invocked during  product :", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": true, "msg": err.Error()})
		return
	}

	product.Name = productReq.Name
	product.Description = productReq.Description
	product.ImageUrl = productReq.ImageUrl
	product.Price = productReq.Price
	product.MetaInfo = productReq.MetaInfo
	product.CreatedAt = time.Now().Unix()
	product.UpdatedAt = time.Now().Unix()

	id, err := dbs.Mgr.InsertOne(product, constval.ProductCollection)
	if err != nil {
		log.Println("error invocked during insert product data :", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": true, "msg": err.Error()})
		return
	}
	product.ID = id.(primitive.ObjectID)
	c.JSON(http.StatusOK, gin.H{"error": true, "msg": "success", "data": product})

}

// list the  products
func ListProducts(c *gin.Context) {
	page := c.DefaultQuery("page", "1")
	limit := c.DefaultQuery("limit", "10")
	offset := c.DefaultQuery("offset", "0")

	pageInt, _ := helper.ConvertStrInt(page)
	limitInt, _ := helper.ConvertStrInt(limit)
	offsetInt, _ := helper.ConvertStrInt(offset)

	productResp, count, err := dbs.Mgr.GetListProduct(pageInt, limitInt, offsetInt, constval.ProductCollection)
	if err != nil {
		log.Println("error invocked during get list product data :", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": true, "msg": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"error": false, "msg": "success", "data": map[string]interface{}{
		"products":   productResp,
		"totalCount": count,
	}})
}

// search product
func SearchProduct(c *gin.Context) {
	page := c.DefaultQuery("page", "1")
	limit := c.DefaultQuery("limit", "10")
	offset := c.DefaultQuery("offset", "0")
	s := c.Query("search")

	pageInt, _ := helper.ConvertStrInt(page)
	limitInt, _ := helper.ConvertStrInt(limit)
	offsetInt, _ := helper.ConvertStrInt(offset)

	productResp, count, err := dbs.Mgr.SearchProduct(pageInt, limitInt, offsetInt, s, constval.ProductCollection)
	if err != nil {
		log.Println("error invocked during search  product data :", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": true, "msg": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"error": false, "msg": "success", "data": map[string]interface{}{
		"products":   productResp,
		"totalCount": count,
	}})
}

// update product
func UpdateProduct(c *gin.Context) {
	userEmail, ok := c.Get("email")
	if !ok {
		log.Println("error was invoked  during get email :", userEmail)
		c.JSON(http.StatusInternalServerError, gin.H{"error": true, "msg": constval.NotRegisteredUser})
		return
	}

	userResp := dbs.Mgr.GetSingleRecordByEmailForUser(userEmail.(string), constval.UserCollection)
	if userResp.UserType != constval.AdminUser {
		log.Println("error was invoked  due to user not authorize ")
		c.JSON(http.StatusInternalServerError, gin.H{"error": true, "msg": constval.NotAuthorizedUserError})
		return
	}

	var updateProduct payloads.UpdateProduct
	var product payloads.Product
	err := c.BindJSON(&updateProduct)
	if err != nil {
		log.Println("error was invoked  during bind with product data :", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": true, "msg": err.Error()})
		return
	}
	pId, err := primitive.ObjectIDFromHex(updateProduct.ID)
	if err != nil {
		log.Println("error was invoked  during parse user object id :", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": true, "msg": err.Error()})
		return
	}
	pResp, err := dbs.Mgr.GetSingleProductById(pId, constval.ProductCollection)
	if err != nil {
		log.Println("error was invoked  during get product :", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": true, "msg": err.Error()})
		return
	}

	if pResp.Name == "" {
		log.Println("error was invoked due to emplty product")
		c.JSON(http.StatusBadRequest, gin.H{"error": true, "msg": constval.NoProductAvaliable})
		return
	}
	product.ID = pId
	product.Description = pResp.Description
	product.Name = pResp.Name
	product.Price = pResp.Price
	product.MetaInfo = pResp.MetaInfo
	product.ImageUrl = pResp.ImageUrl
	product.CreatedAt = pResp.CreatedAt
	product.UpdatedAt = time.Now().Unix()

	if updateProduct.Name != "" {
		product.Name = updateProduct.Name
	}

	if updateProduct.Description != "" {
		product.Description = updateProduct.Description
	}

	if updateProduct.Price > 0 {
		product.Price = updateProduct.Price
	}
	err = dbs.Mgr.UpdateProduct(product, constval.ProductCollection)
	if err != nil {
		log.Println("error was invoked  during update product :", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": true, "msg": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"error": false, "msg": "success", "data": product})

}

// delete product by user email
func DelProduct(c *gin.Context) {
	userEmail, ok := c.Get("email")
	if !ok {
		log.Println("error was invoked  during get email :", userEmail)
		c.JSON(http.StatusInternalServerError, gin.H{"error": true, "msg": constval.NotRegisteredUser})
		return
	}

	userResp := dbs.Mgr.GetSingleRecordByEmailForUser(userEmail.(string), constval.UserCollection)
	if userResp.UserType != constval.AdminUser {
		log.Println("error was invoked  due to user not authorize ")
		c.JSON(http.StatusInternalServerError, gin.H{"error": true, "msg": constval.NotAuthorizedUserError})
		return
	}

	id := c.Query("id")
	pId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		log.Println("error was invoked during get product id :", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": true, "msg": err.Error()})
		return
	}

	pResp, err := dbs.Mgr.GetSingleProductById(pId, constval.UserCollection)
	if err != nil {
		log.Println("error was invoked during get product data :", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": true, "msg": err.Error()})
		return
	}

	if pResp.Name == "" {
		log.Println("error was invoked due to product not avaliable")
		c.JSON(http.StatusBadRequest, gin.H{"error": true, "msg": constval.NoProductAvaliable})
		return
	}

	err = dbs.Mgr.DeleteProduct(pId, constval.ProductCollection)
	if err != nil {
		log.Println("error was invoked during delete product :", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": true, "msg": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"error": false, "msg": "success"})
}

// checkout  product by id
func CheckoutProductById(c *gin.Context) {
	cartIdStr := c.Param("id")
	cartId, err := primitive.ObjectIDFromHex(cartIdStr)
	if err != nil {
		log.Println("error was invoked during parse cart object id:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": true, "msg": err.Error()})
		return
	}

	{
		s := dbs.Mgr
		println("dbs manager :", s)
	}
	cart, err := dbs.Mgr.GetCartObjectById(cartId, constval.CartCollection)
	if err != nil {
		log.Println("error was invoked during get cart data drom db :", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": true, "msg": err.Error()})
		return
	}

	cart.Checkout = true

	err = dbs.Mgr.UpdateCartToCheckout(cart, constval.CartCollection)
	if err != nil {
		log.Println("error was invoked during update cart data in drom db :", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": true, "msg": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"error": false, "msg": "success"})
}

package users

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/willqiang/bookstore_oauth-go/oauth"
	"github.com/willqiang/bookstore_users-api/domain/users"
	"github.com/willqiang/bookstore_users-api/services"
	"github.com/willqiang/bookstore_users-api/utils/errors"
	"net/http"
	"strconv"
)

func Creat(c *gin.Context) {
	var user users.User
	if err := c.ShouldBindJSON(&user); err != nil {
		restErr := errors.NewBadRequestError("invalid json body")
		c.JSON(restErr.Status, restErr)
		return
	}

	result, saveErr := services.UsersService.CreateUser(user)
	if saveErr != nil {
		c.JSON(saveErr.Status, saveErr)
		return
	}

	fmt.Println(user)
	c.JSON(http.StatusCreated, result.Marshal(c.GetHeader("X-Public") == "true"))
}

func Get(c *gin.Context) {
	// Authenticate the access token
	if err := oauth.AuthenticateRequest(c.Request); err != nil {
		c.JSON(err.Status, err)
		return
	}
	// if we need force every access is a private access (the access token must be authenticated), uncomment follow code
	/*if oauth.GetCallerId(c.Request) == 0 {
		err := errors.RestErr{
			Status: http.StatusUnauthorized,
			Message: "resource not available",
		}
		c.JSON(err.Status, err)
		return
	}*/
	userId, idErr := getUserId(c.Param("user_id"))
	if idErr != nil {
		c.JSON(idErr.Status, idErr)
		return
	}

	user, getErr := services.UsersService.GetUser(userId)
	if getErr != nil {
		c.JSON(getErr.Status, getErr)
		return
	}
	if oauth.GetCallerId(c.Request) == user.Id {
		c.JSON(http.StatusOK, user.Marshal(false))
		return
	}
	c.JSON(http.StatusOK, user.Marshal(oauth.IsPublic(c.Request)))
}

func Update(c *gin.Context) {
	userId, idErr := getUserId(c.Param("user_id"))
	if idErr != nil {
		c.JSON(idErr.Status, idErr)
		return
	}

	var user users.User
	if err := c.ShouldBindJSON(&user); err != nil {
		restErr := errors.NewBadRequestError("invalid json body")
		c.JSON(restErr.Status, restErr)
		return
	}
	user.Id = userId

	isPartial := c.Request.Method == http.MethodPatch

	result, err := services.UsersService.UpdateUser(user, isPartial)
	if err != nil {
		c.JSON(err.Status, err)
		return
	}

	c.JSON(http.StatusOK, result.Marshal(c.GetHeader("X-Public") == "true"))
}

func Delete(c *gin.Context) {
	userId, idErr := getUserId(c.Param("user_id"))
	if idErr != nil {
		c.JSON(idErr.Status, idErr)
		return
	}

	if err := services.UsersService.DeleteUser(userId); err != nil {
		c.JSON(err.Status, err)
		return
	}
	c.JSON(http.StatusOK, map[string]string{"status": "deleted"})
}

func getUserId(userIdParam string) (int64, *errors.RestErr) {
	userId, err := strconv.ParseInt(userIdParam, 10, 64)
	if err != nil {
		userErr := errors.NewBadRequestError("user id must be a number")
		return 0, userErr
	}
	return userId, nil
}

// GET /internal/users/search?status=active
/*
获取status的值（以?开始的参数）和前面几个controller获取user_id的值不同
user_id 被称为parameter
status 被称为 query parameter
*/
func Search(c *gin.Context) {
	status := c.Query("status")
	users, err := services.UsersService.Search(status)
	if err != nil {
		c.JSON(err.Status, err)
		return
	}
	c.JSON(http.StatusOK, users.Marshal(c.GetHeader("X-Public") == "true"))
}

func Login(c *gin.Context) {
	var request users.LoginRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		restErr := errors.NewBadRequestError("invalid json body")
		c.JSON(restErr.Status, restErr)
		return
	}

	user, err := services.UsersService.LoginUser(request)
	if err != nil {
		c.JSON(err.Status, err)
		return
	}
	c.JSON(http.StatusOK, user.Marshal(c.GetHeader("X-Public") == "true"))
}

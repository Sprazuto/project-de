package controllers

import (
	"net/http"

	"github.com/Massad/gin-boilerplate/db"
	"github.com/Massad/gin-boilerplate/forms"
	"github.com/Massad/gin-boilerplate/models"
	"github.com/go-playground/validator/v10"

	"github.com/gin-gonic/gin"
)

// UserController ...
type UserController struct{}

var userModel = new(models.UserModel)
var userForm = new(forms.UserForm)

// getUserID ...
func getUserID(c *gin.Context) (userID int64) {
	//MustGet returns the value for the given key if it exists, otherwise it panics.
	return c.MustGet("userID").(int64)
}

// Login User godoc
// @Summary Login User example
// @Schemes
// @Description Login User example
// @Tags User
// @Accept json
// @Produce json
// @Param login body forms.LoginForm true "User"
// @Success 	 200  {object}  models.UserLoginResponse
// @Failure      406  {object}  models.MessageResponse
// @Router /user/login [post]
func (ctrl UserController) Login(c *gin.Context) {
	var loginForm forms.LoginForm

	if validationErr := c.ShouldBindJSON(&loginForm); validationErr != nil {
		message := userForm.Login(validationErr)
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": message})
		return
	}

	// Custom validation: ensure at least one of email or username is provided
	if loginForm.Email == "" && loginForm.Username == "" {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": "Email or username is required"})
		return
	}

	user, token, err := userModel.Login(loginForm)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Successfully logged in", "user": user, "token": token})
}

// Register User godoc
// @Summary Register User example
// @Schemes
// @Description Register User example
// @Tags User
// @Accept json
// @Produce json
// @Param register body forms.RegisterForm true "User"
// @Success 	 200  {object}  models.UserLoginResponse
// @Failure      406  {object}  models.MessageResponse
// @Router /user/register [post]
func (ctrl UserController) Register(c *gin.Context) {
	var registerForm forms.RegisterForm

	if validationErr := c.ShouldBindJSON(&registerForm); validationErr != nil {
		message := userForm.Register(validationErr)
		c.AbortWithStatusJSON(http.StatusNotAcceptable, gin.H{"message": message})
		return
	}

	user, err := userModel.Register(registerForm)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusNotAcceptable, gin.H{"message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Successfully registered", "user": user})
}

// Logout User godoc
// @Summary Logout User example
// @Schemes
// @Description Logout User example
// @Tags User
// @Accept json
// @Produce json
// @Success 	 200  {object}  models.MessageResponse
// @Failure      406  {object}  models.MessageResponse
// @Security BearerAuth
// @Router /user/logout [GET]
func (ctrl UserController) Logout(c *gin.Context) {

	au, err := authModel.ExtractTokenMetadata(c.Request)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": "User not logged in"})
		return
	}

	deleted, delErr := authModel.DeleteAuth(au.AccessUUID)
	if delErr != nil || deleted == 0 { //if any goes wrong
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "Invalid request"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Successfully logged out"})
}

// GetProfile godoc
// @Summary Get User Profile
// @Schemes
// @Description Get current user profile
// @Tags User
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} models.User
// @Failure 401 {object} models.MessageResponse
// @Router /user/profile [get]
func (ctrl UserController) GetProfile(c *gin.Context) {
	userID := getUserID(c)
	user, err := userModel.One(userID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "User not found"})
		return
	}

	c.JSON(http.StatusOK, user)
}

// ForgotPassword godoc
// @Summary Forgot Password
// @Schemes
// @Description Send password reset link to email
// @Tags User
// @Accept json
// @Produce json
// @Param forgot body forms.ForgotPasswordForm true "Email"
// @Success 200 {object} models.MessageResponse
// @Failure 400 {object} models.MessageResponse
// @Router /user/forgot-password [post]
func (ctrl UserController) ForgotPassword(c *gin.Context) {
	var form forms.ForgotPasswordForm

	if validationErr := c.ShouldBindJSON(&form); validationErr != nil {
		message := userForm.Email(validationErr.(validator.ValidationErrors)[0].Tag())
		c.AbortWithStatusJSON(http.StatusNotAcceptable, gin.H{"message": message})
		return
	}

	// Check if email exists
	checkEmail, err := db.GetDB().SelectInt(`SELECT count(id) FROM public."user" WHERE LOWER(email)=LOWER($1) LIMIT 1`, form.Email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Something went wrong"})
		return
	}
	if checkEmail == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Email not found"})
		return
	}

	// Stub for sending email
	// In production, generate token, save to db, send email with link
	c.JSON(http.StatusOK, gin.H{"message": "Password reset link sent to your email"})
}

// AssignRole godoc
// @Summary Assign Role to User
// @Schemes
// @Description Assign a role to a user
// @Tags User
// @Accept json
// @Produce json
// @Param assign body forms.AssignRoleForm true "User ID and Role Name"
// @Success 200 {object} models.MessageResponse
// @Failure 400 {object} models.MessageResponse
// @Security BearerAuth
// @Router /user/assign-role [post]
func (ctrl UserController) AssignRole(c *gin.Context) {
	var form forms.AssignRoleForm

	if validationErr := c.ShouldBindJSON(&form); validationErr != nil {
		c.AbortWithStatusJSON(http.StatusNotAcceptable, gin.H{"message": "Invalid form"})
		return
	}

	// Get role ID
	getDb := db.GetDB()
	var roleID int64
	_, err := getDb.Get(&roleID, `SELECT id FROM public.roles WHERE LOWER(name) = LOWER($1) LIMIT 1`, form.RoleName)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Role not found"})
		return
	}

	// Check if already assigned
	check, err := getDb.SelectInt(`SELECT count(id) FROM public.user_roles WHERE user_id = $1 AND role_id = $2`, form.UserID, roleID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Something went wrong"})
		return
	}
	if check > 0 {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Role already assigned"})
		return
	}

	// Assign role
	_, err = getDb.Exec(`INSERT INTO public.user_roles (user_id, role_id) VALUES ($1, $2)`, form.UserID, roleID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to assign role"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Role assigned successfully"})
}

// CreatePermission godoc
// @Summary Create Permission
// @Schemes
// @Description Create a new permission
// @Tags User
// @Accept json
// @Produce json
// @Param create body string true "Permission Name"
// @Success 200 {object} models.MessageResponse
// @Failure 400 {object} models.MessageResponse
// @Security BearerAuth
// @Router /permission/create [post]
func (ctrl UserController) CreatePermission(c *gin.Context) {
	var permName string
	if err := c.ShouldBindJSON(&permName); err != nil {
		c.AbortWithStatusJSON(http.StatusNotAcceptable, gin.H{"message": "Invalid form"})
		return
	}

	// Check if permission exists
	getDb := db.GetDB()
	check, err := getDb.SelectInt(`SELECT count(id) FROM public.permissions WHERE LOWER(name) = LOWER($1) LIMIT 1`, permName)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Something went wrong"})
		return
	}
	if check > 0 {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Permission already exists"})
		return
	}

	// Create permission
	_, err = getDb.Exec(`INSERT INTO public.permissions (name) VALUES ($1)`, permName)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to create permission"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Permission created successfully"})
}

package models

import (
	"errors"
	"time"

	"github.com/Massad/gin-boilerplate/db"
	"github.com/Massad/gin-boilerplate/forms"

	"golang.org/x/crypto/bcrypt"
)

type UserLoginResponse struct {
	User    User   `json:"user"`
	Token   Token  `json:"token"`
	Message string `json:"message"`
}

type MessageResponse struct {
	Message string `json:"message"`
}

// User ...
type User struct {
	ID             int64  `db:"id, primarykey, autoincrement" json:"id"`
	Email          string `db:"email" json:"email"`
	Username       string `db:"username" json:"username"`
	Password       string `db:"password" json:"-"`
	Name           string `db:"name" json:"name"`
	FailedAttempts int64  `db:"failed_attempts" json:"-"`
	LockedUntil    int64  `db:"locked_until" json:"-"`
	UpdatedAt      int64  `db:"updated_at" json:"-"`
	CreatedAt      int64  `db:"created_at" json:"-"`
}

func (u User) TableName() string {
	return "user"
}

// Role ...
type Role struct {
	ID        int64  `db:"id" json:"id"`
	Name      string `db:"name" json:"name"`
	UpdatedAt int64  `db:"updated_at" json:"-"`
	CreatedAt int64  `db:"created_at" json:"-"`
}

// Permission ...
type Permission struct {
	ID        int64  `db:"id" json:"id"`
	Name      string `db:"name" json:"name"`
	UpdatedAt int64  `db:"updated_at" json:"-"`
	CreatedAt int64  `db:"created_at" json:"-"`
}

// UserRole ...
type UserRole struct {
	ID        int64 `db:"id" json:"id"`
	UserID    int64 `db:"user_id" json:"user_id"`
	RoleID    int64 `db:"role_id" json:"role_id"`
	UpdatedAt int64 `db:"updated_at" json:"-"`
	CreatedAt int64 `db:"created_at" json:"-"`
}

// RolePermission ...
type RolePermission struct {
	ID           int64 `db:"id" json:"id"`
	RoleID       int64 `db:"role_id" json:"role_id"`
	PermissionID int64 `db:"permission_id" json:"permission_id"`
	UpdatedAt    int64 `db:"updated_at" json:"-"`
	CreatedAt    int64 `db:"created_at" json:"-"`
}

// LoginAttempt ...
type LoginAttempt struct {
	ID          int64 `db:"id" json:"id"`
	UserID      int64 `db:"user_id" json:"user_id"`
	Success     bool  `db:"success" json:"success"`
	AttemptTime int64 `db:"attempt_time" json:"attempt_time"`
	UpdatedAt   int64 `db:"updated_at" json:"-"`
	CreatedAt   int64 `db:"created_at" json:"-"`
}

// UserModel ...
type UserModel struct{}

var authModel = new(AuthModel)

// Login ...
func (m UserModel) Login(form forms.LoginForm) (user User, token Token, err error) {
	getDb := db.GetDB()
	currentTime := time.Now().Unix()

	// Query by email or username
	row := getDb.Db.QueryRow(`SELECT id, email, username, password, name, failed_attempts, locked_until, updated_at, created_at FROM public."user" WHERE (LOWER(email)=LOWER($1) OR LOWER(username)=LOWER($2)) LIMIT 1`, form.Email, form.Username)
	err = row.Scan(&user.ID, &user.Email, &user.Username, &user.Password, &user.Name, &user.FailedAttempts, &user.LockedUntil, &user.UpdatedAt, &user.CreatedAt)
	if err != nil {
		return user, token, errors.New("invalid login details")
	}

	// Check if account is locked
	if user.LockedUntil > currentTime {
		return user, token, errors.New("account is locked. Try again later")
	}

	// Compare password
	bytePassword := []byte(form.Password)
	byteHashedPassword := []byte(user.Password)
	err = bcrypt.CompareHashAndPassword(byteHashedPassword, bytePassword)
	if err != nil {
		// Log failed attempt
		m.LogLoginAttempt(user.ID, false)

		// Increment failed attempts
		failed := user.FailedAttempts + 1
		if failed >= 5 {
			lockUntil := time.Now().Add(1 * time.Minute).Unix()
			_, updateErr := getDb.Exec(`UPDATE public."user" SET failed_attempts=$1, locked_until=$2 WHERE id=$3`, failed, lockUntil, user.ID)
			if updateErr != nil {
				return user, token, updateErr
			}
		} else {
			_, updateErr := getDb.Exec(`UPDATE public."user" SET failed_attempts=$1 WHERE id=$2`, failed, user.ID)
			if updateErr != nil {
				return user, token, updateErr
			}
		}
		return user, token, errors.New("invalid login details")
	}

	// Success: reset failed attempts, log success
	_, resetErr := getDb.Exec(`UPDATE public."user" SET failed_attempts=0, locked_until=0 WHERE id=$1`, user.ID)
	if resetErr != nil {
		return user, token, resetErr
	}
	m.LogLoginAttempt(user.ID, true)

	// Generate token
	tokenDetails, err := authModel.CreateToken(user.ID)
	if err != nil {
		return user, token, err
	}

	saveErr := authModel.CreateAuth(user.ID, tokenDetails)
	if saveErr == nil {
		token.AccessToken = tokenDetails.AccessToken
		token.RefreshToken = tokenDetails.RefreshToken
	}

	return user, token, nil
}

// Register ...
func (m UserModel) Register(form forms.RegisterForm) (user User, err error) {
	getDb := db.GetDB()

	// Check if email or username exists
	checkEmail, err := getDb.SelectInt(`SELECT count(id) FROM public."user" WHERE LOWER(email)=LOWER($1) LIMIT 1`, form.Email)
	if err != nil {
		return user, errors.New("something went wrong, please try again later")
	}
	checkUsername, err := getDb.SelectInt(`SELECT count(id) FROM public."user" WHERE LOWER(username)=LOWER($1) LIMIT 1`, form.Username)
	if err != nil {
		return user, errors.New("something went wrong, please try again later")
	}
	if checkEmail > 0 {
		return user, errors.New("email already exists")
	}
	if checkUsername > 0 {
		return user, errors.New("username already exists")
	}

	bytePassword := []byte(form.Password)
	hashedPassword, err := bcrypt.GenerateFromPassword(bytePassword, bcrypt.DefaultCost)
	if err != nil {
		return user, errors.New("something went wrong, please try again later")
	}

	// Create user
	row := getDb.QueryRow(`INSERT INTO public."user"(email, username, password, name, failed_attempts, locked_until) VALUES($1, $2, $3, $4, 0, 0) RETURNING id`, form.Email, form.Username, string(hashedPassword), form.Name)
	err = row.Scan(&user.ID)
	if err != nil {
		return user, errors.New("something went wrong, please try again later")
	}

	user.Name = form.Name
	user.Email = form.Email
	user.Username = form.Username

	// Assign default 'user' role
	roleID, err := m.getOrCreateRoleID("user")
	if err != nil {
		return user, err
	}
	_, err = getDb.Exec(`INSERT INTO public.user_roles (user_id, role_id) VALUES ($1, $2)`, user.ID, roleID)
	if err != nil {
		return user, err
	}

	return user, nil
}

// One ...
func (m UserModel) One(userID int64) (user User, err error) {
	row := db.GetDB().Db.QueryRow(`SELECT id, email, username, name, failed_attempts, locked_until FROM public."user" WHERE id=$1 LIMIT 1`, userID)
	err = row.Scan(&user.ID, &user.Email, &user.Username, &user.Name, &user.FailedAttempts, &user.LockedUntil)
	return user, err
}

// LogLoginAttempt ...
func (m UserModel) LogLoginAttempt(userID int64, success bool) error {
	_, err := db.GetDB().Exec(`INSERT INTO public.login_attempts (user_id, success, attempt_time) VALUES ($1, $2, $3)`, userID, success, time.Now().Unix())
	return err
}

// GetUserRoles ...
func (m UserModel) GetUserRoles(userID int64) ([]Role, error) {
	var roles []Role
	_, err := db.GetDB().Select(&roles, "SELECT r.* FROM public.roles r JOIN public.user_roles ur ON r.id = ur.role_id WHERE ur.user_id = $1", userID)
	return roles, err
}

// HasPermission ...
func (m UserModel) HasPermission(userID int64, permName string) (bool, error) {
	var count int
	_, err := db.GetDB().Get(&count, `SELECT COUNT(*) FROM public.permissions p
		JOIN public.role_permissions rp ON p.id = rp.permission_id
		JOIN public.user_roles ur ON rp.role_id = ur.role_id
		WHERE ur.user_id = $1 AND p.name = $2`, userID, permName)
	return count > 0, err
}

// getOrCreateRoleID ...
func (m UserModel) getOrCreateRoleID(roleName string) (int64, error) {
	getDb := db.GetDB()

	// Check if role exists
	var roleID int64
	_, err := getDb.Get(&roleID, "SELECT id FROM public.roles WHERE LOWER(name) = LOWER($1) LIMIT 1", roleName)
	if err == nil {
		return roleID, nil
	}

	// Create role
	row := getDb.QueryRow("INSERT INTO public.roles (name) VALUES ($1) RETURNING id", roleName)
	err = row.Scan(&roleID)
	if err != nil {
		return 0, err
	}
	return roleID, nil
}

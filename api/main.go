package main

import (
	"log"
	"net/http"
	"os"
	"runtime"

	"github.com/Massad/gin-boilerplate/controllers"
	"github.com/Massad/gin-boilerplate/db"
	_ "github.com/Massad/gin-boilerplate/docs"
	"github.com/Massad/gin-boilerplate/forms"
	"github.com/Massad/gin-boilerplate/models"
	"github.com/gin-contrib/gzip"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/google/uuid"
	"github.com/joho/godotenv"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger" // gin-swagger middleware
	"golang.org/x/crypto/bcrypt"
)

// RequestIDMiddleware ...
// Generate a unique ID and attach it to each request for future reference or use
func RequestIDMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		uuid := uuid.New()
		c.Writer.Header().Set("X-Request-Id", uuid.String())
		c.Next()
	}
}

// TokenAuthMiddleware ...
// JWT Authentication middleware attached to each request that needs to be authenticated to validate the access_token in the header
func TokenAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		auth := new(controllers.AuthController)
		auth.TokenValid(c)
		c.Next()
	}
}

// @title           Golang Gin Boilerplate
// @version         1.0
// @description     A RESTful API boilerplate with Gin Framework, PostgreSQL, Redis and JWT authentication
// @termsOfService  http://swagger.io/terms/

// @contact.name   API Support
// @contact.url    http://www.swagger.io/support
// @contact.email  support@swagger.io

// @license.name  MIT License
// @license.url   https://github.com/Massad/gin-boilerplate/blob/master/LICENSE

// @host      localhost:9000
// @BasePath  /v1

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization

// @externalDocs.description  OpenAPI
// @externalDocs.url          https://swagger.io/resources/open-api/

// setupCORS configures CORS middleware with proper security
func setupCORS() gin.HandlerFunc {
	return func(c *gin.Context) {
		origin := c.Request.Header.Get("Origin")
		env := os.Getenv("ENV")

		var allowedOrigins []string
		if env == "PRODUCTION" {
			productionDomain := os.Getenv("FRONTEND_DOMAIN")
			if productionDomain != "" {
				allowedOrigins = []string{
					"https://" + productionDomain,
					"http://" + productionDomain,
				}
			} else {
				allowedOrigins = []string{}
			}
		} else {
			allowedOrigins = []string{
				"http://localhost:5173",
				"http://localhost:8000",
				"http://localhost:3000",
				"http://127.0.0.1:5173",
				"http://127.0.0.1:3000",
			}
		}

		allowOrigin := ""
		for _, allowed := range allowedOrigins {
			if origin == allowed {
				allowOrigin = allowed
				break
			}
		}

		if allowOrigin != "" {
			c.Writer.Header().Set("Access-Control-Allow-Origin", allowOrigin)
		}

		c.Writer.Header().Set("Access-Control-Max-Age", "86400")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE, UPDATE")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "X-Requested-With, Content-Type, Origin, Authorization, Accept, Client-Security-Token, Accept-Encoding, x-access-token")
		c.Writer.Header().Set("Access-Control-Expose-Headers", "Content-Length")

		if allowOrigin != "" {
			c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		}

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(200)
		} else {
			c.Next()
		}
	}
}

func main() {
	//Load the .env file
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("error: failed to load the env file")
	}

	if os.Getenv("ENV") == "PRODUCTION" {
		gin.SetMode(gin.ReleaseMode)
	}

	//Start the default gin server
	r := gin.Default()

	//Custom form validator
	binding.Validator = new(forms.DefaultValidator)

	r.Use(setupCORS())
	r.Use(RequestIDMiddleware())
	r.Use(gzip.Gzip(gzip.DefaultCompression))

	//Start PostgreSQL database
	//Example: db.GetDB() - More info in the models folder
	db.Init()

	// Run migrations
	err = models.RunMigrations()
	if err != nil {
		log.Fatal("Failed to run migrations:", err)
	}

	// Create SPSE tables
	err = models.CreateSPASETables()
	if err != nil {
		log.Printf("Warning: Failed to create SPSE tables: %v", err)
	}

	// Seed admin user if not exists
	getDb := db.GetDB()
	adminCount, err := getDb.SelectInt(`SELECT count(id) FROM public."user" WHERE LOWER(username)=LOWER($1)`, "admin")
	if err != nil {
		log.Fatal("Failed to check admin user:", err)
	}
	if adminCount == 0 {
		bytePassword := []byte("Admin777")
		hashedPassword, hashErr := bcrypt.GenerateFromPassword(bytePassword, bcrypt.DefaultCost)
		if hashErr != nil {
			log.Fatal("Failed to hash admin password:", hashErr)
		}

		var userID int64
		err = getDb.QueryRow(`INSERT INTO public."user" (email, username, password, name, failed_attempts, locked_until) VALUES ($1, $2, $3, $4, 0, 0) RETURNING id`, "admin@project.de", "admin", string(hashedPassword), "Administrator").Scan(&userID)
		if err != nil {
			log.Fatal("Failed to create admin user:", err)
		}

		// Get or create admin role
		var roleID int64
		roleCount, err := getDb.SelectInt(`SELECT count(id) FROM public.roles WHERE LOWER(name) = LOWER($1)`, "admin")
		if err != nil {
			log.Fatal("Failed to check admin role:", err)
		}
		if roleCount == 0 {
			err = getDb.QueryRow("INSERT INTO public.roles (name) VALUES ($1) RETURNING id", "admin").Scan(&roleID)
			if err != nil {
				log.Fatal("Failed to create admin role:", err)
			}
		} else {
			roleID, err = getDb.SelectInt(`SELECT id FROM public.roles WHERE LOWER(name) = LOWER($1) LIMIT 1`, "admin")
			if err != nil {
				log.Fatal("Failed to get admin role ID:", err)
			}
		}

		// Assign role to user
		_, assignErr := getDb.Exec(`INSERT INTO public.user_roles (user_id, role_id) VALUES ($1, $2)`, userID, roleID)
		if assignErr != nil {
			log.Fatal("Failed to assign admin role:", assignErr)
		}

	}

	//Start Redis on database 1 - it's used to store the JWT but you can use it for anythig else
	//Example: db.GetRedis().Set(KEY, VALUE, at.Sub(now)).Err()
	db.InitRedis(1)

	v1 := r.Group("/v1")
	{
		/*** START USER ***/
		user := new(controllers.UserController)
		auth := new(controllers.AuthController)

		v1.POST("/user/login", user.Login)
		v1.POST("/user/register", user.Register)
		v1.GET("/user/logout", user.Logout)
		v1.GET("/user/profile", TokenAuthMiddleware(), user.GetProfile)
		v1.POST("/user/forgot-password", user.ForgotPassword)
		v1.POST("/user/assign-role", TokenAuthMiddleware(), auth.HasPermission("manage_users"), user.AssignRole)

		/*** START AUTH ***/
		//Refresh the token when needed to generate new access_token and refresh_token for the user
		v1.POST("/token/refresh", auth.Refresh)

		/*** START Permission ***/
		v1.POST("/permission/create", TokenAuthMiddleware(), auth.HasPermission("manage_users"), user.CreatePermission)

		/*** START Article ***/
		article := new(controllers.ArticleController)

		v1.POST("/article", TokenAuthMiddleware(), auth.HasPermission("write_article"), article.Create)
		v1.GET("/articles", TokenAuthMiddleware(), auth.HasPermission("read_article"), article.All)
		v1.GET("/article/:id", TokenAuthMiddleware(), auth.HasPermission("read_article"), article.One)
		v1.PUT("/article/:id", TokenAuthMiddleware(), auth.HasPermission("write_article"), article.Update)
		v1.DELETE("/article/:id", TokenAuthMiddleware(), auth.HasPermission("write_article"), article.Delete)

		/*** START Sijagur ***/
		sijagur := new(controllers.SijagurController)

		// Realisasi endpoints (existing)
		v1.GET("/realisasi-bulan", TokenAuthMiddleware(), sijagur.GetRealisasiBulan)
		v1.GET("/realisasi-tahun", TokenAuthMiddleware(), sijagur.GetRealisasiTahun)
		v1.GET("/realisasi-perbulan", TokenAuthMiddleware(), sijagur.GetRealisasiPerbulan)

		// Peringkat Kinerja (alias-based ranking, scoped by jenis_opd via ?scope=skpd|kecamatan)
		// Uses models.SijagurData.GetPeringkatKinerja and returns models.RankingResponse
		v1.GET("/sijagur/peringkat-kinerja", TokenAuthMiddleware(), sijagur.GetPeringkatKinerja)

		/*** START SPSE Procurement Scraper ***/
		spse := new(controllers.SPSEController)

		// Scraping endpoints (public access for testing)
		v1.POST("/spse/scraper/perencanaan", spse.ScrapePerencanaan)
		v1.POST("/spse/scraper/persiapan", spse.ScrapePersiapan)
		v1.POST("/spse/scraper/pemilihan", spse.ScrapePemilihan)
		v1.POST("/spse/scraper/hasilpemilihan", spse.ScrapeHasilPemilihan)
		v1.POST("/spse/scraper/kontrak", spse.ScrapeKontrak)
		v1.POST("/spse/scraper/serahterima", spse.ScrapeSerahTerima)
		v1.POST("/spse/scraper/all", spse.ScrapeAll)

		// Data retrieval endpoints (protected)
		v1.GET("/spse/statistics", spse.GetStatistics)
		v1.GET("/spse/data/perencanaan", spse.GetPerencanaan)
		v1.GET("/spse/data/persiapan", spse.GetPersiapan)
		v1.GET("/spse/data/pemilihan", spse.GetPemilihan)
		v1.GET("/spse/data/hasilpemilihan", spse.GetHasilPemilihan)
		v1.GET("/spse/data/kontrak", spse.GetKontrak)
		v1.GET("/spse/data/serahterima", spse.GetSerahTerima)
	}

	// Swagger docs
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))

	r.LoadHTMLGlob("./public/html/*")

	r.Static("/public", "./public")

	r.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", gin.H{
			"ginBoilerplateVersion": "v0.03",
			"goVersion":             runtime.Version(),
		})
	})

	r.NoRoute(func(c *gin.Context) {
		c.HTML(404, "404.html", gin.H{})
	})

	port := os.Getenv("PORT")

	log.Printf("\n\n PORT: %s \n ENV: %s \n SSL: %s \n Version: %s \n\n", port, os.Getenv("ENV"), os.Getenv("SSL"), os.Getenv("API_VERSION"))

	if os.Getenv("SSL") == "TRUE" {

		//Generated using sh generate-certificate.sh
		SSLKeys := &struct {
			CERT string
			KEY  string
		}{
			CERT: "./cert/myCA.cer",
			KEY:  "./cert/myCA.key",
		}

		r.RunTLS(":"+port, SSLKeys.CERT, SSLKeys.KEY)
	} else {
		r.Run(":" + port)
	}
}

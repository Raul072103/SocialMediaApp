package main

import (
	"SocialMediaApp/internal/auth"
	"SocialMediaApp/internal/db"
	"SocialMediaApp/internal/env"
	mailer2 "SocialMediaApp/internal/mailer"
	"SocialMediaApp/internal/store"
	cache2 "SocialMediaApp/internal/store/cache"
	"github.com/go-redis/redis/v8"
	"go.uber.org/zap"
	"time"
)

const version = "0.0.2"

//	@title			SocialMedia API
//	@description	API for a social media app, simulating real-world scenarios used by social media apps.
//	@termsOfService	http://swagger.io/terms/

//	@contact.name	API Support
//	@contact.url	http://www.swagger.io/support
//	@contact.email	support@swagger.io

//	@license.name	Apache 2.0
//	@license.url	http://www.apache.org/licenses/LICENSE-2.0.html

// @BasePath					/v1
//
// @securityDefinitions.apikey	ApiKeyAuth
// @in							header
// @name						Authorization
// @description
func main() {
	cfg := config{
		addr: env.GetString("ADDR", ":8080"),
		db: dbConfig{
			addr:         env.GetString("DB_ADDR", "postgres://admin:adminpassword@localhost:5434/social?sslmode=disable"),
			maxOpenConns: env.GetInt("DB_MAX_OPEN_CONNS", 30),
			maxIdleConns: env.GetInt("DB_MAX_IDLE_CONNS", 30),
			maxIdleTime:  env.GetString("DB_MAX_IDLE_TIME", "15m"),
		},
		redisCfg: redisConfig{
			addr:    env.GetString("REDIS_ADDR", "localhost:6379"),
			pw:      env.GetString("REDIS_PW", ""),
			db:      env.GetInt("REDIS_DB", 1),
			enabled: env.GetBool("REDIS_ENABLED", true),
		},
		env:    env.GetString("DEV", "development"),
		apiURL: env.GetString("EXTERNAL_URL", "localhost:8080"),
		mail: mailConfig{
			exp:       time.Hour * 24 * 3, // 3 days
			fromEmail: env.GetString("FROM_EMAIL", "raulsocialmedia@demomailtrap.com"),
			toEmail:   env.GetString("TO_EMAIL_DEFAULT", "raulandrei2019@gmail.com"),
			sendGrid: sendGridConfig{
				apiKey: env.GetString("SENDGRID_API_KEY", ""),
			},
			mailTrap: mailTrapConfig{
				apiKey: env.GetString("MAILTRAP_API_KEY", "39368e1ef343a7c84489ba5c81a79f94"),
			},
		},
		frontendURL: env.GetString("FRONTEND_URL", "http://localhost:5174"),
		auth: authConfig{
			basic: basicConfig{
				user: env.GetString("AUTH_BASIC_USER", "admin"),
				pass: env.GetString("AUTH_BASIC_PASS", "admin"),
			},
			token: tokenConfig{
				secret: env.GetString("AUTH_TOKEN_SECRET", "example"),
				exp:    time.Hour * 24 * 3, // 3 days
				iss:    "raulsocialmedia",
			},
		},
	}

	// Logger
	logger := zap.Must(zap.NewProduction()).Sugar()
	defer logger.Sync()

	// Database
	database, err := db.New(
		cfg.db.addr,
		cfg.db.maxOpenConns,
		cfg.db.maxIdleConns,
		cfg.db.maxIdleTime)

	defer database.Close()
	logger.Info("database connection pool established")

	// Cache
	var redisDB *redis.Client
	if cfg.redisCfg.enabled {
		redisDB = cache2.NewRedisClient(cfg.redisCfg.addr, cfg.redisCfg.pw, cfg.redisCfg.db)
		logger.Info("redis cache connection pool established")
	}

	if err != nil {
		logger.Fatal(err)
	}

	store := store.NewStorage(database)
	cacheStorage := cache2.NewRedisStorage(redisDB)

	// Mailer
	mailtrap, err := mailer2.NewMailTrapClient(cfg.mail.mailTrap.apiKey, cfg.mail.fromEmail, cfg.mail.toEmail)
	if err != nil {
		logger.Fatal(err)
	}

	jwtAuthenticator := auth.NewJWTAuthenticator(
		cfg.auth.token.secret,
		cfg.auth.token.iss,
		cfg.auth.token.iss)

	app := &application{
		config:        cfg,
		store:         store,
		cacheStorage:  cacheStorage,
		logger:        logger,
		mailer:        mailtrap,
		authenticator: jwtAuthenticator,
	}

	mux := app.mount()
	logger.Fatal(app.run(mux))
}

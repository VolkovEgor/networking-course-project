package yak

import (
	"github.com/architectv/networking-course-project/backend/pkg/handlers"
	"github.com/architectv/networking-course-project/backend/pkg/repositories"
	"github.com/architectv/networking-course-project/backend/pkg/repositories/postgres"
	"github.com/architectv/networking-course-project/backend/pkg/services"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

func CreateApp() {
	logrus.SetFormatter(new(logrus.JSONFormatter))

	if err := initConfig(); err != nil {
		logrus.Fatalf("error initializing configs: %s", err.Error())
	}

	// if err := godotenv.Load(); err != nil {
	// 	logrus.Fatalf("error loading env variables: %s", err.Error())
	// }

	// db, err := mongoDB.NewMongoDB()
	db, err := postgres.NewPostgresDB(postgres.Config{
		Host:     viper.GetString("db.host"),
		Port:     viper.GetString("db.port"),
		Username: viper.GetString("db.username"),
		DBName:   viper.GetString("db.dbname"),
		SSLMode:  viper.GetString("db.sslmode"),
		Password: viper.GetString("db.password"), //os.Getenv("DB_PASSWORD"),
	})
	if err != nil {
		logrus.Fatalf("failed to initialize db: %s", err.Error())
	}

	repos := repositories.NewRepository(db)
	services := services.NewService(repos)
	handlers := handlers.NewHandler(services)

	app := fiber.New()
	app.Use(logger.New())
	handlers.RegisterHandlers(app)
	app.Listen(viper.GetString("port"))
}

func initConfig() error {
	viper.AddConfigPath("config")
	viper.SetConfigName("config")
	return viper.ReadInConfig()
}

package cmd

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"gorm.io/gorm"

	"github.com/u2u-labs/layerg-crawler/db"
)

func startApi(cmd *cobra.Command, args []string) {
	logger, _ := zap.NewProduction()
	defer logger.Sync() // flushes buffer, if any
	sugar := logger.Sugar()

	gdb, err := db.NewCockroachDbClient(&db.DbConfig{
		Url:  viper.GetString("COCKROACH_DB_URL"),
		Name: viper.GetString("COCKROACH_DB_NAME"),
	})
	if err != nil {
		panic(err)
	}
	rdb, err := db.NewRedisClient(&db.RedisConfig{
		Url:      viper.GetString("REDIS_DB_URL"),
		Db:       viper.GetInt("REDIS_DB"),
		Password: viper.GetString("REDIS_DB_PASSWORD"),
	})
	if err != nil {
		panic(err)
	}

	serveApi(sugar, gdb, rdb)
}

func serveApi(sugar *zap.SugaredLogger, gdb *gorm.DB, rdb *redis.Client) error {

	// Create a default Gin router
	router := gin.Default()

	// Define a simple GET route
	router.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "Hello, World!",
		})
	})

	// Start the server on port 8080
	router.Run(":8085")
	return nil
}

module database-camp

go 1.17

require (
	cloud.google.com/go/storage v1.18.2
	github.com/go-redis/redis/v8 v8.11.4
	github.com/go-sql-driver/mysql v1.6.0
	github.com/gofiber/fiber/v2 v2.24.0
	github.com/golang-jwt/jwt v3.2.2+incompatible
	github.com/satori/go.uuid v1.2.0
	go.uber.org/zap v1.20.0
	golang.org/x/crypto v0.0.0-20210513164829-c07d793c2f9a
	golang.org/x/oauth2 v0.0.0-20211104180415-d3ed0bb246c8
	gorm.io/driver/mysql v1.2.2
	gorm.io/gorm v1.22.4
)

require github.com/joho/godotenv v1.4.0

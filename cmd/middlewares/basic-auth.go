package middleware

// ApiKeyAuth is a middleware for Basic Authentication
// func ApiKeyAuth(db *dbCon.Queries) gin.HandlerFunc {
// 	return func(c *gin.Context) {
// 		// Get the Authorization header

// 		apiHeaderKey := viper.GetString("API_KEY_HEADER_NAME")

// 		authApikey := c.GetHeader(apiHeaderKey)
// 		if authApikey == "" {
// 			c.JSON(http.StatusUnauthorized, gin.H{"error": "API key is required"})
// 			c.Abort()
// 			return
// 		}

// 		appId := viper.GetString("APP_ID")

// 		// Check if the API key is valid
// 		uuidAppId, err := uuid.Parse(appId)
// 		if err != nil {
// 			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid app Id format"})
// 			c.Abort()
// 			return
// 		}
// 		app, err := db.GetAppById(c, uuidAppId)
// 		if err != nil {
// 			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid API key"})
// 			c.Abort()
// 			return
// 		}

// 		// Check if the username and password are set
// 		if app.SecretKey != authApikey {
// 			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid API key"})
// 			c.Abort()
// 			return
// 		}

// 		// Continue down the chain to handler etc
// 		c.Next()
// 	}
// }

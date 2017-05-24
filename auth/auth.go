package auth

import (
	"errors"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

//Declaring errors and errors messages
var (
	unexpectedMethod = errors.New("Unexpected signing method")
	invalidTimeToken = "Token time invalid"
	timeOutToken     = "Token expired"
	invalidToken     = "The token is not valid"
)

//This function return two things,
//the first is a token with its respective time
//and second is a possible error in another case
func CreateToken(mail string) (string, error) {
	token := jwt.New(jwt.SigningMethodHS512)
	claims := token.Claims.(jwt.MapClaims)
	claims["mail"] = mail
	//Adding 30 days to expiration time
	claims["exp"] = time.Now().Add(time.Hour * 24 * 30).Unix()
	return token.SignedString(mySigningKey)
}

//This Middleware function,
//it has the purpose of validate a token to access to API
func ValidateToken() gin.HandlerFunc {
	log.Println("Validator JWT listening")
	return func(c *gin.Context) {
		//Parsing header request
		tokenArray := strings.Split(c.Request.Header["Authorization"][0], " ")
		tokenString := tokenArray[1]
		token, err := jwt.Parse(tokenString,
			func(token *jwt.Token) (interface{}, error) {
				// Validating the algorithm used
				if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
					return nil, unexpectedMethod
				}
				return mySigningKey, nil
			})
		//If parsing ending with error
		if err != nil {
			response := gin.H{
				"status":  "error",
				"data":    nil,
				"message": err.Error(),
			}
			c.JSON(http.StatusUnauthorized, response)
			c.AbortWithStatus(http.StatusUnauthorized)
		} else {
			//Obtains data claims
			if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
				expf, ok := claims["exp"].(float64)
				exp := int64(expf)
				//Obtaining the expiration time generates an error
				if !ok {
					response := gin.H{
						"status":  "error",
						"data":    nil,
						"message": invalidTimeToken,
					}
					c.JSON(http.StatusUnauthorized, response)
				}
				//Token time is expired
				if exp < time.Now().Unix() {
					response := gin.H{
						"status":  "error",
						"data":    nil,
						"message": timeOutToken,
					}
					c.JSON(http.StatusUnauthorized, response)
					c.Abort()
				} else {
					// Good case! :)
					c.Next()
				}
			} else {
				//Token creation error
				response := gin.H{
					"status":  "error",
					"data":    nil,
					"message": invalidToken,
				}
				c.JSON(http.StatusUnauthorized, response)
			}
		}

	}
}

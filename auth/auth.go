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

//my secret signature
var mySigningKey = []byte("3f37722f5feb721161c13df2acf554e5")

//This function return two things,
//the first is a token and second is a possible error
func CreateToken(mail string) (string, error) {
	token := jwt.New(jwt.SigningMethodHS512)
	claims := token.Claims.(jwt.MapClaims)

	claims["mail"] = mail
	claims["exp"] = time.Now().Add(time.Hour * 24 * 30).Unix()
	log.Println(time.Now().Unix())
	log.Println(claims["exp"])

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
				if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
					return nil, errors.New("Unexpected signing method")
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
						"message": "Token time invalid",
					}
					c.JSON(http.StatusUnauthorized, response)
				}
				//Token time is expired
				if exp < time.Now().Unix() {
					response := gin.H{
						"status":  "error",
						"data":    nil,
						"message": "Token expired",
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
					"message": "The token is not valid",
				}
				c.JSON(http.StatusUnauthorized, response)
			}
		}

	}
}

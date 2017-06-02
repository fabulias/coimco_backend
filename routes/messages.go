package routes

//Messages to return request information
var (
	GetMessageErrorPlural   = "There are no"
	GetMessageErrorSingular = "There is not"
	PostMessageError        = "Error inserting"
	ErrorParams             = "Error in query params"
	BindJson                = "Error binding json"
	LoginOK                 = "Mail and pass are correct, token it's OK"
	LoginError              = "Mail or pass aren't correct"
	TokenError              = "Error creating token"
	ErrorHashPassword       = "Error hashing password"
)

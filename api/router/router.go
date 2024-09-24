package router

import (
	"database/sql"
	"net/http"

	"github.com/thaironsilva/messenger/api/resource/user"
	"github.com/thaironsilva/messenger/cognitoClient"
)

func New(db *sql.DB) *http.ServeMux {
	router := http.NewServeMux()

	repository := user.NewRepository(db)
	cognito := cognitoClient.NewCognitoClient()
	userHandler := user.NewHandler(repository, cognito)
	router.HandleFunc("GET /user", user.GetUser(userHandler))
	router.HandleFunc("GET /users", user.GetUsers(userHandler))
	router.HandleFunc("POST /users", user.CreateUser(userHandler))
	router.HandleFunc("POST /users/confirmation", user.ConfirmAccount(userHandler))
	router.HandleFunc("POST /users/login", user.SignIn(userHandler))
	router.HandleFunc("PUT /users/password", user.UpdatePassword(userHandler))

	return router
}

// curl http://localhost:8080/users/login \
//     --include \
//     --header "Content-Type: application/json" \
//     --request "POST" \
//     --data '{"email": "thairon.ssilva@gmail.com","password": "Password!1"}'

// curl http://localhost:8080/user \
// 	--include \
// 	--header "Content-Type: application/json" \
// 	--header "Authorization: Bearer eyJraWQiOiJ4SlwvKzhUQXVJSXRXT05heEEyOGxOeGRGVE5KRTZVVzR3bzhDUFI2NE5GST0iLCJhbGciOiJSUzI1NiJ9.eyJzdWIiOiI5MTliMzU5MC0wMGMxLTcwNzItNzliOS1kODM3YTlhNDA3MWYiLCJpc3MiOiJodHRwczpcL1wvY29nbml0by1pZHAudXMtZWFzdC0yLmFtYXpvbmF3cy5jb21cL3VzLWVhc3QtMl9iR1RQTEZnTTciLCJjbGllbnRfaWQiOiJxaXFvZW1wOGhqbHQxYXIxcWYxMjMzbzd0Iiwib3JpZ2luX2p0aSI6IjhiY2ViMTViLWU3OTgtNDdmNy05ZmNmLThhYThlOGM5Y2E3ZSIsImV2ZW50X2lkIjoiNDY2NDRhZjItMWRkNC00YmY4LWE0NzQtZjE1ZDg3YjBlODQzIiwidG9rZW5fdXNlIjoiYWNjZXNzIiwic2NvcGUiOiJhd3MuY29nbml0by5zaWduaW4udXNlci5hZG1pbiIsImF1dGhfdGltZSI6MTcyNzE4ODg4MCwiZXhwIjoxNzI3MTkyNDgwLCJpYXQiOjE3MjcxODg4ODAsImp0aSI6IjI4ZmZlMzU2LTEzZjQtNGE0Yi04YWQwLWY3MDJlMDRlM2M5NiIsInVzZXJuYW1lIjoiOTE5YjM1OTAtMDBjMS03MDcyLTc5YjktZDgzN2E5YTQwNzFmIn0.RaYgkL_XPSi6miznkM9e7bflYyy_pktsQQreFjv9rq3yQazxeTNFIg_v_-7aDwHKS8t9D5eFjGq4H7eNKcX1MTc3Trgh6rpOCX14yE8FL4wJYxboA1uKvEoxK7JL_liTe0SzcsgDo5zAfcPCUu2AACn1Q45anX-_rxfIt5XfaGXyxFJ1mxVmeGSTeJNapwl7u0WUpBAc39z4YZS0MzggSBsfW4UoV_O0Xi6Z4ieK8rW720GOnWZdfyQnmG1oRcCH80hEtpQoWFKjARitV5sWbegeeYw60DDSjMSnygVHlpYI02mM-pqo4Jm4iTXTKlSdSO3RKLJDQ961VEWJnB02WQ" \
// 	--request "GET"

// curl http://localhost:8080/users/password \
// 	--include \
// 	--header "Content-Type: application/json" \
// 	--header "Authorization: Bearer eyJraWQiOiJ4SlwvKzhUQXVJSXRXT05heEEyOGxOeGRGVE5KRTZVVzR3bzhDUFI2NE5GST0iLCJhbGciOiJSUzI1NiJ9.eyJzdWIiOiI5MTliMzU5MC0wMGMxLTcwNzItNzliOS1kODM3YTlhNDA3MWYiLCJpc3MiOiJodHRwczpcL1wvY29nbml0by1pZHAudXMtZWFzdC0yLmFtYXpvbmF3cy5jb21cL3VzLWVhc3QtMl9iR1RQTEZnTTciLCJjbGllbnRfaWQiOiJxaXFvZW1wOGhqbHQxYXIxcWYxMjMzbzd0Iiwib3JpZ2luX2p0aSI6IjhiY2ViMTViLWU3OTgtNDdmNy05ZmNmLThhYThlOGM5Y2E3ZSIsImV2ZW50X2lkIjoiNDY2NDRhZjItMWRkNC00YmY4LWE0NzQtZjE1ZDg3YjBlODQzIiwidG9rZW5fdXNlIjoiYWNjZXNzIiwic2NvcGUiOiJhd3MuY29nbml0by5zaWduaW4udXNlci5hZG1pbiIsImF1dGhfdGltZSI6MTcyNzE4ODg4MCwiZXhwIjoxNzI3MTkyNDgwLCJpYXQiOjE3MjcxODg4ODAsImp0aSI6IjI4ZmZlMzU2LTEzZjQtNGE0Yi04YWQwLWY3MDJlMDRlM2M5NiIsInVzZXJuYW1lIjoiOTE5YjM1OTAtMDBjMS03MDcyLTc5YjktZDgzN2E5YTQwNzFmIn0.RaYgkL_XPSi6miznkM9e7bflYyy_pktsQQreFjv9rq3yQazxeTNFIg_v_-7aDwHKS8t9D5eFjGq4H7eNKcX1MTc3Trgh6rpOCX14yE8FL4wJYxboA1uKvEoxK7JL_liTe0SzcsgDo5zAfcPCUu2AACn1Q45anX-_rxfIt5XfaGXyxFJ1mxVmeGSTeJNapwl7u0WUpBAc39z4YZS0MzggSBsfW4UoV_O0Xi6Z4ieK8rW720GOnWZdfyQnmG1oRcCH80hEtpQoWFKjARitV5sWbegeeYw60DDSjMSnygVHlpYI02mM-pqo4Jm4iTXTKlSdSO3RKLJDQ961VEWJnB02WQ" \
// 	--request "PUT" \
// 	--data '{"email": "thairon.ssilva@gmail.com","password": "Password!2"}'

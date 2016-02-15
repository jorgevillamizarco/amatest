package main

import "github.com/gin-gonic/gin"
import "github.com/jorgevillamizarco/amadeus21/auth"

func main() {
	r := gin.Default()
	r.Static("/wsdl", "./wsdl")
	r.GET("/auth", auth.Run)
	//r.GET("/fmptbs", fmptbs)
	r.Run(":8080")
}

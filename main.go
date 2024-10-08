package main

import (
    "github.com/gin-gonic/gin"
    "github.com/gin-contrib/cors"
)

func main() {
    r := gin.Default()

    // ตั้งค่า CORS
    config := cors.DefaultConfig()
    config.AllowOrigins = []string{"*"} // อนุญาตให้ทุกโดเมนสามารถเข้าถึง
    config.AllowMethods = []string{"GET", "POST", "PATCH", "PUT", "DELETE", "OPTIONS"}
    config.AllowHeaders = []string{"Origin", "Content-Type", "X-Auth-Token", "Authorization"}
    r.Use(cors.New(config))

    // เส้นทางสำหรับการยืนยันตัวตนและรีเฟรช Token
    r.POST("/login", loginHandler)
    r.POST("/refresh", refreshHandler)

    // เส้นทางที่ต้องการการยืนยันตัวตน
    auth:= r.Group("/auth")
    auth.Use(AuthMiddleware()) // ตรวจสอบสิทธิ์
    {
        auth.GET("/resource", authHandler)
    }

    // เริ่มเซิร์ฟเวอร์ที่พอร์ต 7089
    r.Run(":7089")
}

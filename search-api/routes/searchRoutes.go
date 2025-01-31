package routes

import (
	"search-api/controllers"

	"github.com/gin-gonic/gin"
)

// SetupRouter configura las rutas del servidor
func SetupRouter() *gin.Engine {
	r := gin.Default()

	// Rutas existentes
	r.POST("/hotels/index", controllers.IndexHotel)
	r.GET("/hotels/search", controllers.SearchHotels)

	// Ruta de prueba para verificar que el servidor esté corriendo
	r.GET("/test", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "Server is running!"})
	})

	return r
}

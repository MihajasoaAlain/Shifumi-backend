package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// HealthCheck godoc
// @Summary      Vérifier la santé de l'API
// @Description  Renvoie un statut 200 si le service Shifumi est opérationnel.
// @Tags         System
// @Produce      json
// @Success      200  {object}  map[string]string "Exemple: {"message": "work well"}"
// @Router       /health [get]
func HealthHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "work well",
	})
}

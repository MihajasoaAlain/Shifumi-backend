package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// HealthCheck godoc
// @Summary      Vérifier la santé de l'API
// @Description  Renvoie un statut 200 si le service Shifumi est opérationnel.
// @Tags         System
// @Produce      plain
// @Success      200  {string}  string  "work well"
// @Router       /health [get]
func HealthHandler(c *gin.Context) {
	c.String(http.StatusOK, "work well")
}
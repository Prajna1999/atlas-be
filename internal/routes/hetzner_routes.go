package routes

import (
	"net/http"
	"strconv"

	"github.com/Prajna1999/atlas-be/internal/service"

	"github.com/gin-gonic/gin"
)

type HetznerRoutes struct {
	service *service.HetznerService
}

// initialise NewHetznerRoutes
func NewHetznerRoutes(service *service.HetznerService) *HetznerRoutes {
	return &HetznerRoutes{service: service}
}

// Setup defined the routes related to Hetzner operations
func (h *HetznerRoutes) Setup(router *gin.RouterGroup) {
	hetzner := router.Group("/hetzner")
	{
		hetzner.GET("/servers", h.GetAllServers)
		hetzner.POST("/servers", h.CreateServer)
		hetzner.PUT("/servers/:id", h.UpdateServer)
		hetzner.DELETE("/servers/:id", h.DeleteServer)
		hetzner.GET("/servers/:id", h.GetOneServerByID)
		hetzner.GET("/servers/:id/metrics", h.GetServerMetricsByID)
	}

}

func (h *HetznerRoutes) GetAllServers(c *gin.Context) {
	servers, err := h.service.GetAllServers()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, servers)
}

func (h *HetznerRoutes) GetOneServerByID(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid server id"})
		return
	}
	server, err := h.service.GetServerByID(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, server)
}

func (h *HetznerRoutes) CreateServer(c *gin.Context) {
	var requestBody struct {
		Name       string `json:"name" binding:"required"`
		ServerType string `json:"server_type" binding:"required"`
		Image      string `json:"image" binding:"required"`
	}

	if err := c.ShouldBindBodyWithJSON(&requestBody); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	server, err := h.service.CreateServer(
		requestBody.Name, requestBody.Image, requestBody.ServerType)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	go h.service.LogOperation("create", server)
	c.JSON(http.StatusCreated, server)
}

func (h *HetznerRoutes) UpdateServer(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid server ID"})
		return
	}

	var requestBody struct {
		Name   string            `json:"name,omitempty"`
		Labels map[string]string `json:"labels,omitempty"`
	}
	if err := c.ShouldBindBodyWithJSON(&requestBody); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}
	server, err := h.service.UpdateServer(id, requestBody.Name, requestBody.Labels)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	go h.service.LogOperation("update", server)
	c.JSON(http.StatusOK, server)
}

func (h *HetznerRoutes) DeleteServer(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid server ID"})
		return
	}
	if err := h.service.DeleteServer(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	go h.service.LogOperation("delete", id)
	c.JSON(http.StatusOK, gin.H{"message": "Server deleted successfully"})

}

func (h *HetznerRoutes) GetServerMetricsByID(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid server ID"})
		return
	}

	start := c.Query("start")
	end := c.Query("end")
	if start == "" || end == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Start and end times are required"})
		return
	}

	metrics, err := h.service.GetServerMetricsByID(id, start, end)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	go h.service.LogOperation("get_metrics", id)
	c.JSON(http.StatusOK, metrics)
}

package handlers

import (
	"app/domain/models"
	"app/helpers"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func (r *appRoute) TicketStatusRoutes(rg *gin.RouterGroup) {
	api := rg.Group("/ticket-statuses")
	api.POST("", r.createTicketStatus)
	api.GET("", r.getTicketStatuses)
	api.GET("/:id", r.getTicketStatusByID)
	api.PUT("/:id", r.updateTicketStatus)
	api.DELETE("/:id", r.deleteTicketStatus)
}

// CreateTicketStatus godoc
// @Summary Create a new ticket status
// @Description Create a new ticket status
// @Tags ticket-statuses
// @Accept json
// @Produce json
// @Param status body models.TicketStatus true "Ticket Status Data"
// @Success 201 {object} helpers.Response{data=models.TicketStatus}
// @Router /ticket-statuses [post]
func (r *appRoute) createTicketStatus(c *gin.Context) {
	var status models.TicketStatus
	if err := c.ShouldBindJSON(&status); err != nil {
		response := helpers.NewResponse(http.StatusBadRequest, "Invalid request body", nil, nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	if err := r.Service.CreateTicketStatus(&status); err != nil {
		response := helpers.NewResponse(http.StatusInternalServerError, "Failed to create ticket status", nil, nil)
		c.JSON(http.StatusInternalServerError, response)
		return
	}

	response := helpers.NewResponse(http.StatusCreated, "Ticket status created successfully", nil, status)
	c.JSON(http.StatusCreated, response)
}

// GetTicketStatuses godoc
// @Summary Get all ticket statuses
// @Description Get a list of all ticket statuses
// @Tags ticket-statuses
// @Produce json
// @Success 200 {object} helpers.Response{data=[]models.TicketStatus}
// @Router /ticket-statuses [get]
func (r *appRoute) getTicketStatuses(c *gin.Context) {
	statuses, err := r.Service.GetTicketStatuses()
	if err != nil {
		response := helpers.NewResponse(http.StatusInternalServerError, "Failed to get ticket statuses", nil, nil)
		c.JSON(http.StatusInternalServerError, response)
		return
	}

	response := helpers.NewResponse(http.StatusOK, "Ticket statuses retrieved successfully", nil, statuses)
	c.JSON(http.StatusOK, response)
}

// GetTicketStatusByID godoc
// @Summary Get a ticket status by ID
// @Description Get a ticket status by its ID
// @Tags ticket-statuses
// @Produce json
// @Param id path int true "Status ID"
// @Success 200 {object} helpers.Response{data=models.TicketStatus}
// @Router /ticket-statuses/{id} [get]
func (r *appRoute) getTicketStatusByID(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		response := helpers.NewResponse(http.StatusBadRequest, "Invalid status ID", nil, nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	status, err := r.Service.GetTicketStatusByID(id)
	if err != nil {
		response := helpers.NewResponse(http.StatusNotFound, "Status not found", nil, nil)
		c.JSON(http.StatusNotFound, response)
		return
	}

	response := helpers.NewResponse(http.StatusOK, "Ticket status retrieved successfully", nil, status)
	c.JSON(http.StatusOK, response)
}

// UpdateTicketStatus godoc
// @Summary Update a ticket status
// @Description Update a ticket status by its ID
// @Tags ticket-statuses
// @Accept json
// @Produce json
// @Param id path int true "Status ID"
// @Param status body models.TicketStatus true "Updated Status Data"
// @Success 200 {object} helpers.Response{data=models.TicketStatus}
// @Router /ticket-statuses/{id} [put]
func (r *appRoute) updateTicketStatus(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		response := helpers.NewResponse(http.StatusBadRequest, "Invalid status ID", nil, nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	var status models.TicketStatus
	if err := c.ShouldBindJSON(&status); err != nil {
		response := helpers.NewResponse(http.StatusBadRequest, "Invalid request body", nil, nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	status.ID = id
	if err := r.Service.UpdateTicketStatus(&status); err != nil {
		response := helpers.NewResponse(http.StatusInternalServerError, "Failed to update ticket status", nil, nil)
		c.JSON(http.StatusInternalServerError, response)
		return
	}

	response := helpers.NewResponse(http.StatusOK, "Ticket status updated successfully", nil, status)
	c.JSON(http.StatusOK, response)
}

// DeleteTicketStatus godoc
// @Summary Delete a ticket status
// @Description Delete a ticket status by its ID
// @Tags ticket-statuses
// @Produce json
// @Param id path int true "Status ID"
// @Success 200 {object} helpers.Response
// @Router /ticket-statuses/{id} [delete]
func (r *appRoute) deleteTicketStatus(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		response := helpers.NewResponse(http.StatusBadRequest, "Invalid status ID", nil, nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	if err := r.Service.DeleteTicketStatus(id); err != nil {
		response := helpers.NewResponse(http.StatusInternalServerError, "Failed to delete ticket status", nil, nil)
		c.JSON(http.StatusInternalServerError, response)
		return
	}

	response := helpers.NewResponse(http.StatusOK, "Ticket status deleted successfully", nil, nil)
	c.JSON(http.StatusOK, response)
}

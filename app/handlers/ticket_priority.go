package handlers

import (
	"app/domain/models"
	"app/helpers"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func (r *appRoute) TicketPriorityRoutes(rg *gin.RouterGroup) {
	api := rg.Group("/ticket-priorities")
	api.POST("", r.createTicketPriority)
	api.GET("", r.getTicketPriorities)
	api.GET("/:id", r.getTicketPriorityByID)
	api.PUT("/:id", r.updateTicketPriority)
	api.DELETE("/:id", r.deleteTicketPriority)
}

// CreateTicketPriority godoc
// @Summary Create a new ticket priority
// @Description Create a new ticket priority
// @Tags ticket-priorities
// @Accept json
// @Produce json
// @Param priority body models.TicketPriority true "Ticket Priority Data"
// @Success 201 {object} helpers.Response{data=models.TicketPriority}
// @Router /ticket-priorities [post]
func (r *appRoute) createTicketPriority(c *gin.Context) {
	var priority models.TicketPriority
	if err := c.ShouldBindJSON(&priority); err != nil {
		response := helpers.NewResponse(http.StatusBadRequest, "Invalid request body", nil, nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	if err := r.Service.CreateTicketPriority(&priority); err != nil {
		response := helpers.NewResponse(http.StatusInternalServerError, "Failed to create ticket priority", nil, nil)
		c.JSON(http.StatusInternalServerError, response)
		return
	}

	response := helpers.NewResponse(http.StatusCreated, "Ticket priority created successfully", nil, priority)
	c.JSON(http.StatusCreated, response)
}

// GetTicketPriorities godoc
// @Summary Get all ticket priorities
// @Description Get a list of all ticket priorities
// @Tags ticket-priorities
// @Produce json
// @Success 200 {object} helpers.Response{data=[]models.TicketPriority}
// @Router /ticket-priorities [get]
func (r *appRoute) getTicketPriorities(c *gin.Context) {
	priorities, err := r.Service.GetTicketPriorities()
	if err != nil {
		response := helpers.NewResponse(http.StatusInternalServerError, "Failed to get ticket priorities", nil, nil)
		c.JSON(http.StatusInternalServerError, response)
		return
	}

	response := helpers.NewResponse(http.StatusOK, "Ticket priorities retrieved successfully", nil, priorities)
	c.JSON(http.StatusOK, response)
}

// GetTicketPriorityByID godoc
// @Summary Get a ticket priority by ID
// @Description Get a ticket priority by its ID
// @Tags ticket-priorities
// @Produce json
// @Param id path int true "Priority ID"
// @Success 200 {object} helpers.Response{data=models.TicketPriority}
// @Router /ticket-priorities/{id} [get]
func (r *appRoute) getTicketPriorityByID(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		response := helpers.NewResponse(http.StatusBadRequest, "Invalid priority ID", nil, nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	priority, err := r.Service.GetTicketPriorityByID(id)
	if err != nil {
		response := helpers.NewResponse(http.StatusNotFound, "Priority not found", nil, nil)
		c.JSON(http.StatusNotFound, response)
		return
	}

	response := helpers.NewResponse(http.StatusOK, "Ticket priority retrieved successfully", nil, priority)
	c.JSON(http.StatusOK, response)
}

// UpdateTicketPriority godoc
// @Summary Update a ticket priority
// @Description Update a ticket priority by its ID
// @Tags ticket-priorities
// @Accept json
// @Produce json
// @Param id path int true "Priority ID"
// @Param priority body models.TicketPriority true "Updated Priority Data"
// @Success 200 {object} helpers.Response{data=models.TicketPriority}
// @Router /ticket-priorities/{id} [put]
func (r *appRoute) updateTicketPriority(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		response := helpers.NewResponse(http.StatusBadRequest, "Invalid priority ID", nil, nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	var priority models.TicketPriority
	if err := c.ShouldBindJSON(&priority); err != nil {
		response := helpers.NewResponse(http.StatusBadRequest, "Invalid request body", nil, nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	priority.IDPriority = id
	if err := r.Service.UpdateTicketPriority(&priority); err != nil {
		response := helpers.NewResponse(http.StatusInternalServerError, "Failed to update ticket priority", nil, nil)
		c.JSON(http.StatusInternalServerError, response)
		return
	}

	response := helpers.NewResponse(http.StatusOK, "Ticket priority updated successfully", nil, priority)
	c.JSON(http.StatusOK, response)
}

// DeleteTicketPriority godoc
// @Summary Delete a ticket priority
// @Description Delete a ticket priority by its ID
// @Tags ticket-priorities
// @Produce json
// @Param id path int true "Priority ID"
// @Success 200 {object} helpers.Response
// @Router /ticket-priorities/{id} [delete]
func (r *appRoute) deleteTicketPriority(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		response := helpers.NewResponse(http.StatusBadRequest, "Invalid priority ID", nil, nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	if err := r.Service.DeleteTicketPriority(id); err != nil {
		response := helpers.NewResponse(http.StatusInternalServerError, "Failed to delete ticket priority", nil, nil)
		c.JSON(http.StatusInternalServerError, response)
		return
	}

	response := helpers.NewResponse(http.StatusOK, "Ticket priority deleted successfully", nil, nil)
	c.JSON(http.StatusOK, response)
}

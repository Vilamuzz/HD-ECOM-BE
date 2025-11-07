package handlers

import (
	"app/domain/models"
	"app/helpers"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func (r *appRoute) TicketRoutes(rg *gin.RouterGroup) {
	api := rg.Group("/tickets")
	api.POST("", r.createTicket)
	api.GET("", r.getTickets)
	api.GET("/:id", r.getTicketByID)
	api.PUT("/:id", r.updateTicket)
	api.DELETE("/:id", r.deleteTicket)
}

// CreateTicket godoc
// @Summary Create a new ticket
// @Description Create a new ticket
// @Tags tickets
// @Accept json
// @Produce json
// @Param ticket body models.Ticket true "Ticket Data"
// @Success 201 {object} helpers.Response{data=models.Ticket}
// @Router /tickets [post]
func (r *appRoute) createTicket(c *gin.Context) {
	var ticket models.Ticket
	if err := c.ShouldBindJSON(&ticket); err != nil {
		response := helpers.NewResponse(http.StatusBadRequest, "Invalid request body", nil, nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	if err := r.Service.CreateTicket(&ticket); err != nil {
		response := helpers.NewResponse(http.StatusInternalServerError, "Failed to create ticket", nil, nil)
		c.JSON(http.StatusInternalServerError, response)
		return
	}

	response := helpers.NewResponse(http.StatusCreated, "Ticket created successfully", nil, ticket)
	c.JSON(http.StatusCreated, response)
}

// GetTickets godoc
// @Summary Get all tickets
// @Description Get a list of all tickets
// @Tags tickets
// @Produce json
// @Success 200 {object} helpers.Response{data=[]models.Ticket}
// @Router /tickets [get]
func (r *appRoute) getTickets(c *gin.Context) {
	tickets, err := r.Service.GetTickets()
	if err != nil {
		response := helpers.NewResponse(http.StatusInternalServerError, "Failed to get tickets", nil, nil)
		c.JSON(http.StatusInternalServerError, response)
		return
	}

	response := helpers.NewResponse(http.StatusOK, "Tickets retrieved successfully", nil, tickets)
	c.JSON(http.StatusOK, response)
}

// GetTicketByID godoc
// @Summary Get a ticket by ID
// @Description Get a ticket by its ID
// @Tags tickets
// @Produce json
// @Param id path int true "Ticket ID"
// @Success 200 {object} helpers.Response{data=models.Ticket}
// @Router /tickets/{id} [get]
func (r *appRoute) getTicketByID(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		response := helpers.NewResponse(http.StatusBadRequest, "Invalid ticket ID", nil, nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	ticket, err := r.Service.GetTicketByID(id)
	if err != nil {
		response := helpers.NewResponse(http.StatusNotFound, "Ticket not found", nil, nil)
		c.JSON(http.StatusNotFound, response)
		return
	}

	response := helpers.NewResponse(http.StatusOK, "Ticket retrieved successfully", nil, ticket)
	c.JSON(http.StatusOK, response)
}

// UpdateTicket godoc
// @Summary Update a ticket
// @Description Update a ticket by its ID
// @Tags tickets
// @Accept json
// @Produce json
// @Param id path int true "Ticket ID"
// @Param ticket body models.Ticket true "Updated Ticket Data"
// @Success 200 {object} helpers.Response{data=models.Ticket}
// @Router /tickets/{id} [put]
func (r *appRoute) updateTicket(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		response := helpers.NewResponse(http.StatusBadRequest, "Invalid ticket ID", nil, nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	var ticket models.Ticket
	if err := c.ShouldBindJSON(&ticket); err != nil {
		response := helpers.NewResponse(http.StatusBadRequest, "Invalid request body", nil, nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	ticket.IDTicket = id
	if err := r.Service.UpdateTicket(&ticket); err != nil {
		response := helpers.NewResponse(http.StatusInternalServerError, "Failed to update ticket", nil, nil)
		c.JSON(http.StatusInternalServerError, response)
		return
	}

	response := helpers.NewResponse(http.StatusOK, "Ticket updated successfully", nil, ticket)
	c.JSON(http.StatusOK, response)
}

// DeleteTicket godoc
// @Summary Delete a ticket
// @Description Delete a ticket by its ID
// @Tags tickets
// @Produce json
// @Param id path int true "Ticket ID"
// @Success 200 {object} helpers.Response
// @Router /tickets/{id} [delete]
func (r *appRoute) deleteTicket(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		response := helpers.NewResponse(http.StatusBadRequest, "Invalid ticket ID", nil, nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	if err := r.Service.DeleteTicket(id); err != nil {
		response := helpers.NewResponse(http.StatusInternalServerError, "Failed to delete ticket", nil, nil)
		c.JSON(http.StatusInternalServerError, response)
		return
	}

	response := helpers.NewResponse(http.StatusOK, "Ticket deleted successfully", nil, nil)
	c.JSON(http.StatusOK, response)
}

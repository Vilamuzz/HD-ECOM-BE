package handlers

import (
	"app/domain/models"
	"app/helpers"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func (r *appRoute) TicketCategoryRoutes(rg *gin.RouterGroup) {
	api := rg.Group("/ticket-categories")
	api.POST("", r.Middleware.Auth(), r.Middleware.RequireRole(models.RoleAdmin), r.createTicketCategory)
	api.GET("", r.getTicketCategories)
	api.GET("/:id", r.getTicketCategoryByID)
	api.PUT("/:id", r.Middleware.Auth(), r.Middleware.RequireRole(models.RoleAdmin),  r.updateTicketCategory)
	api.DELETE("/:id", r.Middleware.Auth(), r.Middleware.RequireRole(models.RoleAdmin),  r.deleteTicketCategory)
}

// CreateTicketCategory godoc
// @Summary Create a new ticket category
// @Description Create a new ticket category
// @Tags ticket-categories
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param category body models.TicketCategory true "Ticket Category Data"
// @Success 201 {object} helpers.Response{data=models.TicketCategory}
// @Router /ticket-categories [post]
func (r *appRoute) createTicketCategory(c *gin.Context) {
	var category models.TicketCategory
	if err := c.ShouldBindJSON(&category); err != nil {
		response := helpers.NewResponse(http.StatusBadRequest, "Invalid request body", nil, nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	if err := r.Service.CreateTicketCategory(&category); err != nil {
		response := helpers.NewResponse(http.StatusInternalServerError, "Failed to create ticket category", nil, nil)
		c.JSON(http.StatusInternalServerError, response)
		return
	}

	response := helpers.NewResponse(http.StatusCreated, "Ticket category created successfully", nil, category)
	c.JSON(http.StatusCreated, response)
}

// GetTicketCategories godoc
// @Summary Get all ticket categories
// @Description Get a list of all ticket categories
// @Tags ticket-categories
// @Produce json
// @Success 200 {object} helpers.Response{data=[]models.TicketCategory}
// @Router /ticket-categories [get]
func (r *appRoute) getTicketCategories(c *gin.Context) {
	categories, err := r.Service.GetTicketCategories()
	if err != nil {
		response := helpers.NewResponse(http.StatusInternalServerError, "Failed to get ticket categories", nil, nil)
		c.JSON(http.StatusInternalServerError, response)
		return
	}

	response := helpers.NewResponse(http.StatusOK, "Ticket categories retrieved successfully", nil, categories)
	c.JSON(http.StatusOK, response)
}

// GetTicketCategoryByID godoc
// @Summary Get a ticket category by ID
// @Description Get a ticket category by its ID
// @Tags ticket-categories
// @Produce json
// @Param id path int true "Category ID"
// @Success 200 {object} helpers.Response{data=models.TicketCategory}
// @Router /ticket-categories/{id} [get]
func (r *appRoute) getTicketCategoryByID(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		response := helpers.NewResponse(http.StatusBadRequest, "Invalid category ID", nil, nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	category, err := r.Service.GetTicketCategoryByID(id)
	if err != nil {
		response := helpers.NewResponse(http.StatusNotFound, "Category not found", nil, nil)
		c.JSON(http.StatusNotFound, response)
		return
	}

	response := helpers.NewResponse(http.StatusOK, "Ticket category retrieved successfully", nil, category)
	c.JSON(http.StatusOK, response)
}

// UpdateTicketCategory godoc
// @Summary Update a ticket category
// @Description Update a ticket category by its ID
// @Tags ticket-categories
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Category ID"
// @Param category body models.TicketCategory true "Updated Category Data"
// @Success 200 {object} helpers.Response{data=models.TicketCategory}
// @Router /ticket-categories/{id} [put]
func (r *appRoute) updateTicketCategory(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		response := helpers.NewResponse(http.StatusBadRequest, "Invalid category ID", nil, nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	var category models.TicketCategory
	if err := c.ShouldBindJSON(&category); err != nil {
		response := helpers.NewResponse(http.StatusBadRequest, "Invalid request body", nil, nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	category.ID = id
	if err := r.Service.UpdateTicketCategory(&category); err != nil {
		response := helpers.NewResponse(http.StatusInternalServerError, "Failed to update ticket category", nil, nil)
		c.JSON(http.StatusInternalServerError, response)
		return
	}

	response := helpers.NewResponse(http.StatusOK, "Ticket category updated successfully", nil, category)
	c.JSON(http.StatusOK, response)
}

// DeleteTicketCategory godoc
// @Summary Delete a ticket category
// @Description Delete a ticket category by its ID
// @Tags ticket-categories
// @Produce json
// @Security BearerAuth
// @Param id path int true "Category ID"
// @Success 200 {object} helpers.Response
// @Router /ticket-categories/{id} [delete]
func (r *appRoute) deleteTicketCategory(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		response := helpers.NewResponse(http.StatusBadRequest, "Invalid category ID", nil, nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	if err := r.Service.DeleteTicketCategory(id); err != nil {
		response := helpers.NewResponse(http.StatusInternalServerError, "Failed to delete ticket category", nil, nil)
		c.JSON(http.StatusInternalServerError, response)
		return
	}

	response := helpers.NewResponse(http.StatusOK, "Ticket category deleted successfully", nil, nil)
	c.JSON(http.StatusOK, response)
}

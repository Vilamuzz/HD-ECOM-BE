package handlers

import (
	"app/domain/models"
	"app/helpers"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func (r *appRoute) TicketLogRoutes(rg *gin.RouterGroup) {
	api := rg.Group("/ticket-logs")
	api.POST("", r.createTicketLog)
	api.GET("", r.getTicketLogs)
	api.GET("/:id", r.getTicketLogByID)
}

// CreateTicketLog godoc
// @Summary Create a new ticket log
// @Description Create a new activity log for a ticket
// @Tags ticket-logs
// @Accept json
// @Produce json
// @Param log body models.TicketLog true "Log Data"
// @Success 201 {object} helpers.Response{data=models.TicketLog}
// @Failure 400 {object} helpers.Response
// @Failure 500 {object} helpers.Response
// @Router /ticket-logs [post]
func (r *appRoute) createTicketLog(c *gin.Context) {
	var log models.TicketLog
	if err := c.ShouldBindJSON(&log); err != nil {
		c.JSON(http.StatusBadRequest, helpers.NewResponse(http.StatusBadRequest, "Invalid request body", nil, nil))
		return
	}

	if err := r.Service.CreateTicketLog(&log); err != nil {
		c.JSON(http.StatusInternalServerError, helpers.NewResponse(http.StatusInternalServerError, "Failed to create log", nil, nil))
		return
	}

	c.JSON(http.StatusCreated, helpers.NewResponse(http.StatusCreated, "Log created successfully", nil, log))
}

// GetTicketLogs godoc
// @Summary Get all ticket logs
// @Description Get all activity logs, optionally filtered by ticket ID
// @Tags ticket-logs
// @Produce json
// @Param ticket_id query int false "Filter by Ticket ID"
// @Success 200 {object} helpers.Response{data=[]models.TicketLog}
// @Failure 500 {object} helpers.Response
// @Router /ticket-logs [get]
func (r *appRoute) getTicketLogs(c *gin.Context) {
	ticketID, _ := strconv.Atoi(c.Query("ticket_id"))

	var logs []models.TicketLog
	var err error

	if ticketID > 0 {
		logs, err = r.Service.GetTicketLogsByTicketID(ticketID)
	} else {
		logs, err = r.Service.GetTicketLogs()
	}

	if err != nil {
		c.JSON(http.StatusInternalServerError, helpers.NewResponse(http.StatusInternalServerError, "Failed to get logs", nil, nil))
		return
	}

	c.JSON(http.StatusOK, helpers.NewResponse(http.StatusOK, "Logs retrieved successfully", nil, logs))
}

// GetTicketLogByID godoc
// @Summary Get a ticket log by ID
// @Description Get a specific log entry by its ID
// @Tags ticket-logs
// @Produce json
// @Param id path int true "Log ID"
// @Success 200 {object} helpers.Response{data=models.TicketLog}
// @Failure 404 {object} helpers.Response
// @Router /ticket-logs/{id} [get]
func (r *appRoute) getTicketLogByID(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, helpers.NewResponse(http.StatusBadRequest, "Invalid log ID", nil, nil))
		return
	}

	log, err := r.Service.GetTicketLogByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, helpers.NewResponse(http.StatusNotFound, "Log not found", nil, nil))
		return
	}

	c.JSON(http.StatusOK, helpers.NewResponse(http.StatusOK, "Log retrieved successfully", nil, log))
}
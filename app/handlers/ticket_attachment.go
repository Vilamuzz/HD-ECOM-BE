package handlers

import (
	"app/domain/models"
	"app/helpers"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func (r *appRoute) TicketAttachmentRoutes(rg *gin.RouterGroup) {
	api := rg.Group("/ticket-attachments")
	api.POST("", r.createTicketAttachment)
	api.GET("", r.getTicketAttachments)
	api.GET("/:id", r.getTicketAttachmentByID)
	api.PUT("/:id", r.updateTicketAttachment)
	api.DELETE("/:id", r.deleteTicketAttachment)
}

// CreateTicketAttachment godoc
// @Summary Create a new ticket attachment
// @Description Create a new ticket attachment
// @Tags ticket-attachments
// @Accept multipart/form-data
// @Produce json
// @Param id_ticket formData int true "Ticket ID"
// @Param file formData file true "Attachment file"
// @Success 201 {object} helpers.Response{data=models.TicketAttachment}
// @Failure 400 {object} helpers.Response
// @Failure 500 {object} helpers.Response
// @Router /ticket-attachments [post]
func (r *appRoute) createTicketAttachment(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		response := helpers.NewResponse(http.StatusBadRequest, "No file uploaded", nil, nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	ticketID, err := strconv.Atoi(c.PostForm("id_ticket"))
	if err != nil {
		response := helpers.NewResponse(http.StatusBadRequest, "Invalid ticket ID", nil, nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	// Save file and create attachment record
	filePath := "uploads/" + file.Filename // You might want to generate a unique filename
	if err := c.SaveUploadedFile(file, filePath); err != nil {
		response := helpers.NewResponse(http.StatusInternalServerError, "Failed to save file", nil, nil)
		c.JSON(http.StatusInternalServerError, response)
		return
	}

	attachment := &models.TicketAttachment{
		IDTicket: ticketID,
		FilePath: filePath,
	}

	if err := r.Service.CreateTicketAttachment(attachment); err != nil {
		response := helpers.NewResponse(http.StatusInternalServerError, "Failed to create attachment record", nil, nil)
		c.JSON(http.StatusInternalServerError, response)
		return
	}

	response := helpers.NewResponse(http.StatusCreated, "Attachment created successfully", nil, attachment)
	c.JSON(http.StatusCreated, response)
}

// GetTicketAttachments godoc
// @Summary Get all ticket attachments
// @Description Get a list of all ticket attachments
// @Tags ticket-attachments
// @Produce json
// @Param ticket_id query int false "Filter by ticket ID"
// @Success 200 {object} helpers.Response{data=[]models.TicketAttachment}
// @Failure 500 {object} helpers.Response
// @Router /ticket-attachments [get]
func (r *appRoute) getTicketAttachments(c *gin.Context) {
	ticketID, _ := strconv.Atoi(c.Query("ticket_id"))

	var attachments []models.TicketAttachment
	var err error

	if ticketID > 0 {
		attachments, err = r.Service.GetTicketAttachmentsByTicketID(ticketID)
	} else {
		attachments, err = r.Service.GetTicketAttachments()
	}

	if err != nil {
		response := helpers.NewResponse(http.StatusInternalServerError, "Failed to get attachments", nil, nil)
		c.JSON(http.StatusInternalServerError, response)
		return
	}

	response := helpers.NewResponse(http.StatusOK, "Attachments retrieved successfully", nil, attachments)
	c.JSON(http.StatusOK, response)
}

// GetTicketAttachmentByID godoc
// @Summary Get a ticket attachment by ID
// @Description Get a ticket attachment by its ID
// @Tags ticket-attachments
// @Produce json
// @Param id path int true "Attachment ID"
// @Success 200 {object} helpers.Response{data=models.TicketAttachment}
// @Failure 404 {object} helpers.Response
// @Router /ticket-attachments/{id} [get]
func (r *appRoute) getTicketAttachmentByID(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		response := helpers.NewResponse(http.StatusBadRequest, "Invalid attachment ID", nil, nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	attachment, err := r.Service.GetTicketAttachmentByID(id)
	if err != nil {
		response := helpers.NewResponse(http.StatusNotFound, "Attachment not found", nil, nil)
		c.JSON(http.StatusNotFound, response)
		return
	}

	response := helpers.NewResponse(http.StatusOK, "Attachment retrieved successfully", nil, attachment)
	c.JSON(http.StatusOK, response)
}

// UpdateTicketAttachment godoc
// @Summary Update a ticket attachment
// @Description Update a ticket attachment by its ID
// @Tags ticket-attachments
// @Accept multipart/form-data
// @Produce json
// @Param id path int true "Attachment ID"
// @Param id_ticket formData int false "Ticket ID"
// @Param file formData file false "New attachment file"
// @Success 200 {object} helpers.Response{data=models.TicketAttachment}
// @Failure 400 {object} helpers.Response
// @Failure 404 {object} helpers.Response
// @Failure 500 {object} helpers.Response
// @Router /ticket-attachments/{id} [put]
func (r *appRoute) updateTicketAttachment(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		response := helpers.NewResponse(http.StatusBadRequest, "Invalid attachment ID", nil, nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	attachment, err := r.Service.GetTicketAttachmentByID(id)
	if err != nil {
		response := helpers.NewResponse(http.StatusNotFound, "Attachment not found", nil, nil)
		c.JSON(http.StatusNotFound, response)
		return
	}

	// Update ticket ID if provided
	if ticketID := c.PostForm("id_ticket"); ticketID != "" {
		if id, err := strconv.Atoi(ticketID); err == nil {
			attachment.IDTicket = id
		}
	}

	// Update file if provided
	if file, err := c.FormFile("file"); err == nil {
		filePath := "uploads/" + file.Filename // You might want to generate a unique filename
		if err := c.SaveUploadedFile(file, filePath); err != nil {
			response := helpers.NewResponse(http.StatusInternalServerError, "Failed to save file", nil, nil)
			c.JSON(http.StatusInternalServerError, response)
			return
		}
		attachment.FilePath = filePath
	}

	if err := r.Service.UpdateTicketAttachment(attachment); err != nil {
		response := helpers.NewResponse(http.StatusInternalServerError, "Failed to update attachment", nil, nil)
		c.JSON(http.StatusInternalServerError, response)
		return
	}

	response := helpers.NewResponse(http.StatusOK, "Attachment updated successfully", nil, attachment)
	c.JSON(http.StatusOK, response)
}

// DeleteTicketAttachment godoc
// @Summary Delete a ticket attachment
// @Description Delete a ticket attachment by its ID
// @Tags ticket-attachments
// @Produce json
// @Param id path int true "Attachment ID"
// @Success 200 {object} helpers.Response
// @Failure 400 {object} helpers.Response
// @Failure 500 {object} helpers.Response
// @Router /ticket-attachments/{id} [delete]
func (r *appRoute) deleteTicketAttachment(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		response := helpers.NewResponse(http.StatusBadRequest, "Invalid attachment ID", nil, nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	if err := r.Service.DeleteTicketAttachment(id); err != nil {
		response := helpers.NewResponse(http.StatusInternalServerError, "Failed to delete attachment", nil, nil)
		c.JSON(http.StatusInternalServerError, response)
		return
	}

	response := helpers.NewResponse(http.StatusOK, "Attachment deleted successfully", nil, nil)
	c.JSON(http.StatusOK, response)
}

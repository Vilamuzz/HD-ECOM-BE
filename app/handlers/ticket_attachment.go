package handlers

import (
	"app/domain/models"
	"app/helpers"
	"log"
	"mime/multipart"
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
		log.Printf("[ticket-attachment] no file uploaded: %v", err)
		response := helpers.NewResponse(http.StatusBadRequest, "No file uploaded", nil, nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	// Log incoming file info for debugging
	log.Printf("[ticket-attachment] incoming file: name=%s size=%d header=%v", file.Filename, file.Size, file.Header)

	ticketID, err := strconv.Atoi(c.PostForm("id_ticket"))
	if err != nil {
		log.Printf("[ticket-attachment] invalid ticket id: %v", err)
		response := helpers.NewResponse(http.StatusBadRequest, "Invalid ticket ID", nil, nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	attachment, err := r.Service.CreateTicketAttachment(ticketID, file)
	if err != nil {
		response := helpers.NewResponse(http.StatusInternalServerError, "Failed to create attachment", nil, nil)
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

	attachment, downloadURL, err := r.Service.GetTicketAttachmentByID(id)
	if err != nil {
		response := helpers.NewResponse(http.StatusNotFound, "Attachment not found", nil, nil)
		c.JSON(http.StatusNotFound, response)
		return
	}

	if downloadURL == "" {
		response := helpers.NewResponse(http.StatusInternalServerError, "Failed to generate download URL", nil, nil)
		c.JSON(http.StatusInternalServerError, response)
		return
	}

	data := map[string]interface{}{
		"attachment":   attachment,
		"download_url": downloadURL,
	}

	response := helpers.NewResponse(http.StatusOK, "Attachment retrieved successfully", nil, data)
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

	var ticketID *int
	if ticketIDStr := c.PostForm("id_ticket"); ticketIDStr != "" {
		if parsedID, err := strconv.Atoi(ticketIDStr); err == nil {
			ticketID = &parsedID
		}
	}

	var file *multipart.FileHeader
	if f, err := c.FormFile("file"); err == nil {
		file = f
	}

	attachment, err := r.Service.UpdateTicketAttachment(id, ticketID, file)
	if err != nil {
		if err.Error() == "record not found" {
			response := helpers.NewResponse(http.StatusNotFound, "Attachment not found", nil, nil)
			c.JSON(http.StatusNotFound, response)
		} else {
			response := helpers.NewResponse(http.StatusInternalServerError, "Failed to update attachment", nil, nil)
			c.JSON(http.StatusInternalServerError, response)
		}
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

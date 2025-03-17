package handler

import (
	"net/http"
	"strconv"
	"strings"

	"mailbox-api/api/middleware"
	"mailbox-api/logger"
	"mailbox-api/model"
	"mailbox-api/service"

	"github.com/gin-gonic/gin"
)

type MailboxHandler struct {
	service service.MailboxService
	logger  *logger.Logger
}

func NewMailboxHandler(service service.MailboxService, logger *logger.Logger) *MailboxHandler {
	return &MailboxHandler{
		service: service,
		logger:  logger,
	}
}

func (h *MailboxHandler) GetMailboxes(c *gin.Context) {
	filter, err := parseMailboxFilter(c)
	if err != nil {
		h.logger.Error("Failed to parse mailbox filter", "error", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var response interface{}
	role, _ := c.Get("role")
	userRole, ok := role.(middleware.Role)

	if !ok {
		h.logger.Error("Failed to get user role")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	if userRole == middleware.RoleCEO {
		response, err = h.service.GetMailboxes(c.Request.Context(), filter)
	} else if userRole == middleware.RoleCTO {
		response, err = h.service.GetMailboxesInSubOrg(c.Request.Context(), string(userRole), filter)
	} else {
		c.JSON(http.StatusForbidden, gin.H{"error": "Access denied"})
		return
	}

	if err != nil {
		h.logger.Error("Failed to get mailboxes", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get mailboxes"})
		return
	}

	c.JSON(http.StatusOK, response)
}

func (h *MailboxHandler) GetMailbox(c *gin.Context) {
	identifier := c.Param("id")

	mailbox, err := h.service.GetMailboxByIdentifier(c.Request.Context(), identifier)
	if err != nil {
		h.logger.Error("Failed to get mailbox", "error", err, "identifier", identifier)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get mailbox"})
		return
	}

	if mailbox == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Mailbox not found"})
		return
	}

	role, _ := c.Get("role")
	userRole, ok := role.(middleware.Role)
	if !ok {
		h.logger.Error("Failed to get user role")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	if userRole == middleware.RoleCEO {
		c.JSON(http.StatusOK, mailbox)
		return
	}

	if userRole == middleware.RoleCTO {
		ctoIdentifier := "david.brown@falafel.org"

		isUnderCTO, err := h.service.IsMailboxInSubOrg(c.Request.Context(), ctoIdentifier, identifier)
		if err != nil {
			h.logger.Error("Failed to check if mailbox is in CTO's sub-org", "error", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
			return
		}

		if !isUnderCTO && identifier != ctoIdentifier {
			c.JSON(http.StatusForbidden, gin.H{"error": "Access denied"})
			return
		}

		c.JSON(http.StatusOK, mailbox)
		return
	}

	c.JSON(http.StatusForbidden, gin.H{"error": "Access denied"})
}

func (h *MailboxHandler) CalculateOrgMetrics(c *gin.Context) {
	err := h.service.CalculateOrgMetrics(c.Request.Context())
	if err != nil {
		h.logger.Error("Failed to calculate org metrics", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to calculate org metrics"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Org metrics calculated successfully"})
}

func parseMailboxFilter(c *gin.Context) (model.MailboxFilter, error) {
	var filter model.MailboxFilter

	filter.SearchTerm = c.Query("search")

	if departmentStr := c.Query("department"); departmentStr != "" {
		department, err := strconv.Atoi(departmentStr)
		if err != nil {
			return filter, err
		}
		filter.Department = department
	}

	if orgDepthStr := c.Query("org_depth_exact"); orgDepthStr != "" {
		orgDepth, err := strconv.Atoi(orgDepthStr)
		if err != nil {
			return filter, err
		}
		filter.OrgDepthExact = &orgDepth
	}

	if orgDepthGtStr := c.Query("org_depth_gt"); orgDepthGtStr != "" {
		orgDepthGt, err := strconv.Atoi(orgDepthGtStr)
		if err != nil {
			return filter, err
		}
		filter.OrgDepthGt = &orgDepthGt
	}

	if orgDepthLtStr := c.Query("org_depth_lt"); orgDepthLtStr != "" {
		orgDepthLt, err := strconv.Atoi(orgDepthLtStr)
		if err != nil {
			return filter, err
		}
		filter.OrgDepthLt = &orgDepthLt
	}

	if subOrgSizeMinStr := c.Query("sub_org_size_min"); subOrgSizeMinStr != "" {
		subOrgSizeMin, err := strconv.Atoi(subOrgSizeMinStr)
		if err != nil {
			return filter, err
		}
		filter.SubOrgSizeMin = &subOrgSizeMin
	}

	if subOrgSizeMaxStr := c.Query("sub_org_size_max"); subOrgSizeMaxStr != "" {
		subOrgSizeMax, err := strconv.Atoi(subOrgSizeMaxStr)
		if err != nil {
			return filter, err
		}
		filter.SubOrgSizeMax = &subOrgSizeMax
	}

	if sortBy := c.QueryArray("sort_by"); len(sortBy) > 0 {
		filter.SortBy = sortBy

		sortDirs := c.QueryArray("sort_dir")
		if len(sortDirs) < len(sortBy) {
			for i := len(sortDirs); i < len(sortBy); i++ {
				sortDirs = append(sortDirs, "asc")
			}
		}
		filter.SortDirections = sortDirs
	}

	if fields := c.Query("fields"); fields != "" {
		filter.Fields = strings.Split(fields, ",")
	}

	if pageStr := c.Query("page"); pageStr != "" {
		page, err := strconv.Atoi(pageStr)
		if err != nil {
			return filter, err
		}
		filter.Page = page
	} else {
		filter.Page = 1
	}

	if pageSizeStr := c.Query("page_size"); pageSizeStr != "" {
		pageSize, err := strconv.Atoi(pageSizeStr)
		if err != nil {
			return filter, err
		}
		filter.PageSize = pageSize
	} else {
		filter.PageSize = 10
	}

	return filter, nil
}

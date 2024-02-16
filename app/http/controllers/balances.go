package controllers

import (
	"bank_file_analyser/domain"
	"mime/multipart"
	"net/http"

	"github.com/gin-gonic/gin"
)

type HTTPHandler struct {
	AccBalanceService domain.BalanceGeneratorService
}

func NewHTTPHandler(srvc domain.BalanceGeneratorService) *HTTPHandler {
	return &HTTPHandler{AccBalanceService: srvc}
}

type ProcessFileRequest struct {
	File *multipart.FileHeader `form:"file" binding:"required"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}

type SuccessResponse struct {
	Message string                           `json:"message"`
	Result  *domain.FormattedBankAccBalances `json:"result"`
}

// ProcessStatement godoc
// @Summary Use to upload statement file for processing
// @Description Use to upload master files to generate account balances
// @Tags Accounts
// @Produce json
// @Param file formData file true "file to process"
// @Success 200 {object} SuccessResponse
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /process_statement [post]
func (h *HTTPHandler) ProcessStatement(c *gin.Context) {
	req := ProcessFileRequest{}
	err := c.Bind(&req)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
		return
	}

	if req.File == nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "invalid file"})
		return
	}

	file, err := req.File.Open()
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		return
	}

	accBalances, err := h.AccBalanceService.GenerateAccBalancesFromFile(file)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, SuccessResponse{Message: "success", Result: h.AccBalanceService.FormatAccountBalances(accBalances)})
}

package controllers

import (
	"bank_file_analyser/config"
	"bank_file_analyser/domain"
	"bank_file_analyser/fileparser"
	"mime/multipart"
	"net/http"

	"github.com/gin-gonic/gin"
)

// Http handler for balances API
type AccountsHandler struct {
	AccBalanceService domain.BalanceGeneratorService
	Config            *config.AppConfig
}

// initializes new Http handler for balances API
func NewAccountsHandler(srvc domain.BalanceGeneratorService, conf *config.AppConfig) *AccountsHandler {
	return &AccountsHandler{AccBalanceService: srvc, Config: conf}
}

// ProcessFile API request struct
type ProcessFileRequest struct {
	File *multipart.FileHeader `form:"file" binding:"required"`
	Date string                `form:"date"`
}

// Error response struct
type ErrorResponse struct {
	Error string `json:"error"`
}

// ProcessFile API success response struct
type ProcessFileSuccessResponse struct {
	Message string                           `json:"message"`
	Result  *domain.FormattedBankAccBalances `json:"result"`
}

// ProcessStatement godoc
// @Summary Use to upload statement file for processing
// @Description Use to upload master files to generate account balances
// @Tags Accounts
// @Produce json
// @Param file formData file true "file to process"
// @Param date formData string false "filter date format:DD/MM/YYYY"
// @Success 200 {object} ProcessFileSuccessResponse
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /process_statement [post]
func (h *AccountsHandler) ProcessStatement(c *gin.Context) {
	req := ProcessFileRequest{}

	//parse request
	err := c.Bind(&req)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
		return
	}

	//check if file is not empty
	if req.File == nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "invalid file"})
		return
	}

	//initialize file parser
	parser := fileparser.NewCSVParser(rune(h.Config.FileColumnSeparator[0]), h.Config.DecimalPrecision, h.Config.FileDateFormat)

	//open uploaded file
	file, err := req.File.Open()
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		return
	}
	defer file.Close()

	//process file data
	accBalances, err := h.AccBalanceService.GenerateAccBalancesFromFile(parser, file, req.Date)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		return
	}

	//return processed result
	c.JSON(http.StatusOK, ProcessFileSuccessResponse{Message: "success", Result: h.AccBalanceService.FormatAccountBalances(accBalances)})
}

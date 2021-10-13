package md

import (
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/go-playground/validator"
	"github.com/gofiber/fiber/v2"
)

type MDHandlers struct {
	mdService *MDService
}

// InitMDHandlers Creates an instance of a MDHandlers
// Requires a reference to a md.Service instance
func InitMDHandlers(mdService *MDService) *MDHandlers {
	return &MDHandlers{mdService: mdService}
}

// ConfiugureRoutes
// Imports and configures various routes for
// all modules.
func (m *MDHandlers) ConfigureRoutes(app *fiber.App) {
	app.Post("/md", m.CreateMDHandler)
	app.Patch("/md", m.UpdateMDHandler)
	app.Get("/md/search", m.SearchMDHandler)
	app.Get("/md/:id", m.GetMDHandler)
	app.Get("/md", m.GetAllMDHandler)
	app.Delete("/md/:id", m.DeleteMDHandler)
}

// CreateMDHandler POST - creates a MarkdownSnippet from the provided body
// @Summary Create new a markdown snippet
// @Accept json
// @Produce json
// @Tags md
// @Success 201 {object} MarkdownSnippet
// @Failure 400 {object} api.ErrorResponse
// @Failure 400 {object} api.ErrorResponse
// @Failure 500 {object} api.ErrorResponse
// @Router /md [post]
// @Param message body CreateMDReq true "Post Body"
func (m *MDHandlers) CreateMDHandler(ctx *fiber.Ctx) error {
	snippetRequest := new(CreateMDReq)
	if err := ctx.BodyParser(snippetRequest); err != nil {
		log.Printf("Failed to parse CreateMarkdownSnippet: %s", err)
		return fiber.NewError(http.StatusBadRequest, err.Error())
	}

	if errs := ValidateStruct(snippetRequest); errs != nil {
		ctx.Status(http.StatusBadRequest)
		return ctx.JSON(errs)
	}

	newSnippet, err := m.mdService.CreateMarkdownSnippet(snippetRequest)
	if err != nil {
		log.Printf("Failed in insert new MarkdownSnippet: %s", err)
		return fiber.NewError(http.StatusInternalServerError, err.Error())
	}

	ctx.Status(http.StatusCreated)
	return ctx.JSON(newSnippet)
}

// GetMDHandler GET - MarkdownSnippet Retrieval
// @Summary Retrieve Markdown Snippet
// @Accept json
// @Produce json
// @Tags md
// @Success 200 {object} MarkdownSnippet
// @Failure 400 {object} api.ErrorResponse
// @Failure 500 {object} api.ErrorResponse
// @Router /md/{id} [get]
// @Param id path string true "Snippet uuid"
func (m *MDHandlers) GetMDHandler(ctx *fiber.Ctx) error {
	id := ctx.Params("id")

	snippet, err := m.mdService.GetMarkdownSnippet(id)
	if err != nil {
		log.Printf("Error Retrieving Markdown Snippet %s: %s", id, err)
		ctx.Status(http.StatusInternalServerError)
		return fiber.NewError(http.StatusInternalServerError, err.Error())
	}
	if snippet == nil && err == nil {
		ctx.Status(http.StatusNotFound)
		return fiber.NewError(http.StatusNotFound, "Markdown Snippet Not Found")
	}

	return ctx.JSON(snippet)
}

// GetAllMDHandler GET - Get All MarkdownSnippets Retrieval
// @Deprecated
// @Summary Retrieve All Markdown Snippets
// @Accept json
// @Produce json
// @Tags md
// @Success 200 {object} []MDListItem
// @Failure 400 {object} api.ErrorResponse
// @Failure 500 {object} api.ErrorResponse
// @Router /md [get]
func (m *MDHandlers) GetAllMDHandler(ctx *fiber.Ctx) error {
	snippets, err := m.mdService.GetAllMarkdownSnippets()
	if err != nil {
		log.Printf("Error Retrieving All Markdown Snippet: %s", err)
		ctx.Status(http.StatusInternalServerError)
		return fiber.NewError(http.StatusBadRequest, err.Error())
	}

	return ctx.JSON(snippets)
}

// SearchMDHandler GET - Search for MarkdownSnippets
// @Summary Search, Sort, and paginate through MarkdownSnippets.
// @Param text query string false "Search Term"
// @Param limit query int false "Number of Snippets" default(10)
// @Param skip query int false "Skip Number of Snippets" default(0)
// @Param sort query string false "Sort By" Enums(createDate_ASC, createDate_DESC)
// @Accept json
// @Produce json
// @Tags md
// @Success 200 {object} []MDListItem
// @Failure 400 {object} api.ErrorResponse
// @Failure 500 {object} api.ErrorResponse
// @Router /md/search [get]
func (m *MDHandlers) SearchMDHandler(ctx *fiber.Ctx) error {
	limit, err := strconv.ParseInt(ctx.Query("limit", "10"), 10, 64)
	if err != nil {
		return fiber.NewError(http.StatusBadRequest, "limit: invalid value")
	}
	skip, err := strconv.ParseInt(ctx.Query("skip", "0"), 10, 64)
	if err != nil {
		return fiber.NewError(http.StatusBadRequest, "skip: invalid value")
	}

	sort := SortBy(ctx.Query("sort", CreateDate_DESC))
	if err := sort.validate(); err != nil {
		return fiber.NewError(http.StatusBadRequest, fmt.Sprintf("`%s` is not a valid value for sort", sort))
	}

	params := MDSearchParams{
		Text:   ctx.Query("text"),
		Limit:  limit,
		Skip:   skip,
		SortBy: sort,
	}

	snippets, err := m.mdService.SearchMarkdownSnippets(params)
	if err != nil {
		log.Printf("Error Searching for Markdown Snippets: %s", err)
		ctx.Status(http.StatusInternalServerError)
		return fiber.NewError(http.StatusBadRequest, err.Error())
	}

	return ctx.JSON(snippets)
}

// UpdateMDHandler PATCH - Updates a MarkdownSnippet
// @Summary Updates a markdown snippet
// @Accept json
// @Produce json
// @Tags md
// @Success 200 {object} MarkdownSnippet
// @Failure 400 {object} api.ErrorResponse
// @Failure 404 {object} api.ErrorResponse
// @Failure 500 {object} api.ErrorResponse
// @Router /md [patch]
// @Param message body UpdateMDReq true "Patch Body"
func (m *MDHandlers) UpdateMDHandler(ctx *fiber.Ctx) error {
	patchSnippet := new(UpdateMDReq)
	if err := ctx.BodyParser(patchSnippet); err != nil {
		ctx.Status(http.StatusBadRequest)
		return fiber.NewError(http.StatusBadRequest, err.Error())
	}

	if errs := ValidateStruct(patchSnippet); errs != nil {
		ctx.Status(http.StatusBadRequest)
		return ctx.JSON(errs)
	}

	if !m.mdService.ValidateIdAndKey(patchSnippet.ID, patchSnippet.UpdateKey) {
		return fiber.NewError(http.StatusUnauthorized, "Invalid Update Key")
	}

	updatedSnippet, err := m.mdService.UpdateMarkdownSnippet(patchSnippet)
	if err != nil {
		log.Printf("Failed in update MarkdownSnippet %s: %s", patchSnippet.ID, err)
		ctx.Status(http.StatusBadRequest)
		return fiber.NewError(http.StatusInternalServerError, err.Error())
	}

	return ctx.JSON(updatedSnippet)
}

// DeleteMDHandler DELETE - Removes MarkdownSnippet permanantly
// @Summary Removes MarkdownSnippet permanantly
// @Accept json
// @Produce json
// @Tags md
// @Success 204
// @Failure 400 {object} api.ErrorResponse
// @Failure 500 {object} api.ErrorResponse
// @Router /md/{id} [delete]
// @Param id path string true "Snippet uuid"
func (m *MDHandlers) DeleteMDHandler(ctx *fiber.Ctx) error {
	id := ctx.Params("id")
	log.Printf("Deleting %s", id)
	// TODO: Make service call to remove snippet.

	ctx.Status(fiber.StatusNoContent)
	return nil
}

type ErrorResponse struct {
	FailedField string
	Tag         string
	Value       string
}

func ValidateStruct(v interface{}) []*ErrorResponse {
	var errors []*ErrorResponse
	validate := validator.New()
	err := validate.Struct(v)
	if err != nil {
		for _, err := range err.(validator.ValidationErrors) {
			var element ErrorResponse
			element.FailedField = err.StructNamespace()
			element.Tag = err.Tag()
			element.Value = err.Param()
			errors = append(errors, &element)
		}
	}
	return errors
}

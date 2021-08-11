package handler

import (
	"log"
	"net/http"
	"soulxsnips/src/model"
	"soulxsnips/src/service"

	"github.com/go-playground/validator"
	"github.com/gofiber/fiber/v2"
)

// CreateMD POST - creates a MarkdownSnippet from the provided body
// @Summary Create new a markdown snippet
// @Accept json
// @Produce json
// @Tags md
// @Success 201 {object} model.MarkdownSnippet
// @Failure 400 {object} model.ApiResponse
// @Failure 400 {object} model.ApiResponse
// @Failure 500 {object} model.ApiResponse
// @Router /md [post]
// @Param message body model.CreateMarkdownSnippet true "Post Body"
func CreateMD(ctx *fiber.Ctx) error {
	snippetRequest := new(model.CreateMarkdownSnippet)
	if err := ctx.BodyParser(snippetRequest); err != nil {
		log.Printf("Failed to parse CreateMarkdownSnippet: %s", err)
		return ctx.JSON(model.ApiResponse{Code: http.StatusBadRequest, Message: err.Error()})
	}

	if errs := ValidateStruct(snippetRequest); errs != nil {
		ctx.Status(http.StatusBadRequest)
		return ctx.JSON(errs)
	}

	newSnippet, err := service.CreateMarkdownSnippet(snippetRequest)
	if err != nil {
		log.Printf("Failed in insert new MarkdownSnippet: %s", err)
		return ctx.JSON(model.ApiResponse{Code: http.StatusInternalServerError, Message: err.Error()})
	}

	ctx.Status(http.StatusCreated)
	return ctx.JSON(newSnippet)
}

// GetMD GET - MarkdownSnippet Retrieval
// @Summary Retrieve Markdown Snippet
// @Accept json
// @Produce json
// @Tags md
// @Success 200 {object} model.MarkdownSnippet
// @Failure 400 {object} model.ApiResponse
// @Failure 500 {object} model.ApiResponse
// @Router /md/{id} [get]
// @Param id path string true "Snippet uuid"
func GetMD(ctx *fiber.Ctx) error {
	id := ctx.Params("id")

	snippet, err := service.GetMarkdownSnippet(id)
	if err != nil {
		log.Printf("Error Retrieving Markdown Snippet %s: %s", id, err)
		ctx.Status(http.StatusInternalServerError)
		return ctx.JSON(model.ApiResponse{Code: http.StatusInternalServerError, Message: err.Error()})
	}
	if snippet == nil && err == nil {
		ctx.Status(http.StatusNotFound)
		return ctx.JSON(model.ApiResponse{Code: http.StatusNotFound, Message: "Markdown Snippet Not Found"})
	}

	return ctx.JSON(snippet)
}

// GetAllMD GET - Get All MarkdownSnippets Retrieval
// @Summary Retrieve All Markdown Snippets
// @Accept json
// @Produce json
// @Tags md
// @Success 200 {object} []model.MarkdownSnippetListItem
// @Failure 400 {object} model.ApiResponse
// @Failure 500 {object} model.ApiResponse
// @Router /md [get]
func GetAllMD(ctx *fiber.Ctx) error {
	snippets, err := service.GetAllMarkdownSnippets()
	if err != nil {
		log.Printf("Error Retrieving All Markdown Snippet: %s", err)
		ctx.Status(http.StatusInternalServerError)
		return ctx.JSON(model.ApiResponse{Code: http.StatusBadRequest, Message: err.Error()})
	}

	return ctx.JSON(snippets)
}

// UpdateMD PATCH - Updates a MarkdownSnippet
// @Summary Updates a markdown snippet
// @Accept json
// @Produce json
// @Tags md
// @Success 200 {object} model.MarkdownSnippet
// @Failure 400 {object} model.ApiResponse
// @Failure 404 {object} model.ApiResponse
// @Failure 500 {object} model.ApiResponse
// @Router /md [patch]
// @Param message body model.UpdateMarkdownSnippet true "Patch Body"
func UpdateMD(ctx *fiber.Ctx) error {
	patchSnippet := new(model.UpdateMarkdownSnippet)
	if err := ctx.BodyParser(patchSnippet); err != nil {
		ctx.Status(http.StatusBadRequest)
		return ctx.JSON(model.ApiResponse{Code: http.StatusBadRequest, Message: err.Error()})
	}

	if errs := ValidateStruct(patchSnippet); errs != nil {
		ctx.Status(http.StatusBadRequest)
		return ctx.JSON(errs)
	}

	updatedSnippet, err := service.UpdateMarkdownSnippet(patchSnippet)
	if err != nil {
		log.Printf("Failed in update MarkdownSnippet %s: %s", patchSnippet.ID, err)
		ctx.Status(http.StatusBadRequest)
		return ctx.JSON(model.ApiResponse{Code: http.StatusInternalServerError, Message: err.Error()})
	}

	return ctx.JSON(updatedSnippet)
}

// DeleteMD DELETE - Removes MarkdownSnippet permanantly
// @Summary Removes MarkdownSnippet permanantly
// @Accept json
// @Produce json
// @Tags md
// @Success 204
// @Failure 400 {object} model.ApiResponse
// @Router /md/{id} [delete]
// @Param id path string true "Snippet uuid"
func DeleteMD(ctx *fiber.Ctx) error {
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

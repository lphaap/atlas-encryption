package api

import (
	"time"

	"github.com/gofiber/fiber/v2"
)

type ResponseData struct {
	At      string `json:"at"`
	Success bool   `json:"success"`
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    any    `json:"data"`
}

func (response ResponseData) Send(fiber *fiber.Ctx) error {
	return fiber.Status(response.Code).JSON(response)
}

func (response ResponseData) SendError(fiber *fiber.Ctx, code int, message string, data any) error {
	response.Success = false
	response.Code = code
	response.Message = message
	response.Data = data
	response.At = time.Now().Format("2006-01-02 15:04:05")
	return response.Send(fiber)
}

func (response ResponseData) SendSuccess(fiber *fiber.Ctx, code int, message string, data any) error {
	response.Success = true
	response.Code = code
	response.Message = message
	response.Data = data
	response.At = time.Now().Format("2006-01-02 15:04:05")
	return response.Send(fiber)
}

func (response ResponseData) OK(fiber *fiber.Ctx, data any) error {
	return response.SendSuccess(fiber, 200, "OK", data)
}

func (response ResponseData) Created(fiber *fiber.Ctx, data any) error {
	return response.SendSuccess(fiber, 201, "Created", data)
}

func (response ResponseData) Accepted(fiber *fiber.Ctx, data any) error {
	return response.SendSuccess(fiber, 202, "Accepted", data)
}

func (response ResponseData) BadRequest(fiber *fiber.Ctx, data any) error {
	return response.SendError(fiber, 400, "Bad Request", data)
}

func (response ResponseData) Unauthorized(fiber *fiber.Ctx, data any) error {
	return response.SendError(fiber, 401, "Unauthorized", data)
}

func (response ResponseData) NotFound(fiber *fiber.Ctx, data any) error {
	return response.SendError(fiber, 404, "Not Found", data)
}

func (response ResponseData) InternalServerError(fiber *fiber.Ctx, data any) error {
	return response.SendError(fiber, 500, "Internal Server Error", data)
}

var Response = ResponseData{}

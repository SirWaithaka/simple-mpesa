package error_handlers

import (
	"log"

	"simple-mpesa/app/errors"

	"github.com/gofiber/fiber/v2"
)

// ErrorHandler provides a custom error handling mechanism for fiber framework
func ErrorHandler(ctx *fiber.Ctx, err error) error {

	// if error corresponds to unauthorized
	if e, ok := err.(errors.Unauthorized); ok {
		log.Println(err)
		res := errors.UnauthorizedResponse(e.Error())
		return ctx.Status(res.Status).JSON(res)
	}

	// if error is our custom validation errors slice type
	if e, ok := err.(errors.ValidationErrors); ok {
		log.Println(err)
		res := errors.BadRequestResponse(e.Error())
		return ctx.Status(res.Status).JSON(res)
	}

	if e, ok := err.(errors.Error); ok {
		// we first log the error
		log.Println(e)

		if errors.ErrorCode(e) == errors.EINTERNAL {
			res := errors.InternalServerError(e.Error())
			return ctx.Status(res.Status).JSON(res)
		} else if _, ok := e.Err.(errors.Unauthorized); ok {
			res := errors.UnauthorizedResponse(e.Error())
			return ctx.Status(res.Status).JSON(res)
		} else {
			res := errors.BadRequestResponse(e.Error())
			return ctx.Status(res.Status).JSON(res)
		}
	}

	// if its a fiber error we send back the status code and empty response
	if e, ok := err.(*fiber.Error); ok {
		ctx.Status(e.Code)
		return nil
	}

	// will catch any other error we dont process here and return status 500
	if err != nil {
		log.Println(err)
		msg := "Something has happened. Report Issue."
		res := errors.InternalServerError(msg)
		return ctx.Status(res.Status).JSON(res)
	}

	// Return from handler
	return nil
}

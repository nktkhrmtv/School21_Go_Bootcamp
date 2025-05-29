package handlers

/*
#include <stdlib.h>
#include "cow.h" 
*/
import "C"
import (
    "unsafe"
    "github.com/go-openapi/runtime/middleware"
    "day04/restapi/operations"
)

func Handle(params operations.BuyCandyParams) middleware.Responder {
    if *params.Order.Money < 0 || *params.Order.CandyCount < 0 || !isValidCandyType(*params.Order.CandyType) {
        return operations.NewBuyCandyBadRequest().WithPayload(&operations.BuyCandyBadRequestBody{Error: "Invalid input data"})
    }

    price := getPriceByCandyType(*params.Order.CandyType)
    totalCost := price * int(*params.Order.CandyCount)
    if int(*params.Order.Money) < totalCost {
        return operations.NewBuyCandyPaymentRequired().WithPayload(&operations.BuyCandyPaymentRequiredBody{Error: "Not enough money"})
    }

    change := int(*params.Order.Money) - totalCost

    phrase := C.CString("Thank you!") 
    defer C.free(unsafe.Pointer(phrase)) 
    cowMessage := C.ask_cow(phrase)
    defer C.free(unsafe.Pointer(cowMessage)) 
    thanksMessage := C.GoString(cowMessage) 
    return operations.NewBuyCandyCreated().WithPayload(&operations.BuyCandyCreatedBody{
        Thanks: thanksMessage,
        Change: int64(change),
    })
    // return operations.NewBuyCandyCreated().WithPayload(&operations.BuyCandyCreatedBody{
    //     Thanks: "Thank you!",
    //     Change: int64(change),
    // })
}

func isValidCandyType(candyType string) bool {
    switch candyType {
    case "CE", "AA", "NT", "DE", "YR":
        return true
    default:
        return false
    }
}

func getPriceByCandyType(candyType string) int {
    switch candyType {
		case "CE":
			return 10
		case "AA":
			return 15
		case "NT":
			return 17
		case "DE":
			return 21
		case "YR":
			return 23
		default:
			return 0 
    }
}


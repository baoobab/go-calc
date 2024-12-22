package pkg

import (
	"calc/enums"
)

type CalcResponse struct { // Body ответа от сервера
	Result float64         `json:"result,omitempty"` // тк поле может быть опциональным - либо ошибки нет, есть результат
	Error  enums.ErrorCode `json:"error,omitempty"`  // тк поле может быть опциональным - либо ошибка есть, нет результата
}

type CalcRequest struct { // Body запроса на вход
	Expression string `json:"expression"`
}

package builtin_tool

import (
	"context"
	"errors"
	"math/big"
	"rag-new/internal/schema"
)

var calculatorAllowedMethods = []string{"add", "subtract", "multiply", "divide"}

type calculateParams struct {
	A      string `json:"a"`
	B      string `json:"b"`
	Method string `json:"method"`
}

func (s *Service) Calculator(_ context.Context, args schema.FunctionCallArguments) (*schema.CallBuiltInResponse, error) {
	var response = &schema.CallBuiltInResponse{}
	var params calculateParams
	err := args.Unmarshal(&params)
	if err != nil {
		return response, err
	}

	a := new(big.Float)
	b := new(big.Float)

	a, _, err = big.ParseFloat(params.A, 10, 0, big.ToZero)
	if err != nil {
		return response, errors.New("invalid value for A")
	}

	b, _, err = big.ParseFloat(params.B, 10, 0, big.ToZero)
	if err != nil {
		return response, errors.New("invalid value for B")
	}

	var result *big.Float

	switch params.Method {
	case "add":
		result = new(big.Float).Add(a, b)
	case "subtract":
		result = new(big.Float).Sub(a, b)
	case "multiply":
		result = new(big.Float).Mul(a, b)
	case "divide":
		if b.Cmp(big.NewFloat(0)) == 0 {
			response.Content = "cannot divide by zero"
			return response, errors.New(response.Content)
		}
		result = new(big.Float).Quo(a, b)
	default:
		response.Content = "invalid method"
		return response, errors.New(response.Content)
	}

	response.Content = result.String()
	return response, nil
}

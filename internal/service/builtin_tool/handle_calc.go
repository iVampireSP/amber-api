package builtin_tool

import (
	"context"
	"errors"
	"fmt"
	"math/big"
	"rag-new/internal/schema"
)

var calculatorAllowedMethods = []string{"add", "subtract", "multiply", "divide", "modulus"}

type calculateParams struct {
	NumberA  string `json:"number_a"  mapstructure:"number_a"`
	NumberB  string `json:"number_b"  mapstructure:"number_b"`
	Operator string `json:"operator"  mapstructure:"operator"`
}

func (s *Service) Calculator(_ context.Context, args schema.FunctionCallArguments) (*schema.CallBuiltInResponse, error) {
	var response = &schema.CallBuiltInResponse{}
	var params calculateParams
	err := args.Unmarshal(&params)
	if err != nil {
		return response, err
	}

	var additionalInfo string

	// 取模需要单独处理
	if params.Operator == "modulus" {
		a := new(big.Int)
		b := new(big.Int)

		// 将字符串转换为大整数
		a.SetString(params.NumberA, 10)
		b.SetString(params.NumberB, 10)

		// 进行取模运算
		result := new(big.Int).Mod(a, b)

		response.Content = fmt.Sprintf("%d", result)

	} else {
		a := new(big.Float)
		b := new(big.Float)

		a, _, err = big.ParseFloat(params.NumberA, 10, 0, big.ToZero)
		if err != nil {
			return response, errors.New("invalid value for Number A")
		}

		b, _, err = big.ParseFloat(params.NumberB, 10, 0, big.ToZero)
		if err != nil {
			return response, errors.New("invalid value for Number B")
		}

		var result *big.Float

		switch params.Operator {
		case "add":
			result = new(big.Float).Add(a, b)
		case "subtract":
			result = new(big.Float).Sub(a, b)
			if a.Cmp(b) == 0 {
				additionalInfo = "equal"
			} else if a.Cmp(b) > 0 {
				additionalInfo = params.NumberA + " greater than " + params.NumberB
			} else {
				additionalInfo = params.NumberA + " less than " + params.NumberB
			}

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

		response.Content = fmt.Sprintf("%.2f", result)
	}

	if additionalInfo != "" {
		response.Content += fmt.Sprintf(" (%s)", additionalInfo)
	}

	return response, nil
}

type compareParams struct {
	NumberA string `json:"number_a"  mapstructure:"number_a"`
	NumberB string `json:"number_b"  mapstructure:"number_b"`
}

//
//func (s *Service) Compare(_ context.Context, args schema.FunctionCallArguments) (*schema.CallBuiltInResponse, error) {
//	var response = &schema.CallBuiltInResponse{}
//	var params compareParams
//	err := args.Unmarshal(&params)
//	if err != nil {
//		return response, err
//	}
//
//	a := new(big.Float)
//	b := new(big.Float)
//
//	a, _, err = big.ParseFloat(params.NumberA, 10, 0, big.ToZero)
//	if err != nil {
//		return response, errors.New("invalid value for A")
//	}
//
//	b, _, err = big.ParseFloat(params.NumberB, 10, 0, big.ToZero)
//	if err != nil {
//		return response, errors.New("invalid value for B")
//	}
//
//	if a.Cmp(b) == 0 {
//		response.Content = "equal"
//	} else if a.Cmp(b) > 0 {
//		response.Content = params.NumberA + " greater than " + params.NumberA
//	} else {
//		response.Content = params.NumberB + " less than " + params.NumberB
//	}
//
//	return response, nil
//}

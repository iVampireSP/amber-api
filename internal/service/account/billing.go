package account

import (
	"bytes"
	"errors"
	"github.com/bytedance/sonic"
	"io"
	"net/http"
	"rag-new/internal/schema"
	"strconv"
)

func (s *Service) UnitReduce(userId schema.UserId, num int, unit string) error {
	err := s.request("/applications/balances/unit_reduce", UnitReduceRequest{
		Amount: strconv.Itoa(num),
		Reason: "UnitReduce",
		Unit:   unit,
		UserId: userId.String(),
	}, nil)

	if err != nil {
		s.logger.Sugar.Error(err)
		return err
	}

	return nil
}

func (s *Service) CanBillUnit(userId schema.UserId, unit string) (bool, error) {
	var canBill = CanBill{}

	err := s.request("/applications/balances/can_bill_unit", CanBillUnitRequest{
		Unit:   unit,
		UserId: userId.String(),
	}, &canBill)

	if err != nil {
		s.logger.Sugar.Error(err)
		return false, err
	}

	return canBill.CanBill, nil
}

func (s *Service) request(path string, data interface{}, output interface{}) error {
	dataJson, err := sonic.Marshal(data)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("POST", s.config.Account.Host+path, bytes.NewBuffer(dataJson))
	if err != nil {
		return err
	}

	req.Header.Set("Authorization", "Bearer "+s.config.Account.ApplicationKey)
	req.Header.Set("Content-Type", "application/json")

	rsp, err := http.DefaultClient.Do(req)
	if err != nil {
		s.logger.Sugar.Error(err)
		return err
	}

	defer func() {
		err = rsp.Body.Close()
		if err != nil {
			s.logger.Sugar.Error(err)
		}
	}()

	body, err := io.ReadAll(rsp.Body)
	if err != nil {
		return err
	}

	if output == nil {
		return nil
	}

	//if rsp.StatusCode  not 20x
	if rsp.StatusCode >= 200 && rsp.StatusCode < 300 {
		return sonic.Unmarshal(body, output)
	}

	return errors.New(string(body))
}

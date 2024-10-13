package batch

import (
	"rag-new/internal/base/conf"
	"rag-new/internal/dao"
	"rag-new/internal/service/account"
	"rag-new/internal/service/unsettled_token"
)

type UnsettedTokenBilling struct {
	AccountService       *account.Service
	UnsettedTokenService *unsettled_token.Service
	Config               *conf.Config
	DAO                  *dao.Query
}

func (*Batch) UnsettedTokenBilling(utb *UnsettedTokenBilling) error {
	unsettedTokens, err := utb.UnsettedTokenService.GetUnsettledTokenLargerThan(utb.Config.Account.UnitStart)
	if err != nil {
		return err
	}

	for _, ut := range unsettedTokens {
		// 除以 UnitStart
		iter := ut.Count / utb.Config.Account.UnitStart

		// 扣费
		err = utb.AccountService.UnitReduce(ut.UserId, int(iter), utb.Config.Account.Unit)
		if err != nil {
			return err
		}

		// 计算剩余的
		remain := ut.Count - (iter * utb.Config.Account.UnitStart)

		err = utb.UnsettedTokenService.DecreaseUnsettledToken(ut.UserId, remain)
		if err != nil {
			return err
		}
	}

	return nil
}

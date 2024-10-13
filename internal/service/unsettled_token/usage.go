package unsettled_token

import (
	"rag-new/internal/entity"
	"rag-new/internal/schema"
)

func (s *Service) IncreaseUnsettledToken(userId schema.UserId, count int64) error {
	var err error

	err = s.createRecordIfNotExists(userId)
	if err != nil {
		return err
	}

	current, err := s.GetUnsettledToken(userId)
	if err != nil {
		return err
	}

	if count < 0 {
		count = 0
	}

	count += current

	_, err = s.dao.UnsettledToken.Where(
		s.dao.UnsettledToken.UserId.Eq(userId.String()),
	).UpdateSimple(s.dao.UnsettledToken.Count_.Value(count))

	return nil
}

func (s *Service) DecreaseUnsettledToken(userId schema.UserId, count int64) error {
	var err error

	err = s.createRecordIfNotExists(userId)
	if err != nil {
		return err
	}

	current, err := s.GetUnsettledToken(userId)
	if err != nil {
		return err
	}

	if count < 0 {
		count = 0
	}

	count -= current

	if count < 0 {
		count = 0
	}

	_, err = s.dao.UnsettledToken.Where(
		s.dao.UnsettledToken.UserId.Eq(userId.String()),
	).UpdateSimple(s.dao.UnsettledToken.Count_.Value(count))

	return nil
}

func (s *Service) GetUnsettledToken(userId schema.UserId) (int64, error) {
	err := s.createRecordIfNotExists(userId)
	if err != nil {
		return 0, err
	}

	first, err := s.dao.UnsettledToken.Where(
		s.dao.UnsettledToken.UserId.Eq(userId.String()),
	).First()
	if err != nil {
		return 0, err
	}

	return first.Count, nil
}

// GetUnsettledTokenLargerThan count
func (s *Service) GetUnsettledTokenLargerThan(count int64) ([]*entity.UnsettledToken, error) {
	unsettledTokens, err := s.dao.UnsettledToken.Where(
		s.dao.UnsettledToken.Count_.Gte(count),
	).Find()
	if err != nil {
		return nil, err
	}

	return unsettledTokens, nil
}

func (s *Service) createRecordIfNotExists(userId schema.UserId) error {
	count, err := s.dao.UnsettledToken.Where(
		s.dao.UnsettledToken.UserId.Eq(userId.String()),
	).Count()
	if err != nil {
		return err
	}

	if count == 0 {
		return s.dao.UnsettledToken.Create(&entity.UnsettledToken{
			UserId: userId,
			Count:  0,
		})
	}

	return nil
}

package embedding

import (
	"context"
	"rag-new/internal/entity"
	"rag-new/pkg/md5"
)

func (s *Service) TextEmbedding(ctx context.Context, input []string) ([][]float32, error) {

	var r = make([][]float32, len(input)-1)

	for _, v := range input {
		embedding, err := s.getCache(ctx, v)
		if err != nil {
			return r, err
		}

		if embedding != nil {
			r = append(r, embedding)
			continue
		}

		embedding2, err := s.OpenAI.CreateEmbedding(ctx, []string{v})
		if err != nil {
			return nil, err
		}

		r = append(r, embedding2[0])

		err = s.setCache(ctx, v, embedding2[0])
		if err != nil {
			return nil, err
		}
	}

	return r, nil
}

func (s *Service) getCache(ctx context.Context, input string) ([]float32, error) {
	md5Str, err := md5.Md5(input)
	if err != nil {
		return nil, err
	}

	c, err := s.dao.WithContext(ctx).Embedding.Where(s.dao.Embedding.TextMd5.Eq(md5Str)).
		Where(s.dao.Embedding.EmbeddingModel.Eq(s.config.OpenAI.EmbeddingModel)).
		Count()
	if c == 0 {
		return nil, err
	}

	first, err := s.dao.WithContext(ctx).Embedding.Where(s.dao.Embedding.TextMd5.Eq(md5Str)).
		Where(s.dao.Embedding.EmbeddingModel.Eq(s.config.OpenAI.EmbeddingModel)).
		First()
	if err != nil {
		return nil, err
	}

	// byte to float32
	return first.Vector, nil
}

func (s *Service) setCache(ctx context.Context, input string, embedding []float32) error {
	md5Str, err := md5.Md5(input)
	if err != nil {
		return err
	}

	return s.dao.WithContext(ctx).Embedding.Create(&entity.Embedding{
		Text:           input,
		TextMd5:        md5Str,
		Vector:         embedding,
		EmbeddingModel: s.config.OpenAI.EmbeddingModel,
	})
}

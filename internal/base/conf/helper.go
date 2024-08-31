package conf

func (o *OpenAI) CanUse(model string) bool {
	if o.AllowedModels == nil {
		return true
	}

	for _, allowedModel := range o.AllowedModels {
		if allowedModel == model {
			return true
		}
	}

	return false
}

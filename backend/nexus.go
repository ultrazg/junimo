package backend

func NewNexus() *Nexus {
	return &Nexus{}
}

func (n *Nexus) ValidateUser(apiKey string) (*NexusUserValidateResult, error) {
	client, err := NewClient(apiKey)
	if err != nil {
		return nil, err
	}

	result := &NexusUserValidateResult{}
	err = client.GetJSON(ApiUsersValidate, result)
	if err != nil {
		return nil, err
	}

	return result, nil
}

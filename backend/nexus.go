package backend

import (
	"fmt"

	"github.com/spf13/viper"
)

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

func (n *Nexus) ViewSpecifiedModFile(modID string) (*NexusViewSpecifiedModFileResult, error) {
	apiKey := viper.GetString("nexus_api_key")
	if apiKey == "" {
		return nil, fmt.Errorf("请检查 Nexus Mods API Key 是否正确")
	}

	client, err := NewClient(apiKey)
	if err != nil {
		return nil, err
	}

	result := &NexusViewSpecifiedModFileResult{}
	err = client.GetJSON(fmt.Sprintf(ApiViewSpecifiedModFile, modID), result)
	if err != nil {
		return nil, err
	}

	return result, nil
}

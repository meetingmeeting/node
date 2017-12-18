package nats_dialog

import (
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"testing"
	"github.com/mysterium/node/identity"
)

func TestRequestSerialize(t *testing.T) {
	var identity = identity.NewIdentity("123")
	var tests = []struct {
		model        dialogCreateRequest
		expectedJson string
	}{
		{
			dialogCreateRequest{
				IdentityId: identity.Id,
			},
			`{
				"identity_id": "123"
			}`,
		},
		{
			dialogCreateRequest{},
			`{
				"identity_id": ""
			}`,
		},
	}

	for _, test := range tests {
		jsonBytes, err := json.Marshal(test.model)

		assert.NoError(t, err)
		assert.JSONEq(t, test.expectedJson, string(jsonBytes))
	}
}

func TestRequestUnserialize(t *testing.T) {
	var tests = []struct {
		json          string
		expectedModel dialogCreateRequest
		expectedError error
	}{
		{
			`{
				"identity_id": "123"
			}`,
			dialogCreateRequest{
				IdentityId: identity.NewIdentity("123").Id,
			},
			nil,
		},
		{
			`{}`,
			dialogCreateRequest{
				IdentityId: identity.NewIdentity("").Id,
			},
			nil,
		},
	}

	for _, test := range tests {
		var model dialogCreateRequest
		err := json.Unmarshal([]byte(test.json), &model)

		assert.Exactly(t, test.expectedModel, model)
		if test.expectedError != nil {
			assert.EqualError(t, err, test.expectedError.Error())
		} else {
			assert.NoError(t, err)
		}
	}
}

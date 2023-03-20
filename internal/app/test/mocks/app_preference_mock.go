package mocks

import context "context"

type mockAppPreference struct {
	values map[string]interface{}
}

func GetAppPreferenceMock(value map[string]interface{}) mockAppPreference {
	return mockAppPreference{value}
}

func (mock mockAppPreference) GetValue(ctx context.Context, key string, value interface{}) interface{} {
	if v, f := mock.values[key]; f {
		return v
	}
	return value
}

func (mock mockAppPreference) GetBool(ctx context.Context, key string) bool {
	if v, f := mock.values[key]; f {
		return v.(bool)
	}
	return false
}

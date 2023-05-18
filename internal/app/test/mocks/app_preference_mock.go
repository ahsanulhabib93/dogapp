package mocks

import "context"

type mockAppPreference struct {
	values map[string]interface{}
}

func GetAppPreferenceMock(value map[string]interface{}) mockAppPreference {
	return mockAppPreference{value}
}

func (mock mockAppPreference) GetBool(ctx context.Context, key string) bool {
	// TODO
	// need to define the mock function which got added due to goFramework update
	return true
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

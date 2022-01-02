package inmemory

import (
	"context"
	"testing"
)

func TestSave(t *testing.T) {

	valueStore := NewInMemory()

	for _, tc := range []struct {
		name    string
		input   string
		wantErr bool
	}{
		{
			name:  "Save simple value",
			input: "some random text input",
		},
		{
			name:  "empty string test case",
			input: "",
		},
		{
			name:  "url input",
			input: "https://hunterheston.com",
		},
		{
			name:  "more input",
			input: "https://hunterheston.com",
		},
		{
			name:  "another input",
			input: "https://hunterheston.com",
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			ctx := context.Background()
			got, err := valueStore.Save(ctx, []byte(tc.input))
			if (err != nil) != tc.wantErr {
				t.Errorf("Unexpected InMemory.Save(%q) got err: %v", tc.input, err)
			}
			if !tc.wantErr && got == "" {
				t.Errorf("InMemroy.Save(%q)=%q want=not an empty string", tc.input, got)
			}
		})
	}
}

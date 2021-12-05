package validation

import "testing"

func TestValidateURL(t *testing.T) {

	for _, tc := range []struct {
		name  string
		input string
		want  bool
	}{
		{
			name:  "valid url",
			input: "https://hunterheston.com",
			want:  true,
		},
		{
			name:  "empty url",
			input: "",
		},
		{
			name:  "missing : after https",
			input: "https://hunterheston.com",
			want:  true,
		},
		{
			name:  "http protocol",
			input: "http://hunterheston.com",
			want:  true,
		},
		{
			name:  "random text",
			input: "kjhsdaflkadflkjhas;ldkfj",
		},
		{
			name:  "email",
			input: "testemail@gmail.com",
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			got := ValidateURL(tc.input)
			if got != tc.want {
				t.Errorf("ValidateURL(%v)=%v, want=%v", tc.input, got, tc.want)
			}
		})
	}

}

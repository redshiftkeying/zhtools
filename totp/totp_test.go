package totp

import (
	"testing"
)

func TestGenerateCode(t *testing.T) {
	tests := []struct {
		name       string
		secret     string
		digitCount int
		timeStep   uint64
		t0         uint64
		unixTime   uint64
		wantErr    bool
	}{
		{
			name:       "Valid code generation",
			secret:     "JBSWY3DPEHPK3PXP",
			digitCount: 6,
			timeStep:   30,
			t0:         0,
			unixTime:   1660000000,
			wantErr:    false,
		},
		{
			name:       "Empty secret",
			secret:     "",
			digitCount: 6,
			timeStep:   30,
			t0:         0,
			unixTime:   1660000000,
			wantErr:    true,
		},
		{
			name:       "Invalid secret base32",
			secret:     "INVALIDSECRET",
			digitCount: 6,
			timeStep:   30,
			t0:         0,
			unixTime:   1660000000,
			wantErr:    true,
		},
		{
			name:       "DigitCount less than 1",
			secret:     "JBSWY3DPEHPK3PXP",
			digitCount: 0,
			timeStep:   30,
			t0:         0,
			unixTime:   1660000000,
			wantErr:    true,
		},
		{
			name:       "DigitCount exceeds 10",
			secret:     "JBSWY3DPEHPK3PXP",
			digitCount: 11,
			timeStep:   30,
			t0:         0,
			unixTime:   1660000000,
			wantErr:    true,
		},
		{
			name:       "UnixTime before T0",
			secret:     "JBSWY3DPEHPK3PXP",
			digitCount: 6,
			timeStep:   30,
			t0:         1000000000,
			unixTime:   999999999,
			wantErr:    false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			totp := &TOTP{
				Secret:     tt.secret,
				DigitCount: tt.digitCount,
				TimeStep:   tt.timeStep,
				T0:         tt.t0,
			}
			got, err := totp.GenerateCode(tt.unixTime)
			if (err != nil) != tt.wantErr {
				t.Errorf("GenerateCode() error = %v, wantErr = %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && len(got) != tt.digitCount {
				t.Errorf("GenerateCode() = %v, expected code length = %d", got, tt.digitCount)
			}
		})
	}
}

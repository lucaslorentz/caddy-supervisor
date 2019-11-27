package supervisor

import (
	"github.com/caddyserver/caddy"
	"github.com/go-errors/errors"
	"testing"
)

func TestParseOption(t *testing.T) {
	type args struct {
		c       *caddy.Controller
		options *Options
	}
	tests := []struct {
		name                      string
		args                      args
		expectedParseOptionReturn bool
		assertOptions             func(*Options) error
	}{
		{"env with variable", args{caddy.NewTestController("http", "supervisor {\nenv NUMBER={{ add 1 .Replica }}\n}"), new(Options)}, true, func(options *Options) error {
			if len(options.Env) != 1 {
				return errors.Errorf("Env should be exactly one and returned %v", len(options.Env))
			}

			expected := "NUMBER={{ add 1 .Replica }}"
			returned := options.Env[0]
			if expected != returned {
				return errors.Errorf("env should be %v but returned %v", expected, returned)
			}

			return nil
		}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.args.c.RemainingArgs()
			for tt.args.c.NextBlock() {
				if parseOptionReturn := ParseOption(tt.args.c, tt.args.options); parseOptionReturn != tt.expectedParseOptionReturn {
					t.Errorf("ParseOption should return %v and returned %v", tt.expectedParseOptionReturn, parseOptionReturn)
				}

				if err := tt.assertOptions(tt.args.options); err != nil {
					t.Error(err)
				}
			}
		})
	}
}

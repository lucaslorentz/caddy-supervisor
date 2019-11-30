package supervisor

import (
	"testing"

	"github.com/caddyserver/caddy"
	"github.com/go-errors/errors"
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
		{"env with variable 2 args", args{caddy.NewTestController("http", "supervisor {\nenv NUMBER \"{{ add 1 .Replica }}\"\n}"), CreateOptions()}, true, func(options *Options) error {
			if len(options.Env) != 1 {
				return errors.Errorf("Env should be exactly one and returned %v", len(options.Env))
			}

			expected := "NUMBER={{ add 1 .Replica }}"
			returned := options.Env[0]
			if expected != returned {
				return errors.Errorf("env should be %v but returned %v", expected, returned)
			}

			supervisors := CreateSupervisors(options)
			expected = "NUMBER=1"
			returned = supervisors[0].options.Env[0]
			if expected != returned {
				return errors.Errorf("env should be %v but returned %v", expected, returned)
			}

			return nil
		}},
		{"env with variable 1 arg", args{caddy.NewTestController("http", "supervisor {\nenv \"NUMBER={{ add 1 .Replica }}\"\n}"), CreateOptions()}, true, func(options *Options) error {
			if len(options.Env) != 1 {
				return errors.Errorf("Env should be exactly one and returned %v", len(options.Env))
			}

			expected := "NUMBER={{ add 1 .Replica }}"
			returned := options.Env[0]
			if expected != returned {
				return errors.Errorf("env should be %v but returned %v", expected, returned)
			}

			supervisors := CreateSupervisors(options)
			expected = "NUMBER=1"
			returned = supervisors[0].options.Env[0]
			if expected != returned {
				return errors.Errorf("env should be %v but returned %v", expected, returned)
			}

			return nil
		}},
		{"env variable 1 arg without =", args{caddy.NewTestController("http", "supervisor {\nenv \"NUMBER {{ add 1 .Replica }}\"\n}"), CreateOptions()}, false, func(options *Options) error {
			if len(options.Env) != 0 {
				return errors.Errorf("Env should be exactly zero and returned %v", len(options.Env))
			}

			return nil
		}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.args.c.RemainingArgs()
			processedBlocks := 0
			for tt.args.c.NextBlock() {
				processedBlocks++
				if parseOptionReturn := ParseOption(tt.args.c, tt.args.options); parseOptionReturn != tt.expectedParseOptionReturn {
					t.Errorf("ParseOption should return %v and returned %v", tt.expectedParseOptionReturn, parseOptionReturn)
				}

				if err := tt.assertOptions(tt.args.options); err != nil {
					t.Error(err)
				}
			}

			if processedBlocks == 0 {
				t.Errorf("this test case seems to be empty as no block was processed")
			}
		})
	}
}

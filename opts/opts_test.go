package opts

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type TC struct {
	name           string
	inputBuilder   ConfigBuilder
	inputOpts      []Option
	expectedConfig Config
}

func TestInit(t *testing.T) {
	testcases := []TC{
		{
			name:           "default",
			inputBuilder:   ConfigBuilder{},
			inputOpts:      []Option{},
			expectedConfig: Config{},
		},
		{
			name:           "bulider sets name",
			inputBuilder:   ConfigBuilder{}.AppName("foo"),
			inputOpts:      []Option{},
			expectedConfig: Config{appName: "foo"},
		},
		{
			name:         "bulider sets name, option overwrites",
			inputBuilder: ConfigBuilder{}.AppName("foo"),
			inputOpts: []Option{
				WithApp("bar"),
			},
			expectedConfig: Config{appName: "bar"},
		},
		{
			name:         "option sets name",
			inputBuilder: ConfigBuilder{},
			inputOpts: []Option{
				WithApp("baz"),
			},
			expectedConfig: Config{appName: "baz"},
		},
		{
			name:           "bulider sets wait",
			inputBuilder:   ConfigBuilder{}.Wait(true),
			inputOpts:      []Option{},
			expectedConfig: Config{wait: true},
		},
		{
			name:         "option sets wait",
			inputBuilder: ConfigBuilder{},
			inputOpts: []Option{
				Wait(),
			},
			expectedConfig: Config{wait: true},
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			actualConfig := Init(tc.inputBuilder, tc.inputOpts...)
			assert.Equal(t, tc.expectedConfig.AppName(), actualConfig.AppName())
			assert.Equal(t, tc.expectedConfig.Wait(), actualConfig.Wait())
		})
	}
}

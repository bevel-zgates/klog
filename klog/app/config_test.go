package app

import (
	"testing"

	"github.com/jotaen/klog/klog"
	tf "github.com/jotaen/klog/klog/app/cli/terminalformat"
	"github.com/jotaen/klog/klog/service"
	"github.com/stretchr/testify/assert"
)

func createMockConfigFromEnv(vs map[string]string) FromEnvVars {
	return FromEnvVars{GetVar: func(n string) string {
		return vs[n]
	}}
}

func TestCreatesNewDefaultConfig(t *testing.T) {
	c := NewDefaultConfig(tf.COLOUR_THEME_NO_COLOUR)
	assert.Equal(t, c.IsDebug.Value(), false)
	assert.Equal(t, c.Editor.UnwrapOr(""), "")
	assert.Equal(t, c.CpuKernels.Value(), 1)
	assert.Equal(t, c.ColourScheme.Value(), tf.COLOUR_THEME_NO_COLOUR)
	assert.False(t, c.HideWarnings.UnwrapOr(false))

	isRoundingSet := false
	c.DefaultRounding.Unwrap(func(_ service.Rounding) {
		isRoundingSet = true
	})
	assert.False(t, isRoundingSet)

	isShouldTotalSet := false
	c.DefaultShouldTotal.Unwrap(func(_ klog.ShouldTotal) {
		isShouldTotalSet = true
	})
	assert.False(t, isShouldTotalSet)
}

func TestSetsParamsMetadataIsHandledCorrectly(t *testing.T) {
	{
		c := NewDefaultConfig(tf.COLOUR_THEME_NO_COLOUR)
		assert.Equal(t, c.IsDebug.Value(), false)
	}
	{
		c, _ := NewConfig(
			FromDeterminedValues{NumCpus: 1},
			createMockConfigFromEnv(map[string]string{
				"KLOG_DEBUG": "1",
			}),
			FromConfigFile{""},
		)
		assert.Equal(t, c.IsDebug.Value(), true)
	}
}

func TestSetsParamsFromEnv(t *testing.T) {
	t.Run("Read plain environment variables.", func(t *testing.T) {
		c, _ := NewConfig(
			FromDeterminedValues{NumCpus: 1},
			createMockConfigFromEnv(map[string]string{
				"EDITOR":     "subl",
				"KLOG_DEBUG": "1",
				"NO_COLOR":   "1",
			}),
			FromConfigFile{""},
		)
		assert.Equal(t, c.IsDebug.Value(), true)
		assert.Equal(t, c.ColourScheme.Value(), tf.COLOUR_THEME_NO_COLOUR)
		assert.Equal(t, c.Editor.UnwrapOr(""), "subl")
	})

	t.Run("`$EDITOR` env variable trumps `editor` setting from config file.", func(t *testing.T) {
		c, _ := NewConfig(
			FromDeterminedValues{NumCpus: 1},
			createMockConfigFromEnv(map[string]string{
				"EDITOR": "subl",
			}),
			FromConfigFile{"editor = vi"},
		)
		assert.Equal(t, "subl", c.Editor.UnwrapOr(""))
	})

	t.Run("`$NO_COLOR` env variable trumps `colour_scheme = dark` from config file.", func(t *testing.T) {
		// This is important, otherwise you wouldn’t be able to override the colour scheme
		// e.g. for programmatic usage of klog.
		c, _ := NewConfig(
			FromDeterminedValues{NumCpus: 1},
			createMockConfigFromEnv(map[string]string{
				"NO_COLOR": "1",
			}),
			FromConfigFile{"colour_scheme = dark"},
		)
		assert.Equal(t, tf.COLOUR_THEME_NO_COLOUR, c.ColourScheme.Value())
	})
}

func TestSetsColourSchemeParamFromConfigFile(t *testing.T) {
	for _, x := range []struct {
		cfg string
		exp tf.ColourTheme
	}{
		{`colour_scheme = dark`, tf.COLOUR_THEME_DARK},
		{`colour_scheme = light`, tf.COLOUR_THEME_LIGHT},
		{`colour_scheme = no_colour`, tf.COLOUR_THEME_NO_COLOUR},
	} {
		c, _ := NewConfig(
			FromDeterminedValues{NumCpus: 1},
			createMockConfigFromEnv(map[string]string{}),
			FromConfigFile{x.cfg},
		)
		assert.Equal(t, x.exp, c.ColourScheme.Value())
	}
}

func TestSetsDefaultRoundingParamFromConfigFile(t *testing.T) {
	for _, x := range []struct {
		cfg string
		exp int
	}{
		{`default_rounding = 5m`, 5},
		{`default_rounding = 10m`, 10},
		{`default_rounding = 15m`, 15},
		{`default_rounding = 30m`, 30},
		{`default_rounding = 60m`, 60},
	} {
		c, _ := NewConfig(
			FromDeterminedValues{NumCpus: 1},
			createMockConfigFromEnv(map[string]string{}),
			FromConfigFile{x.cfg},
		)
		var value int
		c.DefaultRounding.Unwrap(func(r service.Rounding) {
			value = r.ToInt()
		})
		assert.Equal(t, x.exp, value)
	}
}

func TestSetsDefaultShouldTotalParamFromConfigFile(t *testing.T) {
	for _, x := range []struct {
		cfg string
		exp string
	}{
		{`default_should_total = 8h30m!`, "8h30m!"},
	} {
		c, _ := NewConfig(
			FromDeterminedValues{NumCpus: 1},
			createMockConfigFromEnv(map[string]string{}),
			FromConfigFile{x.cfg},
		)
		var value string
		c.DefaultShouldTotal.Unwrap(func(s klog.ShouldTotal) {
			value = s.ToString()
		})
		assert.Equal(t, x.exp, value)
	}
}

func TestSetsDateFormatParamFromConfigFile(t *testing.T) {
	for _, x := range []struct {
		cfg string
		exp bool
	}{
		{`date_format = YYYY-MM-DD`, true},
		{`date_format = YYYY/MM/DD`, false},
	} {
		c, _ := NewConfig(
			FromDeterminedValues{NumCpus: 1},
			createMockConfigFromEnv(map[string]string{}),
			FromConfigFile{x.cfg},
		)
		var value bool
		c.DateUseDashes.Unwrap(func(s bool) {
			value = s
		})
		assert.Equal(t, x.exp, value)
	}
}

func TestSetTimeFormatParamFromConfigFile(t *testing.T) {
	for _, x := range []struct {
		cfg string
		exp bool
	}{
		{`time_convention = 24h`, true},
		{`time_convention = 12h`, false},
	} {
		c, _ := NewConfig(
			FromDeterminedValues{NumCpus: 1},
			createMockConfigFromEnv(map[string]string{}),
			FromConfigFile{x.cfg},
		)
		var value bool
		c.TimeUse24HourClock.Unwrap(func(s bool) {
			value = s
		})
		assert.Equal(t, x.exp, value)
	}
}

func TestIgnoresUnknownPropertiesInConfigFile(t *testing.T) {
	for _, tml := range []string{`
unknown_property = 1
what_is_this = true
`,
	} {
		_, err := NewConfig(
			FromDeterminedValues{NumCpus: 1},
			createMockConfigFromEnv(map[string]string{}),
			FromConfigFile{tml},
		)
		assert.Nil(t, err)
	}
}

func TestIgnoresEmptyConfigFileOrEmptyParameters(t *testing.T) {
	for _, tml := range []string{
		``,
		`editor = `,
		`colour_scheme = `,
		`default_rounding =`,
		`default_should_total = `,
		`date_format = `,
		`time_convention = `,
	} {
		_, err := NewConfig(
			FromDeterminedValues{NumCpus: 1},
			createMockConfigFromEnv(map[string]string{}),
			FromConfigFile{tml},
		)
		assert.Nil(t, err)
	}
}

func TestRejectsInvalidConfigFile(t *testing.T) {
	for _, tml := range []string{
		`default_rounding = true`,              // Wrong type
		`default_rounding = 25m`,               // Invalid value
		`colour_scheme = true`,                 // Wrong type
		`colour_scheme = yellow`,               // Invalid value
		`default_should_total = [true, false]`, // Wrong type
		`default_should_total = 15`,            // Invalid value
		`date_format = [true, false]`,          // Wrong type
		`date_format = YYYY.MM.DD`,             // Invalid value
		`time_convention = [true, false]`,      // Wrong type
		`time_convention = 2h`,                 // Invalid value
	} {
		_, err := NewConfig(
			FromDeterminedValues{NumCpus: 1},
			createMockConfigFromEnv(map[string]string{}),
			FromConfigFile{tml},
		)
		assert.Error(t, err)
	}
}

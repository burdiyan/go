package cobrautil

import (
	"strings"

	"github.com/octago/sflags"
	"github.com/octago/sflags/gen/gpflag"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

type options struct {
	Strict bool
	Viper *viper.Viper
}

// Option is a type for options.
type Option func(*options)

// WithStrict enables the strict mode for configuration binding.
// This means that binding will fail if specified configuration file doesn't exist.
func WithStrict(enable bool) Option {
	return func(opts *options) {
		opts.Strict = enable
	}
}

// WithViper providers a custom Viper instance.
func WithViper(v *viper.Viper) Option {
	return func(opts *options) {
		opts.Viper = v
	}
}

// BindConfig creates cobra flags based on cfg, which must be a pointer to a configuration struct.
// Struct is then populated from flags -> env variables -> configuration file in this order.
// Configuration file must be named as the cmd name with one of the extensions supported by viper.
func BindConfig(cfg interface{}, cmd *cobra.Command, opt ...Option) error {
	opts := options {
		Viper: viper.New(),
	}
	
	for _, o := range opt {
		o(&opts)
	}

	fs := cmd.PersistentFlags()
	v := opts.Viper
	if err := gpflag.ParseTo(cfg, fs, sflags.FlagTag("mapstructure"), sflags.FlagDivider("."), sflags.Flatten(false)); err != nil {
		return err
	}

	v.AutomaticEnv()
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	if err := v.BindPFlags(fs); err != nil {
		return err
	}

	cmd.SetGlobalNormalizationFunc(func(fs *pflag.FlagSet, name string) pflag.NormalizedName {
		return pflag.NormalizedName(strings.Replace(name, "_", "-", -1))
	})

	var cfgFile string

	fs.StringVar(&cfgFile, "config", "", "path to the config file")

	cmd.PersistentPreRunE = func(cmd *cobra.Command, args []string) error {
		if cfgFile != "" {
			v.SetConfigFile(cfgFile)
		} else {
			v.SetConfigName(cmd.Name())
			v.AddConfigPath(".")
			v.AddConfigPath("/etc/")
		}


		if err := v.ReadInConfig(); opts.Strict && err != nil {
			return err
		}

		if err := v.Unmarshal(cfg); err != nil {
			return err
		}

		return nil
	}

	return nil
}
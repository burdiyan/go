package kongcli

import (
	"context"
	"io"
	"os"
	"strings"

	"github.com/alecthomas/kong"
	"gopkg.in/yaml.v2"
)

// DefaultOptions creates a set of opinionated options.
func DefaultOptions(ctx context.Context, name, description, version, envPrefix string, opts ...kong.Option) []kong.Option {
	os := []kong.Option{
		kong.Name(name),
		kong.Description(description),
		kong.UsageOnError(),
		BindContext(ctx),
		kong.Configuration(YAML, configsForApp(name)...),
		kong.Resolvers(EnvResolver(envPrefix)),
		kong.Vars{
			"version": version,
		},
	}

	os = append(os, opts...)

	return os
}

// YAML implements configuration loader for YAML files.
func YAML(r io.Reader) (kong.Resolver, error) {
	values := map[string]interface{}{}

	err := yaml.NewDecoder(r).Decode(&values)
	if err != nil {
		return nil, err
	}

	var f kong.ResolverFunc = func(context *kong.Context, parent *kong.Path, flag *kong.Flag) (interface{}, error) {
		name := strings.Replace(flag.Name, "-", "_", -1)

		raw, ok := values[name]
		if ok {
			return raw, nil
		}

		if !strings.Contains(name, ".") {
			return nil, nil
		}

		nested := strings.Split(name, ".")
		top, ok := values[nested[0]].(map[interface{}]interface{})
		if !ok {
			return nil, nil
		}

		raw, ok = top[nested[1]]
		if ok {
			return raw, nil
		}

		return nil, nil
	}

	return f, nil
}

// EnvResolver resolves flag values from the environment using and optional prefix.
func EnvResolver(prefix string) kong.Resolver {
	return kong.ResolverFunc(func(ctx *kong.Context, parent *kong.Path, flag *kong.Flag) (interface{}, error) {
		env := strings.ToUpper(flag.Name)
		env = strings.Replace(env, "-", "_", -1)
		env = strings.Replace(env, ".", "_", -1)

		if prefix != "" {
			env = prefix + env
		}

		v := os.Getenv(env)

		if v == "" {
			return nil, nil
		}

		return v, nil
	})
}

// BindContext binds the provided context so that it can be used as a parameter for `Run()` method in the commands.
func BindContext(ctx context.Context) kong.Option {
	return kong.BindTo(ctx, (*context.Context)(nil))
}

func configsForApp(name string) []string {
	return []string{
		"./" + name + ".yaml",
		"/etc/" + name + "/config.yaml",
	}
}

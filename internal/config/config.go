package config

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"gopkg.in/yaml.v3"
)

const (
	ConfigDirEnv         = "RSS_CLI_CONFIG_DIR"
	XDGConfigHomeEnv     = "XDG_CONFIG_HOME"
	AppDirName           = "rss-cli"
	ContextsDirName      = "contexts"
	ServicesDirName      = "services"
	CurrentFileName      = "current"
	ConfigFileName       = "config.yaml"
	CredentialsName      = "credentials.enc"
	CurrentConfigVersion = 1
)

type ContextConfig struct {
	Version int            `yaml:"version"`
	Tracing *TracingConfig `yaml:"tracing,omitempty"`
}

type TracingConfig struct {
	Enabled  bool   `yaml:"enabled"`
	Endpoint string `yaml:"endpoint"`
}

type ServiceConfig struct {
	Version  int    `yaml:"version"`
	Provider string `yaml:"provider"`
	BaseURL  string `yaml:"base_url"`
}

func BaseDirectory() (string, error) {
	if dir := strings.TrimSpace(os.Getenv(ConfigDirEnv)); dir != "" {
		return dir, nil
	}
	configHome := strings.TrimSpace(os.Getenv(XDGConfigHomeEnv))
	if configHome == "" {
		home, err := os.UserHomeDir()
		if err != nil {
			return "", fmt.Errorf("resolve home directory: %w", err)
		}
		configHome = filepath.Join(home, ".config")
	}
	return filepath.Join(configHome, AppDirName), nil
}

func ContextsDirectory() (string, error) {
	base, err := BaseDirectory()
	if err != nil {
		return "", err
	}
	return filepath.Join(base, ContextsDirName), nil
}
func ContextDirectory(name string) (string, error) {
	if err := ValidateName(name, "context"); err != nil {
		return "", err
	}
	base, err := ContextsDirectory()
	if err != nil {
		return "", err
	}
	return filepath.Join(base, name), nil
}
func ContextConfigPath(name string) (string, error) {
	dir, err := ContextDirectory(name)
	if err != nil {
		return "", err
	}
	return filepath.Join(dir, ConfigFileName), nil
}
func ServiceDirectory(contextName, service string) (string, error) {
	if err := ValidateName(service, "service"); err != nil {
		return "", err
	}
	dir, err := ContextDirectory(contextName)
	if err != nil {
		return "", err
	}
	return filepath.Join(dir, ServicesDirName, service), nil
}
func ServiceConfigPath(contextName, service string) (string, error) {
	dir, err := ServiceDirectory(contextName, service)
	if err != nil {
		return "", err
	}
	return filepath.Join(dir, ConfigFileName), nil
}
func CredentialsPath(contextName, service string) (string, error) {
	dir, err := ServiceDirectory(contextName, service)
	if err != nil {
		return "", err
	}
	return filepath.Join(dir, CredentialsName), nil
}

func ValidateName(name, kind string) error {
	if name == "" || name == "." || name == ".." || strings.ContainsAny(name, `/\\`) {
		return fmt.Errorf("invalid %s name %q", kind, name)
	}
	return nil
}

func SaveContext(name string, cfg ContextConfig) error {
	if cfg.Version == 0 {
		cfg.Version = CurrentConfigVersion
	}
	if cfg.Version != CurrentConfigVersion {
		return fmt.Errorf("unsupported context config version %d", cfg.Version)
	}
	path, err := ContextConfigPath(name)
	if err != nil {
		return err
	}
	if _, err := os.Stat(path); err == nil {
		return fmt.Errorf("context %q already exists", name)
	} else if !errors.Is(err, os.ErrNotExist) {
		return err
	}
	return writeYAML(path, cfg)
}

func LoadContext(name string) (ContextConfig, error) {
	path, err := ContextConfigPath(name)
	if err != nil {
		return ContextConfig{}, err
	}
	var cfg ContextConfig
	if err := readYAML(path, &cfg); err != nil {
		return cfg, err
	}
	if cfg.Version != CurrentConfigVersion {
		return cfg, fmt.Errorf("unsupported context config version %d", cfg.Version)
	}
	return cfg, nil
}

func UpdateContext(name string, cfg ContextConfig) error {
	if cfg.Version == 0 {
		cfg.Version = CurrentConfigVersion
	}
	if cfg.Version != CurrentConfigVersion {
		return fmt.Errorf("unsupported context config version %d", cfg.Version)
	}
	path, err := ContextConfigPath(name)
	if err != nil {
		return err
	}
	if _, err := os.Stat(path); errors.Is(err, os.ErrNotExist) {
		return fmt.Errorf("context %q does not exist", name)
	} else if err != nil {
		return err
	}
	return writeYAML(path, cfg)
}

func SaveService(contextName, service string, cfg ServiceConfig, force bool) error {
	if cfg.Version == 0 {
		cfg.Version = CurrentConfigVersion
	}
	if cfg.Version != CurrentConfigVersion {
		return fmt.Errorf("unsupported service config version %d", cfg.Version)
	}
	if err := ValidateName(cfg.Provider, "provider"); err != nil {
		return err
	}
	if cfg.Provider != service {
		return fmt.Errorf("service %q and provider %q must match", service, cfg.Provider)
	}
	path, err := ServiceConfigPath(contextName, service)
	if err != nil {
		return err
	}
	if !force {
		if _, statErr := os.Stat(path); statErr == nil {
			return fmt.Errorf("service %q is already configured", service)
		}
	}
	return writeYAML(path, cfg)
}

func LoadService(contextName, service string) (ServiceConfig, error) {
	path, err := ServiceConfigPath(contextName, service)
	if err != nil {
		return ServiceConfig{}, err
	}
	var cfg ServiceConfig
	if err := readYAML(path, &cfg); err != nil {
		return cfg, err
	}
	if cfg.Version != CurrentConfigVersion {
		return cfg, fmt.Errorf("unsupported service config version %d", cfg.Version)
	}
	if cfg.Provider != service {
		return cfg, fmt.Errorf("service %q and provider %q must match", service, cfg.Provider)
	}
	return cfg, nil
}

func SetCurrent(name string) error {
	if _, err := LoadContext(name); err != nil {
		return err
	}
	base, err := BaseDirectory()
	if err != nil {
		return err
	}
	if err := os.MkdirAll(base, 0o700); err != nil {
		return fmt.Errorf("create config directory: %w", err)
	}
	return os.WriteFile(filepath.Join(base, CurrentFileName), []byte(name+"\n"), 0o600)
}
func Current() (string, error) {
	base, err := BaseDirectory()
	if err != nil {
		return "", err
	}
	data, err := os.ReadFile(filepath.Join(base, CurrentFileName))
	if errors.Is(err, os.ErrNotExist) {
		return "", nil
	}
	if err != nil {
		return "", fmt.Errorf("read current context: %w", err)
	}
	return strings.TrimSpace(string(data)), nil
}

func ResolveContext(override string) (string, error) {
	name := override
	if name == "" {
		var err error
		name, err = Current()
		if err != nil {
			return "", err
		}
	}
	if name == "" {
		return "", fmt.Errorf("no current context; create one with: rss-cli context create NAME")
	}
	if _, err := LoadContext(name); err != nil {
		return "", err
	}
	return name, nil
}

func ListContexts() ([]string, error) {
	dir, err := ContextsDirectory()
	if err != nil {
		return nil, err
	}
	entries, err := os.ReadDir(dir)
	if errors.Is(err, os.ErrNotExist) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	names := make([]string, 0)
	for _, entry := range entries {
		if entry.IsDir() {
			names = append(names, entry.Name())
		}
	}
	sort.Strings(names)
	return names, nil
}
func DeleteContext(name string) error {
	dir, err := ContextDirectory(name)
	if err != nil {
		return err
	}
	if _, err := os.Stat(dir); errors.Is(err, os.ErrNotExist) {
		return fmt.Errorf("context %q does not exist", name)
	} else if err != nil {
		return err
	}
	if err := os.RemoveAll(dir); err != nil {
		return err
	}
	current, _ := Current()
	if current == name {
		base, baseErr := BaseDirectory()
		if baseErr != nil {
			return baseErr
		}
		_ = os.Remove(filepath.Join(base, CurrentFileName))
	}
	return nil
}

func ListServices(contextName string) ([]string, error) {
	dir, err := ContextDirectory(contextName)
	if err != nil {
		return nil, err
	}
	entries, err := os.ReadDir(filepath.Join(dir, ServicesDirName))
	if errors.Is(err, os.ErrNotExist) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	names := make([]string, 0)
	for _, entry := range entries {
		if entry.IsDir() {
			names = append(names, entry.Name())
		}
	}
	sort.Strings(names)
	return names, nil
}
func HasCredentials(contextName, service string) (bool, error) {
	path, err := CredentialsPath(contextName, service)
	if err != nil {
		return false, err
	}
	_, err = os.Stat(path)
	if errors.Is(err, os.ErrNotExist) {
		return false, nil
	}
	return err == nil, err
}

func writeYAML(path string, value any) error {
	data, err := yaml.Marshal(value)
	if err != nil {
		return err
	}
	if err := os.MkdirAll(filepath.Dir(path), 0o700); err != nil {
		return err
	}
	return os.WriteFile(path, data, 0o600)
}
func readYAML(path string, value any) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("read config: %w", err)
	}
	decoder := yaml.NewDecoder(strings.NewReader(string(data)))
	decoder.KnownFields(true)
	if err := decoder.Decode(value); err != nil {
		return fmt.Errorf("parse config: %w", err)
	}
	return nil
}

package matreshka

import (
	"os"
	"path"
	"sort"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"go.vervstack.ru/matreshka/pkg/matreshka/environment"
	"go.vervstack.ru/matreshka/pkg/matreshka/resources"
	"go.vervstack.ru/matreshka/pkg/matreshka/server"
)

func Test_ParseConfig(t *testing.T) {
	tmpDirPath := path.Join(os.TempDir(), t.Name())
	require.NoError(t, os.MkdirAll(tmpDirPath, os.ModePerm))

	t.Run("OK_EMPTY", func(t *testing.T) {
		cfgPath := path.Join(tmpDirPath, path.Base(t.Name())+".yaml")
		defer func() {
			require.NoError(t, os.RemoveAll(cfgPath))
		}()

		require.NoError(t,
			os.WriteFile(
				cfgPath,
				emptyConfig,
				os.ModePerm))

		cfg, err := getFromFile(cfgPath)
		require.NoError(t, err)
		require.Equal(t, cfg, NewEmptyConfig())
	})

	t.Run("OK_FULL_FROM_FILE", func(t *testing.T) {
		t.Parallel()

		cfgActual, err := ParseConfig(fullConfig)
		require.NoError(t, err)

		cfgExpect := getFullConfigTest()
		for _, s := range cfgExpect.Servers {
			s.Name = ""
		}

		require.Equal(t, cfgExpect.AppInfo, cfgActual.AppInfo)
		require.Equal(t, cfgExpect.Environment, cfgActual.Environment)
		require.Equal(t, cfgExpect.DataSources, cfgActual.DataSources)
		require.Equal(t, cfgExpect.ServiceDiscovery, cfgActual.ServiceDiscovery)
		require.Equal(t, cfgExpect.Servers, cfgActual.Servers)
	})

	t.Run("ERROR_UNKNOWN_FILE)", func(t *testing.T) {
		cfg, err := getFromFile("unreadable config path")
		require.ErrorIs(t, err, os.ErrNotExist)
		require.Equal(t, cfg, NewEmptyConfig())
	})

	t.Run("ERROR_UNMARSHALLING_CONFIG", func(t *testing.T) {
		cfgPath := path.Join(tmpDirPath, path.Base(t.Name())+".yaml")
		defer func() {
			require.NoError(t, os.RemoveAll(cfgPath))
		}()
		require.NoError(t,
			os.WriteFile(
				cfgPath,
				[]byte("1f!cked #p\nc0nfig"),
				os.ModePerm))

		cfg, err := getFromFile(cfgPath)
		require.Equal(t, err.Error(), "yaml: unmarshal errors:\n  line 1: cannot unmarshal !!str `1f!cked` into matreshka.AppConfig;error decoding config to struct")
		require.Equal(t, cfg, NewEmptyConfig())
	})
}

func Test_ReadConfigs(t *testing.T) {
	t.Parallel()

	tmpDirPath := path.Join(os.TempDir(), t.Name())
	require.NoError(t, os.MkdirAll(tmpDirPath, os.ModePerm))

	t.Run("OK", func(t *testing.T) {
		// preparing empty config
		emptyConfigPath := path.Join(tmpDirPath, path.Base(t.Name()+"_empty")+".yaml")
		{
			defer func() {
				require.NoError(t, os.RemoveAll(emptyConfigPath))
			}()
			require.NoError(t,
				os.WriteFile(
					emptyConfigPath,
					emptyConfig,
					os.ModePerm))
		}

		fullConfigPath := path.Join(tmpDirPath, path.Base(t.Name()+"_full")+".yaml")
		{
			defer func() {
				require.NoError(t, os.RemoveAll(fullConfigPath))
			}()
			require.NoError(t,
				os.WriteFile(
					fullConfigPath,
					fullConfig,
					os.ModePerm))
		}

		expectedCfg := AppConfig{
			AppInfo: AppInfo{
				Name:            "matreshka",
				Version:         "v0.0.1",
				StartupDuration: time.Second * 10,
			},
			DataSources: []resources.Resource{
				getPostgresClientTest(),
				getRedisClientTest(),
				getTelegramClientTest(),
				getGRPCClientTest(),
			},
			Servers:          getConfigServersFull(),
			Environment:      Environment(getEnvironmentVariables()),
			ServiceDiscovery: getConfigServiceDiscovery(),
		}

		sort.Slice(expectedCfg.Environment, func(i, j int) bool {
			return expectedCfg.Environment[i].Name > expectedCfg.Environment[j].Name
		})

		t.Run("EMPTY_TO_FULL", func(t *testing.T) {
			// empty and full config merge
			actualCfg, err := ReadConfigs(emptyConfigPath, fullConfigPath)

			sort.Slice(actualCfg.Environment, func(i, j int) bool {
				return actualCfg.Environment[i].Name > actualCfg.Environment[j].Name
			})

			require.NoError(t, err)
			require.Equal(t, expectedCfg, actualCfg)
		})

		t.Run("FULL_TO_EMPTY", func(t *testing.T) {
			actualCfg, err := ReadConfigs(fullConfigPath, emptyConfigPath)

			sort.Slice(actualCfg.Environment, func(i, j int) bool {
				return actualCfg.Environment[i].Name > actualCfg.Environment[j].Name
			})

			require.NoError(t, err)
			require.Equal(t, expectedCfg, actualCfg)
		})

		t.Run("EMPTY_TO_FULL_TO_ENV", func(t *testing.T) {
			// empty and full config merge

			expectedCfgWithEnv := AppConfig{
				AppInfo: AppInfo{
					Name:            "matreshka",
					Version:         "v0.0.1",
					StartupDuration: time.Second * 10,
				},
				DataSources: []resources.Resource{
					getPostgresClientTest(),
					getRedisClientTest(),
					getTelegramClientTest(),
					getGRPCClientTest(),
				},
				Servers:          getConfigServersFull(),
				Environment:      Environment(getEnvironmentVariables()),
				ServiceDiscovery: getConfigServiceDiscovery(),
			}

			sort.Slice(expectedCfgWithEnv.Environment, func(i, j int) bool {
				return expectedCfgWithEnv.Environment[i].Name > expectedCfgWithEnv.Environment[j].Name
			})

			// Evon style only
			require.NoError(t, os.Setenv("database-max-connections", "3"))
			expectedCfgWithEnv.Environment[5] = environment.MustNewVariable("database_max_connections", 3)

			// Service name + env part + evon style
			require.NoError(t, os.Setenv("matreshka_environment_welcome-string", "wel-cum"))
			expectedCfgWithEnv.Environment[0] = environment.MustNewVariable("welcome_string", "wel-cum")

			// Service name + evon style
			require.NoError(t, os.Setenv("matreshka_one-of-welcome-string", "three"))
			expectedCfgWithEnv.Environment[4] = environment.MustNewVariable(
				"one_of_welcome_string", "three",
				environment.WithEnum("one", "two", "three"))

			// Service name + env part + default env style
			require.NoError(t, os.Setenv("matreshka_environment_request_timeout", "10s"))
			expectedCfgWithEnv.Environment[3] = environment.MustNewVariable("request_timeout", time.Second*10)

			// Default env style
			require.NoError(t, os.Setenv("available_ports", "[12:18,20]"))
			expectedCfgWithEnv.Environment[8] = environment.MustNewVariable("available_ports", []int{12, 13, 14, 15, 16, 17, 18, 20})

			actualCfg, err := ReadConfigs(emptyConfigPath, fullConfigPath)
			require.NoError(t, err)

			sort.Slice(actualCfg.Environment, func(i, j int) bool {
				return actualCfg.Environment[i].Name > actualCfg.Environment[j].Name
			})

			require.NoError(t, err)
			require.Equal(t, expectedCfgWithEnv, actualCfg)
		})
	})

	t.Run("INVALID_READING_ONE_CONFIG", func(t *testing.T) {
		cfgPath := path.Join(tmpDirPath, path.Base(t.Name())+".yaml")

		defer func() {
			require.NoError(t, os.RemoveAll(cfgPath))
		}()

		require.NoError(t,
			os.WriteFile(
				cfgPath,
				emptyConfig,
				os.ModePerm))

		cfg, err := ReadConfigs(cfgPath, "unreadable config path")
		require.ErrorIs(t, err, os.ErrNotExist)
		require.Equal(t, cfg, NewEmptyConfig())
	})

	t.Run("INVALID_READING_FIRST_CONFIG", func(t *testing.T) {
		cfg, err := ReadConfigs("unreadable config path")
		require.ErrorIs(t, err, os.ErrNotExist)
		require.Equal(t, cfg, NewEmptyConfig())
	})
}

func Test_ReadConfig(t *testing.T) {
	// Not t.Parallel(): the WITH_CONFIG_BYTES_ONLY subtest uses t.Setenv,
	// which panics if any ancestor test is marked parallel.

	tmpDirPath := path.Join(os.TempDir(), t.Name())
	require.NoError(t, os.MkdirAll(tmpDirPath, os.ModePerm))

	t.Run("WITH_CONFIG_PATHS_ONLY", func(t *testing.T) {
		emptyConfigPath := path.Join(tmpDirPath, path.Base(t.Name()+"_empty")+".yaml")
		defer func() {
			require.NoError(t, os.RemoveAll(emptyConfigPath))
		}()
		require.NoError(t,
			os.WriteFile(
				emptyConfigPath,
				emptyConfig,
				os.ModePerm))

		fullConfigPath := path.Join(tmpDirPath, path.Base(t.Name()+"_full")+".yaml")
		defer func() {
			require.NoError(t, os.RemoveAll(fullConfigPath))
		}()
		require.NoError(t,
			os.WriteFile(
				fullConfigPath,
				fullConfig,
				os.ModePerm))

		expectedCfg, err := ReadConfigs(emptyConfigPath, fullConfigPath)
		require.NoError(t, err)

		actualCfg, err := ReadConfig(WithConfigPaths(emptyConfigPath, fullConfigPath))
		require.NoError(t, err)

		require.Equal(t, expectedCfg, actualCfg)
	})

	t.Run("WITH_CONFIG_BYTES_ONLY", func(t *testing.T) {
		bytesCfg := []byte("app_info:\n  name: test\ndata_sources:\n  - resource_name: postgres\n    host: localhost\n")

		actualCfg, err := ReadConfig(WithConfigBytes(bytesCfg))
		require.NoError(t, err)

		pg, ok := actualCfg.DataSources.get("postgres").(*resources.Postgres)
		require.True(t, ok)
		require.Equal(t, "localhost", pg.Host)

		t.Setenv("DATA-SOURCES_POSTGRES_HOST", "envhost")

		actualCfgWithEnv, err := ReadConfig(WithConfigBytes(bytesCfg))
		require.NoError(t, err)

		pgWithEnv, ok := actualCfgWithEnv.DataSources.get("postgres").(*resources.Postgres)
		require.True(t, ok)
		require.Equal(t, "envhost", pgWithEnv.Host)
	})

	t.Run("WITH_CONFIG_BYTES_ONLY_FULL_UNDERSCORE_ENV", func(t *testing.T) {
		bytesCfg := []byte("app_info:\n  name: test\ndata_sources:\n  - resource_name: postgres\n    host: localhost\n")

		// Fully underscored variant of DATA-SOURCES_POSTGRES_HOST: real env var
		// names can't contain '-', so this must resolve the same way.
		t.Setenv("DATA_SOURCES_POSTGRES_HOST", "envhost")

		actualCfgWithEnv, err := ReadConfig(WithConfigBytes(bytesCfg))
		require.NoError(t, err)

		pgWithEnv, ok := actualCfgWithEnv.DataSources.get("postgres").(*resources.Postgres)
		require.True(t, ok)
		require.Equal(t, "envhost", pgWithEnv.Host)
	})

	t.Run("WITH_CONFIG_PATHS_AND_BYTES", func(t *testing.T) {
		pathConfigPath := path.Join(tmpDirPath, path.Base(t.Name())+".yaml")
		defer func() {
			require.NoError(t, os.RemoveAll(pathConfigPath))
		}()
		require.NoError(t,
			os.WriteFile(
				pathConfigPath,
				[]byte("app_info:\n  name: test\ndata_sources:\n  - resource_name: postgres\n    host: path-host\n"),
				os.ModePerm))

		bytesCfg := []byte("app_info:\n  name: test\n  version: v0.0.1\ndata_sources:\n  - resource_name: postgres\n    host: bytes-host\n  - resource_name: redis\n    host: bytes-redis-host\n")

		actualCfg, err := ReadConfig(WithConfigPaths(pathConfigPath), WithConfigBytes(bytesCfg))
		require.NoError(t, err)

		// path source registered first -> wins for a name present in both
		pg, ok := actualCfg.DataSources.get("postgres").(*resources.Postgres)
		require.True(t, ok)
		require.Equal(t, "path-host", pg.Host)

		// bytes source fills in a name the path source did not declare
		redis, ok := actualCfg.DataSources.get("redis").(*resources.Redis)
		require.True(t, ok)
		require.Equal(t, "bytes-redis-host", redis.Host)

		// path source's app_info wins for Name; bytes-only field (Version) fills in
		require.Equal(t, "test", actualCfg.AppInfo.Name)
		require.Equal(t, "v0.0.1", actualCfg.AppInfo.Version)
	})

	t.Run("READ_CONFIGS_WRAPPER_REGRESSION", func(t *testing.T) {
		fullConfigPath := path.Join(tmpDirPath, path.Base(t.Name())+".yaml")
		defer func() {
			require.NoError(t, os.RemoveAll(fullConfigPath))
		}()
		require.NoError(t,
			os.WriteFile(
				fullConfigPath,
				fullConfig,
				os.ModePerm))

		expectedCfg := getFullConfigTest()

		actualCfg, err := ReadConfigs(fullConfigPath)
		require.NoError(t, err)

		sort.Slice(expectedCfg.Environment, func(i, j int) bool {
			return expectedCfg.Environment[i].Name < expectedCfg.Environment[j].Name
		})
		sort.Slice(actualCfg.Environment, func(i, j int) bool {
			return actualCfg.Environment[i].Name < actualCfg.Environment[j].Name
		})

		require.Equal(t, expectedCfg.AppInfo, actualCfg.AppInfo)
		require.Equal(t, expectedCfg.DataSources, actualCfg.DataSources)
		require.Equal(t, expectedCfg.Servers, actualCfg.Servers)
		require.Equal(t, expectedCfg.Environment, actualCfg.Environment)
		require.Equal(t, expectedCfg.ServiceDiscovery, actualCfg.ServiceDiscovery)
	})

	t.Run("ERROR_NONEXISTENT_FILE", func(t *testing.T) {
		cfg, err := ReadConfig(WithConfigPaths("unreadable config path"))
		require.Error(t, err)
		require.ErrorIs(t, err, os.ErrNotExist)
		require.Equal(t, NewEmptyConfig(), cfg)
	})

	t.Run("WITH_ENV_FILE_OVERRIDES_YAML", func(t *testing.T) {
		fullConfigPath := path.Join(tmpDirPath, path.Base(t.Name())+".yaml")
		defer func() {
			require.NoError(t, os.RemoveAll(fullConfigPath))
		}()
		require.NoError(t,
			os.WriteFile(
				fullConfigPath,
				fullConfig,
				os.ModePerm))

		actualCfg, err := ReadConfig(
			WithConfigPaths(fullConfigPath),
			WithEnvFile(path.Join("tests", "sample.env")),
		)
		require.NoError(t, err)

		// .env overrides a field the YAML fixture already set
		require.Equal(t, "v9.9.9", actualCfg.AppInfo.Version)

		pg, ok := actualCfg.DataSources.get("postgres").(*resources.Postgres)
		require.True(t, ok)
		require.Equal(t, "envfile-host", pg.Host)

		// .env fills in a field the YAML fixture did not set at all
		require.Equal(t, "envfile:1281", actualCfg.ServiceDiscovery.MakoshUrl)

		// fields untouched by .env still come from YAML
		require.Equal(t, "matreshka", actualCfg.AppInfo.Name)
		require.Equal(t, uint64(5432), pg.Port)
	})

	t.Run("REAL_ENV_OVERRIDES_ENV_FILE", func(t *testing.T) {
		envFilePath := path.Join(tmpDirPath, path.Base(t.Name())+".env")
		defer func() {
			require.NoError(t, os.RemoveAll(envFilePath))
		}()
		require.NoError(t,
			os.WriteFile(
				envFilePath,
				[]byte("DATA-SOURCES_POSTGRES_HOST=envfilehost\n"),
				os.ModePerm))

		bytesCfg := []byte("app_info:\n  name: test\ndata_sources:\n  - resource_name: postgres\n    host: localhost\n")

		actualCfg, err := ReadConfig(WithConfigBytes(bytesCfg), WithEnvFile(envFilePath))
		require.NoError(t, err)

		pg, ok := actualCfg.DataSources.get("postgres").(*resources.Postgres)
		require.True(t, ok)
		require.Equal(t, "envfilehost", pg.Host)

		// Real OS env var (full-underscore style, mirrors
		// WITH_CONFIG_BYTES_ONLY_FULL_UNDERSCORE_ENV) must win over the .env file value.
		t.Setenv("DATA_SOURCES_POSTGRES_HOST", "realenvhost")

		actualCfgWithRealEnv, err := ReadConfig(WithConfigBytes(bytesCfg), WithEnvFile(envFilePath))
		require.NoError(t, err)

		pgWithRealEnv, ok := actualCfgWithRealEnv.DataSources.get("postgres").(*resources.Postgres)
		require.True(t, ok)
		require.Equal(t, "realenvhost", pgWithRealEnv.Host)
	})

	t.Run("ERROR_MISSING_ENV_FILE", func(t *testing.T) {
		bytesCfg := []byte("app_info:\n  name: test\n")

		cfg, err := ReadConfig(WithConfigBytes(bytesCfg), WithEnvFile("unreadable env file path"))
		require.Error(t, err)
		require.ErrorIs(t, err, os.ErrNotExist)
		// the YAML/bytes sources still get merged before the env-file error is returned
		require.Equal(t, "test", cfg.AppInfo.Name)
	})
}

func Test_MergeConfigs_ServerNameDedup(t *testing.T) {
	t.Parallel()

	master := NewEmptyConfig()
	masterServer := &server.Server{
		Name: "MASTER",
		Port: "8080",
	}
	master.Servers[8080] = masterServer

	slave := NewEmptyConfig()
	slaveServer := &server.Server{
		Name: "",
		Port: "9090",
	}
	slave.Servers[9090] = slaveServer

	merged := MergeConfigs(master, slave)

	require.Len(t, merged.Servers, 1)
	require.Equal(t, masterServer, merged.Servers[8080])
}

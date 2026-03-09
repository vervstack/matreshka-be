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

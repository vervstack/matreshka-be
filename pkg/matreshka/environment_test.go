package matreshka

import (
	"sort"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"go.redsock.ru/evon"

	"go.vervstack.ru/matreshka/pkg/matreshka/environment"
	testCfg "go.vervstack.ru/matreshka/pkg/matreshka/tests"
)

func Test_Environment(t *testing.T) {
	t.Parallel()

	t.Run("PARSE_ENV_TO_STRUCT", func(t *testing.T) {
		t.Parallel()

		env := Environment(getEnvironmentVariables())

		customEnvConf := &testCfg.EnvironmentConfig{}

		err := env.ParseToStruct(customEnvConf)
		require.NoError(t, err)

		expected := &testCfg.EnvironmentConfig{
			AvailablePorts:                   []int{10, 12, 34, 35, 36, 37, 38, 39, 40},
			CreditPercent:                    0.01,
			CreditPercentsBasedOnYearOfBirth: []float64{0.01, 0.02, 0.03, 0.04},
			DatabaseMaxConnections:           1,
			OneOfWelcomeString:               "one",
			RequestTimeout:                   time.Second * 10,
			TrueFalser:                       true,
			UsernamesToBan:                   []string{"hacker228", "mothe4acker"},
			WelcomeString:                    "not so basic 🤡 string",
		}
		require.Equal(t, expected, customEnvConf)
	})

	t.Run("PARSE_ENV_MORE_THAN_HAVE_IN_STRUCT", func(t *testing.T) {
		t.Parallel()

		env := Environment([]*environment.Variable{
			environment.MustNewVariable("new_unknown", "nil"),
		})

		customEnvConf := &testCfg.EnvironmentConfig{}

		err := env.ParseToStruct(customEnvConf)
		require.ErrorIs(t, err, ErrNotFound)

		expected := &testCfg.EnvironmentConfig{}
		require.Equal(t, expected, customEnvConf)
	})

	t.Run("MARSHAL_UNMARSHAL_YAML", func(t *testing.T) {
		t.Run("OK", func(t *testing.T) {
			ac := AppConfig{
				Environment: getEnvironmentVariables(),
			}

			bytes, err := ac.Marshal()
			require.NoError(t, err)

			var ac2 AppConfig
			err = ac2.Unmarshal(bytes)
			require.NoError(t, err)

			require.Equal(t, ac.Environment, ac2.Environment)
		})

		t.Run("ERROR", func(t *testing.T) {
			t.Run("UNKNOWN_TYPE", func(t *testing.T) {
				bytes := []byte(`
environment:
    - name: database_max_connections
      type: integer
      value: 1
`)

				var ac2 AppConfig
				err := ac2.Unmarshal(bytes)
				require.ErrorIs(t, err, environment.ErrUnknownEnvVariableType)
			})
		})
	})

	t.Run("MARSHAL_UNMARSHAL_EVON", func(t *testing.T) {
		env := Environment(getEnvironmentVariables())
		sort.Slice(env, func(i, j int) bool {
			return env[i].Name < env[j].Name
		})

		nodes, err := env.MarshalEnv("MATRESHKA")
		require.NoError(t, err)

		root := &evon.Node{
			Name:       "MATRESHKA",
			InnerNodes: nodes,
		}

		var env2 Environment
		err = env2.UnmarshalEnv(root)
		require.NoError(t, err)

		require.Equal(t, env, env2)
	})
}

package matreshka

import (
	"testing"

	"github.com/stretchr/testify/require"
	"go.redsock.ru/evon"
	"gopkg.in/yaml.v3"

	"go.vervstack.ru/matreshka/pkg/matreshka/server"
)

func Test_Servers(t *testing.T) {
	t.Run("YAML", func(t *testing.T) {
		t.Run("Marshal_OK", func(t *testing.T) {
			t.Parallel()

			var cfgIn AppConfig
			cfgIn.Servers = getConfigServersFull()

			marshaled, err := cfgIn.Marshal()
			require.NoError(t, err)

			var actual map[any]any

			require.NoError(t, yaml.Unmarshal(marshaled, &actual))

			expected := map[any]any{
				"servers": map[any]any{
					8080: map[string]any{
						"/{FS}": map[string]any{
							"dist": "web/dist",
						},
						"name": "MASTER",
					},
					50051: map[string]any{
						"/{GRPC}": map[string]any{
							"module":  "pkg/matreshka_be_api",
							"gateway": "/api",
						},
						"name": "MASTER2",
					},
				},
			}

			require.Equal(t, expected, actual)
		})
		t.Run("Marshal_With_Name_OK", func(t *testing.T) {
			t.Parallel()

			var cfgIn AppConfig
			cfgIn.Servers = getConfigServersFull()
			cfgIn.Servers[8080].Name = "Main"
			cfgIn.Servers[50051].Name = "Grpc"

			marshaled, err := cfgIn.Marshal()
			require.NoError(t, err)

			var actual map[any]any

			require.NoError(t, yaml.Unmarshal(marshaled, &actual))

			expected := map[any]any{
				"servers": map[any]any{
					8080: map[string]any{
						"name": "Main",
						"/{FS}": map[string]any{
							"dist": "web/dist",
						},
					},
					50051: map[string]any{
						"name": "Grpc",
						"/{GRPC}": map[string]any{
							"module":  "pkg/matreshka_be_api",
							"gateway": "/api",
						},
					},
				},
			}

			require.Equal(t, expected, actual)
		})
		t.Run("Unmarshal_OK", func(t *testing.T) {
			t.Parallel()

			cfg, err := ParseConfig(apiConfig)
			require.NoError(t, err)

			servers := getConfigServersFull()
			require.Equal(t, cfg.Servers, servers)
		})
		t.Run("Unmarshal_With_Name_OK", func(t *testing.T) {
			t.Parallel()

			cfg, err := ParseConfig(apiConfigWithName)
			require.NoError(t, err)

			servers := getConfigServersFull()
			servers[8080].Name = "Main"
			require.Equal(t, cfg.Servers, servers)
		})

		t.Run("Unmarshal_Error", func(t *testing.T) {
			t.Run("Port_Is_Not_Int", func(t *testing.T) {
				cfg := AppConfig{}
				err := cfg.Unmarshal(apiInvalidPortConfig)
				require.Equal(t, err.Error(), "strconv.Atoi: parsing \"string\": invalid syntax;error converting port to int")
			})
			t.Run("Invalid_Struct", func(t *testing.T) {
				cfg := AppConfig{}
				err := cfg.Unmarshal(apiInvalidStructConfig)
				require.Equal(t, err.Error(), "yaml: unmarshal errors:\n  line 2: cannot unmarshal !!seq into map[string]yaml.Node;error unmarshalling to yaml.Nodes")
			})
			t.Run("Invalid_Item", func(t *testing.T) {
				cfg := AppConfig{}
				err := cfg.Unmarshal(apiInvalidItemConfig)
				require.Equal(t, err.Error(), "yaml: unmarshal errors:\n  line 3: cannot unmarshal !!seq into map[string]yaml.Node;error unmarshalling YAML;error decoding server")
			})
		})
	})

	t.Run("ParseToStruct", func(t *testing.T) {
		t.Run("Duplicate_Servers_Error", func(t *testing.T) {
			t.Parallel()

			type dst struct {
				MASTER *server.Server
			}

			for i := 0; i < 20; i++ {
				servers := Servers{
					8080: {Name: "", Port: "8080"},
					9090: {Name: "MASTER", Port: "9090"},
				}

				var d dst
				err := servers.ParseToStruct(&d)
				require.Error(t, err)
				require.Contains(t, err.Error(), "duplicate")
			}
		})

		t.Run("Distinct_Servers_OK", func(t *testing.T) {
			t.Parallel()

			type dst struct {
				MASTER  *server.Server
				MASTER2 *server.Server
			}

			servers := getConfigServersFull()

			var d dst
			err := servers.ParseToStruct(&d)
			require.NoError(t, err)
			require.Equal(t, servers[8080], d.MASTER)
			require.Equal(t, servers[50051], d.MASTER2)
		})
	})

	t.Run("ENV", func(t *testing.T) {
		t.Run("Marshal", func(t *testing.T) {
			t.Parallel()

			var cfgIn AppConfig
			cfgIn.Servers = getConfigServersFull()

			marshaledNodes, err := cfgIn.Servers.MarshalEnv("MATRESHKA_SERVERS")
			require.NoError(t, err)

			marshalledBytes := evon.Marshal(marshaledNodes)
			require.Equal(t, string(apiEnvConfig), string(marshalledBytes))
		})
		t.Run("Unmarshal", func(t *testing.T) {
			t.Parallel()
			cfg := NewEmptyConfig()
			err := evon.UnmarshalWithPrefix("MATRESHKA", apiEnvConfig, &cfg)
			require.NoError(t, err)

			servers := getConfigServersFull()
			require.Equal(t, cfg.Servers, servers)
		})
	})

}

func getConfigServersFull() Servers {
	return Servers{
		8080: {
			Name: "MASTER",
			Port: "8080",
			GRPC: map[string]*server.GRPC{},
			FS: map[string]*server.FS{
				"/{FS}": {
					Dist: "web/dist",
				},
			},
			HTTP: map[string]*server.HTTP{},
		},
		50051: {
			Name: "MASTER2",
			Port: "50051",
			GRPC: map[string]*server.GRPC{
				"/{GRPC}": {
					Module:  "pkg/matreshka_be_api",
					Gateway: "/api",
				},
			},
			FS:   map[string]*server.FS{},
			HTTP: map[string]*server.HTTP{},
		},
	}
}

package web_test

import (
	"bytes"
	"context"
	"encoding/json"
	"mtui/types"
	"mtui/web"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/stretchr/testify/assert"
)

type bearerTransport struct {
	token string
}

func (t *bearerTransport) RoundTrip(r *http.Request) (*http.Response, error) {
	r = r.Clone(r.Context())
	r.Header.Set("Authorization", "Bearer "+t.token)
	return http.DefaultTransport.RoundTrip(r)
}

func loginToken(t *testing.T, api *web.Api) string {
	buf, err := json.Marshal(&web.LoginRequest{Username: "singleplayer", Password: "mypass"})
	assert.NoError(t, err)
	w := httptest.NewRecorder()
	api.DoLogin(w, httptest.NewRequest("POST", "http://", bytes.NewBuffer(buf)))
	assert.Equal(t, 200, w.Result().StatusCode)

	for _, c := range w.Result().Cookies() {
		if c.Name == web.TOKEN_COOKIE_NAME {
			return c.Value
		}
	}
	t.Fatal("no token cookie")
	return ""
}

func TestMCP(t *testing.T) {
	api, app := CreateTestApi(t)
	token := loginToken(t, api)

	srv := httptest.NewServer(api.MCPHandler())
	defer srv.Close()

	// api feature disabled: endpoint not available
	resp, err := http.Post(srv.URL, "application/json", strings.NewReader("{}"))
	assert.NoError(t, err)
	assert.Equal(t, 500, resp.StatusCode)

	assert.NoError(t, app.Repos.FeatureRepository.Set(&types.Feature{Name: types.FEATURE_API, Enabled: true}))

	// no token
	resp, err = http.Post(srv.URL, "application/json", strings.NewReader("{}"))
	assert.NoError(t, err)
	assert.Equal(t, 401, resp.StatusCode)

	// invalid token
	ctx := context.Background()
	client := mcp.NewClient(&mcp.Implementation{Name: "test", Version: "1"}, nil)
	_, err = client.Connect(ctx, &mcp.StreamableClientTransport{
		Endpoint:   srv.URL,
		HTTPClient: &http.Client{Transport: &bearerTransport{token: "invalid"}},
	}, nil)
	assert.Error(t, err)

	// valid token, "interact" priv only
	session, err := client.Connect(ctx, &mcp.StreamableClientTransport{
		Endpoint:   srv.URL,
		HTTPClient: &http.Client{Transport: &bearerTransport{token: token}},
	}, nil)
	assert.NoError(t, err)
	defer session.Close()

	tools, err := session.ListTools(ctx, nil)
	assert.NoError(t, err)
	names := map[string]bool{}
	for _, tool := range tools.Tools {
		names[tool.Name] = true
		assert.True(t, tool.Annotations != nil, tool.Name)
	}
	assert.True(t, names["get_server_stats"])
	assert.True(t, names["search_players"])
	assert.True(t, names["get_player_info"])
	assert.False(t, names["search_logs"], "requires the ban priv")
	assert.False(t, names["list_mods"], "requires the server priv")
	assert.False(t, names["list_areas"], "requires the areas feature")

	// player info
	res, err := session.CallTool(ctx, &mcp.CallToolParams{
		Name:      "get_player_info",
		Arguments: map[string]any{"name": "singleplayer"},
	})
	assert.NoError(t, err)
	assert.False(t, res.IsError)
	info := &web.PlayerInfo{}
	assert.NoError(t, json.Unmarshal([]byte(res.Content[0].(*mcp.TextContent).Text), info))
	assert.Equal(t, "singleplayer", info.Name)
	assert.Equal(t, []string{"interact"}, info.Privs)

	// unknown player: tool error, not a protocol error
	res, err = session.CallTool(ctx, &mcp.CallToolParams{
		Name:      "get_player_info",
		Arguments: map[string]any{"name": "nobody"},
	})
	assert.NoError(t, err)
	assert.True(t, res.IsError)

	// search
	res, err = session.CallTool(ctx, &mcp.CallToolParams{
		Name:      "search_players",
		Arguments: map[string]any{"name_like": "single%"},
	})
	assert.NoError(t, err)
	assert.False(t, res.IsError)
	assert.Contains(t, res.Content[0].(*mcp.TextContent).Text, `"count":1`)
}

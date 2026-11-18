package web_test

import (
	"bytes"
	"context"
	"encoding/json"
	"mtui/types"
	"mtui/web"
	"net/http"
	"net/http/httptest"
	"os"
	"path"
	"strings"
	"testing"

	authdb "github.com/luanti-golang/mtdb/auth"
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

func connectMCP(t *testing.T, url, token string) *mcp.ClientSession {
	client := mcp.NewClient(&mcp.Implementation{Name: "test", Version: "1"}, nil)
	session, err := client.Connect(context.Background(), &mcp.StreamableClientTransport{
		Endpoint:   url,
		HTTPClient: &http.Client{Transport: &bearerTransport{token: token}},
	}, nil)
	assert.NoError(t, err)
	return session
}

func callMCP(t *testing.T, session *mcp.ClientSession, name string, args map[string]any) (string, bool) {
	res, err := session.CallTool(context.Background(), &mcp.CallToolParams{Name: name, Arguments: args})
	assert.NoError(t, err)
	return res.Content[0].(*mcp.TextContent).Text, res.IsError
}

func TestMCPWriteTools(t *testing.T) {
	api, app := CreateTestApi(t)

	// grant all relevant privs and features
	auth_entry, err := app.DBContext.Auth.GetByUsername("singleplayer")
	assert.NoError(t, err)
	for _, priv := range []string{"shout", "ban", "server"} {
		assert.NoError(t, app.DBContext.Privs.Create(&authdb.PrivilegeEntry{ID: *auth_entry.ID, Privilege: priv}))
	}
	for _, f := range []types.FeatureName{types.FEATURE_API, types.FEATURE_MAIL, types.FEATURE_XBAN, types.FEATURE_MODMANAGEMENT,
		types.FEATURE_LUANTI_CONFIG, types.FEATURE_SHELL, "chat"} {
		assert.NoError(t, app.Repos.FeatureRepository.Set(&types.Feature{Name: f, Enabled: true}))
	}
	app.Config.LuantiConfig = path.Join(app.WorldDir, "minetest.conf")
	assert.NoError(t, os.WriteFile(app.Config.LuantiConfig, []byte("max_users = 10\n"), 0644))

	srv := httptest.NewServer(api.MCPHandler())
	defer srv.Close()
	session := connectMCP(t, srv.URL, loginToken(t, api))
	defer session.Close()

	// all tools listed, write tools are not read-only
	tools, err := session.ListTools(context.Background(), nil)
	assert.NoError(t, err)
	byName := map[string]*mcp.Tool{}
	for _, tool := range tools.Tools {
		byName[tool.Name] = tool
	}
	for _, name := range []string{"send_chat", "ban_player", "unban_player", "list_inbox", "send_mail", "search_contentdb",
		"resolve_contentdb_dependencies", "install_contentdb_package", "install_git_mod", "update_mod_version",
		"set_luanti_setting", "delete_luanti_setting", "execute_chatcommand"} {
		tool := byName[name]
		if assert.NotNil(t, tool, name) {
			assert.NotNil(t, tool.Annotations, name)
		}
	}
	assert.False(t, byName["ban_player"].Annotations.ReadOnlyHint)
	assert.True(t, *byName["ban_player"].Annotations.DestructiveHint)
	assert.True(t, byName["list_inbox"].Annotations.ReadOnlyHint)

	// mail roundtrip
	_, isErr := callMCP(t, session, "send_mail", map[string]any{"to": "singleplayer", "subject": "hi", "body": "hello there"})
	assert.False(t, isErr)
	text, isErr := callMCP(t, session, "list_inbox", map[string]any{"unread_only": true})
	assert.False(t, isErr)
	assert.Contains(t, text, "hello there")
	assert.Contains(t, text, `"total":1`)
	text, isErr = callMCP(t, session, "send_mail", map[string]any{"to": "nobody", "subject": "x", "body": "y"})
	assert.True(t, isErr)
	assert.Contains(t, text, "invalid recipient")

	// cc/bcc are delivered once to every unique recipient, unknown ones are all reported
	for _, name := range []string{"alice", "bob", "carol"} {
		assert.NoError(t, app.DBContext.Auth.Create(&authdb.AuthEntry{Name: name, Password: "x"}))
	}
	text, isErr = callMCP(t, session, "send_mail", map[string]any{"to": "alice", "cc": "bob, ghost1", "bcc": "ghost2", "subject": "x", "body": "y"})
	assert.True(t, isErr)
	assert.Contains(t, text, "'ghost1', 'ghost2'")

	text, isErr = callMCP(t, session, "send_mail", map[string]any{"to": "alice,bob", "cc": "bob carol", "bcc": "singleplayer", "subject": "multi", "body": "to all"})
	assert.False(t, isErr)
	assert.Contains(t, text, `"to":"alice, bob"`)
	assert.Contains(t, text, `"cc":"bob, carol"`)
	for _, name := range []string{"alice", "bob", "carol", "singleplayer"} {
		e, err := app.Mail.GetEntry(name)
		assert.NoError(t, err)
		count := 0
		for _, m := range e.Inbox {
			if m.Subject == "multi" {
				count++
			}
		}
		assert.Equal(t, 1, count, name)
	}
	e, err := app.Mail.GetEntry("singleplayer")
	assert.NoError(t, err)
	assert.Equal(t, 2, len(e.Outbox), "one outbox entry per sent mail ('hi' and 'multi')")

	// input validation before anything reaches the engine
	text, isErr = callMCP(t, session, "ban_player", map[string]any{"name": "bob", "reason": "grief", "duration": "forever"})
	assert.True(t, isErr)
	assert.Contains(t, text, "invalid duration")
	_, isErr = callMCP(t, session, "ban_player", map[string]any{"name": "bob", "reason": ""})
	assert.True(t, isErr)
	_, isErr = callMCP(t, session, "execute_chatcommand", map[string]any{"command": " / "})
	assert.True(t, isErr)
	_, isErr = callMCP(t, session, "send_chat", map[string]any{"message": "  "})
	assert.True(t, isErr)

	// settings
	_, isErr = callMCP(t, session, "set_luanti_setting", map[string]any{"key": "my_custom", "value": "a \" quote"})
	assert.False(t, isErr)
	text, isErr = callMCP(t, session, "get_luanti_settings", map[string]any{"filter": "my_custom"})
	assert.False(t, isErr)
	assert.Contains(t, text, `a \" quote`)
	_, isErr = callMCP(t, session, "delete_luanti_setting", map[string]any{"key": "my_custom"})
	assert.False(t, isErr)
	text, _ = callMCP(t, session, "get_luanti_settings", map[string]any{"filter": "my_custom"})
	assert.Equal(t, "{}", text)

	// git install with an existing name is rejected
	assert.NoError(t, app.Repos.ModRepo.Create(&types.Mod{Name: "existing", SourceType: types.SourceTypeGIT, ModType: types.ModTypeMod}))
	text, isErr = callMCP(t, session, "install_git_mod", map[string]any{"name": "existing", "url": "https://example.com/x.git"})
	assert.True(t, isErr)
	assert.Contains(t, text, "already installed")
}

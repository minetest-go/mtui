package web

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"mtui/types"
	"net/http"
	"path"
	"strings"
	"time"

	"github.com/luanti-golang/areasparser"
	"github.com/luanti-golang/mtdb/auth"
	mcpauth "github.com/modelcontextprotocol/go-sdk/auth"
	"github.com/modelcontextprotocol/go-sdk/mcp"
	"mtui/luanticonfig/depanalyzer"
)

// MCP (model context protocol) endpoint, mounted at /api/mcp
//
// Authentication: "Authorization: Bearer <jwt>" (an api-token from /api/token or a session token)
// The tool-list is built per session from the privileges in the token and the enabled features,
// every call re-checks both in case a feature got disabled in the meantime.

const mcpClaimsKey = "claims"

// guard for a single tool
type mcpGuard struct {
	Priv    string            // required priv, empty for "any logged in user"
	Feature types.FeatureName // required feature, empty for "always on"
}

func (api *Api) mcpCheck(claims *types.Claims, g mcpGuard) error {
	if claims == nil {
		return errors.New("unauthorized")
	}
	if g.Priv != "" && !claims.HasPriv(g.Priv) {
		return fmt.Errorf("forbidden, missing priv: '%s'", g.Priv)
	}
	if g.Feature != "" && !api.app.IsFeatureEnabled(g.Feature) {
		return fmt.Errorf("feature '%s' not enabled", g.Feature)
	}
	return nil
}

func mcpClaims(req *mcp.CallToolRequest) *types.Claims {
	if req.Extra == nil || req.Extra.TokenInfo == nil {
		return nil
	}
	claims, _ := req.Extra.TokenInfo.Extra[mcpClaimsKey].(*types.Claims)
	return claims
}

var mcpReadOnly = &mcp.ToolAnnotations{ReadOnlyHint: true}

// addMCPTool registers the tool on the server if the user is allowed to use it
func addMCPTool[In any](api *Api, s *mcp.Server, claims *types.Claims, g mcpGuard, t *mcp.Tool, fn func(ctx context.Context, claims *types.Claims, in In) (any, error)) {
	if api.mcpCheck(claims, g) != nil {
		// not available for this user (or feature disabled), hide it
		return
	}
	mcp.AddTool(s, t, func(ctx context.Context, req *mcp.CallToolRequest, in In) (*mcp.CallToolResult, any, error) {
		c := mcpClaims(req)
		if err := api.mcpCheck(c, g); err != nil {
			return nil, nil, err
		}
		out, err := fn(ctx, c, in)
		if err != nil {
			return nil, nil, err
		}
		data, err := json.Marshal(out)
		if err != nil {
			return nil, nil, fmt.Errorf("marshal result: %w", err)
		}
		res := &mcp.CallToolResult{Content: []mcp.Content{&mcp.TextContent{Text: string(data)}}}
		if bytes.HasPrefix(bytes.TrimSpace(data), []byte("{")) {
			// structuredContent must be a json object per spec, arrays/scalars are text-only
			res.StructuredContent = json.RawMessage(data)
		}
		return res, nil, nil
	})
}

// newMCPServer creates a server instance with the tools available to the given user
func (api *Api) newMCPServer(claims *types.Claims) *mcp.Server {
	s := mcp.NewServer(&mcp.Implementation{
		Name:    "mtui",
		Title:   "Luanti server ui",
		Version: api.app.Version,
	}, &mcp.ServerOptions{
		Instructions: "Tools to inspect a Luanti game server: online stats, players, logs, chat, bans, areas, mods and settings. " +
			"Timestamps are unix milliseconds unless noted otherwise.",
	})

	api.registerMCPTools(s, claims)
	api.registerMCPWriteTools(s, claims)
	return s
}

// MCPHandler returns the http handler for the streamable mcp transport
func (api *Api) MCPHandler() http.Handler {
	verifier := func(ctx context.Context, token string, r *http.Request) (*mcpauth.TokenInfo, error) {
		claims, err := api.parseClaims(token)
		if err != nil {
			return nil, fmt.Errorf("%w: %v", mcpauth.ErrInvalidToken, err)
		}
		ti := &mcpauth.TokenInfo{
			Scopes: claims.Privileges,
			UserID: claims.Username,
			Extra:  map[string]any{mcpClaimsKey: claims},
		}
		if claims.ExpiresAt != nil {
			ti.Expiration = claims.ExpiresAt.Time
		}
		return ti, nil
	}

	h := mcp.NewStreamableHTTPHandler(func(r *http.Request) *mcp.Server {
		ti := mcpauth.TokenInfoFromContext(r.Context())
		if ti == nil {
			return nil
		}
		claims, _ := ti.Extra[mcpClaimsKey].(*types.Claims)
		return api.newMCPServer(claims)
	}, &mcp.StreamableHTTPOptions{
		SessionTimeout: 30 * time.Minute,
	})

	return api.Feature(types.FEATURE_API, mcpauth.RequireBearerToken(verifier, nil)(h).ServeHTTP)
}

// helpers

func clampLimit(v, def, max int) int {
	if v <= 0 {
		return def
	}
	if v > max {
		return max
	}
	return v
}

// parses an RFC3339 timestamp into unix millis, empty string yields nil
func parseMCPTime(name, v string) (*int64, error) {
	if v == "" {
		return nil, nil
	}
	t, err := time.Parse(time.RFC3339, v)
	if err != nil {
		return nil, fmt.Errorf("invalid '%s' timestamp, expected RFC3339 (2006-01-02T15:04:05Z): %v", name, err)
	}
	ms := t.UnixMilli()
	return &ms, nil
}

func optString(v string) *string {
	if v == "" {
		return nil
	}
	return &v
}

// tool inputs

type mcpEmpty struct{}

type mcpPlayerSearchInput struct {
	NameLike string `json:"name_like,omitempty" jsonschema:"SQL LIKE pattern for the player name, for example 'bob%'"`
	Limit    int    `json:"limit,omitempty" jsonschema:"max number of results (default 20, max 100)"`
}

type mcpPlayerNameInput struct {
	Name string `json:"name" jsonschema:"the exact player name"`
}

type mcpLogSearchInput struct {
	From          string `json:"from,omitempty" jsonschema:"start of the time range (RFC3339)"`
	To            string `json:"to,omitempty" jsonschema:"end of the time range (RFC3339)"`
	Category      string `json:"category,omitempty" jsonschema:"log category: 'ui', 'minetest' (ingame events) or 'service'"`
	Event         string `json:"event,omitempty" jsonschema:"event name, see list_log_events for the available ones"`
	Username      string `json:"username,omitempty" jsonschema:"exact player name"`
	IPAddress     string `json:"ip_address,omitempty"`
	GeoCountry    string `json:"geo_country,omitempty" jsonschema:"country name as resolved by geoip"`
	Nodename      string `json:"nodename,omitempty" jsonschema:"node name, for example 'default:chest'"`
	Limit         int    `json:"limit,omitempty" jsonschema:"max number of results (default 50, max 500)"`
	SortAscending bool   `json:"sort_ascending,omitempty" jsonschema:"oldest entries first (default: newest first)"`
}

type mcpLogEventsInput struct {
	Category string `json:"category" jsonschema:"log category: 'ui', 'minetest' or 'service'"`
}

type mcpChatInput struct {
	Channel string `json:"channel,omitempty" jsonschema:"chat channel (default 'main')"`
	From    string `json:"from,omitempty" jsonschema:"start of the time range (RFC3339), omit for the latest messages"`
	To      string `json:"to,omitempty" jsonschema:"end of the time range (RFC3339), defaults to now"`
	Limit   int    `json:"limit,omitempty" jsonschema:"max number of (most recent) messages (default 100, max 1000)"`
}

type mcpIPInput struct {
	IP string `json:"ip" jsonschema:"IPv4 or IPv6 address"`
}

type mcpAreasInput struct {
	Owner string `json:"owner,omitempty" jsonschema:"only return areas owned by this player"`
	Limit int    `json:"limit,omitempty" jsonschema:"max number of results (default 100, max 1000)"`
}

type mcpSettingsInput struct {
	Filter string `json:"filter,omitempty" jsonschema:"only return settings whose key contains this substring"`
}

// tool definitions

func (api *Api) registerMCPTools(s *mcp.Server, claims *types.Claims) {
	// general

	addMCPTool(api, s, claims, mcpGuard{}, &mcp.Tool{
		Name:        "get_server_stats",
		Description: "Current server stats: uptime, max lag, time of day, memory and online players (positions only with the 'ban' priv). Empty if the engine has not reported stats within the last 10 seconds.",
		Annotations: mcpReadOnly,
	}, func(ctx context.Context, c *types.Claims, in mcpEmpty) (any, error) {
		return buildStats(c), nil
	})

	addMCPTool(api, s, claims, mcpGuard{}, &mcp.Tool{
		Name:        "healthcheck",
		Description: "Checks the database connections of the ui and the game, returns true if healthy.",
		Annotations: mcpReadOnly,
	}, func(ctx context.Context, c *types.Claims, in mcpEmpty) (any, error) {
		if err := api.app.Healthcheck(); err != nil {
			return nil, err
		}
		return true, nil
	})

	// players

	addMCPTool(api, s, claims, mcpGuard{Priv: types.PRIV_INTERACT}, &mcp.Tool{
		Name:        "search_players",
		Description: "Searches registered players by name, ordered by last login (newest first). Returns the total match count and the player infos (privileges, first/last login, stats).",
		Annotations: mcpReadOnly,
	}, func(ctx context.Context, c *types.Claims, in mcpPlayerSearchInput) (any, error) {
		limit := clampLimit(in.Limit, 20, 100)
		col := auth.LastLogin
		dir := auth.Descending
		search := &auth.AuthSearch{
			Usernamelike:   optString(in.NameLike),
			Limit:          &limit,
			OrderColumn:    &col,
			OrderDirection: &dir,
		}
		count, err := api.app.DBContext.Auth.Count(search)
		if err != nil {
			return nil, err
		}
		players, err := api.searchPlayers(search, c)
		if err != nil {
			return nil, err
		}
		return map[string]any{"count": count, "players": players}, nil
	})

	addMCPTool(api, s, claims, mcpGuard{Priv: types.PRIV_INTERACT}, &mcp.Tool{
		Name:        "get_player_info",
		Description: "Details for a single player: privileges, first/last login (unix seconds), health, stats and position (position only with the 'ban' priv).",
		Annotations: mcpReadOnly,
	}, func(ctx context.Context, c *types.Claims, in mcpPlayerNameInput) (any, error) {
		return api.getPlayerInfo(in.Name, c)
	})

	addMCPTool(api, s, claims, mcpGuard{Priv: types.PRIV_INTERACT}, &mcp.Tool{
		Name:        "list_privs",
		Description: "All privileges registered on the server with their descriptions.",
		Annotations: mcpReadOnly,
	}, func(ctx context.Context, c *types.Claims, in mcpEmpty) (any, error) {
		return api.readModStorageJSON("priv_info", map[string]*types.PrivInfo{})
	})

	addMCPTool(api, s, claims, mcpGuard{Priv: types.PRIV_INTERACT}, &mcp.Tool{
		Name:        "list_chatcommands",
		Description: "All chat commands registered on the server with their parameters, descriptions and required privileges.",
		Annotations: mcpReadOnly,
	}, func(ctx context.Context, c *types.Claims, in mcpEmpty) (any, error) {
		return api.readModStorageJSON("chatcommand_info", map[string]*types.ChatcommandInfo{})
	})

	// logs and moderation

	addMCPTool(api, s, claims, mcpGuard{Priv: types.PRIV_BAN}, &mcp.Tool{
		Name:        "search_logs",
		Description: "Searches the event log (joins, leaves, chat, node digs/places, ui actions, ...) with optional filters. Returns the total match count and the entries.",
		Annotations: mcpReadOnly,
	}, func(ctx context.Context, c *types.Claims, in mcpLogSearchInput) (any, error) {
		from, err := parseMCPTime("from", in.From)
		if err != nil {
			return nil, err
		}
		to, err := parseMCPTime("to", in.To)
		if err != nil {
			return nil, err
		}
		limit := clampLimit(in.Limit, 50, 500)
		search := &types.LogSearch{
			FromTimestamp: from,
			ToTimestamp:   to,
			Event:         optString(in.Event),
			Username:      optString(in.Username),
			IPAddress:     optString(in.IPAddress),
			GeoCountry:    optString(in.GeoCountry),
			Nodename:      optString(in.Nodename),
			Limit:         &limit,
			SortAscending: in.SortAscending,
		}
		if in.Category != "" {
			cat := types.LogCategory(in.Category)
			search.Category = &cat
		}

		count, err := api.app.Repos.LogRepository.Count(search)
		if err != nil {
			return nil, err
		}
		logs, err := api.app.Repos.LogRepository.Search(search)
		if err != nil {
			return nil, err
		}
		for _, l := range logs {
			// binary payloads are of no use here
			l.Attachment = nil
		}
		return map[string]any{"count": count, "logs": logs}, nil
	})

	addMCPTool(api, s, claims, mcpGuard{Priv: types.PRIV_BAN}, &mcp.Tool{
		Name:        "list_log_events",
		Description: "Lists the distinct event names recorded for a log category, useful as filter values for search_logs.",
		Annotations: mcpReadOnly,
	}, func(ctx context.Context, c *types.Claims, in mcpLogEventsInput) (any, error) {
		return api.app.Repos.LogRepository.GetEvents(types.LogCategory(in.Category))
	})

	addMCPTool(api, s, claims, mcpGuard{Priv: types.PRIV_BAN}, &mcp.Tool{
		Name:        "geoip_lookup",
		Description: "Resolves an IP address to country, city and ASN (if a geoip source is configured).",
		Annotations: mcpReadOnly,
	}, func(ctx context.Context, c *types.Claims, in mcpIPInput) (any, error) {
		return api.app.GeoipResolver.Resolve(in.IP), nil
	})

	addMCPTool(api, s, claims, mcpGuard{Priv: types.PRIV_BAN, Feature: types.FEATURE_XBAN}, &mcp.Tool{
		Name:        "get_ban_record",
		Description: "The xban record of a player: ban status, reason, expiry, known names/IPs and ban history.",
		Annotations: mcpReadOnly,
	}, func(ctx context.Context, c *types.Claims, in mcpPlayerNameInput) (any, error) {
		return api.app.GetOnlineXBanEntry(in.Name)
	})

	addMCPTool(api, s, claims, mcpGuard{Priv: types.PRIV_BAN, Feature: types.FEATURE_XBAN}, &mcp.Tool{
		Name:        "list_bans",
		Description: "All currently banned xban records.",
		Annotations: mcpReadOnly,
	}, func(ctx context.Context, c *types.Claims, in mcpEmpty) (any, error) {
		return api.getBannedRecords()
	})

	// chat

	addMCPTool(api, s, claims, mcpGuard{Priv: "shout", Feature: types.FeatureName("chat")}, &mcp.Tool{
		Name:        "get_chat_history",
		Description: "Chat messages of a channel, either the latest ones or within a time range, oldest first.",
		Annotations: mcpReadOnly,
	}, func(ctx context.Context, c *types.Claims, in mcpChatInput) (any, error) {
		channel := in.Channel
		if channel == "" {
			channel = "main"
		}
		limit := clampLimit(in.Limit, 100, 1000)

		var list []*types.ChatLog
		from, err := parseMCPTime("from", in.From)
		if err != nil {
			return nil, err
		}
		if from == nil {
			list, err = api.app.Repos.ChatLogRepo.GetLatest(channel, limit)
		} else {
			to, err2 := parseMCPTime("to", in.To)
			if err2 != nil {
				return nil, err2
			}
			if to == nil {
				now := time.Now().UnixMilli()
				to = &now
			}
			list, err = api.app.Repos.ChatLogRepo.Search(channel, *from, *to)
		}
		if err != nil {
			return nil, err
		}
		if len(list) > limit {
			list = list[len(list)-limit:]
		}
		return list, nil
	})

	// world

	addMCPTool(api, s, claims, mcpGuard{Feature: types.FEATURE_AREAS}, &mcp.Tool{
		Name:        "list_areas",
		Description: "Protected areas (from the areas mod) with owner, name, parent and corner positions.",
		Annotations: mcpReadOnly,
	}, func(ctx context.Context, c *types.Claims, in mcpAreasInput) (any, error) {
		areas, err := areasparser.ParseFile(path.Join(api.app.WorldDir, "areas.json"))
		if err != nil {
			return nil, err
		}
		limit := clampLimit(in.Limit, 100, 1000)
		result := make([]*areasparser.Area, 0)
		for _, a := range areas {
			if in.Owner != "" && a.Owner != in.Owner {
				continue
			}
			result = append(result, a)
			if len(result) >= limit {
				break
			}
		}
		return result, nil
	})

	// administration

	modGuard := mcpGuard{Priv: types.PRIV_SERVER, Feature: types.FEATURE_MODMANAGEMENT}

	addMCPTool(api, s, claims, modGuard, &mcp.Tool{
		Name:        "list_mods",
		Description: "Mods managed by the ui: name, author, source (git/contentdb), installed and latest known version, status.",
		Annotations: mcpReadOnly,
	}, func(ctx context.Context, c *types.Claims, in mcpEmpty) (any, error) {
		return api.app.Repos.ModRepo.GetAll()
	})

	addMCPTool(api, s, claims, modGuard, &mcp.Tool{
		Name:        "check_mod_updates",
		Description: "Queries the mod sources for new versions and stores them as 'latest_version' (does not install anything), returns the updated mod list.",
		Annotations: &mcp.ToolAnnotations{
			DestructiveHint: new(bool),
			IdempotentHint:  true,
		},
	}, func(ctx context.Context, c *types.Claims, in mcpEmpty) (any, error) {
		if err := api.app.ModManager.CheckUpdates(); err != nil {
			return nil, err
		}
		return api.app.Repos.ModRepo.GetAll()
	})

	addMCPTool(api, s, claims, modGuard, &mcp.Tool{
		Name:        "validate_mods",
		Description: "Analyzes the installed mods for missing or unsatisfied dependencies.",
		Annotations: mcpReadOnly,
	}, func(ctx context.Context, c *types.Claims, in mcpEmpty) (any, error) {
		return depanalyzer.AnalyzeDeps(
			path.Join(api.app.WorldDir, "worldmods"),
			path.Join(api.app.WorldDir, "game/mods"),
		)
	})

	addMCPTool(api, s, claims, mcpGuard{Priv: types.PRIV_SERVER, Feature: types.FEATURE_LUANTI_CONFIG}, &mcp.Tool{
		Name:        "get_luanti_settings",
		Description: "The current settings of the engine config (minetest.conf), optionally filtered by key.",
		Annotations: mcpReadOnly,
	}, func(ctx context.Context, c *types.Claims, in mcpSettingsInput) (any, error) {
		sts, err := api.app.GetSettingTypes()
		if err != nil {
			return nil, err
		}
		settings, err := api.app.ReadMTConfig(sts)
		if err != nil {
			return nil, err
		}
		if in.Filter == "" {
			return settings, nil
		}
		filtered := map[string]any{}
		for k, v := range settings {
			if strings.Contains(k, in.Filter) {
				filtered[k] = v
			}
		}
		return filtered, nil
	})
}

func (api *Api) readModStorageJSON(key string, target any) (any, error) {
	entry, err := api.app.DBContext.ModStorage.Get("mtui", []byte(key))
	if err != nil {
		return nil, err
	}
	if entry != nil {
		err = json.Unmarshal(entry.Value, &target)
	}
	return target, err
}

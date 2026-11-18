package web

import (
	"context"
	"errors"
	"fmt"
	"mtui/api/cdb"
	"mtui/luanticonfig"
	"mtui/luanticonfig/depanalyzer"
	"mtui/mail"
	"mtui/types"
	"mtui/types/command"
	"path"
	"slices"
	"strconv"
	"strings"
	"time"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

// MCP tools that change state (tier 2), see mcp.go for the auth and registration logic

// annotation helpers, the client decides about confirmations based on these
func mcpWrite(destructive, idempotent, openWorld bool) *mcp.ToolAnnotations {
	return &mcp.ToolAnnotations{
		DestructiveHint: &destructive,
		IdempotentHint:  idempotent,
		OpenWorldHint:   &openWorld,
	}
}

var mcpCDBRead = &mcp.ToolAnnotations{ReadOnlyHint: true, OpenWorldHint: &mcpTrue}

var mcpTrue = true

// parses durations like "90s", "30m", "12h", "7d" or "2w"
func parseBanDuration(v string) (time.Duration, error) {
	v = strings.TrimSpace(v)
	if v == "" {
		return 0, errors.New("empty duration")
	}
	unit := map[byte]time.Duration{'d': 24 * time.Hour, 'w': 7 * 24 * time.Hour}[v[len(v)-1]]
	if unit > 0 {
		n, err := strconv.ParseFloat(v[:len(v)-1], 64)
		if err != nil {
			return 0, fmt.Errorf("invalid duration '%s'", v)
		}
		return time.Duration(n * float64(unit)), nil
	}
	return time.ParseDuration(v)
}

// validates a primitive setting value against its type
func validateSettingValue(st *luanticonfig.SettingType, value string) error {
	switch st.Type {
	case "v3f", "noise_params_2d", "noise_params_3d":
		return fmt.Errorf("settings of type '%s' can only be changed in the web ui", st.Type)
	case "bool":
		if value != "true" && value != "false" {
			return errors.New("expected 'true' or 'false'")
		}
	case "int", "float":
		f, err := strconv.ParseFloat(value, 64)
		if err != nil || (st.Type == "int" && strings.ContainsAny(value, ".eE")) {
			return fmt.Errorf("expected a number of type '%s'", st.Type)
		}
		if st.Min != nil && f < *st.Min {
			return fmt.Errorf("value below minimum of %v", *st.Min)
		}
		if st.Max != nil && f > *st.Max {
			return fmt.Errorf("value above maximum of %v", *st.Max)
		}
	case "enum":
		if !slices.Contains(st.Choices, value) {
			return fmt.Errorf("expected one of: %s", strings.Join(st.Choices, ", "))
		}
	}
	return nil
}

// tool inputs

type mcpSendChatInput struct {
	Message string `json:"message" jsonschema:"the chat message, sent with your player name"`
	Channel string `json:"channel,omitempty" jsonschema:"chat channel (default 'main')"`
}

type mcpBanInput struct {
	Name     string `json:"name" jsonschema:"the exact player name"`
	Reason   string `json:"reason" jsonschema:"ban reason, shown to the player"`
	Duration string `json:"duration,omitempty" jsonschema:"ban duration like '30m', '12h', '7d' or '2w', omit for a permanent ban"`
}

type mcpInboxInput struct {
	UnreadOnly bool `json:"unread_only,omitempty"`
	Limit      int  `json:"limit,omitempty" jsonschema:"max number of (newest) messages (default 20, max 200)"`
}

type mcpSendMailInput struct {
	To      string `json:"to" jsonschema:"recipient player name"`
	Subject string `json:"subject"`
	Body    string `json:"body"`
}

type mcpCDBSearchInput struct {
	Query  string `json:"query,omitempty" jsonschema:"search text"`
	Author string `json:"author,omitempty"`
	Type   string `json:"type,omitempty" jsonschema:"package type: 'mod' (default), 'game' or 'txp'"`
	Sort   string `json:"sort,omitempty" jsonschema:"'score' (default), 'downloads', 'reviews' or 'last_release'"`
	Limit  int    `json:"limit,omitempty" jsonschema:"max number of results (default 20, max 100)"`
}

type mcpCDBPackageInput struct {
	Author           string   `json:"author" jsonschema:"package author on ContentDB"`
	Name             string   `json:"name" jsonschema:"package name on ContentDB"`
	SelectedPackages []string `json:"selected_packages,omitempty" jsonschema:"preferred 'author/name' packages for dependencies with multiple choices"`
}

type mcpGitModInput struct {
	Name    string `json:"name" jsonschema:"mod name (directory name, must match the mod's mod.conf name)"`
	URL     string `json:"url" jsonschema:"git repository url"`
	Branch  string `json:"branch,omitempty" jsonschema:"git branch (default 'master')"`
	ModType string `json:"mod_type,omitempty" jsonschema:"'mod' (default), 'game' or 'txp'"`
}

type mcpModVersionInput struct {
	ID      string `json:"id" jsonschema:"mod id, see list_mods"`
	Version string `json:"version" jsonschema:"target version: git commit hash or ContentDB release id, see 'latest_version' in list_mods"`
}

type mcpSetSettingInput struct {
	Key   string `json:"key" jsonschema:"setting key, for example 'max_users'"`
	Value string `json:"value" jsonschema:"new value as string (bool: 'true'/'false')"`
}

type mcpSettingKeyInput struct {
	Key string `json:"key" jsonschema:"setting key"`
}

type mcpChatcommandInput struct {
	Command string `json:"command" jsonschema:"chat command with parameters, for example 'time 6000' (leading slash optional)"`
}

func (api *Api) registerMCPWriteTools(s *mcp.Server, claims *types.Claims) {
	// chat

	addMCPTool(api, s, claims, mcpGuard{Priv: "shout", Feature: types.FeatureName("chat")}, &mcp.Tool{
		Name:        "send_chat",
		Description: "Sends a chat message to the game as your player.",
		Annotations: mcpWrite(false, false, true),
	}, func(ctx context.Context, c *types.Claims, in mcpSendChatInput) (any, error) {
		if strings.TrimSpace(in.Message) == "" {
			return nil, errors.New("empty message")
		}
		msg := &command.ChatMessage{Channel: in.Channel, Message: in.Message}
		if msg.Channel == "" {
			msg.Channel = "main"
		}
		return msg, api.sendChat(c, msg, nil)
	})

	// moderation

	xbanGuard := mcpGuard{Priv: types.PRIV_BAN, Feature: types.FEATURE_XBAN}

	addMCPTool(api, s, claims, xbanGuard, &mcp.Tool{
		Name:        "ban_player",
		Description: "Bans a player (and their known IPs) with xban, permanently or for a duration. Kicks the player if online.",
		Annotations: mcpWrite(true, true, false),
	}, func(ctx context.Context, c *types.Claims, in mcpBanInput) (any, error) {
		if in.Reason == "" {
			return nil, errors.New("a reason is required")
		}
		var seconds *int64
		if in.Duration != "" {
			d, err := parseBanDuration(in.Duration)
			if err != nil {
				return nil, err
			}
			if d < time.Second {
				return nil, errors.New("duration must be at least one second")
			}
			s := int64(d.Seconds())
			seconds = &s
		}
		if _, err := api.banPlayer(c, in.Name, in.Reason, seconds); err != nil {
			return nil, err
		}
		return api.app.GetOnlineXBanEntry(in.Name)
	})

	addMCPTool(api, s, claims, xbanGuard, &mcp.Tool{
		Name:        "unban_player",
		Description: "Lifts the xban ban of a player.",
		Annotations: mcpWrite(false, true, false),
	}, func(ctx context.Context, c *types.Claims, in mcpPlayerNameInput) (any, error) {
		if _, err := api.unbanPlayer(c, in.Name); err != nil {
			return nil, err
		}
		return api.app.GetOnlineXBanEntry(in.Name)
	})

	// mail

	mailGuard := mcpGuard{Priv: types.PRIV_INTERACT, Feature: types.FEATURE_MAIL}

	addMCPTool(api, s, claims, mailGuard, &mcp.Tool{
		Name:        "list_inbox",
		Description: "Ingame mails in your inbox, newest first. 'time' is in unix seconds.",
		Annotations: mcpReadOnly,
	}, func(ctx context.Context, c *types.Claims, in mcpInboxInput) (any, error) {
		e, err := api.app.Mail.GetEntry(c.Username)
		if err != nil {
			return nil, err
		}
		limit := clampLimit(in.Limit, 20, 200)
		result := make([]*mail.Message, 0)
		for i := len(e.Inbox) - 1; i >= 0 && len(result) < limit; i-- {
			if in.UnreadOnly && e.Inbox[i].Read {
				continue
			}
			result = append(result, e.Inbox[i])
		}
		return map[string]any{"total": len(e.Inbox), "messages": result}, nil
	})

	addMCPTool(api, s, claims, mailGuard, &mcp.Tool{
		Name:        "send_mail",
		Description: "Sends an ingame mail from your player to another player.",
		Annotations: mcpWrite(false, false, true),
	}, func(ctx context.Context, c *types.Claims, in mcpSendMailInput) (any, error) {
		m := &mail.Message{To: in.To, Subject: in.Subject, Body: in.Body}
		return m, api.sendMail(c, m)
	})

	// mods

	modGuard := mcpGuard{Priv: types.PRIV_SERVER, Feature: types.FEATURE_MODMANAGEMENT}

	addMCPTool(api, s, claims, modGuard, &mcp.Tool{
		Name:        "search_contentdb",
		Description: "Searches packages on ContentDB (content.luanti.org).",
		Annotations: mcpCDBRead,
	}, func(ctx context.Context, c *types.Claims, in mcpCDBSearchInput) (any, error) {
		q := &cdb.PackageQuery{
			Query:  in.Query,
			Author: in.Author,
			Type:   []cdb.PackageType{cdb.PackageTypeMod},
			Limit:  clampLimit(in.Limit, 20, 100),
			Sort:   cdb.PackageSortScore,
			Order:  cdb.PackageSortOrderDescending,
		}
		if in.Type != "" {
			q.Type = []cdb.PackageType{cdb.PackageType(in.Type)}
		}
		if in.Sort != "" {
			q.Sort = cdb.PackageSortType(in.Sort)
		}
		return cached_cdbcli.SearchPackages(q)
	})

	addMCPTool(api, s, claims, modGuard, &mcp.Tool{
		Name:        "resolve_contentdb_dependencies",
		Description: "Resolves the dependency tree of a ContentDB package against the installed mods: per dependency the installation candidates, the selected one and whether it is already installed.",
		Annotations: mcpCDBRead,
	}, func(ctx context.Context, c *types.Claims, in mcpCDBPackageInput) (any, error) {
		return api.resolveCDBDeps(in)
	})

	addMCPTool(api, s, claims, modGuard, &mcp.Tool{
		Name: "install_contentdb_package",
		Description: "Installs a ContentDB package with all missing dependencies (blocks until done). " +
			"Fails without installing anything if a dependency has no installation candidate. Use resolve_contentdb_dependencies first to review the plan.",
		Annotations: mcpWrite(false, false, true),
	}, func(ctx context.Context, c *types.Claims, in mcpCDBPackageInput) (any, error) {
		existing, err := api.app.Repos.ModRepo.GetByName(in.Name)
		if err != nil {
			return nil, err
		}
		if existing != nil {
			return nil, fmt.Errorf("'%s' is already installed (id: %s)", in.Name, existing.ID)
		}

		pkg, err := cached_cdbcli.GetDetails(in.Author, in.Name)
		if err != nil {
			return nil, fmt.Errorf("package lookup: %v", err)
		}
		if pkg == nil || pkg.Package == nil {
			return nil, fmt.Errorf("package '%s/%s' not found", in.Author, in.Name)
		}

		deps, err := api.resolveCDBDeps(in)
		if err != nil {
			return nil, err
		}

		// collect what needs to be installed, root package last
		todo := []*types.Mod{}
		missing := []string{}
		for _, d := range deps {
			if d.Installed {
				continue
			}
			parts := strings.SplitN(d.Selected, "/", 2)
			if len(d.Choices) == 0 || len(parts) != 2 {
				missing = append(missing, d.Name)
				continue
			}
			todo = append(todo, &types.Mod{Author: parts[0], Name: parts[1], ModType: types.ModTypeMod, SourceType: types.SourceTypeCDB})
		}
		if len(missing) > 0 {
			return nil, fmt.Errorf("no installation candidate for: %s", strings.Join(missing, ", "))
		}
		todo = append(todo, &types.Mod{Author: pkg.Author, Name: pkg.Name, ModType: types.ModType(pkg.Type), SourceType: types.SourceTypeCDB})

		installed := []*types.Mod{}
		for _, m := range todo {
			api.logModCreate(c, m, nil)
			if err := api.createMod(m); err != nil {
				return map[string]any{"installed": installed, "failed": m, "error": err.Error()}, fmt.Errorf("installing '%s/%s' failed after %d packages: %v", m.Author, m.Name, len(installed), err)
			}
			installed = append(installed, m)
		}
		return map[string]any{"installed": installed}, nil
	})

	addMCPTool(api, s, claims, modGuard, &mcp.Tool{
		Name:        "install_git_mod",
		Description: "Installs a mod, game or texturepack from a git repository (blocks until done). Dependencies are not resolved, check with validate_mods afterwards.",
		Annotations: mcpWrite(false, false, true),
	}, func(ctx context.Context, c *types.Claims, in mcpGitModInput) (any, error) {
		if in.Name == "" || in.URL == "" {
			return nil, errors.New("name and url are required")
		}
		existing, err := api.app.Repos.ModRepo.GetByName(in.Name)
		if err != nil {
			return nil, err
		}
		if existing != nil {
			return nil, fmt.Errorf("'%s' is already installed (id: %s)", in.Name, existing.ID)
		}
		m := &types.Mod{
			Name:       in.Name,
			URL:        in.URL,
			Branch:     in.Branch,
			ModType:    types.ModType(in.ModType),
			SourceType: types.SourceTypeGIT,
		}
		if m.Branch == "" {
			m.Branch = "master"
		}
		if m.ModType == "" {
			m.ModType = types.ModTypeMod
		}
		api.logModCreate(c, m, nil)
		return m, api.createMod(m)
	})

	addMCPTool(api, s, claims, modGuard, &mcp.Tool{
		Name:        "update_mod_version",
		Description: "Updates (or downgrades) an installed mod to the given version.",
		Annotations: mcpWrite(true, true, true),
	}, func(ctx context.Context, c *types.Claims, in mcpModVersionInput) (any, error) {
		return api.updateModVersion(c, in.ID, in.Version, nil)
	})

	// settings

	settingsGuard := mcpGuard{Priv: types.PRIV_SERVER, Feature: types.FEATURE_LUANTI_CONFIG}

	addMCPTool(api, s, claims, settingsGuard, &mcp.Tool{
		Name: "set_luanti_setting",
		Description: "Sets a setting in the engine config (minetest.conf). Values are validated against the setting type. " +
			"Primitive settings are applied at runtime as well, others (and 'secure.*') need a server restart.",
		Annotations: mcpWrite(true, true, false),
	}, func(ctx context.Context, c *types.Claims, in mcpSetSettingInput) (any, error) {
		if in.Key == "" {
			return nil, errors.New("empty key")
		}
		sts, err := api.app.GetSettingTypes()
		if err != nil {
			return nil, err
		}
		if err := validateSettingValue(settingTypeOrDefault(sts, in.Key), in.Value); err != nil {
			return nil, fmt.Errorf("invalid value for '%s': %v", in.Key, err)
		}
		st, err := api.setMTSetting(in.Key, &luanticonfig.Setting{Value: in.Value}, c, nil)
		if err != nil {
			return nil, err
		}
		runtime := !strings.HasPrefix(in.Key, "secure.") && runtime_set_allowed_types[st.Type]
		return map[string]any{"key": in.Key, "value": in.Value, "type": st.Type, "applied_at_runtime": runtime}, nil
	})

	addMCPTool(api, s, claims, settingsGuard, &mcp.Tool{
		Name:        "delete_luanti_setting",
		Description: "Removes a setting from the engine config (minetest.conf), the engine default applies afterwards.",
		Annotations: mcpWrite(true, true, false),
	}, func(ctx context.Context, c *types.Claims, in mcpSettingKeyInput) (any, error) {
		if in.Key == "" {
			return nil, errors.New("empty key")
		}
		return map[string]any{"key": in.Key, "deleted": true}, api.deleteMTSetting(in.Key, c, nil)
	})

	// shell

	addMCPTool(api, s, claims, mcpGuard{Feature: types.FEATURE_SHELL}, &mcp.Tool{
		Name:        "execute_chatcommand",
		Description: "Executes an ingame chat command as your player, with your ingame privileges. See list_chatcommands for the available commands.",
		Annotations: mcpWrite(true, false, false),
	}, func(ctx context.Context, c *types.Claims, in mcpChatcommandInput) (any, error) {
		cmd := strings.TrimPrefix(strings.TrimSpace(in.Command), "/")
		if cmd == "" {
			return nil, errors.New("empty command")
		}
		return api.executeChatcommand(c, cmd, nil)
	})
}

func (api *Api) resolveCDBDeps(in mcpCDBPackageInput) ([]*cdb.ResolvedDependency, error) {
	if in.Author == "" || in.Name == "" {
		return nil, errors.New("author and name are required")
	}
	ad, err := depanalyzer.AnalyzeDeps(
		path.Join(api.app.WorldDir, "worldmods"),
		path.Join(api.app.WorldDir, "game/mods"),
	)
	if err != nil {
		return nil, err
	}
	return cdb.ResolveDependencies(cached_cdbcli, in.Author+"/"+in.Name, in.SelectedPackages, ad.Installed)
}

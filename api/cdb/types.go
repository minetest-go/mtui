package cdb

import "time"

type PackageType string

const (
	PackageTypeMod         PackageType = "mod"
	PackageTypeGame        PackageType = "game"
	PackageTypeTexturepack PackageType = "txp"
)

type Package struct {
	Author           string      `json:"author"`
	Name             string      `json:"name"`
	Release          int         `json:"release"`
	ShortDescription string      `json:"short_description"`
	Thumbnail        string      `json:"thumbnail"` //https://content.minetest.net/thumbnails/1/6b28be927c.jpg
	Title            string      `json:"title"`
	Type             PackageType `json:"type"`
}

type DevStateType string

type PackageStateType string

type PackageDetails struct {
	*Package
	ContentWarnings []string         `json:"content_warnings"`
	CreatedAt       time.Time        `json:"created_at"`
	DevStage        DevStateType     `json:"dev_state"`
	DonateURL       string           `json:"donate_url"`
	Downloads       int              `json:"downloads"`
	Forums          int              `json:"forums"`
	GameSupport     []string         `json:"game_support"`
	IssueTracker    string           `json:"issue_tracker"`
	License         string           `json:"license"`
	LongDescription string           `json:"long_description"`
	Maintainers     []string         `json:"maintainers"`
	MediaLicense    string           `json:"media_license"`
	Provides        []string         `json:"provides"`
	Repo            string           `json:"repo"`
	Score           float64          `json:"score"`
	Screenshots     []string         `json:"screenshots"`
	State           PackageStateType `json:"state"`
	Tags            []string         `json:"tags"`
	URL             string           `json:"url"` //https://content.minetest.net/packages/Warr1024/nodecore/download/
	VideoURL        string           `json:"video_url"`
	Website         string           `json:"website"`
}

type DependencyInfo struct {
	IsOptional bool     `json:"is_optional"`
	Name       string   `json:"name"`
	Packages   []string `json:"packages"`
}

type PackageDependency map[string]*DependencyInfo

package cdb

import (
	"fmt"
	"net/url"
)

type PackageType string

const (
	PackageTypeMod         PackageType = "mod"
	PackageTypeGame        PackageType = "game"
	PackageTypeTexturepack PackageType = "txp"
)

type PackageSortType string

const (
	PackageSortScore       PackageSortType = "score"
	PackageSortReviews     PackageSortType = "reviews"
	PackageSortDownloads   PackageSortType = "downloads"
	PackageSortLastRelease PackageSortType = "last_release"
)

type PackageSortOrderType string

const (
	PackageSortOrderAscending  PackageSortOrderType = "asc"
	PackageSortOrderDescending PackageSortOrderType = "desc"
)

type PackageQuery struct {
	Type            []PackageType        `json:"type"`
	Query           string               `json:"query"`
	Author          string               `json:"author"`
	Limit           int                  `json:"limit"`
	Hide            []ContentWarning     `json:"hide"`
	Sort            PackageSortType      `json:"sort"`
	Order           PackageSortOrderType `json:"order"`
	ProtocolVersion int                  `json:"protocol_version"`
	EngineVersion   string               `json:"engine_version"`
}

func (q *PackageQuery) Params() url.Values {
	params := url.Values{}
	for _, t := range q.Type {
		params.Add("type", string(t))
	}
	if q.Query != "" {
		params.Add("q", q.Query)
	}
	if q.Author != "" {
		params.Add("author", q.Author)
	}
	if q.Limit > 0 {
		params.Add("limit", fmt.Sprintf("%d", q.Limit))
	}
	for _, cw := range q.Hide {
		params.Add("hide", cw.Name)
	}
	if q.Sort != "" {
		params.Add("sort", string(q.Sort))
	}
	if q.Order != "" {
		params.Add("order", string(q.Order))
	}
	if q.ProtocolVersion > 0 {
		params.Add("protocol_version", fmt.Sprintf("%d", q.ProtocolVersion))
	}
	if q.EngineVersion != "" {
		params.Add("engine_version", q.EngineVersion)
	}
	return params
}

type Package struct {
	Author           string      `json:"author"`
	Name             string      `json:"name"`
	Release          int         `json:"release"`
	ShortDescription string      `json:"short_description"`
	Thumbnail        string      `json:"thumbnail"` // https://content.minetest.net/thumbnails/1/6b28be927c.jpg
	Title            string      `json:"title"`
	Type             PackageType `json:"type"`
}

type Tag struct {
	Description string `json:"description"`
	IsProtected bool   `json:"is_protected"`
	Name        string `json:"name"`
	Title       string `json:"title"`
	Views       int    `json:"views"`
}

type ContentWarning struct {
	Description string `json:"description"`
	Name        string `json:"name"`
	Title       string `json:"title"`
}

type DevStateType string

type PackageStateType string

type PackageDetails struct {
	*Package
	ContentWarnings []string         `json:"content_warnings"`
	CreatedAt       string           `json:"created_at"`
	DevStage        DevStateType     `json:"dev_state"`
	DonateURL       string           `json:"donate_url"`
	Downloads       int              `json:"downloads"`
	Forums          int              `json:"forums"`
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
	URL             string           `json:"url"` // https://content.minetest.net/packages/Warr1024/nodecore/download/
	VideoURL        string           `json:"video_url"`
	Website         string           `json:"website"`
}

type DependencyInfo struct {
	IsOptional bool     `json:"is_optional"`
	Name       string   `json:"name"`
	Packages   []string `json:"packages"`
}

type PackageDependency map[string][]*DependencyInfo

type MinetestVersion struct {
	IsDev           bool   `json:"is_dev"`
	Name            string `json:"name"`
	ProtocolVersion int    `json:"protocol_version"`
}

type PackageRelease struct {
	Commit             string           `json:"commit"`
	Downloads          int              `json:"downloads"`
	ID                 int              `json:"id"`
	MinMinetestVersion *MinetestVersion `json:"min_minetest_version"`
	MaxMinetestVersion *MinetestVersion `json:"max_minetest_version"`
	ReleaseDate        string           `json:"release_date"`
	Title              string           `json:"title"`
	URL                string           `json:"url"` // /uploads/e4e8d405b0.zip
}

type PackageScreenshot struct {
	Approved     bool   `json:"approved"`
	CreatedAt    string `json:"created_at"`
	Height       int    `json:"height"`
	Width        int    `json:"width"`
	ID           int    `json:"id"`
	IsCoverImage bool   `json:"is_cover_image"`
	Order        int    `json:"order"`
	Title        string `json:"title"`
	URL          string `json:"url"` // https://content.minetest.net/uploads/Bz7IWGEnCH.png
}

type PackageThumbnails struct {
	Small  string `json:"small"`
	Medium string `json:"medium"`
	Large  string `json:"large"`
}

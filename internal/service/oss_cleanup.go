package service

import (
	"context"
	"strings"

	"github.com/gogf/gf/v2/frame/g"
)

// collectDualURLs appends non-empty thumb/original URLs from duals into dst.
func collectDualURLs(dst []string, duals ...DualURL) []string {
	for _, d := range duals {
		if u := strings.TrimSpace(d.Thumb); u != "" {
			dst = append(dst, u)
		}
		if u := strings.TrimSpace(d.Original); u != "" {
			dst = append(dst, u)
		}
	}
	return dst
}

// collectDualURLLists flattens multiple DualURL slices.
func collectDualURLLists(lists ...[]DualURL) []string {
	var out []string
	for _, list := range lists {
		out = collectDualURLs(out, list...)
	}
	return out
}

// urlsNotIn returns URLs present in old but not in keep (set difference).
func urlsNotIn(old, keep []string) []string {
	kept := make(map[string]struct{}, len(keep))
	for _, u := range keep {
		u = strings.TrimSpace(u)
		if u == "" {
			continue
		}
		kept[u] = struct{}{}
	}
	seen := make(map[string]struct{})
	var out []string
	for _, u := range old {
		u = strings.TrimSpace(u)
		if u == "" {
			continue
		}
		if _, ok := kept[u]; ok {
			continue
		}
		if _, dup := seen[u]; dup {
			continue
		}
		seen[u] = struct{}{}
		out = append(out, u)
	}
	return out
}

// deleteOSSURLsBestEffort deletes each URL via DeleteObject; failures are logged only.
func deleteOSSURLsBestEffort(ctx context.Context, urls []string) {
	for _, u := range urls {
		u = strings.TrimSpace(u)
		if u == "" {
			continue
		}
		if err := DeleteObject(u); err != nil {
			g.Log().Warningf(ctx, "oss cleanup delete failed url=%s err=%v", u, err)
		}
	}
}

// deleteReplacedDualsBestEffort deletes OSS objects for dual URLs in old but not in neu.
func deleteReplacedDualsBestEffort(ctx context.Context, oldDuals, newDuals []DualURL) {
	oldURLs := collectDualURLs(nil, oldDuals...)
	newURLs := collectDualURLs(nil, newDuals...)
	deleteOSSURLsBestEffort(ctx, urlsNotIn(oldURLs, newURLs))
}

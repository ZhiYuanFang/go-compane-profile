package service

import (
	"reflect"
	"testing"
)

func TestCollectDualURLs(t *testing.T) {
	got := collectDualURLs(nil,
		DualURL{Thumb: "a-t", Original: "a-o"},
		DualURL{Thumb: "", Original: "b-o"},
		DualURL{Thumb: "  ", Original: ""},
	)
	want := []string{"a-t", "a-o", "b-o"}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("collectDualURLs = %#v, want %#v", got, want)
	}
}

func TestUrlsNotIn(t *testing.T) {
	old := []string{"keep", "drop", "drop", "", "also-drop"}
	keep := []string{"keep", "new"}
	got := urlsNotIn(old, keep)
	want := []string{"drop", "also-drop"}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("urlsNotIn = %#v, want %#v", got, want)
	}
}

func TestUrlsNotInClearedGallery(t *testing.T) {
	old := collectDualURLLists(
		[]DualURL{{Thumb: "r-t", Original: "r-o"}},
		[]DualURL{{Thumb: "a-t", Original: "a-o"}},
	)
	newURLs := collectDualURLLists(
		[]DualURL{{Thumb: "r-t", Original: "r-o"}},
		[]DualURL{{Thumb: "", Original: ""}},
	)
	got := urlsNotIn(old, newURLs)
	want := []string{"a-t", "a-o"}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("cleared real urlsNotIn = %#v, want %#v", got, want)
	}
}

func TestDeleteOSSURLsBestEffortSkipsEmpty(t *testing.T) {
	// Must not panic; empty URLs are skipped before DeleteObject.
	deleteOSSURLsBestEffort(nil, []string{"", "  "})
}

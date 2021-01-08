package ucommon

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

//////////////////////////////////////////////////////////////////////////////
//
//
//
// Constants
//
//
//
//////////////////////////////////////////////////////////////////////////////

const (
	// AtomAuthorName is the name of the author to include in Atom feeds.
	AtomAuthorName = "Brandur Leach"

	// AtomAbsoluteURL is the absolute URL to use when generating Atom feeds.
	AtomAbsoluteURL = "https://mutelight.org"

	// AtomTag is a stable constant to use in Atom tags.
	AtomTag = "mutelight.org"

	// LayoutsDir is the source directory for view layouts.
	LayoutsDir = "./layouts"

	// MainLayout is the site's main layout.
	MainLayout = LayoutsDir + "/main.ace"

	// TitleSuffix is the suffix to add to the end of page and Atom titles.
	TitleSuffix = " — mutelight.org"

	// ViewsDir is the source directory for views.
	ViewsDir = "./views"
)

//////////////////////////////////////////////////////////////////////////////
//
//
//
// Functions
//
//
//
//////////////////////////////////////////////////////////////////////////////

// ExitWithError prints the given error to stderr and exits with a status of 1.
func ExitWithError(err error) {
	fmt.Fprintf(os.Stderr, "error: %v\n", err)
	os.Exit(1)
}

// ExtractSlug gets a slug for the given filename by using its basename
// stripped of file extension.
func ExtractSlug(source string) string {
	return strings.TrimSuffix(filepath.Base(source), filepath.Ext(source))
}

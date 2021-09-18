package main

import (
	"fmt"
	"html/template"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"sort"
	"strings"
	"sync"
	texttemplate "text/template"
	"time"

	"github.com/yosssi/ace"
	"golang.org/x/xerrors"

	"github.com/brandur/modulir"
	"github.com/brandur/modulir/modules/mace"
	"github.com/brandur/modulir/modules/matom"
	"github.com/brandur/modulir/modules/mfile"
	"github.com/brandur/modulir/modules/mmarkdownext"
	"github.com/brandur/modulir/modules/mtemplate"
	"github.com/brandur/modulir/modules/mtemplatemd"
	"github.com/brandur/modulir/modules/mtoml"
	"github.com/brandur/mutelight/modules/ucommon"
)

//////////////////////////////////////////////////////////////////////////////
//
//
//
// Variables
//
//
//
//////////////////////////////////////////////////////////////////////////////

// These are all objects that are persisted between build loops so that if
// necessary we can rebuild jobs that depend on them like index pages without
// reparsing all the source material. In each case we try to only reparse the
// sources if those source files actually changed.
var (
	articles []*Article
)

// A function map of template helpers which is the combined version of the maps
// from ftemplate, mtemplate, and mtemplatemd.
var htmlTemplateFuncMap template.FuncMap = mtemplate.CombineFuncMaps(
	mtemplate.FuncMap,
	mtemplatemd.FuncMap,
)

// Same as above, but for text templates.
var textTemplateFuncMap texttemplate.FuncMap = mtemplate.HTMLFuncMapToText(htmlTemplateFuncMap)

// List of common build dependencies, a change in any of which will trigger a
// rebuild on everything: partial views, JavaScripts, and stylesheets. Even
// though some of those changes will false positives, these sources are
// pervasive enough, and changes infrequent enough, that it's worth the
// tradeoff. This variable is a global because so many render functions access
// it.
var universalSources []string

//////////////////////////////////////////////////////////////////////////////
//
//
//
// Init
//
//
//
//////////////////////////////////////////////////////////////////////////////

func init() {
	mmarkdownext.FuncMap = textTemplateFuncMap
}

//////////////////////////////////////////////////////////////////////////////
//
//
//
// Build function
//
//
//
//////////////////////////////////////////////////////////////////////////////

func build(c *modulir.Context) []error {
	//
	// PHASE 0: Setup
	//
	// (No jobs should be enqueued here.)
	//

	c.Log.Debugf("Running build loop")

	// This is where we stored "versioned" assets like compiled JS and CSS.
	// These assets have a release number that we can increment and by
	// extension quickly invalidate.
	versionedAssetsDir := path.Join(c.TargetDir, "assets", Release)

	// A set of source paths that rebuild everything when any one of them
	// changes. These are dependencies that are included in more or less
	// everything: common partial views, JavaScript sources, and stylesheet
	// sources.
	universalSources = nil

	// Generate a list of partial views to add to universal sources.
	{
		sources, err := mfile.ReadDirCached(c, c.SourceDir+"/views",
			&mfile.ReadDirOptions{ShowMeta: true})
		if err != nil {
			return []error{err}
		}

		var partialViews []string
		for _, source := range sources {
			if strings.HasPrefix(filepath.Base(source), "_") {
				partialViews = append(partialViews, source)
			}
		}

		universalSources = append(universalSources, partialViews...)
	}

	// Generate a set of stylesheet sources to add to universal sources.
	{
		stylesheetSources, err := mfile.ReadDirCached(c, c.SourceDir+"/content/stylesheets",
			&mfile.ReadDirOptions{ShowMeta: true})
		if err != nil {
			return []error{err}
		}
		universalSources = append(universalSources, stylesheetSources...)
	}

	//
	// PHASE 1
	//
	// The build is broken into phases because some jobs depend on jobs that
	// ran before them. For example, we need to parse all our article metadata
	// before we can create an article index and render the home page (which
	// contains a short list of articles).
	//
	// After each phase, we call `Wait` on our context which will wait for the
	// worker pool to finish all its current work and restart it to accept new
	// jobs after it has.
	//
	// The general rule is to make sure that work is done as early as it
	// possibly can be. e.g. Jobs with no dependencies should always run in
	// phase 1. Try to make sure that as few phases as necessary.
	//

	{
		commonDirs := []string{
			c.TargetDir + "/a",
			versionedAssetsDir,
		}
		for _, dir := range commonDirs {
			err := mfile.EnsureDir(c, dir)
			if err != nil {
				return []error{nil}
			}
		}
	}

	//
	// Symlinks
	//

	{
		commonSymlinks := [][2]string{
			{c.SourceDir + "/content/images", c.TargetDir + "/assets/images"},
			{c.SourceDir + "/content/javascripts", versionedAssetsDir + "/javascripts"},
			{c.SourceDir + "/content/stylesheets", versionedAssetsDir + "/stylesheets"},
		}
		for _, link := range commonSymlinks {
			err := mfile.EnsureSymlink(c, link[0], link[1])
			if err != nil {
				return []error{nil}
			}
		}
	}

	//
	// Articles
	//

	var articlesChanged bool
	var articlesMu sync.Mutex

	{
		sources, err := mfile.ReadDirCached(c, c.SourceDir+"/content/articles", nil)
		if err != nil {
			return []error{err}
		}

		for _, s := range sources {
			source := s

			name := fmt.Sprintf("article: %s", filepath.Base(source))
			c.AddJob(name, func() (bool, error) {
				return renderArticle(c, source,
					&articles, &articlesChanged, &articlesMu)
			})
		}
	}

	//
	// Robots.txt
	//

	{
		c.AddJob("robots.txt", func() (bool, error) {
			return renderRobotsTxt(c)
		})
	}

	//
	//
	//
	// PHASE 2
	//
	//
	//

	if errors := c.Wait(); errors != nil {
		c.Log.Errorf("Cancelling next phase due to build errors")
		return errors
	}

	// Various sorts for anything that might need it.
	{
		sortArticles(articles)
	}

	// Index
	{
		c.AddJob("index", func() (bool, error) {
			return renderIndex(c, articles, articlesChanged)
		})
	}

	//
	// Articles
	//

	// Articles index (archive)
	{
		c.AddJob("articles index (Archive)", func() (bool, error) {
			return renderArticlesIndex(c, articles, articlesChanged)
		})
	}

	// Articles feed
	{
		c.AddJob("articles feed", func() (bool, error) {
			return renderArticlesFeed(c, articles, articlesChanged)
		})
	}

	return nil
}

//////////////////////////////////////////////////////////////////////////////
//
//
//
// Types
//
//
//
//////////////////////////////////////////////////////////////////////////////

// Article represents an article to be rendered.
type Article struct {
	// Content is the HTML content of the article. It isn't included as TOML
	// frontmatter, and is rather split out of an article's Markdown file,
	// rendered, and then added separately.
	Content string `toml:"-"`

	// Location is the place where the article was published. It may be empty.
	Location string `toml:"location"`

	// PublishedAt is when the article was published.
	PublishedAt *time.Time `toml:"published_at"`

	// Slug is a unique identifier for the article that also helps determine
	// where it's addressable by URL.
	Slug string `toml:"-"`

	// TinySlug is a short URL assigned to the article at `/a/<tiny slug>`
	// which redirects to the main article.
	//
	// This was almost certainly something that was never needed, but I added
	// it way back near 2010 when I was obsessed with URL shorteners, one of
	// the internet's worst ideas.
	TinySlug string `toml:"tiny_slug"`

	// Title is the article's title.
	Title string `toml:"title"`
}

func (a *Article) validate(source string) error {
	if a.Title == "" {
		return xerrors.Errorf("no title for article: %v", source)
	}

	if a.PublishedAt == nil {
		return xerrors.Errorf("no publish date for article: %v", source)
	}

	return nil
}

// articleYear holds a collection of articles grouped by year.
type articleYear struct {
	Year     int
	Articles []*Article
}

//////////////////////////////////////////////////////////////////////////////
//
//
//
// Private
//
//
//
//////////////////////////////////////////////////////////////////////////////

// getAceOptions gets a good set of default options for Ace template rendering
// for the project.
func getAceOptions(dynamicReload bool) *ace.Options {
	options := &ace.Options{FuncMap: htmlTemplateFuncMap}

	if dynamicReload {
		options.DynamicReload = true
	}

	return options
}

// Gets a map of local values for use while rendering a template and includes
// a few "special" values that are globally relevant to all templates.
func getLocals(title string, locals map[string]interface{}) map[string]interface{} {
	defaults := map[string]interface{}{
		"GoogleAnalyticsID": conf.GoogleAnalyticsID,
		"MetaDescription":   "",
		"Release":           Release,
		"MutelightEnv":      conf.MutelightEnv,
		"Title":             title,
		"TitleSuffix":       ucommon.TitleSuffix,
	}

	for k, v := range locals {
		defaults[k] = v
	}

	return defaults
}

func groupArticlesByYear(articles []*Article) []*articleYear {
	var year *articleYear
	var years []*articleYear

	for _, article := range articles {
		if year == nil || year.Year != article.PublishedAt.Year() {
			year = &articleYear{article.PublishedAt.Year(), nil}
			years = append(years, year)
		}

		year.Articles = append(year.Articles, article)
	}

	return years
}

func insertOrReplaceArticle(articles *[]*Article, article *Article) {
	for i, a := range *articles {
		if article.Slug == a.Slug {
			(*articles)[i] = article
			return
		}
	}

	*articles = append(*articles, article)
}

func renderArticle(c *modulir.Context, source string,
	articles *[]*Article, articlesChanged *bool, mu *sync.Mutex) (bool, error) {
	sourceChanged := c.Changed(source)
	viewsChanged := c.ChangedAny(append(
		[]string{
			ucommon.MainLayout,
			ucommon.ViewsDir + "/articles/show.ace",
		},
		universalSources...,
	)...)
	if !sourceChanged && !viewsChanged {
		return false, nil
	}

	var article Article
	data, err := mtoml.ParseFileFrontmatter(c, source, &article)
	if err != nil {
		return true, err
	}

	err = article.validate(source)
	if err != nil {
		return true, err
	}

	article.Slug = ucommon.ExtractSlug(source)

	content, err := mmarkdownext.Render(string(data), &mmarkdownext.RenderOptions{NoRetina: true})
	if err != nil {
		return true, err
	}
	article.Content = content

	locals := getLocals(article.Title, map[string]interface{}{
		"Article": article,
	})

	// Always use force context because if we made it to here we know that our
	// sources have changed.
	err = mace.RenderFile(c, ucommon.MainLayout, ucommon.ViewsDir+"/articles/show.ace",
		path.Join(c.TargetDir, article.Slug), getAceOptions(viewsChanged), locals)
	if err != nil {
		return true, err
	}

	// Ideally, this would be an actual redirect, but the combination of S3 +
	// CloudFront makes those somewhat difficult. Also, these should never
	// really get used from anywhere anymore so it doesn't actually matter that
	// much.
	if article.TinySlug != "" {
		filename := path.Join(c.TargetDir, "a", article.TinySlug)
		err := ioutil.WriteFile(
			filename,
			[]byte(fmt.Sprintf(
				`<!DOCTYPE html><html>please click through to: <strong><a href="/%s">/%s</a></strong></html>`,
				article.Slug, article.Slug,
			)),
			0o600,
		)
		if err != nil {
			return true, xerrors.Errorf("error writing file '%s': %w", filename, err)
		}
	}

	mu.Lock()
	insertOrReplaceArticle(articles, &article)
	*articlesChanged = true
	mu.Unlock()

	return true, nil
}

func renderArticlesFeed(c *modulir.Context, articles []*Article, articlesChanged bool) (bool, error) {
	if !articlesChanged {
		return false, nil
	}

	return renderFeed(c, "articles", "Articles", articles)
}

func renderArticlesIndex(c *modulir.Context, articles []*Article, articlesChanged bool) (bool, error) {
	viewsChanged := c.ChangedAny(append(
		[]string{
			ucommon.MainLayout,
			ucommon.ViewsDir + "/articles/index.ace",
		},
		universalSources...,
	)...)
	if !articlesChanged && !viewsChanged {
		return false, nil
	}

	articlesByYear := groupArticlesByYear(articles)

	locals := getLocals("Articles", map[string]interface{}{
		"ArticlesByYear": articlesByYear,
	})

	return true, mace.RenderFile(c, ucommon.MainLayout, ucommon.ViewsDir+"/articles/index.ace",
		c.TargetDir+"/archive", getAceOptions(viewsChanged), locals)
}

func renderFeed(_ *modulir.Context, slug, title string, articles []*Article) (bool, error) {
	filename := slug + ".atom"
	title += ucommon.TitleSuffix

	feed := &matom.Feed{
		Title: title,
		ID:    "tag:" + ucommon.AtomTag + ",2009:/" + slug,

		Links: []*matom.Link{
			{Rel: "self", Type: "application/atom+xml", Href: ucommon.AtomAbsoluteURL + "/" + filename},
			{Rel: "alternate", Type: "text/html", Href: ucommon.AtomAbsoluteURL},
		},
	}

	if len(articles) > 0 {
		feed.Updated = *articles[0].PublishedAt
	}

	for i, article := range articles {
		if i >= conf.NumAtomEntries {
			break
		}

		atomEntry := &matom.Entry{
			Title:     article.Title,
			Content:   &matom.EntryContent{Content: article.Content, Type: "html"},
			Published: *article.PublishedAt,
			Updated:   *article.PublishedAt,
			Link:      &matom.Link{Href: conf.AbsoluteURL + article.Slug},
			ID:        "tag:" + ucommon.AtomTag + "," + article.PublishedAt.Format("2006-01-02") + ":/" + article.Slug,

			AuthorName: ucommon.AtomAuthorName,
			AuthorURI:  conf.AbsoluteURL,
		}
		feed.Entries = append(feed.Entries, atomEntry)
	}

	f, err := os.Create(path.Join(conf.TargetDir, filename))
	if err != nil {
		return true, xerrors.Errorf("error creating file '%s': %w", filename, err)
	}
	defer f.Close()

	return true, feed.Encode(f, "  ")
}

func renderIndex(c *modulir.Context, articles []*Article, articlesChanged bool) (bool, error) {
	viewsChanged := c.ChangedAny(append(
		[]string{
			ucommon.MainLayout,
			ucommon.ViewsDir + "/index.ace",
		},
		universalSources...,
	)...)
	if !articlesChanged && !viewsChanged {
		return false, nil
	}

	// Cut off the number of items on the main page at some point.
	const numTopArticles = 10
	var topArticles []*Article
	for i := 0; i < numTopArticles; i++ {
		if i >= len(articles) {
			break
		}

		topArticles = append(topArticles, articles[i])
	}

	locals := getLocals("Mutelight", map[string]interface{}{
		"TopArticles": topArticles,
	})

	return true, mace.RenderFile(c, ucommon.MainLayout, ucommon.ViewsDir+"/index.ace",
		c.TargetDir+"/index.html", getAceOptions(viewsChanged), locals)
}

func renderRobotsTxt(c *modulir.Context) (bool, error) {
	if !c.FirstRun && !c.Forced {
		return false, nil
	}

	var content string
	if conf.Drafts {
		// Allow Twitterbot so that we can preview card images on dev.
		//
		// Disallow everything else.
		content = `User-agent: Twitterbot
Disallow:

User-agent: *
Disallow: /
`
	}

	filename := c.TargetDir + "/robots.txt"
	outFile, err := os.Create(filename)
	if err != nil {
		return true, xerrors.Errorf("error creating file '%s': %w", filename, err)
	}
	if _, err := outFile.WriteString(content); err != nil {
		return true, xerrors.Errorf("error writing file '%s': %w", filename, err)
	}
	outFile.Close()

	return true, nil
}

func sortArticles(articles []*Article) {
	sort.Slice(articles, func(i, j int) bool {
		return articles[j].PublishedAt.Before(*articles[i].PublishedAt)
	})
}

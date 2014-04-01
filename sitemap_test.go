package sitemap

import (
	"fmt"
	"log"
	"os"
	"path"
	"testing"
	"time"
)

var (
	itemResult = `
	<url>
		<loc>http://www.google.com</loc>
		<lastmod>2014-03-31T15:00:00+01:00</lastmod>
		<changefreq>hourly</changefreq>
		<priority>0.5</priority>
	</url>`
	sitemapResult = fmt.Sprintf(`<?xml version="1.0" encoding="UTF-8"?>
<urlset xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance"
	xsi:schemaLocation="http://www.sitemaps.org/schemas/sitemap/0.9 http://www.sitemaps.org/schemas/sitemap/0.9/sitemap.xsd"
	xmlns="http://www.sitemaps.org/schemas/sitemap/0.9">%s
</urlset>`, itemResult)

	sitemapIndexItemResult = `
	<sitemap>
		<loc>http://www.google.com/sitemap.xml.gz</loc>
		<lastmod>2014-03-31T15:00:00+01:00</lastmod>
	</sitemap>`

	sitemapIndexResult = fmt.Sprintf(`<?xml version="1.0" encoding="UTF-8"?>
<sitemapindex xmlns="http://www.sitemaps.org/schemas/sitemap/0.9">%s
</sitemapindex>
`, sitemapIndexItemResult)
)

func TestFileGeneration(t *testing.T) {
	testDir := os.TempDir() + "/sitemap"
	err := os.Mkdir(testDir, os.ModeDir)
	if err != nil {
		log.Fatalf("could not create temporary test directory: %v", err)
	}
	defer func() {
		os.RemoveAll(testDir)
	}()

	lastMod, _ := time.Parse(time.RFC3339, "2014-03-31T15:00:00+01:00")

	// Sitemap item
	item := SitemapItem{
		"http://www.google.com",
		lastMod,
		"hourly",
		0.5,
	}

	if item.String() != itemResult {
		t.Errorf("Expected sitemap item to be %s, actual: %s", itemResult, item.String())
	}

	// Sitemap
	sitemap := Sitemap{
		[]SitemapItem{
			item,
		},
	}

	if sitemap.String() != sitemapResult {
		t.Errorf("Expected sitemap to be %s, actual: %s", sitemapResult, sitemap.String())
	}

	// Save sitemap to test directory
	err = sitemap.ToFile(testDir + "/sitemap.xml.gz")
	if err != nil {
		t.Errorf("Could not save the sitemap to a file: %v", err)
	}

	// SitemapIndexItem
	sitemapIndexItem := SitemapIndexItem{
		"http://www.google.com/sitemap.xml.gz",
		lastMod,
	}

	if sitemapIndexItem.String() != sitemapIndexItemResult {
		t.Errorf("Expected sitemap index item to be %s, actual: %s", sitemapIndexItemResult, sitemapIndexItem.String())
	}

	// SitemapIndex
	sitemapIndex := SitemapIndex{
		[]SitemapIndexItem{
			sitemapIndexItem,
		},
	}

	if sitemapIndex.String() != sitemapIndexResult {
		t.Errorf("Expected sitemap index to be %s, actual: %s", sitemapIndexResult, sitemapIndex.String())
	}

	sitemapIndex2, err := NewIndexFromDir(testDir, "http://www.google.com/")
	if err != nil {
		log.Fatalf("could not create sitemap index from directory: %v", err)
	}

	file, err := os.Open(path.Join(testDir, "sitemap.xml.gz"))
	if err != nil {
		log.Fatalf("could not open file 'sitemap.xml.gz' in test dir: %v", err)
	}

	fileinfo, err := file.Stat()
	if err != nil {
		log.Fatalf("could not stat file 'sitemap.xml.gz' in test dir: %v", err)
	}

	sitemapIndexItem2 := SitemapIndexItem{
		"http://www.google.com/sitemap.xml.gz",
		fileinfo.ModTime(),
	}
	sitemapIndexResult2 := fmt.Sprintf(`<?xml version="1.0" encoding="UTF-8"?>
<sitemapindex xmlns="http://www.sitemaps.org/schemas/sitemap/0.9">%s
</sitemapindex>
`, sitemapIndexItem2.String())

	if sitemapIndex2.String() != sitemapIndexResult2 {
		t.Errorf("Expected sitemap index created from dir '%s' to be %s, actual: %s", testDir, sitemapIndexResult2, sitemapIndex2.String())
	}

	// Save sitemap index to test directory
	err = sitemap.ToFile(testDir + "/sitemap-index.xml.gz")
	if err != nil {
		t.Errorf("Could not save the sitemap index to a file: %v", err)
	}

}

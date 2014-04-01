package sitemap

import (
	"compress/gzip"
	"fmt"
	"os"
	"path"
	"path/filepath"
	"strings"
	"time"
)

const (
	// MaxSitemapItems is the maximum number of items for a single sitemap
	MaxSitemapItems = 50000

	// SitemapXML is the XML structure for urlset in sitemaps
	SitemapXML = `<?xml version="1.0" encoding="UTF-8"?>
<urlset xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance"
	xsi:schemaLocation="http://www.sitemaps.org/schemas/sitemap/0.9 http://www.sitemaps.org/schemas/sitemap/0.9/sitemap.xsd"
	xmlns="http://www.sitemaps.org/schemas/sitemap/0.9">%s</urlset>`

	// SitemapItemXML is the XML format for the URL item in sitemap
	SitemapItemXML = `
	<url>
		<loc>%s</loc>
		<lastmod>%s</lastmod>
		<changefreq>%s</changefreq>
		<priority>%.1f</priority>
	</url>
`

	// SitemapIndexXML is the XML structure of a sitemap index
	SitemapIndexXML = `<?xml version="1.0" encoding="UTF-8"?>
<sitemapindex xmlns="http://www.sitemaps.org/schemas/sitemap/0.9">%s</sitemapindex>
`

	// SitemapIndexItemXML is the XML structure of a sitemap index item
	SitemapIndexItemXML = `
	<sitemap>
		<loc>%s</loc>
		<lastmod>%s</lastmod>
	</sitemap>
`
)

// Sitemap represent a sitemap
type Sitemap struct {
	items []SitemapItem
}

// Add adds a sitemap item to the sitemap
func (s *Sitemap) Add(item SitemapItem) error {
	if len(s.items) >= MaxSitemapItems {
		return fmt.Errorf("your sitemap has reached the maximum number of items which is %v", MaxSitemapItems)
	}

	s.items = append(s.items, item)

	return nil
}

// String return the string format of the sitemap
func (s *Sitemap) String() string {
	var items []string
	for _, item := range s.items {
		items = append(items, item.String())
	}
	return fmt.Sprintf(SitemapXML, strings.Join(items, `
`))
}

// ToFile saves a sitemap to a file with either extension .xml or .gz.
// If extension is .gz, the file will be gzipped.
func (s *Sitemap) ToFile(path string) error {
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()

	ext := filepath.Ext(file.Name())
	if ext != ".xml" && ext != ".gz" {
		return fmt.Errorf("filename %s does not have extension .xml or .gz, extension %s given", file.Name(), ext)
	}

	// Gzip
	if ext == ".gz" {
		zip := gzip.NewWriter(file)
		defer zip.Close()

		_, err = zip.Write([]byte(s.String()))
		if err != nil {
			return err
		}
	} else {
		_, err = file.Write([]byte(s.String()))
		if err != nil {
			return err
		}
	}

	return nil
}

// SitemapItem represents an item in the sitemap
type SitemapItem struct {
	Loc        string
	LastMod    time.Time
	ChangeFreq string
	Priority   float32
}

// String return the string format of the sitemap item
func (i *SitemapItem) String() string {
	return fmt.Sprintf(SitemapItemXML, i.Loc, i.LastMod.Format(time.RFC3339), i.ChangeFreq, i.Priority)
}

// SitemapIndex is an index for multiple sitemaps
type SitemapIndex struct {
	items []SitemapIndexItem
}

// Add adds a sitemap to the sitemap index
func (s *SitemapIndex) Add(item SitemapIndexItem) {
	s.items = append(s.items, item)
}

// String return the string format of the sitemap index
func (s *SitemapIndex) String() string {
	var items []string
	for _, item := range s.items {
		items = append(items, item.String())
	}

	return fmt.Sprintf(SitemapIndexXML, strings.Join(items, `
`))
}

// SitemapIndexItem represents an item in the sitemap index
type SitemapIndexItem struct {
	Loc     string
	LastMod time.Time
}

// String return the string format of the sitemap item
func (i *SitemapIndexItem) String() string {
	return fmt.Sprintf(SitemapIndexItemXML, i.Loc, i.LastMod.Format(time.RFC3339))
}

// ToFile saves a sitemap index to a file with either extension .xml or .gz.
// If extension is .gz, the file will be gzipped.
func (s *SitemapIndex) ToFile(path string) error {
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()

	ext := filepath.Ext(file.Name())
	if ext != ".xml" && ext != ".gz" {
		return fmt.Errorf("filename %s does not have extension .xml or .gz, extension %s given", file.Name(), ext)
	}

	// Gzip
	if ext == ".gz" {
		zip := gzip.NewWriter(file)
		defer zip.Close()

		_, err = zip.Write([]byte(s.String()))
		if err != nil {
			return err
		}
	} else {
		_, err = file.Write([]byte(s.String()))
		if err != nil {
			return err
		}
	}

	return nil
}

// NewIndexFromDir creates a sitemap index by scanning a folder for files.
// The files modified time will be used as LastMod.
func NewIndexFromDir(dir, pathPrefix string) (*SitemapIndex, error) {
	s := &SitemapIndex{
		make([]SitemapIndexItem, 0),
	}

	f, err := os.Open(dir)
	if err != nil {
		return s, err
	}

	list, err := f.Readdir(-1)
	f.Close()
	if err != nil {
		return s, err
	}

	for _, file := range list {
		ext := filepath.Ext(file.Name())
		if ext == ".xml" || ext == ".gz" {
			var sitemapPath string
			if pathPrefix != "" {
				sitemapPath = pathPrefix + file.Name()
			} else {
				sitemapPath = path.Join(dir, file.Name())
			}
			item := SitemapIndexItem{
				sitemapPath,
				file.ModTime(),
			}

			s.Add(item)
		}
	}

	return s, nil
}

## A Go package to generate sitemaps
Follows the formats and guidelines for [sitemaps.org](http://www.sitemaps.org/). More information at [Google support's answer](https://support.google.com/webmasters/answer/183668)

### Installation
```shell
go get github.com/pengux/sitemap
```

### Usage
```go
// Import package
import (
	...
	"github.com/pengux/sitemap"
)

// Sitemap item
item := SitemapItem{
	"http://www.google.com",
	time.Now(),
	"hourly",
	0.5,
}

// Sitemap
sitemap := Sitemap{
	[]SitemapItem{
		item,
	},
}

fmt.Print(sitemap.String()) // Output the sitemap as string
sitemap.ToFile("sitemap.xml.gz") // Save sitemap to a gzipped file


// SitemapIndexItem
sitemapIndexItem := SitemapIndexItem{
	"http://www.google.com/sitemap.xml.gz",
	time.Now(),
}

// SitemapIndex
sitemapIndex := SitemapIndex{
	[]SitemapIndexItem{
		sitemapIndexItem,
	},
}

fmt.Print(sitemapIndex.String()) // Output the sitemap index as string
sitemapIndex.ToFile("sitemap.xml.gz") // Save sitemap to a gzipped file

// Create sitemap index from a directory containing sitemap files
sitemapIndex, err := NewIndexFromDir(path, "http://www.google.com/")
```

### TODO
- Support sitemap extensions (images, videos, mobile, news)


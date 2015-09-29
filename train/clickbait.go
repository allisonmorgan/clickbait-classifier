package clickbait

import (
	"bytes"
	"io/ioutil"
	"net/http"

	"github.com/PuerkitoBio/goquery"
)

var (
	// these are big files with the urls of many article
	BUZZFEED_SITEMAP  = "http://www.buzzfeed.com/go/sitemap_buzz.xml"     // ~20,000
	REUTERS_SITEMAP   = "http://www.reuters.com/sitemap_news_index.xml"   //   1,000
	BLOOMBERG_SITEMAP = "http://www.bloomberg.com/feeds/sitemap_news.xml" //     150
	ALJAZEERA_SITEMAP = "http://america.aljazeera.com/sitemap.0.xml"      //   1,000
	CNN_SITEMAP       = "http://www.cnn.com/sitemaps/sitemap-news.xml"    // 	 250

	// classes
	CLICKBAIT     = "clickbait"
	NOT_CLICKBAIT = "not_clickbait"
)

type Document struct {
	Text    string
	Classes []string
}

// scrape buzzfeed and npr for headlines
func BuzzfeedScrape(n int) ([]Document, error) {
	// hit their sitemap feed (may take a while)
	body, status, err := Get(BUZZFEED_SITEMAP)
	if err != nil || status != 200 {
		return nil, err
	}

	doc, err := NewReaderFromBody(body)
	if err != nil {
		return nil, err
	}

	// parse the sitemap and grab the first 100 urls
	urls := make([]string, 0, n)
	newDoc := doc.Find("loc")
	newDoc.Each(func(i int, s *goquery.Selection) {
		if i > n {
			return
		}
		title := s.Text()
		urls = append(urls, title)
	})

	// for each doc, grab the title
	documents := make([]Document, 0)
	for _, urlStr := range urls {
		body, status, err = Get(urlStr)
		if err != nil || status != 200 {
			continue
		}

		doc, err = NewReaderFromBody(body)
		if err != nil {
			continue
		}

		newDoc = doc.Find("title")
		newDoc.Each(func(i int, s *goquery.Selection) {
			// articles are marked up with more than one title: actual title, and "Buzzfeed"
			if i > 0 {
				return
			}
			title := s.Text()
			documents = append(documents, Document{title, []string{CLICKBAIT}})
		})
	}

	return documents, nil
}

func ReutersScrape(n int) ([]Document, error) {
	// hit their sitemap feed (may take a while)
	body, status, err := Get(REUTERS_SITEMAP)
	if err != nil || status != 200 {
		return nil, err
	}

	doc, err := NewReaderFromBody(body)
	if err != nil {
		return nil, err
	}

	documents := make([]Document, 0)
	// unlike Buzzfeed, Reuters' sitemap has the title in it
	newDoc := doc.Find("news\\:title")
	newDoc.Each(func(i int, s *goquery.Selection) {
		if i > n {
			return
		}

		title := s.Text()
		documents = append(documents, Document{title, []string{NOT_CLICKBAIT}})
	})

	return documents, nil
}

func BloombergScrape(n int) ([]Document, error) {
	// hit their sitemap feed (may take a while)
	body, status, err := Get(BLOOMBERG_SITEMAP)
	if err != nil || status != 200 {
		return nil, err
	}

	doc, err := NewReaderFromBody(body)
	if err != nil {
		return nil, err
	}

	documents := make([]Document, 0)
	// similar to Reuters', Bloomberg's sitemap has the title in it
	newDoc := doc.Find("news\\:title")
	newDoc.Each(func(i int, s *goquery.Selection) {
		if i > n {
			return
		}

		title := s.Text()
		documents = append(documents, Document{title, []string{NOT_CLICKBAIT}})
	})

	return documents, nil
}

func CNNScrape(n int) ([]Document, error) {
	// hit their sitemap feed (may take a while)
	body, status, err := Get(CNN_SITEMAP)
	if err != nil || status != 200 {
		return nil, err
	}

	doc, err := NewReaderFromBody(body)
	if err != nil {
		return nil, err
	}

	documents := make([]Document, 0)
	// similar to Reuters', Bloomberg's sitemap has the title in it
	newDoc := doc.Find("news\\:title")
	newDoc.Each(func(i int, s *goquery.Selection) {
		if i > n {
			return
		}

		title := s.Text()
		documents = append(documents, Document{title, []string{NOT_CLICKBAIT}})
	})

	return documents, nil
}

func AljazeeraScrape(n int) ([]Document, error) {
	// hit their sitemap feed (may take a while)
	body, status, err := Get(ALJAZEERA_SITEMAP)
	if err != nil || status != 200 {
		return nil, err
	}

	doc, err := NewReaderFromBody(body)
	if err != nil {
		return nil, err
	}

	// parse the sitemap and grab the first 100 urls
	urls := make([]string, 0, n)
	newDoc := doc.Find("loc")
	newDoc.Each(func(i int, s *goquery.Selection) {
		if i > n {
			return
		}
		title := s.Text()
		urls = append(urls, title)
	})

	// for each doc, grab the title
	documents := make([]Document, 0)
	for _, urlStr := range urls {
		body, status, err = Get(urlStr)
		if err != nil || status != 200 {
			continue
		}

		doc, err = NewReaderFromBody(body)
		if err != nil {
			continue
		}

		newDoc = doc.Find("title")
		newDoc.Each(func(i int, s *goquery.Selection) {
			// articles are marked up with more than one title: actual title, and "Buzzfeed"
			if i > 0 {
				return
			}
			title := s.Text()
			documents = append(documents, Document{title, []string{CLICKBAIT}})
		})
	}

	return documents, nil
}

func NewReaderFromBody(body []byte) (*goquery.Document, error) {
	uReader := bytes.NewReader(body)
	doc, err := goquery.NewDocumentFromReader(uReader)
	if err != nil {
		return nil, err
	}
	return doc, nil
}

func Get(urlString string) ([]byte, int, error) {
	resp, err := http.Get(urlString)
	if err != nil {
		return nil, 0, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	return body, resp.StatusCode, err
}

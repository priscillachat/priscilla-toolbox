# image command

```bash
Usage:
  -cseid string
        Google custom search ID
  -csekey string
        Google custom search Key
  -giphykey string
        Giphy api key (leave default to use public beta key) (default "dc6zaTOxFJmzC")
  -mode string
        Mode: image or gif (default "image")
  -query string
        search query
  -safe string
        (image mode only) Safe search: high/medium/off (default "medium")
```

# Important note about google custom search

Since google has shutdown their free image search API, you'll need to create an
google custom search engine (https://cse.google.com/) allows 100 search per day
free and fill in the `cseid` with the 'search engine ID' (can be found
under 'Basic'->'Details'->'Search engine ID'). Then you need to go create an
api key from Google Developers Console (instructions here:
https://developers.google.com/custom-search/json-api/v1/introduction) and fill
that into the `csekey`.


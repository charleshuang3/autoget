package indexers

import (
	"bytes"
	_ "embed"
	"text/template"

	"github.com/rs/zerolog/log"
)

//go:embed message_templates/rss.md
var rssTemplateContent string

var rssTemplate *template.Template

func init() {
	var err error
	rssTemplate, err = template.New("rss").Parse(rssTemplateContent)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to parse RSS template")
	}
}

type RSSResultTemplateData struct {
	Indexer                string
	DownloadStarted        []string
	DownloadPendingToStart []string
}

func RenderRSSResult(indexer string, downloadStarted []string, downloadPendingToStart []string) (string, error) {
	data := RSSResultTemplateData{
		Indexer:                indexer,
		DownloadStarted:        downloadStarted,
		DownloadPendingToStart: downloadPendingToStart,
	}

	var buf bytes.Buffer
	err := rssTemplate.Execute(&buf, data)
	if err != nil {
		return "", err
	}
	return buf.String(), nil
}

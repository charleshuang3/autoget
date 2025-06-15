# {{.Indexer}} RSS

{{if .DownloadStarted}}
## Download Started
{{range .DownloadStarted}}
- {{.}}
{{end}}
{{end}}

{{if .DownloadPendingToStart}}
## Download Pending to Start
{{range .DownloadPendingToStart}}
- {{.}}
{{end}}
{{end}}

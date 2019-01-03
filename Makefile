github-contribution-collector: $(wildcard *.go)
	go build

github-contribution-collector.zip: github-contribution-collector
	zip github-contribution-collector.zip github-contribution-collector

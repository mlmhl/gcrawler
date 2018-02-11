package types

// An Item represents an effective record crawled by Spider.
type Item interface {
	Content() string
}

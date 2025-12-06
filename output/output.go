package output

import "jellynfo/jellydata"

// Outputter defines the interface for displaying Jellyfin data.
type Outputter interface {
	Display(data *jellydata.JellyfinData) error
}
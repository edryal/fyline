package protocol

type Channel struct {
	// ID is the stable identifier. Never changes for the life of the channel.
	// Messages reference this, selection keys off this, renames don't touch it.
	ID string `json:"id"`

	// Name is the label shown in the sidebar. Mutable.
	Name string `json:"name"`

	// Position controls sidebar ordering (lower = higher)
	Position int `json:"position"`
}

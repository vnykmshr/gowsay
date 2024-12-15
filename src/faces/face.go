package faces

type Face struct {
	Eyes     string
	Tongue   string
	Thoughts string
	Cow      string
}

func New(cow, mood string) (*Face, error) {
	f := &Face{
		Cow:    cow,
		Eyes:   "oo",
		Tongue: "  ",
	}

	switch mood {
	case "borg":
		f.Eyes = "=="
	case "dead":
		f.Eyes = "xx"
		f.Tongue = "U "
	case "greedy":
		f.Eyes = "$$"
	case "paranoid":
		f.Eyes = "@@"
	case "stoned":
		f.Eyes = "**"
		f.Tongue = "U "
	case "tired":
		f.Eyes = "--"
	case "wired":
		f.Eyes = "OO"
	case "young":
		f.Eyes = ".."
	default:
		// proceed with default mood
	}

	return f, nil
}

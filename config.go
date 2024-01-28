package nxboundary

type Config struct {
	DepConstraints map[string][]string
}

func (c *Config) IsTagAllowed(tag string, tagToCheck string) bool {
	if tag == tagToCheck {
		return true
	}

	if allowedTags, ok := c.DepConstraints[tag]; ok {
		for _, allowedTag := range allowedTags {
			if allowedTag == tagToCheck {
				return true
			}
		}
	}

	return false
}

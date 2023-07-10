package railgun

type RailGun struct {
	client         *request
	authentication *authenticationConfig
}

// Todo: support non-singleton mode
func Init() *RailGun {
	return &RailGun{
		client: newRequest(),
	}
}

func (r *RailGun) Parse(token string, queryCollection bool) (Parser, error) {
	parser, err := r.getMediaTokenType(token, queryCollection)
	if err != nil {
		return nil, err
	}

	parser.setRequestClient(r.client)

	return parser, nil
}

func (r *RailGun) GetConfig() *conf {
	return &Config
}

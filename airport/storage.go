package airport

type Storage struct {
	Datas [][2]string
}

func (s Storage) Read(key string) (Airport, bool) {

	for _, data := range s.Datas {
		if data[0] == key {
			return Airport{
				Code: data[0],
				Name: data[1],
			}, true
		}
	}

	return Airport{}, false
}

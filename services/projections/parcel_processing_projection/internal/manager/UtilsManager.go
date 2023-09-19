package manager

// il calcolo viene effettuato solo per le facility
func CalculatePathPercentage(path []string, actualLocationID string) int {
	pathLen := len(path)

	for i, locationID := range path {
		if locationID == actualLocationID {
			// Se la stringa target è stata trovata, calcola la posizione percentuale
			locationPathIndex := i + 1
			position := (float64(locationPathIndex) / float64(pathLen)) * 100
			return int(position)
		}
	}

	return 0
}

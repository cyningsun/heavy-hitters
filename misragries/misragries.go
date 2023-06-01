package misragries

type MisraGries struct {
	k       int
	counter map[string]int
}

func NewMisraGries(k int) *MisraGries {
	return &MisraGries{
		k:       k,
		counter: make(map[string]int),
	}
}

func (mg *MisraGries) ProcessElement(element string) {
	if count, ok := mg.counter[element]; ok {
		mg.counter[element] = count + 1
	} else if len(mg.counter) < mg.k {
		mg.counter[element] = 1
	} else {
		for key := range mg.counter {
			mg.counter[key]--
			if mg.counter[key] == 0 {
				delete(mg.counter, key)
			}
		}
	}
}

func (mg *MisraGries) TopK() map[string]int {
	return mg.counter
}

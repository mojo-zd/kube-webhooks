package pkg

type Patch struct {
	OP    string      `json:"op"`
	Path  string      `json:"path"`
	Value interface{} `json:"value"`
}

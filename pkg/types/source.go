package types

type Secrets struct {
	Secret      string `json:"secret"`
	AccessToken string `json:"accessToken"`
}

type Resources struct {
	Id string `json:"id"`
}

type Credential struct {
	Secrets   Secrets   `json:"secrets"`
	Resources Resources `json:"resources"`
}

type SourceReference struct {
	Id     string `json:"id"`
	Bucket string `json:"bucket"`
	Prefix string `json:"prefix"`
	Region string `json:"region"`
}

type Input struct {
	StorageType string          `json:"storageType"`
	Reference   SourceReference `json:"reference"`
	Credential  Credential      `json:"credential"`
}

type Output struct {
	StorageType string          `json:"storageType"`
	Reference   SourceReference `json:"reference"`
	Credential  Credential      `json:"credential"`
}

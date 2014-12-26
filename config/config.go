package config

type MatchItem struct {
	Type         string `json:"type"`
	MatchData    string `json:"match_data"`
	ResponseFile string `json:"response_file"`
}

type MockServerConfig struct {
	Host        string      `json:"host"`
	Port        int32       `json:"port"`
	DumpRequest bool        `json:"dump_request"`
	Matches     []MatchItem `json:"matchs"`
}

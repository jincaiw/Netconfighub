package model

type DiffResult struct {
	FromHash     string `json:"from_hash"`
	ToHash       string `json:"to_hash"`
	DiffContent  string `json:"diff_content"`
	AddedLines   int    `json:"added_lines"`
	RemovedLines int    `json:"removed_lines"`
}

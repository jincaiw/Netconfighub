package baseline

import (
	"bytes"
	"context"
	"fmt"
	"strings"

	"github.com/netconfighub/netconfighub/internal/gitstore"
	"github.com/netconfighub/netconfighub/internal/model"
)

type Engine interface {
	Compare(ctx context.Context, baselineContent, currentContent string, baselineID uint) (*CompareResult, error)
	CompareWithLatest(ctx context.Context, device *model.Device, baseline *model.Baseline) (*CompareResult, error)
}

type CompareResult struct {
	HasDeviation bool   `json:"has_deviation"`
	DiffContent  string `json:"diff_content"`
	AddedLines   int    `json:"added_lines"`
	RemovedLines int    `json:"removed_lines"`
	BaselineID   uint   `json:"baseline_id"`
}

type baselineEngine struct {
	gitStore *gitstore.Store
}

func NewEngine(gs *gitstore.Store) Engine {
	return &baselineEngine{gitStore: gs}
}

func (e *baselineEngine) Compare(ctx context.Context, baselineContent, currentContent string, baselineID uint) (*CompareResult, error) {
	result := &CompareResult{
		BaselineID: baselineID,
	}

	diffContent, added, removed := generateDiff(baselineContent, currentContent)
	result.DiffContent = diffContent
	result.AddedLines = added
	result.RemovedLines = removed
	result.HasDeviation = added > 0 || removed > 0

	return result, nil
}

func (e *baselineEngine) CompareWithLatest(ctx context.Context, device *model.Device, baseline *model.Baseline) (*CompareResult, error) {
	if e.gitStore == nil {
		return nil, fmt.Errorf("Git 存储未初始化")
	}
	if device == nil {
		return nil, fmt.Errorf("设备不能为空")
	}
	if baseline == nil {
		return nil, fmt.Errorf("基线不能为空")
	}

	groupName := string(device.Vendor)
	configContent, err := e.gitStore.ReadLatestConfig(groupName, device.Name)
	if err != nil {
		return nil, fmt.Errorf("读取设备最新配置失败: %w", err)
	}

	return e.Compare(ctx, baseline.Content, configContent, baseline.ID)
}

func generateDiff(oldText, newText string) (string, int, int) {
	if oldText == newText {
		return "", 0, 0
	}

	oldLines := splitLines(oldText)
	newLines := splitLines(newText)

	lcs := computeLCS(oldLines, newLines)

	var buf bytes.Buffer
	buf.WriteString("--- 基线\n")
	buf.WriteString("+++ 当前配置\n")

	added := 0
	removed := 0

	type diffEntry struct {
		line string
		op   byte
	}

	entries := backtrackLCS(oldLines, newLines, lcs, &added, &removed)

	oldStart := 1
	newStart := 1
	hunkLines := make([]string, 0, len(entries))

	for _, e := range entries {
		switch e.op {
		case ' ':
			hunkLines = append(hunkLines, " "+e.line)
			oldStart++
			newStart++
		case '-':
			hunkLines = append(hunkLines, "-"+e.line)
			oldStart++
		case '+':
			hunkLines = append(hunkLines, "+"+e.line)
			newStart++
		}
	}

	oldCount := 0
	newCount := 0
	for _, e := range entries {
		switch e.op {
		case ' ', '-':
			oldCount++
		}
		switch e.op {
		case ' ', '+':
			newCount++
		}
	}

	fmt.Fprintf(&buf, "@@ -1,%d +1,%d @@\n", oldCount, newCount)
	for _, l := range hunkLines {
		buf.WriteString(l + "\n")
	}

	return buf.String(), added, removed
}

func splitLines(text string) []string {
	if text == "" {
		return nil
	}
	lines := strings.Split(text, "\n")
	if len(lines) > 0 && lines[len(lines)-1] == "" {
		lines = lines[:len(lines)-1]
	}
	return lines
}

func computeLCS(a, b []string) [][]int {
	m := len(a)
	n := len(b)
	dp := make([][]int, m+1)
	for i := range dp {
		dp[i] = make([]int, n+1)
	}

	for i := 1; i <= m; i++ {
		for j := 1; j <= n; j++ {
			if a[i-1] == b[j-1] {
				dp[i][j] = dp[i-1][j-1] + 1
			} else if dp[i-1][j] >= dp[i][j-1] {
				dp[i][j] = dp[i-1][j]
			} else {
				dp[i][j] = dp[i][j-1]
			}
		}
	}

	return dp
}

type diffEntry struct {
	line string
	op   byte
}

func backtrackLCS(oldLines, newLines []string, dp [][]int, added, removed *int) []diffEntry {
	var entries []diffEntry
	i := len(oldLines)
	j := len(newLines)

	for i > 0 || j > 0 {
		if i > 0 && j > 0 && oldLines[i-1] == newLines[j-1] {
			entries = append(entries, diffEntry{line: oldLines[i-1], op: ' '})
			i--
			j--
		} else if j > 0 && (i == 0 || dp[i][j-1] >= dp[i-1][j]) {
			entries = append(entries, diffEntry{line: newLines[j-1], op: '+'})
			j--
			*added++
		} else {
			entries = append(entries, diffEntry{line: oldLines[i-1], op: '-'})
			i--
			*removed++
		}
	}

	for k := 0; k < len(entries)/2; k++ {
		entries[k], entries[len(entries)-1-k] = entries[len(entries)-1-k], entries[k]
	}

	return entries
}

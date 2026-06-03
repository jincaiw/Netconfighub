package baseline

import (
	"context"
	"strings"
	"testing"
)

func TestCompareIdenticalContent(t *testing.T) {
	engine := NewEngine(nil)
	content := "line1\nline2\nline3\n"

	result, err := engine.Compare(context.Background(), content, content, 1)
	if err != nil {
		t.Fatalf("比对失败: %v", err)
	}

	if result.HasDeviation {
		t.Error("相同内容不应有偏差")
	}
	if result.AddedLines != 0 {
		t.Errorf("期望新增行数为 0，实际为 %d", result.AddedLines)
	}
	if result.RemovedLines != 0 {
		t.Errorf("期望删除行数为 0，实际为 %d", result.RemovedLines)
	}
	if result.DiffContent != "" {
		t.Errorf("相同内容不应有差异输出，实际为: %s", result.DiffContent)
	}
	if result.BaselineID != 1 {
		t.Errorf("期望 BaselineID 为 1，实际为 %d", result.BaselineID)
	}
}

func TestCompareDifferentContent(t *testing.T) {
	engine := NewEngine(nil)
	baseline := "line1\nline2\nline3\n"
	current := "line1\nline2_modified\nline3\n"

	result, err := engine.Compare(context.Background(), baseline, current, 1)
	if err != nil {
		t.Fatalf("比对失败: %v", err)
	}

	if !result.HasDeviation {
		t.Error("不同内容应有偏差")
	}
	if result.DiffContent == "" {
		t.Error("不同内容应有差异输出")
	}
	if !strings.Contains(result.DiffContent, "--- 基线") {
		t.Error("差异输出应包含 --- 基线")
	}
	if !strings.Contains(result.DiffContent, "+++ 当前配置") {
		t.Error("差异输出应包含 +++ 当前配置")
	}
	if !strings.Contains(result.DiffContent, "@@") {
		t.Error("差异输出应包含 @@ hunk 标记")
	}
}

func TestCompareAddedLines(t *testing.T) {
	engine := NewEngine(nil)
	baseline := "line1\nline2\n"
	current := "line1\nline2\nline3\nline4\n"

	result, err := engine.Compare(context.Background(), baseline, current, 1)
	if err != nil {
		t.Fatalf("比对失败: %v", err)
	}

	if !result.HasDeviation {
		t.Error("有新增行应有偏差")
	}
	if result.AddedLines != 2 {
		t.Errorf("期望新增行数为 2，实际为 %d", result.AddedLines)
	}
	if result.RemovedLines != 0 {
		t.Errorf("期望删除行数为 0，实际为 %d", result.RemovedLines)
	}
}

func TestCompareRemovedLines(t *testing.T) {
	engine := NewEngine(nil)
	baseline := "line1\nline2\nline3\nline4\n"
	current := "line1\nline2\n"

	result, err := engine.Compare(context.Background(), baseline, current, 1)
	if err != nil {
		t.Fatalf("比对失败: %v", err)
	}

	if !result.HasDeviation {
		t.Error("有删除行应有偏差")
	}
	if result.RemovedLines != 2 {
		t.Errorf("期望删除行数为 2，实际为 %d", result.RemovedLines)
	}
	if result.AddedLines != 0 {
		t.Errorf("期望新增行数为 0，实际为 %d", result.AddedLines)
	}
}

func TestCompareMixedChanges(t *testing.T) {
	engine := NewEngine(nil)
	baseline := "line1\nline2\nline3\n"
	current := "line1\nline2_modified\nline3\nline4\n"

	result, err := engine.Compare(context.Background(), baseline, current, 5)
	if err != nil {
		t.Fatalf("比对失败: %v", err)
	}

	if !result.HasDeviation {
		t.Error("有变更应有偏差")
	}
	if result.AddedLines < 1 {
		t.Errorf("期望至少有 1 行新增，实际为 %d", result.AddedLines)
	}
	if result.RemovedLines < 1 {
		t.Errorf("期望至少有 1 行删除，实际为 %d", result.RemovedLines)
	}
	if result.BaselineID != 5 {
		t.Errorf("期望 BaselineID 为 5，实际为 %d", result.BaselineID)
	}
}

func TestCompareEmptyBaseline(t *testing.T) {
	engine := NewEngine(nil)
	baseline := ""
	current := "line1\nline2\n"

	result, err := engine.Compare(context.Background(), baseline, current, 1)
	if err != nil {
		t.Fatalf("比对失败: %v", err)
	}

	if !result.HasDeviation {
		t.Error("空基线与非空配置比对应有偏差")
	}
	if result.AddedLines != 2 {
		t.Errorf("期望新增行数为 2，实际为 %d", result.AddedLines)
	}
}

func TestCompareEmptyCurrent(t *testing.T) {
	engine := NewEngine(nil)
	baseline := "line1\nline2\n"
	current := ""

	result, err := engine.Compare(context.Background(), baseline, current, 1)
	if err != nil {
		t.Fatalf("比对失败: %v", err)
	}

	if !result.HasDeviation {
		t.Error("非空基线与空配置比对应有偏差")
	}
	if result.RemovedLines != 2 {
		t.Errorf("期望删除行数为 2，实际为 %d", result.RemovedLines)
	}
}

func TestCompareBothEmpty(t *testing.T) {
	engine := NewEngine(nil)

	result, err := engine.Compare(context.Background(), "", "", 1)
	if err != nil {
		t.Fatalf("比对失败: %v", err)
	}

	if result.HasDeviation {
		t.Error("两个空内容不应有偏差")
	}
}

func TestDiffOutputFormat(t *testing.T) {
	engine := NewEngine(nil)
	baseline := "hostname Router1\ninterface Gig0/1\n ip address 10.0.0.1\n"
	current := "hostname Router1\ninterface Gig0/1\n ip address 10.0.0.2\n shutdown\n"

	result, err := engine.Compare(context.Background(), baseline, current, 1)
	if err != nil {
		t.Fatalf("比对失败: %v", err)
	}

	if !result.HasDeviation {
		t.Error("应有偏差")
	}

	lines := strings.Split(result.DiffContent, "\n")
	hasMinus := false
	hasPlus := false
	for _, line := range lines {
		if strings.HasPrefix(line, "-") && !strings.HasPrefix(line, "---") {
			hasMinus = true
		}
		if strings.HasPrefix(line, "+") && !strings.HasPrefix(line, "+++") {
			hasPlus = true
		}
	}
	if !hasMinus {
		t.Error("差异输出应包含 - 行（删除行）")
	}
	if !hasPlus {
		t.Error("差异输出应包含 + 行（新增行）")
	}
}

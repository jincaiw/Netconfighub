package gitstore

import (
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/ProtonMail/go-crypto/openpgp"
	"github.com/ProtonMail/go-crypto/openpgp/armor"
	"github.com/ProtonMail/go-crypto/openpgp/packet"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/object"
)

type writeRequest struct {
	groupName  string
	deviceName string
	content    string
	resultChan chan writeResult
}

type writeResult struct {
	configHash string
	err        error
}

type DiffResult struct {
	FromHash     string `json:"from_hash"`
	ToHash       string `json:"to_hash"`
	DiffContent  string `json:"diff_content"`
	AddedLines   int    `json:"added_lines"`
	RemovedLines int    `json:"removed_lines"`
}

type Store struct {
	repoPath        string
	authorName      string
	authorEmail     string
	repo            *git.Repository
	signKey         *openpgp.Entity
	gitCryptEnabled bool
	gitCryptKeyFile string
	writeQueue      chan writeRequest
	wg              sync.WaitGroup
	cancel          context.CancelFunc
}

func NewStore(repoPath, authorName, authorEmail, gpgSignKey string, gitCryptEnabled bool, gitCryptKeyFile string) (*Store, error) {
	s := &Store{
		repoPath:        repoPath,
		authorName:      authorName,
		authorEmail:     authorEmail,
		gitCryptEnabled: gitCryptEnabled,
		gitCryptKeyFile: gitCryptKeyFile,
		writeQueue:      make(chan writeRequest, 256),
	}

	if err := os.MkdirAll(repoPath, 0755); err != nil {
		return nil, fmt.Errorf("创建 Git 仓库目录失败: %w", err)
	}

	repo, err := git.PlainOpen(repoPath)
	if err != nil {
		if err == git.ErrRepositoryNotExists {
			repo, err = git.PlainInit(repoPath, false)
			if err != nil {
				return nil, fmt.Errorf("初始化 Git 仓库失败: %w", err)
			}
		} else {
			return nil, fmt.Errorf("打开 Git 仓库失败: %w", err)
		}
	}

	s.repo = repo

	if gpgSignKey != "" {
		entity, err := loadSignKey(gpgSignKey)
		if err != nil {
			return nil, fmt.Errorf("加载 GPG 签名密钥失败: %w", err)
		}
		s.signKey = entity
	}

	return s, nil
}

func (s *Store) Start() {
	ctx, cancel := context.WithCancel(context.Background())
	s.cancel = cancel
	s.wg.Add(1)
	go s.writeLoop(ctx)
}

func (s *Store) Stop() {
	if s.cancel != nil {
		s.cancel()
	}
	close(s.writeQueue)
	s.wg.Wait()
}

func (s *Store) writeLoop(_ context.Context) {
	defer s.wg.Done()
	for req := range s.writeQueue {
		configHash, err := s.WriteConfig(req.groupName, req.deviceName, req.content)
		req.resultChan <- writeResult{configHash: configHash, err: err}
	}
}

func (s *Store) WriteConfigAsync(groupName, deviceName, content string) (string, error) {
	if s.cancel == nil {
		return s.WriteConfig(groupName, deviceName, content)
	}
	req := writeRequest{
		groupName:  groupName,
		deviceName: deviceName,
		content:    content,
		resultChan: make(chan writeResult, 1),
	}
	s.writeQueue <- req
	result := <-req.resultChan
	return result.configHash, result.err
}

func loadSignKey(keyPath string) (*openpgp.Entity, error) {
	f, err := os.Open(keyPath)
	if err != nil {
		return nil, fmt.Errorf("打开 GPG 密钥文件失败: %w", err)
	}
	defer f.Close()

	block, err := armor.Decode(f)
	if err != nil {
		f.Seek(0, io.SeekStart)
		entities, err := openpgp.ReadKeyRing(f)
		if err != nil {
			return nil, fmt.Errorf("解析 GPG 密钥失败: %w", err)
		}
		if len(entities) == 0 {
			return nil, fmt.Errorf("GPG 密钥文件中未找到密钥")
		}
		return entities[0], nil
	}

	entity, err := openpgp.ReadEntity(packet.NewReader(block.Body))
	if err != nil {
		return nil, fmt.Errorf("读取 GPG 密钥实体失败: %w", err)
	}

	return entity, nil
}

func (s *Store) gitCryptUnlock() error {
	if !s.gitCryptEnabled || s.gitCryptKeyFile == "" {
		return nil
	}
	cmd := exec.Command("git-crypt", "unlock", s.gitCryptKeyFile)
	cmd.Dir = s.repoPath
	if output, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("git-crypt unlock 失败: %s: %w", string(output), err)
	}
	return nil
}

func (s *Store) gitCryptLock() error {
	if !s.gitCryptEnabled {
		return nil
	}
	cmd := exec.Command("git-crypt", "lock")
	cmd.Dir = s.repoPath
	if output, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("git-crypt lock 失败: %s: %w", string(output), err)
	}
	return nil
}

func (s *Store) WriteConfig(groupName, deviceName, content string) (string, error) {
	groupName, deviceName, err := safeConfigPathComponents(groupName, deviceName)
	if err != nil {
		return "", err
	}

	if err := s.gitCryptUnlock(); err != nil {
		return "", err
	}

	filePath := filepath.Join(groupName, deviceName, "running-config.txt")
	fullPath := filepath.Join(s.repoPath, filePath)

	dir := filepath.Dir(fullPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		s.gitCryptLock()
		return "", fmt.Errorf("创建配置目录失败: %w", err)
	}

	hash := sha256.Sum256([]byte(content))
	configHash := hex.EncodeToString(hash[:])

	existingContent, err := s.ReadLatestConfig(groupName, deviceName)
	if err == nil && existingContent == content {
		s.gitCryptLock()
		return "", nil
	}

	if err := os.WriteFile(fullPath, []byte(content), 0644); err != nil {
		s.gitCryptLock()
		return "", fmt.Errorf("写入配置文件失败: %w", err)
	}

	worktree, err := s.repo.Worktree()
	if err != nil {
		s.gitCryptLock()
		return configHash, fmt.Errorf("获取工作树失败: %w", err)
	}

	_, err = worktree.Add(filePath)
	if err != nil {
		s.gitCryptLock()
		return configHash, fmt.Errorf("暂存文件失败: %w", err)
	}

	timestamp := time.Now().Format(time.RFC3339)
	commitMsg := fmt.Sprintf("backup: %s/%s - %s", groupName, deviceName, timestamp)

	commitOpts := &git.CommitOptions{
		Author: &object.Signature{
			Name:  s.authorName,
			Email: s.authorEmail,
			When:  time.Now(),
		},
	}
	if s.signKey != nil {
		commitOpts.SignKey = s.signKey
	}

	_, err = worktree.Commit(commitMsg, commitOpts)
	if err != nil {
		if err.Error() == "nothing to commit" {
			s.gitCryptLock()
			return "", nil
		}
		s.gitCryptLock()
		return configHash, fmt.Errorf("提交变更失败: %w", err)
	}

	s.gitCryptLock()
	return configHash, nil
}

func (s *Store) ConfigHistory(groupName, deviceName string, limit int) ([]*object.Commit, error) {
	groupName, deviceName, err := safeConfigPathComponents(groupName, deviceName)
	if err != nil {
		return nil, err
	}

	filePath := filepath.Join(groupName, deviceName, "running-config.txt")

	logOpts := &git.LogOptions{
		PathFilter: func(path string) bool {
			return path == filePath
		},
	}

	refs, err := s.repo.Log(logOpts)
	if err != nil {
		return nil, fmt.Errorf("获取提交日志失败: %w", err)
	}

	var commits []*object.Commit
	count := 0
	err = refs.ForEach(func(c *object.Commit) error {
		if limit > 0 && count >= limit {
			return fmt.Errorf("stop iteration")
		}
		commits = append(commits, c)
		count++
		return nil
	})

	if err != nil && err.Error() != "stop iteration" {
		return nil, fmt.Errorf("遍历提交日志失败: %w", err)
	}

	return commits, nil
}

func (s *Store) ReadConfigAtRevision(groupName, deviceName, hash string) (string, error) {
	groupName, deviceName, err := safeConfigPathComponents(groupName, deviceName)
	if err != nil {
		return "", err
	}

	if err := s.gitCryptUnlock(); err != nil {
		return "", err
	}
	defer s.gitCryptLock()

	filePath := filepath.Join(groupName, deviceName, "running-config.txt")

	rev, err := s.repo.ResolveRevision(plumbing.Revision(hash))
	if err != nil {
		return "", fmt.Errorf("解析版本号失败: %w", err)
	}

	commit, err := s.repo.CommitObject(*rev)
	if err != nil {
		return "", fmt.Errorf("获取提交对象失败: %w", err)
	}

	tree, err := commit.Tree()
	if err != nil {
		return "", fmt.Errorf("获取文件树失败: %w", err)
	}

	file, err := tree.File(filePath)
	if err != nil {
		return "", fmt.Errorf("获取文件内容失败: %w", err)
	}

	content, err := file.Contents()
	if err != nil {
		return "", fmt.Errorf("读取文件内容失败: %w", err)
	}

	return content, nil
}

func (s *Store) ReadLatestConfig(groupName, deviceName string) (string, error) {
	groupName, deviceName, err := safeConfigPathComponents(groupName, deviceName)
	if err != nil {
		return "", err
	}

	if err := s.gitCryptUnlock(); err != nil {
		return "", err
	}
	defer s.gitCryptLock()

	filePath := filepath.Join(groupName, deviceName, "running-config.txt")
	fullPath := filepath.Join(s.repoPath, filePath)

	data, err := os.ReadFile(fullPath)
	if err != nil {
		if os.IsNotExist(err) {
			return "", fmt.Errorf("配置文件不存在")
		}
		return "", fmt.Errorf("读取配置文件失败: %w", err)
	}

	return string(data), nil
}

func (s *Store) ListDeviceConfigs() ([]DeviceInfo, error) {
	var devices []DeviceInfo

	groupDirs, err := os.ReadDir(s.repoPath)
	if err != nil {
		return nil, fmt.Errorf("读取仓库目录失败: %w", err)
	}

	for _, groupEntry := range groupDirs {
		if !groupEntry.IsDir() || groupEntry.Name() == ".git" {
			continue
		}

		groupName := groupEntry.Name()
		deviceDirPath := filepath.Join(s.repoPath, groupName)

		deviceEntries, err := os.ReadDir(deviceDirPath)
		if err != nil {
			continue
		}

		for _, deviceEntry := range deviceEntries {
			if !deviceEntry.IsDir() {
				continue
			}

			deviceName := deviceEntry.Name()
			configFile := filepath.Join(deviceDirPath, deviceName, "running-config.txt")

			if _, err := os.Stat(configFile); err == nil {
				devices = append(devices, DeviceInfo{
					GroupName:  groupName,
					DeviceName: deviceName,
				})
			}
		}
	}

	return devices, nil
}

func safeConfigPathComponents(groupName, deviceName string) (string, string, error) {
	groupName = strings.TrimSpace(groupName)
	deviceName = strings.TrimSpace(deviceName)
	if groupName == "" {
		groupName = "ungrouped"
	}
	if err := safePathComponent("分组名称", groupName); err != nil {
		return "", "", err
	}
	if err := safePathComponent("设备名称", deviceName); err != nil {
		return "", "", err
	}
	return groupName, deviceName, nil
}

func safePathComponent(label, value string) error {
	if value == "" {
		return fmt.Errorf("%s不能为空", label)
	}
	if value == "." || value == ".." {
		return fmt.Errorf("%s不能为 . 或 ..", label)
	}
	for _, r := range value {
		if r < 0x20 || r == 0x7f {
			return fmt.Errorf("%s不能包含控制字符", label)
		}
		switch r {
		case '/', '\\', ':':
			return fmt.Errorf("%s不能包含路径分隔符或冒号", label)
		}
	}
	return nil
}

func (s *Store) Push(remote, branch string) error {
	cmd := exec.Command("git", "push", remote, branch)
	cmd.Dir = s.repoPath
	if output, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("Git push 失败: %s: %w", string(output), err)
	}
	return nil
}

func (s *Store) DiffConfigs(groupName, deviceName, fromHash, toHash string) (*DiffResult, error) {
	fromContent, err := s.ReadConfigAtRevision(groupName, deviceName, fromHash)
	if err != nil {
		return nil, fmt.Errorf("读取源版本配置失败: %w", err)
	}
	toContent, err := s.ReadConfigAtRevision(groupName, deviceName, toHash)
	if err != nil {
		return nil, fmt.Errorf("读取目标版本配置失败: %w", err)
	}
	diffContent, added, removed := generateDiff(fromContent, toContent)
	return &DiffResult{
		FromHash:     fromHash,
		ToHash:       toHash,
		DiffContent:  diffContent,
		AddedLines:   added,
		RemovedLines: removed,
	}, nil
}

type DeviceInfo struct {
	GroupName  string `json:"group_name"`
	DeviceName string `json:"device_name"`
}

func generateDiff(oldText, newText string) (string, int, int) {
	if oldText == newText {
		return "", 0, 0
	}
	oldLines := splitLines(oldText)
	newLines := splitLines(newText)
	lcs := computeLCS(oldLines, newLines)
	added := 0
	removed := 0
	entries := backtrackLCS(oldLines, newLines, lcs, &added, &removed)
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
	var buf bytes.Buffer
	fmt.Fprintf(&buf, "@@ -1,%d +1,%d @@\n", oldCount, newCount)
	for _, e := range entries {
		switch e.op {
		case ' ':
			buf.WriteString(" " + e.line + "\n")
		case '-':
			buf.WriteString("-" + e.line + "\n")
		case '+':
			buf.WriteString("+" + e.line + "\n")
		}
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

//go:build windows

package main

import (
	"strings"
	"sync"
)

type viewState struct {
	Ready     bool         `json:"ready"`
	Status    string       `json:"status"`
	Busy      bool         `json:"busy"`
	Bots      []botItem    `json:"bots"`
	Selected  string       `json:"selected"`
	Checks    []checkRow   `json:"checks"`
	Memories  []memoryItem `json:"memories"`
	Source    string       `json:"source"`
	Note      string       `json:"note"`
	AllPassed bool         `json:"allPassed"`
}

type appModel struct {
	mu       sync.Mutex
	report   health
	selected string
	memories []memoryItem
	source   string
	note     string
	status   string
	busy     bool
}

func (m *appModel) snapshot() viewState {
	m.mu.Lock()
	defer m.mu.Unlock()
	passed := len(m.report.Rows) > 0
	for _, row := range m.report.Rows {
		if !row.OK {
			passed = false
		}
	}
	memories := append([]memoryItem(nil), m.memories...)
	bots := append([]botItem(nil), m.report.Bots...)
	checks := append([]checkRow(nil), m.report.Rows...)
	return viewState{
		Ready: m.report.Ready, Status: m.status, Busy: m.busy,
		Bots: bots, Selected: m.selected, Checks: checks,
		Memories: memories, Source: m.source, Note: m.note, AllPassed: passed,
	}
}

func (m *appModel) refresh() viewState {
	m.mu.Lock()
	m.busy = true
	m.status = "正在检测…"
	m.mu.Unlock()
	report := detect()
	m.mu.Lock()
	prev := m.selected
	m.report = report
	m.busy = false
	if !report.Ready {
		m.status = "还不能用。看左边的检查。"
		m.memories = nil
		m.source = ""
	} else if m.status == "正在检测…" || m.status == "" {
		m.status = "可以读写记忆。"
	}
	if prev == "" && len(report.Bots) > 0 {
		m.selected = report.Bots[0].ID
	}
	m.mu.Unlock()
	return m.snapshot()
}

func (m *appModel) selectBot(id string) viewState {
	m.mu.Lock()
	if m.selected != "" && m.selected != id {
		m.memories = nil
		m.source = ""
		m.note = ""
		m.status = "换了 Bot。拉取之后显示这一份。"
	}
	m.selected = id
	m.mu.Unlock()
	return m.snapshot()
}

func (m *appModel) pull() viewState {
	m.mu.Lock()
	if m.busy || !m.report.Ready || m.selected == "" || m.report.Session == nil {
		m.mu.Unlock()
		return m.snapshot()
	}
	sess := *m.report.Session
	id := m.selected
	m.busy = true
	m.status = "正在让 Bot 导出记忆…"
	m.mu.Unlock()

	items, source, note, err := pullMemories(&sess, id)
	m.mu.Lock()
	m.busy = false
	if m.selected == id {
		if err != nil {
			m.status = err.Error()
		} else {
			m.memories = items
			m.source = source
			m.note = note
			m.status = "已拉取 " + itoa(len(items)) + " 条。"
			if note != "" {
				m.status += note
			}
		}
	}
	m.mu.Unlock()
	return m.snapshot()
}

func (m *appModel) write(text, scope, tier string) viewState {
	m.mu.Lock()
	if m.busy || !m.report.Ready || m.selected == "" || m.report.Session == nil {
		m.mu.Unlock()
		return m.snapshot()
	}
	sess := *m.report.Session
	id := m.selected
	m.busy = true
	m.status = "正在让 Bot 写入…"
	m.mu.Unlock()

	err := writeMemory(&sess, id, text, tier, scope)
	m.mu.Lock()
	m.busy = false
	if err != nil {
		m.status = err.Error()
	} else if m.selected == id {
		place := "账号共用"
		if scope == "agent" {
			place = "仅此 Bot"
		}
		kind := "普通事实"
		if tier == "profile" {
			kind = "指令"
		}
		m.status = "已写入：" + place + " · " + kind + "。"
		m.memories = append([]memoryItem{decorate(memoryItem{Text: strings.TrimSpace(text), Scope: scope, Tier: tier})}, m.memories...)
	}
	m.mu.Unlock()
	return m.snapshot()
}

func (m *appModel) forget(text string) viewState {
	m.mu.Lock()
	if m.busy || !m.report.Ready || m.selected == "" || m.report.Session == nil {
		m.mu.Unlock()
		return m.snapshot()
	}
	sess := *m.report.Session
	id := m.selected
	m.busy = true
	m.status = "正在让 Bot 忘掉这一条…"
	m.mu.Unlock()

	err := forgetMemory(&sess, id, text)
	m.mu.Lock()
	m.busy = false
	if err != nil {
		m.status = err.Error()
	} else if m.selected == id {
		next := make([]memoryItem, 0, len(m.memories))
		removed := false
		for _, item := range m.memories {
			if !removed && item.Text == text {
				removed = true
				continue
			}
			next = append(next, item)
		}
		m.memories = next
		m.status = "已忘掉。"
	}
	m.mu.Unlock()
	return m.snapshot()
}

func itoa(n int) string {
	if n == 0 {
		return "0"
	}
	var b [16]byte
	i := len(b)
	for n > 0 {
		i--
		b[i] = byte('0' + n%10)
		n /= 10
	}
	return string(b[i:])
}

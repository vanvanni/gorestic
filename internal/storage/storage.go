package storage

import (
	"encoding/json"
	"os"
	"sync"
)

type Manager struct {
	path    string
	storage *Storage
	mu      sync.RWMutex
}

func NewManager(path string) (*Manager, error) {
	manager := &Manager{
		path:    path,
		storage: &Storage{Stats: []BackupStats{}},
	}

	if err := manager.load(); err != nil {
		return nil, err
	}

	return manager, nil
}

func (m *Manager) load() error {
	m.mu.Lock()
	defer m.mu.Unlock()

	data, err := os.ReadFile(m.path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return err
	}

	return json.Unmarshal(data, &m.storage)
}

func (m *Manager) save() error {
	m.mu.Lock()
	defer m.mu.Unlock()

	data, err := json.MarshalIndent(m.storage, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(m.path, data, 0644)
}

func (m *Manager) AddStats(stats BackupStats) error {
	m.mu.Lock()
	m.storage.Stats = append(m.storage.Stats, stats)
	m.mu.Unlock()

	return m.save()
}

func (m *Manager) GetAllStats() []BackupStats {
	m.mu.RLock()
	defer m.mu.RUnlock()

	stats := make([]BackupStats, len(m.storage.Stats))
	copy(stats, m.storage.Stats)
	return stats
}

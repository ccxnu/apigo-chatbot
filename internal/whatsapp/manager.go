package whatsapp

import (
	"context"
	"fmt"
	"sync"
)

// Manager manages the global WhatsApp service instance
type Manager struct {
	mu      sync.RWMutex
	service *Service
}

var globalManager = &Manager{}

// GetManager returns the global WhatsApp service manager
func GetManager() *Manager {
	return globalManager
}

// SetService sets the current WhatsApp service
func (m *Manager) SetService(service *Service) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.service = service
}

// GetService returns the current WhatsApp service (may be nil)
func (m *Manager) GetService() *Service {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.service
}

// GetCurrentQR returns the current QR code from the active service
func (m *Manager) GetCurrentQR() string {
	m.mu.RLock()
	defer m.mu.RUnlock()
	if m.service == nil {
		return ""
	}
	return m.service.GetCurrentQR()
}

// IsConnected returns whether the WhatsApp service is connected
func (m *Manager) IsConnected() bool {
	m.mu.RLock()
	defer m.mu.RUnlock()
	if m.service == nil {
		return false
	}
	return m.service.IsConnected()
}

// Logout logs out from WhatsApp and clears device pairing
func (m *Manager) Logout(ctx context.Context) error {
	m.mu.RLock()
	defer m.mu.RUnlock()
	if m.service == nil {
		return fmt.Errorf("WhatsApp service not initialized")
	}
	return m.service.Logout(ctx)
}

// Reconnect disconnects and reconnects to generate new QR code
func (m *Manager) Reconnect(ctx context.Context) error {
	m.mu.RLock()
	defer m.mu.RUnlock()
	if m.service == nil {
		return fmt.Errorf("WhatsApp service not initialized")
	}
	return m.service.Reconnect(ctx)
}

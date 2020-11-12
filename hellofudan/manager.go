package hellofudan

import (
	"log"
	"sync"
)

// Manager control the hello fudan goroutine
type Manager struct {
	wg      sync.WaitGroup
	stuList []Student
}

// NewManager return a hello fudan manager
func NewManager(students []Student) *Manager {
	return &Manager{
		wg:      sync.WaitGroup{},
		stuList: students,
	}
}

// Start Manager
func (m *Manager) Start() {
	helloFudan := `
	 _   _      _ _         _____          _
	| | | | ___| | | ___   |  ___|   _  __| | __ _ _ __
	| |_| |/ _ \ | |/ _ \  | |_ | | | |/ _  |/ _  |  _ \
	|  _  |  __/ | | (_) | |  _|| |_| | (_| | (_| | | | |
	|_| |_|\___|_|_|\___/  |_|   \__,_|\__,_|\__,_|_| |_|
	`
	log.Println(helloFudan)
	for _, stu := range m.stuList {
		m.wg.Add(1)
		go m.start(stu)
	}

	m.wg.Wait()
}

func (m *Manager) start(stu Student) {

	defer m.wg.Done()

	hf := newHelloFudan(stu)

	hf.login()
	if !hf.checkStatus() {
		hf.checkIn()
	}
	hf.logout()

}

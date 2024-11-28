// pkg/spinner/spinner.go
package spinner

import (
	"fmt"
	"sync"
	"time"

	"github.com/briandowns/spinner"
)

type CLI struct {}

type Spinner struct {
	message string
	frames  []string
	stop    chan struct{}
	wg      sync.WaitGroup
}

func New(message string) *Spinner {
	return &Spinner{
		message: message,
		frames:  []string{"⠋", "⠙", "⠹", "⠸", "⠼", "⠴", "⠦", "⠧", "⠇", "⠏"},
		stop:    make(chan struct{}),
	}
}

func (s *Spinner) Start() {
	s.wg.Add(1)
	go func() {
		defer s.wg.Done()
		for i := 0; ; i = (i + 1) % len(s.frames) {
			select {
			case <-s.stop:
				return
			default:
				fmt.Printf("\r%s %s", s.frames[i], s.message)
				time.Sleep(100 * time.Millisecond)
			}
		}
	}()
}

func (s *Spinner) Stop() {
	close(s.stop)
	s.wg.Wait()
	fmt.Printf("\r\033[K") // Clear the line
}

func NewSpinner(message string) *spinner.Spinner {
	s := spinner.New(spinner.CharSets[9], 100*time.Millisecond)
	s.Suffix = " " + message
	return s
}

func RunWithSpinner(cli *CLI, provider string, prompt string, noStream bool, fn func() error) error {
	s := NewSpinner("Thinking...")
	s.Start()
	defer s.Stop()

	err := fn()
	if err != nil {
		return fmt.Errorf("error during execution: %w", err)
	}

	return nil
}

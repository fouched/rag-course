package ingest

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"rag-course/llm"
	"rag-course/vector"
	"strings"
	"sync"
	"time"

	"github.com/fsnotify/fsnotify"
)

const debounceDelay = 500 * time.Millisecond

func Watch(ctx context.Context, opts Options, embedder llm.Embedder, store vector.Store, logger *log.Logger) error {
	if filepath.Clean(opts.SourceDir) == filepath.Clean(opts.ProcessedDir) {
		return errors.New("source and processed directories must differ")
	}
	if err := os.MkdirAll(opts.SourceDir, 0o755); err != nil {
		return fmt.Errorf("could not create source directory: %w", err)
	}
	if err := os.MkdirAll(opts.ProcessedDir, 0o755); err != nil {
		return fmt.Errorf("could not create processed directory: %w", err)
	}

	w, err := fsnotify.NewWatcher()
	if err != nil {
		return fmt.Errorf("could not create watcher: %w", err)
	}
	defer w.Close()

	if err := w.Add(opts.SourceDir); err != nil {
		return fmt.Errorf("could not add source directory to watcher: %w", err)
	}

	processedAbs, err := filepath.Abs(opts.ProcessedDir)
	if err != nil {
		return fmt.Errorf("could not resolve processed directory: %w", err)
	}

	handle := func(path string) {
		if err := processOne(ctx, path, opts, embedder, store); err != nil {
			logger.Printf("could not process %s: %v", filepath.Base(path), err)
		}
		dst := filepath.Join(opts.ProcessedDir, filepath.Base(path))
		if err := os.Rename(path, dst); err != nil {
			logger.Printf("could not move %s to processed: %v", filepath.Base(path), err)
		}
		logger.Printf("ingested %s", filepath.Base(path))
	}

	entries, err := os.ReadDir(opts.SourceDir)
	if err != nil {
		return fmt.Errorf("could not read source dir: %w", err)
	}

	for _, e := range entries {
		if e.IsDir() || strings.HasPrefix(e.Name(), ".") {
			continue
		}
		go handle(filepath.Join(opts.SourceDir, e.Name()))
	}

	var (
		timerMu sync.Mutex
		timers  = make(map[string]*time.Timer)
	)

	schedule := func(path string) {
		timerMu.Lock()
		defer timerMu.Unlock()
		if t, ok := timers[path]; ok {
			t.Reset(debounceDelay)
			return
		}

		timers[path] = time.AfterFunc(debounceDelay, func() {
			timerMu.Lock()
			delete(timers, path)
			timerMu.Unlock()
			handle(path)
		})
	}

	for {
		select {
		case <-ctx.Done():
			return nil
		case ev, ok := <-w.Events:
			if !ok {
				return nil
			}
			if ev.Op&(fsnotify.Create|fsnotify.Write) == 0 {
				continue
			}
			if !shouldProcess(ev.Name, processedAbs) {
				continue
			}
			schedule(ev.Name)
		case err, ok := <-w.Errors:
			if !ok {
				return nil
			}
			logger.Printf("watch error: %v", err)
		}
	}
}

func shouldProcess(path, processedAbs string) bool {
	if strings.HasPrefix(filepath.Base(path), ".") {
		return false
	}
	info, err := os.Stat(path)
	if err != nil || !info.IsDir() {
		return false
	}
	abs, err := filepath.Abs(path)
	if err != nil {
		return false
	}
	if processedAbs != "" && strings.HasPrefix(abs, processedAbs+string(filepath.Separator)) {
		return false
	}
	return true
}

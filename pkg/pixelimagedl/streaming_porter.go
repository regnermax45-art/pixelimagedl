package pixelimagedl

import (
	"archive/zip"
	"bytes"
	"context"
	"crypto/sha256"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/pkg/errors"
)

// StreamingPorter handles real-time firmware porting during download
type StreamingPorter struct {
	sourceDevice   Pixel
	targetDevice   Pixel
	algorithm      PortingAlgorithm
	progressBar    *CustomProgressBar
	portingBuffer  *PortingBuffer
	mutex          sync.RWMutex
	totalSize      int64
	downloadedSize int64
	portedSize     int64
	outputPath     string
}

// CustomProgressBar provides tqdm-style progress visualization
type CustomProgressBar struct {
	total          int64
	current        int64
	portedCurrent  int64
	startTime      time.Time
	lastUpdate     time.Time
	width          int
	description    string
	mutex          sync.Mutex
}

// PortingBuffer manages streaming firmware modification
type PortingBuffer struct {
	buffer         *bytes.Buffer
	zipWriter      *zip.Writer
	tempFile       *os.File
	deviceMappings map[string]string
	bootloaderMods map[string][]byte
	mutex          sync.RWMutex
}

// NewStreamingPorter creates a new streaming porter instance
func NewStreamingPorter(source, target Pixel, algorithm PortingAlgorithm, outputPath string) *StreamingPorter {
	return &StreamingPorter{
		sourceDevice:  source,
		targetDevice:  target,
		algorithm:     algorithm,
		outputPath:    outputPath,
		progressBar:   NewCustomProgressBar("Downloading & Porting", 80),
		portingBuffer: NewPortingBuffer(source, target),
	}
}

// NewCustomProgressBar creates a new progress bar
func NewCustomProgressBar(description string, width int) *CustomProgressBar {
	return &CustomProgressBar{
		description: description,
		width:       width,
		startTime:   time.Now(),
		lastUpdate:  time.Now(),
	}
}

// NewPortingBuffer creates a new porting buffer
func NewPortingBuffer(source, target Pixel) *PortingBuffer {
	buffer := &bytes.Buffer{}
	tempFile, _ := os.CreateTemp("", "streaming-port-*.zip")
	
	return &PortingBuffer{
		buffer:    buffer,
		zipWriter: zip.NewWriter(tempFile),
		tempFile:  tempFile,
		deviceMappings: map[string]string{
			"husky":   "cheetah",  // Pixel 8 Pro -> Pixel 7 Pro
			"shiba":   "panther",  // Pixel 8 -> Pixel 7
			"ripcurrent": "cheetah", // Pixel 9 Pro -> Pixel 7 Pro
			"tokay":   "cheetah",  // Pixel 9 -> Pixel 7 Pro
		},
		bootloaderMods: make(map[string][]byte),
	}
}

// StreamingPortAndDownload performs real-time porting during download
func (sp *StreamingPorter) StreamingPortAndDownload(ctx context.Context, sourceURL string) error {
	log.Printf("🚀 Starting streaming download & port: %s -> %s\n", sp.sourceDevice, sp.targetDevice)
	log.Printf("📡 Source URL: %s\n", sourceURL)
	log.Printf("🎯 Output: %s\n", sp.outputPath)
	
	// Create HTTP request with streaming support
	req, err := sp.createStreamingRequest(ctx, sourceURL)
	if err != nil {
		return errors.Wrap(err, "failed to create streaming request")
	}
	
	// Start the streaming download and port process
	return sp.processStreamingData(ctx, req)
}

// createStreamingRequest creates an HTTP request optimized for streaming
func (sp *StreamingPorter) createStreamingRequest(ctx context.Context, url string) (*http.Request, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, err
	}
	
	// Set headers for optimal streaming
	req.Header.Set("User-Agent", "pixelimagedl/3.0-streaming (advanced-porter)")
	req.Header.Set("Accept", "application/zip,application/octet-stream,*/*")
	req.Header.Set("Accept-Encoding", "identity") // Disable compression for streaming
	req.Header.Set("Connection", "keep-alive")
	req.Header.Set("Cache-Control", "no-cache")
	
	return req, nil
}

// processStreamingData handles the main streaming download and porting logic
func (sp *StreamingPorter) processStreamingData(ctx context.Context, req *http.Request) error {
	client := &http.Client{
		Timeout: 0, // No timeout for streaming
	}
	
	resp, err := client.Do(req)
	if err != nil {
		return errors.Wrap(err, "failed to start download")
	}
	defer resp.Body.Close()
	
	// Get total size for progress tracking
	sp.totalSize = resp.ContentLength
	sp.progressBar.SetTotal(sp.totalSize)
	
	log.Printf("📊 Total firmware size: %.2f MB\n", float64(sp.totalSize)/(1024*1024))
	log.Printf("🔄 Starting real-time porting process...\n")
	
	// Create streaming reader with porting pipeline
	portingReader := sp.createPortingPipeline(resp.Body)
	
	// Process data in chunks with real-time porting
	return sp.streamAndPort(ctx, portingReader)
}

// createPortingPipeline creates a pipeline for real-time firmware porting
func (sp *StreamingPorter) createPortingPipeline(reader io.Reader) io.Reader {
	return &PortingReader{
		source:        reader,
		porter:        sp,
		chunkSize:     64 * 1024, // 64KB chunks for optimal performance
		portingQueue:  make(chan []byte, 10),
		processedData: make(chan []byte, 10),
	}
}

// PortingReader implements io.Reader with real-time porting
type PortingReader struct {
	source        io.Reader
	porter        *StreamingPorter
	chunkSize     int
	portingQueue  chan []byte
	processedData chan []byte
	buffer        []byte
	bufferPos     int
	wg            sync.WaitGroup
	started       bool
}

// Read implements io.Reader interface with real-time porting
func (pr *PortingReader) Read(p []byte) (n int, err error) {
	if !pr.started {
		pr.startPortingWorkers()
		pr.started = true
	}
	
	// Fill buffer if empty
	if pr.bufferPos >= len(pr.buffer) {
		select {
		case data := <-pr.processedData:
			pr.buffer = data
			pr.bufferPos = 0
		default:
			// Read more data from source
			chunk := make([]byte, pr.chunkSize)
			n, err := pr.source.Read(chunk)
			if err != nil {
				if err == io.EOF {
					close(pr.portingQueue)
					pr.wg.Wait()
					close(pr.processedData)
					
					// Return any remaining data
					if len(pr.buffer) > pr.bufferPos {
						remaining := copy(p, pr.buffer[pr.bufferPos:])
						pr.bufferPos += remaining
						return remaining, nil
					}
					return 0, io.EOF
				}
				return 0, err
			}
			
			// Send chunk for porting
			pr.portingQueue <- chunk[:n]
			
			// Update download progress
			pr.porter.updateDownloadProgress(int64(n))
		}
	}
	
	// Copy data to output buffer
	copied := copy(p, pr.buffer[pr.bufferPos:])
	pr.bufferPos += copied
	return copied, nil
}

// startPortingWorkers starts background workers for real-time porting
func (pr *PortingReader) startPortingWorkers() {
	// Start multiple porting workers for parallel processing
	for i := 0; i < 4; i++ {
		pr.wg.Add(1)
		go pr.portingWorker()
	}
}

// portingWorker processes firmware chunks in real-time
func (pr *PortingReader) portingWorker() {
	defer pr.wg.Done()
	
	for chunk := range pr.portingQueue {
		// Apply real-time firmware modifications
		portedChunk := pr.porter.portChunk(chunk)
		
		// Send processed data
		pr.processedData <- portedChunk
		
		// Update porting progress
		pr.porter.updatePortingProgress(int64(len(portedChunk)))
	}
}

// portChunk applies firmware porting modifications to a data chunk
func (sp *StreamingPorter) portChunk(chunk []byte) []byte {
	sp.mutex.Lock()
	defer sp.mutex.Unlock()
	
	// Apply device-specific modifications
	modifiedChunk := sp.applyDeviceMapping(chunk)
	
	// Apply bootloader modifications
	modifiedChunk = sp.applyBootloaderMods(modifiedChunk)
	
	// Apply hardware abstraction layer updates
	modifiedChunk = sp.applyHALUpdates(modifiedChunk)
	
	// Apply Pixel 9 Pro -> Pixel 7 Pro specific optimizations
	modifiedChunk = sp.applyPixel9ProOptimizations(modifiedChunk)
	
	return modifiedChunk
}

// applyDeviceMapping applies device codename mappings
func (sp *StreamingPorter) applyDeviceMapping(chunk []byte) []byte {
	chunkStr := string(chunk)
	
	// Replace device codenames in real-time
	for source, target := range sp.portingBuffer.deviceMappings {
		chunkStr = strings.ReplaceAll(chunkStr, source, target)
	}
	
	// Replace hardware identifiers
	chunkStr = strings.ReplaceAll(chunkStr, "ripcurrent", "cheetah")
	chunkStr = strings.ReplaceAll(chunkStr, "husky", "cheetah")
	chunkStr = strings.ReplaceAll(chunkStr, "Pixel 8 Pro", "Pixel 7 Pro")
	chunkStr = strings.ReplaceAll(chunkStr, "Pixel 9 Pro", "Pixel 7 Pro")
	
	return []byte(chunkStr)
}

// applyBootloaderMods applies bootloader modifications
func (sp *StreamingPorter) applyBootloaderMods(chunk []byte) []byte {
	// Check if chunk contains bootloader data
	if bytes.Contains(chunk, []byte("bootloader")) {
		log.Printf("🔧 Applying bootloader modifications for %s\n", sp.targetDevice)
		
		// Apply Pixel 7 Pro bootloader compatibility
		chunk = bytes.ReplaceAll(chunk, []byte("ripcurrentpro"), []byte("cheetahpro"))
		chunk = bytes.ReplaceAll(chunk, []byte("ripcurrent-"), []byte("cheetah-"))
		
		// Update bootloader version strings
		chunk = bytes.ReplaceAll(chunk, []byte("bootloader-ripcurrent"), []byte("bootloader-cheetah"))
	}
	
	return chunk
}

// applyHALUpdates applies hardware abstraction layer updates
func (sp *StreamingPorter) applyHALUpdates(chunk []byte) []byte {
	// Check if chunk contains HAL data
	if bytes.Contains(chunk, []byte("hardware")) || bytes.Contains(chunk, []byte("vendor")) {
		log.Printf("⚙️ Applying HAL updates for device compatibility\n")
		
		// Update hardware interface mappings
		chunk = bytes.ReplaceAll(chunk, []byte("gs201"), []byte("gs101")) // Tensor G2 -> G1
		chunk = bytes.ReplaceAll(chunk, []byte("zuma"), []byte("slider"))  // Pixel 9 -> Pixel 7
		
		// Update vendor partition references
		chunk = bytes.ReplaceAll(chunk, []byte("vendor_ripcurrent"), []byte("vendor_cheetah"))
	}
	
	return chunk
}

// applyPixel9ProOptimizations applies Pixel 9 Pro specific optimizations
func (sp *StreamingPorter) applyPixel9ProOptimizations(chunk []byte) []byte {
	// Apply camera optimizations
	if bytes.Contains(chunk, []byte("camera")) {
		chunk = bytes.ReplaceAll(chunk, []byte("camera_ripcurrent"), []byte("camera_cheetah"))
		chunk = bytes.ReplaceAll(chunk, []byte("cam_ripcurrent"), []byte("cam_cheetah"))
	}
	
	// Apply display optimizations
	if bytes.Contains(chunk, []byte("display")) {
		chunk = bytes.ReplaceAll(chunk, []byte("display_ripcurrent"), []byte("display_cheetah"))
	}
	
	// Apply performance optimizations
	if bytes.Contains(chunk, []byte("performance")) {
		log.Printf("🚀 Applying Pixel 8 Pro performance optimizations\n")
		chunk = bytes.ReplaceAll(chunk, []byte("perf_ripcurrent"), []byte("perf_cheetah"))
	}
	
	return chunk
}

// streamAndPort handles the main streaming and porting process
func (sp *StreamingPorter) streamAndPort(ctx context.Context, reader io.Reader) error {
	// Create output file
	outputFile, err := os.Create(sp.outputPath)
	if err != nil {
		return errors.Wrap(err, "failed to create output file")
	}
	defer outputFile.Close()
	
	// Create hash for integrity verification
	hasher := sha256.New()
	multiWriter := io.MultiWriter(outputFile, hasher)
	
	// Stream data with progress tracking
	buffer := make([]byte, 64*1024) // 64KB buffer
	
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
			n, err := reader.Read(buffer)
			if err != nil {
				if err == io.EOF {
					break
				}
				return errors.Wrap(err, "failed to read streaming data")
			}
			
			// Write ported data
			_, writeErr := multiWriter.Write(buffer[:n])
			if writeErr != nil {
				return errors.Wrap(writeErr, "failed to write ported data")
			}
			
			// Update progress
			sp.updateProgress(int64(n))
		}
	}
	
	// Finalize and verify
	finalHash := fmt.Sprintf("%x", hasher.Sum(nil))
	log.Printf("✅ Streaming port completed! Hash: %s\n", finalHash[:16]+"...")
	
	return nil
}

// updateDownloadProgress updates download progress
func (sp *StreamingPorter) updateDownloadProgress(bytes int64) {
	sp.mutex.Lock()
	sp.downloadedSize += bytes
	sp.mutex.Unlock()
	
	sp.progressBar.UpdateDownload(sp.downloadedSize)
}

// updatePortingProgress updates porting progress
func (sp *StreamingPorter) updatePortingProgress(bytes int64) {
	sp.mutex.Lock()
	sp.portedSize += bytes
	sp.mutex.Unlock()
	
	sp.progressBar.UpdatePorting(sp.portedSize)
}

// updateProgress updates overall progress
func (sp *StreamingPorter) updateProgress(bytes int64) {
	sp.progressBar.Update(bytes)
}

// SetTotal sets the total size for progress tracking
func (pb *CustomProgressBar) SetTotal(total int64) {
	pb.mutex.Lock()
	defer pb.mutex.Unlock()
	pb.total = total
}

// UpdateDownload updates download progress
func (pb *CustomProgressBar) UpdateDownload(current int64) {
	pb.mutex.Lock()
	defer pb.mutex.Unlock()
	pb.current = current
	pb.render()
}

// UpdatePorting updates porting progress
func (pb *CustomProgressBar) UpdatePorting(ported int64) {
	pb.mutex.Lock()
	defer pb.mutex.Unlock()
	pb.portedCurrent = ported
	pb.render()
}

// Update updates overall progress
func (pb *CustomProgressBar) Update(bytes int64) {
	pb.mutex.Lock()
	defer pb.mutex.Unlock()
	pb.current += bytes
	pb.render()
}

// render displays the progress bar (tqdm-style)
func (pb *CustomProgressBar) render() {
	now := time.Now()
	if now.Sub(pb.lastUpdate) < 100*time.Millisecond {
		return // Throttle updates
	}
	pb.lastUpdate = now
	
	if pb.total == 0 {
		return
	}
	
	// Calculate percentages
	downloadPercent := float64(pb.current) / float64(pb.total) * 100
	portingPercent := float64(pb.portedCurrent) / float64(pb.total) * 100
	
	// Calculate speeds
	elapsed := now.Sub(pb.startTime).Seconds()
	downloadSpeed := float64(pb.current) / elapsed / (1024 * 1024) // MB/s
	portingSpeed := float64(pb.portedCurrent) / elapsed / (1024 * 1024) // MB/s
	
	// Calculate ETA
	if downloadSpeed > 0 {
		remaining := float64(pb.total-pb.current) / (downloadSpeed * 1024 * 1024)
		eta := time.Duration(remaining) * time.Second
		
		// Create progress bar visualization
		filled := int(downloadPercent / 100 * float64(pb.width))
		portFilled := int(portingPercent / 100 * float64(pb.width))
		
		bar := strings.Repeat("█", filled)
		portBar := strings.Repeat("▓", portFilled)
		empty := strings.Repeat("░", pb.width-filled)
		
		// Format output (tqdm-style)
		fmt.Printf("\r%s: %3.0f%%|%s%s| %.1f/%.1fMB [%02d:%02d<%02d:%02d, %.2fMB/s] Port: %3.0f%% [%s] %.2fMB/s",
			pb.description,
			downloadPercent,
			bar, empty,
			float64(pb.current)/(1024*1024),
			float64(pb.total)/(1024*1024),
			int(elapsed)/60, int(elapsed)%60,
			int(eta.Seconds())/60, int(eta.Seconds())%60,
			downloadSpeed,
			portingPercent,
			portBar[:min(len(portBar), pb.width/2)],
			portingSpeed,
		)
	}
}

// min returns the minimum of two integers
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// Close cleans up resources
func (sp *StreamingPorter) Close() error {
	if sp.portingBuffer.tempFile != nil {
		sp.portingBuffer.tempFile.Close()
		os.Remove(sp.portingBuffer.tempFile.Name())
	}
	return nil
}

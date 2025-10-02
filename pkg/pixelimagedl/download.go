package pixelimagedl

import (
	"context"
	"crypto/sha256"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/pkg/errors"

	"github.com/jalavosus/pixelimagedl/internal"
	"github.com/jalavosus/pixelimagedl/pkg/download"
)

func DownloadLatest(ctx context.Context, device Pixel, downloadType DownloadType, outDir string) error {
	listCtx, listCancel := context.WithTimeout(ctx, 30*time.Second)
	defer listCancel()

	images, err := ListDeviceImages(listCtx, device, downloadType)
	if err != nil {
		err = errors.WithMessagef(err, "error scraping available %[1]s images for device %[2]s", downloadType.String(), device.String())
		return err
	}

	latest := internal.SliceLast(images)

	log.Printf("latest stable %[1]s image for %[2]s is %[3]s (%[4]s)\n", downloadType.String(), device.String(), latest.Version, latest.BuildNumber)

	var filename string

	downloadUri := latest.DownloadURI
	split := strings.Split(downloadUri, "/")

	// Use custom filename format for Android 17 DP: result_firmware_devicename.zip
	if downloadType == Android17DP {
		deviceName := strings.ToLower(device.String())
		deviceName = strings.ReplaceAll(deviceName, " ", "")
		filename = fmt.Sprintf("result_firmware_%s.zip", deviceName)
	} else {
		filename = internal.SliceLast(split)
	}
	
	if !filepath.IsAbs(outDir) {
		outDir, err = filepath.Abs(outDir)
		if err != nil {
			return err
		}
	}
	filename = filepath.Join(outDir, filename)

	log.Printf("downloading %[1]s image from %[2]s\n", downloadType.String(), downloadUri)

	var resp *http.Response
	var numBytes int64
	
	// For Android 17 DP, implement bypass logic to try multiple servers
	if downloadType == Android17DP {
		// Extract device codename from the URL or use device parameter
		deviceCodename := strings.ToLower(device.String())
		if strings.Contains(deviceCodename, " ") {
			// Convert device name to codename (e.g., "Pixel 7 Pro" -> "cheetah")
			deviceCodename = deviceToCodename(deviceCodename)
		}
		
		bypassURLs := getAndroid17BypassURLs(deviceCodename, latest.BuildNumber)
		log.Printf("attempting bypass download with %d URL combinations\n", len(bypassURLs))
		
		var lastErr error
		for i, url := range bypassURLs {
			log.Printf("trying URL %d/%d: %s\n", i+1, len(bypassURLs), url)
			
			req, _ := http.NewRequestWithContext(ctx, http.MethodGet, url, http.NoBody)
			resp, err = http.DefaultClient.Do(req)
			
			if err != nil {
				lastErr = err
				log.Printf("URL %d failed: %v\n", i+1, err)
				continue
			}
			
			if resp.StatusCode == http.StatusOK {
				log.Printf("SUCCESS: Found working URL %d: %s\n", i+1, url)
				break
			} else {
				resp.Body.Close()
				log.Printf("URL %d returned status %d\n", i+1, resp.StatusCode)
				lastErr = errors.Errorf("HTTP %d", resp.StatusCode)
				continue
			}
		}
		
		if resp == nil || resp.StatusCode != http.StatusOK {
			return errors.WithMessagef(lastErr, "all bypass URLs failed for Android 17 DP")
		}
	} else {
		// Standard download for other types
		req, _ := http.NewRequestWithContext(ctx, http.MethodGet, downloadUri, http.NoBody)
		resp, err = http.DefaultClient.Do(req)
		if err != nil {
			err = errors.WithMessagef(err, "error downloading file at url %[1]s", downloadUri)
			return err
		}
	}

	defer func() {
		if closeErr := resp.Body.Close(); closeErr != nil {
			log.Printf("error closing file download body: %v\n", closeErr)
		}
	}()

	log.Printf("saving %[1]s image to %[2]s\n", downloadType.String(), filename)
	numBytes, err = download.ReadData(resp, filename, dlBufSize())
	if err != nil {
		return err
	}

	log.Printf("saved %-.1[1]fGb to %[2]s", download.GbFromBytes(numBytes), filename)

	// Perform SHA256 verification if hash is available
	if latest.SHA256Sum != "" {
		gotSha, shaMatch := checkSha(filename, latest.SHA256Sum)
		if !shaMatch {
			return errors.Errorf("SHA256 mismatch; expected %[1]s, sum of downloaded file is %[2]s", latest.SHA256Sum, gotSha)
		} else {
			log.Printf("SHA256 sum %[1]s of downloaded file matches expected\n", gotSha)
		}
	} else {
		log.Printf("SHA256 verification skipped (no hash available from source)\n")
	}

	return nil
}

func checkSha(filename, wantSha string) (string, bool) {
	f, err := os.Open(filename)
	if err != nil {
		panic(err)
	}

	defer func() {
		_ = f.Close()
	}()

	h := sha256.New()
	if _, err = io.Copy(h, f); err != nil {
		panic(err)
	}

	check := fmt.Sprintf("%x", h.Sum(nil))

	return check, check == wantSha
}

// deviceToCodename converts device display names to codenames for URL generation
func deviceToCodename(deviceName string) string {
	deviceName = strings.ToLower(deviceName)
	deviceName = strings.ReplaceAll(deviceName, " ", "")
	
	// Map device names to codenames
	deviceMap := map[string]string{
		"pixel7pro":     "cheetah",
		"pixel7":        "panther", 
		"pixel7a":       "lynx",
		"pixel8pro":     "husky",
		"pixel8":        "shiba",
		"pixel8a":       "akita",
		"pixel9pro":     "caiman",
		"pixel9":        "tokay",
		"pixel9proxl":   "komodo",
		"pixel9profold": "comet",
	}
	
	if codename, exists := deviceMap[deviceName]; exists {
		return codename
	}
	
	// If no mapping found, return the cleaned device name
	return deviceName
}

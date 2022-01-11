package redditdl

import (
	"errors"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/gofrs/uuid"
	"github.com/tidwall/gjson"
)

func Download(postURL string, writer io.Writer) error {
	if !strings.HasSuffix(postURL, ".json") {
		postURL = strings.TrimSuffix(postURL, "/") + ".json"
	}

	postResponse, postRequestError := get(postURL)
	if postRequestError != nil {
		return postRequestError
	}

	if postResponse.StatusCode != http.StatusOK {
		return errors.New("error retrieving post, it might've been deleted")
	}

	bytes, readError := ioutil.ReadAll(postResponse.Body)
	if readError != nil {
		return readError
	}

	videoURL := gjson.GetBytes(bytes, "0.data.children.0.data.secure_media.reddit_video.fallback_url").Value().(string)
	videoResponse, videoResponseError := get(videoURL)
	if videoResponseError != nil {
		return videoResponseError
	}

	audioURL := regexp.MustCompile(`_.*\.`).ReplaceAllString(videoURL, "_audio.")
	audioResponse, audioResponseError := get(audioURL)
	if audioResponseError != nil {
		return audioResponseError
	}

	//Audio available
	if audioResponse.StatusCode != http.StatusForbidden {
		tempVideo, tempError := os.CreateTemp(os.TempDir(), "*.video.mp4")
		if tempError != nil {
			return tempError
		}
		defer tempVideo.Close()

		tempAudio, tempError := os.CreateTemp(os.TempDir(), "*.audio.mp4")
		if tempError != nil {
			return tempError
		}
		defer tempAudio.Close()

		var copyError error
		_, copyError = io.Copy(tempVideo, videoResponse.Body)
		if copyError != nil {
			return copyError
		}

		_, copyError = io.Copy(tempAudio, audioResponse.Body)
		if copyError != nil {
			return copyError
		}

		tempCombinedName := filepath.Join(os.TempDir(), uuid.Must(uuid.NewV4()).String()+".combined.mp4")
		ffmpegError := exec.
			Command("ffmpeg", "-i", tempVideo.Name(), "-i", tempAudio.Name(), "-c:v", "copy", "-c:a", "aac", tempCombinedName).
			Run()
		if ffmpegError != nil {
			return ffmpegError
		}

		tempCombinedFile, openError := os.Open(tempCombinedName)
		if openError != nil {
			return openError
		}

		_, copyError = io.Copy(writer, tempCombinedFile)
		return copyError
	}

	//Only video available
	_, copyError := io.Copy(writer, videoResponse.Body)
	return copyError
}

func get(url string) (*http.Response, error) {
	request, requestError := http.NewRequest(http.MethodGet, url, nil)
	if requestError != nil {
		return nil, requestError
	}

	// Believable Desktop User-Agent, so reddit doesn't discriminate against our robot.
	request.Header.Add("User-Agent", "Mozilla/5.0 (X11; Ubuntu; Linux x86_64; rv:95.0) Gecko/20100101 Firefox/95.0")

	return http.DefaultClient.Do(request)
}

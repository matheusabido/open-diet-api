package whisper

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

const (
	timeout  = 60 * time.Second
	language = "pt"

	// whisper.cpp only reads 16 kHz mono 16-bit PCM. The filters drop the
	// frequencies outside the voice range.
	filters = "highpass=f=80,lowpass=f=7600"

	// The encoder always processes a full 30 s window. The audio context cuts
	// that window to what the sentence needs: 768 frames are 15.36 s of audio,
	// so the audio is limited to 15 s.
	maxDuration  = "15"
	audioContext = "768"

	// Biases the model towards digits and the "amount, unit, food" format.
	prompt = "100 gramas de frango. 200 ml de leite. 2 onças de tomate. 1 unidade de ovo."
)

type Transcriber struct {
	cliPath   string
	modelPath string
}

func New(cliPath, modelPath string) *Transcriber {
	return &Transcriber{cliPath: cliPath, modelPath: modelPath}
}

func (t *Transcriber) Configured() bool {
	return t.cliPath != "" && t.modelPath != ""
}

func (t *Transcriber) Transcribe(ctx context.Context, audio io.Reader) (string, error) {
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	dir, err := os.MkdirTemp("", "open-diet-voice-")
	if err != nil {
		return "", err
	}
	defer os.RemoveAll(dir)

	input := filepath.Join(dir, "input")
	if err := save(input, audio); err != nil {
		return "", err
	}

	wav := filepath.Join(dir, "audio.wav")
	if _, err := run(ctx, "ffmpeg", "-nostdin", "-loglevel", "error", "-y",
		"-i", input,
		"-t", maxDuration,
		"-vn", "-ac", "1", "-ar", "16000", "-af", filters, "-c:a", "pcm_s16le",
		wav,
	); err != nil {
		return "", fmt.Errorf("optimize audio: %w", err)
	}

	text, err := run(ctx, t.cliPath,
		"-m", t.modelPath,
		"-f", wav,
		"-l", language,
		"--prompt", prompt,
		"-ac", audioContext,
		"-nt", "-np",
	)
	if err != nil {
		return "", fmt.Errorf("whisper: %w", err)
	}

	return strings.Join(strings.Fields(text), " "), nil
}

func save(path string, audio io.Reader) error {
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = io.Copy(file, audio)
	return err
}

func run(ctx context.Context, name string, args ...string) (string, error) {
	var stdout, stderr bytes.Buffer

	cmd := exec.CommandContext(ctx, name, args...)
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("%w: %s", err, strings.TrimSpace(stderr.String()))
	}

	return stdout.String(), nil
}

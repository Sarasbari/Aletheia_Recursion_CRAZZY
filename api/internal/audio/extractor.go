package audio

import (
	"context"
	"fmt"
	"os/exec"
)

func ExtractMP3(ctx context.Context, inputVideoPath, outputAudioPath string) error {
	cmd := exec.CommandContext(
		ctx,
		"ffmpeg",
		"-y",
		"-i", inputVideoPath,
		"-q:a", "0",
		"-map", "a",
		outputAudioPath,
	)

	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("ffmpeg extract failed: %w: %s", err, string(output))
	}
	return nil
}

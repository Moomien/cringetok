package controllers

import (
	"fmt"
	"os"
	"os/exec"
)

type MoodAssets struct {
	Music      string
	Emoji      string
	Gif        string
	Atmosphere string
}

var moods = map[string]MoodAssets{
	"sad": {
		Music:      "sad_music.mp3",
		Emoji:      "cry_emoji.png",
		Gif:        "rain.mp4",
		Atmosphere: "format=gray",
	},
	"happy": {
		Music:      "happy.mp3",
		Emoji:      "happy.png",
		Gif:        "happy.mp4",
		Atmosphere: "eq=brightness=0.06:contrast=1.15:saturation=1.35,colorbalance=rs=0.05:gs=0.02:bs=-0.03",
	},
	"playful": {
		Music:      "playful.mp3",
		Emoji:      "playful.png",
		Gif:        "playful.mp4",
		Atmosphere: "colorchannelmixer=1:0:0:0:0:0:0:0:0:0:1:0",
	},
}

// берем скрин по destpath из uploadhandler
// берем музыку по assets/music и делаем 10 сек видео
// потом берем это видео и добавляем на него эмодзи и гифку в зависимости от настроения
func GenerateVideo(destPath string, mood string, uuid string) string {
	asset, ok := moods[mood]
	if !ok {
		fmt.Println("Error undefined mood: ", mood)
		return ""
	}

	outputVideo := "./assets/output/" + uuid + "output.mp4"
	emojiPath := "./assets/emojis/" + asset.Emoji
	gifPath := "./assets/gifs/" + asset.Gif
	musicPath := "./assets/music/" + asset.Music
	finalPath := "./assets/output/" + uuid + "final.mp4"
	finishPath := "./assets/output/" + uuid + "finish.mp4"

	workDir, _ := os.Getwd()

	//создание видео с музыкой
	cmd1 := exec.Command("ffmpeg",
		"-loop", "1",
		"-i", destPath, //screenshot
		"-i", musicPath,
		"-vf", "scale=trunc(iw/2)*2:trunc(ih/2)*2",
		"-c:v", "libx264",
		"-t", "10",
		"-pix_fmt", "yuv420p",
		"-c:a", "aac",
		"-b:a", "192k",
		outputVideo,
	)
	cmd1.Dir = workDir
	if err := cmd1.Run(); err != nil {
		fmt.Println("err1:", err)
		return ""
	}

	filter := "[1:v] format=rgba, colorchannelmixer=aa=0.5 [emoji_alpha]; " +
		"[2:v] scale='min(iw, 304/4)':-1 [small_scaled]; " +
		"[0:v][emoji_alpha] overlay=10:10 [tmp]; " +
		"[tmp][small_scaled] overlay=W-w-10:H-h-10"

	//гифка + картинка
	cmd2 := exec.Command("ffmpeg",
		"-i", outputVideo,
		"-i", emojiPath,
		"-i", gifPath,
		"-filter_complex", filter,
		"-c:a", "copy",
		"-y",
		finalPath,
	)
	cmd2.Dir = workDir
	if err := cmd2.Run(); err != nil {
		fmt.Println("err2:", err)
		return ""
	}

	//спецэффекты
	cmd3 := exec.Command("ffmpeg",
		"-i", finalPath,
		"-vf", asset.Atmosphere,
		"-c:v", "libx264", "-pix_fmt", "yuv420p", "-profile:v", "high", "-level", "4.2", "-c:a",
		"aac", "-b:a", "192k", "-movflags", "+faststart",
		"-y",
		finishPath,
	)
	cmd3.Dir = workDir
	if err := cmd3.Run(); err != nil {
		fmt.Println("err3:", err)
		return ""
	}

	cleaning(finalPath, outputVideo, destPath)

	return finishPath
}

func cleaning(finalPath string, outputVideo string, destPath string) {
	err := os.Remove(finalPath)
	if err != nil {
		fmt.Println("Error deleting:", err)
		return
	}

	err = os.Remove(outputVideo)
	if err != nil {
		fmt.Println("Error deleting:", err)
		return
	}

	err = os.Remove(destPath)
	if err != nil {
		fmt.Println("Error deleting:", err)
		return
	}
}

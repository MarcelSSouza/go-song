package main
import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"time"
	"strings"

	"github.com/faiface/beep"
	"github.com/faiface/beep/mp3"
	"github.com/faiface/beep/speaker"
	"github.com/kkdai/youtube/v2"
)

func directory(dirname string) {
	files, err := ioutil.ReadDir(dirname)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Music files found in the directory:")
	for _, file := range files {
		fmt.Println(file.Name())
	}

	for _, file := range files {
		playMusic(dirname + "/" + file.Name())
	}

	if len(files) == 0 {
		fmt.Println("No music files found in the directory")
	}
}


func playMusic(filename string) {
	f, err := os.Open(filename)
	if err != nil {
		fmt.Println("Error opening file:", err)
		return
	}

	streamer, format, err := mp3.Decode(f)
	if err != nil {
		fmt.Println("Error decoding MP3:", err)
		return
	}
	defer streamer.Close()

	err = speaker.Init(format.SampleRate, format.SampleRate.N(time.Second/10))
	if err != nil {
		fmt.Println("Error initializing speaker:", err)
		return
	}

	done := make(chan bool)
	fmt.Println("Playing", filename + "...")
	speaker.Play(beep.Seq(streamer, beep.Callback(func() {
		fmt.Println("Music Finished. Playing next music in 5 seconds...")
		time.Sleep(5 * time.Second)
		cmd := exec.Command("clear")
		cmd.Run()
		done <- true //  song has finished
	})))

	<-done
}

func main() {
	if len(os.Args) < 2 {
		cmd := exec.Command("clear")
		cmd.Stdout = os.Stdout
		cmd.Run()
		directory("./musics")
	} else if len(os.Args) == 2 {
		fmt.Println("Downloading music from youtube...")
		youtubeURL := os.Args[1]
		println(youtubeURL)

		videoID := strings.Split(youtubeURL,"=")[1]
		println(videoID)
		client := youtube.Client{}

		video, err := client.GetVideo(videoID)
		if err != nil {
			panic(err)
		}
	
		formats := video.Formats.WithAudioChannels() // only get videos with audio
		for _, format := range formats {
			fmt.Println(format)
		}
		stream, _, err := client.GetStream(video, &formats[3])
		if err != nil {
			panic(err)
		}
		defer stream.Close()
	
		outputFileName := "./musics/" + videoID + ".mp3"
		cmd := exec.Command("ffmpeg", "-i", "pipe:0", "-vn", "-acodec", "libmp3lame", "-q:a", "0", outputFileName)
		cmd.Stdin = stream
		err = cmd.Run()
		if err != nil {
			panic(err)
		}
	
		fmt.Println("Conversion completed. Saved as", outputFileName)
	}
}


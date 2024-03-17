package main
import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"io"
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
		//get only ID form https://www.youtube.com/watch?v=69RdQFDuYPI
		// 69RdQFDuYPI is the ID
		videoID := strings.Split(youtubeURL,"=")[1]

		client := youtube.Client{}

		video, err := client.GetVideo(videoID)
		if err != nil {
			panic(err)
		}
	
		formats := video.Formats.WithAudioChannels() // only get videos with audio
		stream, _, err := client.GetStream(video, &formats[0])
		if err != nil {
			panic(err)
		}
		defer stream.Close()
	
		file, err := os.Create("video.mp4")
		if err != nil {
			panic(err)
		}
		defer file.Close()
	
		_, err = io.Copy(file, stream)
		if err != nil {
			panic(err)
		}
	}
}


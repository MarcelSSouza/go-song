package main
import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"time"

	"github.com/faiface/beep"
	"github.com/faiface/beep/mp3"
	"github.com/faiface/beep/speaker"
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

func downloadAudio(url string) error {
    cmd := exec.Command("youtube-dl", "-x", "--audio-format", "mp3", "-o", " ./musics/"+url+".mp3" , url)
    err := cmd.Run()
    if err != nil {
        return err
    }
    return nil
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

		if youtubeURL != "" {
			err := downloadAudio(youtubeURL)
			if err != nil {
				panic(err)
			}
			println("Audio downloaded successfully!")
		}
	}
}


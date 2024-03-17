import os
from pytube import YouTube
import argparse
from pydub import AudioSegment

def download_mp3(url, save_path):
    yt = YouTube(url)
    title = yt.title
    audio_stream = yt.streams.filter(only_audio=True, file_extension='mp4').first()
    audio_file = audio_stream.download(output_path=save_path)
    
    # Convert the downloaded MP4 audio to MP3
    audio_file_mp4 = os.path.join(save_path, audio_stream.default_filename)
    audio_file_mp3 = os.path.join(save_path, f"{title}.mp3")
    
    audio = AudioSegment.from_file(audio_file_mp4)
    audio.export(audio_file_mp3, format="mp3")
    
    # Clean up the temporary MP4 file
    os.remove(audio_file_mp4)


def main():
    parser = argparse.ArgumentParser(description="Convert YouTube video to MP3")
    parser.add_argument("url", help="YouTube video URL")
    parser.add_argument("save_path", help="Path to save the MP3 file")
    args = parser.parse_args()

    if not os.path.exists(args.save_path):
        os.makedirs(args.save_path)

    download_mp3(args.url, args.save_path)

if __name__ == "__main__":
    main()

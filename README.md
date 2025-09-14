# Media Orientation (mediaorient)

<p align="center">
<img src="docs/images/icon.avif" width="240" alt="mediaorient"/>
<br/>
<strong>mediaorient</strong> is a CLI tool and Go library to calculate/fix the orientation of images & videos.
</p>

## 🖼️ Usage

You can use **mediaorient** in two ways: as a command-line interface (CLI) tool or a Go library.

The CLI tool is a standalone application that can be used to calculate and fix the orientation of media files, while the library can be integrated into your own Go projects.

The binaries are available for Windows, macOS, and Linux. Download the [latest release](https://github.com/vegidio/mediaorient/releases) that matches your computer architecture and operating system.

### CLI

<p align="center">
<img src="docs/images/screenshot.avif" width="80%" alt="mediaorient"/>
</p>

<details>
<summary>Calculating the similarity score of two files</summary>

#### Run the command below in the terminal:

```bash
$ mediaorient score <media1> <media2>
```
</details>

<details>
<summary>Comparing two or more files</summary>

#### Run the command below in the terminal:

```bash
$ mediaorient files <media1> <media2> [<media3> ...]
```

Where:

- `files` (mandatory): the path to the media files you want to compare. You must pass at least two files, separated by space.
</details>

<details>
<summary>Comparing multiple files in a directory</summary>

#### Run the command below in the terminal:

```bash
$ mediaorient dir <directory> [-r] [--mt <media-type>]
```

Where:

- `directory` (mandatory): the path to the directory where the media files are located.
- `-r` (optional): recursively search for files in subdirectories to include in the comparison.
- `--mt` (optional): the file types to be included in the comparison. You can choose between `image`, `video`, or `all` (default).
</details>

<details>
<summary>Renaming files based on similarity</summary>

#### Run the command below in the terminal:

```bash
$ mediaorient rename <directory> [-r] [--mt <media-type>]
```

Where:

- `directory` (mandatory): the path to the directory where the media files are located.
- `-r` (optional): recursively search for files in subdirectories to include in the comparison.
- `--mt` (optional): the file types to be included in the comparison. You can choose between `image`, `video`, or `all` (default).
</details>

---

Other parameters you can use:

- `-t` (optional): the threshold for the similarity score; a value between 0–1, where 0 is completely different and 1 is identical. The default value is `0.8`, which means only similarities of 80% or higher will be reported.
- `-o` (optional): the output format; you can choose `report` (default) or, if you prefer a raw output, `json` or `csv`.
- `--ie` (optional): ignores errors and continues the comparison even if some files are not valid.
- `--ff` (optional; images only): flips the frames vertically and horizontally during the comparison.
- `--fr` (optional; images only): rotates the frames in multiple angles during the comparison.

For the full list of parameters, type `mediaorient --help` in the terminal.

## 🎞️ Supported media types

In its default configuration, **mediaorient** supports media files with the following extensions:

- Images: `.bmp`, `.gif`, `.jpg` (`.jpeg`), `.png`, `.tiff`, `.webp`
- Videos: `.avi`, `.mp4` (`.m4v`), `.mkv`, `.mov`, `.webm`

If you want to work with additional file extensions, you can use the functions `AddImageType` or `AddVideoType` before performing any similarity comparisons. This allows **mediaorient** to include these file types during calculations.

When adding support for new media formats, it's essential to load a 3rd party library capable of decoding them. For example, to enable AVIF image comparison in **mediaorient**, you could use a library like [avif-go](https://github.com/vegidio/avif-go) to do this:

```go
import _ "github.com/vegidio/avif-go"
mediaorient.AddImageType(".avif")
```

## 💣 Troubleshooting

### Miscalculation of Media Orientation

This app uses a neural network to determine the orientation of images and videos. While this method is accurate about 99% of the time, it isn’t perfect. It works very well for everyday media—such as photos of people, animals, documents and common scenes, but may perform less reliably with certain special cases.

For example, if you have a media of a medical procedure, the app may not detect the correct orientation because the model was not trained (or was only minimally trained), on this type of media. You can help improve the app by contributing to the training dataset. Simply upload your own media [here](https://mega.nz/filerequest/eGNGUqolkGI), so that in future versions of the app this kind of media will be properly detected.

Any media you upload will not be shared with anyone else and will be used solely to train the neural network. I’m not here to judge anyone, so feel free to upload any *legal* content.

However, I cannot emphasize this enough: **do NOT upload any illegal content**. If you upload illegal content, I will not only block your access to this project, but I will also report you to the authorities using your IP address and any other information available.

### The App Binary Is Too Large

There are two main reasons for this:

1. This app depends on [ONNX Runtime](https://onnxruntime.ai). I could ask users to install this dependency separately, but not everyone is tech-savvy, and that would create unnecessary difficulties. To ensure a smoother experience, I chose to embed ONNX Runtime directly in the app.
2. The pre-trained neural network is quite large. (Generally, the larger the model, the more accurate it is at determining media orientation.) This contributes significantly to the app’s binary size.

Unfortunately, there isn’t much I can do about this. If you are using this project as a Go library, expect your final binary to increase by at least 100 MB.

### Video Orientation Doesn't Work

If the orientation calculation/fix of videos is not working, it may be because you don't have [FFmpeg](https://www.ffmpeg.org/download.html) working in your computer, which is required to extract frames from the video files.

When FFmpeg is not found, **mediaorient** will try to automatically download and install it for you. Even though this will work in most cases, it may fail for unpredictable reasons.

The best option to have the video comparison working is to install FFmpeg yourself in your computer and make sure it is available in your `PATH`.

### "App Is Damaged/Blocked..." (Windows & macOS only)

For a couple of years now, Microsoft and Apple have required developers to join their "Developer Program" to gain the pretentious status of an _identified developer_ 😛.

Translating to non-BS language, this means that if you’re not registered with them (i.e., paying the fee), you can’t freely distribute Windows or macOS software. Apps from unidentified developers will display a message saying the app is damaged or blocked and can’t be opened.

To bypass this, open the Terminal and run one of the commands below (depending on your operating system), replacing `<path-to-app>` with the correct path to where you’ve installed the app:

- Windows: `Unblock-File -Path <path-to-app>`
- macOS: `xattr -d com.apple.quarantine <path-to-app>`

## 🛠️ Build

### Dependencies

To build this project, you will need the following dependencies installed in your computer:

- [Golang](https://go.dev/doc/install)
- [Task](https://taskfile.dev/installation/)

### Compiling

With all the dependencies installed, in the project's root folder run the command:

```bash
$ task cli os=<operating-system> arch=<architecture>
```

Where:

- `<operating-system>`: can be `windows`, `darwin` (macOS), or `linux`.
- `<architecture>`: can be `amd64` or `arm64`.

For example, if I wanted to build the CLI for Windows, on architecture AMD64, I would run the command:

```bash
$ task cli os=windows arch=amd64
```

## 📝 License

**mediaorient** is released under the MIT License. See [LICENSE](LICENSE) for details.

## 👨🏾‍💻 Author

Vinicius Egidio ([vinicius.io](http://vinicius.io))

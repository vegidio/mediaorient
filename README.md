# Media Orientation (mediaorient)

<p align="center">
<img src="docs/images/icon.avif" width="240" alt="mediaorient"/>
<br/>
<strong>mediaorient</strong> is a CLI tool and Go library to calculate the orientation of images & videos.
</p>

## 🖼️ Usage

You can use **mediaorient** in two ways: as a command-line interface (CLI) tool or a Go library.

The CLI tool is a standalone application that can be used to calculate and fix the orientation of media files, while the library can be integrated into your own Go projects.

The binaries are available for Windows, macOS, and Linux. Download the [latest release](https://github.com/vegidio/mediaorient/releases) that matches your computer architecture and operating system.

### CLI

<p align="center">
<img src="docs/images/screenshot.avif" width="80%" alt="mediaorient"/>
</p>

#### Calculate the orientation of one or more files

Run the command below in the terminal:

```bash
$ mediaorient files <media1> [<media2> ...]
```

#### Calculate the orientation of a directory

Run the command below in the terminal:

```bash
$ mediaorient dir <directory> [-r] [--mt <media-type>]
```

Where:

- `dir` (mandatory): the path to the directory where the media files are located.
- `-r` (optional): recursively search for files in subdirectories to include in the calculation.
- `--mt` (optional): the file types to be included in the calculation. You can choose between `image`, `video`, or `all` (default).

---

Other parameters you can use:

- `-o` (optional): the output format; you can choose `report` (default) or, if you prefer a raw output, `json` or `csv`.
- `--ie` (optional): ignores errors and continues the calculation even if some files are not valid.

For the full list of parameters, type `mediaorient --help` in the terminal.

## 🎞️ Supported media types

In its default configuration, the **mediaorient** library supports media files with the following extensions:

- Images: `.bmp`, `.gif`, `.jpg` (`.jpeg`), `.png`, `.tiff`, `.webp`
- Videos: `.avi`, `.mp4` (`.m4v`), `.mkv`, `.mov`, `.webm`

The CLI supports two additional image formats: `.avif` and `.heic`.

If you want to work with additional file extensions in the library, like those two above, you can use the functions `AddImageType` or `AddVideoType` before performing any orientation calculation. This allows **mediaorient** to include these file types during calculations.

When adding support for new media formats, it's essential to load a 3rd party library capable of decoding them. For example, to enable AVIF image calculation in **mediaorient**, you could use a library like [avif-go](https://github.com/vegidio/avif-go) to do this:

```go
import _ "github.com/vegidio/avif-go"
mediasim.AddImageType(".avif")
```

## 💣 Troubleshooting

### Internet Connection Required On The First Run

This app uses a Large Language Model (LLM) that was trained with almost 1 million images, creating a neural network of nearly 100 MB. Unfortunately, it becomes impractical to embed this data in the executable, otherwise the file would become to big.

For this reason, the first time you run the app, it will download the model from the internet; this may take a few seconds, depending on your internet connection.

Besides preventing the executable from being too big, this also allows app to use the most recent version of the model, which may contain improvements and bug fixes.

### Miscalculation of Media Orientation

While the LLM has accuracy of about 99%, it isn’t perfect. It works very well for everyday media—such as photos of people, animals, documents and common scenes, but may perform less reliably with certain special cases.

For example, if you have a media of a medical procedure, the app may not detect the correct orientation because the model was not trained (or was only minimally trained), on this type of media. You can help improve the app by contributing to the training dataset. Simply upload your own media [here](https://mega.nz/filerequest/eGNGUqolkGI), so that in future versions of the app this kind of media will be properly detected.

Any media you upload will not be shared with anyone else and will be used solely to train the neural network. I’m not here to judge anyone, so feel free to upload any *legal* content.

However, I cannot emphasize this enough: **do NOT upload any illegal content**. If you upload illegal content, I will not only block your access to this project, but I will also report you to the authorities using your IP address and any other information available.

### Video Orientation Doesn't Work

If the orientation calculation/fix of videos is not working, it may be because you don't have [FFmpeg](https://www.ffmpeg.org/download.html) working in your computer, which is required to extract frames from the video files.

When FFmpeg is not found, **mediaorient** will try to automatically download and install it for you. Even though this will work in most cases, it may fail for unpredictable reasons.

The best option to have the video calculation working is to install FFmpeg yourself in your computer and make sure it is available in your `PATH`.

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

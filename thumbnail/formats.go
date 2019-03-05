package thumbnail

type ImageFormat uint8

const (
    Unknown ImageFormat = iota
    JPEG
    PNG
)

func Format(ext string) ImageFormat {
    switch ext {
    case ".jpg", ".jpeg":
        return JPEG
    case ".png":
        return PNG
    default:
        return Unknown
    }
}
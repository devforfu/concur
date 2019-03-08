package thumbnail

import (
    "fmt"
    "image"
    "image/jpeg"
    "image/png"
    "io"
    "os"
    "path/filepath"
    "strings"
)

// Maker converts an image or a set of images into thumbnails.
type Maker struct {
    // Format defines generated thumbnail format.
    Format ImageFormat
    // Width and Height define thumbnail maximal sizes.
    Width, Height int
}

func NewMaker(format ImageFormat, width, height int) *Maker {
    return &Maker{format, width, height}
}

// Thumbnail converts filename image into smaller image file appending .thumb
// suffix before image extension.
func (t *Maker) Thumbnail(filename string) (string, error) {
    infile, outfile := filename, Name(filename)
    err := t.Create(infile, outfile)
    if err != nil {
        return "", err
    }
    return outfile, nil
}

// Create takes infile image, converts it into thumbnail and saves into outfile.
func (t *Maker) Create(infile, outfile string) (err error) {
    in, err := os.Open(infile)
    if err != nil { return }
    defer in.Close()

    out, err := os.Create(outfile)
    if err != nil { return }

    if err := t.Convert(out, in); err != nil {
        _ = out.Close()
        return fmt.Errorf("scaling %s to %s: %s", infile, outfile, err)
    }

    return out.Close()
}

// Convert reads image from r and writes thumbnail into w.
func (t *Maker) Convert(w io.Writer, r io.Reader) (err error) {
    src, _, err := image.Decode(r)
    if err != nil { return }
    dst := CreateThumbnailImage(src, t.Width, t.Height)
    switch t.Format {
    case JPEG:
        return jpeg.Encode(w, dst, nil)
    case PNG:
        return png.Encode(w, dst)
    default:
        return fmt.Errorf("unknown image format")
    }
}

// ConvertFolder traverses the dirname folder and converts all images discovered
// images into thumbnails. The pattern string define which image extensions to
// look for and should be a pipe-separated string, like "jpeg|png".
func (t *Maker) ConvertFolder(dirname string, pattern string) ([]string, error) {
    thumbs := make([]string, 0)
    files, err := DiscoverImages(dirname, pattern)
    if err != nil { return nil, err }
    for _, file := range files {
        out := Name(file)
        err = t.Create(file, out)
        if err != nil {
            return nil, fmt.Errorf("%s: cannot create thumbnail: %s", err, file)
        }
        thumbs = append(thumbs, out)
    }
    return thumbs, nil
}

// CreateThumbnailImage rescales the src image into thumbnail with maximal
// dimensions defined by width and height.
func CreateThumbnailImage(src image.Image, width, height int) image.Image {
    xs := src.Bounds().Size().X
    ys := src.Bounds().Size().Y

    if aspect := float64(xs)/float64(ys); aspect < 1.0 {
        width = int(float64(width) * aspect)
    } else {
        height = int(float64(height) / aspect)
    }

    xScale := float64(xs)/float64(width)
    yScale := float64(ys)/float64(height)
    dst := image.NewRGBA(image.Rect(0, 0, width, height))

    for x := 0; x < width; x++ {
        for y := 0; y < height; y++ {
            xSrc := int(float64(x)*xScale)
            ySrc := int(float64(y)*yScale)
            dst.Set(x, y, src.At(xSrc, ySrc))
        }
    }

    return dst
}

func Name(file string) string {
    ext := filepath.Ext(file)
    thumb := strings.TrimSuffix(file, ext) + ".thumb" + ext
    return thumb
}

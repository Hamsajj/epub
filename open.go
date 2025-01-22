package epub

import (
	"archive/zip"
	"bytes"
	"fmt"
	"io"
)

// Open open a epub file
func Open(fn string) (*Book, func() error, error) {
	fd, err := zip.OpenReader(fn)
	if err != nil {
		return nil, nil, err
	}

	bk, err := loadBook(&fd.Reader)
	if err != nil {
		fd.Close()
		return nil, nil, err
	}

	return bk, fd.Close, nil
}

func OpenFromReader(reader io.Reader) (*Book, error) {
	buff := bytes.NewBuffer([]byte{})
	_, err := io.Copy(buff, reader)
	if err != nil {
		return nil, fmt.Errorf("failed to copy reader to buffer: %v", err)
	}
	return OpenFromBytes(buff.Bytes())
}

func OpenFromBytes(data []byte) (*Book, error) {
	fd, err := zip.NewReader(bytes.NewReader(data), int64(len(data)))
	if err != nil {
		return nil, fmt.Errorf("failed to open zip reader: %v", err)
	}
	book, err := loadBook(fd)
	if err != nil {
		return nil, err
	}
	return book, nil
}

func loadBook(fd *zip.Reader) (*Book, error) {
	bk := Book{fd: fd}
	mt, err := bk.readBytes("mimetype")
	if err == nil {
		bk.Mimetype = string(mt)
		err = bk.readXML("META-INF/container.xml", &bk.Container)
	}
	if err == nil {
		err = bk.readXML(bk.Container.Rootfile.Path, &bk.Opf)
	}

	for _, mf := range bk.Opf.Manifest {
		if mf.ID == bk.Opf.Spine.Toc {
			err = bk.readXML(bk.filename(mf.Href), &bk.Ncx)
			break
		}
	}

	if err != nil {
		return nil, err
	}

	return &bk, nil
}

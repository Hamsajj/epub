package epub_test

import (
	"io"
	"os"
	"testing"

	"github.com/Hamsajj/epub"
)

func TestOpen(t *testing.T) {
	bk, cancel, err := epub.Open("testdata/alice_adventures.epub")
	t.Cleanup(func() {
		_ = cancel()
	})
	if err != nil {
		t.Fatal(err)
	}

	t.Logf("files: %+v", bk.Files())
	t.Logf("book: %+v", bk)
}

func TestOpenFromReader(t *testing.T) {
	book, err := os.Open("testdata/alice_adventures.epub")
	if err != nil {
		t.Fatal(err)
	}
	defer book.Close()
	bk, err := epub.OpenFromReader(book)
	if err != nil {
		t.Fatal(err)
	}

	if got, want := bk.Opf.Metadata.Title[0], "Alice's Adventures in Wonderland"; got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

func TestOpenBytes(t *testing.T) {
	book, err := os.Open("testdata/alice_adventures.epub")
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		_ = book.Close()
	})
	data, err := io.ReadAll(book)
	if err != nil {
		t.Fatal(err)
	}
	bk, err := epub.OpenFromBytes(data)
	if err != nil {
		t.Fatal(err)
	}

	if got, want := bk.Opf.Metadata.Title[0], "Alice's Adventures in Wonderland"; got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

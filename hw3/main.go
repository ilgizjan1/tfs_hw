package main

import (
	"bytes"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
	"unicode"
	"unicode/utf8"
)

const (
	nameExecutableFile = "main.go"
)

const (
	defaultFrom      = ""
	defaultTo        = ""
	defaultOffset    = 0
	defaultLimit     = 1
	defaultBlockSize = 1024
	defaultConv      = ""
)

const (
	minValueOffset    = 0
	minValueLimit     = 1
	minValueBlockSize = 1
)

const (
	flagFrom       = "from"
	flatTo         = "to"
	flagOffset     = "offset"
	flagLimit      = "limit"
	flagBlockSize  = "block-size"
	flagConv       = "conv"
	flagUpperCase  = "upper_case"
	flagLowerCase  = "lower_case"
	flagTrimSpaces = "trim_spaces"
)

type Options struct {
	from       string
	to         string
	offset     int
	limit      int
	blockSize  int
	conv       string
	upper      bool
	lower      bool
	rightSpace bool
	leftSpace  bool
}

func ParseFlags() (*Options, error) {
	var opts Options
	flag.StringVar(&opts.from, flagFrom, defaultFrom, "file to read. by default - stdin")
	flag.StringVar(&opts.to, flatTo, defaultTo, "file to write. by default - stdout")
	flag.IntVar(&opts.offset, flagOffset, defaultOffset, "number of bytes to skip. by default - zero")
	flag.IntVar(&opts.limit, flagLimit, defaultLimit, "maximum number of bytes to read. by default - read all")
	flag.IntVar(&opts.blockSize, flagBlockSize, defaultBlockSize, "byte size when reading and writing. by default - 1024")
	flag.StringVar(&opts.conv, flagConv, defaultConv, "some of possible converts text. by default - without them")

	flag.Parse()
	if err := opts.Validate(); err != nil {
		return nil, fmt.Errorf("validating options: %w", err)
	}
	return &opts, nil
}

func (opts *Options) Validate() error {
	if opts.to == nameExecutableFile {
		return fmt.Errorf("can not change executable file")
	}
	if opts.blockSize < minValueBlockSize {
		return fmt.Errorf("block-size less than zero")
	}
	if opts.offset < minValueOffset {
		return fmt.Errorf("offset less than zero")
	}
	if opts.limit < minValueLimit {
		return fmt.Errorf("limit less than zero")
	}
	if opts.conv == defaultConv {
		return nil
	}
	convSlice := strings.Split(opts.conv, ",")
	for _, val := range convSlice {
		switch val {
		case flagLowerCase:
			opts.lower = true
		case flagUpperCase:
			opts.upper = true
		case flagTrimSpaces:
			opts.rightSpace = true
			opts.leftSpace = true
		default:
			return fmt.Errorf("incorrect conv format")
		}
	}
	if opts.lower && opts.upper {
		return fmt.Errorf("upper and lower can not live together in conv")
	}
	return nil
}

type ByteConverter interface {
	Convert([]byte) []byte
}

type UpperCaseConverter struct{}

func (c *UpperCaseConverter) Convert(b []byte) []byte {
	return bytes.ToUpper(b)
}

type LowerCaseConverter struct{}

func (c *LowerCaseConverter) Convert(b []byte) []byte {
	return bytes.ToLower(b)
}

type OffsetConverter struct{ offset *int }

func (c *OffsetConverter) Convert(byteSlice []byte) []byte {
	if *c.offset > len(byteSlice) {
		*c.offset -= len(byteSlice)
		return nil
	}
	res := byteSlice[*c.offset:]
	*c.offset = 0
	return res
}

type LimitConverter struct{ limit *int }

func (c *LimitConverter) Convert(byteSlice []byte) []byte {
	if *c.limit > len(byteSlice) {
		*c.limit -= len(byteSlice)
		return byteSlice
	}
	limit := *c.limit
	*c.limit = 0
	return byteSlice[:limit]
}

type LeftSpaceConverter struct{ foundToTrim bool }

func (c *LeftSpaceConverter) Convert(byteSlice []byte) []byte {
	if c.foundToTrim {
		res := bytes.TrimLeftFunc(byteSlice, unicode.IsSpace)
		c.foundToTrim = len(res) == 0
		return res
	}
	return byteSlice
}

type RightSpaceConverter struct {
	spaces []byte
}

func (c *RightSpaceConverter) Convert(byteSlice []byte) []byte {
	res := make([]byte, 0, len(c.spaces)+len(byteSlice))
	res = append(res, c.spaces...)
	res = append(res, byteSlice...)
	n := len(bytes.TrimRightFunc(byteSlice, unicode.IsSpace))
	switch {
	case n < len(byteSlice):
		c.spaces = append(c.spaces, byteSlice[n:]...)
	default:
		c.spaces = make([]byte, 0)
	}
	return bytes.TrimRightFunc(res, unicode.IsSpace)
}

type Processer struct {
	converters []ByteConverter
}

func (opts *Options) newProcesser(rr *ReadWriter) Processer {
	convRaw := make([]ByteConverter, 0, 5)
	convRaw = append(convRaw,
		&OffsetConverter{offset: &rr.offset},
		&LimitConverter{limit: &rr.limit},
		&LeftSpaceConverter{foundToTrim: true},
		&RightSpaceConverter{spaces: make([]byte, 0)},
		&UpperCaseConverter{}, &LowerCaseConverter{},
	)
	conv := make([]ByteConverter, 0)
	if opts.offset != 0 {
		conv = append(conv, convRaw[0])
	}
	if opts.limit != 1 {
		conv = append(conv, convRaw[1])
	}
	if opts.leftSpace {
		conv = append(conv, convRaw[2], convRaw[3])
	}
	if opts.upper {
		conv = append(conv, convRaw[4])
	}
	if opts.lower {
		conv = append(conv, convRaw[5])
	}
	return Processer{converters: conv}
}

func (p *Processer) process(inputBytes []byte) []byte {
	for _, c := range p.converters {
		inputBytes = c.Convert(inputBytes)
	}
	return inputBytes
}

type ReadWriter struct {
	reader    io.Reader
	writer    io.Writer
	limit     int
	offset    int
	blockSize int
}

func (opts *Options) NewReadWriter(reader io.Reader, writer io.Writer) *ReadWriter {
	return &ReadWriter{
		reader:    reader,
		writer:    writer,
		limit:     opts.limit,
		offset:    opts.offset,
		blockSize: opts.blockSize,
	}
}

func (r *ReadWriter) copy(opts *Options) error {
	processer := opts.newProcesser(r)
	data := make([]byte, r.blockSize)
	buf := make([]byte, 0, r.blockSize)
	for r.limit != 0 {
		n, err := r.reader.Read(data)
		if err == io.EOF {
			break
		}
		if err != nil {
			return fmt.Errorf("unable read fromFile %w", err)
		}
		buf = append(buf, data[:n]...)
		if !utf8.Valid(buf) {
			continue
		}
		buf = processer.process(buf)
		if err := Write(r.writer, buf, r.blockSize); err != nil {
			return fmt.Errorf("unable to write the fromFile %w", err)
		}
		buf = make([]byte, 0)
	}
	if r.offset != 0 {
		return fmt.Errorf("fromFile size is less than offset")
	}
	return nil
}

func newReader(nameFile string) (*os.File, error) {
	if nameFile == defaultFrom {
		return os.Stdin, nil
	}
	return os.Open(nameFile)
}

func newWriter(nameFile string) (*os.File, error) {
	if nameFile == defaultTo {
		return os.Stdout, nil
	}
	return os.Create(nameFile)
}

func Write(writer io.Writer, data []byte, blockSize int) error {
	left, right := 0, blockSize
	for left < len(data) {
		if right > len(data) {
			right = len(data)
		}
		_, err := writer.Write(data[left:right])
		if err != nil {
			return err
		}
		left += blockSize
		right += blockSize
	}
	return nil
}

func main() {
	opts, err := ParseFlags()
	if err != nil {
		fmt.Fprintln(os.Stderr, "error in entering flags %w", err)
		log.Fatal(1)
	}
	fromFile, err := newReader(opts.from)
	if err != nil {
		fmt.Fprintln(os.Stderr, "newReader %w", err)
		log.Fatal(1)
	}
	defer fromFile.Close()
	toFile, err := newWriter(opts.to)
	if err != nil {
		fmt.Fprintln(os.Stderr, "newWriter %w", err)
		log.Fatal(1)
	}
	defer toFile.Close()
	readWriter := opts.NewReadWriter(fromFile, toFile)
	err = readWriter.copy(opts)
	if err != nil {
		fmt.Fprintln(os.Stderr, "error in Reading %w", err)
		log.Fatal(1)
	}
}

package gma

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"encoding/json"
	"io"
	"os"
)

func OpenGMA(reader io.ReadSeekCloser) (GMA, error) {
	buf := [4]byte{}

	// Read file signature
	if _, err := reader.Read(buf[:4]); err != nil {
		return nil, err
	}

	// Check if the file has the signature 'GMAD'
	if binary.LittleEndian.Uint32(buf[:4]) != 0x44414D47 {
		return nil, ErrInvalidSignature
	}

	// Skip 18 bytes
	if _, err := reader.Seek(18, io.SeekCurrent); err != nil {
		return nil, err
	}

	bufferedReader := bufio.NewReader(reader)

	name, err := bufferedReader.ReadBytes('\x00')
	if err != nil {
		return nil, err
	}

	desc, err := bufferedReader.ReadBytes('\x00')
	if err != nil {
		return nil, err
	}

	author, err := bufferedReader.ReadBytes('\x00')
	if err != nil {
		return nil, err
	}

	meta := metadata{}

	if err := json.NewDecoder(bytes.NewReader(desc[:len(desc)-1])).Decode(&meta); err != nil {
		meta.Description = string(desc[:len(desc)-1])
		meta.Type = "unknown"
	}

	ad := &gma{
		stream: reader,

		name:   string(name[:len(name)-1]),
		author: string(author[:len(author)-1]),
		meta:   &meta,

		pathMap: make(map[string]*entry),
	}

	if _, err := bufferedReader.Discard(4); err != nil {
		return nil, err
	}

	var (
		offs  int64
		start int64 = 26 + 4 + int64(len(name)+len(desc)+len(author))
	)

	// Read the file
	for {
		if _, err := io.ReadFull(bufferedReader, buf[:4]); err != nil {
			return nil, err
		}
		num := binary.LittleEndian.Uint32(buf[:4])

		if num == 0 {
			break
		}

		path, err := bufferedReader.ReadBytes('\x00')
		if err != nil {
			return nil, err
		}

		if _, err := io.ReadFull(bufferedReader, buf[:4]); err != nil {
			return nil, err
		}
		size := binary.LittleEndian.Uint32(buf[:4])

		if _, err := bufferedReader.Discard(8); err != nil {
			return nil, err
		}

		e := &entry{
			parent: ad,

			name: string(path[:len(path)-1]),
			size: size,
			offs: offs,
		}
		ad.files = append(ad.files, e)
		ad.pathMap[e.name] = e

		offs += int64(size)

		start += 16 + int64(len(path))
	}
	ad.start = start

	return ad, nil
}

func OpenFile(path string) (GMA, error) {
	fs, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	return OpenGMA(fs)
}

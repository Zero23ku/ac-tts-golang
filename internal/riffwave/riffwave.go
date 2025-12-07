package riffwave

import (
	"bytes"
	"encoding/base64"
	"encoding/binary"
)

type RIFFWAVE struct {
	Data    []int16 // soporta 16-bit samples
	Wav     []byte  // archivo WAV generado
	DataURI string  // data URI con base64
	Header  WAVHeader
}

// WAVHeader contiene los metadatos del archivo WAV
type WAVHeader struct {
	ChunkID       [4]byte
	ChunkSize     uint32
	Format        [4]byte
	SubChunk1ID   [4]byte
	SubChunk1Size uint32
	AudioFormat   uint16
	NumChannels   uint16
	SampleRate    uint32
	ByteRate      uint32
	BlockAlign    uint16
	BitsPerSample uint16
	SubChunk2ID   [4]byte
	SubChunk2Size uint32
}

func (r *RIFFWAVE) Make() {
	r.Header.BlockAlign = (r.Header.NumChannels * r.Header.BitsPerSample) >> 3
	r.Header.ByteRate = r.Header.SampleRate * uint32(r.Header.BlockAlign)
	r.Header.SubChunk2Size = uint32(len(r.Data)) * uint32(r.Header.BlockAlign)
	r.Header.ChunkSize = 36 + r.Header.SubChunk2Size

	var buf bytes.Buffer
	// Escribir cabecera RIFF
	binary.Write(&buf, binary.LittleEndian, r.Header.ChunkID)
	binary.Write(&buf, binary.LittleEndian, r.Header.ChunkSize)
	binary.Write(&buf, binary.LittleEndian, r.Header.Format)
	binary.Write(&buf, binary.LittleEndian, r.Header.SubChunk1ID)
	binary.Write(&buf, binary.LittleEndian, r.Header.SubChunk1Size)
	binary.Write(&buf, binary.LittleEndian, r.Header.AudioFormat)
	binary.Write(&buf, binary.LittleEndian, r.Header.NumChannels)
	binary.Write(&buf, binary.LittleEndian, r.Header.SampleRate)
	binary.Write(&buf, binary.LittleEndian, r.Header.ByteRate)
	binary.Write(&buf, binary.LittleEndian, r.Header.BlockAlign)
	binary.Write(&buf, binary.LittleEndian, r.Header.BitsPerSample)
	binary.Write(&buf, binary.LittleEndian, r.Header.SubChunk2ID)
	binary.Write(&buf, binary.LittleEndian, r.Header.SubChunk2Size)

	// Escribir samples (16-bit PCM)
	for _, sample := range r.Data {
		binary.Write(&buf, binary.LittleEndian, sample)
	}

	r.Wav = buf.Bytes()
	r.DataURI = "data:audio/wav;base64," + base64.StdEncoding.EncodeToString(r.Wav)
}

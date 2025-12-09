package animalese

import (
	"bytes"
	"encoding/binary"
	"errors"
	"io"
	"os"
	"strings"
	"unicode"

	"ac-tts/internal/logging"
)

type Animalese struct {
	letterLibrary []byte
}

func shortenWord(str string) string {
	if len(str) > 1 {
		return string(str[0]) + string(str[len(str)-1])
	}
	return str
}

func (a *Animalese) AnimaleseFunc(script string, shorten bool, pitch float64) []byte {
	processedScript := script
	if shorten {
		words := strings.FieldsFunc(script, func(r rune) bool {
			return !unicode.IsLetter(r)
		})
		shortened := make([]string, len(words))
		for i, w := range words {
			shortened[i] = shortenWord(w)
		}
		processedScript = strings.Join(shortened, "")
	}

	sampleFreq := 44100
	libraryLetterSecs := 0.15
	librarySamplesPerLetter := int(libraryLetterSecs * float64(sampleFreq))
	outputLetterSecs := 0.075
	outputSamplesPerLetter := int(outputLetterSecs * float64(sampleFreq))

	data := make([]byte, len(processedScript)*outputSamplesPerLetter)

	for cIndex, r := range strings.ToUpper(processedScript) {
		if r >= 'A' && r <= 'Z' {
			libraryLetterStart := librarySamplesPerLetter * (int(r) - int('A'))
			for i := 0; i < outputSamplesPerLetter; i++ {
				data[cIndex*outputSamplesPerLetter+i] =
					a.letterLibrary[libraryLetterStart+int(float64(i)*pitch)]
			}
		} else {
			for i := 0; i < outputSamplesPerLetter; i++ {
				data[cIndex*outputSamplesPerLetter+i] = 127
			}
		}
	}

	var buf bytes.Buffer
	buf.WriteString("RIFF")
	binary.Write(&buf, binary.LittleEndian, int32(36+len(data)))
	buf.WriteString("WAVEfmt ")
	binary.Write(&buf, binary.LittleEndian, int32(16))
	binary.Write(&buf, binary.LittleEndian, int16(1))
	binary.Write(&buf, binary.LittleEndian, int16(1))
	binary.Write(&buf, binary.LittleEndian, int32(sampleFreq))
	binary.Write(&buf, binary.LittleEndian, int32(sampleFreq))
	binary.Write(&buf, binary.LittleEndian, int16(1))
	binary.Write(&buf, binary.LittleEndian, int16(8))
	buf.WriteString("data")
	binary.Write(&buf, binary.LittleEndian, int32(len(data)))
	buf.Write(data)

	return buf.Bytes()
}

func NewAnimalese(lettersFile string, onload func()) (*Animalese, error) {
	file, err := os.Open(lettersFile)
	if err != nil {
		logging.CreateLog("couldn't open letters file", err)
		return nil, err
	}
	defer file.Close()

	header := make([]byte, 44)
	_, err = io.ReadFull(file, header)
	if err != nil {
		logging.CreateLog("couldn't read wav header file", err)
		return nil, err
	}

	if string(header[0:4]) != "RIFF" || string(header[8:12]) != "WAVE" {
		return nil, errors.New("archivo no es WAV válido")
	}

	subChunk2Size := binary.LittleEndian.Uint32(header[40:44])

	data := make([]byte, subChunk2Size)
	_, err = io.ReadFull(file, data)
	if err != nil {
		logging.CreateLog("couldn't read wav body file", err)
		return nil, err
	}

	a := &Animalese{letterLibrary: data}
	onload()
	return a, nil
}

func NewAnimaleseFromBytes(wavData []byte, onload func()) (*Animalese, error) {
	if len(wavData) < 44 {
		return nil, errors.New("archivo demasiado pequeño para ser WAV válido")
	}

	header := wavData[:44]

	if string(header[0:4]) != "RIFF" || string(header[8:12]) != "WAVE" {
		return nil, errors.New("archivo no es WAV válido")
	}

	subChunk2Size := binary.LittleEndian.Uint32(header[40:44])

	if int(subChunk2Size) > len(wavData)-44 {
		return nil, errors.New("tamaño inconsistente en WAV")
	}

	data := wavData[44 : 44+subChunk2Size]

	a := &Animalese{letterLibrary: data}
	onload()
	return a, nil
}

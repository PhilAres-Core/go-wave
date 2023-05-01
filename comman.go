package wave

import (
	"errors"
	"log"
	"reflect"
)

type waveHeader struct {
	SampleRate     int
	NumChannels    int
	BitDepth       int
	ByteRate       int
	AudioLen       int
	chunkID        []byte // 4
	chunkSize      []byte // 4
	format         []byte // 4
	subChunk1ID    []byte // 4
	subChunk1Size  []byte // 4
	audioFormat    []byte // 2
	numChannels    []byte // 2
	sampleRate     []byte // 4
	byteRate       []byte // 4
	blockAlign     []byte // 2
	bitesPerSample []byte // 2
	subChunk2ID    []byte // 4
	subChunk2Size  []byte // 4
}

func isWave(data []byte) bool {
	if !reflect.DeepEqual(data[0:4], []byte{'R', 'I', 'F', 'F'}) || !reflect.DeepEqual(data[8:12], []byte{'W', 'A', 'V', 'E'}) || !reflect.DeepEqual(data[12:16], []byte{'f', 'm', 't', ' '}) || !reflect.DeepEqual(data[36:40], []byte{'d', 'a', 't', 'a'}) {
		return false
	}
	return true
}

func _byte2header(data []byte) (wave *waveHeader) {
	if !isWave(data) {
		log.Fatalln("the format of the file is not wave.")
	}
	numChannels := int(data[22]) + int(data[23])<<8
	samPleRate := int(data[24]) + int(data[25])<<8 + int(data[26])<<16 + int(data[27])<<24
	byteRate := int(data[28]) + int(data[29])<<8 + int(data[30])<<16 + int(data[31])<<24
	audioLen := int(data[40]) + int(data[41])<<8 + int(data[42])<<16 + int(data[43])<<24
	bitDepth := byteRate * 8 / samPleRate / numChannels
	wave = &waveHeader{
		SampleRate:     samPleRate,
		NumChannels:    numChannels,
		BitDepth:       bitDepth,
		ByteRate:       byteRate,
		AudioLen:       audioLen,
		chunkID:        data[0:4],
		chunkSize:      data[4:8],
		format:         data[8:12],
		subChunk1ID:    data[12:16],
		subChunk1Size:  data[16:20],
		audioFormat:    data[20:22],
		numChannels:    data[22:24],
		sampleRate:     data[24:28],
		byteRate:       data[28:32],
		blockAlign:     data[32:34],
		bitesPerSample: data[34:36],
		subChunk2ID:    data[36:40],
		subChunk2Size:  data[40:44],
	}
	return
}

func _format2header(numChannels, sampleRate, bitDepth int) (wave *waveHeader) {
	byteRate := sampleRate * numChannels * bitDepth / 8
	totalAudioLen := 0
	totalDataLen := totalAudioLen + 36
	perSampleSize := numChannels * bitDepth / 8
	wave = &waveHeader{
		SampleRate:     sampleRate,
		NumChannels:    numChannels,
		BitDepth:       bitDepth,
		ByteRate:       byteRate,
		AudioLen:       totalAudioLen,
		chunkID:        []byte{'R', 'I', 'F', 'F'},
		chunkSize:      []byte{byte(totalDataLen & 0xff), byte((totalDataLen >> 8) & 0xff), byte((totalDataLen >> 16) & 0xff), byte((totalDataLen >> 24) & 0xff)},
		format:         []byte{'W', 'A', 'V', 'E'},
		subChunk1ID:    []byte{'f', 'm', 't', ' '},
		subChunk1Size:  []byte{16, 0, 0, 0},
		audioFormat:    []byte{1, 0},
		numChannels:    []byte{byte(numChannels & 0xff), byte((numChannels >> 8) & 0xff)},
		sampleRate:     []byte{byte(sampleRate & 0xff), byte((sampleRate >> 8) & 0xff), byte((sampleRate >> 16) & 0xff), byte((sampleRate >> 24) & 0xff)},
		byteRate:       []byte{byte(byteRate & 0xff), byte((byteRate >> 8) & 0xff), byte((byteRate >> 16) & 0xff), byte((byteRate >> 24) & 0xff)},
		blockAlign:     []byte{byte(perSampleSize & 0xff), byte((perSampleSize >> 8) & 0xff)},
		bitesPerSample: []byte{byte(bitDepth & 0xff), byte((bitDepth >> 8) & 0xff)},
		subChunk2ID:    []byte{'d', 'a', 't', 'a'},
		subChunk2Size:  []byte{byte(totalAudioLen & 0xff), byte((totalAudioLen >> 8) & 0xff), byte((totalAudioLen >> 16) & 0xff), byte((totalAudioLen >> 24) & 0xff)},
	}
	return
}

// mode:mode=0 >> only read, mode=1 >> read write
// args: [[]byte],["numChannels","sampleRate","bitDepth"]
func newWaveHandler(mode int, args ...interface{}) (wave *waveHeader, err error) {
	if mode == 0 {
		if len(args) != 1 {
			err = errors.New("the count of arguments is error")
			return
		}
		arg := args[0]
		switch arg.(type) {
		case []byte:
			if len(arg.([]byte)) != 44 {
				err = errors.New("the length of the argument is error")
				return
			}
			wave = _byte2header(arg.([]byte))
		default:
			err = errors.New("the type argument is error")
			return
		}
	}
	if mode == 1 {
		var (
			numChannels = 0
			samPleRate  = 0
			bitDepth    = 0
		)
		if len(args) != 3 {
			err = errors.New("the count of arguments is error")
			return
		}
		for ind, arg := range args {
			switch arg.(type) {
			case int:
				if ind == 0 {
					numChannels = arg.(int)
				}
				if ind == 1 {
					samPleRate = arg.(int)
				}
				if ind == 2 {
					bitDepth = arg.(int)
				}
			case float64:
				if ind == 0 {
					numChannels = arg.(int)
				}
				if ind == 1 {
					samPleRate = arg.(int)
				}
				if ind == 2 {
					bitDepth = arg.(int)
				}
			default:
				err = errors.New("the type argument is error")
				return
			}
		}
		wave = _format2header(numChannels, samPleRate, bitDepth)
	}
	return
}

func (wave *waveHeader) getHeader(audioLen int) (header []byte) {
	if audioLen > 0 {
		wave.AudioLen = wave.AudioLen + audioLen
		totalDataLen := wave.AudioLen + 36
		wave.chunkSize = []byte{byte(totalDataLen), byte((totalDataLen >> 8) & 0xff), byte((totalDataLen >> 16) & 0xff), byte((totalDataLen >> 24) & 0xff)}
		wave.subChunk2Size = []byte{byte(wave.AudioLen), byte((wave.AudioLen >> 8) & 0xff), byte((wave.AudioLen >> 16) & 0xff), byte((wave.AudioLen >> 24) & 0xff)}
	}
	header = make([]byte, 0)
	header = append(header, wave.chunkID...)
	header = append(header, wave.chunkSize...)
	header = append(header, wave.format...)
	header = append(header, wave.subChunk1ID...)
	header = append(header, wave.subChunk1Size...)
	header = append(header, wave.audioFormat...)
	header = append(header, wave.numChannels...)
	header = append(header, wave.sampleRate...)
	header = append(header, wave.byteRate...)
	header = append(header, wave.blockAlign...)
	header = append(header, wave.bitesPerSample...)
	header = append(header, wave.subChunk2ID...)
	header = append(header, wave.subChunk2Size...)
	return
}

func (wave *waveHeader) getNumChannels() int {
	return wave.NumChannels
}

func (wave *waveHeader) getSampleRate() int {
	return wave.SampleRate
}

func (wave *waveHeader) getBitDepth() int {
	return wave.BitDepth
}

func (wave *waveHeader) getByteRate() int {
	return wave.ByteRate
}

func (wave *waveHeader) getAudioLen() int {
	return wave.AudioLen
}

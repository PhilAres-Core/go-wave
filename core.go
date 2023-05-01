package wave

import (
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"regexp"
)

type waveHandler struct {
	file   *os.File
	header *waveHeader
	offset int64
}

func openWave(waveName, mode string) (wave *waveHandler, err error) {
	var (
		f       *os.File
		h       *waveHeader
		regC, _ = regexp.Compile(".wav$")
	)
	if !regC.MatchString(waveName) {
		log.Fatalln("This file address is incorrect")
	}
	if mode == "RDONLY" {
		var (
			data = make([]byte, 44)
			row  = 0
		)
		if f, err = os.OpenFile(waveName, os.O_RDONLY, 0644); err != nil {
			return
		}
		if row, err = f.Read(data); err != nil {
			return
		}
		if row != 44 {
			err = errors.New(fmt.Sprintf("%v: Format is not wave.", f.Name()))
			return
		}
		if h, err = newWaveHandler(0, data); err != nil {
			_ = f.Close()
			return
		}
	}
	if mode == "WRONLY" {
		if f, err = os.OpenFile(waveName, os.O_WRONLY|os.O_CREATE, 0644); err != nil {
			return
		}
	}
	if mode == "RDWR" {
		var (
			status os.FileInfo
			data   = make([]byte, 44)
		)
		if f, err = os.OpenFile(waveName, os.O_RDWR|os.O_CREATE, 0644); err != nil {
			return
		}
		if status, err = f.Stat(); err != nil {
			return
		}
		if status.Size() >= 44 {
			if _, err = f.Read(data); err != nil {
				return
			}
			if h, err = newWaveHandler(0, data); err != nil {
				_ = f.Close()
				return
			}
		}
	}
	if mode != "RDONLY" && mode != "WRONLY" && mode != "RDWR" {
		log.Fatalln("Access is denied.The mode is error")
	}
	wave = &waveHandler{
		file:   f,
		header: h,
		offset: 0,
	}
	return
}

func (wave *waveHandler) setWaveFormat(numChannels, sampleRate, bitDepth int) (err error) {
	defer func() {
		if e := recover(); e != nil {
			_ = wave.closeWave()
			log.Fatalln(e)
		}
	}()
	if wave.header, err = newWaveHandler(1, numChannels, sampleRate, bitDepth); err != nil {
		return
	}
	header := wave.header.getHeader(0)
	if _, err = wave.file.WriteAt(header, 0); err != nil {
		return
	}
	return err
}

func (wave *waveHandler) write(data []byte) (n int, err error) {
	defer func() {
		if e := recover(); e != nil {
			_ = wave.closeWave()
			log.Fatalln(e)
		}
	}()
	var (
		header []byte
		ret    int64
	)
	if wave.header == nil {
		log.Fatalln(fmt.Sprintf("write %v: Access is denied.Please set the wave format first.", wave.file.Name()))
	}
	header = wave.header.getHeader(len(data))
	if _, err = wave.file.WriteAt(header, 0); err != nil {
		return
	}
	if ret, err = wave.file.Seek(0, io.SeekEnd); err != nil {
		return
	}
	if n, err = wave.file.WriteAt(data, ret); err != nil {
		return
	}
	return
}

func (wave *waveHandler) read(n int) (data []byte, err error) {
	defer func() {
		if e := recover(); e != nil {
			_ = wave.closeWave()
			log.Fatalln(e)
		}
	}()
	var (
		ret int64
		num int
	)
	if n == -1 {
		if ret, err = wave.file.Seek(0, io.SeekEnd); err != nil {
			return
		}
		data = make([]byte, ret)
		if num, err = wave.file.ReadAt(data, 0); err != nil {
			return
		}
		data = data[0:num]
		return
	}
	if n > 0 {
		data = make([]byte, n)
		offset := wave.offset
		if num, err = wave.file.ReadAt(data, offset); err != nil {
			return
		}
		data = data[0:num]
		wave.offset = offset + int64(n)
		return
	}
	if n < 0 {
		log.Fatalln(fmt.Sprintf("read %v: specific parameter values are incorrect.", wave.file.Name()))
		return
	}
	return
}

func (wave *waveHandler) getNumChannels() (n int) {
	defer func() {
		if e := recover(); e != nil {
			_ = wave.closeWave()
			log.Fatalln(e)
		}
	}()
	if wave.header != nil {
		n = wave.header.getNumChannels()
	} else {
		log.Fatalln("Access is denied.")
		return
	}
	return
}

func (wave *waveHandler) getSampleRate() (n int) {
	defer func() {
		if e := recover(); e != nil {
			_ = wave.closeWave()
			log.Fatalln(e)
		}
	}()
	if wave.header != nil {
		n = wave.header.getSampleRate()
	} else {
		log.Fatalln("Access is denied.")
		return
	}
	return
}

func (wave *waveHandler) getByteRate() (n int) {
	defer func() {
		if e := recover(); e != nil {
			_ = wave.closeWave()
			log.Fatalln(e)
		}
	}()
	if wave.header != nil {
		n = wave.header.getByteRate()
	} else {
		log.Fatalln("Access is denied.")
		return
	}
	return
}

func (wave *waveHandler) getBitDepth() (n int) {
	defer func() {
		if e := recover(); e != nil {
			_ = wave.closeWave()
			log.Fatalln(e)
		}
	}()
	if wave.header != nil {
		n = wave.header.getBitDepth()
	} else {
		log.Fatalln("Access is denied.")
		return
	}
	return
}

func (wave *waveHandler) getAudioLen() (n int) {
	defer func() {
		if e := recover(); e != nil {
			_ = wave.closeWave()
			log.Fatalln(e)
		}
	}()
	if wave.header != nil {
		n = wave.header.getAudioLen()
	} else {
		log.Fatalln("Access is denied.")
		return
	}
	return
}

func (wave *waveHandler) getWaveFormat() (info map[string]int) {
	defer func() {
		if e := recover(); e != nil {
			_ = wave.closeWave()
			log.Fatalln(e)
		}
	}()
	if wave.header != nil {
		info = make(map[string]int)
		info["numChannels"] = wave.header.getNumChannels()
		info["sampleRate"] = wave.header.getSampleRate()
	} else {
		log.Fatalln("Access is denied.")
		return
	}
	return
}

func (wave *waveHandler) closeWave() (err error) {
	if wave.file != nil {
		_ = wave.file.Close()
	} else {
		err = errors.New("the file is closed")
	}
	return
}

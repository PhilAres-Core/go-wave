package wave

type Wave struct {
	w *waveHandler
}

// Open TODO: mode="WRONLY" -> 只读
// Open TODO: mode="RDONLY" -> 只读
// Open TODO: mode=RDWR -> "读写"
func Open(waveName, mode string) (wave *Wave, err error) {
	var (
		w *waveHandler
	)
	if w, err = openWave(waveName, mode); err != nil {
		return
	}
	wave = &Wave{
		w: w,
	}
	return
}

// Close TODO: 关闭wave文件
func (wave *Wave) Close() error {
	return wave.w.closeWave()
}

// SetWaveFormat TODO: 设定wave文件的数据格式
// SetWaveFormat TODO: numChannels 音频的声道数
// SetWaveFormat TODO: sampleRate  音频的采样率
// SetWaveFormat TODO: bitDepth    音频的采样精度
func (wave *Wave) SetWaveFormat(numChannels, sampleRate, bitDepth int) error {
	return wave.w.setWaveFormat(numChannels, sampleRate, bitDepth)
}

// Write TODO: 写入数据
func (wave *Wave) Write(data []byte) (int, error) {
	return wave.w.write(data)
}

// Read TODO: 读取数据，n为读取数据的大小
func (wave *Wave) Read(n int) ([]byte, error) {
	return wave.w.read(n)
}

func (wave *Wave) GetNumChannels() int {
	return wave.w.getNumChannels()
}

func (wave *Wave) GetSampleRate() int {
	return wave.w.getSampleRate()
}

func (wave *Wave) GetByteRate() int {
	return wave.w.getByteRate()
}

func (wave *Wave) GetBitDepth() int {
	return wave.w.getBitDepth()
}

func (wave *Wave) GetAudioLen() int {
	return wave.w.getAudioLen()
}

func (wave *Wave) GetWaveInfo() map[string]int {
	return wave.w.getWaveFormat()
}

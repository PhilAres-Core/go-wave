package wave

import (
	"log"
	"testing"
)

//func TestOpen(t *testing.T) {
//	file, err := Open("test.wav", "RDWR")
//	if err != nil {
//		log.Fatalln(err)
//	}
//	defer func() {
//		_ = file.Close()
//	}()
//	_ = file.SetWaveFormat(1, 16000, 16)
//	data := make([]byte, 32000)
//	for i := 0; i < 60*60*24*10; i++ {
//		_, e := file.Write(data)
//		if e != nil {
//			log.Fatalln(e)
//		}
//	}
//}

func BenchmarkWave_Write(b *testing.B) {
	file, err := Open("test.wav", "RDWR")
	if err != nil {
		log.Fatalln(err)
	}
	defer func() {
		_ = file.Close()
	}()
	b.ReportAllocs()
	b.ResetTimer()
	_ = file.SetWaveFormat(1, 16000, 16)
	data := make([]byte, 32000)
	for i := 0; i < b.N; i++ {
		_, e := file.Write(data)
		if e != nil {
			log.Fatalln(e)
		}
	}
}

func BenchmarkWave_Read(b *testing.B) {
	file, err := Open("test.wav", "RDWR")
	if err != nil {
		log.Fatalln(err)
	}
	defer func() {
		_ = file.Close()
	}()
	b.ReportAllocs()
	b.ResetTimer()
	_ = file.SetWaveFormat(1, 16000, 16)
	//data := make([]byte, 32000)
	for i := 0; i < b.N; i++ {
		_, e := file.Read(32000)
		if e != nil {
			log.Fatalln(e)
		}
	}
}

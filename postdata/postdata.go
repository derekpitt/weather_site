package postdata

import (
	"crypto/hmac"
	"crypto/sha512"
	"encoding/base64"
	"fmt"
	"github.com/derekpitt/weather_station/loop2packet"
	"time"
)

type PostData struct {
	Sig    string
	Time   time.Time
	Sample loop2packet.Loop2Packet
}

var encoder = base64.StdEncoding

func checkMAC(message, messageMAC, key []byte) bool {
	mac := hmac.New(sha512.New, key)
	mac.Write(message)
	expectedMAC := mac.Sum(nil)
	return hmac.Equal(messageMAC, expectedMAC)
}

func makeMAC(message, key []byte) []byte {
	mac := hmac.New(sha512.New, key)
	mac.Write(message)
	return mac.Sum(nil)
}

func makeSig(data, key string) string {
	return encoder.EncodeToString((makeMAC([]byte(data), []byte(key))))
}

func checkSig(data, dataSig, key string) bool {
	dataDecoded, _ := encoder.DecodeString(data)
	dataSigDecoded, _ := encoder.DecodeString(dataSig)
	return checkMAC(dataDecoded, dataSigDecoded, []byte(key))
}

func prepareData(time time.Time, packet loop2packet.Loop2Packet) string {
	return fmt.Sprintf("%s%v", time.Unix(), packet)
}

func (pd *PostData) sign(time time.Time, packet loop2packet.Loop2Packet, key string) {
	pd.Sig = makeSig(prepareData(time, packet), key)
}

func makeEmptyPostData(time time.Time, packet loop2packet.Loop2Packet) PostData {
	return PostData{
		"",
		time,
		packet,
	}
}

func NewData(time time.Time, packet loop2packet.Loop2Packet, key string) PostData {
	data := makeEmptyPostData(time, packet)
	data.sign(time, packet, key)

	return data
}

func VerifyData(pd PostData, key string) bool {
	data := makeEmptyPostData(pd.Time, pd.Sample)
	data.sign(pd.Time, pd.Sample, key)

	return data.Sig == pd.Sig
}

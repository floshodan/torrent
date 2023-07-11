package bencode

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"log"
	"strconv"
)

type Decoder struct {
	buf *bytes.Buffer
}

func NewDecoder(r io.Reader) *Decoder {
	b := new(bytes.Buffer)
	b.ReadFrom(r)
	return &Decoder{b}
}

func (d *Decoder) ReadNxt() string {
	byt, err := d.buf.ReadByte()
	d.buf.UnreadByte()
	if err == io.EOF {
		fmt.Println("finished reading the file")
	}
	if err != nil {
		log.Fatalf("Error %s reading file", err)
	}

	if err != nil {
		return "Error"
	}
	return string(byt)
}

func (d *Decoder) Decode() (map[string]interface{}, error) {

	if firstByte, err := d.buf.ReadByte(); err != nil {
		return make(map[string]interface{}), nil
	} else if firstByte != 'd' {
		return nil, errors.New("Not a real .torrent file Torrents must start with a d directive")
	}
	return d.parseDict()

}

func (d *Decoder) parser() (interface{}, error) {
	nxt, err := d.buf.ReadByte()
	d.buf.UnreadByte()
	if err != nil {
		return nil, err
	}
	switch nxt {
	case 'i':
		d.buf.ReadByte()
		return d.parseInt()
	case 'd':
		d.buf.ReadByte()
		return d.parseDict()
	case 'l':
		d.buf.ReadByte()
		return d.parseList()
	case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
		return d.parseString(), err
	default:
		return d.parseString(), err
	}
}

func (d *Decoder) parseDict() (map[string]interface{}, error) {
	dict := make(map[string]interface{})

	for {
		nxt, err := d.buf.ReadByte()
		d.buf.UnreadByte()
		if err == io.EOF {
			//fmt.Println("finished reading the file")
			break
		}
		if err != nil {
			return nil, err
		}

		//when e is reached then the dicionary ends
		if nxt == 'e' {
			d.buf.ReadByte()
			break
		}
		key := d.parseString()

		value, _ := d.parser()
		dict[string(key)] = value
	}
	return dict, nil
}

func (d *Decoder) parseString() string {

	blen, err := d.buf.ReadBytes(':')
	if err != nil {
		return ""
	}
	x, err := strconv.Atoi(string(blen[:len(blen)-1]))

	bytes := make([]byte, x)
	_, err = d.buf.Read(bytes)
	if err != nil {
		return ""
	}

	return string(bytes)
}

func (d *Decoder) parseList() ([]interface{}, error) {

	list := make([]interface{}, 0)

	for {
		nxt, err := d.buf.ReadByte()
		d.buf.UnreadByte()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalf("Error %s reading file", err)
			break
		}

		//when e is reached then the dicionary ends
		if nxt == 'e' {
			d.buf.ReadByte()
			return list, nil
		}

		value, _ := d.parser()
		list = append(list, value)

	}

	return list, nil
}

func (d *Decoder) parseInt() (int64, error) {
	blen, err := d.buf.ReadBytes('e')
	if err != nil {
		return 0, nil
	}
	x, err := strconv.ParseInt(string(blen[:len(blen)-1]), 10, 64)
	return x, nil

}

func (d *Decoder) expect(expected byte) error {
	b, err := d.buf.ReadByte()
	if err != nil {
		return fmt.Errorf("expected '%c', but reached end of input", expected)
	}
	if b != expected {
		return fmt.Errorf("expected '%c', but got '%c'", expected, b)
	}
	return nil
}

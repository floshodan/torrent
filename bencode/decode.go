package bencode

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"log"
	"strconv"
)

type Decoder struct {
	br  *bufio.Reader
	raw []byte
}

func NewDecoder(r io.Reader) *Decoder {
	b := bufio.NewReader(r)
	return &Decoder{b, []byte{}}
}

func (d *Decoder) ReadNxt() string {
	d.br.UnreadByte()
	byt, err := d.br.Peek(1)
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

	if firstByte, err := d.br.ReadByte(); err != nil {
		return make(map[string]interface{}), nil
	} else if firstByte != 'd' {
		return nil, errors.New("Not a real .torrent file Torrents must start with a d directive")
	}
	return d.parseDict()

}

func (d *Decoder) parser() (interface{}, error) {
	nxt, err := d.br.Peek(1)
	if err != nil {
		return nil, err
	}
	switch nxt[0] {
	case 'i':
		d.br.ReadByte()
		return d.parseInt()
	case 'd':
		d.br.ReadByte()
		return d.parseDict()
	case 'l':
		d.br.ReadByte()
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
		nxt, err := d.br.Peek(1)
		if err == io.EOF {
			//fmt.Println("finished reading the file")
			break
		}
		if err != nil {
			log.Fatalf("Error %s reading file", err)
			break
		}

		//when e is reached then the dicionary ends
		if nxt[0] == 'e' {
			d.br.ReadByte()
			break
		}
		key := d.parseString()

		value, _ := d.parser()
		dict[string(key)] = value
	}
	return dict, nil
}

func (d *Decoder) parseString() string {

	blen, err := d.br.ReadBytes(':')
	if err != nil {
		return ""
	}
	x, err := strconv.Atoi(string(blen[:len(blen)-1]))

	bytes := make([]byte, x)
	_, err = d.br.Read(bytes)
	if err != nil {
		return ""
	}

	return string(bytes)
}

func (d *Decoder) parseList() ([]interface{}, error) {

	list := make([]interface{}, 0)

	for {
		nxt, err := d.br.Peek(1)
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalf("Error %s reading file", err)
			break
		}

		//when e is reached then the dicionary ends
		if nxt[0] == 'e' {
			d.br.ReadByte()
			return list, nil
		}

		value, _ := d.parser()
		list = append(list, value)

	}

	return list, nil
}

func (d *Decoder) parseInt() (int64, error) {
	blen, err := d.br.ReadBytes('e')
	if err != nil {
		return 0, nil
	}
	x, err := strconv.ParseInt(string(blen[:len(blen)-1]), 10, 64)
	return x, nil

}

func (d *Decoder) expect(expected byte) error {
	b, err := d.br.ReadByte()
	if err != nil {
		return fmt.Errorf("expected '%c', but reached end of input", expected)
	}
	if b != expected {
		return fmt.Errorf("expected '%c', but got '%c'", expected, b)
	}
	return nil
}

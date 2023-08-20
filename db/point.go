package db

import (
	"bytes"
	"database/sql/driver"
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"log"
)

// Point represents an x,y coordinate in EPSG:4326 for PostGIS.
type Point [2]float64

func (p *Point) String() string {
	// FIXME Add Postgis first and then SRID=4326;
	return fmt.Sprintf("(%v, %v)", p[0], p[1])
}

// Scan implements the sql.Scanner interface.
func (p *Point) Scan(val interface{}) error {
	log.Println("going over there")
	b, err := hex.DecodeString(string(val.([]uint8)))
	if err != nil {
		return err
	}
	r := bytes.NewReader(b)
	var wkbByteOrder uint8
	if err := binary.Read(r, binary.LittleEndian, &wkbByteOrder); err != nil {
		return err
	}

	var byteOrder binary.ByteOrder
	switch wkbByteOrder {
	case 0:
		byteOrder = binary.BigEndian
	case 1:
		byteOrder = binary.LittleEndian
	default:
		return fmt.Errorf("Invalid byte order %d", wkbByteOrder)
	}

	var wkbGeometryType uint64
	if err := binary.Read(r, byteOrder, &wkbGeometryType); err != nil {
		return err
	}

	if err := binary.Read(r, byteOrder, p); err != nil {
		return err
	}

	return nil
}

// Value impl.
func (p Point) Value() (driver.Value, error) {
	log.Println("going over here")
	return p.String(), nil
}

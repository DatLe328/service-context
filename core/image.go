package core

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"strings"
)

type Image struct {
	Id       int    `json:"id" gorm:"column:id;" db:"id"`
	FileName string `json:"file_name" gorm:"column:file_name;" db:"file_name"`
	Width    int    `json:"width" gorm:"column:width;" db:"width"`
	Height   int    `json:"height" gorm:"column:height;" db:"height"`
	Provider string `json:"provider,omitempty" gorm:"column:provider;" db:"provider"`
}

func (*Image) TableName() string { return "images" }

func (img *Image) Fulfill(domain string) {
	if img == nil {
		return
	}
	img.FileName = fmt.Sprintf("%s/%s", strings.TrimRight(domain, "/"), img.FileName)
}

func (img *Image) URL(domain string) string {
	if img == nil {
		return ""
	}
	return fmt.Sprintf("%s/%s", strings.TrimRight(domain, "/"), img.FileName)
}

func (img *Image) Scan(value interface{}) error {
	if value == nil {
		return nil
	}

	bytes, ok := value.([]byte)
	if !ok {
		return fmt.Errorf("failed to unmarshal Image value: %v", value)
	}

	var temp Image
	if err := json.Unmarshal(bytes, &temp); err != nil {
		return err
	}

	*img = temp
	return nil
}

func (img *Image) Value() (driver.Value, error) {
	if img == nil {
		return nil, nil
	}
	return json.Marshal(img)
}

type Images []Image

func (imgs *Images) Scan(value interface{}) error {
	if value == nil {
		*imgs = []Image{}
		return nil
	}

	bytes, ok := value.([]byte)
	if !ok {
		return fmt.Errorf("failed to unmarshal Images value: %v", value)
	}

	var temp []Image
	if err := json.Unmarshal(bytes, &temp); err != nil {
		return err
	}

	*imgs = temp
	return nil
}

func (imgs *Images) Value() (driver.Value, error) {
	if imgs == nil {
		return nil, nil
	}
	return json.Marshal(imgs)
}

func (imgs *Images) FulfillAll(domain string) {
	if imgs == nil {
		return
	}
	for i := range *imgs {
		(*imgs)[i].Fulfill(domain)
	}
}

func (imgs *Images) IsEmpty() bool {
	return imgs == nil || len(*imgs) == 0
}

func (imgs *Images) First() *Image {
	if imgs.IsEmpty() {
		return nil
	}
	return &(*imgs)[0]
}

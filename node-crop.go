package flownodes

import (
	// blank import to support some image formats
	_ "image/jpeg"
	_ "image/png"
)

type Crop struct {
	frame  Port
	region Port
}

func (cr *Crop) Init(db DB) (err error) {
	cr.frame = Port{Name: PortFrame, Type: PortFrame}
	cr.region = Port{Name: PortRegion, Type: PortRegion}

	return
}

func (cr *Crop) Execute(c Ctx) (out Resp, err error) {

	//if data, ok := c.GetData(PortFrame); ok {}
	return Resp{}, nil
}

func (cr *Crop) GetPorts() (in, out []Port) {
	return in, out
}

/*
Я это вижу так:
- получили из Ctx frame, +
- из frame image, +
image вырезали в image,
сделали новый фрейм из image, без лишней конвертации в байты.
это будет НАМНОГО оптимальнее, и кода меньше станет раза в 2.

*/

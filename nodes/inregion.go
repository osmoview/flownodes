package nodes

import (
	"github.com/osmoview/flowdata"
	"github.com/osmoview/flownodes"
)

type InRegion struct {
	Region string

	storage flownodes.Storage
}

type RegionsDTO struct {
	Objects []*flowdata.Object
}

func (node *InRegion) Init(storage flownodes.Storage) {
	node.storage = storage
}

func (node *InRegion) Run(req RegionsDTO, c flownodes.Ctx) (resp RegionsDTO, err error) {
	regions, err := node.getRegions(c.GetSourceInfo().ID)
	if err != nil {
		return resp, err
	}

	// map to collect only unique objects
	objects := map[*flowdata.Object]struct{}{}
	for _, r := range regions {
		for _, o := range req.Objects {
			if r.ContainsPoint(o.BBox.Center()) {
				objects[o] = struct{}{}
			}
		}
	}

	for o := range objects {
		resp.Objects = append(resp.Objects, o)
	}

	return resp, nil
}

func (node *InRegion) getRegions(cameraID string) (regions []flowdata.Region, err error) {
	// load camera from the storage
	cam, err := node.storage.GetCamera(cameraID)
	if err != nil {
		return regions, err
	}

	return cam.GetRegions(node.Region), nil
}

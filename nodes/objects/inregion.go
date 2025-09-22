package objects

import (
	"github.com/osmoview/flowdata"
	"github.com/osmoview/flownodes"
)

type InRegion struct {
	Regions []string `yaml:"regions"`

	storage flownodes.Storage
}

type InRegionInput struct {
	Objects []*flowdata.Object
}

type InRegionOutput struct {
	Objects []*flowdata.Object
	Regions []flowdata.Region
}

func (node *InRegion) Init(storage flownodes.Storage) {
	node.storage = storage
}

func (node *InRegion) Run(in InRegionInput, c flownodes.Ctx) (out InRegionOutput, err error) {
	regions, err := node.getRegions(c.GetSourceInfo().ID)
	if err != nil {
		return out, err
	}

	// map to collect only unique objects
	objectsMap := map[*flowdata.Object]struct{}{}
	regionsMap := map[string]flowdata.Region{}
	for _, r := range regions {
		for _, o := range in.Objects {
			if r.ContainsPoint(o.BBox.Center()) {
				objectsMap[o] = struct{}{}
				regionsMap[r.ID] = r
			}
		}
	}

	for o := range objectsMap {
		out.Objects = append(out.Objects, o)
	}

	for _, r := range regionsMap {
		out.Regions = append(out.Regions, r)
	}

	return out, nil
}

func (node *InRegion) getRegions(cameraID string) (regions []flowdata.Region, err error) {
	// load camera from the storage
	ds, err := node.storage.GetDataSource(cameraID)
	if err != nil {
		return regions, err
	}

	for _, regionName := range node.Regions {
		regions = append(regions, ds.GetRegions(regionName)...)
	}

	return
}

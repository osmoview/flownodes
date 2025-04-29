package objects

import (
	"github.com/osmoview/flowdata"
	"github.com/osmoview/flownodes"
)

//
// WARN: maybe need a standalone service to search profiles
//

type Profiles struct {
	ProfilesType string  `yaml:"profiles_type"`
	MinThreshold float32 `yaml:"min_threshold"`
	MaxThreshold float32 `yaml:"max_threshold"`

	db flownodes.Database
}

func (node *Profiles) Init(data flownodes.Data, db flownodes.Database) {
	node.db = db

	// set default thresholds
	if node.MinThreshold == 0 {
		node.MinThreshold = 0.6
	}
	if node.MaxThreshold == 0 {
		node.MaxThreshold = 1.0
	}
}

type ProfilesDTO struct {
	Objects []*flowdata.Object
}

func (node *Profiles) Run(req ProfilesDTO) (ProfilesDTO, error) {
	// no objects, no need to touch the database
	if len(req.Objects) == 0 {
		return req, nil
	}

	// search profiles by external id and select objects with vector to
	// the objects map
	objectsWithVector := map[*flowdata.Object]struct{}{}
	for _, o := range req.Objects {
		// ignore objects with profiles
		if o.Profile != nil {
			continue
		}

		// try to fine a profile by an external_id
		if o.ExternalID != "" {
			x := flowdata.Profile{Type: node.ProfilesType, ExternalID: o.ExternalID}
			if err := node.db.Get(&x); err == nil {
				o.Profile = &x
				continue
			}
		}

		// if object has vector then append it to the objects map to search profiles by vector
		if o.Vector.IsValid() {
			objectsWithVector[o] = struct{}{}
		}
	}

	// load profiles from the database one by one and compare it with objects
	err := node.db.Range(&flowdata.Profile{Type: node.ProfilesType}, func(v flownodes.DatabaseTableModel) bool {
		p := v.(*flowdata.Profile)

		for o := range objectsWithVector {
			dist, iter, ok := p.VectorMatch(o.Vector, node.MaxThreshold, node.MinThreshold)
			if ok {
				o.Profile = p
				o.VectorMatch = []float32{dist, float32(iter)}
				o.Name = p.Name
				delete(objectsWithVector, o)
				break
			}
		}

		return len(objectsWithVector) > 0
	})

	return req, err
}

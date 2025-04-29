package objects

import (
	"fmt"
	"sort"
	"strconv"
	"time"

	"github.com/osmoview/flowdata"
	"github.com/osmoview/flownodes"
)

type Tracking struct {
	Dist float32
	Keep float32

	id       int
	prevTime time.Time
	prevList []*trackingObject
}

type TrackingData struct {
	Objects []*flowdata.Object
}

type trackingObject struct {
	*flowdata.Object

	used bool
	time time.Time
	bbox flowdata.BBox
}

type matchObject struct {
	*trackingObject
	dist float32
}

func (node *Tracking) Init() {
	if node.Keep <= 0 {
		node.Keep = 3
	}
}

func (node *Tracking) Run(req TrackingData, ctx flownodes.Ctx) TrackingData {
	currTime := ctx.GetTime()
	currList := node.prepareCurrentObjects(req.Objects, currTime)

	timeFactor := float32(currTime.Sub(node.prevTime).Seconds())

	cleanTime := currTime.Add(-1 * time.Duration(timeFactor*node.Keep*float32(time.Second)))
	node.clean(cleanTime)

	if len(node.prevList) > 0 {
		prevUsed := 0

		// infinite loop protection
		for i := 0; i < 10; i++ {
			used := 0

			for _, item := range currList {
				if item.used {
					used++
					continue
				}

				prevMatches := node.orderByDistance(item, node.prevList, timeFactor)
				if len(prevMatches) == 0 {
					continue
				}
				prev := prevMatches[0]

				currMatches := node.orderByDistance(prev.trackingObject, currList, timeFactor)
				if len(currMatches) == 0 {
					continue
				}
				curr := currMatches[0]

				if curr.Object == item.Object {
					item.ID = prev.ID
					item.used = true
					prev.used = true
				}
			}

			if used >= len(currList) || used == prevUsed {
				break
			}
			prevUsed = used
		}

	}

	var prevList []*trackingObject
	for _, prev := range node.prevList {
		if !prev.used {
			prevList = append(prevList, prev)
		}
	}
	node.prevList = prevList

	for _, obj := range currList {
		if !node.isID(obj.ID) {
			obj.ID = node.nextID()
		}

		if len(obj.Objects) > 0 {
			node.setNestedObjectsID(obj.Object)
		}

		obj.used = false
		node.prevList = append(node.prevList, obj)
	}

	node.prevTime = currTime
	return req
}

func (node *Tracking) prepareCurrentObjects(list []*flowdata.Object, t time.Time) (currList []*trackingObject) {
	for _, o := range list {
		if !o.BBox.IsValid() {
			continue
		}
		to := &trackingObject{
			Object: o,
			time:   t,
			bbox:   o.BBox,
		}

		currList = append(currList, to)
	}

	return
}

func (node *Tracking) orderByDistance(obj *trackingObject, list []*trackingObject, timeFactor float32) (matches []matchObject) {
	size := (obj.BBox[2] + obj.BBox[3]) / 2
	maxdist := size / node.Dist * timeFactor * 16

	for _, to := range list {
		if to.used || to.Type != obj.Type || !to.BBox.IsValid() {
			continue
		}

		dist := to.BBox.Dist(obj.BBox)

		if dist < maxdist {
			matches = append(matches, matchObject{to, dist})
		}
	}

	sort.Slice(matches, func(i, j int) bool {
		return matches[i].dist < matches[j].dist
	})

	return
}

func (node *Tracking) setNestedObjectsID(obj *flowdata.Object) {
	for i, no := range obj.Objects {
		no.ID = fmt.Sprintf("%s-%s:%d", obj.ID, no.Type, i)

		if len(no.Objects) > 0 {
			node.setNestedObjectsID(no)
		}
	}
}

//
// utils
//

func (node *Tracking) clean(t time.Time) {
	var list []*trackingObject

	for _, o := range node.prevList {
		if o.time.After(t) {
			list = append(list, o)
		}
	}

	node.prevList = list
}

func (node *Tracking) nextID() string {
	node.id++
	return strconv.Itoa(node.id)
}

func (node *Tracking) isID(id string) bool {
	_, err := strconv.Atoi(id)
	return err == nil
}

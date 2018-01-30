package main

import (
	"sort"
	"time"
)

type osm struct {
	Bounds *struct {
		Minlat string `xml:"minlat,attr"`
		Minlon string `xml:"minlon,attr"`
		Maxlat string `xml:"maxlat,attr"`
		Maxlon string `xml:"maxlon,attr"`
	} `xml:"bounds"`

	NodeList []struct {
		ID        int       `xml:"id,attr"`
		UID       int       `xml:"uid,attr"`
		Changeset string    `xml:"changeset,attr"`
		Timestamp time.Time `xml:"timestamp,attr"`
		Version   int       `xml:"version,attr"`
		Visible   bool      `xml:"visible,attr"`
		User      string    `xml:"user,attr"`
		Lat       float64   `xml:"lat,attr"`
		Lon       float64   `xml:"lon,attr"`

		Tag []struct {
			K string `xml:"k,attr"`
			V string `xml:"v,attr"`
		} `xml:"tag"`
	} `xml:"node"`

	WayList []struct {
		ID        int       `xml:"id,attr"`
		UID       int       `xml:"uid,attr"`
		Changeset string    `xml:"changeset,attr"`
		Timestamp time.Time `xml:"timestamp,attr"`
		Version   int       `xml:"version,attr"`
		Visible   bool      `xml:"visible,attr"`
		User      string    `xml:"user,attr"`
		Lat       float64   `xml:"lat,attr"`
		Lon       float64   `xml:"lon,attr"`

		Nd []struct {
			Ref int `xml:"ref,attr"`
		} `xml:"nd"`

		Tag []struct {
			K string `xml:"k,attr"`
			V string `xml:"v,attr"`
		} `xml:"tag"`
	} `xml:"way"`

	RelationList []struct {
		ID        int       `xml:"id,attr"`
		UID       int       `xml:"uid,attr"`
		Changeset string    `xml:"changeset,attr"`
		Timestamp time.Time `xml:"timestamp,attr"`
		Version   int       `xml:"version,attr"`
		Visible   bool      `xml:"visible,attr"`
		User      string    `xml:"user,attr"`

		Member []struct {
			Type string `xml:"type,attr"`
			Ref  int    `xml:"ref,attr"`
			Role string `xml:"role,attr"`
		} `xml:"member"`

		Tag []struct {
			K string `xml:"k,attr"`
			V string `xml:"v,attr"`
		} `xml:"tag"`
	} `xml:"relation"`
}

func (o *osm) sort() {
	o.Bounds = nil // temporary solution

	for k := range o.NodeList {
		sort.Slice(o.NodeList[k].Tag, func(i, j int) bool {
			if o.NodeList[k].Tag[i].K < o.NodeList[k].Tag[j].K {
				return true
			}
			if o.NodeList[k].Tag[i].K > o.NodeList[k].Tag[j].K {
				return false
			}
			return o.NodeList[k].Tag[i].V < o.NodeList[k].Tag[j].V
		})
	}

	for k := range o.WayList {
		sort.Slice(o.WayList[k].Nd, func(i, j int) bool {
			return o.WayList[k].Nd[i].Ref < o.WayList[k].Nd[j].Ref
		})
		sort.Slice(o.WayList[k].Tag, func(i, j int) bool {
			if o.WayList[k].Tag[i].K < o.WayList[k].Tag[j].K {
				return true
			}
			if o.WayList[k].Tag[i].K > o.WayList[k].Tag[j].K {
				return false
			}
			return o.WayList[k].Tag[i].V < o.WayList[k].Tag[j].V
		})
	}

	for k := range o.RelationList {
		sort.Slice(o.RelationList[k].Member, func(i, j int) bool {
			if o.RelationList[k].Member[i].Ref < o.RelationList[k].Member[j].Ref {
				return true
			}
			if o.RelationList[k].Member[i].Ref > o.RelationList[k].Member[j].Ref {
				return false
			}
			if o.RelationList[k].Member[i].Type < o.RelationList[k].Member[j].Type {
				return true
			}
			if o.RelationList[k].Member[i].Type > o.RelationList[k].Member[j].Type {
				return false
			}
			return o.RelationList[k].Member[i].Role > o.RelationList[k].Member[j].Role
		})
		sort.Slice(o.RelationList[k].Tag, func(i, j int) bool {
			if o.RelationList[k].Tag[i].K < o.RelationList[k].Tag[j].K {
				return true
			}
			if o.RelationList[k].Tag[i].K > o.RelationList[k].Tag[j].K {
				return false
			}
			return o.RelationList[k].Tag[i].V < o.RelationList[k].Tag[j].V
		})
	}

	sort.Slice(o.NodeList, func(i, j int) bool { return o.NodeList[i].ID < o.NodeList[j].ID })
	sort.Slice(o.WayList, func(i, j int) bool { return o.WayList[i].ID < o.WayList[j].ID })
	sort.Slice(o.RelationList, func(i, j int) bool { return o.RelationList[i].ID < o.RelationList[j].ID })
}

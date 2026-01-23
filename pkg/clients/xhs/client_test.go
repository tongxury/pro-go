package xhs

import (
	"context"
	"testing"
)

func TestName(t *testing.T) {

	c := NewClient()
	//{"keyword":"","page":1,"page_size":20,"search_id":"2eplbdefxvxnbifcux7lw@2eplbeb8ai21fzhmv8flu","sort":"popularity_descending","note_type":1,"ext_flags":[],
	//"filters":[{"tags":["general"],"type":"sort_type"},{"tags":["不限"],"type":"filter_note_type"},{"tags":["不限"],"type":"filter_note_time"},{"tags":["不限"],"type":"filter_note_range"},{"tags":["不限"],"type":"filter_pos_distance"}],"geo":"","image_formats":["jpg","webp","avif"]}
	c.ListNotes(context.Background(), ListNotesParams{
		Keyword:  "热点",
		Page:     1,
		PageSize: 20,
		SearchId: "2eplbdefxvxnbifcux7lw@2eplbeb8ai21fzhmv8flu",
		Sort:     "popularity_descending",
		NoteType: 1,
		ExtFlags: nil,
		Filters: []Filter{
			{
				Tags: []string{"general"},
				Type: "sort_type",
			},
			{
				Tags: []string{"不限"},
				Type: "filter_note_type",
			},
			{
				Tags: []string{"不限"},
				Type: "filter_note_time",
			},
			{
				Tags: []string{"不限"},
				Type: "filter_note_range",
			}, {
				Tags: []string{"不限"},
				Type: "filter_pos_distance",
			},
		},
		Geo:          "",
		ImageFormats: []string{"jpg", "webp", "avif"},
	})

}

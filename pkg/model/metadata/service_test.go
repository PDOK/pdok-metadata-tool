package metadata

import (
	"testing"
)

func TestOperatesOnRef_GetID1(t *testing.T) {
	type fields struct {
		UUIDRef string
		Href    string
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			name: "Test get id from escaped Href",
			fields: fields{
				UUIDRef: "",
				Href:    "https://www.nationaalgeoregister.nl/geonetwork/srv/dut/csw?service=CSW&amp;request=GetRecordById&amp;version=2.0.2&amp;outputSchema=http://www.isotc211.org/2005/gmd&amp;elementSetName=full&amp;id=290d86cf-9791-4c1a-8bf2-c1520cf3563a#MD_DataIdentification",
			},
			want: "290d86cf-9791-4c1a-8bf2-c1520cf3563a",
		},
		{
			name: "Test get ID from escaped Href",
			fields: fields{
				UUIDRef: "",
				Href:    "https://nationaalgeoregister.nl/geonetwork/srv/dut/csw?SERVICE=CSW&amp;version=2.0.2&amp;REQUEST=GetRecordById&amp;ID=a29917b9-3426-4041-a11b-69bcb2256904&amp;OUTPUTSCHEMA=http://www.isotc211.org/2005/gmd&amp;ELEMENTSETNAME=full#MD_DataIdentification",
			},
			want: "a29917b9-3426-4041-a11b-69bcb2256904",
		},
		{
			name: "Test get id from escaped Href with whitespace",
			fields: fields{
				UUIDRef: "",
				Href:    "https://nationaalgeoregister.nl/geonetwork/srv/dut/csw?service=CSW&amp;request=GetRecordById&amp;version=2.0.2&amp;outputSchema=http://www.isotc211.org/2005/gmd&amp;elementSetName=full&amp;id=a29917b9-3426-4041-a11b-69bcb2256904 #MD_DataIdentification",
			},
			want: "a29917b9-3426-4041-a11b-69bcb2256904",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			o := &OperatesOnRef{
				UUIDRef: tt.fields.UUIDRef,
				Href:    tt.fields.Href,
			}
			if got := o.GetID(); got != tt.want {
				t.Errorf("GetID() = %v, want %v", got, tt.want)
			}
		})
	}
}

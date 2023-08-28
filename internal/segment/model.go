package segment

type Segment struct {
	SegmentId   int    `json:"segment_id"`
	SegmentName string `json:"segment_name"`
	Active      bool   `json:"active"`
}

type SegmentDTO struct {
	SegmentName string `json:"segment_name"`
}

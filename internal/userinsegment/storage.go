package userinsegment

import (
	"AvitoTechTask/internal/segment"
	"context"
)

type Repository interface {
	AddSegments(ctx context.Context, user *UserDTO, segment *segment.SegmentDTO) error
	AddSegmentsWithPeriod(ctx context.Context, user *UserDTO, segment *segment.SegmentDTO) error
	DeleteFromSegment(ctx context.Context, user *UserInSegmentDTO) error
	GetSegments(ctx context.Context, user *UserInSegment) (UserSegmentsListDTO, error)
	GetUserHistory(ctx context.Context, user *UserInSegment) ([]UserInSegmentsHistory, error)
}

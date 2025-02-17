// Copyright (C) 2020 Storj Labs, Inc.
// See LICENSE for copying information.

package metabase_test

import (
	"testing"
	"time"

	"storj.io/common/storj"
	"storj.io/common/testcontext"
	"storj.io/common/testrand"
	"storj.io/storj/satellite/metabase"
	"storj.io/storj/satellite/metabase/metabasetest"
)

func TestDeleteExpiredObjects(t *testing.T) {
	metabasetest.Run(t, func(ctx *testcontext.Context, t *testing.T, db *metabase.DB) {
		obj1 := metabasetest.RandObjectStream()
		obj2 := metabasetest.RandObjectStream()
		obj3 := metabasetest.RandObjectStream()

		now := time.Now()
		zombieDeadline := now.Add(24 * time.Hour)
		pastTime := now.Add(-1 * time.Hour)
		futureTime := now.Add(1 * time.Hour)

		t.Run("none", func(t *testing.T) {
			defer metabasetest.DeleteAll{}.Check(ctx, t, db)

			metabasetest.DeleteExpiredObjects{
				Opts: metabase.DeleteExpiredObjects{
					ExpiredBefore: time.Now(),
				},
			}.Check(ctx, t, db)
			metabasetest.Verify{}.Check(ctx, t, db)
		})

		t.Run("partial objects", func(t *testing.T) {
			defer metabasetest.DeleteAll{}.Check(ctx, t, db)

			// pending object without expiration time
			metabasetest.BeginObjectExactVersion{
				Opts: metabase.BeginObjectExactVersion{
					ObjectStream: obj1,
					Encryption:   metabasetest.DefaultEncryption,
				},
				Version: 1,
			}.Check(ctx, t, db)

			// pending object with expiration time in the past
			metabasetest.BeginObjectExactVersion{
				Opts: metabase.BeginObjectExactVersion{
					ObjectStream: obj2,
					ExpiresAt:    &pastTime,
					Encryption:   metabasetest.DefaultEncryption,
				},
				Version: 1,
			}.Check(ctx, t, db)

			// pending object with expiration time in the future
			metabasetest.BeginObjectExactVersion{
				Opts: metabase.BeginObjectExactVersion{
					ObjectStream: obj3,
					ExpiresAt:    &futureTime,
					Encryption:   metabasetest.DefaultEncryption,
				},
				Version: 1,
			}.Check(ctx, t, db)

			metabasetest.DeleteExpiredObjects{
				Opts: metabase.DeleteExpiredObjects{
					ExpiredBefore: time.Now(),
				},
			}.Check(ctx, t, db)

			metabasetest.Verify{ // the object with expiration time in the past is gone
				Objects: []metabase.RawObject{
					{
						ObjectStream: obj1,
						CreatedAt:    now,
						Status:       metabase.Pending,

						Encryption:             metabasetest.DefaultEncryption,
						ZombieDeletionDeadline: &zombieDeadline,
					},
					{
						ObjectStream: obj3,
						CreatedAt:    now,
						ExpiresAt:    &futureTime,
						Status:       metabase.Pending,

						Encryption:             metabasetest.DefaultEncryption,
						ZombieDeletionDeadline: &zombieDeadline,
					},
				},
			}.Check(ctx, t, db)
		})

		t.Run("batch size", func(t *testing.T) {
			expiresAt := time.Now().Add(-30 * 24 * time.Hour)
			for i := 0; i < 32; i++ {
				_ = metabasetest.CreateExpiredObject(ctx, t, db, metabasetest.RandObjectStream(), 3, expiresAt)
			}
			metabasetest.DeleteExpiredObjects{
				Opts: metabase.DeleteExpiredObjects{
					ExpiredBefore: time.Now().Add(time.Hour),
					BatchSize:     4,
				},
			}.Check(ctx, t, db)

			metabasetest.Verify{}.Check(ctx, t, db)
		})

		t.Run("committed objects", func(t *testing.T) {
			defer metabasetest.DeleteAll{}.Check(ctx, t, db)

			object1 := metabasetest.CreateTestObject{}.Run(ctx, t, db, obj1, 1)
			metabasetest.CreateTestObject{
				BeginObjectExactVersion: &metabase.BeginObjectExactVersion{
					ObjectStream: obj2,
					ExpiresAt:    &pastTime,
					Encryption:   metabasetest.DefaultEncryption,
				},
			}.Run(ctx, t, db, obj2, 1)
			object3 := metabasetest.CreateTestObject{
				BeginObjectExactVersion: &metabase.BeginObjectExactVersion{
					ObjectStream: obj3,
					ExpiresAt:    &futureTime,
					Encryption:   metabasetest.DefaultEncryption,
				},
			}.Run(ctx, t, db, obj3, 1)

			expectedObj1Segment := metabase.Segment{
				StreamID:          obj1.StreamID,
				RootPieceID:       storj.PieceID{1},
				CreatedAt:         now,
				EncryptedKey:      []byte{3},
				EncryptedKeyNonce: []byte{4},
				EncryptedETag:     []byte{5},
				EncryptedSize:     1060,
				PlainSize:         512,
				Pieces:            metabase.Pieces{{Number: 0, StorageNode: storj.NodeID{2}}},
				Redundancy:        metabasetest.DefaultRedundancy,
			}

			expectedObj3Segment := expectedObj1Segment
			expectedObj3Segment.StreamID = obj3.StreamID
			expectedObj3Segment.ExpiresAt = &futureTime

			metabasetest.DeleteExpiredObjects{
				Opts: metabase.DeleteExpiredObjects{
					ExpiredBefore: time.Now(),
				},
			}.Check(ctx, t, db)

			metabasetest.Verify{ // the object with expiration time in the past is gone
				Objects: []metabase.RawObject{
					metabase.RawObject(object1),
					metabase.RawObject(object3),
				},
				Segments: []metabase.RawSegment{
					metabase.RawSegment(expectedObj1Segment),
					metabase.RawSegment(expectedObj3Segment),
				},
			}.Check(ctx, t, db)
		})
	})
}

func TestDeleteZombieObjects(t *testing.T) {
	metabasetest.Run(t, func(ctx *testcontext.Context, t *testing.T, db *metabase.DB) {
		obj1 := metabasetest.RandObjectStream()
		obj2 := metabasetest.RandObjectStream()
		obj3 := metabasetest.RandObjectStream()

		now := time.Now()
		zombieDeadline := now.Add(24 * time.Hour)
		pastTime := now.Add(-1 * time.Hour)
		futureTime := now.Add(1 * time.Hour)

		t.Run("none", func(t *testing.T) {
			defer metabasetest.DeleteAll{}.Check(ctx, t, db)

			metabasetest.DeleteZombieObjects{
				Opts: metabase.DeleteZombieObjects{
					DeadlineBefore: time.Now(),
				},
			}.Check(ctx, t, db)
			metabasetest.Verify{}.Check(ctx, t, db)
		})

		t.Run("partial objects", func(t *testing.T) {
			defer metabasetest.DeleteAll{}.Check(ctx, t, db)

			// zombie object with default deadline
			metabasetest.BeginObjectExactVersion{
				Opts: metabase.BeginObjectExactVersion{
					ObjectStream: obj1,
					Encryption:   metabasetest.DefaultEncryption,
				},
				Version: 1,
			}.Check(ctx, t, db)

			// zombie object with deadline time in the past
			metabasetest.BeginObjectExactVersion{
				Opts: metabase.BeginObjectExactVersion{
					ObjectStream:           obj2,
					ZombieDeletionDeadline: &pastTime,
					Encryption:             metabasetest.DefaultEncryption,
				},
				Version: 1,
			}.Check(ctx, t, db)

			// pending object with expiration time in the future
			metabasetest.BeginObjectExactVersion{
				Opts: metabase.BeginObjectExactVersion{
					ObjectStream:           obj3,
					ZombieDeletionDeadline: &futureTime,
					Encryption:             metabasetest.DefaultEncryption,
				},
				Version: 1,
			}.Check(ctx, t, db)

			metabasetest.DeleteZombieObjects{
				Opts: metabase.DeleteZombieObjects{
					DeadlineBefore: time.Now(),
				},
			}.Check(ctx, t, db)

			metabasetest.Verify{ // the object with zombie deadline time in the past is gone
				Objects: []metabase.RawObject{
					{
						ObjectStream: obj1,
						CreatedAt:    now,
						Status:       metabase.Pending,

						Encryption:             metabasetest.DefaultEncryption,
						ZombieDeletionDeadline: &zombieDeadline,
					},
					{
						ObjectStream: obj3,
						CreatedAt:    now,
						Status:       metabase.Pending,

						Encryption:             metabasetest.DefaultEncryption,
						ZombieDeletionDeadline: &futureTime,
					},
				},
			}.Check(ctx, t, db)
		})

		t.Run("batch size", func(t *testing.T) {
			deadlineAt := time.Now().Add(-30 * 24 * time.Hour)
			for i := 0; i < 33; i++ {
				obj := metabasetest.RandObjectStream()

				metabasetest.BeginObjectExactVersion{
					Opts: metabase.BeginObjectExactVersion{
						ObjectStream:           obj,
						Encryption:             metabasetest.DefaultEncryption,
						ZombieDeletionDeadline: &deadlineAt,
					},
					Version: obj.Version,
				}.Check(ctx, t, db)

				for i := byte(0); i < 3; i++ {
					metabasetest.BeginSegment{
						Opts: metabase.BeginSegment{
							ObjectStream: obj,
							Position:     metabase.SegmentPosition{Part: 0, Index: uint32(i)},
							RootPieceID:  storj.PieceID{i + 1},
							Pieces: []metabase.Piece{{
								Number:      1,
								StorageNode: testrand.NodeID(),
							}},
						},
					}.Check(ctx, t, db)

					metabasetest.CommitSegment{
						Opts: metabase.CommitSegment{
							ObjectStream: obj,
							Position:     metabase.SegmentPosition{Part: 0, Index: uint32(i)},
							RootPieceID:  storj.PieceID{1},
							Pieces:       metabase.Pieces{{Number: 0, StorageNode: storj.NodeID{2}}},

							EncryptedKey:      []byte{3},
							EncryptedKeyNonce: []byte{4},
							EncryptedETag:     []byte{5},

							EncryptedSize: 1024,
							PlainSize:     512,
							PlainOffset:   0,
							Redundancy:    metabasetest.DefaultRedundancy,
						},
					}.Check(ctx, t, db)
				}
			}

			metabasetest.DeleteZombieObjects{
				Opts: metabase.DeleteZombieObjects{
					DeadlineBefore: time.Now().Add(time.Hour),
					BatchSize:      4,
				},
			}.Check(ctx, t, db)

			metabasetest.Verify{}.Check(ctx, t, db)
		})

		t.Run("committed objects", func(t *testing.T) {
			defer metabasetest.DeleteAll{}.Check(ctx, t, db)

			object1 := metabasetest.CreateTestObject{}.Run(ctx, t, db, obj1, 1)

			object2 := object1
			object2.ObjectStream = obj2
			metabasetest.CreateTestObject{
				BeginObjectExactVersion: &metabase.BeginObjectExactVersion{
					ObjectStream:           object2.ObjectStream,
					ZombieDeletionDeadline: &pastTime,
					Encryption:             metabasetest.DefaultEncryption,
				},
			}.Run(ctx, t, db, object2.ObjectStream, 1)

			object3 := metabasetest.CreateTestObject{
				BeginObjectExactVersion: &metabase.BeginObjectExactVersion{
					ObjectStream:           obj3,
					ZombieDeletionDeadline: &futureTime,
					Encryption:             metabasetest.DefaultEncryption,
				},
			}.Run(ctx, t, db, obj3, 1)

			expectedObj1Segment := metabase.Segment{
				StreamID:          obj1.StreamID,
				RootPieceID:       storj.PieceID{1},
				CreatedAt:         now,
				EncryptedKey:      []byte{3},
				EncryptedKeyNonce: []byte{4},
				EncryptedETag:     []byte{5},
				EncryptedSize:     1060,
				PlainSize:         512,
				Pieces:            metabase.Pieces{{Number: 0, StorageNode: storj.NodeID{2}}},
				Redundancy:        metabasetest.DefaultRedundancy,
			}

			expectedObj2Segment := expectedObj1Segment
			expectedObj2Segment.StreamID = object2.StreamID
			expectedObj3Segment := expectedObj1Segment
			expectedObj3Segment.StreamID = object3.StreamID

			metabasetest.DeleteZombieObjects{
				Opts: metabase.DeleteZombieObjects{
					DeadlineBefore: time.Now(),
				},
			}.Check(ctx, t, db)

			metabasetest.Verify{ // all committed objects should NOT be deleted
				Objects: []metabase.RawObject{
					metabase.RawObject(object1),
					metabase.RawObject(object2),
					metabase.RawObject(object3),
				},
				Segments: []metabase.RawSegment{
					metabase.RawSegment(expectedObj1Segment),
					metabase.RawSegment(expectedObj2Segment),
					metabase.RawSegment(expectedObj3Segment),
				},
			}.Check(ctx, t, db)
		})
	})
}

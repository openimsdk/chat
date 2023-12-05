package chat

import (
	"context"
	table "github.com/OpenIMSDK/chat/pkg/common/db/table/chat"
	"github.com/OpenIMSDK/chat/pkg/proto/emoticon_pack"
	"github.com/OpenIMSDK/tools/errs"
	"github.com/OpenIMSDK/tools/log"
	"sync"
	"time"
)

const (
	epoch            int64 = 1609459200000 // 设置起始时间 (例如 2021-01-01)
	machineIDBits    uint8 = 5             // 机器ID所占的位数
	datacenterIDBits uint8 = 5             // 数据中心ID所占的位数
	sequenceBits     uint8 = 12            // 序列号所占的位数

	maxMachineID    int64 = -1 ^ (-1 << machineIDBits)    // 最大机器ID
	maxDatacenterID int64 = -1 ^ (-1 << datacenterIDBits) // 最大数据中心ID
	maxSequence     int64 = -1 ^ (-1 << sequenceBits)     // 最大序列号

	timeLeft    uint8 = 22 // 时间戳向左的位移
	dataLeft    uint8 = 17 // 数据中心ID向左的位移
	machineLeft uint8 = 12 // 机器ID向左的位移
)

type Snowflake struct {
	mutex         sync.Mutex // 保护同时访问
	lastTimestamp int64
	datacenterID  int64
	machineID     int64
	sequence      int64
}

func NewSnowflake(datacenterID, machineID int64) (*Snowflake, error) {
	if datacenterID < 0 || datacenterID > maxDatacenterID {
		return nil, errs.ErrData.Wrap("datacenter ID out of rang")
	}
	if machineID < 0 || machineID > maxMachineID {
		return nil, errs.ErrData.Wrap("machine ID out of range")
	}

	return &Snowflake{
		lastTimestamp: 0,
		datacenterID:  datacenterID,
		machineID:     machineID,
		sequence:      0,
	}, nil
}

func (s *Snowflake) Generate() (int64, error) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	now := time.Now().UnixNano() / 1e6 // 当前时间的毫秒时间戳
	if s.lastTimestamp == now {
		s.sequence = (s.sequence + 1) & maxSequence
		if s.sequence == 0 {
			for now <= s.lastTimestamp {
				now = time.Now().UnixNano() / 1e6
			}
		}
	} else {
		s.sequence = 0
	}

	if now < s.lastTimestamp {
		return 0, errs.ErrData.Wrap("clock moved backwards")
	}
	s.lastTimestamp = now

	id := ((now - epoch) << timeLeft) | (s.datacenterID << dataLeft) | (s.machineID << machineLeft) | s.sequence
	return id, nil
}

func (o *chatSvr) AddEmoticon(ctx context.Context, req *emoticon_pack.AddEmoticonReq) (*emoticon_pack.AddEmoticonResp, error) {

	log.ZDebug(ctx, "hello here rpc", "add Emoticon")
	sf, err := NewSnowflake(1, 1)
	if err != nil {
		return nil, err
	}
	result, err := sf.Generate()
	if err != nil {
		return nil, err
	}
	image := &table.Image{
		ID:       result,
		ImageURL: req.ImageData,
		OwnerID:  req.OwnerId,
	}
	err = o.Database.AddImage(ctx, image)
	if err != nil {
		return nil, err
	}

	return &emoticon_pack.AddEmoticonResp{}, nil
}
func (o *chatSvr) RemoveEmoticon(ctx context.Context, req *emoticon_pack.RemoveEmoticonReq) (*emoticon_pack.RemoveEmoticonResp, error) {
	//userID, _, err := mctx.Check(ctx)
	//if _, err := o.Database.GetUser(ctx, userID); err != nil {
	//	return nil, err
	//}

	err := o.Database.RemoveImage(ctx, req.UserId, req.EmoticonId)
	if err != nil {
		return nil, err
	}

	return &emoticon_pack.RemoveEmoticonResp{}, nil
}
func (o *chatSvr) GetEmoticon(ctx context.Context, req *emoticon_pack.GetEmoticonReq) (*emoticon_pack.GetEmoticonResp, error) {
	results, err := o.Database.GetImages(ctx, req.UserId)
	if err != nil {
		return nil, err
	}

	var pbEmoticons []*emoticon_pack.Emoticon
	for _, result := range results {
		pbEmoticons = append(pbEmoticons, &emoticon_pack.Emoticon{
			ImageURL:   result.ImageURL,
			EmoticonId: result.ID,
			UserId:     result.OwnerID,
		})
	}

	return &emoticon_pack.GetEmoticonResp{E: pbEmoticons}, nil
}
